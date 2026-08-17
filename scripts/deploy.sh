#!/usr/bin/env bash
set -euo pipefail

node scripts/generate-wrangler.mjs

wrangler deploy
