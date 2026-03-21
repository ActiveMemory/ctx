---
title: Migrate system cmd.Print* calls to internal/write/
date: 2026-03-21
status: ready
---

# System write/ Migration

## Problem

`internal/cli/system/core/` has 34 `cmd.Print*` calls across 10 files.
`internal/cli/system/cmd/*/run.go` has 32 more across 18 files.

Per the project taxonomy decision: output functions that call
`cmd.Print*` belong in `write/`, not `core/`. `core/` should own
logic and types only.

## Scope

### Priority 1: core/ functions that accept *cobra.Command (21 functions)

These are the worst violations — `core/` should never import cobra
for output purposes. Each needs splitting: logic stays in core/,
output moves to write/.

| File | Function | Print calls |
|------|----------|-------------|
| `input.go` | `PrintHookContext` | 1 |
| `ceremony.go` | `EmitCeremonyNudge` | 1 |
| `context_size.go` | `EmitCheckpoint` | 2 |
| `context_size.go` | `EmitWindowWarning` | 2 |
| `context_size.go` | `EmitBillingWarning` | 2 |
| `events.go` | `OutputEventsJSON` | 1 |
| `events.go` | `OutputEventsHuman` | 2 |
| `knowledge.go` | `EmitKnowledgeWarning` | 1 |
| `map_staleness.go` | `EmitMapStalenessWarning` | 1 |
| `message_cmd.go` | `PrintTemplateVars` | 2 |
| `resources.go` | `OutputResourcesText` | 13 |
| `resources.go` | `OutputResourcesJSON` | 0 (uses encoder) |
| `stats.go` | `DumpStats` | 1 |
| `stats.go` | `OutputStatsJSON` | 1 |
| `stats.go` | `PrintStatsHeader` | 2 |
| `stats.go` | `PrintStatsLine` | 1 |
| `stats.go` | `StreamStats` | 1 |
| `version.go` | `CheckKeyAge` | 1 |
| `version_drift.go` | `CheckVersionDrift` | 0 (delegates) |
| `backup.go` | `addEntry` | 0 (uses tw.Writer) |
| `knowledge.go` | `CheckKnowledgeHealth` | 0 (delegates) |

### Priority 2: cmd/ run.go direct Print calls (32 calls, 18 files)

These are less critical (run.go is allowed some output) but should
still use write/ functions for consistency and localization.

| File | Calls |
|------|-------|
| `block_dangerous_command/run.go` | 1 |
| `block_non_path_ctx/run.go` | 1 |
| `bootstrap/run.go` | 1 |
| `check_backup_age/run.go` | 1 |
| `check_context_size/run.go` | 1 |
| `check_freshness/run.go` | 1 |
| `check_journal/run.go` | 1 |
| `check_memory_drift/run.go` | 1 |
| `check_persistence/run.go` | 3 |
| `check_reminder/run.go` | 1 |
| `check_version/run.go` | 1 |
| `events/run.go` | 1 |
| `mark_journal/run.go` | 2 |
| `message/cmd/*/run.go` | 12 |
| `pause/run.go` | 1 |

## Approach

1. Create `internal/write/system/` package
2. For each core/ function that does output:
   - Extract the `cmd.Print*` calls into a write/ function
   - Core function returns data; caller passes to write/ function
   - OR core function takes an `io.Writer` instead of `*cobra.Command`
3. For cmd/ run.go calls: create write/ helpers where the same
   output pattern repeats (NudgeBox, JSON response, etc.)

## Non-goals

- Moving NudgeBox itself (it returns a string, doesn't print)
- Changing the Relay/NudgeAndRelay functions (they don't print)
- Moving JSON encoder calls (they use cmd.OutOrStdout(), not Print)
