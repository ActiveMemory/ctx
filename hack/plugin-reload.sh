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

# Clear old cache.
PARENT_DIR="$HOME/.claude/plugins/cache/activememory-ctx"
if [ -d "$PARENT_DIR" ]; then
    rm -rf "$PARENT_DIR"
    echo "Cleared old cache: $PARENT_DIR"
fi

# Rebuild from source assets.
mkdir -p "$CACHE_DIR/.claude-plugin"
mkdir -p "$CACHE_DIR/hooks"
mkdir -p "$CACHE_DIR/skills"

cp "$ASSETS_DIR/.claude-plugin/plugin.json" "$CACHE_DIR/.claude-plugin/"
cp "$ASSETS_DIR/hooks/hooks.json" "$CACHE_DIR/hooks/"

# Copy all skills.
for skill_dir in "$ASSETS_DIR"/skills/*/; do
    skill_name="$(basename "$skill_dir")"
    mkdir -p "$CACHE_DIR/skills/$skill_name"
    cp "$skill_dir"SKILL.md "$CACHE_DIR/skills/$skill_name/"
done

echo "Rebuilt plugin cache at: $CACHE_DIR"
echo "  .claude-plugin/plugin.json"
echo "  hooks/hooks.json"
echo "  skills/ ($(ls -d "$CACHE_DIR"/skills/*/ | wc -l) skills)"
echo ""
echo "IMPORTANT: Claude Code snapshots hooks at session startup."
echo "You must restart your Claude Code session for changes to take effect."
echo "New sessions will pick up the updated plugin automatically."
