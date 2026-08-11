// Package main is the TinyGo guest for the ru-fakers Mockarty plugin.
//
// It is compiled to a freestanding WebAssembly module (no WASI, no imports)
// with:
//
//	tinygo build -target=wasm-unknown -no-debug -o ../faker.wasm .
//
// and loaded by Mockarty's WASM extension layer at the "faker-provider" point.
// The host runs it in a deny-by-default sandbox: the module can touch nothing
// but its own linear memory.
//
// # The Mockarty guest ABI (stable public contract)
//
//	memory                        exported linear memory (TinyGo does this for us)
//	mk_alloc(size i32) -> i32     returns a guest pointer for the input
//	<fn>(ptr i32, len i32) -> i64 packs (outPtr<<32 | outLen) of the JSON result
//
// For the faker point the host passes {"name":"<faker>"} and expects
// {"value":"<string>"} back. One exported function serves every faker name.
//
// Built only under TinyGo (build tag `tinygo`): the guest does WASM
// linear-memory pointer arithmetic (uintptr→unsafe.Pointer) that the host `go
// vet` flags as "possible misuse" — correct there, a false positive here. A
// trivial host-side stub (stub.go, `!tinygo`) keeps the module buildable/vettable
// on the host without pulling this file in.
//
//go:build tinygo

// # Zero allocation
//
// A freestanding TinyGo module ships with a leaking GC (no scheduler to run a
// real one), so anything that allocates on the heap — make(), string
// concatenation, []byte(...) — is never reclaimed and, after enough calls, the
// growing heap marches over the module's own buffers and corrupts output. So
// this guest allocates nothing on the hot path: fixed buffers, results written
// byte-by-byte. It stays correct for unbounded calls.
package main

import "unsafe"

func main() {}

// Fixed, separate input and output buffers. The host calls mk_alloc once per
// call and reads the result at the returned pointer, so fixed slots are correct
// and cannot overlap. Each module INSTANCE has its own memory and the host runs
// one call at a time per instance, so no locking is needed.
var inBuf [256]byte
var outBuf [64]byte

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

// ru_faker reads {"name":...}, generates the requested identifier straight into
// a fixed digit buffer, and writes {"value":"<digits>"} into outBuf with no heap
// allocation.
//
//export ru_faker
func ru_faker(ptr, length uint32) uint64 { //nolint:revive // ABI export name
	name := memAt(ptr, length)
	seed(name)

	var d [13]byte // digit scratch (max length: OGRN = 13)
	who := matchName(name)
	n := 0
	switch {
	case bytesEq(who, "ru_inn"):
		n = genINN10(&d)
	case bytesEq(who, "ru_snils"):
		n = genSNILS(&d)
	case bytesEq(who, "ru_ogrn"):
		n = genOGRN(&d)
	case bytesEq(who, "ru_kpp"):
		n = genKPP(&d)
	default:
		return writeValue(d[:0]) // unknown → empty value, host falls through
	}
	return writeValue(d[:n])
}

// writeValue emits {"value":"<digits>"} into outBuf and returns the packed ptr/len.
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
	return (uint64(ptrOf(&outBuf[0])) << 32) | uint64(i)
}

// --- tiny deterministic PRNG (no entropy source in a freestanding module) ---
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

// --- Russian identifier generators, writing digits into d, returning length ---

// genINN10 — 10-digit legal-entity INN. Check digit over 9 with canonical
// weights, sum mod 11 mod 10.
func genINN10(d *[13]byte) int {
	var w = [9]int{2, 4, 10, 3, 5, 9, 4, 6, 8}
	sum := 0
	for i := 0; i < 9; i++ {
		v := nextDigit()
		d[i] = byte('0' + v)
		sum += v * w[i]
	}
	d[9] = byte('0' + (sum%11)%10)
	return 10
}

// genSNILS — 9 body digits weighted 9..1 + a 2-digit control (sum mod 101, with
// 100/101 → 00).
func genSNILS(d *[13]byte) int {
	sum := 0
	for i := 0; i < 9; i++ {
		v := nextDigit()
		d[i] = byte('0' + v)
		sum += v * (9 - i)
	}
	ctrl := sum % 101
	if ctrl == 100 || ctrl == 101 {
		ctrl = 0
	}
	d[9] = byte('0' + ctrl/10)
	d[10] = byte('0' + ctrl%10)
	return 11
}

// genOGRN — 13 digits; the 13th is the 12-digit number mod 11, then mod 10.
func genOGRN(d *[13]byte) int {
	mod := 0
	for i := 0; i < 12; i++ {
		v := nextDigit()
		d[i] = byte('0' + v)
		mod = (mod*10 + v) % 11
	}
	d[12] = byte('0' + mod%10)
	return 13
}

// genKPP — 4-digit tax authority + "01" reason + 3-digit serial.
func genKPP(d *[13]byte) int {
	for i := 0; i < 4; i++ {
		d[i] = byte('0' + nextDigit())
	}
	d[4], d[5] = '0', '1'
	for i := 6; i < 9; i++ {
		d[i] = byte('0' + nextDigit())
	}
	return 9
}

// matchName extracts the value of "name" from the flat {"name":"..."} input,
// returning a sub-slice (no allocation). Good enough for the host-generated input.
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
