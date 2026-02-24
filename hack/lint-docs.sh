#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# lint-docs.sh — verify doc.go file organization listings match reality.
#
# For every doc.go that contains a "File Organization" section, checks
# that every non-test .go file in the same directory is listed, and that
# no listed file is missing from disk.
#
# Exit code: number of issues found (0 = clean).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

issues=0

for docfile in $(find . -name doc.go -not -path './vendor/*'); do
  # Skip if no File Organization section
  grep -q 'File Organization' "$docfile" || continue

  dir=$(dirname "$docfile")

  # Get actual .go files (excluding doc.go itself and _test.go)
  actual=$(ls "$dir"/*.go 2>/dev/null \
    | xargs -n1 basename \
    | grep -v '_test\.go$' \
    | grep -v '^doc\.go$' \
    | sort)

  # Get listed files from doc.go (lines matching "  - filename.go:")
  # Exclude _test.go entries — tests aren't part of file organization.
  listed=$(grep -oE '^\s*//\s+-\s+\S+\.go' "$docfile" \
    | sed 's|.*- ||' \
    | sed 's|:.*||' \
    | grep -v '_test\.go$' \
    | sort)

  # Find files on disk but not in doc
  missing=$(comm -23 <(echo "$actual") <(echo "$listed"))
  # Find files in doc but not on disk
  extra=$(comm -13 <(echo "$actual") <(echo "$listed"))

  if [ -n "$missing" ]; then
    for f in $missing; do
      echo "$docfile: missing from listing: $f"
      issues=$((issues + 1))
    done
  fi

  if [ -n "$extra" ]; then
    for f in $extra; do
      echo "$docfile: listed but not on disk: $f"
      issues=$((issues + 1))
    done
  fi
done

if [ "$issues" -eq 0 ]; then
  echo "lint-docs: clean"
else
  echo ""
  echo "lint-docs: $issues issues found"
fi

exit "$issues"
