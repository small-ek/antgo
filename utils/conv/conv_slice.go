package conv

import (
	"encoding/json"
)

// Ints converts `any` to []int.<将“any”转换为[]int。>
func Ints(any interface{}) []int {
	if any == nil {
		return nil
	}
	var array []int
	switch value := any.(type) {
	case []string:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int(value[i])
		}
	case []int:
		array = value
	case []int8:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []int16:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []int32:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []int64:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []uint:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []uint8:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []uint16:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []uint32:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []uint64:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int(value[i])
		}
	case []bool:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			if value[i] {
				array[i] = 1
			} else {
				array[i] = 0
			}
		}
	case []float32:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int(value[i])
		}
	case []float64:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int(value[i])
		}
	case []interface{}:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int(value[i])
		}
	case [][]byte:
		array = make([]int, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int(value[i])
		}
	case string:
		err := json.Unmarshal(Bytes(value), &array)
		if err != nil {
			panic(err)
		}
	}
	return array
}

// Int32s converts `any` to []int32.<将“any”转换为[]int32。>
func Int32s(any interface{}) []int32 {
	if any == nil {
		return nil
	}
	var array []int32
	switch value := any.(type) {
	case []string:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int32(value[i])
		}
	case []int:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []int8:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []int16:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []int32:
		array = value
	case []int64:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []uint:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []uint8:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []uint16:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []uint32:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []uint64:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int32(value[i])
		}
	case []bool:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int32(value[i])
			if value[i] {
				array[i] = 1
			} else {
				array[i] = 0
			}
		}
	case []float32:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case []float64:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int32(value[i])
		}
	case []interface{}:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int32(value[i])
		}
	case [][]byte:
		array = make([]int32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int32(value[i])
		}
	case string:
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
	}
	return array
}

// Int64s converts `any` to []int64.<将“any”转换为[]int64。>
func Int64s(any interface{}) []int64 {
	if any == nil {
		return nil
	}
	var array []int64
	switch value := any.(type) {
	case []string:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int64(value[i])
		}
	case []int:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []int8:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []int16:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []int32:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []int64:
		array = value
	case []uint:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []uint8:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []uint16:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []uint32:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []uint64:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = int64(value[i])
		}
	case []bool:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			if value[i] {
				array[i] = 1
			} else {
				array[i] = 0
			}
		}
	case []float32:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int64(value[i])
		}
	case []float64:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int64(value[i])
		}
	case []interface{}:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int64(value[i])
		}
	case [][]byte:
		array = make([]int64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Int64(value[i])
		}
	case string:
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
	}

	return array
}

// Strings converts `any` to []string.<将“any”转换为[]string。>
func Strings(any interface{}) []string {
	if any == nil {
		return nil
	}
	var array []string
	switch value := any.(type) {
	case []int:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []int8:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []int16:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []int32:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []int64:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []uint:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []uint8:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []uint16:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []uint32:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []uint64:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []bool:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []float32:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []float64:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []interface{}:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case []string:
		array = value
	case [][]byte:
		array = make([]string, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = String(value[i])
		}
	case string:
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
	}
	return array
}

// Uints converts `any` to []uint.<将“any”转换为[]uint。>
func Uints(any interface{}) []uint {
	if any == nil {
		return nil
	}
	var array []uint
	switch value := any.(type) {
	case []string:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint(value[i])
		}
	case []int8:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []int16:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []int32:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []int64:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []uint:
		array = value
	case []uint8:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []uint16:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []uint32:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []uint64:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint(value[i])
		}
	case []bool:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			if value[i] {
				array[i] = 1
			} else {
				array[i] = 0
			}
		}
	case []float32:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint(value[i])
		}
	case []float64:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint(value[i])
		}
	case []interface{}:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint(value[i])
		}
	case [][]byte:
		array = make([]uint, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint(value[i])
		}
	case string:
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
	}
	return array
}

// Uint32s converts `any` to []uint32.<将“any”转换为[]uint32。>
func Uint32s(any interface{}) []uint32 {
	if any == nil {
		return nil
	}
	var array []uint32
	switch value := any.(type) {
	case []string:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint32(value[i])
		}
	case []int8:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []int16:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []int32:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []int64:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []uint:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []uint8:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []uint16:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []uint32:
		array = value
	case []uint64:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint32(value[i])
		}
	case []bool:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			if value[i] {
				array[i] = 1
			} else {
				array[i] = 0
			}
		}
	case []float32:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint32(value[i])
		}
	case []float64:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint32(value[i])
		}
	case []interface{}:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint32(value[i])
		}
	case [][]byte:
		array = make([]uint32, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint32(value[i])
		}
	case string:
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
	}
	return array
}

// Uint64s converts `any` to []uint64.<将“any”转换为[]uint64。>
func Uint64s(any interface{}) []uint64 {
	if any == nil {
		return nil
	}
	var array []uint64
	switch value := any.(type) {
	case []string:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint64(value[i])
		}
	case []int8:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []int16:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []int32:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []int64:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []uint:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []uint8:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []uint16:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []uint32:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = uint64(value[i])
		}
	case []uint64:
		array = value
	case []bool:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			if value[i] {
				array[i] = 1
			} else {
				array[i] = 0
			}
		}
	case []float32:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint64(value[i])
		}
	case []float64:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint64(value[i])
		}
	case []interface{}:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint64(value[i])
		}
	case [][]byte:
		array = make([]uint64, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = Uint64(value[i])
		}
	case string:
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
	}
	return array
}

// Interfaces converts `any` to []interface{}.<将“any”转换为[]interface{}。>
func Interfaces(any interface{}) []interface{} {
	if any == nil {
		return nil
	}
	var array []interface{}
	switch value := any.(type) {
	case []int:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []int8:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []int16:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []int32:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []int64:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []uint:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []uint8:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []uint16:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []uint32:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []uint64:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []bool:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []float32:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []float64:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case []interface{}:
		array = value
	case []string:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case [][]byte:
		array = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			array[i] = value[i]
		}
	case string:
		err := json.Unmarshal([]byte(value), &array)
		if err != nil {
			panic(err)
		}
	}
	return array
}
