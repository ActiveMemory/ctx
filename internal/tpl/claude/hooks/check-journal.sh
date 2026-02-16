#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Journal export/enrich reminder hook for Claude Code.
# Detects unexported sessions and unenriched journal entries, then prints
# actionable commands. Runs once per day (throttled by marker file).
#
# Two-stage check:
#   Stage 1 — Unexported sessions: compare newest JSONL mtime across
#             ~/.claude/projects/ vs newest .context/journal/*.md mtime.
#   Stage 2 — Unenriched entries: journal files without YAML frontmatter
#             (don't start with "---").
#
# Output: Reminder messages to stdout (prepended as context for Claude)
# Exit: Always 0 (never blocks execution)

# Use a user-specific temp directory to prevent symlink race attacks (M-3).
CTX_TMPDIR="${XDG_RUNTIME_DIR:-/tmp}/ctx-$(id -u)"
mkdir -p "$CTX_TMPDIR" && chmod 700 "$CTX_TMPDIR"

REMINDED="${CTX_TMPDIR}/journal-reminded"
JOURNAL_DIR=".context/journal"
CLAUDE_PROJECTS_DIR="$HOME/.claude/projects"

# Only remind once per day to avoid spam
if [ -f "$REMINDED" ]; then
  REMINDED_DATE=$(date -r "$REMINDED" +%Y%m%d 2>/dev/null || echo "")
  if [ "$REMINDED_DATE" = "$(date +%Y%m%d)" ]; then
    exit 0
  fi
fi

# Bail out early if there's no journal directory or Claude projects directory
if [ ! -d "$JOURNAL_DIR" ] || [ ! -d "$CLAUDE_PROJECTS_DIR" ]; then
  exit 0
fi

# --- Stage 1: Unexported sessions ---

# Find the newest journal entry mtime (0 if no entries)
NEWEST_JOURNAL=0
for f in "$JOURNAL_DIR"/*.md; do
  [ -f "$f" ] || continue
  mtime=$(stat -c %Y "$f" 2>/dev/null || stat -f %m "$f" 2>/dev/null || echo 0)
  if [ "$mtime" -gt "$NEWEST_JOURNAL" ]; then
    NEWEST_JOURNAL="$mtime"
  fi
done

# Count JSONL files newer than the newest journal entry
UNEXPORTED=0
if [ -d "$CLAUDE_PROJECTS_DIR" ]; then
  while IFS= read -r -d '' jsonl; do
    mtime=$(stat -c %Y "$jsonl" 2>/dev/null || stat -f %m "$jsonl" 2>/dev/null || echo 0)
    if [ "$mtime" -gt "$NEWEST_JOURNAL" ]; then
      UNEXPORTED=$((UNEXPORTED + 1))
    fi
  done < <(find "$CLAUDE_PROJECTS_DIR" -name '*.jsonl' -print0 2>/dev/null)
fi

# --- Stage 2: Unenriched entries ---

UNENRICHED=0
for f in "$JOURNAL_DIR"/*.md; do
  [ -f "$f" ] || continue
  # Check if the file starts with "---" (YAML frontmatter = enriched)
  first_line=$(head -n 1 "$f" 2>/dev/null || echo "")
  if [ "$first_line" != "---" ]; then
    UNENRICHED=$((UNENRICHED + 1))
  fi
done

# --- Output ---

if [ "$UNEXPORTED" -eq 0 ] && [ "$UNENRICHED" -eq 0 ]; then
  exit 0
fi

echo "IMPORTANT: Relay this journal reminder to the user VERBATIM before answering their question."
echo ""
echo "┌─ Journal Reminder ─────────────────────────────"

if [ "$UNEXPORTED" -gt 0 ] && [ "$UNENRICHED" -gt 0 ]; then
  echo "│ You have ${UNEXPORTED} new session(s) not yet exported."
  echo "│ ${UNENRICHED} existing entries need enrichment."
  echo "│"
  echo "│ Export & enrich:"
  echo "│   ctx recall export --all"
  echo "│   /ctx-journal-enrich-all"
elif [ "$UNEXPORTED" -gt 0 ]; then
  echo "│ You have ${UNEXPORTED} new session(s) not yet exported."
  echo "│"
  echo "│ Export:"
  echo "│   ctx recall export --all"
else
  echo "│ ${UNENRICHED} journal entries need enrichment."
  echo "│"
  echo "│ Enrich:"
  echo "│   /ctx-journal-enrich-all"
fi

echo "└────────────────────────────────────────────────"

touch "$REMINDED"

exit 0
