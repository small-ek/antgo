package request

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/ginp/conv"
	"io/ioutil"
	"log"
	"net/http"
)

//PageParam Paging parameters
type PageParam struct {
	CurrentPage int         `form:"current_page"`
	PageSize    int         `form:"page_size"`
	Total       int64       `form:"total"`
	Filter      interface{} `form:"filter"`
	Order       string      `form:"order"`
	Details     interface{} `json:"details"`
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
		log.Println(err.Error())
	}
	return request
}

//Input ...
func Input(name string, c *gin.Context) interface{} {
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
