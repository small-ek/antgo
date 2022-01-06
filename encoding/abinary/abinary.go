package abinary

import (
	"bytes"
	"encoding/binary"
)

//Encode
func Encode(value interface{}) *bytes.Buffer {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		panic(err)
	}
	return buf
}

//Encode
func Decode(b []byte, value interface{}) error {
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, value)
	return err
}
