#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# detect-ai-typography.sh — find docs likely AI-edited but not human-reviewed.
#
# Scans markdown files for em-dashes, smart quotes, and "--" (double hyphen
# used as dash). These are heuristic signals: almost all LLM output uses
# them because the training corpus overwhelmingly contains typographic
# punctuation from published English text.
#
# False positives are possible (em-dash is valid typography). False negatives
# are unlikely given current model behavior.
#
# Usage:
#   ./hack/detect-ai-typography.sh [dir]   # default: docs/
#   ./hack/detect-ai-typography.sh --stat  # summary only (no line detail)

set -euo pipefail

STAT_ONLY=false
DIR=""

for arg in "$@"; do
  case "$arg" in
    --stat) STAT_ONLY=true ;;
    *) DIR="$arg" ;;
  esac
done

DIR="${DIR:-docs}"

if [[ ! -d "$DIR" ]]; then
  echo "Directory not found: $DIR" >&2
  exit 1
fi

# Patterns (PCRE):
#   \xe2\x80\x94  = em-dash (—)
#   \xe2\x80\x93  = en-dash (–)
#   \xe2\x80\x9c  = left double quote (")
#   \xe2\x80\x9d  = right double quote (")
#   \xe2\x80\x98  = left single quote (')
#   \xe2\x80\x99  = right single quote (')
#    --            = space-padded double hyphen (" -- ") used as dash.
#                    Bare -- without spaces is excluded (CLI flags, YAML
#                    frontmatter, table separators). AI almost always
#                    space-pads its dashes. "| -- " is excluded (empty
#                    table cells).
PATTERN='\xe2\x80[\x93\x94\x98\x99\x9c\x9d]|(?<!\|) -- '

file_count=0
hit_count=0

while IFS= read -r -d '' file; do
  # Skip files inside code fences — match only outside fences.
  # Simple approach: grep the whole file; code-fence false positives
  # are acceptable for a heuristic tool.
  matches=$(grep -cP "$PATTERN" "$file" 2>/dev/null || true)
  if [[ "$matches" -gt 0 ]]; then
    file_count=$((file_count + 1))
    hit_count=$((hit_count + matches))

    rel="${file#./}"
    if [[ "$STAT_ONLY" == true ]]; then
      printf "  %3d  %s\n" "$matches" "$rel"
    else
      echo ""
      echo "--- $rel ($matches matches) ---"
      grep -nP "$PATTERN" "$file" 2>/dev/null | while IFS= read -r line; do
        echo "  $line"
      done
    fi
  fi
done < <(find "$DIR" -name '*.md' -print0 | sort -z)

echo ""
if [[ "$file_count" -eq 0 ]]; then
  echo "No AI typography signals found."
else
  echo "Found $hit_count matches across $file_count file(s)."
fi

exit 0
