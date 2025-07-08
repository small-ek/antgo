package conv

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

// Float32 converts any type to float32.
// Note that this conversion may lose precision for large integers or high-precision decimals.
// 将任意类型转换为 float32。注意，转换大整数或高精度小数时可能丢失精度。
func Float32(any interface{}) float32 {
	if any == nil {
		return 0
	}

	switch value := any.(type) {
	case float32:
		return value
	case float64:
		// Direct conversion from float64 may lose precision
		// 直接从 float64 转换可能会丢失精度
		return float32(value)
	case string:
		// Parse with 32-bit precision, returns 0 on failure
		// 使用 32 位精度解析字符串，解析失败返回 0
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0
		}
		return float32(result)
	case []byte:
		// Requires exactly 4 bytes in little-endian IEEE 754 format
		// Panics if length is invalid (considered programming error)
		// 需要恰好 4 字节小端序 IEEE 754 格式，长度不符会 直接返回0
		if len(value) != 4 {
			return 0
		}
		return math.Float32frombits(binary.LittleEndian.Uint32(value))
	default:
		// Convert to string first, may lose precision for complex types
		// 先转换为字符串，复杂类型可能丢失精度
		v, err := strconv.ParseFloat(fmt.Sprintf("%v", any), 32)
		if err != nil {
			return 0
		}
		return float32(v)
	}
}

// Float64 converts any type to float64.
// Provides higher precision than Float32 but still has limitations with very large numbers.
// 将任意类型转换为 float64。相比 Float32 精度更高，但对极大数值仍有限制。
func Float64(any interface{}) float64 {
	if any == nil {
		return 0
	}

	switch value := any.(type) {
	case float64:
		return value
	case float32:
		// Direct widening conversion preserves original precision
		// 直接扩展转换，保留原始精度
		return float64(value)
	case string:
		// Parse with full 64-bit precision, returns 0 on failure
		// 使用完整的 64 位精度解析，解析失败返回 0
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0
		}
		return result
	case []byte:
		// Requires exactly 8 bytes in little-endian IEEE 754 format
		// Panics if length is invalid (considered programming error)
		// 需要恰好 8 字节小端序 IEEE 754 格式，长度不符会 panic（视为编程错误）
		if len(value) != 8 {
			return 0
		}
		return math.Float64frombits(binary.LittleEndian.Uint64(value))
	default:
		// Convert to string first, may lose precision for non-primitive types
		// 先转换为字符串，非基本类型可能丢失精度
		v, err := strconv.ParseFloat(fmt.Sprintf("%v", any), 64)
		if err != nil {
			return 0
		}
		return v
	}
}
