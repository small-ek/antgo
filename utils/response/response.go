package response

// Write Return parameter 返回参数
type Write struct {
	Code    int         `json:"code"`            //Status Code 业务逻辑的状态码
	Message string      `json:"message"`         //Msg Prompt message 简短描述
	Data    interface{} `json:"data,omitempty"`  //Data 实际数据
	Error   interface{} `json:"error,omitempty"` //Error message 错误相关的信息
}

// Page Pagination return 分页返回
type Page struct {
	Total int64       `json:"total"` //Total pages 总数量
	Items interface{} `json:"items"` //Data set of the current page 数据集合
}

// Success Successfully returned 成功返回
func Success(msg string, code int, data ...interface{}) *Write {
	var lenData = len(data)
	if lenData == 1 {
		return &Write{Code: code, Message: msg, Data: data[0]}
	} else if lenData > 1 {
		return &Write{Code: code, Message: msg, Data: data}
	}

	return &Write{Code: code, Message: msg}
}

// Fail Error return 错误返回
func Fail(msg string, code int, err ...string) *Write {
	var lenErr = len(err)

	if lenErr == 1 {
		return &Write{Code: code, Message: msg, Error: err[0]}
	} else if lenErr > 1 {
		return &Write{Code: code, Message: msg, Error: err}
	}

	return &Write{Code: code, Message: msg}
}
