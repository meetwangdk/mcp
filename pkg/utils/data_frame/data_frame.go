package dataframe

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

// GrafanaDataFrame is a struct that represents a Grafana data frame
// generally, we can use https://pkg.go.dev/github.com/grafana/grafana-plugin-sdk-go/data,
// but there are some compatibility issues with the lunettes api response.
type GrafanaDataFrame struct {
	Schema struct {
		Fields []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"fields"`
		Name string `json:"name"`
	} `json:"schema"`
	Data struct {
		Values [][]any `json:"values"`
	} `json:"data"`
}

// TableData is a struct that represents a table
type TableData struct {
	Name string  `json:"name"`
	Data [][]any `json:"Data"`
}

// ToTable converts Grafana data frames to a tables
// if data is a single frame, it will be converted to a table
func ToTable(data string, opts ...Option) string {
	if len(data) > 0 && data[0] == '{' {
		return FrameToTable(data, opts...)
	}

	frames := make([]GrafanaDataFrame, 0)
	err := json.Unmarshal([]byte(data), &frames)
	if err != nil {
		return data
	}
	if len(frames) == 0 {
		return data
	}

	opt := &options{}
	for _, o := range opts {
		o.apply(opt)
	}

	framesTable := make([]TableData, 0)
	for _, frame := range frames {
		table := dataFrameToTable(frame, opt)
		if table == nil {
			return data
		}
		framesTable = append(framesTable, *table)
	}

	converted, _ := json.Marshal(framesTable)
	return string(converted)
}

// FrameToTable converts a Grafana data frame to a table
func FrameToTable(data string, opts ...Option) string {
	frame := GrafanaDataFrame{}
	err := json.Unmarshal([]byte(data), &frame)
	if err != nil {
		fmt.Println(err)
		return data
	}

	opt := &options{}
	for _, o := range opts {
		o.apply(opt)
	}
	table := dataFrameToTable(frame, opt)
	if table == nil {
		return data
	}

	converted, _ := json.Marshal(table)
	return string(converted)
}

func dataFrameToTable(frame GrafanaDataFrame, opts *options) *TableData {
	// if frame is empty, return nil
	if frame.Schema.Name == "" && len(frame.Schema.Fields) == 0 &&
		len(frame.Data.Values) == 0 {
		return nil
	}
	table := TableData{
		Name: frame.Schema.Name,
		Data: make([][]any, 0),
	}

	// add header
	row := make([]any, 0)
	includeFieldIndexs := make([]int, 0)
	for i, field := range frame.Schema.Fields {
		if len(opts.includeFields) > 0 && !slices.Contains(opts.includeFields, field.Name) {
			continue
		}
		if len(opts.excludeFields) > 0 && slices.Contains(opts.excludeFields, field.Name) {
			continue
		}
		includeFieldIndexs = append(includeFieldIndexs, i)
		row = append(row, field.Name)
	}
	table.Data = append(table.Data, row)

	// add data values
	fields := len(frame.Schema.Fields)
	for i := 0; i < len(frame.Data.Values[0]); i++ {
		row := make([]any, 0, fields)
		for j := 0; j < fields; j++ {
			if !slices.Contains(includeFieldIndexs, j) {
				continue
			}

			if i >= len(frame.Data.Values[j]) {
				row = append(row, nil)
				continue
			}
			v := frame.Data.Values[j][i]

			// format time
			if opts.formatTime && frame.Schema.Fields[j].Type == "time" {
				if tv, ok := v.(float64); ok {
					t := time.UnixMilli(int64(tv)).Format(time.RFC3339)
					v = t
				}
			}
			row = append(row, v)
		}
		table.Data = append(table.Data, row)
	}

	// rename fields
	if len(opts.renameFields) > 0 {
		for i, field := range table.Data[0] {
			if name, ok := opts.renameFields[field.(string)]; ok {
				table.Data[0][i] = name
			}
		}
	}

	// sort by field
	if opts.sortByField != "" {
		sortTable(table, opts.sortByField, opts.sortDesc)
	}
	return &table
}

func sortTable(table TableData, field string, desc bool) {
	// find field index
	fieldIndex := -1
	for i, f := range table.Data[0] {
		if f == field {
			fieldIndex = i
			break
		}
	}
	if fieldIndex == -1 {
		return
	}

	slices.SortStableFunc(table.Data[1:], func(a, b []any) int {
		if desc {
			return compare(b[fieldIndex], a[fieldIndex])
		}
		return compare(a[fieldIndex], b[fieldIndex])
	})
}

func compare(a, b any) int {
	if a == b {
		return 0
	}
	switch va := a.(type) {
	case float64:
		vb, ok := b.(float64)
		if ok {
			return int(va - vb)
		}
	case string:
		vb, ok := b.(string)
		if ok {
			return strings.Compare(va, vb)
		}
	}
	return -1
}
