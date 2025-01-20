package errors

import (
	"fmt"
	"io"
)

// Errorf 格式化并返回一个带有堆栈跟踪信息的错误对象。
// 这个函数接收格式化字符串和参数，并将其与当前堆栈信息一起封装为一个错误。
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...), // 格式化错误消息
		stack: callers(),                    // 获取当前的堆栈信息
	}
}

// fundamental 结构体用于表示一个基础错误，它包含错误消息和堆栈信息。
type fundamental struct {
	msg    string // 错误消息
	*stack        // 嵌入堆栈信息
}

// Error 方法实现了 error 接口，返回错误的消息。
func (f *fundamental) Error() string {
	return f.msg
}

// Format 方法实现了 fmt.Formatter 接口，允许自定义堆栈信息的格式化输出。
func (f *fundamental) Format(s fmt.State, verb rune) {
	f.StackTrace().Format(s, verb) // 格式化堆栈信息
}

// withStack 结构体表示一个带有堆栈信息的错误，它包装了一个原始错误。
type withStack struct {
	error  // 内嵌原始错误
	*stack // 嵌入堆栈信息
}

// Cause 方法返回原始错误，用于链式错误的获取。
func (w *withStack) Cause() error {
	return w.error
}

// Unwrap 提供 Go 1.13 错误链兼容性，返回原始错误。
func (w *withStack) Unwrap() error {
	return w.error
}

// Format 方法实现了 fmt.Formatter 接口，允许自定义格式化输出。
// 支持的格式：
//
//	%v	输出错误信息和堆栈跟踪。
//	%s	只输出错误消息。
//	%q	输出错误消息的引用格式。
func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// 使用 %+v 格式时输出详细的错误信息和堆栈跟踪
		fmt.Fprintf(s, "%+v", w.Cause())
		w.stack.StackTrace().Format(s, verb)
		return
	case 's':
		io.WriteString(s, w.Error()) // 只输出错误消息
	case 'q':
		fmt.Fprintf(s, "%q", w.Error()) // 输出错误消息的引用格式
	}
}

// Wrapf 返回一个错误，它会将原始错误（err）与当前堆栈信息以及格式化的消息一起包装。
// 如果原始错误为 nil，则返回 nil。
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	type wrapper interface {
		Unwrap() error
	}

	// 如果错误已经被包装，则返回带有新的格式化消息的包装错误
	if _, ok := err.(wrapper); ok {
		return WithMessagef(err, format, args...)
	}

	// 否则，创建一个新的带有堆栈跟踪的包装错误
	return &withStack{
		WithMessagef(err, format, args...),
		callers(),
	}
}

// WithStack 返回一个带有堆栈信息的包装错误。
// 如果原始错误已经被包装，则直接返回该错误。
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	type wrapper interface {
		Unwrap() error
	}

	// 如果错误已经被包装，则返回原始错误
	if _, ok := err.(wrapper); ok {
		return err
	}

	// 否则，创建一个新的带堆栈跟踪的包装错误
	return &withStack{
		err,
		callers(),
	}
}

// WithMessagef 将格式化的消息与错误一起包装。
// 如果原始错误为 nil，则返回 nil。
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

// withMessage 结构体表示一个带有消息的错误，它包装了原始错误。
type withMessage struct {
	cause error  // 原始错误
	msg   string // 错误消息
}

// Error 方法实现了 error 接口，返回带有格式化消息的错误信息。
func (w *withMessage) Error() string {
	if w.msg == "" {
		return w.cause.Error() // 如果消息为空，则只返回原始错误的信息
	}
	return w.msg + ": " + w.cause.Error() // 返回格式化的错误消息
}

// Cause 返回原始错误，供错误链使用。
func (w *withMessage) Cause() error {
	return w.cause
}

// Unwrap 提供 Go 1.13 错误链兼容性，返回原始错误。
func (w *withMessage) Unwrap() error {
	return w.cause
}

// Format 方法实现了 fmt.Formatter 接口，用于格式化错误消息的输出。
// 支持的格式：
//
//	%v	输出详细的错误信息，包含消息和堆栈跟踪。
//	%s	输出错误消息。
//	%q	输出错误消息的引用格式。
func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if w.msg != "" { // 如果消息不为空，则输出消息
			io.WriteString(s, fmt.Sprintln(w.msg))
		}
		// 输出堆栈跟踪
		fmt.Fprintf(s, "%+v\n", w.Cause())
		return
	case 's', 'q':
		io.WriteString(s, w.Error()) // 输出错误消息
	}
}

// Cause 返回错误的根本原因。
// 如果错误实现了 Cause() 方法，则返回底层的错误，否则返回原始错误。
// 如果错误为 nil，则返回 nil。
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer) // 如果错误实现了 Cause 方法，则获取底层错误
		if !ok {
			break // 如果错误不支持 Cause 方法，则退出
		}
		err = cause.Cause() // 继续查找下一个错误
	}
	return err
}
