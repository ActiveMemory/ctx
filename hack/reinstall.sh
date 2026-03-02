#!/usr/bin/env bash
# Quick rebuild + reinstall for development.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
make build
cp -f ctx /usr/local/bin/ctx 2>/dev/null || sudo cp -f ctx /usr/local/bin/ctx
echo "ctx $(ctx version 2>/dev/null || echo '(installed)') → /usr/local/bin/ctx"
