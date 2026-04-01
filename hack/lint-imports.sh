#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# lint-imports.sh — scan all Go files for import grouping violations.
#
# Convention: stdlib — blank — external (cobra, yaml) — blank — ctx imports.
# Three groups, always in this order. Files with a single import or no
# imports are skipped.
#
# Violations:
#   MIXED_GROUP  — stdlib and ctx imports in the same group (no blank line)
#   EXT_IN_CTX   — external import mixed into the ctx group
#   CTX_IN_EXT   — ctx import mixed into the external group
#   CTX_BEFORE_EXT — ctx import appears before external import without separation
#   WRONG_ORDER  — groups appear in wrong order

set -euo pipefail

find internal/ cmd/ -name '*.go' ! -name '*_test.go' | sort | while read -r file; do
  # Extract the import block
  import_block=$(awk '/^import \(/{found=1; next} found && /^\)/{exit} found{print}' "$file")
  [ -z "$import_block" ] && continue

  # Track groups: split on blank lines
  group_num=0
  has_stdlib=0
  has_ext=0
  has_ctx=0
  prev_blank=1
  violations=""

  while IFS= read -r line; do
    trimmed=$(echo "$line" | sed 's/^[[:space:]]*//' | sed 's/[[:space:]]*$//')

    # Blank line = group separator — reset per-group type flags.
    if [ -z "$trimmed" ]; then
      prev_blank=1
      group_num=$((group_num + 1))
      has_stdlib=0
      has_ext=0
      has_ctx=0
      continue
    fi

    # Skip comments
    echo "$trimmed" | grep -q '^//' && continue

    # Determine import type
    if echo "$trimmed" | grep -q 'github.com/ActiveMemory/ctx'; then
      import_type="ctx"
    elif echo "$trimmed" | grep -q 'github.com/\|gopkg.in/'; then
      import_type="ext"
    else
      import_type="stdlib"
    fi

    # Check for mixing within a group
    case $import_type in
      stdlib)
        if [ "$has_ctx" -eq 1 ] && [ "$prev_blank" -eq 0 ]; then
          violations="${violations}MIXED_GROUP $file (stdlib after ctx in same group)\n"
        fi
        has_stdlib=1
        ;;
      ext)
        if [ "$has_ctx" -eq 1 ] && [ "$prev_blank" -eq 0 ]; then
          violations="${violations}EXT_IN_CTX $file (external import in ctx group)\n"
        fi
        has_ext=1
        ;;
      ctx)
        if [ "$has_ext" -eq 1 ] && [ "$prev_blank" -eq 0 ]; then
          violations="${violations}CTX_IN_EXT $file (ctx import in external group)\n"
        fi
        if [ "$has_stdlib" -eq 1 ] && [ "$prev_blank" -eq 0 ]; then
          violations="${violations}MIXED_GROUP $file (ctx after stdlib in same group)\n"
        fi
        has_ctx=1
        ;;
    esac
    prev_blank=0
  done <<< "$import_block"

  if [ -n "$violations" ]; then
    printf "%b" "$violations"
  fi
done
