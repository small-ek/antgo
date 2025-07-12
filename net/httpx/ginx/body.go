package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/net/httpx"
)

// Body 是 Gin 的请求体读取封装，自动重设 c.Request.Body
func Body(c *gin.Context) ([]byte, error) {
	maxSize := httpx.CalculateMaxSize(c.Request.ContentLength)

	body, newRC, err := httpx.ReadBody(c.Request.Body, maxSize)
	if err != nil {
		return nil, err
	}

	c.Request.Body = newRC
	return body, nil
}
