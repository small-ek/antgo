package acharset

import (
	"bytes"
	"errors"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io"
	"strings"
	"sync"
)

// charsetAliases 定义了字符集的别名映射，将常见别名映射到标准IANA名称
// charsetAliases defines a mapping of common aliases to standard IANA names
var charsetAliases = map[string]string{
	"hzgb2312": "HZ-GB-2312", // HZ-GB-2312 编码的常见别名
	"gb2312":   "HZ-GB-2312", // Common aliases for HZ-GB-2312 encoding
}

// encodingCache 用于缓存已查询的编码器，提升并发访问性能
// encodingCache caches queried encodings for better concurrent performance
var encodingCache sync.Map

// Decode 将指定字符集的字符串解码为UTF-8字节序列
// 参数：
//
//	value:  需要解码的原始字符串
//	charset: 源字符集名称
//
// 返回：
//
//	[]byte: 解码后的UTF-8字节序列
//	error:  解码过程中遇到的错误
//
// Decode converts a string from specified charset to UTF-8 byte sequence
// Parameters:
//
//	value:  The raw string to be decoded
//	charset: Source charset name
//
// Returns:
//
//	[]byte: Decoded UTF-8 byte sequence
//	error:  Any decoding error encountered
func Decode(value string, charset string) ([]byte, error) {
	enc := getEncoding(charset)
	if enc == nil {
		return nil, errors.New("unsupported charset")
	}

	// 使用转换流进行解码
	// Using transform stream for decoding
	reader := transform.NewReader(
		bytes.NewReader([]byte(value)),
		enc.NewDecoder(),
	)

	// 读取解码后的全部数据
	// Reading all decoded data
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// getEncoding 获取指定字符集的编码器，优先使用缓存
// 参数：
//
//	charset: 字符集名称（支持别名）
//
// 返回：
//
//	encoding.Encoding: 字符集编码器对象
//
// getEncoding gets the encoding for specified charset, using cache first
// Parameters:
//
//	charset: Charset name (supports aliases)
//
// Returns:
//
//	encoding.Encoding: Charset encoder object
func getEncoding(charset string) encoding.Encoding {
	// 统一转换为小写进行别名匹配
	// Convert to lowercase for alias matching
	lowerCharset := strings.ToLower(charset)

	// 检查并替换已知别名
	// Check and replace known aliases
	if alias, ok := charsetAliases[lowerCharset]; ok {
		charset = alias
	}

	// 首先尝试从缓存获取编码器
	// Try getting encoding from cache first
	if enc, ok := encodingCache.Load(charset); ok {
		return enc.(encoding.Encoding)
	}

	// 从IANA官方注册表查询编码
	// Query encoding from IANA registry
	enc, err := ianaindex.MIB.Encoding(charset)
	if err != nil {
		return nil
	}

	// 将查询结果存入缓存
	// Store result in cache
	encodingCache.Store(charset, enc)
	return enc
}
