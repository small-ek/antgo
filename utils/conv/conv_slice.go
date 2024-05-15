package conv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Ints 将 `any` 转换为 []int。
func Ints(any interface{}) []int {
	if any == nil {
		return nil
	}

	switch value := any.(type) {
	case []string:
		return convertStringSliceToIntSlice(value)
	case []int:
		return value
	case string:
		var array []int
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
		return array
	default:
		// 检查 'value' 的类型是否为切片
		valueType := reflect.TypeOf(value)
		if valueType.Kind() != reflect.Slice {
			panic(fmt.Sprintf("Unwilling type：%T", any))
		}

		// 获取切片的基础类型
		elemType := valueType.Elem()

		// 检查基础类型是否为数值类型或布尔类型
		switch elemType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return convertNumericSliceToIntSlice(reflect.ValueOf(value))
		case reflect.Bool:
			return convertBoolSliceToIntSlice(reflect.ValueOf(value))
		case reflect.Interface:
			return convertInterfaceSliceToIntSlice(reflect.ValueOf(value))
		default:
			panic(fmt.Sprintf("Unwilling type：%T", any))
		}
	}
}

// convertStringSliceToIntSlice 将字符串切片转换为 []int
func convertStringSliceToIntSlice(value []string) []int {
	array := make([]int, len(value))
	for i, v := range value {
		array[i], _ = strconv.Atoi(v) // 处理错误
	}
	return array
}

// convertNumericSliceToIntSlice 将数值类型切片转换为 []int
func convertNumericSliceToIntSlice(value reflect.Value) []int {
	length := value.Len()
	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = int(value.Index(i).Convert(reflect.TypeOf(0)).Int())
	}
	return array
}

// convertBoolSliceToIntSlice 将布尔切片转换为 []int
func convertBoolSliceToIntSlice(value reflect.Value) []int {
	length := value.Len()
	array := make([]int, length)
	for i := 0; i < length; i++ {
		if value.Index(i).Bool() {
			array[i] = 1
		} else {
			array[i] = 0
		}
	}
	return array
}

// 将接口切片转换为 []int
func convertInterfaceSliceToIntSlice(value reflect.Value) []int {
	length := value.Len()
	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = Int(value.Index(i).Interface())
	}
	return array
}

// Strings 将 `any` 转换为 []string。
func Strings(any interface{}) []string {
	if any == nil {
		return nil
	}

	switch value := any.(type) {
	case []string:
		return value
	case string:
		var array []string
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
		return array
	}

	v := reflect.ValueOf(any)
	if v.Kind() != reflect.Slice {
		panic("Unwilling type")
	}

	length := v.Len()
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result[i] = convertToString(v.Index(i))
	}
	return result
}

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
	}

	panic("Unwilling type")
}

// Interfaces converts `any` to []interface{}.<将“any”转换为[]interface{}。>
func Interfaces(any interface{}) []interface{} {
	if any == nil {
		return nil
	}
	var array []interface{}
	switch value := any.(type) {
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []bool, []float32, []float64, []interface{}, []string, [][]byte:
		if err := json.Unmarshal(Bytes(value), &array); err != nil {
			panic(err)
		}

	case string:
		if err := json.Unmarshal([]byte(value), &array); err != nil {
			panic(err)
		}
	}
	return array
}
