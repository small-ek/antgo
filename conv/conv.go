package conv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

//基本类型转换

//将<i>转换为Rune.
func Rune(i interface{}) rune {
	if v, ok := i.(rune); ok {
		return v
	}
	return rune(Int32(i))
}

//将<i>转换为Runes.
func Runes(i interface{}) []rune {
	if v, ok := i.([]rune); ok {
		return v
	}
	return []rune(String(i))
}

//将<i>转换为Byte.
func Byte(i interface{}) byte {
	if v, ok := i.(byte); ok {
		return v
	}
	return Uint8(i)
}

// 将<i>转换为Bytes.
func Bytes(i interface{}) []byte {
	if i == nil {
		return nil
	}

	switch value := i.(type) {

	case string:
		return []byte(value)
	case []byte:
		return value
	case int:
		return intToBytes(value)
	case int32:
		return intToBytes(value)
	case int64:
		return intToBytes(value)
	case float32:
		return Float32ToBytes(value)
	case float64:
		return Float64ToBytes(value)
	default:
		result, err := json.Marshal(value)
		if err != nil {
			log.Println("类型不正确" + err.Error())
		}
		return result
	}
	return i.([]byte)
}

//将<i>转换为String
func String(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	}
	return i.(string)
}

//将<i>转换为Bool
func Bool(i interface{}) bool {
	if i == nil {
		return false
	}
	switch value := i.(type) {
	case bool:
		return value
	case []byte:
		if strings.ToLower(string(value)) == "false" {
			return false
		}
		return true
	case string:
		if strings.ToLower(value) == "false" {
			return false
		}
		return true
	}

	return false
}

//将<i>转换为Int
func Int(i interface{}) int {
	if i == nil {
		return 0
	}
	if v, ok := i.(int); ok {
		return v
	}
	return int(Int64(i))
}

//将<i>转换为Int8
func Int8(i interface{}) int8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int8); ok {
		return v
	}
	return int8(Int64(i))
}

//将<i>转换为Int16
func Int16(i interface{}) int16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int16); ok {
		return v
	}
	return int16(Int64(i))
}

//将<i>转换为Int32
func Int32(i interface{}) int32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int32); ok {
		return v
	}
	return int32(Int64(i))
}

//将<i>转换为Int64
func Int64(i interface{}) int64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return int64(binary.BigEndian.Uint64(value))
	case string:
		int64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Println("数据类型不正确" + err.Error())
		}
		return int64
	}
	return i.(int64)
}

//将<i>转换为Uint
func Uint(i interface{}) uint {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint); ok {
		return v
	}
	return uint(Uint64(i))
}

//将<i>转换为Uint8
func Uint8(i interface{}) uint8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint8); ok {
		return v
	}
	return uint8(Uint64(i))
}

//将<i>转换为Uint16
func Uint16(i interface{}) uint16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint16); ok {
		return v
	}
	return uint16(Uint64(i))
}

//将<i>转换为Uint32
func Uint32(i interface{}) uint32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint32); ok {
		return v
	}
	return uint32(Uint64(i))
}

//将<i>转换为Uint64
func Uint64(i interface{}) uint64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
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
	case uint64:
		return value
	case float32:
		return uint64(value)
	case float64:
		return uint64(value)
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return binary.LittleEndian.Uint64(value)

	}
	return i.(uint64)
}

// 将<i>转换为Float32
func Float32(i interface{}) float32 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case []byte:
		return math.Float32frombits(binary.LittleEndian.Uint32(value))

	default:
		v, _ := strconv.ParseFloat(String(i), 64)
		return float32(v)
	}
}

// 将<i>转换为Float64
func Float64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case []byte:
		return math.Float64frombits(binary.LittleEndian.Uint64(value))
	default:
		v, _ := strconv.ParseFloat(String(i), 64)
		return v
	}
}

// int将<i>转换为Bytes.
func intToBytes(i interface{}) []byte {
	x := Int64(i)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// float64将<i>转换为Bytes.
func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

// float64将<i>转换为Bytes.
func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
