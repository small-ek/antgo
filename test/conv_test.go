package test

import (
	"fmt"
	"github.com/small-ek/antgo/utils/conv"
	"testing"
	"time"
)

func TestRune(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected rune
	}{
		{rune('a'), 'a'},
		{int(97), 'a'},
		{int8(97), 'a'},
		{int16(97), 'a'},
		{int64(97), 'a'},
		{nil, 0}, // Edge case for nil
	}

	for _, test := range tests {
		t.Run("TestRune", func(t *testing.T) {
			result := conv.Rune(test.input)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestRunes(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []rune
	}{
		{"hello", []rune{'h', 'e', 'l', 'l', 'o'}},
		{[]rune{'h', 'e', 'l', 'l', 'o'}, []rune{'h', 'e', 'l', 'l', 'o'}},
		{123, []rune("123")}, // Integer as string
		{nil, []rune("")},    // Edge case for nil
	}

	for _, test := range tests {
		t.Run("TestRunes", func(t *testing.T) {
			result := conv.Runes(test.input)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
			for i, r := range result {
				if r != test.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i, test.expected[i], r)
				}
			}
		})
	}
}

func TestByte(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected byte
	}{
		{byte(97), 97},
		{int(97), 97},
		{int8(97), 97},
		{int16(97), 97},
		{int64(97), 97},
		{nil, 0}, // Edge case for nil
	}

	for _, test := range tests {
		t.Run("TestByte", func(t *testing.T) {
			result := conv.Byte(test.input)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []byte
	}{
		{"hello", []byte("hello")},
		{123, []byte{0, 0, 0, 123}},            // Example byte slice for int
		{3.14, []byte{0xC0, 0x48, 0xF5, 0x3F}}, // Example byte slice for float64
		{nil, nil},                             // Edge case for nil
	}

	for _, test := range tests {
		t.Run("TestBytes", func(t *testing.T) {
			result := conv.Bytes(test.input)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
			for i, b := range result {
				if b != test.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i, test.expected[i], b)
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	ts := time.Now()
	fmt.Println(conv.String(&ts))
}
