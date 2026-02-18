#!/usr/bin/env bash
# Replace smart quotes and em-dashes with plain equivalents in all markdown files under docs/
set -euo pipefail

DOCS_DIR="${1:-docs}"

if [[ ! -d "$DOCS_DIR" ]]; then
  echo "Directory not found: $DOCS_DIR" >&2
  exit 1
fi

count=0
while IFS= read -r -d '' file; do
  if grep -qP '\xe2\x80[\x98\x99\x9c\x9d\x94]' "$file" 2>/dev/null; then
    sed -i \
      -e "s/\xe2\x80\x98/'/g" \
      -e "s/\xe2\x80\x99/'/g" \
      -e "s/\xe2\x80\x9c/\"/g" \
      -e "s/\xe2\x80\x9d/\"/g" \
      -e "s/\xe2\x80\x94/--/g" \
      "$file"
    echo "fixed: $file"
    ((count++))
  fi
done < <(find "$DOCS_DIR" -name '*.md' -print0)

echo "Done. Fixed $count file(s)."
