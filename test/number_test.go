package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"math"
	"testing"
)

func TestInt(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  int
	}{
		{"nil", nil, 0},
		{"int", 42, 42},
		{"int8", int8(8), 8},
		{"int16", int16(16), 16},
		{"int32", int32(32), 32},
		{"int64", int64(64), 64},
		{"uint", uint(42), 42},
		{"uint8", uint8(8), 8},
		{"uint16", uint16(16), 16},
		{"uint32", uint32(32), 32},
		{"uint64", uint64(64), 64},
		{"float32", float32(3.14), 3},
		{"float64", float64(2.718), 2},
		{"string", "123", 123},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"[]byte", []byte{0, 0, 0, 0, 0, 0, 0, 42}, 42},
		{"invalid string", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conv.Int(tt.input); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt8(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  int8
	}{
		{"nil", nil, 0},
		{"int8", int8(8), 8},
		{"int16 overflow", int16(128), 127}, // 超过int8范围
		{"int16", int16(16), 16},
		{"int32", int32(32), 32},
		{"int64", int64(64), 64},
		{"uint8", uint8(8), 8},
		{"uint16", uint16(16), 16},
		{"uint32", uint32(32), 32},
		{"uint64", uint64(64), 64},
		{"float32", float32(3.14), 3},
		{"float64", float64(2.718), 2},
		{"string", "123", 123},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"[]byte", []byte{0, 0, 0, 0, 0, 0, 0, 42}, 42},
		{"invalid string", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conv.Int8(tt.input); got != tt.want {
				t.Errorf("Int8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  int64
	}{
		{"nil", nil, 0},
		{"int64", int64(64), 64},
		{"int64 max", int64(math.MaxInt64), math.MaxInt64},
		{"int64 min", int64(math.MinInt64), math.MinInt64},
		{"uint64 overflow", uint64(math.MaxUint64), math.MaxInt64}, // 超过int64范围
		{"float64", float64(2.718), 2},
		{"float64 overflow", float64(1e20), math.MaxInt64}, // 超过int64范围
		{"string", "123", 123},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"[]byte", []byte{0, 0, 0, 0, 0, 0, 0, 42}, 42},
		{"invalid string", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conv.Int64(tt.input); got != tt.want {
				t.Errorf("Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  uint
	}{
		{"nil", nil, 0},
		{"uint", uint(42), 42},
		{"uint8", uint8(8), 8},
		{"uint16", uint16(16), 16},
		{"uint32", uint32(32), 32},
		{"uint64", uint64(64), 64},
		{"int64 overflow", int64(-1), 0}, // 负数转uint
		{"float32", float32(3.14), 3},
		{"float64", float64(2.718), 2},
		{"string", "123", 123},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"[]byte", []byte{0, 0, 0, 0, 0, 0, 0, 42}, 42},
		{"invalid string", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conv.Uint(tt.input); got != tt.want {
				t.Errorf("Uint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  uint64
	}{
		{"nil", nil, 0},
		{"uint64", uint64(64), 64},
		{"uint64 max", uint64(math.MaxUint64), math.MaxUint64},
		{"int64 overflow", int64(-1), 0}, // 负数转uint64
		{"float64", float64(2.718), 2},
		{"float64 overflow", float64(1e20), math.MaxUint64}, // 超过uint64范围
		{"string", "123", 123},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"[]byte", []byte{0, 0, 0, 0, 0, 0, 0, 42}, 42},
		{"invalid string", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conv.Uint64(tt.input); got != tt.want {
				t.Errorf("Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
