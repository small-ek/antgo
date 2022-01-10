package request

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/utils/conv"
	"io/ioutil"
)

//PageParam Paging parameters
type PageParam struct {
	CurrentPage int                    `form:"current_page" json:"current_page" bson:"current_page" xml:"current_page" yaml:"current_page"`
	PageSize    int                    `form:"page_size" json:"page_size" bson:"page_size" xml:"page_size" yaml:"page_size"`
	Total       int64                  `form:"total" json:"total" bson:"total" xml:"total" yaml:"total"`
	Filter      []string               `form:"filter[]" json:"filter[]" bson:"filter[]" xml:"filter[]" yaml:"filter[]"`
	Order       string                 `form:"order" json:"order" bson:"order" xml:"order" yaml:"order"`
	Select      []string               `form:"select[]" json:"select[]" bson:"select[]" xml:"select[]" yaml:"select[]"`
	Group       string                 `form:"group" json:"group" bson:"group" xml:"group" yaml:"group"`
	Omit        string                 `form:"omit" json:"omit" bson:"omit" xml:"omit" yaml:"omit"`
	Extra       map[string]interface{} `form:"extra" json:"extra" bson:"extra" xml:"extra" yaml:"extra"`
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
		panic(err)
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
