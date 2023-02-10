package abinary

import (
	"bytes"
	"encoding/binary"
)

// Encode
func Encode(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encode
func Decode(b []byte, value interface{}) error {
	return binary.Read(bytes.NewReader(b), binary.LittleEndian, value)
}
