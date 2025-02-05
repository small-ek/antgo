package conv

import (
	"encoding/binary"
	"math"
	"strconv"
)

// Int converts "any" to int.
// Rules: <nil> => 0, bool => 0/1, numeric => value truncation, []byte => big-endian conversion,
// string => base10 parsing, others => 0
// 转换规则：<nil>转0，布尔转0/1，数字类型截断转换，字节切片按大端序转，字符串十进制解析，其他类型返回0
func Int(any interface{}) int {
	if any == nil {
		return 0
	}
	if v, ok := any.(int); ok { // 快速路径
		return v
	}
	return int(Int64(any)) // 统一通过Int64转换
}

// Int8 converts "any" to int8 with overflow checking.
// 带溢出检查的int8转换，处理方式与Int类似
func Int8(any interface{}) (v int8) {
	if any == nil {
		return 0
	}
	if tv, ok := any.(int8); ok {
		return tv
	}
	return clampInt64ToInt8(Int64(any))
}

// Int16 converts "any" to int16 with overflow checking.
// 带溢出检查的int16转换
func Int16(any interface{}) (v int16) {
	if any == nil {
		return 0
	}
	if tv, ok := any.(int16); ok {
		return tv
	}
	return clampInt64ToInt16(Int64(any))
}

// Int32 converts "any" to int32 with overflow checking.
// 带溢出检查的int32转换
func Int32(any interface{}) (v int32) {
	if any == nil {
		return 0
	}
	if tv, ok := any.(int32); ok {
		return tv
	}
	return clampInt64ToInt32(Int64(any))
}

// Int64 converts "any" to int64 with value clamping and safe parsing.
// Supported types: all numeric types, bool, []byte, string
// 安全转换函数，支持所有数字类型、布尔、字节切片和字符串
func Int64(any interface{}) int64 {
	if any == nil {
		return 0
	}

	switch val := any.(type) {
	// 按出现频率排序优化
	case int64:
		return val
	case int:
		return int64(val)
	case uint64:
		return clampUint64ToInt64(val)
	case float64:
		return clampFloatToInt(val, 64)
	case string:
		return parseIntString(val)
	case []byte:
		return parseBigEndianBytes(val)
	case bool:
		return boolToInt(val)
	case int32:
		return int64(val)
	case uint32:
		return int64(val)
	case float32:
		return clampFloatToInt(float64(val), 32)
	case int16:
		return int64(val)
	case uint16:
		return int64(val)
	case int8:
		return int64(val)
	case uint8:
		return int64(val)
	case uint:
		return int64(val)
	default:
		return 0
	}
}

/*--------------------------  Unsigned Conversions --------------------------*/

// Uint converts "any" to uint with overflow checking.
// 带溢出检查的uint转换
func Uint(any interface{}) uint {
	if any == nil {
		return 0
	}
	if tv, ok := any.(uint); ok {
		return tv
	}
	return clampUint64ToUint(Uint64(any))
}

// Uint8 converts "any" to uint8 with overflow checking.
// 带溢出检查的uint8转换
func Uint8(any interface{}) uint8 {
	if any == nil {
		return 0
	}
	if tv, ok := any.(uint8); ok {
		return tv
	}
	return clampUint64ToUint8(Uint64(any))
}

// Uint16 converts "any" to uint16 with overflow checking.
// 带溢出检查的uint16转换
func Uint16(any interface{}) uint16 {
	if any == nil {
		return 0
	}
	if tv, ok := any.(uint16); ok {
		return tv
	}
	return clampUint64ToUint16(Uint64(any))
}

// Uint32 converts "any" to uint32 with overflow checking.
// 带溢出检查的uint32转换
func Uint32(any interface{}) uint32 {
	if any == nil {
		return 0
	}
	if tv, ok := any.(uint32); ok {
		return tv
	}
	return clampUint64ToUint32(Uint64(any))
}

// Uint64 converts "any" to uint64 with value clamping and safe parsing.
// 安全转换函数，支持所有数字类型、布尔、字节切片和字符串
func Uint64(any interface{}) uint64 {
	if any == nil {
		return 0
	}

	switch val := any.(type) {
	// 按出现频率排序优化
	case uint64:
		return val
	case uint:
		return uint64(val)
	case int64:
		return clampInt64ToUint64(val)
	case float64:
		return clampFloatToUint(val, 64)
	case string:
		return parseUintString(val)
	case []byte:
		return parseBigEndianBytesToUint(val)
	case bool:
		return boolToUint(val)
	case uint32:
		return uint64(val)
	case int:
		return clampInt64ToUint64(int64(val))
	case float32:
		return clampFloatToUint(float64(val), 32)
	case uint16:
		return uint64(val)
	case int32:
		return clampInt64ToUint64(int64(val))
	case uint8:
		return uint64(val)
	case int16:
		return clampInt64ToUint64(int64(val))
	case int8:
		return clampInt64ToUint64(int64(val))
	default:
		return 0
	}
}

