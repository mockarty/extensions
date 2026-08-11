// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package redactiontransformer

import (
	"context"
	"encoding/base64"
	"os"
	"testing"

	"mockarty/internal/wasmext"
)

func load(t *testing.T) *wasmext.Module {
	t.Helper()
	w, err := os.ReadFile("transform.wasm")
	if err != nil {
		t.Fatalf("read transform.wasm (./build.sh): %v", err)
	}
	m, err := wasmext.Compile(context.Background(), "redaction-transformer", w, wasmext.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestMaskDigits(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	cases := []struct{ in, want string }{
		{`{"card":"4111111111111111","name":"Ann"}`, `{"card":"****************","name":"Ann"}`},
		{`no digits here`, `no digits here`},
		{`{"pin":1234}`, `{"pin":****}`},
	}
	for _, c := range cases {
		out, err := m.Call(context.Background(), "mask_digits", []byte(c.in))
		if err != nil {
			t.Fatalf("mask_digits: %v", err)
		}
		if string(out) != c.want {
			t.Fatalf("mask_digits(%q) = %q, want %q", c.in, out, c.want)
		}
	}
}

func TestBase64Wrap(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	for _, in := range []string{`{"a":1}`, ``, `x`, `hello world!`, `{"nested":{"k":[1,2,3]}}`} {
		out, err := m.Call(context.Background(), "base64_wrap", []byte(in))
		if err != nil {
			t.Fatalf("base64_wrap: %v", err)
		}
		want := `{"b64":"` + base64.StdEncoding.EncodeToString([]byte(in)) + `"}`
		if string(out) != want {
			t.Fatalf("base64_wrap(%q) = %q, want %q", in, out, want)
		}
	}
}
