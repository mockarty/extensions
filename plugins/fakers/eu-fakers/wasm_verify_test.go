// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package eufakers

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"mockarty/internal/wasmext"
)

func load(t *testing.T) *wasmext.Module {
	t.Helper()
	w, err := os.ReadFile("faker.wasm")
	if err != nil {
		t.Fatalf("read faker.wasm (./build.sh): %v", err)
	}
	m, err := wasmext.Compile(context.Background(), "eu-fakers", w, wasmext.Limits{})
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

// ibanMod97 independently validates an IBAN: move the first 4 chars to the end,
// map letters A..Z→10..35, and check the whole number ≡ 1 (mod 97).
func ibanMod97(iban string) bool {
	if len(iban) < 5 {
		return false
	}
	rearr := iban[4:] + iban[:4]
	rem := 0
	for i := 0; i < len(rearr); i++ {
		c := rearr[i]
		switch {
		case c >= '0' && c <= '9':
			rem = (rem*10 + int(c-'0')) % 97
		case c >= 'A' && c <= 'Z':
			v := int(c-'A') + 10
			rem = (rem*100 + v) % 97
		default:
			return false
		}
	}
	return rem == 1
}

func TestEUFakers(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	for i := 0; i < 40; i++ {
		iban := call(t, m, "iban")
		if len(iban) != 22 || iban[:2] != "DE" {
			t.Fatalf("iban shape wrong: %q", iban)
		}
		if !ibanMod97(iban) {
			t.Fatalf("iban fails mod-97: %q", iban)
		}
		bic := call(t, m, "bic")
		if len(bic) != 8 || bic[4:6] != "DE" {
			t.Fatalf("bic wrong: %q", bic)
		}
		for j := 0; j < 4; j++ {
			if bic[j] < 'A' || bic[j] > 'Z' {
				t.Fatalf("bic bank code not letters: %q", bic)
			}
		}
		vat := call(t, m, "vat_de")
		if len(vat) != 11 || vat[:2] != "DE" {
			t.Fatalf("vat_de wrong: %q", vat)
		}
	}
	if call(t, m, "nope") != "" {
		t.Fatal("unknown faker should be empty")
	}
}
