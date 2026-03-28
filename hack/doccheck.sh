#!/bin/bash
# Scan all Go files for exported functions with docstring violations.
# Convention: // FunctionName does X. (first line starts with function name)
# Must have Parameters: and Returns: sections if function has params/returns.

violations=0

find internal/ cmd/ -name '*.go' ! -name '*_test.go' ! -name 'doc.go' | sort | while read -r file; do
  # Extract exported function signatures with line numbers
  grep -n '^func [A-Z]' "$file" | while IFS=: read -r lineno rest; do
    funcname=$(echo "$rest" | sed 's/^func \([A-Za-z0-9_]*\).*/\1/')
    
    # Check the line before for docstring
    prev=$((lineno - 1))
    if [ "$prev" -lt 1 ]; then
      echo "MISSING_DOC $file:$lineno $funcname (no preceding line)"
      continue
    fi
    
    prevline=$(sed -n "${prev}p" "$file")
    
    # Walk back to find first line of comment block
    docstart=$prev
    while [ "$docstart" -gt 1 ]; do
      checkline=$(sed -n "$((docstart-1))p" "$file")
      if echo "$checkline" | grep -q '^//'; then
        docstart=$((docstart-1))
      else
        break
      fi
    done
    
    # Check if there's a docstring at all
    if ! echo "$prevline" | grep -q '^//'; then
      echo "MISSING_DOC $file:$lineno $funcname"
      continue
    fi
    
    # Check first line starts with function name
    firstdoc=$(sed -n "${docstart}p" "$file")
    if ! echo "$firstdoc" | grep -q "^// $funcname "; then
      echo "BAD_FIRST_LINE $file:$lineno $funcname -- got: $firstdoc"
      continue
    fi
    
    # Check for Parameters section (if function has params)
    params=$(echo "$rest" | sed 's/^func [A-Za-z0-9_]*(//; s/).*//')
    if [ -n "$params" ] && [ "$params" != ")" ]; then
      docblock=$(sed -n "${docstart},$((lineno-1))p" "$file")
      if ! echo "$docblock" | grep -q '// Parameters:'; then
        echo "MISSING_PARAMS $file:$lineno $funcname"
        continue
      fi
    fi
    
    # Check for Returns section (if function has return values)
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
