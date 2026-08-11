// Copyright (c) 2026 Mockarty. All rights reserved.
//go:build !tinygo

// Host-side stub so the guest module builds/vets on the host toolchain; the real
// guest (main.go) is TinyGo-only. See ru-fakers for the rationale.
package main

func main() {}
