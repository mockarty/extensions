// Package main is the TinyGo guest for the barcode-fakers plugin. It adds
// dynamic-response generators for valid retail barcodes:
//
//	$.fake.ean13   — 13-digit EAN with correct check digit
//	$.fake.upc     — 12-digit UPC-A with correct check digit
//	$.fake.isbn13  — 978-prefixed 13-digit ISBN with correct check digit
//
// Freestanding WASM (no WASI, no imports), zero-allocation — a leaking GC would
// corrupt output after enough calls, so the hot path uses only fixed buffers.
// See examples/plugins/ru-fakers for the full rationale.
//
//go:build tinygo

package main

import "unsafe"

func main() {}

var inBuf [256]byte
var outBuf [64]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 { _ = size; return ptrOf(&inBuf[0]) }

func ptrOf(b *byte) uint32 { return uint32(uintptr(unsafe.Pointer(b))) }
func memAt(ptr, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

//export faker
func faker(ptr, length uint32) uint64 {
	name := matchName(memAt(ptr, length))
	seed(name)
	var d [13]byte
	n := 0
	switch {
	case bytesEq(name, "ean13"):
		n = genGTIN(&d, 13)
	case bytesEq(name, "upc"):
		n = genGTIN(&d, 12)
	case bytesEq(name, "isbn13"):
		d[0], d[1], d[2] = '9', '7', '8'
		n = genGTINPrefixed(&d, 13, 3)
	default:
		return writeValue(d[:0])
	}
	return writeValue(d[:n])
}

// genGTIN fills length-1 random digits then the GS1 mod-10 check digit.
func genGTIN(d *[13]byte, length int) int { return genGTINPrefixed(d, length, 0) }

// genGTINPrefixed keeps the first `fixed` digits (already set) and fills the rest
// up to length-1, then appends the GS1 check digit. GS1 weights alternate 1,3
// from the RIGHT of the payload; sum, then (10 - sum%10) % 10.
func genGTINPrefixed(d *[13]byte, length, fixed int) int {
	for i := fixed; i < length-1; i++ {
		d[i] = byte('0' + nextDigit())
	}
	sum := 0
	// payload is d[0 .. length-2]; rightmost payload digit gets weight 3.
	for i := 0; i < length-1; i++ {
		v := int(d[i] - '0')
		// position from the right within the payload
		if (length-2-i)%2 == 0 {
			sum += v * 3
		} else {
			sum += v
		}
	}
	check := (10 - sum%10) % 10
	d[length-1] = byte('0' + check)
	return length
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

var rngState uint64 = 0x2545F4914F6CDD1D
var callCounter uint64

func seed(name []byte) {
	callCounter++
	rngState ^= callCounter * 0x9E3779B97F4A7C15
	for i := 0; i < len(name); i++ {
		rngState = rngState*1099511628211 ^ uint64(name[i])
	}
}

func nextDigit() int {
	rngState = rngState*6364136223846793005 + 1442695040888963407
	return int((rngState >> 33) % 10)
}
