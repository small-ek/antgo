package abinary

import (
	"bytes"
	"encoding/binary"
	"sync"
)

// bufferPool 用于复用 bytes.Buffer 减少内存分配
// bufferPool is used to reuse bytes.Buffer and reduce memory allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Encode 将任意类型序列化为小端格式的字节切片
// 注意：返回的字节切片是独立副本，修改不会影响缓冲区
//
// Encode serializes the given value into a little-endian byte slice
// Note: The returned byte slice is an independent copy, modifications won't affect the buffer
//
// 参数：
//
//	value - 需要序列化的数据（必须是固定大小的类型或包含固定大小类型的结构体）
//
// Parameters:
//
//	value - The data to be encoded (must be fixed-size type or struct containing fixed-size types)
//
// 返回值：
//
//	[]byte - 序列化后的字节切片
//	error  - 编码过程中遇到的错误（如类型不支持）
//
// Returns:
//
//	[]byte - Serialized byte slice
//	error  - Encoding errors (e.g. unsupported type)
func Encode(value interface{}) ([]byte, error) {
	// 从对象池获取缓冲区
	// Get buffer from pool
	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()         // 重置缓冲区以供复用
		bufferPool.Put(buf) // 放回对象池
	}()

	// 序列化数据到缓冲区
	// Serialize data to buffer
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, err
	}

	// 创建独立副本返回，保证数据安全性
	// Create independent copy for return data safety
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// Decode 将小端格式的字节切片反序列化到指定变量
//
// # Decode deserializes little-endian byte slice into the specified variable
//
// 参数：
//
//	b     - 需要反序列化的字节切片
//	value - 目标变量指针（必须与编码时类型匹配）
//
// Parameters:
//
//	b     - Byte slice to be deserialized
//	value - Pointer to destination variable (must match encoded type)
//
// 返回值：
//
//	error - 解码过程中遇到的错误（如数据不完整或类型不匹配）
//
// Returns:
//
//	error - Decoding errors (e.g. incomplete data or type mismatch)
func Decode(b []byte, value interface{}) error {
	// 直接使用字节切片创建 Reader，无需额外内存分配
	// Create reader directly from byte slice without additional allocation
	return binary.Read(bytes.NewReader(b), binary.LittleEndian, value)
}
