#!/bin/sh
# ctx: post-commit hook for recording commit context history.
# Requires: ctx on $PATH
# Installed by: ctx trace hook enable
# Remove with:  ctx trace hook disable

COMMIT_HASH=$(git rev-parse HEAD)
ctx trace collect --record "$COMMIT_HASH" 2>/dev/null || true
