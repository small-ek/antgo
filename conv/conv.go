package conv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/small-ek/antgo/os/logs"
	"math"
	"strconv"
	"time"
)

//Rune converts `any` to rune.<将“any”转换为rune。>
func Rune(any interface{}) rune {
	if v, ok := any.(rune); ok {
		return v
	}
	return Int32(any)
}

//Runes converts `any` to []rune.<将“any”转换为[]rune。>
func Runes(any interface{}) []rune {
	if v, ok := any.([]rune); ok {
		return v
	}
	return []rune(String(any))
}

//Byte converts `any` to byte.<将“any”转换为byte。>
func Byte(any interface{}) byte {
	if v, ok := any.(byte); ok {
		return v
	}
	return Uint8(any)
}

//Bytes converts `any` to []byte.<将“any”转换为[]byte。>
func Bytes(any interface{}) []byte {
	if any == nil {
		return nil
	}

	switch value := any.(type) {
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
		return float32ToBytes(value)
	case float64:
		return float64ToBytes(value)
	default:
		result, err := json.Marshal(value)
		if err != nil {
			logs.Error(err.Error())
		}
		return result
	}
}

//String converts `any` to string.<将“any”转换为string。>
func String(any interface{}) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
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
	return any.(string)
}

//intToBytes converts `any` to []byte.<将“any”转换为[]byte。>
func intToBytes(any interface{}) []byte {
	x := Int64(any)
	bytesBuffer := bytes.NewBuffer([]byte{})
	var err = binary.Write(bytesBuffer, binary.BigEndian, x)
	if err != nil {
		logs.Error(err.Error())
	}
	return bytesBuffer.Bytes()
}

//float32ToBytes converts `float` to []byte.<将“float”转换为[]byte。>
func float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, bits)
	return result
}

//float64ToBytes converts `float` to []byte.<将“float”转换为[]byte。>
func float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, bits)
	return result
}
