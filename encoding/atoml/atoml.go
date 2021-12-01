package atoml

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
)

//Encode
func Encode(v interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := toml.NewEncoder(buffer).Encode(v); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

//Decode
func Decode(v []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := toml.Unmarshal(v, &result); err != nil {
		return nil, err
	}
	return result, nil
}

//DecodeTo
func DecodeTo(v []byte, result interface{}) error {
	return toml.Unmarshal(v, result)
}

//ToJson
func ToJson(v []byte) ([]byte, error) {
	if r, err := Decode(v); err != nil {
		return nil, err
	} else {
		return json.Marshal(r)
	}
}
