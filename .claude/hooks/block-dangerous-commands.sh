#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Block commands that Claude cannot or should not run.
#
# BLOCKED:
# - sudo *              (cannot enter password)
# - cp/install to ~/.local/bin  (workaround for PATH ctx rules)
#
# NOT BLOCKED (intentionally):
# - rm -rf              (legitimate in cleanup/tests, too restrictive)

HOOK_INPUT=$(cat)
COMMAND=$(echo "$HOOK_INPUT" | jq -r '.tool_input.command // empty')

if [ -z "$COMMAND" ]; then
  exit 0
fi

BLOCKED_REASON=""

# sudo — Claude cannot enter a password, this will always hang or fail
if echo "$COMMAND" | grep -qE '(^|\s|;|&&|\|\|)sudo\s'; then
  BLOCKED_REASON="Cannot use sudo (no password access). Use 'make build && sudo make install' manually if needed."
fi

# cp/install to ~/.local/bin — known workaround that breaks PATH ctx rules
if echo "$COMMAND" | grep -qE '(cp|install)\s.*~/\.local/bin'; then
  BLOCKED_REASON="Do not copy binaries to ~/.local/bin — this overrides the system ctx in /usr/local/bin. Use 'ctx' from PATH."
fi

if [ -n "$BLOCKED_REASON" ]; then
  cat << EOF
{"decision": "block", "reason": "$BLOCKED_REASON"}
EOF
  exit 0
fi

exit 0
