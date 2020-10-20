package response

import (
	"github.com/small-ek/ginp/os/logger"
	"go.uber.org/zap"
)

const (
	//ERROR Default error code returned
	ERROR = 403
	//SUCCESS Default success code return
	SUCCESS = 200
)

//Write Return parameter
type Write struct {
	Code  int         `json:"code"`  //Code
	Msg   string      `json:"msg"`   //Msg Prompt message
	Data  interface{} `json:"data"`  //Data
	Error string      `json:"error"` //Error message
}

//Page Pagination return
type Page struct {
	Total int         `json:"total"` //Total total pages
	List  interface{} `json:"list"`  //List json data
}

//ErrorResponse Error output
func ErrorResponse(err error) *Write {
	return &Write{
		Code:  ERROR,
		Msg:   "错误",
		Error: err.Error(),
	}
}

//Success Successfully returned
func Success(msg string, data ...interface{}) *Write {
	var lenData = len(data)
	if lenData == 1 {
		return &Write{Code: SUCCESS, Msg: msg, Data: data[0]}
	} else if lenData > 1 {
		return &Write{Code: SUCCESS, Msg: msg, Data: data}
	}
	return &Write{Code: SUCCESS, Msg: msg}
}

//Fail Error return, the second parameter is passed back to the front end and printed
func Fail(msg string, err ...interface{}) *Write {
	if len(err) > 0 {
		logger.Write.Error("错误", zap.Any("error", err[0].(string)))
		return &Write{Code: ERROR, Msg: msg, Error: err[0].(string), Data: ""}
	}
	return &Write{Code: ERROR, Msg: msg}
}
