package aurl

import (
	"net"
	"net/url"
)

// Encode escapes the string so it can be safely placed inside a URL query.
// It uses url.QueryEscape which encodes spaces as '+' symbols.
// Encode 转义字符串，以便安全地放在 URL 查询中。
// 使用 url.QueryEscape 实现，空格会被编码为 '+' 符号。
func Encode(str string) string {
	return url.QueryEscape(str)
}

// Decode converts URL-encoded query strings back to original form.
// It handles '+' as spaces and hex-encoded characters like %AB.
// Decode 将 URL 编码的查询字符串转换回原始形式。
// 处理 '+' 为空格，以及 %AB 形式的十六进制编码字符。
func Decode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// RawEncode URL-encodes according to RFC 3986.
// Spaces are encoded as %20 and special characters like ~ are preserved.
// RawEncode 根据 RFC 3986 进行 URL 编码。
// 空格被编码为 %20，特殊字符如 ~ 保留不转义。
func RawEncode(str string) string {
	return url.PathEscape(str)
}

// RawDecode decodes strings encoded with RFC 3986.
// It reverses the transformations done by RawEncode.
// RawDecode 解码 RFC 3986 编码的字符串。
// 还原 RawEncode 的转换操作。
func RawDecode(str string) (string, error) {
	return url.PathUnescape(str)
}

// BuildQuery generates URL-encoded query string from map values.
// Uses url.Values.Encode() which sorts parameters by key.
// BuildQuery 从映射值生成 URL 编码的查询字符串。
// 使用 url.Values.Encode() 实现，参数按键排序。
func BuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

// splitHostPort 拆分主机地址中的主机名和端口
// 优先使用标准库的 net.SplitHostPort 以处理复杂场景（如IPv6）
func splitHostPort(host string) (hostname, port string) {
	var err error
	hostname, port, err = net.SplitHostPort(host)
	if err != nil {
		// 当无端口或解析错误时返回完整主机名
		hostname = host
	}
	return
}

// ParseURL parses a URL and returns specified components.
// Component flags: -1=all, 1=scheme, 2=host, 4=port, 8=user,
// 16=pass, 32=path, 64=query, 128=fragment.
// Optimized to avoid redundant host parsing operations.
// ParseURL 解析 URL 并返回指定组件。
// 组件标志位：-1=全部，1=协议，2=主机，4=端口，8=用户，
// 16=密码，32=路径，64=查询，128=片段。
// 优化避免重复的主机解析操作。
func ParseURL(str string, component int) (map[string]string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}

	// 预分配结果map容量（最多包含8个元素）
	components := make(map[string]string, 8)
	var hostname, port string

	// 合并处理host和port的解析
	if component == -1 || (component&(2|4)) != 0 {
		hostname, port = splitHostPort(u.Host)
	}

	// 根据标志位填充组件
	if component == -1 || (component&1) != 0 {
		components["scheme"] = u.Scheme
	}
	if component == -1 || (component&2) != 0 {
		components["host"] = hostname
	}
	if component == -1 || (component&4) != 0 {
		components["port"] = port
	}
	if component == -1 || (component&8) != 0 {
		components["user"] = u.User.Username()
	}
	if component == -1 || (component&16) != 0 {
		components["pass"], _ = u.User.Password()
	}
	if component == -1 || (component&32) != 0 {
		components["path"] = u.Path
	}
	if component == -1 || (component&64) != 0 {
		components["query"] = u.RawQuery
	}
	if component == -1 || (component&128) != 0 {
		components["fragment"] = u.Fragment
	}

	return components, nil
}
