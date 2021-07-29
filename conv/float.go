package conv

import (
	"encoding/binary"
	"math"
	"strconv"
)

//Float32 converts `any` to float32.<将“any”转换为float32。>
func Float32(any interface{}) float32 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case []byte:
		return math.Float32frombits(binary.LittleEndian.Uint32(value))
	case string:
		var result, _ = strconv.ParseFloat(value, 32/64)
		return Float32(result)
	default:
		v, _ := strconv.ParseFloat(String(any), 64)
		return float32(v)
	}
}

//Float64 converts `any` to float64.<将“any”转换为float64。>
func Float64(any interface{}) float64 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case string:
		var result, _ = strconv.ParseFloat(value, 32/64)
		return result
	case []byte:
		return math.Float64frombits(binary.LittleEndian.Uint64(value))
	default:
		v, _ := strconv.ParseFloat(String(any), 64)
		return v
	}
}
