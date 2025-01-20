package errors

import (
	"fmt"

	// "io"
	"slices"
)

// TODO: ；联动日志
const (
	LevelWarn = iota + 1
	LevelInfo
	LevelTrace
	LevelDebug
	LevelError
	LevelFatal
	LevelBuss
)

type loggableLevel int

var (
	defaultLoggableLevel = loggableLevel(LevelWarn) // 默认WARN级别
)

type ILogLevel interface {
	error
	Level() loggableLevel
	Cause() error
}

// 获取err中包含的日志级别
// error链中存在多个日志级别日志时取最大的级别日志
// 没有设置日志级别日志时返回WARN级别日志
// 返回日志级别的定义与kgo/log中的日志级别定义保持一致
//
// 返回参数解释
//
//	causeErr: 通过liberr.Cause拿到的底层的error
//	level: err所包含的日志级别
//	ok: 仅表示是否存在err是否为ILogLevel, ok==true不代表causeErr!=nil；ok==false时，表明err不是ILogLevel
func LoggableLevel(err error) (cause error, level int, ok bool) {
	cause, l, ok := GetLoggableLevel(err)
	return cause, int(l), ok
}

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

	return causeErr, slices.Max(levels), true
}

type loggableLevelMsg struct {
	withMessage
	level loggableLevel
}

func (llm *loggableLevelMsg) Level() loggableLevel {
	if llm == nil {
		return LevelInfo
	}
	return llm.level
}

// 封装error，返回上层打印TRACE级别日志
func WithTraceLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelTrace,
	}
}

// 封装error，返回上层打印TRACE级别日志
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

// 封装error，返回上层打印DEBUG级别日志
func WithDebugLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelDebug,
	}
}

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

// 封装error，返回到上层打印INFO级别日志
func WithInfoLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelInfo,
	}
}

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

// 封装error，返回到上层打印WARN级别日志
func WithWarnLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelWarn,
	}
}

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

// 封装error，返回到上层打印ERROR级别日志
func WithErrorLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelError,
	}
}

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

// 封装error，返回到上层打印FATAL级别日志
func WithFatalLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelFatal,
	}
}

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

// 封装error，返回到上层打印BUSS级别日志
func WithBussLogLevel(err error) error {
	if err == nil {
		return nil
	}

	return &loggableLevelMsg{
		withMessage: withMessage{cause: err},
		level:       LevelBuss,
	}
}

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
