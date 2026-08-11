// Package main is the TinyGo guest for the crypto-fakers plugin:
//
//	$.fake.eth_address  — 0x-prefixed 40-hex Ethereum-style address
//	$.fake.tx_hash      — 0x-prefixed 64-hex transaction hash
//	$.fake.btc_address  — base58 P2PKH-style address (starts with 1)
//
// Synthetic test data (no real key derivation). Freestanding, zero-allocation.
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

const hexlower = "0123456789abcdef"
const base58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

//export faker
func faker(ptr, length uint32) uint64 {
	name := matchName(memAt(ptr, length))
	seed(name)
	var d [80]byte
	n := 0
	switch {
	case bytesEq(name, "eth_address"):
		d[0], d[1] = '0', 'x'
		n = 2 + fill(d[2:], 40, hexlower)
	case bytesEq(name, "tx_hash"):
		d[0], d[1] = '0', 'x'
		n = 2 + fill(d[2:], 64, hexlower)
	case bytesEq(name, "btc_address"):
		d[0] = '1'
		n = 1 + fill(d[1:], 33, base58)
	default:
		return writeValue(d[:0])
	}
	return writeValue(d[:n])
}

func fill(d []byte, n int, alphabet string) int {
	m := uint64(len(alphabet))
	for i := 0; i < n; i++ {
		d[i] = alphabet[next()%m]
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

var rngState uint64 = 0x100000001B3
var callCounter uint64

func seed(name []byte) {
	callCounter++
	rngState ^= callCounter * 0x9E3779B97F4A7C15
	for i := 0; i < len(name); i++ {
		rngState = rngState*1099511628211 ^ uint64(name[i])
	}
}
func next() uint64 {
	rngState = rngState*6364136223846793005 + 1442695040888963407
	return rngState >> 17
}
