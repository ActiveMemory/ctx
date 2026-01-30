#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Block git push - requires explicit user approval
#

LOG="/tmp/claude-hook-debug.log"

# Read hook input from stdin (JSON) - same as block-non-path-ctx.sh
HOOK_INPUT=$(cat)
COMMAND=$(echo "$HOOK_INPUT" | jq -r '.tool_input.command // empty')

echo "$(date '+%Y-%m-%d %H:%M:%S') block-git-push.sh triggered" >> "$LOG"
echo "  COMMAND: $COMMAND" >> "$LOG"

if echo "$COMMAND" | grep -qE 'git\s+push'; then
  echo "  BLOCKED: git push detected" >> "$LOG"
  echo '{"decision": "block", "reason": "git push requires explicit user approval"}'
  exit 0
fi

echo "  PASSED: no git push" >> "$LOG"
