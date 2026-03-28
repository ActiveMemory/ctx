#!/bin/bash
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

    returnpart=$(echo "$rest" | sed 's/.*) //')
    if echo "$returnpart" | grep -qE '(error|string|int|bool|\*|\[|map)'; then
      docblock=$(sed -n "${docstart},$((lineno-1))p" "$file")
      if ! echo "$docblock" | grep -q '// Returns:'; then
        echo "MISSING_RETURNS $file:$lineno $funcname"
        continue
      fi
    fi
  done
done 2>/dev/null
