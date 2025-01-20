package errors

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

// Frame 表示程序栈中的一个调用帧。它是一个 uintptr 类型，用于存储程序计数器。
type Frame uintptr

// pc 返回帧的程序计数器地址。返回值是当前 Frame 减去 1。
func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}

// file 返回当前 Frame 所在的源文件路径。如果无法找到文件信息，则返回 "unknown"。
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line 返回当前 Frame 所在的行号。如果无法找到行号，则返回 0。
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name 返回当前 Frame 的函数名称。如果无法获取函数名，则返回 "unknown"。
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// callerInfo 返回当前调用帧的信息，包括行号和文件路径/函数名称。
func (f Frame) callerInfo() (line int, pathFileFunc string) {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		pathFileFunc = "unknown"
		return
	}
	pathFileFunc = fn.Name()
	pathFileFunc = leftTrimPath(pathFileFunc) // 去除路径中的多余部分
	_, line = fn.FileLine(f.pc())
	return
}

// leftTrimPath 从函数的路径中去除前三个目录部分，仅保留文件名。
func leftTrimPath(pathFile string) string {
	foundCnt := 0
	i := len(pathFile) - 1
	// 向后遍历路径，直到找到三个 "/" 为止
	for i >= 0 && foundCnt < 3 {
		if pathFile[i] == '/' {
			foundCnt++
		}
		i--
	}
	return pathFile[i+1:]
}

// Format 实现了 fmt.Formatter 接口，用于根据不同的格式化规则输出 Frame 信息。
// 支持格式：
//
//	%s    输出源文件路径
//	%d    输出源行号
//	%n    输出函数名称
//	%v    等效于 %s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		io.WriteString(s, f.name()) // 输出函数名称
		io.WriteString(s, "\n\t")
		io.WriteString(s, f.file()) // 输出源文件路径
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line())) // 输出行号
	case 'n':
		io.WriteString(s, funcname(f.name())) // 输出函数名（不带路径）
	case 'v':
		f.Format(s, 's') // 输出函数名
		io.WriteString(s, ":")
		f.Format(s, 'd') // 输出行号
	}
}

// MarshalText 实现了 fmt.TextMarshaler 接口，
// 将 Frame 格式化为文本字符串，输出格式为函数名、文件路径和行号。
func (f Frame) MarshalText() ([]byte, error) {
	name := f.name()
	if name == "unknown" {
		return []byte(name), nil
	}
	// 格式化输出函数名、文件路径和行号
	return []byte(fmt.Sprintf("%s %s:%d", name, f.file(), f.line())), nil
}

// StackTrace 是一个包含多个 Frame 的栈，表示函数调用栈
type StackTrace []Frame

// Format 实现了 fmt.Formatter 接口，用于格式化栈跟踪。
// 支持格式：
//
//	%s	输出栈中每个帧的源文件路径
//	%v	输出栈中每个帧的源文件路径和行号
//	%+v	输出栈中每个帧的完整信息：文件路径、函数名和行号
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('#'):
			// 带有 '#' 标志时输出详细格式
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			// 默认情况下输出每一帧的文件、函数和行号
			for _, f := range st {
				io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		}
	case 's':
		st.formatSlice(s, verb) // 使用 %s 格式时，调用 formatSlice 方法
	}
}

// formatSlice 将 StackTrace 格式化为一个帧的切片。
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb) // 格式化每个帧
	}
	io.WriteString(s, "]")
}

// stack 表示一个包含程序计数器的栈。
type stack []uintptr

// Format 实现了 fmt.Formatter 接口，用于格式化程序计数器栈。
func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		// 格式化栈中的每一帧并输出
		for _, pc := range *s {
			f := Frame(pc)
			fmt.Fprintf(st, "\n%+v", f)
		}
	}
}

// StackTrace 返回 stack 中的每一帧作为 StackTrace。
func (s *stack) StackTrace() StackTrace {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}
	return f
}

// callers 获取当前的调用栈，返回一个包含调用栈程序计数器的栈。
func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) // 获取调用栈（跳过当前函数）
	var st stack = pcs[0:n]
	return &st
}

// funcname 去掉函数名中的路径前缀，仅保留函数名称。
func funcname(name string) string {
	i := strings.LastIndex(name, "/") // 去除路径
	name = name[i+1:]
	i = strings.Index(name, ".") // 去除包名
	return name[i+1:]
}
