package response

import (
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
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
	StatusCode int         `json:"statusCode"` //Error Code
	Msg        string      `json:"msg"`        //Msg Prompt message
	Result     interface{} `json:"result"`     //Data
	Error      interface{} `json:"error"`      //Error message
}

// Page Pagination return
type Page struct {
	Total int64       `json:"total"` //Total total pages
	List  interface{} `json:"list"`  //List json data
}

// ErrorResponse Error output
func ErrorResponse(err string) *Write {
	return &Write{
		StatusCode: ERROR,
		Msg:        "Error",
		Error:      err,
	}
}

// Success Successfully returned
func Success(msg string, data ...interface{}) *Write {
	var lenData = len(data)
	if lenData == 1 {
		return &Write{StatusCode: SUCCESS, Msg: msg, Result: data[0]}
	} else if lenData > 1 {
		return &Write{StatusCode: SUCCESS, Msg: msg, Result: data}
	}
	
	return &Write{StatusCode: SUCCESS, Msg: msg}
}

// Fail Error return, the second parameter is passed back to the front end and printed
func Fail(msg string, err ...string) *Write {
	var lenErr = len(err)
	if lenErr > 0 && config.GetBool("system.debug") == true {
		alog.Write.Debug("Return error", zap.Any("error", err))
	}

	if lenErr == 1 {
		return &Write{StatusCode: ERROR, Msg: msg, Error: err[0]}
	} else if lenErr > 1 {
		return &Write{StatusCode: ERROR, Msg: msg, Error: err}
	}

	return &Write{StatusCode: ERROR, Msg: msg}
}
