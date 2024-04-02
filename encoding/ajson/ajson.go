package ajson

import (
	"bytes"
	"encoding/json"
	"github.com/small-ek/antgo/utils/conv"
	"io/ioutil"
	"strings"
)

// Json Json parameter structure.
type Json struct {
	Child interface{} //json next level.
}

// Open 读取json文件
func Open(file string) []byte {
	byteValue, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return byteValue
}

// Decode Parse array.<解析json字符串>
func Decode(data []byte) *Json {
	var result interface{}
	data = bytes.Trim(bytes.Trim(data, "\n"), " ")
	str := string(data)
	if str != "" && string(str[0]) == "[" {
		result = []interface{}{}
	} else {
		result = make(map[string]interface{})
	}

	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return &Json{
		Child: result,
	}
}

// Encode Parses map to JSON string<解析map转json字符串>
func Encode(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(result)
}

// Next the next level of array or json.<获取json的节点>
func (j *Json) Next(name interface{}) *Json {
	var child = j.Child
	switch child.(type) {
	case map[string]interface{}:
		return &Json{
			Child: child.(map[string]interface{})[conv.String(name)],
		}
	case map[interface{}]interface{}:
		return &Json{
			Child: child.(map[interface{}]interface{})[conv.String(name)],
		}
	case map[string]string:
		return &Json{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case map[string]int:
		return &Json{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case []interface{}:
		return &Json{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []string:
		return &Json{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int:
		return &Json{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int64:
		return &Json{
			Child: child.([]interface{})[conv.Int(name)],
		}
	}
	return &Json{
		Child: child,
	}
}

// Get Parse json according to point split
func (j *Json) Get(name string) *Json {
	var list = strings.Split(name, ".")
	for i := 0; i < len(list); i++ {
		var value = list[i]
		var result = j.Next(value)
		j.Child = result.Child
	}
	return j
}

// String Data type conversion.
func (j *Json) String() string {
	if j.Child == nil {
		return ""
	}
	defer j.End()
	return conv.String(j.Child)
}

// End Data type conversion.
func (j *Json) End() {
	j.Child = nil
}

// Strings Data type conversion.
func (j *Json) Strings() []string {
	defer j.End()
	return conv.Strings(j.Child)
}

// Byte Data type conversion.
func (j *Json) Byte() byte {
	defer j.End()
	return conv.Byte(j.Child)
}

// Bytes Data type conversion.
func (j *Json) Bytes() []byte {
	defer j.End()
	return conv.Bytes(j.Child)
}

// Int Data type conversion.
func (j *Json) Int() int {
	if j.Child == nil {
		return 0
	}
	defer j.End()
	return conv.Int(j.Child)
}

// Bool Data type conversion.
func (j *Json) Bool() bool {
	defer j.End()
	return conv.Bool(j.Child)
}

// Ints Data type conversion.
func (j *Json) Ints() []int {
	defer j.End()
	return conv.Ints(j.Child)
}

// Int64 Data type conversion.
func (j *Json) Int64() int64 {
	if j.Child == nil {
		return 0
	}
	defer j.End()
	return conv.Int64(j.Child)
}

// Float64 Data type conversion.
func (j *Json) Float64() float64 {
	if j.Child == nil {
		return 0
	}
	defer j.End()
	return conv.Float64(j.Child)
}

// Map Data type conversion.
func (j *Json) Map() map[string]interface{} {
	if j.Child == nil {
		return nil
	}
	defer j.End()
	return j.Child.(map[string]interface{})
}

// Maps Data type conversion.
func (j *Json) Maps() []map[string]interface{} {
	if j.Child == nil {
		return nil
	}
	defer j.End()
	return j.Child.([]map[string]interface{})
}

// Array Data type conversion.
func (j *Json) Array() []interface{} {
	if j.Child == nil {
		return nil
	}
	defer j.End()
	return j.Child.([]interface{})
}

// Uint Data type conversion.
func (j *Json) Uint() uint {
	defer j.End()
	return conv.Uint(j.Child)
}

// Uint8 Data type conversion.
func (j *Json) Uint8() uint8 {
	defer j.End()
	return conv.Uint8(j.Child)
}

// Uint16 Data type conversion.
func (j *Json) Uint16() uint16 {
	defer j.End()
	return conv.Uint16(j.Child)
}

// Uint32 Data type conversion.
func (j *Json) Uint32() uint32 {
	defer j.End()
	return conv.Uint32(j.Child)
}

// Uint64 Data type conversion.
func (j *Json) Uint64() uint64 {
	defer j.End()
	return conv.Uint64(j.Child)
}

// Interfaces Data type conversion.
func (j *Json) Interfaces() []interface{} {
	return conv.Interfaces(j.Child)
}

// Interface Data type conversion.
func (j *Json) Interface() interface{} {
	return j.Child
}
