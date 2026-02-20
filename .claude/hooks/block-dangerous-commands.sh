#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Regex safety net for commands that deny rules cannot express.
#
# The bulk of command blocking is handled by permissions.deny in
# settings.local.json, which is automatically populated during
# `ctx init`. Deny rules cover prefix-based patterns like:
#   sudo, git push, rm -rf, curl, wget, go install, ./ctx, etc.
#
# This hook catches only patterns that require regex matching:
# - Mid-command sudo/git-push (after &&, ||, ;)
# - cp/mv to bin directories
# - cp/install to ~/.local/bin
# - Absolute-path ctx (except /tmp/ctx-test)

HOOK_INPUT=$(cat)
COMMAND=$(echo "$HOOK_INPUT" | jq -r '.tool_input.command // empty')

if [ -z "$COMMAND" ]; then
  exit 0
fi

BLOCKED_REASON=""

# Mid-command sudo — after && || ; (prefix sudo caught by deny rule)
if echo "$COMMAND" | grep -qE '(;|&&|\|\|)\s*sudo\s'; then
  BLOCKED_REASON="Cannot use sudo (no password access). Use 'make build && sudo make install' manually if needed."
fi

# Mid-command git push — after && || ; (prefix git push caught by deny rule)
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(;|&&|\|\|)\s*git\s+push'; then
  BLOCKED_REASON="git push requires explicit user approval."
fi

# cp/mv to specific bin directories — agent must never install binaries
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(cp|mv)\s+\S+\s+(/usr/local/bin|/usr/bin|~/go/bin|~/.local/bin|/home/\S+/go/bin|/home/\S+/.local/bin)'; then
  BLOCKED_REASON="Agent must not copy binaries to bin directories. Ask the user to run 'sudo make install' instead."
fi

# cp/install to ~/.local/bin — known workaround that breaks PATH ctx rules
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(cp|install)\s.*~/\.local/bin'; then
  BLOCKED_REASON="Do not copy binaries to ~/.local/bin — this overrides the system ctx in /usr/local/bin. Use 'ctx' from PATH."
fi

# Absolute paths to ctx binary (except /tmp/ctx-test for integration tests)
if [ -z "$BLOCKED_REASON" ] && echo "$COMMAND" | grep -qE '(^|;|&&|\|\||\|)\s*(/home/|/tmp/|/var/)\S*/ctx(\s|$)' \
   && ! echo "$COMMAND" | grep -qE '/tmp/ctx-test'; then
  BLOCKED_REASON="Use 'ctx' from PATH, not absolute paths. See AGENT_PLAYBOOK.md: Invoking ctx."
fi

if [ -n "$BLOCKED_REASON" ]; then
  cat << EOF
{"decision": "block", "reason": "$BLOCKED_REASON"}
EOF
  exit 0
fi

exit 0
