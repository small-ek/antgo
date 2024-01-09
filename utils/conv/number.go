package conv

import (
	"encoding/binary"
	"strconv"
)

// Int converts "any" to int.<将“any”转换为int。>
func Int(any interface{}) int {
	if any == nil {
		return 0
	}
	if v, ok := any.(int); ok {
		return v
	}
	return int(Int64(any))
}

// Int8 converts "any" to int8.<将“any”转换为int8。>
func Int8(any interface{}) int8 {
	if any == nil {
		return 0
	}
	if v, ok := any.(int8); ok {
		return v
	}
	return int8(Int64(any))
}

// Int16 converts "any" to int16.<将“any”转换为int16。>
func Int16(any interface{}) int16 {
	if any == nil {
		return 0
	}
	if v, ok := any.(int16); ok {
		return v
	}
	return int16(Int64(any))
}

// Int32 converts "any" to int32.<将“any”转换为int32。>
func Int32(any interface{}) int32 {
	if any == nil {
		return 0
	}
	if v, ok := any.(int32); ok {
		return v
	}
	return int32(Int64(any))
}

// Int64 converts "any" to int64.<将“any”转换为int64。>
func Int64(any interface{}) int64 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case int, int8, int16, int32, uint, uint8, uint16, uint32, uint64, float32, float64:
		return int64(any.(float64))
	case int64:
		return value
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return int64(binary.BigEndian.Uint64(value))
	case string:
		str, _ := strconv.ParseInt(value, 10, 64)
		return str
	}
	return any.(int64)
}

// Uint converts "any" to uint.<将“any”转换为uint。>
func Uint(any interface{}) uint {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint); ok {
		return v
	}
	return uint(Uint64(any))
}

// Uint8 converts "any" to uint8.<将“any”转换为uint8。>
func Uint8(any interface{}) uint8 {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint8); ok {
		return v
	}
	return uint8(Uint64(any))
}

// Uint16 converts "any" to uint16.<将“any”转换为uint16。>
func Uint16(any interface{}) uint16 {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint16); ok {
		return v
	}
	return uint16(Uint64(any))
}

// Uint32 converts "any" to uint32.<将“any”转换为uint32。>
func Uint32(any interface{}) uint32 {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint32); ok {
		return v
	}
	return uint32(Uint64(any))
}

// Uint64 converts "any" to uint64.<将“any”转换为uint64。>
func Uint64(any interface{}) uint64 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case int:
		return uint64(value)
	case int8:
		return uint64(value)
	case int16:
		return uint64(value)
	case int32:
		return uint64(value)
	case int64:
		return uint64(value)
	case uint:
		return uint64(value)
	case uint8:
		return uint64(value)
	case uint16:
		return uint64(value)
	case uint32:
		return uint64(value)
	case float32:
		return uint64(value)
	case float64:
		return uint64(value)
	case uint64:
		return value
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return binary.LittleEndian.Uint64(value)
	case string:
		intNum, _ := strconv.Atoi(value)
		return uint64(intNum)
	}
	return any.(uint64)
}
