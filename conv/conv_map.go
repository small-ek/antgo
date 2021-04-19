package conv

import (
	"encoding/json"
	"github.com/small-ek/antgo/os/logs"
)

//Map converts <i> to map[string]interface{}.
func Map(i interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	result, err := json.Marshal(i)
	if err != nil {
		logs.Error(err.Error())
	}

	err = json.Unmarshal(result, &data)
	if err != nil {
		logs.Error(err.Error())
	}
	return data
}

//MapString converts <i> to map[string]string.
func MapString(i interface{}) map[string]string {
	var data = make(map[string]string)
	result, err := json.Marshal(i)

	if err != nil {
		logs.Error(err.Error())
	}

	err = json.Unmarshal(result, &data)

	if err != nil {
		logs.Error(err.Error())
	}
	return data
}

//MapInt converts <i> to map[int]interface{}.
func MapInt(i interface{}) map[int]interface{} {
	var data = make(map[int]interface{})
	result, err := json.Marshal(i)

	if err != nil {
		logs.Error(err.Error())
	}

	err = json.Unmarshal(result, &data)

	if err != nil {
		logs.Error(err.Error())
	}
	return data
}
