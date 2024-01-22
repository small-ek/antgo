package conv

import (
	"encoding/json"
)

// Map converts `any` to map[string]interface{}.<将“any”转换为map[string]interface{}。>
func Map(any interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	if err := json.Unmarshal(Bytes(any), &data); err != nil {
		panic(err)
	}
	return data
}

// Maps converts `any` to []map[string]interface{}.<将“any”转换为[]map[string]interface{}。>
func Maps(any interface{}) []map[string]interface{} {
	var data = []map[string]interface{}{}
	if err := json.Unmarshal(Bytes(any), &data); err != nil {
		panic(err)
	}
	return data
}

// MapString converts `any` to map[string]string.<将“any”转换为map[string]string。>
func MapString(any interface{}) map[string]string {
	var data = make(map[string]string)
	if err := json.Unmarshal(Bytes(any), &data); err != nil {
		panic(err)
	}
	return data
}

// MapInt converts `any` to map[int]interface{}.<将“any”map[int]interface{}。>
func MapInt(any interface{}) map[int]interface{} {
	var data = make(map[int]interface{})
	if err := json.Unmarshal(Bytes(any), &data); err != nil {
		panic(err)
	}
	return data
}
