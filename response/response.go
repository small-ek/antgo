package response

import (
	"github.com/small-ek/ginp/os/logger"
	"go.uber.org/zap"
)

const (
	ERROR   = 403
	SUCCESS = 200
)

//返回
type Write struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

//分页返回
type Page struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

//错误输出
func ErrorResponse(err error) *Write {
	return &Write{
		Code:  ERROR,
		Msg:   "参数错误",
		Error: err.Error(),
	}
}

//成功返回
func Success(msg string, data ...interface{}) *Write {
	var lenData = len(data)

	if lenData == 1 {
		return &Write{Code: SUCCESS, Msg: msg, Data: data[0]}
	} else if lenData > 1 {
		return &Write{Code: SUCCESS, Msg: msg, Data: data}
	}

	return &Write{Code: SUCCESS, Msg: msg}
}

//错误返回,第二个参数传参返回给前端并会打印
func Fail(msg string, err ...interface{}) *Write {

	if len(err) > 0 {
		logger.Write.Error("错误", zap.Any("error", err[0].(string)))
		return &Write{Code: ERROR, Msg: msg, Error: err[0].(string), Data: ""}
	}

	return &Write{Code: ERROR, Msg: msg}
}
