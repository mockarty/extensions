// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package identifierfakers

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"mockarty/internal/wasmext"
)

func load(t *testing.T) *wasmext.Module {
	t.Helper()
	w, err := os.ReadFile("faker.wasm")
	if err != nil {
		t.Fatalf("read faker.wasm (./build.sh): %v", err)
	}
	m, err := wasmext.Compile(context.Background(), "identifier-fakers", w, wasmext.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	return m
}
func call(t *testing.T, m *wasmext.Module, name string) string {
	in, _ := json.Marshal(map[string]string{"name": name})
	out, err := m.Call(context.Background(), "faker", in)
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	var r struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return r.Value
}
func allIn(s, alphabet string) bool {
	for _, c := range s {
		if !strings.ContainsRune(alphabet, c) {
			return false
		}
	}
	return len(s) > 0
}

func TestIdentifierFakers(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	cases := []struct {
		name     string
		length   int
		alphabet string
	}{
		{"ulid", 26, "0123456789ABCDEFGHJKMNPQRSTVWXYZ"},
		{"nanoid", 21, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"},
		{"hex_token", 32, "0123456789abcdef"},
	}
	for _, c := range cases {
		for i := 0; i < 20; i++ {
			v := call(t, m, c.name)
			if len(v) != c.length {
				t.Fatalf("%s len=%d want %d (%q)", c.name, len(v), c.length, v)
			}
			if !allIn(v, c.alphabet) {
				t.Fatalf("%s %q has chars outside alphabet", c.name, v)
			}
		}
	}
	if call(t, m, "nope") != "" {
		t.Fatal("unknown faker should be empty")
	}
}
