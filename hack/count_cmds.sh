#!/bin/bash

# Initialize associative arrays
declare -A all_cmds
declare -A by_skill
declare -A by_hook
declare -A by_recipe
declare -A by_playbook
declare -A by_makefile

# Skills
for f in internal/assets/claude/skills/*/SKILL.md .claude/skills/*/SKILL.md; do
  [ -f "$f" ] || continue
  skill_name=$(basename $(dirname "$f"))
  grep -oP '(?:^|[^a-zA-Z])ctx\s+\K[a-z]+(?:-[a-z]+)*(?=[\s\(]|$)' "$f" | while read cmd; do
    echo "$cmd|skill|$skill_name"
  done
done

# Hooks
if [ -f internal/assets/claude/hooks/hooks.json ]; then
  jq -r '.hooks[].hooks[].command' internal/assets/claude/hooks/hooks.json 2>/dev/null | \
    grep -oP 'ctx\s+\K[a-z]+(?:-[a-z]+)*(?=[\s\(]|$)' | while read cmd; do
    echo "$cmd|hook|hooks.json"
  done
fi

# Recipes
for f in docs/recipes/*.md; do
  [ -f "$f" ] || continue
  recipe_name=$(basename "$f")
  grep -oP '`ctx\s+\K[a-z]+(?:-[a-z]+)*(?=[\s\(]|`|$)' "$f" | while read cmd; do
    echo "$cmd|recipe|$recipe_name"
  done
done

# AGENT_PLAYBOOK
for f in internal/assets/context/AGENT_PLAYBOOK.md .context/AGENT_PLAYBOOK.md; do
  [ -f "$f" ] || continue
  grep -oP '`ctx\s+\K[a-z]+(?:-[a-z]+)*(?=[\s\(]|`|$)' "$f" | while read cmd; do
    echo "$cmd|playbook|AGENT_PLAYBOOK.md"
  done
done

# Makefiles
for f in Makefile internal/assets/project/Makefile.ctx; do
  [ -f "$f" ] || continue
  makefile_name=$(basename "$f")
  grep -oP 'ctx\s+\K[a-z]+(?:-[a-z]+)*(?=[\s\(]|$)' "$f" | while read cmd; do
    echo "$cmd|makefile|$makefile_name"
  done
done
