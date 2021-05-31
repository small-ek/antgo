package i18n

import (
	"encoding/json"
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/logs"
	"io/ioutil"
	"strings"
)

//Datas ...
var Datas map[string]interface{}

//Maps ...
var Maps map[string]interface{}

//Array ...
var Array []interface{}

//SetPath Set path
func SetPath(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err2 := json.Unmarshal(bytes, &Datas)
	if err2 != nil {
		logs.Error(err.Error())
		return
	}
}

//SetLanguage language setting
func SetLanguage(lang string) {
	if Datas[lang] != nil {
		Maps = Datas[lang].(map[string]interface{})
	}
}

//Get language
func Get(name string) string {
	var list = strings.Split(name, ".")
	if len(list) > 1 {
		for _, value := range list {
			if value != "" {
				switch values := Maps[value].(type) {
				case map[string]interface{}:
					Maps = values
					names := name[len(value)+1:]
					return Get(names)
				case string:
					return conv.String(Maps)
				case []interface{}:
					Array = values
					names := name[len(value)+1:]
					return array(names)
				}
			}
		}
	}

	if Maps[name] != nil {
		return conv.String(Maps[name])
	} else {
		return name
	}
}

//array Slice analysis
func array(name string) string {
	var list = strings.Split(name, ".")
	if len(list) > 1 {
		for _, value := range list {
			if value != "" {
				switch values := Array[conv.Int(value)].(type) {
				case map[string]interface{}:
					Maps = values
					names := name[len(value)+1:]
					return Get(names)
				case string:
					return conv.String(values)
				case []interface{}:
					Array = values
					names := name[len(value)+1:]
					return array(names)
				}
			}
		}
	}
	return conv.String(Array[conv.Int(name)])
}
