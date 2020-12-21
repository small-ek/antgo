package json

import (
	"encoding/json"
	"github.com/small-ek/antgo/conv"
	"log"
	"strings"
)

//New Json parameter structure.
type New struct {
	Child interface{} //json next level.
}

//Decode Parse array.<解析json字符串>
func Decode(data string) *New {
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
	return &New{
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
func (get *New) Get(name interface{}) *New {
	var child = get.Child
	switch child.(type) {
	case map[string]interface{}:
		return &New{
			Child: child.(map[string]interface{})[conv.String(name)],
		}
	case map[string]string:
		return &New{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case []interface{}:
		return &New{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []string:
		return &New{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int:
		return &New{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int64:
		return &New{
			Child: child.([]interface{})[conv.Int(name)],
		}
	}
	return &New{
		Child: child,
	}
}

//String Data type conversion.
func (get *New) String() string {
	return conv.String(get.Child)
}

//Int Data type conversion.
func (get *New) Int() int {
	return conv.Int(get.Child)
}

//Int64 Data type conversion.
func (get *New) Int64() int64 {
	return conv.Int64(get.Child)
}

//Float64 Data type conversion.
func (get *New) Float64() float64 {
	return conv.Float64(get.Child)
}

//Map Data type conversion.
func (get *New) Map() map[string]interface{} {
	return get.Child.(map[string]interface{})
}

//Array Data type conversion.
func (get *New) Array() []interface{} {
	return get.Child.([]interface{})
}
