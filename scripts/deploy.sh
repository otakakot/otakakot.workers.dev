#!/usr/bin/env bash
set -euo pipefail

node scripts/generate-wrangler.js

wrangler deploy --config wrangler.generated.jsonc
