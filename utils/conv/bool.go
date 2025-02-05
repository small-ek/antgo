package conv

import "strings"

// Bool converts `any` to bool.
// 将 `any` 转换为 bool。
func Bool(any interface{}) bool {
	if any == nil {
		return false
	}

	switch v := any.(type) {
	case bool:
		return v
	case []byte:
		return !strings.EqualFold(string(v), "false")
	case string:
		return !strings.EqualFold(v, "false")
	}
	return false
}
