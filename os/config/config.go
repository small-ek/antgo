package config

import (
	"github.com/BurntSushi/toml"
	"github.com/small-ek/ginp/conv"
	"log"
	"os"
)

var Data map[string]interface{}

//SetPath Set path.
func SetPath(path string) {
	if _, err := toml.DecodeFile(path, &Data); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

type Result struct {
	Child interface{}
}

//Default config
func Default() *Result {
	return &Result{
		Child: Data,
	}
}

//Get config
func (get *Result) Get(name interface{}) *Result {
	var child = get.Child
	switch child.(type) {
	case map[string]interface{}:
		return &Result{
			Child: child.(map[string]interface{})[conv.String(name)],
		}
	case map[string]string:
		return &Result{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case map[string]int:
		return &Result{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case []interface{}:
		return &Result{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []string:
		return &Result{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int:
		return &Result{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int64:
		return &Result{
			Child: child.([]interface{})[conv.Int(name)],
		}
	}
	return &Result{
		Child: child,
	}
}

//String Data type conversion.
func (get *Result) String() string {
	if get.Child == nil {
		return ""
	}
	return conv.String(get.Child)
}

//String Data type conversion.
func (get *Result) Strings() []string {
	return conv.Strings(get.Child)
}

//Int Data type conversion.
func (get *Result) Int() int {
	if get.Child == nil {
		return 0
	}
	return conv.Int(get.Child)
}

//Ints Data type conversion.
func (get *Result) Ints() []int {
	return conv.Ints(get.Child)
}

//Int64 Data type conversion.
func (get *Result) Int64() int64 {
	if get.Child == nil {
		return 0
	}
	return conv.Int64(get.Child)
}

//Int64 Data type conversion.
func (get *Result) Float64() float64 {
	if get.Child == nil {
		return 0
	}
	return conv.Float64(get.Child)
}

//Map Data type conversion.
func (get *Result) Map() map[string]interface{} {
	if get.Child == nil {
		return nil
	}
	return get.Child.(map[string]interface{})
}

//Array Data type conversion.
func (get *Result) Array() []interface{} {
	if get.Child == nil {
		return nil
	}
	return get.Child.([]interface{})
}
