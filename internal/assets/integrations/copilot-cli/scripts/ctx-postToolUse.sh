#!/usr/bin/env bash
# ctx postToolUse hook for GitHub Copilot CLI
# Reads tool result JSON from stdin and appends to audit log.
set -euo pipefail

# Append tool invocation to audit log if ctx is available.
if command -v ctx >/dev/null 2>&1; then
  INPUT=$(cat)
  LOGDIR=".context/state"
  LOGFILE="$LOGDIR/copilot-cli-audit.jsonl"

  if [ -d ".context" ]; then
    mkdir -p "$LOGDIR"
    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date +"%Y-%m-%dT%H:%M:%S")
    echo "{\"timestamp\":\"$TIMESTAMP\",\"event\":\"postToolUse\",\"data\":$INPUT}" >> "$LOGFILE" 2>/dev/null || true
  fi
fi
