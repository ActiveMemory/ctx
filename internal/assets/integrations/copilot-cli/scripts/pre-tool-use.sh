#!/bin/bash
# ctx pre-tool-use hook for Copilot CLI
# Ensures context is loaded and blocks dangerous commands
set -euo pipefail

TOOL="${1:-}"

# Always check context load gate
ctx system context-load-gate 2>/dev/null || true

# Bash-specific hooks
if [ "$TOOL" = "bash" ] || [ "$TOOL" = "powershell" ]; then
  ctx system block-non-path-ctx 2>/dev/null || true
  ctx system qa-reminder 2>/dev/null || true
fi
