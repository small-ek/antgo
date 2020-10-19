package config

import (
	"github.com/BurntSushi/toml"
	. "github.com/small-ek/ginp/conv"
	"log"
	"os"
)

var Data map[string]interface{}

func SetPath(path string) {
	if _, err := toml.DecodeFile(path, &Data); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

type Result struct {
	Child interface{}
}

func Default() *Result {
	return &Result{
		Child: Data,
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
	if this.Child == nil {
		return ""
	}
	return String(this.Child)
}

func (this *Result) Int() int {
	if this.Child == nil {
		return 0
	}
	return Int(this.Child)
}

func (this *Result) Int64() int64 {
	if this.Child == nil {
		return 0
	}
	return Int64(this.Child)
}

func (this *Result) Float64() float64 {
	if this.Child == nil {
		return 0
	}
	return Float64(this.Child)
}

func (this *Result) Map() map[string]interface{} {
	if this.Child == nil {
		return nil
	}
	return this.Child.(map[string]interface{})
}

func (this *Result) Array() []interface{} {
	if this.Child == nil {
		return nil
	}
	return this.Child.([]interface{})
}
