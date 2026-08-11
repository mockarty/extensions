# ru-fakers — Russian identifier fakers (WASM plugin example)

Adds four faker functions to Mockarty's dynamic-response templater:

| Faker | Produces |
|-------|----------|
| `$.fake.ru_inn` | 10-digit legal-entity **ИНН** with a valid check digit |
| `$.fake.ru_snils` | 11-digit **СНИЛС** with a valid control number |
| `$.fake.ru_ogrn` | 13-digit **ОГРН** with a valid check digit |
| `$.fake.ru_kpp` | 9-character **КПП** |

Once installed and enabled, use them in any mock body just like a built-in
faker:

```json
{ "inn": "$.fake.ru_inn", "snils": "$.fake.ru_snils" }
```

## How it works

The logic lives in a **sandboxed WebAssembly module** (`faker.wasm`). Mockarty
runs it through its built-in wazero runtime with **no access to disk, network,
env or clock** — it can only read and write its own linear memory. That is the
whole security story for code plugins: "the plugin cannot touch anything without
an administrator's grant."

`guest/main.go` is the TinyGo source; it documents the stable guest ABI
(`memory` + `mk_alloc` + a `fn(ptr,len)->i64` packing `outPtr<<32|outLen`).

## Build & install

```bash
./build.sh                                   # needs TinyGo; rebuilds faker.wasm
mockarty-cli plugin pack .                    # → mockarty.ru-fakers-1.0.0.zip
mockarty-cli plugin install mockarty.ru-fakers-1.0.0.zip
mockarty-cli plugin enable mockarty.ru-fakers
```

The committed `faker.wasm` is ready to install as-is — you only need TinyGo to
change the generators.

## Note on randomness

A freestanding WASM module has no entropy source, so the generators use a small
PRNG advanced per call. Values vary call-to-call (never cached) but the sequence
is deterministic per server start — ideal for reproducible test data.
