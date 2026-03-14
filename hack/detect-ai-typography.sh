#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# detect-ai-typography.sh — find files likely AI-edited but not human-reviewed.
#
# Scans files for em-dashes, smart quotes, "--" (double hyphen used as dash),
# and quad backticks (````). These are heuristic signals: almost all LLM output
# uses typographic punctuation from its training corpus, and AI frequently wraps
# code fences in quad backticks which breaks zensical rendering.
#
# False positives are possible (em-dash is valid typography). False negatives
# are unlikely given current model behavior.
#
# Usage:
#   ./hack/detect-ai-typography.sh [dir]              # default: docs/, *.md
#   ./hack/detect-ai-typography.sh internal --ext go   # scan .go files
#   ./hack/detect-ai-typography.sh --ext "go,md,txt"   # multiple extensions
#   ./hack/detect-ai-typography.sh --stat              # summary only (no line detail)

set -euo pipefail

STAT_ONLY=false
DIR=""
EXT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --stat) STAT_ONLY=true; shift ;;
    --ext) EXT="$2"; shift 2 ;;
    *) DIR="$1"; shift ;;
  esac
done

DIR="${DIR:-docs}"
EXT="${EXT:-md}"

if [[ ! -d "$DIR" ]]; then
  echo "Directory not found: $DIR" >&2
  exit 1
fi

# Build find -name arguments from comma-separated extensions.
FIND_ARGS=()
first=true
IFS=',' read -ra EXTS <<< "$EXT"
for ext in "${EXTS[@]}"; do
  ext="${ext## }"  # trim leading space
  ext="${ext%% }"  # trim trailing space
  if [[ "$first" == true ]]; then
    FIND_ARGS+=(-name "*.${ext}")
    first=false
  else
    FIND_ARGS+=(-o -name "*.${ext}")
  fi
done

# Wrap in parens when multiple extensions.
if [[ ${#EXTS[@]} -gt 1 ]]; then
  FIND_ARGS=("(" "${FIND_ARGS[@]}" ")")
fi

# Patterns (PCRE with Unicode escapes):
#   \x{2014}  = em-dash (—)
#   \x{2013}  = en-dash (–)
#   \x{201C}  = left double quote (\u201c)
#   \x{201D}  = right double quote (\u201d)
#   \x{2018}  = left single quote (\u2018)
#   \x{2019}  = right single quote (\u2019)
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
done < <(find "$DIR" "${FIND_ARGS[@]}" -print0 | sort -z)

echo ""
if [[ "$file_count" -eq 0 ]]; then
  echo "No AI typography signals found."
else
  echo "Found $hit_count matches across $file_count file(s)."
fi

exit 0
