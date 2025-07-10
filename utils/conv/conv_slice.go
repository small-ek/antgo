package conv

import (
	"reflect"
	"strconv"
)

// Ints 将 any 转换为 []int，解析失败或不支持类型时返回空切片。
// Ints converts any value to a slice of int; returns nil for unsupported types or parse failures.
func Ints(any interface{}) []int {
	// 空值检查 Empty value check
	if any == nil {
		return nil
	}

	switch v := any.(type) {
	case []int:
		// Already a []int, return directly.
		// 已经是 []int，直接返回
		return v
	case []string:
		// Convert slice of strings to []int
		// 将字符串切片转换为整数切片
		return convertStringSliceToIntSlice(v)
	case string:
		// Parse JSON array string to []int
		// 尝试将 JSON 数组字符串解析为整型切片
		var arr []int
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			return arr
		}
		return nil
	default:
		// Use reflection for other slice types
		// 对其他切片类型使用反射进行处理
		val := reflect.ValueOf(any)
		if val.Kind() != reflect.Slice {
			return nil
		}

		sliceKind := val.Type().Elem().Kind()
		switch sliceKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			// Numeric types conversion
			// 数值类型切片转换
			return convertNumericSliceToIntSlice(val)
		case reflect.Bool:
			// Boolean slice to []int
			// 布尔类型切片转换
			return convertBoolSliceToIntSlice(val)
		case reflect.Interface:
			// Recursive conversion for interface slices
			// 接口类型切片递归处理
			return convertInterfaceSliceToIntSlice(val)
		default:
			return nil
		}
	}
}

// Strings 将 any 转换为 []string，解析失败或不支持类型时返回空切片。
// Strings converts any value to a slice of string; returns nil for unsupported types or parse failures.
func Strings(any interface{}) []string {
	if any == nil {
		return nil
	}

	switch v := any.(type) {
	case []string:
		// Already a []string, return directly.
		// 已经是 []string，直接返回
		return v
	case string:
		// Parse JSON array string to []string
		// 尝试将 JSON 数组字符串解析为字符串切片
		var arr []string
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			return arr
		}
		return nil
	}

	// Reflection for other slice types
	// 对其他切片类型使用反射处理
	val := reflect.ValueOf(any)
	if val.Kind() != reflect.Slice {
		return nil
	}

	n := val.Len()
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = toString(val.Index(i))
	}
	return out
}

// Interfaces 将 any 转换为 []interface{}，解析失败或不支持类型时返回空切片。
// Interfaces converts any value to a slice of interface{}; returns nil for unsupported types or parse failures.
func Interfaces(any interface{}) []interface{} {
	if any == nil {
		return nil
	}

	switch v := any.(type) {
	case []interface{}:
		// Already a []interface{}, return directly.
		// 已经是 []interface{}，直接返回
		return v
	case []byte, []int, []string, []bool, [][]byte, []float64:
		// Marshal then unmarshal to []interface{}
		// 使用 JSON 编码再解码
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		var arr []interface{}
		if err := json.Unmarshal(b, &arr); err == nil {
			return arr
		}
		return nil
	case string:
		// Parse JSON array string to []interface{}
		// 尝试解析 JSON 数组字符串
		var arr []interface{}
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			return arr
		}
		return nil
	}
	return nil
}

// --- Internal utility functions ---

// convertStringSliceToIntSlice 将 []string 转 []int，忽略解析错误。
// convertStringSliceToIntSlice converts a slice of strings to a slice of ints, ignoring parse errors.
func convertStringSliceToIntSlice(src []string) []int {
	n := len(src)
	out := make([]int, n)
	for i, s := range src {
		if x, err := strconv.Atoi(s); err == nil {
			out[i] = x
		}
	}
	return out
}

// convertNumericSliceToIntSlice 将数值类型切片转换为 []int。
// convertNumericSliceToIntSlice converts numeric slices to a slice of ints.
func convertNumericSliceToIntSlice(val reflect.Value) []int {
	intType := reflect.TypeOf(0)
	n := val.Len()
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = int(val.Index(i).Convert(intType).Int())
	}
	return out
}

// convertBoolSliceToIntSlice 将布尔切片转换为 []int，其中 true=>1, false=>0。
// convertBoolSliceToIntSlice converts boolean slices to a slice of ints (true=>1, false=>0).
func convertBoolSliceToIntSlice(val reflect.Value) []int {
	n := val.Len()
	out := make([]int, n)
	for i := 0; i < n; i++ {
		if val.Index(i).Bool() {
			out[i] = 1
		}
	}
	return out
}

// convertInterfaceSliceToIntSlice 将 interface{} 切片递归转换为 []int。
// convertInterfaceSliceToIntSlice recursively converts interface{} slices to a slice of ints.
func convertInterfaceSliceToIntSlice(val reflect.Value) []int {
	n := val.Len()
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = firstInt(val.Index(i).Interface())
	}
	return out
}

// firstInt 尝试将单个任意类型的值转换为 int，优先级：int、float64、string、bool。
// firstInt attempts to convert a single value of various types to int.
func firstInt(any interface{}) int {
	if any == nil {
		return 0
	}
	switch v := any.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if x, err := strconv.Atoi(v); err == nil {
			return x
		}
	case bool:
		if v {
			return 1
		}
	}
	return 0
}

// toString 将 reflect.Value 转换为字符串，支持基本类型及 []byte。
// toString converts a reflect.Value to string, supporting basic types and []byte.
func toString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte 转字符串
			return string(v.Bytes())
		}
	default:
		return ""
	}
	return ""
}
