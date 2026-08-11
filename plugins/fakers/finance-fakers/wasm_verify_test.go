// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package financefakers

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"mockarty/internal/wasmext"
)

func call(t *testing.T, mod *wasmext.Module, name string) string {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"name": name})
	out, err := mod.Call(context.Background(), "fin_faker", in)
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	var res struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("%s output %q not JSON: %v", name, out, err)
	}
	return res.Value
}

func luhnValid(s string) bool {
	sum, alt := 0, false
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
		d := int(s[i] - '0')
		if alt {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		alt = !alt
	}
	return sum%10 == 0
}

func TestFinanceFakers(t *testing.T) {
	w, err := os.ReadFile("faker.wasm")
	if err != nil {
		t.Skipf("build it: cd guest && tinygo build -target=wasm-unknown -no-debug -o ../faker.wasm .: %v", err)
	}
	mod, err := wasmext.Compile(context.Background(), "finance", w, wasmext.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	defer mod.Close(context.Background())

	t.Run("credit_card_luhn", func(t *testing.T) {
		for i := 0; i < 500; i++ { // well past the ~256 zero-alloc threshold
			v := call(t, mod, "credit_card")
			if len(v) != 16 || !luhnValid(v) {
				t.Fatalf("credit_card %q is not a 16-digit Luhn-valid number (call %d)", v, i)
			}
		}
	})
	t.Run("iban_shape", func(t *testing.T) {
		v := call(t, mod, "iban")
		if len(v) != 22 || v[:2] != "DE" {
			t.Fatalf("iban = %q, want DE + 20 digits", v)
		}
	})
	t.Run("unknown_empty", func(t *testing.T) {
		if v := call(t, mod, "nope"); v != "" {
			t.Fatalf("unknown faker returned %q", v)
		}
	})
}
