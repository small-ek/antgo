package conv

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
)

//Struct conversion binding
//model Bound model
//data Data
func Struct(model interface{}, data interface{}) {
	result, err := json.Marshal(data)
	err = json.Unmarshal(result, model)
	if err != nil {
		log.Println(err.Error())
	}
}

//StructToBytes Use gob encoding generally used for "similar" two structure transmission binding or RPC communication
func StructToBytes(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		log.Println(err.Error())
	}
	return buf.Bytes()
}

//BytesToStruct Using gob encoding is generally used for "similar" two structure transmission binding or RPC communication
func BytesToStruct(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
