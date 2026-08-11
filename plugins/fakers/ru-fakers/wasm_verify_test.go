// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

// Package rufakers holds the committed verification for the ru-fakers example
// plugin: it loads the built faker.wasm through the real host WASM layer and
// asserts every faker produces a well-formed Russian identifier with a valid
// check digit. This is the dogfood gate — the shipped .wasm must always pass
// the same sandbox the product runs plugins in.
package rufakers

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"mockarty/internal/wasmext"
)

func loadModule(t *testing.T) *wasmext.Module {
	t.Helper()
	wasm, err := os.ReadFile("faker.wasm")
	if err != nil {
		t.Fatalf("read faker.wasm (build it: cd guest && tinygo build -target=wasm-unknown -no-debug -o ../faker.wasm .): %v", err)
	}
	mod, err := wasmext.Compile(context.Background(), "ru-fakers/faker.wasm", wasm, wasmext.Limits{})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return mod
}

func call(t *testing.T, mod *wasmext.Module, name string) string {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"name": name})
	out, err := mod.Call(context.Background(), "ru_faker", in)
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	var res struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("%s output %q is not {\"value\":...}: %v", name, string(out), err)
	}
	return res.Value
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

func TestRuFakers(t *testing.T) {
	mod := loadModule(t)
	defer mod.Close(context.Background())

	t.Run("inn10_valid_checksum", func(t *testing.T) {
		v := call(t, mod, "ru_inn")
		if len(v) != 10 || !allDigits(v) {
			t.Fatalf("ru_inn = %q, want 10 digits", v)
		}
		weights := []int{2, 4, 10, 3, 5, 9, 4, 6, 8}
		sum := 0
		for i, w := range weights {
			sum += int(v[i]-'0') * w
		}
		if (sum%11)%10 != int(v[9]-'0') {
			t.Fatalf("ru_inn %q has an invalid check digit", v)
		}
	})

	t.Run("snils_valid_control", func(t *testing.T) {
		v := call(t, mod, "ru_snils")
		if len(v) != 11 || !allDigits(v) {
			t.Fatalf("ru_snils = %q, want 11 digits", v)
		}
		sum := 0
		for i := 0; i < 9; i++ {
			sum += int(v[i]-'0') * (9 - i)
		}
		ctrl := sum % 101
		if ctrl == 100 || ctrl == 101 {
			ctrl = 0
		}
		want := ctrl/10*10 + ctrl%10
		got := int(v[9]-'0')*10 + int(v[10]-'0')
		if want != got {
			t.Fatalf("ru_snils %q control %02d, want %02d", v, got, want)
		}
	})

	t.Run("ogrn_valid_checksum", func(t *testing.T) {
		v := call(t, mod, "ru_ogrn")
		if len(v) != 13 || !allDigits(v) {
			t.Fatalf("ru_ogrn = %q, want 13 digits", v)
		}
		mod12 := 0
		for i := 0; i < 12; i++ {
			mod12 = (mod12*10 + int(v[i]-'0')) % 11
		}
		if mod12%10 != int(v[12]-'0') {
			t.Fatalf("ru_ogrn %q has an invalid check digit", v)
		}
	})

	t.Run("kpp_shape", func(t *testing.T) {
		v := call(t, mod, "ru_kpp")
		if len(v) != 9 || !allDigits(v) || v[4:6] != "01" {
			t.Fatalf("ru_kpp = %q, want 9 digits with reason 01", v)
		}
	})

	t.Run("unknown_name_empty", func(t *testing.T) {
		if v := call(t, mod, "not_a_faker"); v != "" {
			t.Fatalf("unknown faker returned %q, want empty (host falls through)", v)
		}
	})

	t.Run("sustained_no_corruption", func(t *testing.T) {
		// Guards the zero-alloc requirement: a heap-allocating guest corrupts its
		// output once TinyGo's leaking GC grows the heap over the buffers (~256
		// calls). Well past that threshold every call must still parse and match
		// the expected shape.
		for i := 0; i < 3000; i++ {
			if v := call(t, mod, "ru_inn"); len(v) != 10 || !allDigits(v) {
				t.Fatalf("ru_inn corrupted at call %d: %q", i, v)
			}
		}
	})

	t.Run("fresh_each_call", func(t *testing.T) {
		// Fakers must not cache: repeated calls should vary. Extremely unlikely
		// to collide across 8 draws if the PRNG advances.
		seen := map[string]bool{}
		for i := 0; i < 8; i++ {
			seen[call(t, mod, "ru_inn")] = true
		}
		if len(seen) < 2 {
			t.Fatalf("ru_inn produced no variation across calls: %v", seen)
		}
	})
}
