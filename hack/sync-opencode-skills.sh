#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# sync-opencode-skills.sh — sync OpenCode skills from canonical ctx skills.
#
# ctx skills (internal/assets/claude/skills/) are the source of truth.
# OpenCode skills (internal/assets/integrations/opencode/skills/) are
# generated from them with the `allowed-tools` frontmatter key stripped
# (Claude Code-specific, not applicable to OpenCode).
#
# Enrollment is opt-in by directory presence: a skill syncs iff its
# directory exists under the OpenCode tree. Skills that exist only in
# the OpenCode directory (no ctx counterpart) are left untouched.

set -euo pipefail

CTX_SKILLS="internal/assets/claude/skills"
OPENCODE_SKILLS="internal/assets/integrations/opencode/skills"

synced=0
skipped=0

for opencode_dir in "$OPENCODE_SKILLS"/*/; do
  skill_name=$(basename "$opencode_dir")
  ctx_skill="$CTX_SKILLS/$skill_name/SKILL.md"
  opencode_skill="$opencode_dir/SKILL.md"

  if [ ! -f "$ctx_skill" ]; then
    # No ctx counterpart — OpenCode-only skill, leave untouched.
    skipped=$((skipped + 1))
    continue
  fi

  # Strip `allowed-tools:` line from frontmatter (Claude Code-specific).
  sed '/^allowed-tools:/d' "$ctx_skill" > "$opencode_skill"
  synced=$((synced + 1))
done

echo "OpenCode skills synced: $synced updated, $skipped OpenCode-only (unchanged)."
