#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Post-commit nudge hook for Claude Code.
# Fires after a successful git commit (PostToolUse on Bash).
# Detects git commit commands and nudges the agent to:
#   1. Offer context capture (decision, learning, or neither)
#   2. Ask the user if they want lints and tests run before pushing
#
# The "before YOU push" framing reinforces that the agent should not
# push — the user does that manually (block-git-push.sh is the hard gate).
#
# Output: Agent directive to stdout (prepended as context for Claude)
# Exit: Always 0 (never blocks execution)

# Read hook input from stdin (JSON)
HOOK_INPUT=$(cat)
COMMAND=$(echo "$HOOK_INPUT" | jq -r '.tool_input.command // empty')

# Only trigger on git commit commands
if ! echo "$COMMAND" | grep -qE 'git\s+commit'; then
  exit 0
fi

# Skip amend commits — those are fixups, not milestones
if echo "$COMMAND" | grep -qE -- '--amend'; then
  exit 0
fi

echo ""
echo "┌─ Post-Commit ──────────────────────────────────────────"
echo "│ Commit succeeded."
echo "│"
echo "│ 1. Offer context capture to the user:"
echo "│    Decision (design choice?), Learning (gotcha?), or Neither."
echo "│"
echo "│ 2. Ask the user:"
echo "│    \"Want me to run lints and tests before you push?\""
echo "│"
echo "│ Do NOT push. The user pushes manually."
echo "└────────────────────────────────────────────────────────"
echo ""

exit 0
