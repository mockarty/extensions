// Package main is the TinyGo guest for the format-matchers plugin. Custom
// condition operators (assertAction plugin:<key>):
//
//	plugin:jwt   — value looks like a JWT (three non-empty base64url parts)
//	plugin:e164  — value is an E.164 phone (+ then 8..15 digits)
//	plugin:ipv4  — value is a dotted-quad IPv4 with each octet 0..255
//
// Input {"key","actual","expected"} → {"match":bool}. Freestanding, zero-alloc.
//
//go:build tinygo

package main

import "unsafe"

func main() {}

var inBuf [1 << 20]byte
var outBuf [32]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 { _ = size; return ptrOf(&inBuf[0]) }
func ptrOf(b *byte) uint32       { return uint32(uintptr(unsafe.Pointer(b))) }
func memAt(ptr, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), int(length))
}

//export match
func match(ptr, length uint32) uint64 {
	in := memAt(ptr, length)
	key := field(in, `"key"`)
	actual := field(in, `"actual"`)
	ok := false
	switch {
	case bytesEq(key, "jwt"):
		ok = isJWT(actual)
	case bytesEq(key, "e164"):
		ok = isE164(actual)
	case bytesEq(key, "ipv4"):
		ok = isIPv4(actual)
	}
	return writeMatch(ok)
}

func isJWT(s []byte) bool {
	parts, cur := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			if cur == 0 {
				return false
			}
			parts++
			cur = 0
			continue
		}
		if !isB64URL(s[i]) {
			return false
		}
		cur++
	}
	return parts == 2 && cur > 0
}
func isB64URL(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
}

func isE164(s []byte) bool {
	if len(s) < 9 || len(s) > 16 || s[0] != '+' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func isIPv4(s []byte) bool {
	octets, val, digits := 0, 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if digits == 0 || val > 255 {
				return false
			}
			octets++
			val, digits = 0, 0
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
		val = val*10 + int(c-'0')
		digits++
		if digits > 3 {
			return false
		}
	}
	return octets == 3 && digits > 0 && val <= 255
}

func writeMatch(ok bool) uint64 {
	s := `{"match":false}`
	if ok {
		s = `{"match":true}`
	}
	for i := 0; i < len(s); i++ {
		outBuf[i] = s[i]
	}
	return uint64(ptrOf(&outBuf[0]))<<32 | uint64(len(s))
}

func field(b []byte, key string) []byte {
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
