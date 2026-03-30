#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# lint-docstrings.sh — scan all Go files for docstring convention violations.
#
# Convention (from CONVENTIONS.md):
#   // FunctionName does X.
#   //
#   // Parameters:
#   //   - param1: Description
#   //
#   // Returns:
#   //   - Type: Description
#
# Checks:
#   MISSING_DOC    — exported function has no preceding comment
#   BAD_FIRST_LINE — docstring doesn't start with "// FunctionName "
#   MISSING_PARAMS — function has parameters but no Parameters: section
#   MISSING_RETURNS — function has return values but no Returns: section
#   MISSING_FIELDS — exported struct with 2+ fields has no Fields: section

set -euo pipefail

find internal/ cmd/ -name '*.go' ! -name '*_test.go' ! -name 'doc.go' | sort | while read -r file; do
  grep -n '^func [a-zA-Z]' "$file" | while IFS=: read -r lineno rest; do
    funcname=$(echo "$rest" | sed 's/^func \([A-Za-z0-9_]*\).*/\1/')

    prev=$((lineno - 1))
    if [ "$prev" -lt 1 ]; then
      echo "MISSING_DOC $file:$lineno $funcname (no preceding line)"
      continue
    fi

    prevline=$(sed -n "${prev}p" "$file")

    docstart=$prev
    while [ "$docstart" -gt 1 ]; do
      checkline=$(sed -n "$((docstart-1))p" "$file")
      if echo "$checkline" | grep -q '^//'; then
        docstart=$((docstart-1))
      else
        break
      fi
    done

    if ! echo "$prevline" | grep -q '^//'; then
      echo "MISSING_DOC $file:$lineno $funcname"
      continue
    fi

    firstdoc=$(sed -n "${docstart}p" "$file")
    if ! echo "$firstdoc" | grep -q "^// $funcname "; then
      echo "BAD_FIRST_LINE $file:$lineno $funcname -- got: $firstdoc"
      continue
    fi

    params=$(echo "$rest" | sed 's/^func [A-Za-z0-9_]*(//; s/).*//')
    if [ -n "$params" ] && [ "$params" != ")" ]; then
      docblock=$(sed -n "${docstart},$((lineno-1))p" "$file")
      if ! echo "$docblock" | grep -q '// Parameters:'; then
        echo "MISSING_PARAMS $file:$lineno $funcname"
        continue
      fi
    fi

    # Skip return check for multi-line signatures (no closing paren on
    # this line) or func-typed params (nested parens confuse extraction).
    if ! echo "$rest" | grep -q ') {' && ! echo "$rest" | grep -q ') ('; then
      continue
    fi
    if echo "$rest" | grep -q ' func('; then
      continue
    fi
    returnpart=$(echo "$rest" | sed 's/^func [A-Za-z0-9_]*([^)]*) //')
    # Guard: if sed didn't match (returnpart unchanged), skip
    if [ "$returnpart" = "$rest" ]; then
      continue
    fi
    if echo "$returnpart" | grep -qE '(error|string|int|bool|\*|\[|map)'; then
      docblock=$(sed -n "${docstart},$((lineno-1))p" "$file")
      if ! echo "$docblock" | grep -q '// Returns:'; then
        echo "MISSING_RETURNS $file:$lineno $funcname"
        continue
      fi
    fi
  done || true

  # Check exported structs for Fields: section (2+ fields required).
  grep -n '^type [A-Z].*struct {' "$file" | while IFS=: read -r lineno rest; do
    typename=$(echo "$rest" | sed 's/^type \([A-Za-z0-9_]*\).*/\1/')

    # Count fields: lines between "struct {" and next "}" that look like fields.
    closing=$(awk -v start="$lineno" 'NR > start && /^}/ { print NR; exit }' "$file")
    if [ -z "$closing" ]; then
      continue
    fi
    fieldcount=$(sed -n "$((lineno+1)),$((closing-1))p" "$file" \
      | grep -cP '^\t[A-Z]' || true)
    if [ "$fieldcount" -lt 2 ]; then
      continue
    fi

    # Walk back from the type line to find the docblock.
    prev=$((lineno - 1))
    if [ "$prev" -lt 1 ]; then
      continue
    fi
    prevline=$(sed -n "${prev}p" "$file")
    if ! echo "$prevline" | grep -q '^//'; then
      continue
    fi

    docstart=$prev
    while [ "$docstart" -gt 1 ]; do
      checkline=$(sed -n "$((docstart-1))p" "$file")
      if echo "$checkline" | grep -q '^//'; then
        docstart=$((docstart-1))
      else
        break
      fi
    done

    docblock=$(sed -n "${docstart},$((lineno-1))p" "$file")
    if echo "$docblock" | grep -q '// Fields:'; then
      continue
    fi
    # Accept inline field comments as alternative to Fields: section.
    # Count fields with a preceding or same-line comment.
    inlinecount=$(sed -n "$((lineno+1)),$((closing-1))p" "$file" \
      | grep -cP '^\t// [A-Z]|^\t[A-Z].*//\s' || true)
    if [ "$inlinecount" -ge "$fieldcount" ]; then
      continue
    fi
    echo "MISSING_FIELDS $file:$lineno $typename ($fieldcount fields)"
  done || true
done 2>/dev/null

# Check doc.go files for minimum depth.
# Count only doc comment lines (// Package ... through package line),
# excluding the copyright header. Require at least 4 doc comment lines
# beyond the one-liner "// Package X does Y." + "package X".
find internal/ cmd/ -name 'doc.go' | sort | while read -r file; do
  # Count lines from "// Package" to end of file (skips copyright header).
  doclines=$(sed -n '/^\/\/ Package/,$p' "$file" | wc -l)
  if [ "$doclines" -lt 6 ]; then
    pkg=$(grep '^package ' "$file" | head -1 | sed 's/package //')
    echo "SHALLOW_DOC $file ($doclines doc lines, package $pkg)"
  fi
done 2>/dev/null
