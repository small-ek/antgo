package response

// Write 结构体用于定义API的返回参数
// Write struct defines the response parameters for API
type Write struct {
	Status  int    `json:"status"`          // 0: 成功, 1: 失败 / 0: Success, 1: Failure
	Code    string `json:"code"`            // 业务状态码 / Business status code
	Message string `json:"message"`         // 用户友好提示信息 / User-friendly message
	Data    any    `json:"data"`            // 实际数据 / Actual data
	Error   any    `json:"error,omitempty"` // 错误信息 / Error message
}

// Page 结构体用于分页返回数据
// Page struct defines the pagination response data
type Page struct {
	Total int64 `json:"total"` // 总数量 / Total number of items
	Items any   `json:"items"` // 当前页的数据集合 / Data set of the current page
}

// Success 函数用于成功返回
// Success function is used for successful response
func Success(code, msg string, data ...any) *Write {
	var resultData interface{}
	if len(data) == 1 {
		resultData = data[0]
	} else if len(data) > 1 {
		resultData = data
	}

	return &Write{
		Status:  0,
		Code:    code,
		Message: msg,
		Data:    resultData,
	}
}

// Fail 函数用于错误返回
// Fail function is used for error response
func Fail(code, msg string, err ...string) *Write {
	var errorData any
	if len(err) == 1 {
		errorData = err[0]
	} else if len(err) > 1 {
		errorData = err
	}

	return &Write{
		Status:  1,
		Code:    code,
		Message: msg,
		Error:   errorData,
	}
}

// PageResponse 函数用于分页数据返回
// PageResponse function is used for paginated data response
func PageResponse(total int64, items any) *Page {
	return &Page{
		Total: total,
		Items: items,
	}
}
