package internal

import (
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// BackTrace returns info about function invocations on the calling goroutine's stack.
func BackTrace(skip int, length int) []runtime.Frame {
	pc := make([]uintptr, length)
	n := runtime.Callers(2+skip, pc[:])

	stack := make([]runtime.Frame, 0, n)
	frames := runtime.CallersFrames(pc[:n])

	for {
		f, more := frames.Next()
		stack = append(stack, f)
		if !more {
			break
		}
	}

	return stack
}

// Caller returns runtime frame with info about func invocation.
func Caller(skip int) runtime.Frame {
	return BackTrace(2+skip, 1)[0]
}

// FramesString ...
func FramesString(frames []runtime.Frame, trim bool) string {
	s := new(strings.Builder)
	if trim {
		TrimFramesPathPrefix(frames)
	}
	for i, f := range frames {
		s.WriteString(f.Function + "\n\t" + f.File)
		s.WriteString(":" + strconv.Itoa(f.Line))
		if i < len(frames)-1 {
			s.WriteByte('\n')
		}
	}
	return s.String()
}

var pathPrefix string
var pathPrefixOnce sync.Once

// TrimFramesPathPrefix ...
func TrimFramesPathPrefix(frames []runtime.Frame) {
	if len(frames) == 0 {
		return
	}
	pathPrefixOnce.Do(func() {
		f := make([]string, len(frames))
		for i := range frames {
			f[i] = frames[i].File
		}
		sort.Strings(f)

		for i := 0; i < len(f[0]); i++ {
			if f[0][i] != f[len(f)-1][i] {
				pathPrefix = f[0][:i]
				break
			}
		}
	})
	for i := range frames {
		frames[i].File = strings.Replace(frames[i].File, pathPrefix, ".../", 1)
	}
}
