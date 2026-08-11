// Package main is the TinyGo guest for the eu-fakers plugin:
//
//	$.fake.iban  — German-format IBAN (DE + 2 check digits + 18-digit BBAN) with a
//	               valid ISO 7064 mod-97 checksum
//	$.fake.bic   — 8-char SWIFT/BIC (4 letters bank + 2 letters country + 2 loc)
//	$.fake.vat_de — German VAT number DE + 9 digits
//
// Freestanding, zero-allocation (see ru-fakers).
//
//go:build tinygo

package main

import "unsafe"

func main() {}

var inBuf [256]byte
var outBuf [40]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 { _ = size; return ptrOf(&inBuf[0]) }
func ptrOf(b *byte) uint32       { return uint32(uintptr(unsafe.Pointer(b))) }
func memAt(ptr, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), int(length))
}

//export faker
func faker(ptr, length uint32) uint64 {
	name := matchName(memAt(ptr, length))
	seed(name)
	var d [40]byte
	n := 0
	switch {
	case bytesEq(name, "iban"):
		n = genIBAN(&d)
	case bytesEq(name, "bic"):
		n = genBIC(&d)
	case bytesEq(name, "vat_de"):
		n = genVATDE(&d)
	default:
		return writeValue(d[:0])
	}
	return writeValue(d[:n])
}

// genIBAN builds a DE IBAN: "DE" + check(2) + 18 BBAN digits, where the check
// digits make the whole thing ≡ 1 (mod 97) per ISO 13616 / ISO 7064.
func genIBAN(d *[40]byte) int {
	var bban [18]byte
	for i := 0; i < 18; i++ {
		bban[i] = byte('0' + nextDigit())
	}
	// Rearranged string for the checksum = BBAN + "DE00" numeric = BBAN + "131400".
	// Compute mod 97 over BBAN digits then the trailing "131400".
	rem := 0
	for i := 0; i < 18; i++ {
		rem = (rem*10 + int(bban[i]-'0')) % 97
	}
	for _, c := range []byte("131400") {
		rem = (rem*10 + int(c-'0')) % 97
	}
	check := 98 - rem
	d[0], d[1] = 'D', 'E'
	d[2] = byte('0' + check/10)
	d[3] = byte('0' + check%10)
	for i := 0; i < 18; i++ {
		d[4+i] = bban[i]
	}
	return 22
}

const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genBIC(d *[40]byte) int {
	for i := 0; i < 4; i++ {
		d[i] = upper[next()%26]
	}
	d[4], d[5] = 'D', 'E' // country
	// location: alphanumeric, avoid 0/O ambiguity is not required here
	for i := 6; i < 8; i++ {
		if next()%2 == 0 {
			d[i] = upper[next()%26]
		} else {
			d[i] = byte('0' + nextDigit())
		}
	}
	return 8
}

func genVATDE(d *[40]byte) int {
	d[0], d[1] = 'D', 'E'
	for i := 0; i < 9; i++ {
		d[2+i] = byte('0' + nextDigit())
	}
	return 11
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

var rngState uint64 = 0x1EED1EED1EED1EED
var callCounter uint64

func seed(name []byte) {
	callCounter++
	rngState ^= callCounter * 0x9E3779B97F4A7C15
	for i := 0; i < len(name); i++ {
		rngState = rngState*1099511628211 ^ uint64(name[i])
	}
}
func nextDigit() int { return int(next() % 10) }
func next() uint64 {
	rngState = rngState*6364136223846793005 + 1442695040888963407
	return rngState >> 17
}
