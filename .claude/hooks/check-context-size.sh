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

COUNTER_FILE="/tmp/ctx-context-check-${SESSION_ID}"

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
fi

exit 0
