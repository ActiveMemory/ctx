#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Quick context usage monitor for Claude Code sessions.
# Usage: ./hack/context-watch.sh [interval_seconds]
#
# Finds the most recently modified session JSONL and estimates token usage.

INTERVAL="${1:-10}"
PROJECT_DIR="$HOME/.claude/projects"
MODEL_LIMIT=200000
AUTOCOMPACT_BUFFER=33000       # reserved by Claude Code, not usable
SYSTEM_OVERHEAD=20000          # system prompt + tools + skills + memory
EFFECTIVE_LIMIT=$((MODEL_LIMIT - AUTOCOMPACT_BUFFER))

while true; do
  clear

  # Find most recently modified JSONL
  JSONL=$(find "$PROJECT_DIR" -name '*.jsonl' -type f -printf '%T@ %p\n' 2>/dev/null \
    | sort -rn | head -1 | cut -d' ' -f2-)

  if [ -z "$JSONL" ]; then
    echo "No active session found."
    sleep "$INTERVAL"
    continue
  fi

  # File stats
  SIZE_BYTES=$(stat -c%s "$JSONL" 2>/dev/null)
  SIZE_KB=$((SIZE_BYTES / 1024))

  # Estimate tokens: JSONL ~30 chars per token (JSON keys, escaping, metadata)
  # Plus system overhead (prompt, tools, skills) not in the JSONL
  CHARS=$(wc -c < "$JSONL")
  MSG_TOKENS=$((CHARS / 30))
  EST_TOKENS=$((MSG_TOKENS + SYSTEM_OVERHEAD))
  EST_TOKENS_K=$((EST_TOKENS / 1000))
  EFFECTIVE_K=$((EFFECTIVE_LIMIT / 1000))
  REMAINING=$((EFFECTIVE_LIMIT - EST_TOKENS))
  REMAINING_K=$((REMAINING / 1000))
  if [ "$REMAINING_K" -lt 0 ]; then REMAINING_K=0; fi
  PCT=$((EST_TOKENS * 100 / EFFECTIVE_LIMIT))

  # Color based on usage
  if [ "$PCT" -lt 50 ]; then
    COLOR="\033[32m" # green
    STATUS="HEALTHY"
  elif [ "$PCT" -lt 75 ]; then
    COLOR="\033[33m" # yellow
    STATUS="MONITOR"
  else
    COLOR="\033[31m" # red
    STATUS="SAVE YOUR CHANGES AND END THE SESSION"
  fi
  RESET="\033[0m"

  # Progress bar
  BAR_WIDTH=40
  FILLED=$((PCT * BAR_WIDTH / 100))
  if [ "$FILLED" -gt "$BAR_WIDTH" ]; then FILLED=$BAR_WIDTH; fi
  EMPTY=$((BAR_WIDTH - FILLED))
  BAR=$(printf '%0.s█' $(seq 1 "$FILLED" 2>/dev/null))
  BAR="${BAR}$(printf '%0.s░' $(seq 1 "$EMPTY" 2>/dev/null))"

  # Message count
  LINES=$(wc -l < "$JSONL")

  # Session name from path
  SESSION=$(basename "$JSONL" .jsonl)

  echo -e "  Context Monitor  ${COLOR}[$STATUS]${RESET}"
  echo ""
  echo -e "  ${COLOR}${BAR}${RESET}  ~${EST_TOKENS_K}k / ${EFFECTIVE_K}k tokens (~${PCT}%)"
  echo -e "  Remaining: ~${REMAINING_K}k usable tokens"
  echo ""
  echo "  Session:  ${SESSION:0:40}"
  echo "  File:     ${SIZE_KB} KB, ${LINES} lines"
  echo "  Updated:  $(date -r "$JSONL" '+%H:%M:%S' 2>/dev/null || stat -c '%y' "$JSONL" 2>/dev/null | cut -d. -f1)"
  echo ""
  echo "  Refreshing every ${INTERVAL}s. Ctrl+C to stop."

  sleep "$INTERVAL"
done