/*--------------------------  Core Conversion Logic --------------------------*/

// 字节转换逻辑（大端序处理，自动填充）
func parseBigEndianBytes(b []byte) int64 {
	if len(b) == 0 {
		return 0
	}

	// 处理超过8字节的情况
	if len(b) > 8 {
		b = b[len(b)-8:] // 取最后8字节
	}

	var buf [8]byte
	copy(buf[8-len(b):], b) // 大端序高位填充
	return int64(binary.BigEndian.Uint64(buf[:]))
}

func parseBigEndianBytesToUint(b []byte) uint64 {
	if len(b) == 0 {
		return 0
	}

	if len(b) > 8 {
		b = b[len(b)-8:]
	}

	var buf [8]byte
	copy(buf[8-len(b):], b)
	return binary.BigEndian.Uint64(buf[:])
}

// 字符串解析（带进制检测）
func parseIntString(s string) int64 {
	if num, err := strconv.ParseInt(s, 0, 64); err == nil {
		return num
	}
	return 0
}

func parseUintString(s string) uint64 {
	if num, err := strconv.ParseUint(s, 0, 64); err == nil {
		return num
	}
	return 0
}

/*--------------------------  Value Clamping Functions --------------------------*/

// 浮点转整型时的安全钳制
func clampFloatToInt(f float64, bitSize int) int64 {
	if math.IsNaN(f) {
		return 0
	}

	// 根据bitSize计算最大最小值
	maxVal := float64(math.MaxInt64)
	minVal := float64(math.MinInt64)
	if bitSize < 64 {
		maxVal = math.Pow(2, float64(bitSize-1)) - 1
		minVal = -math.Pow(2, float64(bitSize-1))
	}

	f = math.Max(math.Min(f, maxVal), minVal)
	return int64(math.Trunc(f))
}

func clampFloatToUint(f float64, bitSize int) uint64 {
	if math.IsNaN(f) || f < 0 {
		return 0
	}

	var maxVal uint64
	switch {
	case bitSize == 64:
		maxVal = math.MaxUint64
	case bitSize >= 0 && bitSize < 64:
		maxVal = (1 << bitSize) - 1
	default:
		return 0
	}

	if f > float64(maxVal) || f >= math.MaxUint64 {
		return maxVal
	}

	return uint64(f)
}

// 带符号整型溢出保护
func clampInt64ToInt8(i int64) int8 {
	if i > math.MaxInt8 {
		return math.MaxInt8
	}
	if i < math.MinInt8 {
		return math.MinInt8
	}
	return int8(i)
}

func clampInt64ToInt16(i int64) int16 {
	if i > math.MaxInt16 {
		return math.MaxInt16
	}
	if i < math.MinInt16 {
		return math.MinInt16
	}
	return int16(i)
}

func clampInt64ToInt32(i int64) int32 {
	if i > math.MaxInt32 {
		return math.MaxInt32
	}
	if i < math.MinInt32 {
		return math.MinInt32
	}
	return int32(i)
}

func clampUint64ToInt64(u uint64) int64 {
	if u > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(u)
}

// 无符号整型溢出保护
func clampUint64ToUint(u64 uint64) uint {
	if u64 > math.MaxUint {
		return math.MaxUint
	}
	return uint(u64)
}

func clampUint64ToUint8(u64 uint64) uint8 {
	if u64 > math.MaxUint8 {
		return math.MaxUint8
	}
	return uint8(u64)
}

func clampUint64ToUint16(u64 uint64) uint16 {
	if u64 > math.MaxUint16 {
		return math.MaxUint16
	}
	return uint16(u64)
}

func clampUint64ToUint32(u64 uint64) uint32 {
	if u64 > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(u64)
}

func clampInt64ToUint64(i int64) uint64 {
	if i < 0 {
		return 0
	}
	return uint64(i)
}

// 布尔值转换
func boolToInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func boolToUint(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}
