package response

// Write Return parameter
type Write struct {
	Status int         `json:"status"`           //Status Code
	Msg    string      `json:"msg"`              //Msg Prompt message
	Result interface{} `json:"result,omitempty"` //Data
	Error  interface{} `json:"error,omitempty"`  //Error message
}

// Page Pagination return
type Page struct {
	Total int64       `json:"total"` //Total total pages
	List  interface{} `json:"list"`  //List json data
}

// Success Successfully returned
func Success(msg string, status int, data ...interface{}) *Write {
	var lenData = len(data)
	if lenData == 1 {
		return &Write{Status: status, Msg: msg, Result: data[0]}
	} else if lenData > 1 {
		return &Write{Status: status, Msg: msg, Result: data}
	}

	return &Write{Status: status, Msg: msg}
}

// Fail Error return, the second parameter is passed back to the front end and printed
func Fail(msg string, status int, err ...string) *Write {
	var lenErr = len(err)

	if lenErr == 1 {
		return &Write{Status: status, Msg: msg, Error: err[0]}
	} else if lenErr > 1 {
		return &Write{Status: status, Msg: msg, Error: err}
	}

	return &Write{Status: status, Msg: msg}
}
