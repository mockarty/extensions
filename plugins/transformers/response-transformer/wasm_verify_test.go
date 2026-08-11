// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

// Package responsetransformer holds the committed verification for the
// response-transformer example plugin: it loads the built transform.wasm through
// the real host WASM layer and asserts the guest envelopes a response body
// exactly as the product's response path will. This is the dogfood gate — the
// shipped .wasm must always pass the same sandbox the product runs plugins in.
package responsetransformer

import (
	"context"
	"os"
	"strings"
	"testing"

	"mockarty/internal/wasmext"
)

func loadModule(t *testing.T) *wasmext.Module {
	t.Helper()
	wasm, err := os.ReadFile("transform.wasm")
	if err != nil {
		t.Fatalf("read transform.wasm (build it: ./build.sh): %v", err)
	}
	mod, err := wasmext.Compile(context.Background(), "response-transformer/transform.wasm", wasm, wasmext.Limits{})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return mod
}

func TestResponseTransformerEnvelope(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())

	cases := []struct {
		name string
		body string
	}{
		{"object", `{"orig":true,"n":7}`},
		{"array", `[1,2,3]`},
		{"string", `"hello"`},
		{"empty-object", `{}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := mod.Call(context.Background(), "transform", []byte(tc.body))
			if err != nil {
				t.Fatalf("call: %v", err)
			}
			got := string(out)
			want := `{"data":` + tc.body + `,"transformedBy":"plugin"}`
			if got != want {
				t.Fatalf("envelope mismatch:\n got  %s\n want %s", got, want)
			}
			// The original body must survive verbatim inside the envelope.
			if !strings.Contains(got, tc.body) {
				t.Fatalf("original body not preserved in %s", got)
			}
		})
	}
}

// TestResponseTransformerLargeBody asserts a large body that still fits the host
// input cap is enveloped correctly — no truncation, no corruption.
func TestResponseTransformerLargeBody(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())
	body := `"` + strings.Repeat("x", 256*1024) + `"` // 256KiB JSON string, well under the 1MiB cap
	out, err := mod.Call(context.Background(), "transform", []byte(body))
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	want := `{"data":` + body + `,"transformedBy":"plugin"}`
	if string(out) != want {
		t.Fatalf("large body envelope mismatch (len got=%d want=%d)", len(out), len(want))
	}
}

// TestResponseTransformerInputCap asserts the host refuses a body larger than the
// input cap instead of overflowing guest memory — the caller then serves the
// body unchanged.
func TestResponseTransformerInputCap(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())
	over := make([]byte, (1<<20)+1) // one byte past the 1MiB default cap
	for i := range over {
		over[i] = 'x'
	}
	if _, err := mod.Call(context.Background(), "transform", over); err == nil {
		t.Fatal("expected the host to reject input past the cap")
	}
}
