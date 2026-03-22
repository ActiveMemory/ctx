---
title: Migrate remaining cmd.Print* to write/ packages
date: 2026-03-22
status: ready
prerequisite: system-write-migration.md, write-system-taxonomy.md
---

# Remaining cmd.Print* Migration

## Problem

86 cmd.Print* calls remain in internal/cli/ outside the system/
subtree (which was already migrated). These should route through
domain-specific write/ packages.

## Scope

### New write/ packages needed

| Package | For |
|---------|-----|
| write/agent/ | agent packet rendering |
| write/change/ | changes output |
| write/config/ | config profile status, schema |
| write/doctor/ | doctor report output |
| write/drift/ | drift report, fix output |

### Existing write/ packages to extend

| Package | Functions to add |
|---------|-----------------|
| write/compact/ | Report (heading, entries, summary) |
| write/deps/ | Mermaid, Table, JSON output |
| write/guide/ | CommandList, DefaultGuide |
| write/hook/ | ToolOutput, Separator |
| write/initialize/ | ConfirmPrompt |
| write/memory/ | DiffOutput, StatusSeparator |
| write/pad/ | ShowEntry, ShowBlob, ListEntry |
| write/prompt/ | ShowContent |
| write/recall/ | ConfirmPrompt |
| write/watch/ | Separator |
| write/why/ | Content, Separator |

## Approach

For each call site:
1. Create or extend the write/ function with proper nil guard
2. Use desc.Text() for any remaining inline strings
3. Update the caller to use the write/ function
4. Remove cobra import from caller if no longer needed
