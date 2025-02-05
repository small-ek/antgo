package conv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Ints 将 `any` 转换为 []int。
// Converts `any` to []int.
func Ints(any interface{}) []int {
	if any == nil {
		return nil
	}

	switch value := any.(type) {
	case []string:
		// 直接将字符串切片转换为整数切片
		// Directly convert a string slice to an int slice
		return convertStringSliceToIntSlice(value)
	case []int:
		// 已经是 []int，直接返回
		// Already []int, return directly
		return value
	case string:
		// 尝试将字符串解析为 JSON 数组
		// Try to parse the string as a JSON array
		var array []int
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(fmt.Sprintf("Error unmarshaling string: %v", err))
		}
		return array
	default:
		// 对未知类型进行反射处理
		// Use reflection for unknown types
		return handleUnknownTypeForInts(any)
	}
}

// handleUnknownTypeForInts 处理其他类型的转换。
// Handles conversion for other types.
func handleUnknownTypeForInts(any interface{}) []int {
	// 反射获取类型
	// Get the type via reflection
	valueType := reflect.TypeOf(any)
	if valueType.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Unwilling type: expected slice, got %T", any))
	}

	// 获取切片元素类型
	// Get the slice element type
	elemType := valueType.Elem()

	// 根据元素类型进行转换
	// Convert based on the element type
	switch elemType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// 数值类型，直接转换
		// Numeric types, convert directly
		return convertNumericSliceToIntSlice(reflect.ValueOf(any))
	case reflect.Bool:
		// 布尔类型，转换为 0 或 1
		// Boolean types, convert to 0 or 1
		return convertBoolSliceToIntSlice(reflect.ValueOf(any))
	case reflect.Interface:
		// 接口类型，递归转换
		// Interface types, convert recursively
		return convertInterfaceSliceToIntSlice(reflect.ValueOf(any))
	default:
		panic(fmt.Sprintf("Unwilling type: %T", any))
	}
}

// convertStringSliceToIntSlice 将字符串切片转换为 []int。
// Converts a string slice to []int.
func convertStringSliceToIntSlice(value []string) []int {
	// 优化：一次性创建切片，避免重复扩容
	// Optimization: Create the slice at once to avoid repeated resizing
	array := make([]int, len(value))
	for i, v := range value {
		// 错误处理：在实际环境中不应忽略错误
		// Error handling: In a real-world scenario, errors should not be ignored
		if n, err := strconv.Atoi(v); err == nil {
			array[i] = n
		} else {
			array[i] = 0 // 默认值，或者记录错误
			// Default value, or log the error
		}
	}
	return array
}

// convertNumericSliceToIntSlice 将数值类型切片转换为 []int。
// Converts a numeric slice to []int.
func convertNumericSliceToIntSlice(value reflect.Value) []int {
	length := value.Len()
	array := make([]int, length)
	for i := 0; i < length; i++ {
		// 优化：减少重复的类型转换操作
		// Optimization: Reduce redundant type conversions
		array[i] = int(value.Index(i).Convert(reflect.TypeOf(0)).Int())
	}
	return array
}

// convertBoolSliceToIntSlice 将布尔类型切片转换为 []int。
// Converts a boolean slice to []int.
func convertBoolSliceToIntSlice(value reflect.Value) []int {
	length := value.Len()
	array := make([]int, length)
	for i := 0; i < length; i++ {
		// 优化：避免反射过多，直接使用 Bool() 方法
		// Optimization: Avoid excessive reflection, directly use Bool() method
		if value.Index(i).Bool() {
			array[i] = 1
		} else {
			array[i] = 0
		}
	}
	return array
}

// convertInterfaceSliceToIntSlice 将接口类型切片转换为 []int。
// Converts an interface slice to []int.
func convertInterfaceSliceToIntSlice(value reflect.Value) []int {
	length := value.Len()
	array := make([]int, length)
	for i := 0; i < length; i++ {
		// 递归处理每个元素，利用现有的 Int 函数进行转换
		// Recursively handle each element, using the existing Int function for conversion
		array[i] = Int(value.Index(i).Interface())
	}
	return array
}

// Strings 将 `any` 转换为 []string。
// Converts `any` to []string.
func Strings(any interface{}) []string {
	if any == nil {
		return nil
	}

	switch value := any.(type) {
	case []string:
		// 已经是 []string，直接返回
		// Already []string, return directly
		return value
	case string:
		// 尝试将字符串解析为 JSON 数组
		// Try to parse the string as a JSON array
		var array []string
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(fmt.Sprintf("Error unmarshaling string: %v", err))
		}
		return array
	}

	// 反射处理切片类型
	// Reflection handling for slice types
	v := reflect.ValueOf(any)
	if v.Kind() != reflect.Slice {
		panic("Unwilling type: expected slice")
	}

	// 将每个元素转换为字符串
	// Convert each element to a string
	length := v.Len()
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result[i] = convertToString(v.Index(i))
	}
	return result
}

// convertToString 将任意类型转换为字符串。
// Converts any type to a string.
func convertToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.String:
		return v.String()
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes())
		}
	default:
		return ""
	}

	panic(fmt.Sprintf("Unwilling type: %T", v.Interface()))
}

// Interfaces 将 `any` 转换为 []interface{}。
// Converts `any` to []interface{}.
func Interfaces(any interface{}) []interface{} {
	if any == nil {
		return nil
	}

	var array []interface{}
	switch value := any.(type) {
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []bool, []float32, []float64, []interface{}, []string, [][]byte:
		// 优化：一次性将切片转换为字节，再解码成 []interface{}
		// Optimization: Convert the slice to bytes at once and then decode into []interface{}
		if err := json.Unmarshal(Bytes(value), &array); err != nil {
			panic(fmt.Sprintf("Error unmarshaling slice: %v", err))
		}

	case string:
		// 字符串解析为 JSON 数组
		// Parse the string as a JSON array
		if err := json.Unmarshal([]byte(value), &array); err != nil {
			panic(fmt.Sprintf("Error unmarshaling string: %v", err))
		}
	}
	return array
}
