package json

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

// 定义JSON操作
var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

// MarshalToString JSON编码为字符串
func MarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}

// ToMapInterface interface转换为map[string]interface{}
func ToMapInterface(v interface{}) map[string]interface{} {
	if v, ok := v.(map[string]interface{}); ok {
		return v
	}

	var m map[string]interface{}
	err := jsoniter.UnmarshalFromString(MarshalToString(v), &m)
	if err != nil {
		return nil
	}
	return m
}

// ToSliceInterface interface转换为[]interface{}
func ToSliceInterface(v interface{}) []interface{} {

	if v, ok := v.([]interface{}); ok {
		return v
	}
	var s []interface{}

	if v := reflect.ValueOf(v); v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			s = append(s, item)
		}
		return s
	}
	_ = jsoniter.UnmarshalFromString(MarshalToString(v), &s)
	return s
}
