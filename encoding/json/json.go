package json

import (
	"encoding/json"
	"github.com/small-ek/ginp/conv"
	"log"
)

//New Json parameter structure.
type New struct {
	Child interface{} //json next level.
}

//DecodeJson Parse json.
func DecodeJson(data string) *New {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println(err.Error())
	}
	return &New{
		Child: result,
	}
}

//DecodeArray Parse array.
func DecodeArray(data string) *New {
	var result []interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println(err.Error())
	}
	return &New{
		Child: result,
	}
}

//Get the next level of array or json.
func (this *New) Get(name interface{}) *New {
	var child = this.Child
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
func (this *New) String() string {
	return conv.String(this.Child)
}

//Int Data type conversion.
func (this *New) Int() int {
	return conv.Int(this.Child)
}

//Int64 Data type conversion.
func (this *New) Int64() int64 {
	return conv.Int64(this.Child)
}

//Float64 Data type conversion.
func (this *New) Float64() float64 {
	return conv.Float64(this.Child)
}

//Map Data type conversion.
func (this *New) Map() map[string]interface{} {
	return this.Child.(map[string]interface{})
}

//Array Data type conversion.
func (this *New) Array() []interface{} {
	return this.Child.([]interface{})
}
