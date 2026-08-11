// Package main is the TinyGo guest for the condition-matcher example plugin.
// It adds two request matchers usable as a mock condition's assertAction:
//
//	plugin:luhn          — the value passes the Luhn checksum (card/IMEI numbers)
//	plugin:divisible_by  — the value is an integer divisible by the operand
//
// The host passes {"key","actual","expected"} and expects {"match":bool}. Same
// sandbox + guest ABI as the faker point; zero-allocation (see ru-fakers).
//
//go:build tinygo

package main

import "unsafe"

func main() {}

const maxInput = 1 << 20

var inBuf [maxInput]byte
var outBuf [32]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 {
	_ = size
	return ptrOf(&inBuf[0])
}

func ptrOf(b *byte) uint32 { return uint32(uintptr(unsafe.Pointer(b))) }

func memAt(ptr, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

//export match
func match(ptr, length uint32) uint64 {
	in := memAt(ptr, length)
	key := field(in, `"key"`)
	actual := field(in, `"actual"`)
	expected := field(in, `"expected"`)

	ok := false
	switch {
	case bytesEq(key, "luhn"):
		ok = luhnValid(actual)
	case bytesEq(key, "divisible_by"):
		ok = divisibleBy(actual, expected)
	}
	return writeMatch(ok)
}

// luhnValid reports whether the digit string passes the Luhn checksum. Non-digit
// characters make it fail. Empty input fails.
func luhnValid(d []byte) bool {
	if len(d) == 0 {
		return false
	}
	sum := 0
	dbl := false
	for i := len(d) - 1; i >= 0; i-- {
		c := d[i]
		if c < '0' || c > '9' {
			return false
		}
		v := int(c - '0')
		if dbl {
			v *= 2
			if v > 9 {
				v -= 9
			}
		}
		sum += v
		dbl = !dbl
	}
	return sum%10 == 0
}

// divisibleBy reports whether actual (integer) is divisible by expected (integer).
// A zero or non-numeric operand never matches.
func divisibleBy(actual, expected []byte) bool {
	a, aok := parseInt(actual)
	b, bok := parseInt(expected)
	if !aok || !bok || b == 0 {
		return false
	}
	return a%b == 0
}

func parseInt(b []byte) (int64, bool) {
	if len(b) == 0 {
		return 0, false
	}
	i := 0
	neg := false
	if b[0] == '-' {
		neg = true
		i = 1
		if len(b) == 1 {
			return 0, false
		}
	}
	var n int64
	for ; i < len(b); i++ {
		if b[i] < '0' || b[i] > '9' {
			return 0, false
		}
		n = n*10 + int64(b[i]-'0')
	}
	if neg {
		n = -n
	}
	return n, true
}

// writeMatch emits {"match":true} or {"match":false}.
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

// field extracts a quoted string field value by key (e.g. `"actual"`) from a
// flat JSON object. Returns an empty slice when absent. Minimal: values are
// assumed not to contain escaped quotes (the host sends plain scalars).
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
