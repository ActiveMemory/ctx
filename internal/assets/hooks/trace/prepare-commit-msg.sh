#!/bin/sh
# ctx: prepare-commit-msg hook for commit context tracing.
# Installed by: ctx trace hook enable
# Remove with:  ctx trace hook disable
# Requires:     ctx on $PATH

COMMIT_MSG_FILE="$1"
COMMIT_SOURCE="$2"

# Only inject on normal commits (not merges, squashes, or amends).
# Amends arrive as COMMIT_SOURCE=commit with a SHA third argument.
case "$COMMIT_SOURCE" in
  merge|squash) exit 0 ;;
  commit) [ -n "${3:-}" ] && exit 0 ;;
esac

# Collect context refs (requires ctx on $PATH)
TRAILER=$(ctx trace collect 2>/dev/null)

# Idempotency: never append a second trailer of the same key
# (protects amends and re-entrant hook runs).
if [ -n "$TRAILER" ] && grep -q "^${TRAILER%%:*}:" "$COMMIT_MSG_FILE" 2>/dev/null; then
  exit 0
fi

if [ -n "$TRAILER" ]; then
  # Append trailer with a blank line separator
  echo "" >> "$COMMIT_MSG_FILE"
  echo "$TRAILER" >> "$COMMIT_MSG_FILE"
fi
