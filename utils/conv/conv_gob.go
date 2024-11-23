package conv

import (
	"bytes"
	"encoding/gob"
	"errors"
)

// GobEncoder encodes a struct into a byte slice using gob.
// Suitable for RPC communication or similar struct data transmission.
func GobEncoder(data any) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecoder decodes a byte slice into a struct using gob.
func GobDecoder(data []byte, target any) error {
	if data == nil || target == nil {
		return errors.New("data and target struct cannot be nil")
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(target)
}
