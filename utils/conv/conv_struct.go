package conv

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
)

// ToStruct converts data from one struct to another using JSON serialization/deserialization. 使用 JSON 序列化-反序列化将数据从一个结构转换为另一个结构
// model: Target struct to bind data into.
// data: Source data to convert. Can be a struct, map, or slice.
// Returns an error if the conversion fails.
func ToStruct(data any, model any) error {
	if model == nil || data == nil {
		return errors.New("model and data cannot be nil")
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// JSON serialization/deserialization for conversion
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(result, model)
}

// UnmarshalJSON Using JSON deserialization. 使用 JSON反序列化
// model: Target struct to bind data into.
// data: Source data to convert. Can be a struct, map, or slice.
// Returns an error if the conversion fails.
func UnmarshalJSON(data []byte, model any) error {
	if model == nil || data == nil {
		return errors.New("data and model cannot be nil")
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	return json.Unmarshal(data, model)
}

// ToJSON Using JSON deserialization. 使用 JSON序列化字符串
// model: Target struct to bind data into.
// data: Source data to convert. Can be a struct, map, or slice.
func ToJSON(data any) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	return json.Marshal(data)
}
