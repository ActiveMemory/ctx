#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# sync-agy-skills.sh — sync Antigravity plugin skills from canonical ctx skills.
#
# ctx skills (internal/assets/claude/skills/) are the source of truth.
# Antigravity skills (internal/assets/agy/skills/) are generated from them
# with the `allowed-tools` frontmatter key stripped (Claude Code-specific).
#
# Skills whose markdown body exceeds Antigravity's 24,000-byte validator cap
# (>22,000 bytes total) are deterministically split at the last `## ` heading
# before line 450 into references/extended.md.

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

CTX_SKILLS="internal/assets/claude/skills"
AGY_SKILLS="internal/assets/agy/skills"
AGY_RULES="internal/assets/agy/rules"

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

mkdir -p "$AGY_RULES"
cp internal/assets/integrations/agents.md "$AGY_RULES/AGENTS.md"

mkdir -p "$AGY_SKILLS"

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

  mkdir -p "$AGY_SKILLS/$skill_name"
  # Strip `allowed-tools:` line from frontmatter (Claude Code-specific).
  sed '/^allowed-tools:/d' "$ctx_skill" > "$AGY_SKILLS/$skill_name/SKILL.md"

  # Mirror the skill's references/ directory (skill bodies cite these
  # files; shipping SKILL.md alone would point agents at 404s).
  rm -rf "$AGY_SKILLS/$skill_name/references"
  if [ -d "$ctx_dir/references" ]; then
    cp -R "$ctx_dir/references" "$AGY_SKILLS/$skill_name/references"
  fi

  # Antigravity strict validator imposes a 24,000-byte limit on SKILL.md body.
  # For files exceeding 22,000 bytes, split at the last ## heading
  # before line 450.
  skill_file="$AGY_SKILLS/$skill_name/SKILL.md"
  file_size=$(wc -c < "$skill_file" | tr -d '[:space:]')
  if [ "$file_size" -gt 22000 ]; then
    split_line=$(awk \
      'NR < 450 && /^## / { split_line = NR } END { print split_line }' \
      "$skill_file")
    if [ -n "$split_line" ] && [ "$split_line" -gt 1 ]; then
      mkdir -p "$AGY_SKILLS/$skill_name/references"
      extended_file="$AGY_SKILLS/$skill_name/references/extended.md"

      # Write overflow section into references/extended.md
      sed -n "${split_line},\$p" "$skill_file" > "$extended_file"

      # Write head section and extended reference pointer into SKILL.md
      head_n=$((split_line - 1))
      tmp_head="${skill_file}.tmp"
      head -n "$head_n" "$skill_file" > "$tmp_head"
      cat << 'EOF' >> "$tmp_head"

## Extended Reference

Read `references/extended.md` in this skill directory for the remaining specification and detailed steps.
EOF
      mv "$tmp_head" "$skill_file"
    fi
  fi

  synced=$((synced + 1))
done

# Remove Antigravity skills whose ctx counterpart is gone or now excluded.
for agy_dir in "$AGY_SKILLS"/*/; do
  [ -d "$agy_dir" ] || continue
  skill_name=$(basename "$agy_dir")
  if [ ! -f "$CTX_SKILLS/$skill_name/SKILL.md" ] || excluded "$skill_name"; then
    rm -rf "$agy_dir"
    removed=$((removed + 1))
  fi
done

echo "Antigravity skills synced: $synced updated, $skipped Claude-only (excluded), $removed removed."
