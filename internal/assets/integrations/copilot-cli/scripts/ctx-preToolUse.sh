#!/usr/bin/env bash
# ctx preToolUse hook for GitHub Copilot CLI
# Reads tool invocation JSON from stdin and blocks dangerous commands.
set -euo pipefail

INPUT=$(cat)

# Extract the tool name from the JSON input.
TOOL=""
if command -v jq >/dev/null 2>&1; then
  TOOL=$(echo "$INPUT" | jq -r '.tool_name // .tool // empty' 2>/dev/null)
fi

# Block dangerous shell commands matching known patterns.
if [ "$TOOL" = "shell" ] || [ "$TOOL" = "bash" ]; then
  COMMAND=""
  if command -v jq >/dev/null 2>&1; then
    COMMAND=$(echo "$INPUT" | jq -r '.input.command // empty' 2>/dev/null)
  fi

  case "$COMMAND" in
    *"sudo "* | *"rm -rf /"* | *"rm -rf ~"* | *"chmod 777"*)
      echo '{"decision":"deny","reason":"ctx: blocked dangerous command"}' >&2
      exit 1
      ;;
    *"git push"* | *"git reset --hard"*)
      echo '{"decision":"deny","reason":"ctx: blocked irreversible git operation — review first"}' >&2
      exit 1
      ;;
  esac
fi
