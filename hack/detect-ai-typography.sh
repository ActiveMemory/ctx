#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# detect-ai-typography.sh — find docs likely AI-edited but not human-reviewed.
#
# Scans markdown files for em-dashes, smart quotes, "--" (double hyphen
# used as dash), and quad backticks (````). These are heuristic signals:
# almost all LLM output uses typographic punctuation from its training
# corpus, and AI frequently wraps code fences in quad backticks which
# breaks zensical rendering.
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

# Patterns (PCRE with Unicode escapes):
#   \x{2014}  = em-dash (—)
#   \x{2013}  = en-dash (–)
#   \x{201C}  = left double quote (")
#   \x{201D}  = right double quote (")
#   \x{2018}  = left single quote (')
#   \x{2019}  = right single quote (')
#    --       = space-padded double hyphen (" -- ") used as dash.
#               Bare -- without spaces is excluded (CLI flags, YAML
#               frontmatter, table separators). AI almost always
#               space-pads its dashes. "| -- " is excluded (empty
#               table cells).
#   ````      = quad backtick. AI wraps code fences in four-backtick
#               blocks; zensical doesn't support them. Triple is the
#               project maximum.
PATTERN='\x{2013}|\x{2014}|\x{2018}|\x{2019}|\x{201C}|\x{201D}|(?<!\|) -- |````'

# Files where typographic punctuation is intentional.
# Add glob patterns here to skip specific paths.
EXCLUDE_PATTERNS=()

file_count=0
hit_count=0

while IFS= read -r -d '' file; do
  # Skip excluded files (formal/academic content where typography is intentional).
  skip=false
  for excl in "${EXCLUDE_PATTERNS[@]}"; do
    # shellcheck disable=SC2254
    case "$file" in $excl) skip=true; break ;; esac
  done
  if [[ "$skip" == true ]]; then continue; fi

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
