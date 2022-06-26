package rt

import (
	"path/filepath"
	"runtime"
)

// StackFrame struct.
type StackFrame struct {
	File     string
	FullFile string
	Function string
	Line     int
}

func frameForPC(pc uintptr) *StackFrame {
	f := runtime.FuncForPC(pc)
	if f == nil {
		return nil
	}
	file, line := f.FileLine(pc)
	return &StackFrame{
		File:     filepath.Base(file),
		FullFile: file,
		Function: f.Name(),
		Line:     line,
	}
}

// BackTrace returns info about function invocations on the calling goroutine's stack.
func BackTrace(skip int) []*StackFrame {
	pc := make([]uintptr, 32)
	n := runtime.Callers(2+skip, pc)
	stack := make([]*StackFrame, n)

	for i := 0; i < n; i++ {
		stack[i] = frameForPC(pc[i])
	}

	return stack
}

// Caller ...
func Caller(skip int) *StackFrame {
	pc, _, _, _ := runtime.Caller(2 + skip)
	return frameForPC(pc)
}
