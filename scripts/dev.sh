#!/usr/bin/env bash
set -euo pipefail

rm -rf .wrangler/state/

set -a
source .env
set +a

node scripts/generate-wrangler.mjs

wrangler dev
