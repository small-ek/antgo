package conv

import (
	"fmt"
	"strconv"
)

// Map converts any type to map[string]interface{} with panic handling.
// Supported types: JSON bytes/string, map, struct
// 将任意类型转换为map[string]interface{}，包含panic处理
// 支持类型：JSON字节/字符串、map、结构体
func Map(any any) (data map[string]interface{}) {
	// Fast path for direct type conversion [[4,15]]
	if m, ok := any.(map[string]interface{}); ok {
		return m
	}

	bytes := Bytes(any)
	data = make(map[string]interface{})
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil
	}
	return
}

// Maps converts any type to []map[string]interface{} with panic handling.
// Auto-converts JSON strings/bytes to slice
// 将任意类型转换为[]map[string]interface{}，包含panic处理
// 自动转换JSON字符串/字节为切片
func Maps(any any) (data []map[string]interface{}) {
	// Fast path for direct type conversion [[4,15]]
	if slice, ok := any.([]map[string]interface{}); ok {
		return slice
	}

	bytes := Bytes(any)
	data = []map[string]interface{}{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil
	}
	return
}

// MapsInt converts any type to []map[string]interface{} with panic handling.
// Auto-converts JSON strings/bytes to slice
// 将任意类型转换为[]map[int]interface{}
// 自动转换JSON字符串/字节为切片
func MapsInt(any any) (data []map[int]interface{}) {
	// Fast path for direct type conversion [[4,15]]
	if slice, ok := any.([]map[int]interface{}); ok {
		return slice
	}

	bytes := Bytes(any)
	data = []map[int]interface{}{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil
	}
	return
}

// MapString converts any type to map[string]string with type compatibility.
// Supports all JSON value types converted to string via fmt.Sprintf
// 将任意类型转换为map[string]string，增强类型兼容性
// 支持通过fmt.Sprintf将所有JSON值类型转为字符串
func MapString(any any) map[string]string {
	bytes := Bytes(any)
	var temp map[string]interface{}
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return nil
	}

	data := make(map[string]string, len(temp))
	for k, v := range temp {
		data[k] = fmt.Sprintf("%v", v) // Universal conversion [[6,15]]
	}
	return data
}

// MapInt converts any type to map[int]interface{} with key validation.
// Requires JSON keys to be numeric strings (e.g., "123")
// 将任意类型转换为map[int]interface{}，包含键值验证
// 要求JSON键为数字字符串（如"123"）
func MapInt(any any) map[int]interface{} {
	bytes := Bytes(any)
	var temp map[string]interface{}
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return nil
	}

	data := make(map[int]interface{}, len(temp))
	for k, v := range temp {
		intKey, err := strconv.Atoi(k)
		if err != nil {
			return nil
		}
		data[intKey] = v
	}
	return data
}
