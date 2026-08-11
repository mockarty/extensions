// Package main is the TinyGo guest for the datetime-fakers plugin:
//
//	$.fake.iso_date      — YYYY-MM-DD
//	$.fake.iso_datetime  — YYYY-MM-DDThh:mm:ssZ (UTC)
//	$.fake.cron          — a 5-field cron expression
//	$.fake.duration      — an ISO-8601 duration like PT2H30M
//
// A freestanding module has no clock, so values are PRNG-derived within plausible
// ranges (deterministic per server start). Zero-allocation (see ru-fakers).
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
	case bytesEq(name, "iso_date"):
		n = writeDate(d[:])
	case bytesEq(name, "iso_datetime"):
		n = writeDate(d[:])
		d[n] = 'T'
		n++
		n += writeTime(d[n:])
		d[n] = 'Z'
		n++
	case bytesEq(name, "cron"):
		n = writeCron(d[:])
	case bytesEq(name, "duration"):
		n = writeDuration(d[:])
	default:
		return writeValue(d[:0])
	}
	return writeValue(d[:n])
}

func writeDate(d []byte) int {
	y := 2020 + nextN(11) // 2020..2030
	mo := 1 + nextN(12)
	day := 1 + nextN(28) // safe for every month
	i := put4(d, y)
	d[i] = '-'
	i++
	i += put2(d[i:], mo)
	d[i] = '-'
	i++
	i += put2(d[i:], day)
	return i
}
func writeTime(d []byte) int {
	i := put2(d, nextN(24))
	d[i] = ':'
	i++
	i += put2(d[i:], nextN(60))
	d[i] = ':'
	i++
	i += put2(d[i:], nextN(60))
	return i
}
func writeCron(d []byte) int {
	// minute hour dom mon dow  (each a plain value for simplicity/validity)
	i := putInt(d, nextN(60))
	i += sp(d[i:])
	i += putInt(d[i:], nextN(24))
	i += sp(d[i:])
	i += putInt(d[i:], 1+nextN(28))
	i += sp(d[i:])
	i += putInt(d[i:], 1+nextN(12))
	i += sp(d[i:])
	i += putInt(d[i:], nextN(7))
	return i
}
func writeDuration(d []byte) int {
	d[0], d[1] = 'P', 'T'
	i := 2
	i += putInt(d[i:], 1+nextN(12))
	d[i] = 'H'
	i++
	i += putInt(d[i:], nextN(60))
	d[i] = 'M'
	i++
	return i
}
func sp(d []byte) int { d[0] = ' '; return 1 }

func put2(d []byte, v int) int { d[0] = byte('0' + v/10%10); d[1] = byte('0' + v%10); return 2 }
func put4(d []byte, v int) int {
	d[0] = byte('0' + v/1000%10)
	d[1] = byte('0' + v/100%10)
	d[2] = byte('0' + v/10%10)
	d[3] = byte('0' + v%10)
	return 4
}
func putInt(d []byte, v int) int {
	if v == 0 {
		d[0] = '0'
		return 1
	}
	var tmp [10]byte
	n := 0
	for v > 0 {
		tmp[n] = byte('0' + v%10)
		v /= 10
		n++
	}
	for i := 0; i < n; i++ {
		d[i] = tmp[n-1-i]
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

var rngState uint64 = 0xC0FFEE123456789
var callCounter uint64

func seed(name []byte) {
	callCounter++
	rngState ^= callCounter * 0x9E3779B97F4A7C15
	for i := 0; i < len(name); i++ {
		rngState = rngState*1099511628211 ^ uint64(name[i])
	}
}
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
