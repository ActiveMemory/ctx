#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# lint-mixed-funcs.sh — find files with both exported and private functions.
#
# Convention: public API in main file, private helpers in separate files.
# Files in cmd/ run.go are exempt (dispatch helpers are expected there).

set -euo pipefail

find internal/ cmd/ -name '*.go' ! -name '*_test.go' ! -name 'doc.go' -print0 | xargs -0 -I{} sh -c '
  file="$1"
  pub=$(grep -c "^func [A-Z]" "$file" 2>/dev/null)
  priv=$(grep -c "^func [a-z]" "$file" 2>/dev/null)
  pub=${pub:-0}; priv=${priv:-0}
  if [ "$pub" -gt 0 ] && [ "$priv" -gt 0 ]; then
    # Exempt cmd/*/run.go files (dispatch helpers are expected)
    case "$file" in */cmd/*/run.go) exit 0 ;; esac
    privnames=$(grep "^func [a-z]" "$file" | sed "s/^func \([a-zA-Z0-9_]*\).*/\1/" | tr "\n" "," | sed "s/,$//" )
    echo "MIXED $file pub=$pub priv=$priv [$privnames]"
  fi
' _ {} 2>/dev/null | sort
