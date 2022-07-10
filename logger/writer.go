package logger

import (
	"io"
	"os"
	"sync"
	"time"
)

// DefaultWriter var.
var DefaultWriter = NewWriter(os.Stderr, false)

// NewWriter makes new writer.
func NewWriter(w io.Writer, useColor bool) *Writer {
	return &Writer{
		out:      w,
		useColor: useColor,
	}
}

// NewWriterFile ...
func NewWriterFile(path string) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	return NewWriter(f, false), nil
}

// Writer struct.
type Writer struct {
	last     time.Time
	out      io.Writer
	buf      []byte
	mut      sync.Mutex
	useColor bool
}

func (w *Writer) lock() {
	w.mut.Lock()
}

func (w *Writer) unlock() {
	w.mut.Unlock()
}

func (w *Writer) writeBuf(text string, col ansiColor) {
	if w.useColor {
		text = col.wrap(text)
	}
	w.buf = append(w.buf, text...)
	w.buf = append(w.buf, ' ')
}

const timecol = magenta

func (w *Writer) write() {
	t := time.Now()
	if t.Minute() != w.last.Minute() {
		return
	}
	w.last = t
	f := t.Format("02.01.2006 15:04\n")
	if w.useColor {
		f = timecol.wrap(f)
	}
	w.out.Write([]byte(f))

	w.buf = append(w.buf, '\n')
	w.out.Write(w.buf)
	w.buf = w.buf[:0]
}

type ansiColor string

func (c ansiColor) wrap(text string) string {
	return string(c) + text + string(resetColor)
}

const (
	red        ansiColor = "\033[1;31m"
	green      ansiColor = "\033[1;32m"
	yellow     ansiColor = "\033[1;33m"
	blue       ansiColor = "\033[1;34m"
	magenta    ansiColor = "\033[1;35m"
	cyan       ansiColor = "\033[1;36m"
	white      ansiColor = "\033[1;37m"
	resetColor ansiColor = "\033[0m"
)
