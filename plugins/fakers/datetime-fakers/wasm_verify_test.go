// Copyright (c) 2026 Mockarty. All rights reserved.
// Licensed under the Mockarty Software License Agreement.
// See LICENSE file in the project root for full license text.

package datetimefakers

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"testing"
	"time"

	"mockarty/internal/wasmext"
)

func load(t *testing.T) *wasmext.Module {
	t.Helper()
	w, err := os.ReadFile("faker.wasm")
	if err != nil {
		t.Fatalf("read faker.wasm (./build.sh): %v", err)
	}
	m, err := wasmext.Compile(context.Background(), "datetime-fakers", w, wasmext.Limits{})
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

func TestDatetimeFakers(t *testing.T) {
	m := load(t)
	defer m.Close(context.Background())
	cron := regexp.MustCompile(`^\d+ \d+ \d+ \d+ \d+$`)
	dur := regexp.MustCompile(`^PT\d+H\d+M$`)
	for i := 0; i < 40; i++ {
		if _, err := time.Parse("2006-01-02", call(t, m, "iso_date")); err != nil {
			t.Fatalf("iso_date not parseable: %v", err)
		}
		if _, err := time.Parse(time.RFC3339, call(t, m, "iso_datetime")); err != nil {
			t.Fatalf("iso_datetime not RFC3339: %v", err)
		}
		if c := call(t, m, "cron"); !cron.MatchString(c) {
			t.Fatalf("cron wrong shape: %q", c)
		}
		if d := call(t, m, "duration"); !dur.MatchString(d) {
			t.Fatalf("duration wrong shape: %q", d)
		}
	}
	if call(t, m, "nope") != "" {
		t.Fatal("unknown faker should be empty")
	}
}
