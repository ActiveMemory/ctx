#!/bin/bash

# Top-level commands from cmd structure
echo "=== TOP-LEVEL COMMANDS (from cmd structure) ==="
defined_cmds=(
  "add"
  "agent"
  "archive"
  "brainstorm"
  "check-links"
  "commit"
  "compact"
  "config"
  "consolidate"
  "decision"
  "deps"
  "doctor"
  "drift"
  "explain"
  "import-plans"
  "init"
  "journal"
  "learning"
  "loop"
  "memory"
  "notify"
  "pad"
  "pause"
  "permission"
  "recall"
  "reflect"
  "remind"
  "resume"
  "site"
  "status"
  "system"
  "task"
  "watch"
  "why"
)

referenced=$(/tmp/count_cmds.sh | awk -F'|' '{print $1}' | sort -u)

echo "Commands with zero references in skills/hooks/recipes:"
for cmd in "${defined_cmds[@]}"; do
  if ! echo "$referenced" | grep -q "^${cmd}$"; then
    echo "  - $cmd"
  fi
done

echo ""
echo "=== COMMAND REFERENCE SUMMARY ==="
cat /tmp/final_counts.txt
