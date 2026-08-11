// Package main is the TinyGo guest for the redaction-transformer plugin. Two
// response transformers, each its own exported function (the response-transformer
// point passes only the body, so one function per key):
//
//	transform:mask_digits  — every ASCII digit in the body becomes '*'
//	transform:base64_wrap  — body is base64-encoded into {"b64":"<...>"}
//
// Freestanding, zero-allocation: fixed arenas sized to the host input cap.
//
//go:build tinygo

package main

import "unsafe"

func main() {}

const maxInput = 1 << 20

var inBuf [maxInput]byte
var outBuf [maxInput*2 + 64]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 { _ = size; return ptrOf(&inBuf[0]) }
func ptrOf(b *byte) uint32       { return uint32(uintptr(unsafe.Pointer(b))) }
func body(ptr, length uint32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), int(length))
}
func pack(n int) uint64 { return uint64(ptrOf(&outBuf[0]))<<32 | uint64(n) }

//export mask_digits
func maskDigits(ptr, length uint32) uint64 {
	src := body(ptr, length)
	for i := 0; i < len(src); i++ {
		c := src[i]
		if c >= '0' && c <= '9' {
			outBuf[i] = '*'
		} else {
			outBuf[i] = c
		}
	}
	return pack(len(src))
}

const b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

//export base64_wrap
func base64Wrap(ptr, length uint32) uint64 {
	src := body(ptr, length)
	const prefix = `{"b64":"`
	const suffix = `"}`
	// Guard: enveloped base64 must fit outBuf.
	if len(prefix)+((len(src)+2)/3)*4+len(suffix) > len(outBuf) {
		copy(outBuf[:len(src)], src)
		return pack(len(src))
	}
	i := 0
	for j := 0; j < len(prefix); j++ {
		outBuf[i] = prefix[j]
		i++
	}
	// standard base64 with padding
	n := len(src)
	k := 0
	for ; k+3 <= n; k += 3 {
		v := uint32(src[k])<<16 | uint32(src[k+1])<<8 | uint32(src[k+2])
		outBuf[i] = b64[(v>>18)&63]
		outBuf[i+1] = b64[(v>>12)&63]
		outBuf[i+2] = b64[(v>>6)&63]
		outBuf[i+3] = b64[v&63]
		i += 4
	}
	switch n - k {
	case 1:
		v := uint32(src[k]) << 16
		outBuf[i] = b64[(v>>18)&63]
		outBuf[i+1] = b64[(v>>12)&63]
		outBuf[i+2] = '='
		outBuf[i+3] = '='
		i += 4
	case 2:
		v := uint32(src[k])<<16 | uint32(src[k+1])<<8
		outBuf[i] = b64[(v>>18)&63]
		outBuf[i+1] = b64[(v>>12)&63]
		outBuf[i+2] = b64[(v>>6)&63]
		outBuf[i+3] = '='
		i += 4
	}
	for j := 0; j < len(suffix); j++ {
		outBuf[i] = suffix[j]
		i++
	}
	return pack(i)
}
