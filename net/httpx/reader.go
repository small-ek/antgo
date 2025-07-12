package httpx

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

const (
	DefaultBodyMaxSize = 1 << 20  // 1MB
	AbsoluteBodyMax    = 10 << 20 // 10MB
	PoolBufferSize     = 128 << 10
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, PoolBufferSize))
	},
}

// ReadBody 读取并限制请求体，返回 []byte 和可复用的 io.ReadCloser（供后续再次读取）
// 用于中间件场景，可记录日志后再绑定 JSON。
func ReadBody(rc io.ReadCloser, maxSize int64) ([]byte, io.ReadCloser, error) {
	if rc == nil {
		return nil, io.NopCloser(bytes.NewReader(nil)), nil
	}
	defer rc.Close()

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	limited := io.LimitReader(rc, maxSize+1)
	n, err := buf.ReadFrom(limited)
	if err != nil && err != io.EOF {
		bufPool.Put(buf)
		return nil, nil, err
	}
	if n > maxSize {
		bufPool.Put(buf)
		return nil, nil, ErrBodySizeExceeded{MaxAllowed: maxSize, Actual: n}
	}

	data := make([]byte, n)
	copy(data, buf.Bytes())
	bufPool.Put(buf)

	newRC := io.NopCloser(bytes.NewReader(data))
	return data, newRC, nil
}

// CalculateMaxSize 依据 Content-Length 计算最大体积
func CalculateMaxSize(contentLength int64) int64 {
	if contentLength > 0 {
		max := contentLength + contentLength/4
		if max > AbsoluteBodyMax {
			return AbsoluteBodyMax
		}
		return max
	}
	return DefaultBodyMaxSize
}

type ErrBodySizeExceeded struct {
	MaxAllowed int64
	Actual     int64
}

func (e ErrBodySizeExceeded) Error() string {
	return fmt.Sprintf("body size exceeded: max=%d actual=%d", e.MaxAllowed, e.Actual)
}
