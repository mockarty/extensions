// Package main is the TinyGo guest for the identifier-fakers plugin:
//
//	$.fake.ulid       — 26-char Crockford base32 identifier
//	$.fake.nanoid     — 21-char URL-safe identifier
//	$.fake.hex_token  — 32-char lowercase hex token
//
// Freestanding, zero-allocation (see ru-fakers).
//
//go:build tinygo

package main

import "unsafe"

func main() {}

var inBuf [256]byte
var outBuf [96]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 { _ = size; return ptrOf(&inBuf[0]) }
func ptrOf(b *byte) uint32       { return uint32(uintptr(unsafe.Pointer(b))) }
func memAt(ptr, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
const urlsafe = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
const hexlower = "0123456789abcdef"

//export faker
func faker(ptr, length uint32) uint64 {
	name := matchName(memAt(ptr, length))
	seed(name)
	var d [64]byte
	n := 0
	switch {
	case bytesEq(name, "ulid"):
		n = fill(d[:], 26, crockford)
	case bytesEq(name, "nanoid"):
		n = fill(d[:], 21, urlsafe)
	case bytesEq(name, "hex_token"):
		n = fill(d[:], 32, hexlower)
	default:
		return writeValue(d[:0])
	}
	return writeValue(d[:n])
}

func fill(d []byte, n int, alphabet string) int {
	m := len(alphabet)
	for i := 0; i < n; i++ {
		d[i] = alphabet[next()%uint64(m)]
	}
	return n
}

func writeValue(digits []byte) uint64 {
	const prefix = `{"value":"`
	const suffix = `"}`
	i := 0
	for j := 0; j < len(prefix); j++ {
		outBuf[i] = prefix[j]
		i++
	}
	for _, c := range digits {
		outBuf[i] = c
		i++
	}
	for j := 0; j < len(suffix); j++ {
		outBuf[i] = suffix[j]
		i++
	}
	return uint64(ptrOf(&outBuf[0]))<<32 | uint64(i)
}

func matchName(b []byte) []byte {
	const key = `"name"`
	i := indexOf(b, key)
	if i < 0 {
		return b[:0]
	}
	i += len(key)
	for i < len(b) && (b[i] == ' ' || b[i] == ':') {
		i++
	}
	if i >= len(b) || b[i] != '"' {
		return b[:0]
	}
	i++
	start := i
	for i < len(b) && b[i] != '"' {
		i++
	}
	return b[start:i]
}
func indexOf(b []byte, s string) int {
	if len(s) == 0 || len(b) < len(s) {
		return -1
	}
	for i := 0; i+len(s) <= len(b); i++ {
		j := 0
		for j < len(s) && b[i+j] == s[j] {
			j++
		}
		if j == len(s) {
			return i
		}
	}
	return -1
}
func bytesEq(b []byte, s string) bool {
	if len(b) != len(s) {
		return false
	}
	for i := 0; i < len(b); i++ {
		if b[i] != s[i] {
			return false
		}
	}
	return true
}

var rngState uint64 = 0x9E3779B97F4A7C15
var callCounter uint64

func seed(name []byte) {
	callCounter++
	rngState ^= callCounter * 0x2545F4914F6CDD1D
	for i := 0; i < len(name); i++ {
		rngState = rngState*1099511628211 ^ uint64(name[i])
	}
}
func next() uint64 {
	rngState = rngState*6364136223846793005 + 1442695040888963407
	return rngState >> 17
}
