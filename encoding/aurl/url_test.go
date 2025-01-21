package aurl

import (
	"net/url"
	"testing"
)

// 单元测试部分
func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"空格处理", "hello world", "hello+world"},
		{"特殊符号", "a&b/c%", "a%26b%2Fc%25"},
		{"中文测试", "你好", "%E4%BD%A0%E5%A5%BD"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 编码测试
			encoded := Encode(tt.input)
			if encoded != tt.expected {
				t.Errorf("Encode() got = %v, want %v", encoded, tt.expected)
			}

			// 解码测试
			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decode() unexpected error: %v", err)
			}
			if decoded != tt.input {
				t.Errorf("Decode() got = %v, want %v", decoded, tt.input)
			}
		})
	}
}

func TestRawEncodeDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"空格处理", "hello world", "hello%20world"},
		{"保留字符", "~hello", "~hello"},
		{"混合字符", "a b~c", "a%20b~c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := RawEncode(tt.input)
			if encoded != tt.expected {
				t.Errorf("RawEncode() got = %v, want %v", encoded, tt.expected)
			}

			decoded, err := RawDecode(encoded)
			if err != nil {
				t.Fatalf("RawDecode() unexpected error: %v", err)
			}
			if decoded != tt.input {
				t.Errorf("RawDecode() got = %v, want %v", decoded, tt.input)
			}
		})
	}
}

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    url.Values
		expected string
	}{
		{"简单参数", url.Values{"a": {"1"}, "b": {"2"}}, "a=1&b=2"},
		{"特殊字符", url.Values{"key": {"value with space"}}, "key=value+with+space"},
		{"多值参数", url.Values{"a": {"1", "2"}}, "a=1&a=2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildQuery(tt.input)
			if result != tt.expected {
				t.Errorf("BuildQuery() got = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	testURL := "https://user:pass@example.com:8080/path?query=param#fragment"

	tests := []struct {
		name       string
		component  int
		expected   map[string]string
		unexpected []string
		input      string
	}{
		{
			name:      "全部组件",
			component: -1,
			expected: map[string]string{
				"scheme":   "https",
				"host":     "example.com",
				"port":     "8080",
				"user":     "user",
				"pass":     "pass",
				"path":     "/path",
				"query":    "query=param",
				"fragment": "fragment",
			},
		},
		{
			name:      "仅host和port",
			component: 2 | 4,
			expected: map[string]string{
				"host": "example.com",
				"port": "8080",
			},
			unexpected: []string{"scheme", "user"},
		},
		{
			name:      "IPv6地址解析",
			component: 2 | 4,
			input:     "http://[::1]:8080/path",
			expected: map[string]string{
				"host": "::1",
				"port": "8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := testURL
			if tt.input != "" {
				input = tt.input
			}

			result, err := ParseURL(input, tt.component)
			if err != nil {
				t.Fatalf("ParseURL() unexpected error: %v", err)
			}

			// 验证预期存在的字段
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("字段 %s 值不符，got = %v, want %v", k, result[k], v)
				}
			}

			// 验证不应存在的字段
			for _, k := range tt.unexpected {
				if _, exists := result[k]; exists {
					t.Errorf("不应存在字段 %s", k)
				}
			}
		})
	}
}

// 性能测试部分
func BenchmarkRawEncode(b *testing.B) {
	testString := "hello world with special ~ characters and spaces"
	b.Run("Optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RawEncode(testString)
		}
	})
}

func BenchmarkParseURL(b *testing.B) {
	testURL := "https://user:pass@example.com:8080/path/to/resource?query=param&another=value#fragment"
	b.Run("WithHostPort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ParseURL(testURL, 2|4)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("AllComponents", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ParseURL(testURL, -1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// 错误处理测试
func TestDecodeErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"无效编码", "%zz", true},
		{"不完整编码", "a%", true},
		{"有效编码", "a%20b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("Decode() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestParseInvalidURL(t *testing.T) {
	_, err := ParseURL(":invalid-url", -1)
	if err == nil {
		t.Error("预期解析错误但未返回")
	}
}
