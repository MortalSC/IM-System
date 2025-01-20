package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error 结构体，嵌套了 ResultMessage 用于表示一个带有详细信息的错误
type Error struct {
	ResultMessage
}

// ResultMessage 结构体，包含了与错误相关的详细信息
// 主要用于存储错误的各类信息，如结果、消息、错误代码等
type ResultMessage struct {
	Data    interface{} `json:"data,omitempty"` // 如果字段为空，JSON 序列化时会被忽略
	Result  string      `json:"result"`         // 表示结果状态（如 "success", "failure"）
	Message string      `json:"message"`        // 错误消息，通常用于描述错误的原因
	Code    int         `json:"code"`           // 错误的通用状态码
	ErrCode int         `json:"err-code"`       // 错误的特定错误代码
	Args    interface{} `json:"-"`              // 不会被序列化的参数，通常用于存储附加的内部信息
	Header  http.Header `json:"-"`              // 不会被序列化的 HTTP 头部信息
}

// Error 方法将 ResultMessage 转换为 JSON 字符串，返回错误的描述信息
func (rm *ResultMessage) Error() string {
	marshal, _ := json.Marshal(rm)
	return string(marshal)
}

// NewError 创建一个简单的错误对象，使用默认的错误消息和结果
func NewError(result string, code int) *Error {
	err := NewErrorEx(result, result, code) // 调用 NewErrorEx 创建错误
	return err
}

// NewErrorEx 创建一个带有自定义消息和状态码的错误对象
func NewErrorEx(result, msg string, code int) *Error {
	err := new(Error)
	err.Result = result // 设置结果
	err.Code = code     // 设置状态码
	err.Message = msg   // 设置消息
	return err
}

// NewErrEx 创建一个包含特定错误代码、状态码和消息的错误对象
func NewErrEx(result string, errCode int, code int, msg string) *Error {
	err := NewErrorEx(result, msg, code)
	err.ErrCode = errCode // 设置特定错误代码
	return err
}

// JsonString 将错误对象转换为 JSON 字符串
func (e *Error) JsonString() string {
	b, _ := json.Marshal(e) // 将 Error 结构体序列化为 JSON
	return string(b)
}

// Msg 设置错误消息并返回新的错误对象
func (e *Error) Msg(msg string) *Error {
	var err = *e
	err.Message = msg // 设置新的错误消息
	return &err
}

// Err 设置错误对象的消息为传入的错误的消息，并返回新的错误对象
func (e *Error) Err(err error) *Error {
	var errr = *e
	errr.Message = err.Error() // 使用传入错误的消息
	return &errr
}

// WithArgs 设置附加参数并返回新的错误对象
func (e *Error) WithArgs(args ...interface{}) *Error {
	var err = *e
	err.Args = args // 设置附加参数
	return &err
}

// WithData 设置附加数据并返回新的错误对象
func (e *Error) WithData(data ...interface{}) *Error {
	var err = *e
	err.Data = data // 设置附加数据
	return &err
}

// Is 判断当前错误与另一个错误是否相等
func (e *Error) Is(other error) bool {
	return e.Equal(other) // 调用 Equal 方法检查是否相等
}

// Equal 判断两个错误是否相等
func (e *Error) Equal(other error) bool {
	other = Cause(other) // 获取原始错误
	if e == other {
		return true // 如果是相同的错误，返回 true
	}

	if other == nil {
		return false // 如果另一个错误是 nil，则不相等
	}

	o, ok := other.(*Error) // 尝试将其他错误转换为 *Error 类型
	if !ok {
		return false // 如果转换失败，返回 false
	}

	// 如果 Result 字段相等，认为是相同的错误
	return e.Result == o.Result
}

// Error 重写 Error 方法，将 Error 对象转换为 JSON 字符串
func (e *Error) Error() string {
	return e.JsonString() // 调用 JsonString 方法返回错误的 JSON 字符串
}

// Errorf 创建一个带有格式化消息的错误对象
func (e *Error) Errorf(format string, args ...interface{}) *Error {
	return NewErrorEx(e.Result, fmt.Sprintf(format, args...), e.Code) // 格式化消息并返回新的错误对象
}
