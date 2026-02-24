#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Import lines from a file into the scratchpad.
# Each non-empty line becomes a separate entry.
#
# Usage: hack/pad-import.sh <file>

set -euo pipefail

if [ $# -ne 1 ]; then
  echo "Usage: $0 <file>" >&2
  exit 1
fi

file="$1"

if [ ! -f "$file" ]; then
  echo "Error: file not found: $file" >&2
  exit 1
fi

count=0
while IFS= read -r line || [ -n "$line" ]; do
  # Skip empty lines
  [ -z "$line" ] && continue
  ctx pad add "$line"
  count=$((count + 1))
done < "$file"

echo "Added $count entries to scratchpad."
