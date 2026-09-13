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
# Like the Codex sync (and unlike the Copilot sync), this is a full
# mirror: every ctx skill that is not on the exclusion list below is
# (re)generated, and OpenCode skill directories whose ctx counterpart
# disappeared are removed. The exclusion list names skills whose body
# only makes sense inside Claude Code (its settings files, plan files,
# or headless runner).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

CTX_SKILLS="internal/assets/claude/skills"
OPENCODE_SKILLS="internal/assets/integrations/opencode/skills"

# Claude-only skills: operate on Claude Code-specific state.
EXCLUDE=(
  ctx-permission-sanitize  # audits .claude/settings.local.json
  ctx-plan-import          # imports ~/.claude/plans/
  ctx-dream                # headless `claude -p` cron + guard.sh
  ctx-skill-create         # authors Claude Code skills/plugins
)

excluded() {
  local name="$1"
  # Guard the expansion: bash 3.2 under set -u aborts on an empty
  # array; see the lint-drift.sh guard and LEARNINGS 2026-08-19.
  for x in ${EXCLUDE[@]+"${EXCLUDE[@]}"}; do
    [ "$x" = "$name" ] && return 0
  done
  return 1
}

mkdir -p "$OPENCODE_SKILLS"

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

  mkdir -p "$OPENCODE_SKILLS/$skill_name"
  # Strip `allowed-tools:` line from frontmatter (Claude Code-specific).
  sed '/^allowed-tools:/d' "$ctx_skill" > "$OPENCODE_SKILLS/$skill_name/SKILL.md"

  # Mirror the skill's references/ directory (skill bodies cite these
  # files; shipping SKILL.md alone would point agents at 404s).
  rm -rf "$OPENCODE_SKILLS/$skill_name/references"
  if [ -d "$ctx_dir/references" ]; then
    cp -R "$ctx_dir/references" "$OPENCODE_SKILLS/$skill_name/references"
  fi
  synced=$((synced + 1))
done

# Remove OpenCode skills whose ctx counterpart is gone or now excluded.
for opencode_dir in "$OPENCODE_SKILLS"/*/; do
  [ -d "$opencode_dir" ] || continue
  skill_name=$(basename "$opencode_dir")
  if [ ! -f "$CTX_SKILLS/$skill_name/SKILL.md" ] || excluded "$skill_name"; then
    rm -rf "$opencode_dir"
    removed=$((removed + 1))
  fi
done

echo "OpenCode skills synced: $synced updated, $skipped Claude-only (excluded), $removed removed."
