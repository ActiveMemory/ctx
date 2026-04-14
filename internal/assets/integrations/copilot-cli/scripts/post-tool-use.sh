#!/bin/bash
# ctx post-tool-use hook for Copilot CLI
# Checks for post-commit context and task completion
set -euo pipefail

TOOL="${1:-}"

if [ "$TOOL" = "bash" ] || [ "$TOOL" = "powershell" ]; then
  ctx system post-commit 2>/dev/null || true
fi

if [ "$TOOL" = "edit" ] || [ "$TOOL" = "write" ]; then
  ctx system check-task-completion 2>/dev/null || true
fi
