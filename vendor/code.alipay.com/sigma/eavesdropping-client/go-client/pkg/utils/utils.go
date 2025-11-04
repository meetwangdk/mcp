package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"code.alipay.com/sigma/eavesdropping-client/go-client/utils"
	"k8s.io/klog/v2"
)

func SetStringParam(values url.Values, name string, f *string) {
	s := values.Get(name)
	if s != "" {
		*f = s
	}
}

func SetBoolParam(values url.Values, name string, f *bool) {
	s := values.Get(name)
	if s == "t" || s == "true" || s == "1" {
		*f = true
	}
}

func SetTimeParam(values url.Values, name string, f *time.Time) {
	s := values.Get(name)
	if s != "" {
		t, err := utils.ParseTime(s)
		if err != nil {
			klog.Errorf("failed to parse time %s for %s: %s", s, name, err.Error())
		} else {
			*f = t
		}
	}
}

func SetTimeLayoutParam(values url.Values, name string, f *time.Time) {
	s := values.Get(name)
	layOut := "2006-01-02T15:04:05"
	if s != "" {
		t, err := time.ParseInLocation(layOut, s, time.Local)
		if err != nil {
			klog.Errorf("failed to parse time %s for %s: %s", s, name, err.Error())
		} else {
			*f = t
		}
	}
}

// trace param parse

func ParseTags(simpleTags []string, jsonTags []string) (map[string]string, error) {
	retMe := make(map[string]string)
	for _, tag := range simpleTags {
		keyAndValue := strings.Split(tag, ":")
		if l := len(keyAndValue); l > 1 {
			retMe[keyAndValue[0]] = strings.Join(keyAndValue[1:], ":")
		} else {
			return nil, fmt.Errorf("malformed 'tag' parameter, expecting key:value, received: %s", tag)
		}
	}
	for _, tags := range jsonTags {
		var fromJSON map[string]string
		if err := json.Unmarshal([]byte(tags), &fromJSON); err != nil {
			return nil, fmt.Errorf("malformed 'tags' parameter, cannot unmarshal JSON: %s", err)
		}
		for k, v := range fromJSON {
			retMe[k] = v
		}
	}
	return retMe, nil
}

type Tabler interface {
	TableName() string
}

type EsTabler interface {
	EsTableName() string
}

type Typer interface {
	TypeName() string
}

// 获取数据的 SqlTableName EsTableName TypeName
func GetMetaName(dest interface{}) (string, string, string, error) {
	value := reflect.ValueOf(dest)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		value = reflect.New(value.Type().Elem())
	}
	modelType := reflect.Indirect(value).Type()

	if modelType.Kind() == reflect.Interface {
		modelType = reflect.Indirect(reflect.ValueOf(dest)).Elem().Type()
	}

	for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return "", "", "", errors.New("data must be struct or slice or array")
	}
	modelValue := reflect.New(modelType)
	tabler, ok := modelValue.Interface().(Tabler)
	if !ok {
		return "", "", "", errors.New("data doesn't have TableName")
	}
	esTabler, ok := modelValue.Interface().(EsTabler)
	if !ok {
		return "", "", "", errors.New("data doesn't have TableName")
	}
	typer, ok := modelValue.Interface().(Typer)
	if !ok {
		return "", "", "", errors.New("data doesn't have TableName")
	}
	return tabler.TableName(), esTabler.EsTableName(), typer.TypeName(), nil
}

type QueryBuilder struct {
	client *elastic.Client
}

func NewQueryBuilder(client *elastic.Client) *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildFromExpressions(expressions []string, baseQuery *elastic.BoolQuery) (*elastic.BoolQuery, error) {
	if baseQuery == nil {
		baseQuery = elastic.NewBoolQuery()
	}

	var err error
	currentQuery := baseQuery

	for _, expr := range expressions {
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}

		currentQuery, err = qb.BuildFromExpression(expr, currentQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expression '%s': %v", expr, err)
		}
	}

	return currentQuery, nil
}

func (qb *QueryBuilder) BuildFromExpression(expression string, baseQuery *elastic.BoolQuery) (*elastic.BoolQuery, error) {
	if baseQuery == nil {
		baseQuery = elastic.NewBoolQuery()
	}

	expression = qb.normalizeExpression(expression)
	if expression == "" {
		return baseQuery, nil
	}

	return qb.parseExpression(expression, baseQuery)
}

