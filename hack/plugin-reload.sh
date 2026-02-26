#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# Clear the cached ctx plugin so Claude Code picks up local changes
# on next restart — no version bump needed.

set -euo pipefail

CACHE_DIR="$HOME/.claude/plugins/cache/activememory-ctx"

if [ -d "$CACHE_DIR" ]; then
    rm -rf "$CACHE_DIR"
    echo "Cleared plugin cache: $CACHE_DIR"
else
    echo "No cache found at $CACHE_DIR (nothing to clear)"
fi

echo ""
echo "Next: restart Claude Code."
echo "The plugin will be re-installed from your local marketplace on startup."
echo ""
echo "If it doesn't load automatically: /plugin -> Install -> ctx"
