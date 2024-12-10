package helper

import (
	"bytes"
	"unsafe"
)

func Bytes2String(b []byte) string {
	trimmedData := bytes.TrimRight(b, "\x00")
	return *(*string)(unsafe.Pointer(&trimmedData))
}
