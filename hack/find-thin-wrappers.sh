#!/bin/bash
# Find thin wrapper functions: exported functions whose body is a single
# return statement delegating to another package's function.

find internal/ cmd/ -name '*.go' ! -name '*_test.go' ! -name 'doc.go' | sort | while read -r file; do
  # Use awk to find functions where the body is just "return pkg.Func(...)"
  awk '
    /^func [A-Z]/ {
      funcline = $0
      funcname = $2
      sub(/\(.*/, "", funcname)
      linenum = NR
      body_lines = 0
      in_func = 1
      brace_depth = 0
      body = ""
      next
    }
    in_func {
      # Track brace depth
      for (i = 1; i <= length($0); i++) {
        c = substr($0, i, 1)
        if (c == "{") brace_depth++
        if (c == "}") brace_depth--
      }
      
      # Collect non-empty, non-brace-only lines
      trimmed = $0
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", trimmed)
      if (trimmed != "{" && trimmed != "}" && trimmed != "") {
        body_lines++
        body = trimmed
      }
      
      if (brace_depth <= 0) {
        # Function ended
        if (body_lines == 1 && body ~ /^return [a-zA-Z]+\./) {
          printf "THIN_WRAPPER %s:%d %s -> %s\n", FILENAME, linenum, funcname, body
        }
        in_func = 0
      }
    }
  ' "$file"
done
