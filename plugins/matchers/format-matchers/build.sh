#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/guest"
tinygo build -target=wasm-unknown -no-debug -o ../match.wasm .
echo "built match.wasm"
