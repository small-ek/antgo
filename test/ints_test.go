package test

import (
	"fmt"
	"github.com/small-ek/antgo/utils/conv"
	"reflect"
	"testing"
)

func TestInts(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []int
	}{
		{[]string{"1", "2", "3"}, []int{1, 2, 3}},     // 测试字符串切片转换
		{[]int{1, 2, 3}, []int{1, 2, 3}},              // 测试 int 切片
		{"[1,2,3]", []int{1, 2, 3}},                   // 测试 JSON 字符串
		{[]interface{}{1, 2.5, true}, []int{1, 2, 1}}, // 测试接口切片
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %v", tt.input), func(t *testing.T) {
			result := conv.Ints(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []string
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},             // 测试字符串切片
		{"[\"a\", \"b\", \"c\"]", []string{"a", "b", "c"}},             // 测试 JSON 字符串
		{[]int{1, 2, 3}, []string{"1", "2", "3"}},                      // 测试 int 切片转换为字符串
		{[]bool{true, false, true}, []string{"true", "false", "true"}}, // 测试布尔切片
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %v", tt.input), func(t *testing.T) {
			result := conv.Strings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInterfaces(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []interface{}
	}{
		{[]int{1, 2, 3}, []interface{}{1, 2, 3}},                         // 测试 int 切片
		{[]string{"a", "b", "c"}, []interface{}{"a", "b", "c"}},          // 测试字符串切片
		{"[1, 2, 3]", []interface{}{1, 2, 3}},                            // 测试 JSON 字符串
		{[]interface{}{1, "test", true}, []interface{}{1, "test", true}}, // 测试接口切片
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %v", tt.input), func(t *testing.T) {
			result := conv.Interfaces(tt.input)

			// 调试：打印实际结果与期望值
			// Debug: Print actual result and expected result
			fmt.Printf("Expected: %v\n", tt.expected)
			fmt.Printf("Result: %v\n", result)

		})
	}
}
