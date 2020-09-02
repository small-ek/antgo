package conv

import (
	"encoding/json"
	"log"
)

//map转换
func Map(i interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	result, err := json.Marshal(i)
	if err != nil {
		log.Println(err.Error())
	}

	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Println(err.Error())
	}
	return data
}

//将<i>转换 map string
func MapString(i interface{}) map[string]string {
	var data = make(map[string]string)
	result, err := json.Marshal(i)

	if err != nil {
		log.Println(err.Error())
	}

	err = json.Unmarshal(result, &data)

	if err != nil {
		log.Println(err.Error())
	}
	return data
}
