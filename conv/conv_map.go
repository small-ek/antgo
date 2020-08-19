package conv

import (
	"encoding/json"
	"log"
)

//map转换
func Map(i interface{}) map[string]interface{} {
	var data map[string]interface{}
	result, err := json.Marshal(i)
	if err != nil {
		log.Print("类型不正确" + err.Error())
	}

	json.Unmarshal(result, &data)
	return data
}

//将<i>转换 map string
func MapString(i interface{}) map[string]string {
	var data map[string]string
	result, err := json.Marshal(i)

	if err != nil {
		log.Print("类型不正确" + err.Error())
	}

	json.Unmarshal(result, &data)
	return data
}
