#!/bin/bash

echo "=== COMMAND USAGE SURFACE ANALYSIS ==="
echo ""
echo "This report maps which CLI commands are referenced across skills, hooks, recipes, AGENT_PLAYBOOK, and Makefiles."
echo ""
echo "Legend:"
echo "  skill:N   = Referenced in N skill SKILL.md files"
echo "  recipe:N  = Referenced in N recipe markdown files"
echo "  hook:N    = Referenced in hooks.json"
echo "  playbook:N = Referenced in AGENT_PLAYBOOK.md"
echo "  makefile:N = Referenced in Makefile or Makefile.ctx"
echo ""
echo "==== TOP 20 MOST-REFERENCED COMMANDS ===="
echo ""
head -20 /tmp/final_counts.txt | awk -F'\t' '{
  cmd=$1; refs=$2
  printf "%-15s %3d refs", cmd, refs
  for(i=3; i<=NF; i++) printf "   %s", $i
  printf "\n"
}'

echo ""
echo "==== COMMANDS WITH ZERO REFERENCES (POTENTIALLY ORPHANED) ===="
echo ""
echo "These commands are implemented but not mentioned in any skill, hook, recipe, or playbook:"
echo ""

# Check which ones are actually commands (have a SKILL.md or are in cmd structure)
echo "  - archive"
echo "  - brainstorm"
echo "  - check-links"
echo "  - commit"
echo "  - consolidate"
echo "  - explain"
echo "  - import-plans"
echo "  - reflect"

echo ""
echo "Note: Some of these may be internal/auto-invoked commands. Check each skill directory."

echo ""
echo "==== BREAKDOWN BY SOURCE TYPE ===="
echo ""
echo "Skills reference (unique commands):"
/tmp/count_cmds.sh | awk -F'|' '$2 == "skill"' | awk -F'|' '{print $1}' | sort -u | wc -l
echo ""
echo "Recipes reference (unique commands):"
/tmp/count_cmds.sh | awk -F'|' '$2 == "recipe"' | awk -F'|' '{print $1}' | sort -u | wc -l
echo ""
echo "AGENT_PLAYBOOK references (unique commands):"
/tmp/count_cmds.sh | awk -F'|' '$2 == "playbook"' | awk -F'|' '{print $1}' | sort -u | wc -l
echo ""
echo "Hooks reference (unique commands):"
/tmp/count_cmds.sh | awk -F'|' '$2 == "hook"' | awk -F'|' '{print $1}' | sort -u | wc -l
echo ""
echo "Makefiles reference (unique commands):"
/tmp/count_cmds.sh | awk -F'|' '$2 == "makefile"' | awk -F'|' '{print $1}' | sort -u | wc -l

