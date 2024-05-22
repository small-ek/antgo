package response

import (
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
)

const (
	//ERROR Default error code returned
	ERROR = 422
	//SUCCESS Default success code return
	SUCCESS = 200
)

// Write Return parameter
type Write struct {
	Status int         `json:"status"` //Status Code
	Msg    string      `json:"msg"`    //Msg Prompt message
	Result interface{} `json:"result"` //Data
	Error  interface{} `json:"error"`  //Error message
}

// Page Pagination return
type Page struct {
	Total int64       `json:"total"` //Total total pages
	List  interface{} `json:"list"`  //List json data
}

// Success Successfully returned
func Success(msg string, data ...interface{}) *Write {
	var lenData = len(data)
	if lenData == 1 {
		return &Write{Status: SUCCESS, Msg: msg, Result: data[0]}
	} else if lenData > 1 {
		return &Write{Status: SUCCESS, Msg: msg, Result: data}
	}

	return &Write{Status: SUCCESS, Msg: msg}
}

// Fail Error return, the second parameter is passed back to the front end and printed
func Fail(msg string, err ...string) *Write {
	var lenErr = len(err)
	if lenErr > 0 {
		alog.Write.Error("Return error", zap.Any("error", err))
	}

	if lenErr == 1 {
		return &Write{Status: ERROR, Msg: msg, Error: err[0]}
	} else if lenErr > 1 {
		return &Write{Status: ERROR, Msg: msg, Error: err}
	}

	return &Write{Status: ERROR, Msg: msg}
}
