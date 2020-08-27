package request

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/ginp/conv"
	"io/ioutil"
	"net/http"
)

//获取请求的数据
func Get(this *gin.Context) map[string]interface{} {
	var request map[string]interface{}
	var body []byte
	if this.Request.Body != nil {
		body, _ = ioutil.ReadAll(this.Request.Body)
		//把刚刚读出来的再写进去其他地方使用没有
		this.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	json.Unmarshal(body, &request)
	return request
}

//获取请求单个数据
func GetString(name string, this *gin.Context) string {
	var request = Get(this)
	return conv.String(request[name])
}

//获取请求单个数据
func GetBool(name string, this *gin.Context) bool {
	var request = Get(this)
	return conv.Bool(request[name])
}

//获取请求单个数据
func GetFloat21(name string, this *gin.Context) float32 {
	var request = Get(this)
	return conv.Float32(request[name])
}

//获取请求单个数据
func GetFloat64(name string, this *gin.Context) float64 {
	var request = Get(this)
	return conv.Float64(request[name])
}

//获取请求单个数据
func GetInt(name string, this *gin.Context) int {
	var request = Get(this)
	return conv.Int(request[name])
}

//获取请求单个数据
func Param(name string, this *gin.Context) interface{} {
	var request = Get(this)
	return request[name]
}

//跨域请求
func Cors(this *gin.Context) {
	this.Header("Access-Control-Allow-Origin", "*")                                      //允许访问所有域
	this.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE,OPTIONS") //允许请求类型
	this.Header("Access-Control-Allow-Credentials", "true")                              //服务器是否接受浏览器发送的Cookie
	this.Header("Connection", "keep-alive")                                              //可以使一次TCP连接为同意用户的多次请求服务,提高了响应速度。
	this.Header("Access-Control-Max-Age", "3600")                                        //多少秒以后再次OPTIONS.默认60分钟
	this.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, If-Match, If-Modified-Since, If-None-Match, If-Unmodified-Since, X-Requested-With")
	this.Header("Strict-Transport-Security", " max-age=63072000; includeSubdomains; preload")
	//放行所有OPTIONS方法
	if this.Request.Method == "OPTIONS" {
		this.AbortWithStatus(http.StatusNoContent)
	}
	//处理请求
	this.Next()
}
