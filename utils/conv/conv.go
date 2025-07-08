package conv

import (
	"encoding/binary"
	"math"
	"strconv"
	"time"
)

// Rune converts any type to rune. Maintains original value for rune types,
// converts to int32 for other types (supports numeric/string/bool compatibility)
// 将任意类型转换为rune类型。保留rune类型的原始值，其他类型转换为int32（支持数字/字符串/布尔型兼容）
func Rune(any interface{}) rune {
	if v, ok := any.(rune); ok {
		return v
	}
	return Int32(any)
}

// Runes converts any type to []rune. Direct conversion for []rune types,
// converts to string first for other types
// 将任意类型转换为[]rune类型。直接转换[]rune类型，其他类型先转换为字符串
func Runes(any interface{}) []rune {
	if v, ok := any.([]rune); ok {
		return v
	}
	return []rune(String(any))
}

// Byte converts any type to byte. Maintains original value for byte types,
// converts to uint8 for other types (supports numeric/string/bool compatibility)
// 将任意类型转换为byte类型。保留byte类型的原始值，其他类型转换为uint8（支持数字/字符串/布尔型兼容）
func Byte(any interface{}) byte {
	if v, ok := any.(byte); ok {
		return v
	}
	return Uint8(any)
}

// Bytes converts any type to []byte with optimized memory allocation.
// Special handling for basic types, JSON serialization for complex types
// 将任意类型转换为[]byte（优化内存分配），基础类型特殊处理，复杂类型使用JSON序列化
func Bytes(any interface{}) []byte {
	if any == nil {
		return nil
	}

	switch v := any.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case int:
		return intToBytes(int64(v))
	case int32:
		return intToBytes(int64(v))
	case int64:
		return intToBytes(v)
	case float32:
		return float32ToBytes(v)
	case float64:
		return float64ToBytes(v)
	default:
		// Fallback to JSON serialization for complex types
		// 复杂类型回退到JSON序列化
		result, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		return result
	}
}

// String converts any type to string with optimized type handling.
// Specialized conversion for basic types, JSON serialization for others
// 将任意类型转换为字符串（优化类型处理），基础类型特殊转换，其他类型使用JSON序列化
func String(any interface{}) string {
	if any == nil {
		return ""
	}

	switch v := any.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case *time.Time:
		if v == nil {
			return ""
		}
		return v.Format("2006-01-02 15:04:05")
	default:
		// JSON serialization for complex types
		// 复杂类型使用JSON序列化
		result, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(result)
	}
}

// intToBytes converts integer to 8-byte slice with zero allocations
// 整型转8字节切片（零内存分配）
func intToBytes(n int64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(n))
	return buf[:]
}

// float32ToBytes converts float32 to 4-byte slice with zero allocations
// float32转4字节切片（零内存分配）
func float32ToBytes(f float32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:]
}

// float64ToBytes converts float64 to 8-byte slice with zero allocations
// float64转8字节切片（零内存分配）
func float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}
