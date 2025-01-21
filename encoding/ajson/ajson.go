package ajson

import (
	"encoding/json"
	"github.com/small-ek/antgo/utils/conv"
	"os"
	"strings"
)

// Json is a wrapper for JSON data parsing and manipulation
// Json 是用于JSON数据解析和操作的封装结构
type Json struct {
	Child interface{} // Current level JSON data (当前层级的JSON数据)
}

// Open reads a JSON file and returns its content as bytes
// Open 读取JSON文件并返回其内容的字节切片
func Open(file string) []byte {
	byteValue, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return byteValue
}

// Decode parses JSON bytes into Json structure
// Decode 将JSON字节解析为Json结构
func Decode(data []byte) *Json {
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}
	return &Json{Child: result}
}

// Encode converts data to JSON string
// Encode 将数据转换为JSON字符串
func Encode(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(result)
}

// Next navigates to next level JSON element
// Next 导航到下一级JSON元素
func (j *Json) Next(name interface{}) *Json {
	switch child := j.Child.(type) {
	case map[string]interface{}: // Handle JSON object (处理JSON对象)
		return &Json{Child: child[conv.String(name)]}
	case []interface{}: // Handle JSON array (处理JSON数组)
		if index := conv.Int(name); index >= 0 && index < len(child) {
			return &Json{Child: child[index]}
		}
		return &Json{} // Return empty for invalid index (索引越界返回空)
	default: // Handle other types (处理其他类型)
		return &Json{}
	}
}

// Get navigates through multi-level JSON path
// Get 通过多级路径导航JSON数据
func (j *Json) Get(path string) *Json {
	current := j
	for _, part := range strings.Split(path, ".") {
		current = current.Next(part)
	}
	return current
}

// String converts current value to string
// String 将当前值转换为字符串
func (j *Json) String() string {
	return conv.String(j.Child)
}

// Strings converts current value to string slice
// Strings 将当前值转换为字符串切片
func (j *Json) Strings() []string {
	return conv.Strings(j.Child)
}

// Int converts current value to int
// Int 将当前值转换为整型
func (j *Json) Int() int {
	return conv.Int(j.Child)
}

// Ints converts current value to int slice
// Ints 将当前值转换为整型切片
func (j *Json) Ints() []int {
	return conv.Ints(j.Child)
}

// Int64 converts current value to int64
// Int64 将当前值转换为64位整型
func (j *Json) Int64() int64 {
	return conv.Int64(j.Child)
}

// Float64 converts current value to float64
// Float64 将当前值转换为64位浮点数
func (j *Json) Float64() float64 {
	return conv.Float64(j.Child)
}

// Bool converts current value to boolean
// Bool 将当前值转换为布尔值
func (j *Json) Bool() bool {
	return conv.Bool(j.Child)
}

// Map converts current value to map[string]interface{}
// Map 将当前值转换为字典类型
func (j *Json) Map() map[string]interface{} {
	if m, ok := j.Child.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// Array converts current value to []interface{}
// Array 将当前值转换为切片类型
func (j *Json) Array() []interface{} {
	if a, ok := j.Child.([]interface{}); ok {
		return a
	}
	return nil
}

// Interface returns raw interface value
// Interface 返回原始接口值
func (j *Json) Interface() interface{} {
	return j.Child
}

// 其他类型转换方法根据需要添加...
