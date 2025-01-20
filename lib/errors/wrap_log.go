package errors

import (
	"fmt"

	// "io"
	"slices"
)

// 日志级别常量定义，后续对接日志
const (
	LevelWarn  = iota + 1 // 1: Warn 级别
	LevelInfo             // 2: Info 级别
	LevelTrace            // 3: Trace 级别
	LevelDebug            // 4: Debug 级别
	LevelError            // 5: Error 级别
	LevelFatal            // 6: Fatal 级别
	LevelBuss             // 7: Buss 业务级别日志
)

// loggableLevel 类型用于表示日志级别。
type loggableLevel int

var (
	defaultLoggableLevel = loggableLevel(LevelWarn) // 默认日志级别为 WARN
)

// ILogLevel 接口定义了一个结构体，该结构体可以提供日志级别并且具有错误链功能。
type ILogLevel interface {
	error                 // 实现了 error 接口
	Level() loggableLevel // 返回日志级别
	Cause() error         // 获取底层错误
}

// LoggableLevel 获取 error 中包含的日志级别。
// 它会遍历错误链，返回包含最大日志级别的错误和该日志级别。
func LoggableLevel(err error) (cause error, level int, ok bool) {
	cause, l, ok := GetLoggableLevel(err)
	return cause, int(l), ok
}

// GetLoggableLevel 返回错误链中的日志级别。遍历整个错误链，获取最大日志级别。
func GetLoggableLevel(err error) (causeErr error, level loggableLevel, ok bool) {
	level = defaultLoggableLevel
	levels := make([]loggableLevel, 0, 4)
	for err != nil {
		loggable, ok := err.(ILogLevel)
		if !ok {
			cause, ok := err.(interface{ Cause() error })
			if !ok {
				causeErr = err
				break
			}
			err = cause.Cause()
			causeErr = err
			continue
		}

		level = loggable.Level()
		levels = append(levels, level)
		err = loggable.Cause()
		causeErr = err
	}

	l := len(levels)
	if l == 0 {
		return causeErr, level, false
	} else if l == 1 {
		return causeErr, levels[0], true
	}

	// 返回最大日志级别
	return causeErr, slices.Max(levels), true
}

// loggableLevelMsg 封装了带有日志级别信息的错误。
type loggableLevelMsg struct {
	withMessage
	level loggableLevel
}

// Level 返回错误的日志级别。
func (llm *loggableLevelMsg) Level() loggableLevel {
	if llm == nil {
		return LevelInfo
	}
	return llm.level
}

// WithTraceLogLevel 封装错误并将其标记为 TRACE 级别日志。
func WithTraceLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelTrace,
	}
}

// WithTraceLogLevelMsg 封装错误并返回 TRACE 级别日志，允许格式化消息。
func WithTraceLogLevelMsg(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		level: LevelTrace,
	}
}

// WithDebugLogLevel 封装错误并将其标记为 DEBUG 级别日志。
func WithDebugLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelDebug,
	}
}

// WithDebugLogLevelMsg 封装错误并返回 DEBUG 级别日志，允许格式化消息。
func WithDebugLogLevelMsg(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		level: LevelDebug,
	}
}

// WithInfoLogLevel 封装错误并将其标记为 INFO 级别日志。
func WithInfoLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelInfo,
	}
}

// WithInfoLogLevelMsg 封装错误并返回 INFO 级别日志，允许格式化消息。
func WithInfoLogLevelMsg(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		level: LevelInfo,
	}
}

// WithWarnLogLevel 封装错误并将其标记为 WARN 级别日志。
func WithWarnLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelWarn,
	}
}

// WithWarnLogLevelMsg 封装错误并返回 WARN 级别日志，允许格式化消息。
func WithWarnLogLevelMsg(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		level: LevelWarn,
	}
}

// WithErrorLogLevel 封装错误并将其标记为 ERROR 级别日志。
func WithErrorLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelError,
	}
}

// WithErrorLogLevelMsg 封装错误并返回 ERROR 级别日志，允许格式化消息。
func WithErrorLogLevelMsg(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		level: LevelError,
	}
}

// WithFatalLogLevel 封装错误并将其标记为 FATAL 级别日志。
func WithFatalLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelFatal,
	}
}

// WithFatalLogLevelMsg 封装错误并返回 FATAL 级别日志，允许格式化消息。
func WithFatalLogLevelMsg(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		level: LevelFatal,
	}
}

// WithBussLogLevel 封装错误并将其标记为 BUSS 业务级别日志。
func WithBussLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelBuss,
	}
}

// WithBussLogLevelMsg 封装错误并返回 BUSS 级别日志，允许格式化消息。
func WithBussLogLevelMsg(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		level: LevelBuss,
	}
}
