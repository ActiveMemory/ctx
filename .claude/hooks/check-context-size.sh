#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Context size checkpoint hook for Claude Code.
# Counts prompts per session and outputs reminders at adaptive intervals,
# prompting Claude to assess remaining context capacity.
#
# Adaptive frequency:
#   Prompts  1-15: silent
#   Prompts 16-30: every 5th prompt
#   Prompts   30+: every 3rd prompt
#
# Output: Checkpoint messages to stderr (non-blocking, visible to Claude)
# Exit: Always 0 (never blocks execution)

# Read hook input from stdin (JSON)
HOOK_INPUT=$(cat)
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id // "unknown"')

# Use a user-specific temp directory to prevent symlink race attacks (M-3).
CTX_TMPDIR="${XDG_RUNTIME_DIR:-/tmp}/ctx-$(id -u)"
mkdir -p "$CTX_TMPDIR" && chmod 700 "$CTX_TMPDIR"

COUNTER_FILE="${CTX_TMPDIR}/context-check-${SESSION_ID}"
LOG_DIR=".context/logs"
LOG_FILE="${LOG_DIR}/check-context-size.log"

# Ensure log directory exists
mkdir -p "$LOG_DIR"

# Log helper: timestamp + message
log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] [session:${SESSION_ID:0:8}] $*" >> "$LOG_FILE"
}

# Initialize or increment counter
if [ -f "$COUNTER_FILE" ]; then
    COUNT=$(cat "$COUNTER_FILE")
    COUNT=$((COUNT + 1))
else
    COUNT=1
fi

echo "$COUNT" > "$COUNTER_FILE"

# Adaptive frequency: check more often as session grows
SHOULD_CHECK=false
if [ "$COUNT" -gt 30 ]; then
    # Every 3rd prompt after 30
    if [ $((COUNT % 3)) -eq 0 ]; then SHOULD_CHECK=true; fi
elif [ "$COUNT" -gt 15 ]; then
    # Every 5th prompt after 15
    if [ $((COUNT % 5)) -eq 0 ]; then SHOULD_CHECK=true; fi
fi

if [ "$SHOULD_CHECK" = true ]; then
    echo "" >&2
    echo "┌─ Context Checkpoint (prompt #${COUNT}) ────────────────" >&2
    echo "│ Assess remaining context capacity." >&2
    echo "│ If usage exceeds ~80%, inform the user." >&2
    echo "└──────────────────────────────────────────────────" >&2
    echo "" >&2
    log "prompt#${COUNT} CHECKPOINT"
else
    log "prompt#${COUNT} silent"
fi

exit 0
