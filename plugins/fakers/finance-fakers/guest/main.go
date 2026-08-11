// Package main is the TinyGo guest for the finance-fakers Mockarty plugin — a
// second WASM faker example alongside ru-fakers. Build:
//
//	tinygo build -target=wasm-unknown -no-debug -o ../faker.wasm .
//
// Fakers: $.fake.credit_card (16 digits, valid Luhn) and $.fake.iban (a
// plausible DE IBAN). Zero-allocation on the hot path (see ru-fakers for the
// why: a freestanding TinyGo module has a leaking GC).
//
//go:build tinygo

package main

import "unsafe"

func main() {}

var inBuf [256]byte
var outBuf [64]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 {
	_ = size
	return uint32(uintptr(unsafe.Pointer(&inBuf[0])))
}

func memAt(ptr, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

//export fin_faker
func fin_faker(ptr, length uint32) uint64 { //nolint:revive // ABI export name
	name := memAt(ptr, length)
	seed(name)
	var d [34]byte
	n := 0
	who := matchName(name)
	switch {
	case bytesEq(who, "credit_card"):
		n = genCard(&d)
	case bytesEq(who, "iban"):
		n = genIBAN(&d)
	}
	return writeValue(d[:n])
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
	return (uint64(uint32(uintptr(unsafe.Pointer(&outBuf[0])))) << 32) | uint64(i)
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

func nextDigit() int {
	rngState = rngState*6364136223846793005 + 1442695040888963407
	return int((rngState >> 33) % 10)
}

// genCard — 16 digits with a valid Luhn check digit (test PAN, never a real card).
func genCard(d *[34]byte) int {
	for i := 0; i < 15; i++ {
		d[i] = byte('0' + nextDigit())
	}
	// Luhn check digit over the first 15.
	sum := 0
	for i := 0; i < 15; i++ {
		v := int(d[i] - '0')
		// From the rightmost of the 15 (position 14), every second digit doubles.
		if (15-i)%2 == 1 {
			v *= 2
			if v > 9 {
				v -= 9
			}
		}
		sum += v
	}
	d[15] = byte('0' + (10-sum%10)%10)
	return 16
}

// genIBAN — "DE" + 20 digits (plausible shape; not a checksum-valid IBAN).
func genIBAN(d *[34]byte) int {
	d[0], d[1] = 'D', 'E'
	for i := 2; i < 22; i++ {
		d[i] = byte('0' + nextDigit())
	}
	return 22
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

func indexOf(b []byte, sub string) int {
	n, m := len(b), len(sub)
	for i := 0; i+m <= n; i++ {
		ok := true
		for j := 0; j < m; j++ {
			if b[i+j] != sub[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}
