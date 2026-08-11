// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package usfakers

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
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
	m, err := wasmext.Compile(context.Background(), "us-fakers", w, wasmext.Limits{})
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

// abaValid recomputes the ABA routing checksum independently.
func abaValid(s string) bool {
	if len(s) != 9 {
		return false
	}
	w := []int{3, 7, 1, 3, 7, 1, 3, 7, 1}
	sum := 0
	for i := 0; i < 9; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
		sum += int(s[i]-'0') * w[i]
	}
	return sum%10 == 0
}

func TestUSFakers(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	for i := 0; i < 40; i++ {
		if r := call(t, m, "us_routing"); !abaValid(r) {
			t.Fatalf("us_routing invalid ABA checksum: %q", r)
		}
		ssn := call(t, m, "ssn")
		parts := strings.Split(ssn, "-")
		if len(parts) != 3 || len(parts[0]) != 3 || len(parts[1]) != 2 || len(parts[2]) != 4 {
			t.Fatalf("ssn wrong shape: %q", ssn)
		}
		if area, _ := strconv.Atoi(parts[0]); area == 0 || area == 666 || area > 899 {
			t.Fatalf("ssn bad area: %q", ssn)
		}
		ein := call(t, m, "ein")
		ep := strings.Split(ein, "-")
		if len(ep) != 2 || len(ep[0]) != 2 || len(ep[1]) != 7 {
			t.Fatalf("ein wrong shape: %q", ein)
		}
	}
	if call(t, m, "nope") != "" {
		t.Fatal("unknown faker should be empty")
	}
}