func (qb *QueryBuilder) normalizeExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return expr
	}

	spaceRegex := regexp.MustCompile(`\s+`)
	expr = spaceRegex.ReplaceAllString(expr, " ")

	colonRegex := regexp.MustCompile(`\s*:\s*`)
	expr = colonRegex.ReplaceAllString(expr, ":")

	expr = strings.ReplaceAll(expr, "( ", "(")
	expr = strings.ReplaceAll(expr, " )", ")")

	return expr
}

func (qb *QueryBuilder) parseExpression(expr string, currentQuery *elastic.BoolQuery) (*elastic.BoolQuery, error) {
	// 使用正则表达式匹配 AND NOT (...) 格式
	andNotParenthesesRegex := regexp.MustCompile(`^AND NOT \((.+)\)$`)
	if matches := andNotParenthesesRegex.FindStringSubmatch(expr); matches != nil {
		return qb.handleAndNotWithParentheses(matches[1], currentQuery)
	}

	// 处理简单的 AND/OR/AND NOT 表达式
	return qb.handleSimpleExpression(expr, currentQuery)
}

func (qb *QueryBuilder) handleAndNotWithParentheses(innerExpr string, currentQuery *elastic.BoolQuery) (*elastic.BoolQuery, error) {
	innerExpr = strings.TrimSpace(innerExpr)

	// 解析括号内的表达式
	subQuery, err := qb.parseInnerExpression(innerExpr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inner expression '%s': %v", innerExpr, err)
	}

	return currentQuery.MustNot(subQuery), nil
}

func (qb *QueryBuilder) parseInnerExpression(innerExpr string) (elastic.Query, error) {
	// 检查是否包含 OR
	orRegex := regexp.MustCompile(` OR `)
	if orRegex.MatchString(innerExpr) {
		parts := orRegex.Split(innerExpr, -1)
		boolQuery := elastic.NewBoolQuery()
		for _, part := range parts {
			termQuery, err := qb.parseSingleTerm(part)
			if err != nil {
				return nil, err
			}
			boolQuery = boolQuery.Should(termQuery)
		}
		boolQuery = boolQuery.MinimumShouldMatch("1")
		return boolQuery, nil
	}

	// 检查是否包含 AND
	andRegex := regexp.MustCompile(` AND `)
	if andRegex.MatchString(innerExpr) {
		parts := andRegex.Split(innerExpr, -1)
		boolQuery := elastic.NewBoolQuery()
		for _, part := range parts {
			termQuery, err := qb.parseSingleTerm(part)
			if err != nil {
				return nil, err
			}
			boolQuery = boolQuery.Must(termQuery)
		}
		return boolQuery, nil
	}

	// 单个条件
	return qb.parseSingleTerm(innerExpr)
}

func (qb *QueryBuilder) handleSimpleExpression(expr string, currentQuery *elastic.BoolQuery) (*elastic.BoolQuery, error) {
	switch {
	case strings.HasPrefix(expr, "AND NOT "):
		termExpr := expr[8:] // 去掉 "AND NOT "
		termQuery, err := qb.parseSingleTerm(termExpr)
		if err != nil {
			return nil, err
		}
		return currentQuery.MustNot(termQuery), nil

	case strings.HasPrefix(expr, "AND "):
		termExpr := expr[4:] // 去掉 "AND "
		termQuery, err := qb.parseSingleTerm(termExpr)
		if err != nil {
			return nil, err
		}
		return currentQuery.Must(termQuery), nil

	case strings.HasPrefix(expr, "OR "):
		termExpr := expr[3:] // 去掉 "OR "
		termQuery, err := qb.parseSingleTerm(termExpr)
		if err != nil {
			return nil, err
		}
		return currentQuery.Should(termQuery), nil

	default:
		return nil, fmt.Errorf("invalid expression format: %s", expr)
	}
}

