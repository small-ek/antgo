package atoml

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"sync"
)

// bufferPool 用于复用 bytes.Buffer 实例，减少内存分配
// bufferPool is used to reuse bytes.Buffer instances, reducing memory allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Encode 将结构体编码为 TOML 格式的字节数据
// Encode encodes a struct into TOML-formatted bytes
// 使用 sync.Pool 复用 bytes.Buffer 以减少内存分配，提升性能
// Uses sync.Pool to reuse bytes.Buffer, reducing memory allocations and improving performance
func Encode(v interface{}) ([]byte, error) {
	// 从对象池获取 buffer
	// Get buffer from pool
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()         // 清空缓冲区以便复用 / Clear buffer for reuse
		bufferPool.Put(buffer) // 放回对象池 / Return to pool
	}()

	// 执行 TOML 编码
	// Perform TOML encoding
	if err := toml.NewEncoder(buffer).Encode(v); err != nil {
		return nil, err
	}

	// 创建新切片并拷贝数据以避免缓冲区复用导致的数据污染
	// Create new slice and copy data to prevent buffer reuse contamination
	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

// Decode 将 TOML 字节数据解码到通用 map 结构
// Decode decodes TOML bytes into a generic map structure
// 注意：高频场景建议使用 DecodeTo 直接解析到具体结构体以获得更好性能
// Note: For high-frequency scenarios, use DecodeTo to parse directly into concrete structs for better performance
func Decode(v []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := toml.Unmarshal(v, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeTo 将 TOML 字节数据直接解析到目标结构体
// DecodeTo decodes TOML bytes directly into the target struct
// 直接解析到结构体比使用通用 map 性能更好，推荐在明确数据结构时使用
// Direct decoding into structs performs better than using generic maps, recommended when data structure is known
func DecodeTo(v []byte, result interface{}) error {
	return toml.Unmarshal(v, result)
}

// jsonConvertPool 用于复用 map 实例，优化 ToJson 性能
// jsonConvertPool is used to reuse map instances to optimize ToJson performance
var jsonConvertPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{})
	},
}

// ToJson 将 TOML 字节数据转换为 JSON 格式
// ToJson converts TOML bytes to JSON format
// 通过复用 map 实例和避免中间数据结构提升转换性能
// Improves conversion performance by reusing map instances and avoiding intermediate data structures
func ToJson(v []byte) ([]byte, error) {
	// 从对象池获取 map 实例
	// Get map instance from pool
	result := jsonConvertPool.Get().(map[string]interface{})
	defer func() {
		// 清空 map 以便复用 / Clear map for reuse
		for key := range result {
			delete(result, key)
		}
		jsonConvertPool.Put(result) // 放回对象池 / Return to pool
	}()

	// 解析 TOML 数据到复用的 map
	// Parse TOML data into reused map
	if err := toml.Unmarshal(v, &result); err != nil {
		return nil, err
	}

	// 将 map 编码为 JSON
	// Encode map to JSON
	return json.Marshal(result)
}
