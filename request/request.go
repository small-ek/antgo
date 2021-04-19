package request

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/logs"
	"io/ioutil"
	"net/http"
)

//PageParam Paging parameters
type PageParam struct {
	CurrentPage int         `form:"current_page" json:"current_page"`
	PageSize    int         `form:"page_size" json:"page_size"`
	Total       int64       `form:"total" json:"total"`
	Filter      []string    `form:"filter[]" json:"filter[]"`
	Order       string      `form:"order" json:"order"`
	Details     interface{} `form:"details" json:"details"`
}

//DefaultPage Default pagination
func DefaultPage() PageParam {
	return PageParam{
		CurrentPage: 1,
		PageSize:    10,
	}
}

//GetBody Get the requested data
func GetBody(c *gin.Context) map[string]interface{} {
	var request map[string]interface{}
	var body []byte
	if c.Request.Body != nil {
		body, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	var err = json.Unmarshal(body, &request)
	if err != nil {
		logs.Error(err.Error())
	}
	return request
}

//Query ...
func Query(name string, c *gin.Context) string {
	return c.Query(name)
}

//QueryArray ...
func QueryArray(name string, c *gin.Context) []string {
	return c.QueryArray(name)
}

//QueryMap ...
func QueryMap(name string, c *gin.Context) map[string]string {
	return c.QueryMap(name)
}

//Param ...
func Param(name string, c *gin.Context) string {
	return c.Param(name)
}

//GetInterface ...
func GetInterface(name string, c *gin.Context) interface{} {
	var request = GetBody(c)
	return request[name]
}

//GetString ...
func GetString(name string, c *gin.Context) string {
	var request = GetBody(c)
	return conv.String(request[name])
}

//GetBool ...
func GetBool(name string, c *gin.Context) bool {
	var request = GetBody(c)
	return conv.Bool(request[name])
}

//GetUint ...
func GetUint(name string, c *gin.Context) uint {
	var request = GetBody(c)
	return conv.Uint(request[name])
}

//GetUint8 ...
func GetUint8(name string, c *gin.Context) uint8 {
	var request = GetBody(c)
	return conv.Uint8(request[name])
}

//GetUint16 ...
func GetUint16(name string, c *gin.Context) uint16 {
	var request = GetBody(c)
	return conv.Uint16(request[name])
}

//GetUint32 ...
func GetUint32(name string, c *gin.Context) uint32 {
	var request = GetBody(c)
	return conv.Uint32(request[name])
}

//GetUint64 ...
func GetUint64(name string, c *gin.Context) uint64 {
	var request = GetBody(c)
	return conv.Uint64(request[name])
}

//GetFloat32 ...
func GetFloat32(name string, c *gin.Context) float32 {
	var request = GetBody(c)
	return conv.Float32(request[name])
}

//GetFloat64 ...
func GetFloat64(name string, c *gin.Context) float64 {
	var request = GetBody(c)
	return conv.Float64(request[name])
}

//GetInt ...
func GetInt(name string, c *gin.Context) int {
	var request = GetBody(c)
	return conv.Int(request[name])
}

//GetInt8 ...
func GetInt8(name string, c *gin.Context) int8 {
	var request = GetBody(c)
	return conv.Int8(request[name])
}

//GetInt16 ...
func GetInt16(name string, c *gin.Context) int16 {
	var request = GetBody(c)
	return conv.Int16(request[name])
}

//GetInt32 ...
func GetInt32(name string, c *gin.Context) int32 {
	var request = GetBody(c)
	return conv.Int32(request[name])
}

//GetInt64 ...
func GetInt64(name string, c *gin.Context) int64 {
	var request = GetBody(c)
	return conv.Int64(request[name])
}

//Cors Cross-domain request
func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")                                      //允许访问所有域
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE,OPTIONS") //允许请求类型
	c.Header("Access-Control-Allow-Credentials", "true")                              //服务器是否接受浏览器发送的Cookie
	c.Header("Connection", "keep-alive")                                              //可以使一次TCP连接为同意用户的多次请求服务,提高了响应速度。
	c.Header("Access-Control-Max-Age", "3600")                                        //多少秒以后再次OPTIONS.默认60分钟
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, If-Match, If-Modified-Since, If-None-Match, If-Unmodified-Since, X-Requested-With")
	c.Header("Strict-Transport-Security", " max-age=63072000; includeSubdomains; preload")
	//放行所有OPTIONS方法
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	//处理请求
	c.Next()
}
