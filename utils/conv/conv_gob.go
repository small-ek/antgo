package conv

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
	"sync"
)

// 包级别错误定义，避免重复创建错误对象 / Package-level error definitions to avoid recreating error objects
var (
	// ErrNilData 表示输入数据为空
	// ErrNilData indicates input data is nil
	ErrNilData = errors.New("conv: data cannot be nil")

	// ErrNilTarget 表示目标对象为空
	// ErrNilTarget indicates target object is nil
	ErrNilTarget = errors.New("conv: target cannot be nil")

	// ErrTargetNotPointer 表示目标对象不是有效指针类型
	// ErrTargetNotPointer indicates target is not a valid pointer type
	ErrTargetNotPointer = errors.New("conv: target must be non-nil pointer")
)

// bufferPool 用于复用 bytes.Buffer 对象，减少内存分配
// bufferPool for reusing bytes.Buffer objects to reduce memory allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GobEncoder 使用gob编码数据（高性能版本）
// 特性：
// 1. 使用内存池减少内存分配
// 2. 严格参数校验
// 3. 自动处理缓冲区生命周期
//
// GobEncoder encodes data using gob (high-performance version)
// Features:
// 1. Uses memory pool to reduce allocations
// 2. Strict parameter validation
// 3. Automatic buffer lifecycle management
func GobEncoder(data any) ([]byte, error) {
	if data == nil {
		return nil, ErrNilData
	}

	// 从内存池获取缓冲区 / Get buffer from pool
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	// 执行编码 / Perform encoding
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	// 创建数据副本防止缓冲区重用导致的数据污染
	// Create data copy to prevent data pollution from buffer reuse
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())

	return result, nil
}

// GobDecoder 使用gob解码数据（增强安全版本）
// 特性：
// 1. 严格参数类型校验
// 2. 详细的错误类型返回
// 3. 自动处理缓冲区
//
// GobDecoder decodes data using gob (enhanced safety version)
// Features:
// 1. Strict parameter type validation
// 2. Detailed error type returns
// 3. Automatic buffer handling
func GobDecoder(data []byte, target any) error {
	if data == nil {
		return ErrNilData
	}
	if target == nil {
		return ErrNilTarget
	}

	// 使用反射验证目标对象类型 / Use reflection to validate target type
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrTargetNotPointer
	}

	// 创建带数据缓存的解码器 / Create decoder with data buffer
	dec := gob.NewDecoder(bytes.NewReader(data))
	return dec.Decode(target)
}
