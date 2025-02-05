package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"math"
	"testing"
)

// TestFloat32 测试 Float32 函数
func TestFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float32
	}{
		// 正常输入
		{"nil input", nil, 0},
		{"float32 input", float32(3.14), 3.14},
		{"float64 input", float64(3.141592653589793), float32(3.1415927)}, // float64 转 float32 会丢失精度
		{"string input", "3.14", 3.14},
		{"integer string input", "42", 42},
		{"byte slice input", []byte{0xdb, 0x0f, 0x49, 0x40}, float32(3.1415927)}, // 小端序 IEEE 754 格式

		// 边界条件
		{"empty string input", "", 0},
		{"invalid string input", "abc", 0},
		{"large integer string input", "12345678901234567890", float32(12345678901234567890)}, // 精度丢失
		{"small byte slice input", []byte{0x01}, 0},                                           // 长度不足，触发 panic

		// 特殊值
		{"float32 max value", float32(math.MaxFloat32), float32(math.MaxFloat32)},
		{"float32 min value", float32(math.SmallestNonzeroFloat32), float32(math.SmallestNonzeroFloat32)},
		{"float64 min value", float64(math.SmallestNonzeroFloat64), 0}, // 超出 float32 范围
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.name != "small byte slice input" {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			actual := conv.Float32(tt.input)
			if actual != tt.expected {
				t.Errorf("Float32(%v) = %v, expected %v", tt.input, actual, tt.expected)
			}
		})
	}
}

// TestFloat64 测试 Float64 函数
func TestFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		// 正常输入
		{"nil input", nil, 0},
		{"float64 input", float64(3.141592653589793), 3.141592653589793},
		{"string input", "3.141592653589793", 3.141592653589793},
		{"integer string input", "42", 42},
		{"byte slice input", []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}, 3.141592653589793}, // 小端序 IEEE 754 格式

		// 边界条件
		{"empty string input", "", 0},
		{"invalid string input", "abc", 0},
		{"large integer string input", "12345678901234567890", 12345678901234567890},
		{"small byte slice input", []byte{0x01}, 0}, // 长度不足，触发 panic

		// 特殊值
		{"float64 max value", float64(math.MaxFloat64), math.MaxFloat64},
		{"float64 min value", float64(math.SmallestNonzeroFloat64), math.SmallestNonzeroFloat64},
		{"float32 max value", float32(math.MaxFloat32), float64(math.MaxFloat32)},
		{"float32 min value", float32(math.SmallestNonzeroFloat32), float64(math.SmallestNonzeroFloat32)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.name != "small byte slice input" {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			actual := conv.Float64(tt.input)
			if actual != tt.expected {
				t.Errorf("Float64(%v) = %v, expected %v", tt.input, actual, tt.expected)
			}
		})
	}
}
