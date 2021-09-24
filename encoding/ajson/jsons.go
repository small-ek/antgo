package ajson

import (
	"bytes"
	"encoding/json"
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/logs"
	"io/ioutil"
	"os"
	"strings"
)

//Json Json parameter structure.
type Json struct {
	Child interface{} //json next level.
}

//Open 读取json文件
func Open(file string) []byte {
	jsonFile, err := os.Open(file)
	if err != nil {
		logs.Error(err.Error())
	}
	defer jsonFile.Close()

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		logs.Error(err2.Error())
	}
	return byteValue
}

//Decode Parse array.<解析json字符串>
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
		logs.Error(err.Error())
	}
	return &Json{
		Child: result,
	}
}

//Encode Parses map to JSON string<解析map转json字符串>
func Encode(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		logs.Error(err.Error())
	}
	return string(result)
}

//Get the next level of array or json.<获取json的节点>
func (j *Json) Get(name interface{}) *Json {
	var child = j.Child
	switch child.(type) {
	case map[string]interface{}:
		return &Json{
			Child: child.(map[string]interface{})[conv.String(name)],
		}
	case map[string]string:
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

//Read Parse json according to point split
func (j *Json) Read(name string) *Json {
	var list = strings.Split(name, ".")
	for i := 0; i < len(list); i++ {
		var value = list[i]
		var result = j.Get(value)
		j.Child = result.Child
	}
	return j
}

//String Data type conversion.
func (j *Json) String() string {
	return conv.String(j.Child)
}

//Int Data type conversion.
func (j *Json) Int() int {
	return conv.Int(j.Child)
}

//Int32 Data type conversion.
func (j *Json) Int32() int32 {
	return conv.Int32(j.Child)
}

//Int64 Data type conversion.
func (j *Json) Int64() int64 {
	return conv.Int64(j.Child)
}

//Float32 Data type conversion.
func (j *Json) Float32() float32 {
	return conv.Float32(j.Child)
}

//Float64 Data type conversion.
func (j *Json) Float64() float64 {
	return conv.Float64(j.Child)
}

//Map Data type conversion.
func (j *Json) Map() map[string]interface{} {
	return j.Child.(map[string]interface{})
}

//MapInt Data type conversion.
func (j *Json) MapInt() map[int]interface{} {
	return j.Child.(map[int]interface{})
}

//MapInt64 Data type conversion.
func (j *Json) MapInt64() map[int64]interface{} {
	return j.Child.(map[int64]interface{})
}

//MapString Data type conversion.
func (j *Json) MapString() map[string]string {
	return j.Child.(map[string]string)
}

//Array Data type conversion.
func (j *Json) Array() []interface{} {
	return j.Child.([]interface{})
}

//Bytes Data type conversion.
func (j *Json) Bytes() []byte {
	return conv.Bytes(j.Child)
}

//Strings Data type conversion.
func (j *Json) Strings() []string {
	return conv.Strings(j.Child)
}

//Ints Data type conversion.
func (j *Json) Ints() []int {
	return conv.Ints(j.Child)
}

//Int32s Data type conversion.
func (j *Json) Int32s() []int32 {
	return conv.Int32s(j.Child)
}

//Int64s Data type conversion.
func (j *Json) Int64s() []int64 {
	return conv.Int64s(j.Child)
}

//Uints Data type conversion.
func (j *Json) Uints() []uint {
	return conv.Uints(j.Child)
}

//Uint32s Data type conversion.
func (j *Json) Uint32s() []uint32 {
	return conv.Uint32s(j.Child)
}

//Uint64s Data type conversion.
func (j *Json) Uint64s() []uint64 {
	return conv.Uint64s(j.Child)
}

//Interfaces Data type conversion.
func (j *Json) Interfaces() []interface{} {
	return conv.Interfaces(j.Child)
}

//Interface Data type conversion.
func (j *Json) Interface() interface{} {
	return j.Child
}
