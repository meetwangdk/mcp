package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

// StructToMapStructJSON 将结构体转换为map[string]any, 并转换为JSON字符串
func StructToMapStructJSON(obj any) string {
	mapstructureObj := make(map[string]any)
	err := mapstructure.Decode(obj, &mapstructureObj)
	if err != nil {
		return ""
	}
	result, _ := json.Marshal(mapstructureObj)
	return string(result)
}

// StructsToMapStructJSON 将结构体数组转换为map[string]any数组, 并转换为JSON字符串
func StructsToMapStructJSON[T any](objs []T) string {
	mapstructureObjs := make([]map[string]any, 0)
	for _, obj := range objs {
		mapstructureObj := make(map[string]any)
		err := mapstructure.Decode(obj, &mapstructureObj)
		if err != nil {
			continue
		}
		mapstructureObjs = append(mapstructureObjs, mapstructureObj)
	}
	result, _ := json.Marshal(mapstructureObjs)
	return string(result)
}

// Dumps dump interface as indented json object
func Dumps(v interface{}) string {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("MarshalError(%s): %+v", err.Error(), v)
		return fmt.Sprintf("%+v", v)
	}
	return string(bs)
}

func RefreshCacheTime(c *cache.Cache, key string, value interface{}, defaultExpiration time.Duration) bool {
	_, found := c.Get(key)
	c.Set(key, value, defaultExpiration)
	return found
}
