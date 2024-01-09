package conv

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"strings"
)

// Struct conversion binding
// model Bound model
// data Data
func Struct(model interface{}, data interface{}) {
	result, err := json.Marshal(data)
	err = json.Unmarshal(result, model)
	if err != nil {
		panic(err)
	}
}

// StructToBytes Use gob encoding generally used for "similar" two structure transmission binding or RPC communication
func StructToBytes(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// BytesToStruct Using gob encoding is generally used for "similar" two structure transmission binding or RPC communication
func BytesToStruct(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// InterfaceToStruct interface conversion Struct
func InterfaceToStruct(data interface{}) interface{} {
	if data == nil {
		return nil
	}
	var result = reflect.New(Indirect(reflect.ValueOf(data)).Type()).Interface()

	jsonStr, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	json.NewDecoder(strings.NewReader(string(jsonStr))).Decode(result)
	return result
}

// Indirect returns last value that v points to
func Indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	return v
}
