package jsons

import (
	"encoding/json"
	"log"
)

type Json struct {
	Json map[string]interface{}
}

func DecodeJson(data string) *Json {
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(data), &mapResult)
	if err != nil {
		log.Println(err.Error())
	}
	return &Json{
		Json: mapResult,
	}
}

/*func (this *Json) Get(name string) *Json {
	var json = this.Json[name]
	log.Println(json)
	log.Println(reflect.TypeOf(json))
	return &Json{
		Json: json,
	}
}*/
