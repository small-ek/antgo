package json

import (
	"encoding/json"
	. "github.com/small-ek/ginp/conv"
	"log"
)

type Result struct {
	Child interface{}
}

func DecodeJson(data string) *Result {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println(err.Error())
	}
	return &Result{
		Child: result,
	}
}

func DecodeArray(data string) *Result {
	var result []interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println(err.Error())
	}
	return &Result{
		Child: result,
	}
}

func (this *Result) Get(name interface{}) *Result {
	var child = this.Child
	switch child.(type) {
	case map[string]interface{}:
		return &Result{
			Child: child.(map[string]interface{})[String(name)],
		}
	case map[string]string:
		return &Result{
			Child: child.(map[string]string)[String(name)],
		}
	case []interface{}:
		return &Result{
			Child: child.([]interface{})[Int(name)],
		}
	case []string:
		return &Result{
			Child: child.([]interface{})[Int(name)],
		}
	case []int:
		return &Result{
			Child: child.([]interface{})[Int(name)],
		}
	case []int64:
		return &Result{
			Child: child.([]interface{})[Int(name)],
		}
	}
	return &Result{
		Child: child,
	}
}

func (this *Result) String() string {
	return String(this.Child)
}

func (this *Result) Int() int {
	return Int(this.Child)
}

func (this *Result) Int64() int64 {
	return Int64(this.Child)
}

func (this *Result) Float64() float64 {
	return Float64(this.Child)
}

func (this *Result) Map() map[string]interface{} {
	return this.Child.(map[string]interface{})
}

func (this *Result) Array() []interface{} {
	return this.Child.([]interface{})
}
