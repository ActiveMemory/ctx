#!/bin/bash

echo "=== VERIFYING ORPHANED COMMAND STATUS ==="
echo ""

orphaned_cmds=(
  "archive"
  "brainstorm"
  "check-links"
  "commit"
  "consolidate"
  "explain"
  "import-plans"
  "reflect"
)

for cmd in "${orphaned_cmds[@]}"; do
  echo "$cmd:"
  if [ -d "internal/assets/claude/skills/ctx-$cmd" ]; then
    echo "  - Has skill: internal/assets/claude/skills/ctx-$cmd/SKILL.md"
    # Check if the skill is self-referential or meta
    if grep -q "## Execution" "internal/assets/claude/skills/ctx-$cmd/SKILL.md" 2>/dev/null; then
      echo "  - Is a proper skill"
    fi
  elif [ -d ".claude/skills/ctx-$cmd" ]; then
    echo "  - Has skill: .claude/skills/ctx-$cmd/SKILL.md"
  fi
  
  # Check if it's referenced in its own SKILL.md
  found_self=$(grep -r "ctx $cmd" internal/assets/claude/skills/ctx-$cmd .claude/skills/ctx-$cmd 2>/dev/null | grep -v "description\|name:" | wc -l)
  if [ "$found_self" -gt 0 ]; then
    echo "  - Self-references in own SKILL.md: Yes ($found_self occurrences)"
  else
    echo "  - Self-references: No (never invoked by its own skill)"
  fi
  echo ""
done

