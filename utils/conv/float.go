package conv

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
)

// Float32 converts `any` to float32. 将任意类型的数据转换为 float32。
func Float32(any interface{}) float32 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		// Already a float32, return as is
		return value
	case float64:
		// Convert float64 to float32
		return float32(value)
	case []byte:
		// Convert IEEE 754 encoded bytes to float32
		if len(value) != 4 {
			panic(errors.New("invalid byte length for float32"))
		}
		return math.Float32frombits(binary.LittleEndian.Uint32(value))
	case string:
		// Parse string to float64, then convert to float32
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0
		}
		return float32(result)
	default:
		// Convert other types to string, then parse as float
		v, err := strconv.ParseFloat(String(any), 64)
		if err != nil {
			return 0
		}
		return float32(v)
	}
}

// Float64 converts `any` to float64. 将任意类型的数据转换为 float64。
func Float64(any interface{}) float64 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		// Convert float32 to float64
		return float64(value)
	case float64:
		// Already a float64, return as is
		return value
	case string:
		// Parse string to float64
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0
		}
		return result
	case []byte:
		// Convert IEEE 754 encoded bytes to float64
		if len(value) != 8 {
			panic(errors.New("invalid byte length for float64"))
		}
		return math.Float64frombits(binary.LittleEndian.Uint64(value))
	default:
		// Convert other types to string, then parse as float
		v, err := strconv.ParseFloat(String(any), 64)
		if err != nil {
			return 0
		}
		return v
	}
}
