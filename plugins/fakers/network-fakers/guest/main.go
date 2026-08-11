// Package main is the TinyGo guest for the network-fakers plugin:
//
//	$.fake.cidr        — IPv4 CIDR block a.b.c.d/nn (nn 8..30)
//	$.fake.private_ip  — RFC1918 private IPv4 (10.x / 172.16-31.x / 192.168.x)
//	$.fake.port        — TCP port 1024..65535
//	$.fake.mac_oui     — MAC address with a locally-administered OUI
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
	var d [24]byte
	n := 0
	switch {
	case bytesEq(name, "cidr"):
		n = putIP(d[:], nextN(256), nextN(256), nextN(256), nextN(256))
		d[n] = '/'
		n++
		n += putInt(d[n:], 8+nextN(23)) // 8..30
	case bytesEq(name, "private_ip"):
		switch nextN(3) {
		case 0:
			n = putIP(d[:], 10, nextN(256), nextN(256), 1+nextN(254))
		case 1:
			n = putIP(d[:], 172, 16+nextN(16), nextN(256), 1+nextN(254))
		default:
			n = putIP(d[:], 192, 168, nextN(256), 1+nextN(254))
		}
	case bytesEq(name, "port"):
		n = putInt(d[:], 1024+nextN(64512))
	case bytesEq(name, "mac_oui"):
		n = putMAC(d[:])
	default:
		return writeValue(d[:0])
	}
	return writeValue(d[:n])
}

func putIP(d []byte, a, b, c, e int) int {
	i := putInt(d, a)
	d[i] = '.'
	i++
	i += putInt(d[i:], b)
	d[i] = '.'
	i++
	i += putInt(d[i:], c)
	d[i] = '.'
	i++
	i += putInt(d[i:], e)
	return i
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

const hexd = "0123456789abcdef"

func putMAC(d []byte) int {
	// first octet: locally-administered + unicast (bit1 set, bit0 clear) → x2/x6/xA/xE
	first := (nextN(16) << 4) | 0x2
	i := 0
	i += putHex2(d[i:], first)
	for k := 0; k < 5; k++ {
		d[i] = ':'
		i++
		i += putHex2(d[i:], nextN(256))
	}
	return i
}
func putHex2(d []byte, v int) int {
	d[0] = hexd[(v>>4)&0xf]
	d[1] = hexd[v&0xf]
	return 2
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

var rngState uint64 = 0xA5A5A5A5F0F0F0F0
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
