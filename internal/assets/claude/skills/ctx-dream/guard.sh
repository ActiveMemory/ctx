#!/usr/bin/env bash
# guard.sh — PreToolUse hook for a headless ctx-dream pass. Bounds the
# dream's blast radius: during a pass the skill only PROPOSES, so writes
# are allowed ONLY under dreams/ (the gitignored notebook), and shell is
# restricted to a small read-only allowlist. This mirrors the Go guards
# in internal/dream (WriteScope + Leak) for the Claude Code executor
# path; other executors call that Go logic directly (see the executor
# contract in docs/reference/dream-executor-contract.md).
#
# Wire it into the dream-specific settings the scheduled run points at
# (NOT the project default settings — the dream is opt-in):
#
#   {
#     "hooks": {
#       "PreToolUse": [
#         { "matcher": "Write|Edit|MultiEdit",
#           "hooks": [{ "type": "command",
#             "command": "<skills>/ctx-dream/guard.sh" }] },
#         { "matcher": "Bash",
#           "hooks": [{ "type": "command",
#             "command": "<skills>/ctx-dream/guard.sh" }] }
#       ]
#     }
#   }
#
# Hook protocol: the tool call JSON arrives on stdin. Exit 0 = allow; a
# JSON block decision = deny.
set -euo pipefail

ALLOW_PREFIX="dreams/"                                       # only writable subtree
ALLOW_CMDS_RE='^(git|grep|rg|sha256sum|cat|ls|find|mkdir)\b' # read-only-ish shell

payload="$(cat)"
tool="$(printf '%s' "$payload" | grep -o '"tool_name"[^,]*' | head -1 | sed 's/.*: *"\(.*\)"/\1/')"

deny() { echo "{\"decision\":\"block\",\"reason\":\"ctx-dream guard: $1\"}"; exit 0; }

case "$tool" in
  Write|Edit|MultiEdit)
    path="$(printf '%s' "$payload" | grep -o '"file_path"[^,]*' | head -1 | sed 's/.*: *"\(.*\)"/\1/')"
    [ -n "$path" ] || deny "could not parse file_path"
    case "$path" in
      "$ALLOW_PREFIX"*|*"/$ALLOW_PREFIX"*) : ;;             # under dreams/ — allow
      *) deny "write outside $ALLOW_PREFIX rejected: $path" ;;
    esac
    ;;
  Bash)
    cmd="$(printf '%s' "$payload" | grep -o '"command"[^}]*' | head -1 | sed 's/.*: *"\(.*\)"$/\1/')"
    printf '%s' "$cmd" | grep -Eq "$ALLOW_CMDS_RE" || deny "shell command not in allowlist: $cmd"
    printf '%s' "$cmd" | grep -Eq 'rm -rf|--hard|push|;|&&|\|\||`|\$\(' && deny "disallowed shell construct"
    ;;
esac
exit 0
