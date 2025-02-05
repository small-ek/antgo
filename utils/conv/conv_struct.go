package conv

import (
	"errors"
	jsoniter "github.com/json-iterator/go" // JSON iterator package with a new name
)

// jsoniterPkg is the shared instance for JSON serialization/deserialization.
// Using this shared instance prevents redundant initialization.
var jsoniterPkg = jsoniter.ConfigCompatibleWithStandardLibrary

// ToStruct converts data from any type to the provided model using JSON serialization and deserialization.
// 使用 JSON 序列化和反序列化将数据从任何类型转换为提供的模型。
// data: The source data that will be converted. Can be a struct, map, or slice.
//
//	data: 需要转换的源数据，可以是结构体、映射或切片。
//
// model: The target model that will hold the converted data.
//
//	model: 目标模型，将存储转换后的数据。
//
// Returns an error if the conversion fails.
// 如果转换失败，返回错误。
func ToStruct(data any, model any) error {
	if data == nil || model == nil {
		return errors.New("data and model cannot be nil") // 错误：数据和模型不能为空
	}
	// Serialize the data to JSON and then deserialize it into the model
	// 将数据序列化为 JSON，然后反序列化为目标模型
	result, err := jsoniterPkg.Marshal(data)
	if err != nil {
		return err // 序列化错误，返回错误
	}
	return jsoniterPkg.Unmarshal(result, model) // 反序列化到目标模型
}

// UnmarshalJSON deserializes JSON byte data into the provided model.
// 使用 JSON 反序列化将字节数据转化为指定的模型。
// data: The JSON byte data to be converted into the model.
//
//	data: 需要转换为模型的 JSON 字节数据。
//
// model: The target model to store the deserialized data.
//
//	model: 目标模型，用于存储反序列化后的数据。
//
// Returns an error if deserialization fails.
// 如果反序列化失败，返回错误。
func UnmarshalJSON(data []byte, model any) error {
	if data == nil || model == nil {
		return errors.New("data and model cannot be nil") // 错误：数据和模型不能为空
	}
	return jsoniterPkg.Unmarshal(data, model) // 反序列化 JSON 数据到目标模型
}

// ToJSON serializes the provided data into JSON byte format.
// 将提供的数据序列化为 JSON 字节格式。
// data: The source data to be serialized into JSON.
//
//	data: 需要序列化为 JSON 的源数据。
//
// Returns the JSON byte slice or an error if serialization fails.
// 返回 JSON 字节切片，如果序列化失败，则返回错误。
func ToJSON(data any) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil") // 错误：数据不能为空
	}
	return jsoniterPkg.Marshal(data) // 序列化数据为 JSON 字节
}
