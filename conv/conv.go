package conv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/small-ek/antgo/os/logs"
	"math"
	"strconv"
	"strings"
	"time"
)

//Basic type conversion

//Rune Convert <i> to Rune.
func Rune(i interface{}) rune {
	if v, ok := i.(rune); ok {
		return v
	}
	return Int32(i)
}

//Runes Convert <i> to []rune.
func Runes(i interface{}) []rune {
	if v, ok := i.([]rune); ok {
		return v
	}
	return []rune(String(i))
}

//Byte Convert <i> to byte.
func Byte(i interface{}) byte {
	if v, ok := i.(byte); ok {
		return v
	}
	return Uint8(i)
}

//Bytes Convert <i> to []byte.
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
			logs.Error(err.Error())
		}
		return result
	}
}

//String Convert <i> to String.
func String(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int8, int16, int32, int:
		return strconv.Itoa(value.(int))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(value.(uint64), 10)
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
	default:
		result, err := json.Marshal(value)
		if err != nil {
			logs.Error(err.Error())
		}
		return string(result)
	}
	return i.(string)
}

//Bool converts <i> to Bool.
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

//Int converts <i> to int.
func Int(i interface{}) int {
	if i == nil {
		return 0
	}
	if v, ok := i.(int); ok {
		return v
	}
	return int(Int64(i))
}

//Int8 converts <i> to int8.
func Int8(i interface{}) int8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int8); ok {
		return v
	}
	return int8(Int64(i))
}

//Int16 converts <i> to int16.
func Int16(i interface{}) int16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int16); ok {
		return v
	}
	return int16(Int64(i))
}

//Int32 converts <i> to int32.
func Int32(i interface{}) int32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int32); ok {
		return v
	}
	return int32(Int64(i))
}

//Int64 converts <i> to int64.
func Int64(i interface{}) int64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int, int8, int16, int32, uint, uint8, uint16, uint32, uint64, float32, float64:
		return int64(i.(float64))
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
		str, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logs.Error(err.Error())
		}
		return str
	}
	return i.(int64)
}

//Uint converts <i> to uint.
func Uint(i interface{}) uint {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint); ok {
		return v
	}
	return uint(Uint64(i))
}

//Uint8 converts <i> to uint8.
func Uint8(i interface{}) uint8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint8); ok {
		return v
	}
	return uint8(Uint64(i))
}

//Uint16 converts <i> to uint16.
func Uint16(i interface{}) uint16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint16); ok {
		return v
	}
	return uint16(Uint64(i))
}

//Uint32 converts <i> to uint32.
func Uint32(i interface{}) uint32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint32); ok {
		return v
	}
	return uint32(Uint64(i))
}

//Uint64 converts <i> to uint64.
func Uint64(i interface{}) uint64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64:
		return value.(uint64)
	case uint64:
		return value
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

//Float32 converts <i> to float32.
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
	case string:
		var result, err = strconv.ParseFloat(value, 32/64)
		if err != nil {
			logs.Error(err.Error())
		}
		return Float32(result)
	default:
		v, _ := strconv.ParseFloat(String(i), 64)
		return float32(v)
	}
}

//Float64 converts <i> to float64.
func Float64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case string:
		var result, err = strconv.ParseFloat(value, 32/64)
		if err != nil {
			logs.Error(err.Error())
		}
		return result
	case []byte:
		return math.Float64frombits(binary.LittleEndian.Uint64(value))
	default:
		v, _ := strconv.ParseFloat(String(i), 64)
		return v
	}
}

//intToBytes converts <i> to []byte.
func intToBytes(i interface{}) []byte {
	x := Int64(i)
	bytesBuffer := bytes.NewBuffer([]byte{})
	var err = binary.Write(bytesBuffer, binary.BigEndian, x)
	if err != nil {
		logs.Error(err.Error())
	}
	return bytesBuffer.Bytes()
}

//Float32ToBytes converts <i> to []byte.
func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, bits)
	return result
}

//Float64ToBytes converts <i> to []byte.
func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, bits)
	return result
}
