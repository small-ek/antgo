package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"testing"
)

func TestBool(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{nil, false},               // nil should return false
		{true, true},               // true should return true
		{false, false},             // false should return false
		{"true", true},             // "true" should return true
		{"false", false},           // "false" should return false
		{"anything", true},         // non-"false" string should return true
		{[]byte("true"), true},     // []byte("true") should return true
		{[]byte("false"), false},   // []byte("false") should return false
		{[]byte("anything"), true}, // non-"false" []byte should return true
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := conv.Bool(tt.input)
			if result != tt.expected {
				t.Errorf("Bool(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
