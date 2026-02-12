#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Persistence nudge hook for Claude Code.
# Counts prompts since the last .context/ file modification and nudges
# the agent to persist learnings, decisions, or task updates.
#
# Nudge frequency:
#   Prompts  1-10:  silent (too early, let the agent work)
#   Prompts 11-25:  nudge once at prompt 20
#   Prompts   25+:  every 15th prompt
#
# The nudge only fires if no .context/ file has been modified since the
# last nudge (or session start). If the agent is already persisting, the
# hook stays silent.
#
# Output: Nudge messages to stderr (non-blocking, visible to Claude)
# Exit: Always 0 (never blocks execution)

# Read hook input from stdin (JSON)
HOOK_INPUT=$(cat)
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id // "unknown"')

# Use a user-specific temp directory to prevent symlink race attacks (M-3).
CTX_TMPDIR="${XDG_RUNTIME_DIR:-/tmp}/ctx-$(id -u)"
mkdir -p "$CTX_TMPDIR" && chmod 700 "$CTX_TMPDIR"

STATE_FILE="${CTX_TMPDIR}/persistence-nudge-${SESSION_ID}"
CONTEXT_DIR=".context"
LOG_DIR="${CONTEXT_DIR}/logs"
LOG_FILE="${LOG_DIR}/check-persistence.log"

# Ensure log directory exists
mkdir -p "$LOG_DIR"

# Log helper: timestamp + message
log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] [session:${SESSION_ID:0:8}] $*" >> "$LOG_FILE"
}

# Get the most recent modification time of any .context/*.md file
get_latest_context_mtime() {
  local latest=0
  for f in "$CONTEXT_DIR"/*.md; do
    [ -f "$f" ] || continue
    mtime=$(stat -c %Y "$f" 2>/dev/null || stat -f %m "$f" 2>/dev/null || echo 0)
    if [ "$mtime" -gt "$latest" ]; then
      latest="$mtime"
    fi
  done
  echo "$latest"
}

# Initialize state file if it doesn't exist
if [ ! -f "$STATE_FILE" ]; then
  INITIAL_MTIME=$(get_latest_context_mtime)
  echo "count=1" > "$STATE_FILE"
  echo "last_nudge=0" >> "$STATE_FILE"
  echo "last_mtime=${INITIAL_MTIME}" >> "$STATE_FILE"
  log "init count=1 mtime=${INITIAL_MTIME}"
  exit 0
fi

# Read state
COUNT=$(grep "^count=" "$STATE_FILE" | cut -d= -f2)
LAST_NUDGE=$(grep "^last_nudge=" "$STATE_FILE" | cut -d= -f2)
LAST_MTIME=$(grep "^last_mtime=" "$STATE_FILE" | cut -d= -f2)
COUNT=$((COUNT + 1))

# Check current mtime
CURRENT_MTIME=$(get_latest_context_mtime)

# If context files were modified since last check, reset the nudge counter
if [ "$CURRENT_MTIME" -gt "$LAST_MTIME" ]; then
  # Agent is persisting — reset and stay silent
  echo "count=${COUNT}" > "$STATE_FILE"
  echo "last_nudge=${COUNT}" >> "$STATE_FILE"
  echo "last_mtime=${CURRENT_MTIME}" >> "$STATE_FILE"
  log "prompt#${COUNT} context-modified, reset nudge counter"
  exit 0
fi

# Calculate prompts since last nudge (or session start)
SINCE_NUDGE=$((COUNT - LAST_NUDGE))

# Determine if we should nudge
SHOULD_NUDGE=false
if [ "$COUNT" -ge 11 ] && [ "$COUNT" -le 25 ] && [ "$SINCE_NUDGE" -ge 20 ]; then
  SHOULD_NUDGE=true
elif [ "$COUNT" -gt 25 ] && [ "$SINCE_NUDGE" -ge 15 ]; then
  SHOULD_NUDGE=true
fi

if [ "$SHOULD_NUDGE" = true ]; then
  echo "" >&2
  echo "┌─ Persistence Checkpoint (prompt #${COUNT}) ───────────" >&2
  echo "│ No context files updated in ${SINCE_NUDGE}+ prompts." >&2
  echo "│ Have you discovered learnings, made decisions," >&2
  echo "│ or completed tasks worth persisting?" >&2
  echo "└──────────────────────────────────────────────────" >&2
  echo "" >&2

  log "prompt#${COUNT} NUDGE since_nudge=${SINCE_NUDGE}"

  # Update last nudge
  echo "count=${COUNT}" > "$STATE_FILE"
  echo "last_nudge=${COUNT}" >> "$STATE_FILE"
  echo "last_mtime=${LAST_MTIME}" >> "$STATE_FILE"
else
  log "prompt#${COUNT} silent since_nudge=${SINCE_NUDGE}"
  # Just update count
  echo "count=${COUNT}" > "$STATE_FILE"
  echo "last_nudge=${LAST_NUDGE}" >> "$STATE_FILE"
  echo "last_mtime=${LAST_MTIME}" >> "$STATE_FILE"
fi

exit 0
