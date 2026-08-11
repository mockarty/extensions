// Package main is the TinyGo guest for the us-fakers plugin:
//
//	$.fake.us_routing — 9-digit ABA routing number with a valid checksum
//	$.fake.ssn        — SSN AAA-GG-SSSS (synthetic, valid area/group/serial ranges)
//	$.fake.ein        — employer identification number NN-NNNNNNN
//
// Freestanding, zero-allocation (see ru-fakers).
//
//go:build tinygo

package main

import "unsafe"

func main() {}

var inBuf [256]byte
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

//export faker
func faker(ptr, length uint32) uint64 {
	name := matchName(memAt(ptr, length))
	seed(name)
	var d [16]byte
	n := 0
	switch {
	case bytesEq(name, "us_routing"):
		n = genRouting(&d)
	case bytesEq(name, "ssn"):
		n = genSSN(&d)
	case bytesEq(name, "ein"):
		n = genEIN(&d)
	default:
		return writeValue(d[:0])
	}
	return writeValue(d[:n])
}

// genRouting writes 8 random digits then the ABA check digit so that
// 3(d1+d4+d7)+7(d2+d5+d8)+(d3+d6+d9) ≡ 0 (mod 10).
func genRouting(d *[16]byte) int {
	var dig [9]int
	for i := 0; i < 8; i++ {
		dig[i] = nextDigit()
	}
	// weights per position 0..8: 3,7,1,3,7,1,3,7,1
	w := [9]int{3, 7, 1, 3, 7, 1, 3, 7, 1}
	sum := 0
	for i := 0; i < 8; i++ {
		sum += dig[i] * w[i]
	}
	// choose d9 so total ≡ 0 mod 10; d9 weight is 1
	dig[8] = (10 - sum%10) % 10
	for i := 0; i < 9; i++ {
		d[i] = byte('0' + dig[i])
	}
	return 9
}

func genSSN(d *[16]byte) int {
	// area 001..899 excluding 666; group 01..99; serial 0001..9999
	area := 1 + nextN(899)
	if area == 666 {
		area = 667
	}
	group := 1 + nextN(99)
	serial := 1 + nextN(9999)
	i := 0
	i += put3(d[i:], area)
	d[i] = '-'
	i++
	i += put2(d[i:], group)
	d[i] = '-'
	i++
	i += put4(d[i:], serial)
	return i
}

func genEIN(d *[16]byte) int {
	prefix := 10 + nextN(89) // 10..98
	rest := nextN(9999999)   // 0..9999999
	i := 0
	i += put2(d[i:], prefix)
	d[i] = '-'
	i++
	i += put7(d[i:], rest)
	return i
}

func put2(d []byte, v int) int { d[0] = byte('0' + v/10%10); d[1] = byte('0' + v%10); return 2 }
func put3(d []byte, v int) int {
	d[0] = byte('0' + v/100%10)
	d[1] = byte('0' + v/10%10)
	d[2] = byte('0' + v%10)
	return 3
}
func put4(d []byte, v int) int {
	for i := 3; i >= 0; i-- {
		d[i] = byte('0' + v%10)
		v /= 10
	}
	return 4
}
func put7(d []byte, v int) int {
	for i := 6; i >= 0; i-- {
		d[i] = byte('0' + v%10)
		v /= 10
	}
	return 7
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

var rngState uint64 = 0xDEADBEEFCAFEBABE
var callCounter uint64

func seed(name []byte) {
	callCounter++
	rngState ^= callCounter * 0x9E3779B97F4A7C15
	for i := 0; i < len(name); i++ {
		rngState = rngState*1099511628211 ^ uint64(name[i])
	}
}
func nextDigit() int { return int(next() % 10) }
func nextN(mod int) int {
	if mod <= 0 {
		return 0
	}
	return int(next() % uint64(mod))
}
func next() uint64 {
	rngState = rngState*6364136223846793005 + 1442695040888963407
	return rngState >> 17
}
