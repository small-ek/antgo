package axml

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"sync"
)

// StringMap 定义字符串映射类型，用于XML与JSON的转换
// StringMap defines a string map type for XML/JSON conversion
type StringMap map[string]string

// MarshalXML 使用Token API优化XML序列化性能（避免反射）
// MarshalXML optimizes XML serialization using Token API (avoids reflection)
func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	// 写入起始标签
	// Write start element
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// 遍历键值对，逐个编码
	// Iterate key-value pairs and encode each
	for k, v := range m {
		elementStart := xml.StartElement{Name: xml.Name{Local: k}}
		if err := e.EncodeToken(elementStart); err != nil {
			return err
		}
		if err := e.EncodeToken(xml.CharData(v)); err != nil {
			return err
		}
		if err := e.EncodeToken(elementStart.End()); err != nil {
			return err
		}
	}

	// 写入父级结束标签并刷新缓冲区
	// Write parent end element and flush
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return e.Flush() // 确保数据写入底层writer | Ensures data is written
}

// UnmarshalXML 使用Token流式解析优化反序列化
// UnmarshalXML optimizes deserialization with token streaming
func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = make(StringMap) // 初始化避免nil map | Initialize to avoid nil map

	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 仅处理开始标签，跳过注释/指令等
		// Process only start elements, skip others like comments
		if se, ok := token.(xml.StartElement); ok {
			var value string
			if err := d.DecodeElement(&value, &se); err != nil {
				return err // 严格错误处理 | Strict error handling
			}
			(*m)[se.Name.Local] = value
		}
	}
	return nil
}

// Encode 直接调用优化后的XML序列化
// Encode directly uses the optimized XML serialization
func Encode(data map[string]string) ([]byte, error) {
	return xml.Marshal(StringMap(data))
}

// Decode 反序列化到新map，确保线程安全
// Decode into a new map for thread safety
func Decode(data []byte) (map[string]string, error) {
	result := make(StringMap)
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeTo 复用现有map空间，清空旧数据
// DecodeTo reuses existing map and clears old data
func DecodeTo(data []byte, result map[string]string) error {
	// 兼容性清空（Go <1.21需手动遍历）
	// Compatibility clear (for Go <1.21)
	for k := range result {
		delete(result, k)
	}

	// 类型转换复用内存
	// Type conversion to reuse memory
	sm := StringMap(result)
	if err := xml.Unmarshal(data, &sm); err != nil {
		return err
	}
	return nil
}

// ToJson 使用对象池减少GC压力
// ToJson uses object pool to reduce GC pressure
var mapPool = sync.Pool{
	New: func() interface{} {
		return make(StringMap)
	},
}

func ToJson(v []byte) ([]byte, error) {
	m := mapPool.Get().(StringMap)
	defer func() {
		// 清空并放回池中
		// Clear and return to pool
		for k := range m {
			delete(m, k)
		}
		mapPool.Put(m)
	}()

	if err := xml.Unmarshal(v, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
