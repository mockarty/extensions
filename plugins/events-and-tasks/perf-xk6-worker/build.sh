#!/usr/bin/env bash
# Builds the worker for THIS machine and writes the manifest that matches it.
#
# The binary is never committed: a managed worker's artifact is identified by
# its sha256, and a checked-in binary would go stale against the source beside
# it the first time anyone edits either. So the manifest is generated here, from
# the bytes that were actually produced.
set -euo pipefail
cd "$(dirname "$0")"

OS="${GOOS:-$(go env GOOS)}"
ARCH="${GOARCH:-$(go env GOARCH)}"

mkdir -p bin
(cd src && GOOS="$OS" GOARCH="$ARCH" go build -trimpath -o ../bin/engine .)

if command -v sha256sum >/dev/null 2>&1; then
    SHA="$(sha256sum bin/engine | cut -d' ' -f1)"
else
    SHA="$(shasum -a 256 bin/engine | cut -d' ' -f1)"
fi

sed -e "s/__OS__/$OS/" -e "s/__ARCH__/$ARCH/" -e "s/__SHA256__/$SHA/" \
    plugin.json.template > plugin.json

echo "built bin/engine for $OS/$ARCH"
echo "sha256 $SHA written to plugin.json"
echo
echo "next:  mockarty-cli plugin test .  &&  mockarty-cli plugin pack ."
