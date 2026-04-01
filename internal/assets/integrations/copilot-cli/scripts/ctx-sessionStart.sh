#!/usr/bin/env bash
# ctx sessionStart hook for GitHub Copilot CLI
# Records session start and loads context status.
set -euo pipefail

if command -v ctx >/dev/null 2>&1; then
  ctx system session-event --type start --caller copilot-cli 2>/dev/null || true
fi
