#!/usr/bin/env bash
# ctx sessionEnd hook for GitHub Copilot CLI
# Records session end event for recall and context persistence.
set -euo pipefail

if command -v ctx >/dev/null 2>&1; then
  ctx system session-event --type end --caller copilot-cli 2>/dev/null || true
fi
