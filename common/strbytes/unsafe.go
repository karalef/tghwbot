package strbytes

import (
	"reflect"
	"unsafe"
)

// B2s converts bytes slice to string without allocation.
func B2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2b converts string to bytes slice without allocation.
func S2b(s string) []byte {
	b := *(*[]byte)(unsafe.Pointer(&s))
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = len(s)
	return b
}
