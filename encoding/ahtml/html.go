package ahtml

import (
	strip "github.com/grokify/html-strip-tags-go"
	"html"
	"strings"
)

var (
	// specialCharsReplacer 预定义HTML特殊字符替换器，提升替换性能
	// Pre-defined HTML special characters replacer for better performance
	specialCharsReplacer = strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&#34;", // 34 = 十六进制的双引号
		"'", "&#39;", // 39 = 十六进制的单引号
	)

	// specialCharsDecodeReplacer 预定义HTML实体反向替换器
	// Pre-defined HTML entity decoder for better performance
	specialCharsDecodeReplacer = strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&#34;", `"`,
		"&#39;", "'",
	)
)

// StripTags 过滤掉HTML标签，只返回text内容
// 注意：使用第三方库实现，性能取决于底层库效率
// 参考：http://php.net/manual/zh/function.strip-tags.php
//
// StripTags removes all HTML tags and returns text content
// Note: Uses third-party library implementation, performance depends on underlying library
// Reference: http://php.net/manual/en/function.strip-tags.php
func StripTags(s string) string {
	return strip.StripTags(s)
}

// Entities 转换所有HTML特殊字符为对应实体（包括引号）
// 注意：与PHP的htmlentities()不同，当前实现仅转义5个基础字符
// 参考：http://php.net/manual/zh/function.htmlentities.php
//
// Entities converts all HTML special characters to entities (including quotes)
// Note: Different from PHP's htmlentities(), current implementation escapes 5 basic characters
// Reference: http://php.net/manual/en/function.htmlentities.php
func Entities(s string) string {
	return html.EscapeString(s)
}

// EntitiesDecode 将HTML实体转换回普通字符
// 参考：http://php.net/manual/zh/function.html-entity-decode.php
//
// EntitiesDecode converts HTML entities back to normal characters
// Reference: http://php.net/manual/en/function.html-entity-decode.php
func EntitiesDecode(s string) string {
	return html.UnescapeString(s)
}

// SpecialChars 转换关键HTML特殊字符为实体
// 特性：
// - 自动处理 &, <, >, ", ' 五个字符
// - 使用预编译的替换器，性能优于动态编译
// - 单引号强制转换（类似PHP的ENT_QUOTES模式）
// 参考：http://php.net/manual/zh/function.htmlspecialchars.php
//
// SpecialChars converts key HTML special characters to entities
// Features:
// - Handles &, <, >, ", ' automatically
// - Uses pre-compiled replacer for better performance
// - Forces single quote conversion (similar to PHP's ENT_QUOTES mode)
// Reference: http://php.net/manual/en/function.htmlspecialchars.php
func SpecialChars(s string) string {
	return specialCharsReplacer.Replace(s)
}

// SpecialCharsDecode 反转SpecialChars的转换操作
// 特性：
// - 使用预编译替换器，高性能解码
// - 完全匹配SpecialChars的编码规则
// 参考：http://php.net/manual/zh/function.htmlspecialchars-decode.php
//
// SpecialCharsDecode reverses SpecialChars conversion
// Features:
// - Uses pre-compiled replacer for high-performance decoding
// - Exactly matches SpecialChars encoding rules
// Reference: http://php.net/manual/en/function.htmlspecialchars-decode.php
func SpecialCharsDecode(s string) string {
	return specialCharsDecodeReplacer.Replace(s)
}
