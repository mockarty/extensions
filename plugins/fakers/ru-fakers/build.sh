#!/usr/bin/env bash
# Build the ru-fakers WASM module from its TinyGo guest source.
# Requires TinyGo (https://tinygo.org). The product itself needs no toolchain —
# it runs the committed faker.wasm through its built-in wazero runtime.
set -euo pipefail
cd "$(dirname "$0")/guest"
tinygo build -target=wasm-unknown -no-debug -o ../faker.wasm .
echo "built faker.wasm"
