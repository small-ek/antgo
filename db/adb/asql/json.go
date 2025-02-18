package asql

import (
	"bytes"
	"database/sql/driver"
	"errors"
)

type Json []byte

// Value
func (j Json) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

// Scan
func (j *Json) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

// MarshalJSON
func (j Json) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON
func (j *Json) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// IsNull
func (j Json) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

// Equals
func (j Json) Equals(j1 Json) bool {
	return bytes.Equal(j, j1)
}
