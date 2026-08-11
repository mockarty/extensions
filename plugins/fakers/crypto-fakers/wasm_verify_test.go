// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package cryptofakers

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
	m, err := wasmext.Compile(context.Background(), "crypto-fakers", w, wasmext.Limits{})
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
func isHex(s string) bool {
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return len(s) > 0
}

func TestCryptoFakers(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	for i := 0; i < 20; i++ {
		eth := call(t, m, "eth_address")
		if len(eth) != 42 || eth[:2] != "0x" || !isHex(eth[2:]) {
			t.Fatalf("eth_address invalid: %q", eth)
		}
		tx := call(t, m, "tx_hash")
		if len(tx) != 66 || tx[:2] != "0x" || !isHex(tx[2:]) {
			t.Fatalf("tx_hash invalid: %q", tx)
		}
		btc := call(t, m, "btc_address")
		if len(btc) != 34 || btc[0] != '1' {
			t.Fatalf("btc_address invalid: %q", btc)
		}
		for _, c := range btc {
			if strings.ContainsRune("0OIl", c) {
				t.Fatalf("btc_address has ambiguous base58 char: %q", btc)
			}
		}
	}
	if call(t, m, "nope") != "" {
		t.Fatal("unknown faker should be empty")
	}
}
