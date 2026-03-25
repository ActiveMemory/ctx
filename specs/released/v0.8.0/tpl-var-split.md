---
title: Split config/tpl into domain config packages
date: 2026-03-22
status: ready
---

# Split config/tpl Template Variables

## Problem

`config/tpl/tpl.go` bags 24 template variable constants from 10
different domains into one package. Each variable should live in
its domain config package.

## Plan

| Variable(s) | Target package | File |
|-------------|---------------|------|
| VarWarnings | config/archive | var.go |
| VarStaleFiles | config/freshness | var.go |
| VarUnexportedCount, VarUnenrichedCount | config/journal | var.go |
| VarPromptCount, VarPromptsSinceNudge | config/nudge | var.go |
| VarReminderList | config/reminder | var.go |
| VarAlertMessages, VarPercentage, VarTokenCount, VarThreshold | config/stats | var.go |
| VarBinaryVersion, VarPluginVersion, VarKeyAgeDays | config/version | var.go |
| VarHeartbeat* (6) | config/heartbeat | var.go |
| VarFileWarnings | config/knowledge | var.go |
| VarLastRefreshDate, VarModuleCount | config/architecture | var.go |

## Execution

1. Add var.go to each target config package with the constants
2. Delete config/tpl/tpl.go and config/tpl/doc.go
3. Let it break
4. Fix each callsite to import the domain package instead of tpl
