#!/bin/bash

# Extract all ctx commands from skill files, hook config, recipes, and agent playbook
# Output format: command\tsource_type\tsource_name

# Skills
for f in internal/assets/claude/skills/*/SKILL.md .claude/skills/*/SKILL.md; do
  [ -f "$f" ] || continue
  # Extract ctx commands (basic pattern)
  grep -oP '`?ctx\s+\K[a-z][a-z-]*' "$f" | while read cmd; do
    echo "$cmd	skill	$(basename $(dirname $f))"
  done
done

# Hooks config
if [ -f internal/assets/claude/hooks/hooks.json ]; then
  jq -r '.hooks[].hooks[].command' internal/assets/claude/hooks/hooks.json 2>/dev/null | \
    grep -oP '`?ctx\s+\K[a-z][a-z-]*' | while read cmd; do
    echo "$cmd	hook	hooks.json"
  done
fi

# Recipes
for f in docs/recipes/*.md; do
  [ -f "$f" ] || continue
  grep -oP '`?ctx\s+\K[a-z][a-z-]*' "$f" | while read cmd; do
    echo "$cmd	recipe	$(basename $f)"
  done
done

# AGENT_PLAYBOOK
for f in internal/assets/context/AGENT_PLAYBOOK.md .context/AGENT_PLAYBOOK.md; do
  [ -f "$f" ] || continue
  grep -oP '`?ctx\s+\K[a-z][a-z-]*' "$f" | while read cmd; do
    echo "$cmd	playbook	$(basename $f)"
  done
done

# Makefiles
for f in Makefile internal/assets/project/Makefile.ctx; do
  [ -f "$f" ] || continue
  grep -oP 'ctx\s+\K[a-z][a-z-]*' "$f" | while read cmd; do
    echo "$cmd	makefile	$(basename $f)"
  done
done
