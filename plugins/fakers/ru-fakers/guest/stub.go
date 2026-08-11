// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

//go:build !tinygo

// This stub exists only so the guest module builds and vets cleanly on the host
// toolchain. The real guest (main.go) is TinyGo-only — it does WASM
// linear-memory pointer arithmetic the host `go vet` would flag. Build the
// module with TinyGo (./build.sh); the host never runs this code.
package main

func main() {}
