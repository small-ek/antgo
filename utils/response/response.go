package response

// Write Return parameter 返回参数
type Write struct {
	Status  int    `json:"status"`          // 0: 成功, 1: 失败
	Code    string `json:"code"`            // 业务状态码
	Message string `json:"message"`         //用户友好提示信息
	Data    any    `json:"data"`            //Data 实际数据
	Error   any    `json:"error,omitempty"` //Error message 错误相关的信息
}

// Page Pagination return 分页返回
type Page struct {
	Total int64 `json:"total"` //Total pages 总数量
	Items any   `json:"items"` //Data set of the current page 数据集合
}

// Success Successfully returned 成功返回
func Success(code, msg string, data ...any) *Write {
	var resultData interface{}
	if len(data) == 1 {
		resultData = data[0]
	} else if len(data) > 1 {
		resultData = data
	}

	return &Write{Status: 0, Code: code, Message: msg, Data: resultData}
}

// Fail Error return 错误返回
func Fail(code, msg string, err ...string) *Write {
	var errorData any
	if len(err) == 1 {
		errorData = err[0]
	} else if len(err) > 1 {
		errorData = err
	}

	return &Write{Status: 1, Code: code, Message: msg, Error: errorData}
}
