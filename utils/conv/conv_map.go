package conv

import (
	"encoding/json"
)

//Map converts `any` to map[string]interface{}.<将“any”转换为map[string]interface{}。>
func Map(any interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	result, err := json.Marshal(any)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(result, &data)
	if err != nil {
		panic(err)
	}
	return data
}

//MapString converts `any` to map[string]string.<将“any”转换为map[string]string。>
func MapString(any interface{}) map[string]string {
	var data = make(map[string]string)
	result, err := json.Marshal(any)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(result, &data)

	if err != nil {
		panic(err)
	}
	return data
}

//MapInt converts `any` to map[int]interface{}.<将“any”map[int]interface{}。>
func MapInt(any interface{}) map[int]interface{} {
	var data = make(map[int]interface{})
	result, err := json.Marshal(any)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(result, &data)

	if err != nil {
		panic(err)
	}
	return data
}
