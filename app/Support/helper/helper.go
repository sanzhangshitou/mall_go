package helper

import (
	"encoding/json"
	"reflect"
)

// ToJSON 将对象转为 JSON 字符串 (用于日志等)
func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// IsEmpty 判断值是否为空 (nil / 零值 / 空字符串 / 空切片)
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.Len() == 0
	case reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

// Ptr 返回一个值的指针 (用于可选字段)
func Ptr[T any](v T) *T {
	return &v
}

// DefaultIfEmpty 如果值为空则返回默认值
func DefaultIfEmpty(val, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}

// InSlice 判断元素是否在切片中
func InSlice[T comparable](item T, slice []T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
