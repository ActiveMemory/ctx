#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# sync-copilot-skills.sh — sync Copilot CLI skills from canonical ctx skills.
#
# ctx skills (internal/assets/claude/skills/) are the source of truth.
# Copilot CLI skills (internal/assets/integrations/copilot-cli/skills/) are
# generated from them with the `allowed-tools` frontmatter key stripped
# (Claude Code-specific, not applicable to Copilot).
#
# Skills that exist only in the Copilot directory (no ctx counterpart)
# are left untouched.

set -euo pipefail

CTX_SKILLS="internal/assets/claude/skills"
COPILOT_SKILLS="internal/assets/integrations/copilot-cli/skills"

synced=0
skipped=0

for copilot_dir in "$COPILOT_SKILLS"/*/; do
  skill_name=$(basename "$copilot_dir")
  ctx_skill="$CTX_SKILLS/$skill_name/SKILL.md"
  copilot_skill="$copilot_dir/SKILL.md"

  if [ ! -f "$ctx_skill" ]; then
    # No ctx counterpart — Copilot-only skill, leave untouched.
    skipped=$((skipped + 1))
    continue
  fi

  # Strip `allowed-tools:` line from frontmatter (Claude Code-specific).
  sed '/^allowed-tools:/d' "$ctx_skill" > "$copilot_skill"
  synced=$((synced + 1))
done

echo "Copilot skills synced: $synced updated, $skipped Copilot-only (unchanged)."
