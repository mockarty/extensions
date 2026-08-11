// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

// Package conditionmatcher holds the committed verification for the
// condition-matcher example plugin: it loads the built match.wasm through the
// real host WASM layer and asserts each matcher decides correctly. This is the
// dogfood gate — the shipped .wasm must always pass the same sandbox the product
// runs plugins in.
package conditionmatcher

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"mockarty/internal/wasmext"
)

func loadModule(t *testing.T) *wasmext.Module {
	t.Helper()
	wasm, err := os.ReadFile("match.wasm")
	if err != nil {
		t.Fatalf("read match.wasm (build it: ./build.sh): %v", err)
	}
	mod, err := wasmext.Compile(context.Background(), "condition-matcher/match.wasm", wasm, wasmext.Limits{})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return mod
}

func match(t *testing.T, mod *wasmext.Module, key, actual, expected string) bool {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"key": key, "actual": actual, "expected": expected})
	out, err := mod.Call(context.Background(), "match", in)
	if err != nil {
		t.Fatalf("call %s(%q,%q): %v", key, actual, expected, err)
	}
	var res struct {
		Match bool `json:"match"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("%s output %q is not {\"match\":...}: %v", key, string(out), err)
	}
	return res.Match
}

func TestLuhn(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())
	cases := []struct {
		num  string
		want bool
	}{
		{"4111111111111111", true},  // valid Visa test PAN
		{"79927398713", true},       // classic Luhn example
		{"1234567812345670", true},  // valid
		{"4111111111111112", false}, // last digit off by one
		{"", false},                 // empty
		{"notanumber", false},       // non-digit
		{"49927398717", false},      // invalid checksum
	}
	for _, c := range cases {
		if got := match(t, mod, "luhn", c.num, ""); got != c.want {
			t.Errorf("luhn(%q) = %v, want %v", c.num, got, c.want)
		}
	}
}

func TestDivisibleBy(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())
	cases := []struct {
		actual, operand string
		want            bool
	}{
		{"100", "5", true},
		{"7", "2", false},
		{"-9", "3", true},
		{"10", "0", false}, // division by zero → never
		{"x", "2", false},  // non-numeric
		{"10", "", false},  // missing operand
	}
	for _, c := range cases {
		if got := match(t, mod, "divisible_by", c.actual, c.operand); got != c.want {
			t.Errorf("divisible_by(%q,%q) = %v, want %v", c.actual, c.operand, got, c.want)
		}
	}
}

// TestUnknownKey asserts an unrecognised matcher key returns match:false rather
// than trapping — a mock naming a key the plugin doesn't serve simply won't match.
func TestUnknownKey(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())
	if match(t, mod, "nope", "whatever", "") {
		t.Error("unknown key should not match")
	}
}
