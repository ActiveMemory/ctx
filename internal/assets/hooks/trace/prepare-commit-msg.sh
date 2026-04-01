#!/bin/sh
# ctx: prepare-commit-msg hook for commit context tracing.
# Installed by: ctx trace hook enable
# Remove with:  ctx trace hook disable
# Requires:     ctx on $PATH

COMMIT_MSG_FILE="$1"
COMMIT_SOURCE="$2"

# Only inject on normal commits (not merges, squashes, or amends)
case "$COMMIT_SOURCE" in
  merge|squash) exit 0 ;;
esac

# Collect context refs (requires ctx on $PATH)
TRAILER=$(ctx trace collect 2>/dev/null)

if [ -n "$TRAILER" ]; then
  # Append trailer with a blank line separator
  echo "" >> "$COMMIT_MSG_FILE"
  echo "$TRAILER" >> "$COMMIT_MSG_FILE"
fi
