---
title: Eliminate write/system/ — split into domain write/ packages
date: 2026-03-21
status: ready
prerequisite: system-write-migration.md
---

# Eliminate write/system/

## Problem

The initial write migration created `write/system/` with generic
wrappers (`Line`, `Lines`, `Raw`, `NudgeBlock`). This is a junk drawer.

The existing write/ taxonomy is flat by domain: `write/backup/`,
`write/bootstrap/`, `write/status/` — each maps to a CLI feature, not
a CLI parent. There is no `write/system/` because "system" is a CLI
namespace, not a domain.

## Solution

Delete `write/system/` entirely. Each domain gets its own write/
package, and existing packages absorb shared patterns.

### Existing packages that absorb work

| Package | Absorbs |
|---------|---------|
| `write/hook/` | Already has `Nudge(cmd, box)`. Add `NudgeBlock`, `HookContext`, `BlockResponse`. All hooks use these. |

### New packages

| Package | Functions | Used by |
|---------|-----------|---------|
| `write/events/` | `JSON`, `Human`, `Empty` | events/run.go |
| `write/stats/` (exists: `write/status/`) | — see below | — |
| `write/resources/` | `Text` | resources/run.go |
| `write/message/` | `TemplateVars`, `CtxSpecificWarning`, `OverrideCreated`, `EditHint`, `SourceHeader`, `ContentBlock`, `NoOverride`, `OverrideRemoved`, `ListHeader`, `ListRow` | message/cmd/*/run.go |
| `write/mark_journal/` | `StageChecked`, `StageMarked` | mark_journal/run.go |
| `write/pause/` | `Confirmed` | pause/run.go |

### Functions routed to existing packages

- **`write/hook.Nudge`** — already exists, used by all nudge-box hooks
- **`write/hook.NudgeBlock`** — new: prints box + empty line (check_context_size, check_persistence)
- **`write/hook.HookContext`** — new: prints JSON hook response (5 hooks)
- **`write/hook.BlockResponse`** — new: prints JSON block response (2 hooks)
- **`write/bootstrap.Dir`** — new: quiet-mode directory output
- **`write/session.PausedMessage`** — check if write/session/ already has this; otherwise add to write/pause/

### What about stats?

`write/status/` already exists for the status command. Stats (token
usage telemetry) is a different domain. Create `write/stats/` with
`Table` function.

### What about check_* hooks?

The check_* hooks (ceremony, knowledge, map_staleness, version,
context_size, backup_age, freshness, journal, memory_drift, reminder,
persistence) all use the same output pattern: `write/hook.Nudge(cmd, box)`
or `write/hook.NudgeBlock(cmd, box)`. They don't need their own write/
packages — the hook package covers their output needs.

### What gets deleted

The entire `write/system/` directory is removed. All files:
- `doc.go`, `system.go`, `nudge.go`, `hook.go`, `events.go`,
  `stats.go`, `resources.go`, `context_size.go`, `ceremony.go`,
  `knowledge.go`, `map_staleness.go`, `version.go`, `message.go`,
  `backup.go`

### Caller migration summary

| Caller pattern | Before | After |
|----------------|--------|-------|
| Nudge box (single) | `systemwrite.Line(cmd, core.NudgeBox(...))` | `writeHook.Nudge(cmd, core.NudgeBox(...))` |
| Nudge box + blank line | `systemwrite.NudgeBlock(cmd, box)` | `writeHook.NudgeBlock(cmd, box)` |
| Hook context JSON | `systemwrite.Line(cmd, core.FormatHookContext(...))` | `writeHook.HookContext(cmd, core.FormatHookContext(...))` |
| Block response JSON | `systemwrite.Line(cmd, string(data))` | `writeHook.BlockResponse(cmd, string(data))` |
| Events output | `systemwrite.Lines(cmd, lines)` | `writeEvents.JSON(cmd, lines)` or `writeEvents.Human(cmd, lines)` |
| Stats output | `systemwrite.Lines(cmd, lines)` | `writeStats.Table(cmd, lines)` |
| Resources output | `systemwrite.Lines(cmd, lines)` | `writeResources.Text(cmd, lines)` |
| Message output | `systemwrite.Line(cmd, fmt.Sprintf(...))` | `writeMessage.OverrideCreated(cmd, path)` etc. |

## Non-goals

- Moving core/ format helpers
- Changing function signatures from the prior migration
