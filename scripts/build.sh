#!/usr/bin/env bash
set -euo pipefail

go run github.com/syumai/workers/cmd/workers-assets-gen -mode=go

GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s" -o ./build/app.wasm .
