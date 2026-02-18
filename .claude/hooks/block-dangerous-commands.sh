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
# - cp/install/mv to any bin directory (agent must not install binaries)
# - go install           (writes to ~/go/bin, shadows system ctx)
# - ./ctx, ./dist/ctx, go run ./cmd/ctx, absolute-path ctx
#   (must use ctx from PATH per AGENT_PLAYBOOK.md)
#
# NOT BLOCKED (intentionally):
# - rm -rf              (legitimate in cleanup/tests, too restrictive)
# - /tmp/ctx-test*      (integration test harness)

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

# cp/mv to specific bin directories — agent must never install binaries
# Use explicit directory list to avoid false positives from "install" inside quoted args.
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(cp|mv)\s+\S+\s+(/usr/local/bin|/usr/bin|~/go/bin|~/.local/bin|/home/\S+/go/bin|/home/\S+/.local/bin)'; then
  BLOCKED_REASON="Agent must not copy binaries to bin directories. Ask the user to run 'sudo make install' instead."
fi

# go install — writes to ~/go/bin, can shadow the system ctx
# Anchor to command start position (after separator or line start) to avoid matching inside quoted text.
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(^|;|&&|\|\|)\s*go install\b'; then
  BLOCKED_REASON="Agent must not run 'go install'. Ask the user to run 'sudo make install' instead."
fi

# ./ctx or ./dist/ctx — must use ctx from PATH, not relative paths
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(^|;|&&|\|\||\||\$\(|`)\s*\./ctx(\s|$)'; then
  BLOCKED_REASON="Use 'ctx' from PATH, not './ctx'. See AGENT_PLAYBOOK.md: Invoking ctx. Run 'which ctx' to verify it is installed."
fi
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(^|;|&&|\|\||\||\$\(|`)\s*\./dist/ctx(\s|$)'; then
  BLOCKED_REASON="Use 'ctx' from PATH, not './dist/ctx'. See AGENT_PLAYBOOK.md: Invoking ctx."
fi

# go run ./cmd/ctx — use the installed binary, not source
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE 'go run \./cmd/ctx'; then
  BLOCKED_REASON="Use 'ctx' from PATH, not 'go run ./cmd/ctx'. See AGENT_PLAYBOOK.md: Invoking ctx."
fi

# Absolute paths to ctx binary (except /tmp/ctx-test for integration tests)
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(^|;|&&|\|\||\|)\s*(/home/|/tmp/|/var/)\S*/ctx(\s|$)' \
   && ! echo "$COMMAND" | grep -qE '/tmp/ctx-test'; then
  BLOCKED_REASON="Use 'ctx' from PATH, not absolute paths. See AGENT_PLAYBOOK.md: Invoking ctx."
fi

if [ -n "$BLOCKED_REASON" ]; then
  jq -n --arg reason "$BLOCKED_REASON" '{"decision":"block","reason":$reason}'
  exit 0
fi

exit 0
