#!/usr/bin/env bash
# Replace smart quotes and em-dashes with plain equivalents in all markdown
# files under docs/
#
# Note that this is a quick hack, and it's not the best way to do this.
#
# Instead ask the agent to do semantic replacements as such:
#
# "Run ./hack/detect-ai-typography <params>; read every file fully and do
# semantic, editorial changes. DO NOT BLINDLY CLEAN THE FILES UP.
# Understand the context and do semantic replacements. -- if you are going to
# write a script to blindly sed them STOP RIGHT THERE! -- sometimes you might
# need to replace it with a ":"; sometimes with a parenthesis, sometimes with a
# regular dash, sometimes rephrasing. -- Your constitution forbids you
# from being lazy!"

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
