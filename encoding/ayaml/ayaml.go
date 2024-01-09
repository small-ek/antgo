package ayaml

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

// Encode
func Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Decode
func Decode(v []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(v, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeTo
func DecodeTo(v []byte, result interface{}) error {
	return yaml.Unmarshal(v, result)
}

// ToJson
func ToJson(v []byte) ([]byte, error) {
	if r, err := Decode(v); err != nil {
		return nil, err
	} else {
		return json.Marshal(r)
	}
}
