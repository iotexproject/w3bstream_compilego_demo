package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Note: In TinyGo "//export" on a func is actually an import!

// main is required for TinyGo to compile to Wasm.
func main() {}

// greet prints a greeting to the console.
func greet(name string) {
	log(fmt.Sprint("wasm >> ", greeting(name)))
}

// log a message to the console using _log.
func log(message string) {
	ptr, size := stringToPtr(message)
	_log(ptr, size)
}

// _log is a wasm import, prints a string
//
//go:wasm-module env
//export log
func _log(ptr uint32, size uint32)

// greeting gets a greeting for the name.
func greeting(name string) string {
	return fmt.Sprint("Hello, ", name, "!")
}

// _greet is wasm export, accepts a string pointer
//
//export greet
func _greet(ptr, size uint32) {
	name := ptrToString(ptr, size)
	greet(name)
}

// _greeting is wasm export that accepts a string pointer and returns a pointer/size pair packed into a uint64. (wasm 1.0 compatibility)
//
//export greeting
func _greeting(ptr, size uint32) (ptrSize uint64) {
	name := ptrToString(ptr, size)
	g := greeting(name)
	ptr, size = stringToPtr(g)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

// ptrToString returns a string from wasm compatible numeric types representing its pointer and length.
func ptrToString(ptr uint32, size uint32) string {
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // requires uintptr https://github.com/tinygo-org/tinygo/issues/1284
		Cap:  uintptr(size),
	}))
}

// stringToPtr returns a pointer and size pair for the given string
func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}