func (qb *QueryBuilder) parseSingleTerm(term string) (elastic.Query, error) {
	term = strings.TrimSpace(term)

	// 使用正则表达式匹配 field:(value) 或 field:value 格式
	termRegex := regexp.MustCompile(`^([^:]+):\(([^)]+)\)$|^([^:]+):(.+)$`)
	matches := termRegex.FindStringSubmatch(term)

	if matches == nil {
		return nil, fmt.Errorf("invalid term format: %s", term)
	}

	var field, value string
	if matches[1] != "" && matches[2] != "" {
		// field:(value) 格式
		field = strings.TrimSpace(matches[1])
		value = strings.TrimSpace(matches[2])
	} else if matches[3] != "" && matches[4] != "" {
		// field:value 格式
		field = strings.TrimSpace(matches[3])
		value = strings.TrimSpace(matches[4])
	} else {
		return nil, fmt.Errorf("invalid term format: %s", term)
	}

	field = strings.ReplaceAll(field, `\/`, `/`)
	value = strings.ReplaceAll(value, `\/`, `/`)

	return elastic.NewTermQuery(field, value), nil
}

func (qb *QueryBuilder) MergeBoolQueries(queries []*elastic.BoolQuery, combineType CombineType) *elastic.BoolQuery {
	if len(queries) == 0 {
		return elastic.NewBoolQuery()
	}
	if len(queries) == 1 {
		return queries[0]
	}

	switch combineType {
	case CombineMust:
		return qb.mergeWithMust(queries)
	case CombineShould:
		return qb.mergeWithShould(queries)
	case CombineMustNot:
		return qb.mergeWithMustNot(queries)
	default:
		return qb.mergeWithMust(queries)
	}
}

// mergeWithMust 以AND关系合并查询
func (qb *QueryBuilder) mergeWithMust(queries []*elastic.BoolQuery) *elastic.BoolQuery {
	result := elastic.NewBoolQuery()
	for _, query := range queries {
		result = result.Must(query)
	}
	return result
}

// mergeWithShould 以OR关系合并查询
func (qb *QueryBuilder) mergeWithShould(queries []*elastic.BoolQuery) *elastic.BoolQuery {
	result := elastic.NewBoolQuery()
	for _, query := range queries {
		result = result.Should(query)
	}
	result = result.MinimumShouldMatch("1")
	return result
}

// mergeWithMustNot 以MustNot关系合并查询
func (qb *QueryBuilder) mergeWithMustNot(queries []*elastic.BoolQuery) *elastic.BoolQuery {
	result := elastic.NewBoolQuery()
	for _, query := range queries {
		result = result.MustNot(query)
	}
	return result
}

// CombineType 合并类型
type CombineType int

const (
	CombineMust CombineType = iota
	CombineShould
	CombineMustNot
)

// ExecuteSearch 执行搜索
func (qb *QueryBuilder) ExecuteSearch(ctx context.Context, index string, query elastic.Query) (*elastic.SearchResult, error) {
	if qb.client == nil {
		return nil, fmt.Errorf("elasticsearch client is not initialized")
	}

	return qb.client.Search().
		Index(index).
		Query(query).
		Pretty(true).
		Do(ctx)
}

// PrintQueryJSON 打印查询的JSON格式（用于调试）
func (qb *QueryBuilder) PrintQueryJSON(query elastic.Query) (string, error) {
	source, err := query.Source()
	if err != nil {
		return "", err
	}

	prettyJSON, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		return "", err
	}

	return string(prettyJSON), nil
}

func (qb *QueryBuilder) CreateRangeQuery(field string, from, to interface{}) elastic.Query {
	return elastic.NewRangeQuery(field).From(from).To(to)
}

func (qb *QueryBuilder) CreateMatchQuery(field, value string) elastic.Query {
	return elastic.NewMatchQuery(field, value)
}

func (qb *QueryBuilder) CreateWildcardQuery(field, pattern string) elastic.Query {
	return elastic.NewWildcardQuery(field, pattern)
}

func (qb *QueryBuilder) AddPagination(searchService *elastic.SearchService, page, size int) *elastic.SearchService {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	from := (page - 1) * size
	return searchService.From(from).Size(size)
}

func (qb *QueryBuilder) AddSort(searchService *elastic.SearchService, field string, ascending bool) *elastic.SearchService {
	return searchService.Sort(field, ascending)
}

func (qb *QueryBuilder) GetTotalPages(totalHits int64, pageSize int) int {
	if pageSize <= 0 {
		return 1
	}
	totalPages := int(totalHits) / pageSize
	if int(totalHits)%pageSize > 0 {
		totalPages++
	}
	return totalPages
}
