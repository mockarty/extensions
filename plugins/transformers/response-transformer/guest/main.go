// Package main is the TinyGo guest for the response-transformer example plugin.
// It reads a mock's rendered response body and returns a transformed body — here
// it wraps the body in an envelope {"data":<body>,"transformedBy":"plugin"}.
//
// # Guest ABI (same as the faker point)
//
//	memory                        exported linear memory (TinyGo provides it)
//	mk_alloc(size i32) -> i32     returns a guest pointer for the input
//	transform(ptr i32, len i32) -> i64 packs (outPtr<<32 | outLen) of the result
//
// The host caps input at wasmext.Limits.MaxInputBytes (1MiB by default) and
// refuses to write more, so the fixed input arena below is always large enough —
// mk_alloc can safely ignore `size`. The output arena is sized to hold the
// largest input plus the envelope, so the transform never overflows.
//
// Built only under TinyGo (build tag `tinygo`): it does WASM linear-memory
// pointer arithmetic the host `go vet` would flag; a host stub (stub.go,
// `!tinygo`) keeps the package buildable/vettable on the host.
//
// # Zero allocation
//
// A freestanding TinyGo module ships with a leaking GC, so the hot path
// allocates nothing on the heap: fixed arenas, bytes copied in place. It stays
// correct for unbounded calls (see ru-fakers for the full rationale).
//
//go:build tinygo

package main

import "unsafe"

func main() {}

// maxInput mirrors the host's default Limits.MaxInputBytes (1MiB). The host
// refuses to write more, so inBuf is always large enough for what it hands us.
const maxInput = 1 << 20

// Fixed, separate input and output arenas — they can never overlap. outBuf holds
// the largest possible envelope: prefix + a full maxInput body + suffix.
var inBuf [maxInput]byte
var outBuf [maxInput + 64]byte

//export mk_alloc
func mkAlloc(size uint32) uint32 {
	_ = size // host guarantees size <= maxInput
	return ptrOf(&inBuf[0])
}

func ptrOf(b *byte) uint32 { return uint32(uintptr(unsafe.Pointer(b))) }

// transform receives the raw response body at (ptr,len) and returns the enveloped
// body {"data":<body>,"transformedBy":"plugin"}.
//
//export transform
func transform(ptr, length uint32) uint64 {
	body := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), int(length))
	const prefix = `{"data":`
	const suffix = `,"transformedBy":"plugin"}`
	total := len(prefix) + len(body) + len(suffix)
	if total > len(outBuf) {
		// Unreachable while the host cap matches maxInput, but stay fail-safe:
		// return the body unchanged rather than overflow.
		dst := unsafe.Slice((*byte)(unsafe.Pointer(&outBuf[0])), len(body))
		copy(dst, body)
		return pack(ptrOf(&outBuf[0]), uint32(len(body)))
	}
	i := 0
	i += copyStr(i, prefix)
	i += copyBytes(i, body)
	i += copyStr(i, suffix)
	return pack(ptrOf(&outBuf[0]), uint32(i))
}

func copyStr(at int, s string) int {
	for j := 0; j < len(s); j++ {
		outBuf[at+j] = s[j]
	}
	return len(s)
}

func copyBytes(at int, b []byte) int {
	for j := 0; j < len(b); j++ {
		outBuf[at+j] = b[j]
	}
	return len(b)
}

func pack(ptr, length uint32) uint64 { return uint64(ptr)<<32 | uint64(length) }
