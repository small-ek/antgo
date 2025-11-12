package asql

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"unicode/utf8"
)

type Json []byte

func (j Json) IsNull() bool {
	b := bytes.TrimSpace(j)
	return len(b) == 0 || bytes.Equal(b, []byte("null"))
}

func (j Json) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	// 保证合法 UTF-8（会去掉/替换非法字节）
	if !utf8.Valid(j) {
		j = bytes.ToValidUTF8(j, []byte{})
	}
	return string(j), nil // 返回 string，避免 MySQL 报 3144
}

func (j *Json) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
		return nil
	case string:
		*j = append((*j)[0:0], v...)
		return nil
	default:
		return fmt.Errorf("asql.Json: unsupported scan type %T", value)
	}
}

func (j Json) MarshalJSON() ([]byte, error) {
	if j == nil || j.IsNull() {
		return []byte("null"), nil
	}
	// 假设 j 已经是合法的 JSON bytes（如果不放心，可以在这里校验）
	return j, nil
}

func (j *Json) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("asql.Json: nil pointer")
	}
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*j = nil
		return nil
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j Json) Equals(j1 Json) bool {
	return bytes.Equal(j, j1)
}
