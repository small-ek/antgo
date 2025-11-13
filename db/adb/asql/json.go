package asql

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"unicode/utf8"
)

var ValidateUTF8 = false // 如果数据来源不受信任，可在 init 或运行时设置为 true

type Json []byte

// IsNull: 使用 bytes.TrimSpace 避免 string 转换分配，同时处理 "null"、空白等
func (j Json) IsNull() bool {
	b := bytes.TrimSpace(j)
	return len(b) == 0 || bytes.Equal(b, []byte("null"))
}

// Value: 返回 string 以兼容大多数驱动与 MySQL JSON 字段。
// 当 ValidateUTF8 为 true 时，会修复非法 UTF-8（发生复制）。
func (j Json) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	if ValidateUTF8 && !utf8.Valid(j) {
		// 仅在需要时做修复，避免每次都额外开销
		j = bytes.ToValidUTF8(j, nil)
	}
	return string(j), nil
}

// Scan: 正确处理类型，返回错误（你原来忘记返回）
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

// MarshalJSON: 若为 nil 或 "null"，输出 null
func (j Json) MarshalJSON() ([]byte, error) {
	if j == nil || j.IsNull() {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON: 将输入原样写入；对 "null" 进行转换
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

// Equals
func (j Json) Equals(j1 Json) bool {
	return bytes.Equal(j, j1)
}
