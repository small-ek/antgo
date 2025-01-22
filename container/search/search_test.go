package search

import (
	"testing"
)

// TestSearch 线性搜索测试用例
// Test cases for linear search
func TestSearch(t *testing.T) {
	// 测试用例表（支持多种类型）
	// Test case table (supports multiple types)
	testCases := []struct {
		name     string // 测试用例名称
		slice    any    // 被搜索的切片
		key      any    // 搜索的键值
		expected int    // 期望结果
	}{
		// 字符串类型测试
		// String type tests
		{name: "StringFound", slice: []string{"A", "B", "C"}, key: "B", expected: 1},
		{name: "StringNotFound", slice: []string{"Go", "Rust"}, key: "C++", expected: -1},
		{name: "EmptyStringSlice", slice: []string{}, key: "test", expected: -1},

		// 整型测试
		// Integer type tests
		{name: "IntFoundFirst", slice: []int{1, 2, 3}, key: 1, expected: 0},
		{name: "IntFoundLast", slice: []int{10, 20, 30}, key: 30, expected: 2},
		{name: "IntNotFound", slice: []int{5, 10, 15}, key: 12, expected: -1},

		// 浮点数测试
		// Floating-point tests
		{name: "Float32Found", slice: []float32{1.1, 2.2, 3.3}, key: float32(2.2), expected: 1},
		{name: "Float64NotFound", slice: []float64{1.5, 2.5, 3.5}, key: 2.500001, expected: -1},

		// 边界测试
		// Edge cases
		{name: "NilSlice", slice: nil, key: "test", expected: -1},
		{name: "SingleElementFound", slice: []rune{'A'}, key: 'A', expected: 0},
		{name: "DuplicateElements", slice: []int{2, 2, 3, 3}, key: 3, expected: 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result int

			// 类型断言执行测试
			// Type assertion to execute test
			switch s := tc.slice.(type) {
			case []string:
				result = Search(s, tc.key.(string))
			case []int:
				result = Search(s, tc.key.(int))
			case []float32:
				result = Search(s, tc.key.(float32))
			case []float64:
				result = Search(s, tc.key.(float64))
			case []rune:
				result = Search(s, tc.key.(rune))
			case nil:
				result = Search([]int(nil), 0) // 测试nil切片
			default:
				t.Fatalf("不支持的切片类型: %T", s)
			}

			if result != tc.expected {
				t.Errorf("预期 %d, 实际 %d (测试用例: %s)",
					tc.expected, result, tc.name)
			}
		})
	}
}

// TestSearchOrdered 有序搜索测试用例
// Test cases for ordered search
func TestSearchOrdered(t *testing.T) {
	testCases := []struct {
		name     string
		slice    any
		key      any
		expected int
	}{
		// 升序排列测试
		// Ascending order tests
		{name: "IntFound", slice: []int{1, 3, 5, 7}, key: 5, expected: 2},
		{name: "IntBeforeFirst", slice: []int{10, 20, 30}, key: 5, expected: -1},
		{name: "IntAfterLast", slice: []int{100, 200}, key: 300, expected: -1},

		// 字符串排序测试
		// String ordering tests
		{name: "StringFound", slice: []string{"apple", "banana", "orange"}, key: "banana", expected: 1},
		{name: "StringCaseSensitive", slice: []string{"Go", "java", "python"}, key: "go", expected: -1},

		// 浮点数精度测试
		// Floating-point precision test
		{name: "FloatExactMatch", slice: []float64{1.1, 2.2, 3.3}, key: 2.2, expected: 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result int

			switch s := tc.slice.(type) {
			case []int:
				result = SearchOrdered(s, tc.key.(int))
			case []string:
				result = SearchOrdered(s, tc.key.(string))
			case []float64:
				result = SearchOrdered(s, tc.key.(float64))
			default:
				t.Fatalf("不支持的切片类型: %T", s)
			}

			if result != tc.expected {
				t.Errorf("预期 %d, 实际 %d (测试用例: %s)",
					tc.expected, result, tc.name)
			}
		})
	}
}

// TestCustomType 测试自定义类型
// Test custom type
func TestCustomType(t *testing.T) {
	type myInt int

	t.Run("CustomIntFound", func(t *testing.T) {
		slice := []myInt{1, 2, 3}
		key := myInt(2)
		result := Search(slice, key)
		if result != 1 {
			t.Errorf("预期 1, 实际 %d", result)
		}
	})

	t.Run("CustomIntOrdered", func(t *testing.T) {
		slice := []myInt{10, 20, 30}
		key := myInt(20)
		result := SearchOrdered(slice, key)
		if result != 1 {
			t.Errorf("预期 1, 实际 %d", result)
		}
	})
}
