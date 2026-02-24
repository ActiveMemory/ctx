#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Export all blob entries from the scratchpad to a directory as files.
# Each blob's label becomes the filename. A unix timestamp is prepended
# to avoid collisions with existing files.
#
# Usage: hack/pad-export-blobs.sh [dir]
#   dir defaults to ./ideas

set -uo pipefail

dir="${1:-./ideas}"
ts="$(date +%s)"

mkdir -p "$dir"

# Parse pad list output: "  N. label [BLOB]"
exported=0
skipped=0

ctx pad 2>/dev/null | while IFS= read -r line; do
  # Only process blob entries
  echo "$line" | grep -q '\[BLOB\]$' || continue

  # Extract entry number and label
  #   "  25. settings.local.json [BLOB]"  →  num=25  label=settings.local.json
  num="$(echo "$line" | sed 's/^ *//' | cut -d. -f1)"
  label="$(echo "$line" | sed 's/^ *[0-9]*\. //' | sed 's/ \[BLOB\]$//')"

  if [ -z "$num" ] || [ -z "$label" ]; then
    echo "  ! could not parse: $line" >&2
    skipped=$((skipped + 1))
    continue
  fi

  # Determine output filename — prepend timestamp to avoid collisions
  outfile="$dir/${ts}-${label}"

  if ctx pad show "$num" --out "$outfile" 2>/dev/null; then
    echo "  + $label → $(basename "$outfile")"
    exported=$((exported + 1))
  else
    echo "  ! failed to export entry $num: $label" >&2
    skipped=$((skipped + 1))
  fi
done

echo ""
echo "Done. Exported $exported, skipped $skipped."
