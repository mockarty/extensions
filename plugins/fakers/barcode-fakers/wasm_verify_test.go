// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

// Package barcodefakers verifies the shipped barcode faker: every generated
// EAN-13 / UPC-A / ISBN-13 must carry a correct GS1 check digit.
package barcodefakers

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"mockarty/internal/wasmext"
)

func loadModule(t *testing.T) *wasmext.Module {
	t.Helper()
	w, err := os.ReadFile("faker.wasm")
	if err != nil {
		t.Fatalf("read faker.wasm (build it: ./build.sh): %v", err)
	}
	mod, err := wasmext.Compile(context.Background(), "barcode-fakers/faker.wasm", w, wasmext.Limits{})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return mod
}

func call(t *testing.T, mod *wasmext.Module, name string) string {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"name": name})
	out, err := mod.Call(context.Background(), "faker", in)
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	var res struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("%s output %q not {\"value\":...}: %v", name, string(out), err)
	}
	return res.Value
}

// gs1CheckValid recomputes the GS1 mod-10 check digit and compares it to the last
// digit — the same rule the guest used, independently implemented here.
func gs1CheckValid(code string) bool {
	if len(code) < 2 {
		return false
	}
	sum := 0
	payload := code[:len(code)-1]
	for i := 0; i < len(payload); i++ {
		if payload[i] < '0' || payload[i] > '9' {
			return false
		}
		v := int(payload[i] - '0')
		if (len(payload)-1-i)%2 == 0 {
			sum += v * 3
		} else {
			sum += v
		}
	}
	check := (10 - sum%10) % 10
	return int(code[len(code)-1]-'0') == check
}

func TestBarcodeFakers(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())
	cases := []struct {
		name   string
		length int
		prefix string
	}{
		{"ean13", 13, ""},
		{"upc", 12, ""},
		{"isbn13", 13, "978"},
	}
	for _, c := range cases {
		for i := 0; i < 30; i++ {
			v := call(t, mod, c.name)
			if len(v) != c.length {
				t.Fatalf("%s length = %d, want %d (%q)", c.name, len(v), c.length, v)
			}
			if c.prefix != "" && v[:len(c.prefix)] != c.prefix {
				t.Fatalf("%s should start with %s, got %q", c.name, c.prefix, v)
			}
			if !gs1CheckValid(v) {
				t.Fatalf("%s has invalid check digit: %q", c.name, v)
			}
		}
	}
	// Unknown name → empty value (host falls through).
	if v := call(t, mod, "nope"); v != "" {
		t.Fatalf("unknown faker should return empty, got %q", v)
	}
}
