// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package formatmatchers

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"mockarty/internal/wasmext"
)

func load(t *testing.T) *wasmext.Module {
	t.Helper()
	w, err := os.ReadFile("match.wasm")
	if err != nil {
		t.Fatalf("read match.wasm (./build.sh): %v", err)
	}
	m, err := wasmext.Compile(context.Background(), "format-matchers", w, wasmext.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	return m
}
func match(t *testing.T, m *wasmext.Module, key, actual string) bool {
	in, _ := json.Marshal(map[string]string{"key": key, "actual": actual, "expected": ""})
	out, err := m.Call(context.Background(), "match", in)
	if err != nil {
		t.Fatalf("match %s: %v", key, err)
	}
	var r struct {
		Match bool `json:"match"`
	}
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("%s: %v", key, err)
	}
	return r.Match
}

func TestFormatMatchers(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	cases := []struct {
		key, val string
		want     bool
	}{
		{"jwt", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.abc-_123", true},
		{"jwt", "notajwt", false},
		{"jwt", "a.b", false},
		{"jwt", "a..c", false},
		{"e164", "+14155552671", true},
		{"e164", "+1", false},
		{"e164", "14155552671", false},
		{"e164", "+1415555267a", false},
		{"ipv4", "192.168.1.1", true},
		{"ipv4", "255.255.255.255", true},
		{"ipv4", "256.1.1.1", false},
		{"ipv4", "1.2.3", false},
		{"ipv4", "1.2.3.4.5", false},
		{"ipv4", "a.b.c.d", false},
		{"unknown", "whatever", false},
	}
	for _, c := range cases {
		if got := match(t, m, c.key, c.val); got != c.want {
			t.Errorf("%s(%q) = %v, want %v", c.key, c.val, got, c.want)
		}
	}
}
