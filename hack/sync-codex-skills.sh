#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# sync-codex-skills.sh — sync Codex plugin skills from canonical ctx skills.
#
# ctx skills (internal/assets/claude/skills/) are the source of truth.
# Codex skills (internal/assets/codex/skills/) are generated from them
# with the `allowed-tools` frontmatter key stripped (Claude Code-specific;
# Codex skills have no tool-permission frontmatter).
#
# Unlike the Copilot sync, this is a full mirror: every ctx skill that
# is not on the exclusion list below is (re)generated, and Codex skill
# directories whose ctx counterpart disappeared are removed. The
# exclusion list names skills whose body only makes sense inside
# Claude Code (its settings files, plan files, or headless runner).

set -euo pipefail

CTX_SKILLS="internal/assets/claude/skills"
CODEX_SKILLS="internal/assets/codex/skills"

# Claude-only skills: operate on Claude Code-specific state.
EXCLUDE=(
  ctx-permission-sanitize  # audits .claude/settings.local.json
  ctx-plan-import          # imports ~/.claude/plans/
  ctx-dream                # headless `claude -p` cron + guard.sh
  ctx-skill-create         # authors Claude Code skills/plugins
)

excluded() {
  local name="$1"
  for x in "${EXCLUDE[@]}"; do
    [ "$x" = "$name" ] && return 0
  done
  return 1
}

mkdir -p "$CODEX_SKILLS"

synced=0
removed=0
skipped=0

for ctx_dir in "$CTX_SKILLS"/*/; do
  skill_name=$(basename "$ctx_dir")
  ctx_skill="$ctx_dir/SKILL.md"
  [ -f "$ctx_skill" ] || continue

  if excluded "$skill_name"; then
    skipped=$((skipped + 1))
    continue
  fi

  mkdir -p "$CODEX_SKILLS/$skill_name"
  # Strip `allowed-tools:` line from frontmatter (Claude Code-specific).
  sed '/^allowed-tools:/d' "$ctx_skill" > "$CODEX_SKILLS/$skill_name/SKILL.md"

  # Mirror the skill's references/ directory (skill bodies cite these
  # files; shipping SKILL.md alone would point agents at 404s).
  rm -rf "$CODEX_SKILLS/$skill_name/references"
  if [ -d "$ctx_dir/references" ]; then
    cp -R "$ctx_dir/references" "$CODEX_SKILLS/$skill_name/references"
  fi
  synced=$((synced + 1))
done

# Remove Codex skills whose ctx counterpart is gone or now excluded.
for codex_dir in "$CODEX_SKILLS"/*/; do
  [ -d "$codex_dir" ] || continue
  skill_name=$(basename "$codex_dir")
  if [ ! -f "$CTX_SKILLS/$skill_name/SKILL.md" ] || excluded "$skill_name"; then
    rm -rf "$codex_dir"
    removed=$((removed + 1))
  fi
done

# The Claude plugin root is dual-manifest: its .codex-plugin/plugin.json
# points Codex at hooks/codex.json, so a Codex that resolves the legacy
# .claude-plugin marketplace still gets working hooks. Keep that file a
# byte-copy of the canonical Codex manifest.
cp internal/assets/codex/hooks/hooks.json internal/assets/claude/hooks/codex.json
jq '.hooks = "./hooks/codex.json"' internal/assets/codex/.codex-plugin/plugin.json \
  > internal/assets/claude/.codex-plugin/plugin.json

echo "Codex skills synced: $synced updated, $skipped Claude-only (excluded), $removed removed."
