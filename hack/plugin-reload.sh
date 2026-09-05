#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# Rebuild the cached ctx plugin from local source so Claude Code
# picks up changes without a version bump or restart.

set -euo pipefail

# Resolve paths.
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ASSETS_DIR="$PROJECT_ROOT/internal/assets/claude"

VERSION="$(cat "$PROJECT_ROOT/VERSION" | tr -d '[:space:]')"
CACHE_DIR="$HOME/.claude/plugins/cache/activememory-ctx/ctx/$VERSION"

# Stage the new cache first, then swap atomically: a failure while
# staging must never destroy the existing cache (a half-built cache
# breaks the next Claude Code session).
STAGE="$(mktemp -d "${TMPDIR:-/tmp}/ctx-plugin-reload.XXXXXX")"
trap 'rm -rf "$STAGE"' EXIT

# Mirror the entire plugin root (plugin.json, hooks/, skills/ with
# references/, .mcp.json, CLAUDE.md, dual-manifest files) so the dev
# cache matches a marketplace install exactly.
cp -R "$ASSETS_DIR/." "$STAGE/"

# Swap.
PARENT_DIR="$HOME/.claude/plugins/cache/activememory-ctx"
rm -rf "$PARENT_DIR"
mkdir -p "$(dirname "$CACHE_DIR")"
mv "$STAGE" "$CACHE_DIR"
trap - EXIT

echo "Rebuilt plugin cache at: $CACHE_DIR"
echo "  .claude-plugin/plugin.json"
echo "  hooks/hooks.json"
echo "  skills/ ($(ls -d "$CACHE_DIR"/skills/*/ | wc -l) skills)"
echo ""
echo "IMPORTANT: Claude Code snapshots hooks at session startup."
echo "You must restart your Claude Code session for changes to take effect."
echo "New sessions will pick up the updated plugin automatically."
