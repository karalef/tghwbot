package logger

import (
	"fmt"
)

// New creates new logger instance.
func New(w *Writer, name string) *Logger {
	return &Logger{
		out:  w,
		name: name,
	}
}

// Logger struct.
type Logger struct {
	out  *Writer
	name string
}

// Child ...
func (l *Logger) Child(name string) *Logger {
	return New(l.out, l.name+"::"+name)
}

const namecol = green

func (l *Logger) log(pref string, pcol ansiColor, f string, v ...interface{}) {
	if f == "" {
		return
	}

	w := l.out
	w.lock()
	w.writeBuf(pref, pcol)
	if l.name != "" {
		w.writeBuf(l.name, namecol)
	}
	w.writeBuf(fmt.Sprintf(f, v...), white)
	w.write()
	w.unlock()
}

// Info ...
func (l *Logger) Info(f string, v ...interface{}) {
	l.log("INFO", white, f, v...)
}

// Warn ...
func (l *Logger) Warn(f string, v ...interface{}) {
	l.log("WARN", yellow, f, v...)
}

// Error ...
func (l *Logger) Error(f string, v ...interface{}) {
	l.log("ERROR", red, f, v...)
}
