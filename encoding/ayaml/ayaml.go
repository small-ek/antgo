package ayaml

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

// Encode 将Go对象编码为YAML格式
// Encode converts Go value to YAML format
// 参数 Parameters:
// v - 要编码的Go对象 (Go value to encode)
// 返回 Returns:
// []byte - YAML格式字节切片 (YAML formatted byte slice)
// error - 编码过程中遇到的错误 (Error during encoding)
func Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Decode 将YAML解码到map[string]interface{}
// Decodes YAML into map[string]interface{}
// 参数 Parameters:
// v - YAML格式字节切片 (YAML formatted byte slice)
// 返回 Returns:
// map[string]interface{} - 解码后的字典数据 (Decoded map data)
// error - 解码过程中遇到的错误 (Error during decoding)
func Decode(v []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(v, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeTo 将YAML解码到指定结构体
// Decodes YAML into specified structure
// 参数 Parameters:
// v - YAML格式字节切片 (YAML formatted byte slice)
// result - 目标结构体指针 (Pointer to target structure)
// 返回 Returns:
// error - 解码过程中遇到的错误 (Error during decoding)
func DecodeTo(v []byte, result interface{}) error {
	return yaml.Unmarshal(v, result)
}

// ToJson 将YAML转换为JSON格式
// Converts YAML to JSON format
// 参数 Parameters:
// v - YAML格式字节切片 (YAML formatted byte slice)
// 返回 Returns:
// []byte - JSON格式字节切片 (JSON formatted byte slice)
// error - 转换过程中遇到的错误 (Error during conversion)
func ToJson(v []byte) ([]byte, error) {
	// 直接解析到interface{}避免二次解码
	// Parse directly to interface{} to avoid double decoding
	var data interface{}
	if err := yaml.Unmarshal(v, &data); err != nil {
		return nil, err
	}
	return json.Marshal(data)
}
