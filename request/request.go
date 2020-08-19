package request

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
