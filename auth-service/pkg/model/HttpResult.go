package model

type HttpResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (h *HttpResult) Success(data any) *HttpResult {
	return &HttpResult{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func (h *HttpResult) Failed(code int, msg string) *HttpResult {
	return &HttpResult{
		Code: code,
		Msg:  msg,
	}
}
