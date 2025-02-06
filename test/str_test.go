package test

import (
	"github.com/small-ek/antgo/utils/str"
	"testing"
)

// TestClearQuotes 测试 ClearQuotes 方法
func TestClearQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello`, "hello"},
		{`hello"`, "hello"},
		{`hello`, "hello"},
		{`""`, ""},
		{`"`, ""},
	}

	for _, tt := range tests {
		actual := str.ClearQuotes(tt.input)
		if actual != tt.expected {
			t.Errorf("ClearQuotes(%q) = %q, want %q", tt.input, actual, tt.expected)
		}
	}
}

// TestUcFirst 测试 UcFirst 方法
func TestUcFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"", ""},
		// 包含 Unicode 字符测试（例如带变音符号的小写字母）
		{"àbc", "Àbc"},
	}

	for _, tt := range tests {
		actual := str.UcFirst(tt.input)
		if actual != tt.expected {
			t.Errorf("UcFirst(%q) = %q, want %q", tt.input, actual, tt.expected)
		}
	}
}

// TestReplaceByMap 测试 ReplaceByMap 方法
func TestReplaceByMap(t *testing.T) {
	origin := "The quick brown fox jumps over the lazy dog."
	replaces := map[string]string{
		"quick": "slow",
		"brown": "white",
		"dog":   "cat",
	}
	expected := "The slow white fox jumps over the lazy cat."
	actual := str.ReplaceByMap(origin, replaces)
	if actual != expected {
		t.Errorf("ReplaceByMap(%q) = %q, want %q", origin, actual, expected)
	}
}

// TestRemoveSymbols 测试 RemoveSymbols 方法
func TestRemoveSymbols(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "HelloWorld"},
		{"123-456_789", "123456789"},
		{"@Go#Lang!", "GoLang"},
		// Unicode 测试（保留字母、数字）
		{"你好，世界123", "你好世界123"},
	}

	for _, tt := range tests {
		actual := str.RemoveSymbols(tt.input)
		if actual != tt.expected {
			t.Errorf("RemoveSymbols(%q) = %q, want %q", tt.input, actual, tt.expected)
		}
	}
}

// TestEqualFoldWithoutChars 测试 EqualFoldWithoutChars 方法
func TestEqualFoldWithoutChars(t *testing.T) {
	tests := []struct {
		s1, s2   string
		expected bool
	}{
		{"Hello-World", "hello world", true},
		{"Go_Lang", "go.lang", true},
		{"Test", "Taste", false},
		{"123.456", "123456", true},
	}

	for _, tt := range tests {
		actual := str.EqualFoldWithoutChars(tt.s1, tt.s2)
		if actual != tt.expected {
			t.Errorf("EqualFoldWithoutChars(%q, %q) = %v, want %v", tt.s1, tt.s2, actual, tt.expected)
		}
	}
}

// TestSplitAndTrim 测试 SplitAndTrim 方法
func TestSplitAndTrim(t *testing.T) {
	input := " apple, banana , , orange,  grape "
	expected := []string{"apple", "banana", "orange", "grape"}
	actual := str.SplitAndTrim(input, ",")
	if len(actual) != len(expected) {
		t.Fatalf("SplitAndTrim(%q) returned %v, want %v", input, actual, expected)
	}
	for i, v := range expected {
		if actual[i] != v {
			t.Errorf("SplitAndTrim(%q)[%d] = %q, want %q", input, i, actual[i], v)
		}
	}

	// 测试带额外 characterMask 参数
	input = "##apple##,##banana##"
	expected = []string{"apple", "banana"}
	actual = str.SplitAndTrim(input, ",", "#")
	if len(actual) != len(expected) {
		t.Fatalf("SplitAndTrim with mask returned %v, want %v", actual, expected)
	}
	for i, v := range expected {
		if actual[i] != v {
			t.Errorf("SplitAndTrim with mask(%q)[%d] = %q, want %q", input, i, actual[i], v)
		}
	}
}

// TestTrim 测试 Trim 方法
func TestTrim(t *testing.T) {
	tests := []struct {
		input    string
		mask     string
		expected string
	}{
		// 仅默认空白字符剥离
		{"  hello  ", "", "hello"},
		// 默认+额外字符剥离
		{"***hello***", "*", "hello"},
		// 无需剥离
		{"hello", "", "hello"},
		// 剥离多个字符
		{"\t\n hello \r\f", "", "hello"},
	}

	for _, tt := range tests {
		var actual string
		if tt.mask == "" {
			actual = str.Trim(tt.input)
		} else {
			actual = str.Trim(tt.input, tt.mask)
		}
		if actual != tt.expected {
			t.Errorf("Trim(%q, %q) = %q, want %q", tt.input, tt.mask, actual, tt.expected)
		}
	}
}

// TestFormatCmdKey 测试 FormatCmdKey 方法
func TestFormatCmdKey(t *testing.T) {
	input := "CMD_KEY_TEST"
	expected := "cmd.key.test"
	actual := str.FormatCmdKey(input)
	if actual != expected {
		t.Errorf("FormatCmdKey(%q) = %q, want %q", input, actual, expected)
	}
}

// TestFormatEnvKey 测试 FormatEnvKey 方法
func TestFormatEnvKey(t *testing.T) {
	input := "env.key.test"
	expected := "ENV_KEY_TEST"
	actual := str.FormatEnvKey(input)
	if actual != expected {
		t.Errorf("FormatEnvKey(%q) = %q, want %q", input, actual, expected)
	}
}

// TestStripSlashes 测试 StripSlashes 方法
func TestStripSlashes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// 测试单个反斜杠移除
		{`Hello\ World`, "Hello World"},
		// 连续的反斜杠被移除
		{`C:\\\\\\\\Path\\\\to\\\\file`, "C:Pathtofile"}, // 修改预期值：原来预期 "C:Pathofile" 是错误的
		// 无需处理
		{"NoSlashes", "NoSlashes"},
		// 复杂情况
		{`\\a\\b\\\c`, "abc"},
	}

	for _, tt := range tests {
		actual := str.StripSlashes(tt.input)
		if actual != tt.expected {
			t.Errorf("StripSlashes(%q) = %q, want %q", tt.input, actual, tt.expected)
		}
	}
}

// TestIsNumeric 测试 IsNumeric 方法
func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"-123", true},
		{"123.456", true},
		{"-123.456", true},
		{"123.", false},
		{".456", false},
		{"abc", false},
		{"", false},
		{"12a3", false},
		{"-12-3", false},
	}

	for _, tt := range tests {
		actual := str.IsNumeric(tt.input)
		if actual != tt.expected {
			t.Errorf("IsNumeric(%q) = %v, want %v", tt.input, actual, tt.expected)
		}
	}
}
