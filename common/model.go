package common

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (r *Result) Success(data any) *Result {
	return &Result{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func (r *Result) Failed(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
	}
}
