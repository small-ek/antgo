package conv

import "strings"

// Bool converts `any` to bool.<将“any”转换为bool。>
func Bool(any interface{}) bool {
	if any == nil {
		return false
	}
	switch value := any.(type) {
	case bool:
		return value
	case []byte:
		if strings.ToLower(string(value)) == "false" {
			return false
		}
		return true
	case string:
		if strings.ToLower(value) == "false" {
			return false
		}
		return true
	}
	return false
}
