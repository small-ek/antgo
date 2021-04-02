package jsons

import (
	"encoding/json"
	"github.com/small-ek/antgo/conv"
	"io/ioutil"
	"log"
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
		log.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		log.Println(err2)
	}
	return byteValue
}

//Decode Parse array.<解析json字符串>
func Decode(data string) *Json {
	var result interface{}
	data = strings.Trim(strings.Trim(data, "\n"), " ")

	if data != "" && string(data[0]) == "[" {
		result = []interface{}{}
	} else {
		result = make(map[string]interface{})
	}

	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println(err.Error())
	}
	return &Json{
		Child: result,
	}
}

//Encode Parses map to JSON string<解析map转json字符串>
func Encode(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
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

//MapString Data type conversion.
func (j *Json) MapString() map[string]string {
	return j.Child.(map[string]string)
}

//Array Data type conversion.
func (j *Json) Array() []interface{} {
	return j.Child.([]interface{})
}

//Strings Data type conversion.
func (j *Json) Strings() []string {
	return conv.Strings(j.Child)
}

//Ints Data type conversion.
func (j *Json) Ints() []int {
	return conv.Ints(j.Child)
}
