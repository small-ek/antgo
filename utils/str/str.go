package str

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// DefaultTrimChars 是默认被 Trim* 方法剥离的字符集合。
// DefaultTrimChars is the set of characters stripped by default in Trim* functions.
var DefaultTrimChars = string([]byte{
	'\t', // 制表符 Tab.
	'\v', // 垂直制表符 Vertical Tab.
	'\n', // 换行符 Newline (line feed).
	'\r', // 回车符 Carriage return.
	'\f', // 换页符 Form feed.
	' ',  // 空格 Space.
	0x00, // NUL 字符.
	0x85, // Next Line (NEL).
	0xA0, // 不换行空格 Non-breaking space.
})

// ClearQuotes 移除字符串首尾的双引号（如果存在）。
// ClearQuotes removes the leading and trailing double quotes from the string, if present.
func ClearQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

// UcFirst 返回将字符串 s 的首字母转换为大写的新字符串。
// UcFirst returns a copy of s with its first character converted to uppercase.
func UcFirst(s string) string {
	if s == "" {
		return s
	}
	// 使用 utf8.DecodeRuneInString 兼容 Unicode 字符
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// ReplaceByMap 根据 replaces 映射中的键值对对 origin 进行替换（不考虑替换顺序，区分大小写）。
// ReplaceByMap returns a copy of origin after performing replacements defined in the map (case-sensitive, unordered).
func ReplaceByMap(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = strings.ReplaceAll(origin, k, v)
	}
	return origin
}

// RemoveSymbols 移除字符串中所有非字母和非数字的字符。
// RemoveSymbols removes all characters from the string except letters and digits.
func RemoveSymbols(s string) string {
	var b strings.Builder
	// 预分配内存以提高性能
	b.Grow(len(s))
	for _, r := range s {
		// 使用 unicode 判断以支持所有 Unicode 字符
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// EqualFoldWithoutChars 忽略 '-'、'_'、'.' 和空格后，比较两个字符串是否相等（不区分大小写）。
// EqualFoldWithoutChars checks if s1 and s2 are equal in a case-insensitive manner after removing symbols like '-', '_', '.', and space.
func EqualFoldWithoutChars(s1, s2 string) bool {
	return strings.EqualFold(RemoveSymbols(s1), RemoveSymbols(s2))
}

// SplitAndTrim 按照 delimiter 分割字符串 str，对每个元素调用 Trim 进行剥离，并忽略剥离后为空的元素。
// SplitAndTrim splits str by delimiter, trims each element (with optional additional mask), and ignores empty results.
func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	parts := strings.Split(str, delimiter)
	results := make([]string, 0, len(parts))
	for _, part := range parts {
		part = Trim(part, characterMask...)
		if part != "" {
			results = append(results, part)
		}
	}
	return results
}

// Trim 去除字符串两端的空白（或其他指定字符）。
// 如果提供了可选参数 characterMask，则同时剥离这些额外的字符。
// Trim trims whitespace (or other specified characters) from both ends of the string.
// If characterMask is provided, those characters will also be trimmed.
func Trim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}

// FormatCmdKey 将字符串 s 格式化为命令键（统一格式为小写，'_' 转为 '.'）。
// FormatCmdKey formats the string s as a command key: converts to lowercase and replaces '_' with '.'.
func FormatCmdKey(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", "."))
}

// FormatEnvKey 将字符串 s 格式化为环境变量键（统一格式为大写，'.' 转为 '_'）。
// FormatEnvKey formats the string s as an environment key: converts to uppercase and replaces '.' with '_'.
func FormatEnvKey(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, ".", "_"))
}

// StripSlashes 移除由 AddSlashes 插入的转义反斜杠。
// StripSlashes un-quotes a quoted string by removing escape backslashes (as added by AddSlashes).
func StripSlashes(s string) string {
	// 为兼容多字节字符，将字符串转换为 rune 切片
	runes := []rune(s)
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' {
			// 如果后面跟着另一个反斜杠，则跳过当前反斜杠，并跳过下一个
			if i+1 < len(runes) && runes[i+1] == '\\' {
				i++ // 跳过下一个反斜杠
			}
			// 不将反斜杠写入输出中
			continue
		}
		b.WriteRune(runes[i])
	}
	return b.String()
}

// IsNumeric 检查字符串 s 是否为数字。
// 注意：形如 "123.456" 的浮点数字符串也会返回 true。
// IsNumeric checks whether the string s represents a numeric value.
// Note: Floating point strings like "123.456" are also considered numeric.
func IsNumeric(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		// 首字符允许负号
		if s[i] == '-' && i == 0 {
			continue
		}
		// 处理小数点：必须不是首尾字符
		if s[i] == '.' {
			if i > 0 && i < len(s)-1 {
				continue
			}
			return false
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
