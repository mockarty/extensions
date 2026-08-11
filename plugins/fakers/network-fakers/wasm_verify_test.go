// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package networkfakers

import (
	"context"
	"encoding/json"
	"net"
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
	m, err := wasmext.Compile(context.Background(), "network-fakers", w, wasmext.Limits{})
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

func TestNetworkFakers(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	for i := 0; i < 40; i++ {
		// cidr must parse
		if _, _, err := net.ParseCIDR(call(t, m, "cidr")); err != nil {
			t.Fatalf("cidr invalid: %v", err)
		}
		// private_ip must be a valid RFC1918 IPv4
		ip := net.ParseIP(call(t, m, "private_ip"))
		if ip == nil || ip.To4() == nil || !ip.IsPrivate() {
			t.Fatalf("private_ip not RFC1918: %v", ip)
		}
		// port in [1024,65535]
		p, err := strconv.Atoi(call(t, m, "port"))
		if err != nil || p < 1024 || p > 65535 {
			t.Fatalf("port out of range: %v", p)
		}
		// mac parses and is locally-administered
		macStr := call(t, m, "mac_oui")
		mac, err := net.ParseMAC(macStr)
		if err != nil || len(mac) != 6 {
			t.Fatalf("mac invalid: %q", macStr)
		}
		if mac[0]&0x02 == 0 {
			t.Fatalf("mac not locally-administered: %q", macStr)
		}
		if !strings.Contains(macStr, ":") {
			t.Fatalf("mac wrong format: %q", macStr)
		}
	}
	if call(t, m, "nope") != "" {
		t.Fatal("unknown faker should be empty")
	}
}
