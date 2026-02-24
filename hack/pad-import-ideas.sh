#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Import first-level files from ./ideas/ into the scratchpad as blobs.
# Each file becomes a pad entry with the filename as the label.
# Skips directories, only processes regular files.
#
# Usage: hack/pad-import-ideas.sh [dir]
#   dir defaults to ./ideas

set -uo pipefail

dir="${1:-./ideas}"

if [ ! -d "$dir" ]; then
  echo "Error: directory not found: $dir" >&2
  exit 1
fi

added=0
skipped=0

for file in "$dir"/*; do
  # Skip directories and non-regular files
  [ -f "$file" ] || continue

  slug="$(basename "$file")"

  if ctx pad add "$slug" --file "$file" 2>/dev/null; then
    echo "  + $slug"
    added=$((added + 1))
  else
    echo "  ! skipped: $slug" >&2
    skipped=$((skipped + 1))
  fi
done

echo ""
echo "Done. Added $added, skipped $skipped."
