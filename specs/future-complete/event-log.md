# Event Log and Doctor

## Problem

ctx hooks emit structured events (context loads, QA gates, nudges,
commits) that currently vanish unless a webhook is configured. Even
with webhooks, the events land in an external sink (Google Sheets, Slack)
that ctx can't query. There's no local, queryable record of what
hooks fired, when, or how often.

This matters for three reasons:

1. **Diagnostics** — "why didn't my hook fire?" or "why did the agent
   ignore the QA gate?" is unanswerable without scrolling Sheets.
2. **Patterns** — "qa-reminder fired 4 times in 7 minutes" reveals the
   agent was editing rapidly without committing. You'd never spot this
   without aggregation.
3. **Session forensics** — "what happened in session X?" currently
   requires correlating Claude transcripts with external logs. A local
   event log makes it self-contained.

## Approach

Two complementary features:

1. **Event log** — append-only JSONL file written by hooks, queryable
   via `ctx system events`. The dumb pipe: filter, output, done.

2. **Doctor** — a thin CLI command (`ctx doctor`) for mechanical health
   checks, paired with a skill (`/ctx-doctor`) that adds semantic
   analysis of event patterns. The CLI does what Go is good at
   (structural checks); the LLM does what it's good at (pattern
   recognition, correlation, natural language explanation). No Grafana
   reinvention.

### Design principles

- **Append-only JSONL** — one JSON object per line, no parsing
  ambiguity, `tail -f` friendly.
- **Opt-in via `.ctxrc`** — `event_log: true`. Default off to avoid
  surprises for existing users.
- **No aggregation engine** — this is `grep` and `jq` territory, not
  ELK. The CLI command provides convenience filters, not analytics.
- **Same payload as webhooks** — reuse `notify.Payload` so the local
  log and webhook carry identical data. One struct, two sinks.
- **Rotation by size** — when the log exceeds a threshold (default
  1MB), rotate to `events.1.jsonl`. Keep one rotated file. Total cap
  ~2MB.
- **Semantic analysis belongs in the skill** — programmatic aggregation
  is a treadmill toward dashboards. The LLM reads the raw events and
  reasons about patterns. Go code filters; the agent interprets.

### Where it fits

```
hook fires
  → notify.Send()        → webhook (existing, unchanged)
  → eventlog.Append()    → .context/state/events.jsonl (new)
```

The append happens in the system hook functions (alongside the
existing `notify.Send` calls), not inside `notify.Send` itself.
This keeps the notify package focused on webhooks and avoids coupling.

## Behavior

### Happy path

1. User sets `event_log: true` in `.ctxrc`.
2. System hooks call `eventlog.Append()` after their logic completes.
3. Events accumulate in `.context/state/events.jsonl`.
4. User runs `ctx system events` to query the raw log.
5. When the file exceeds 1MB, current file rotates to
   `events.1.jsonl` and a new `events.jsonl` starts.
6. User runs `ctx doctor` for a structural health report.
7. User asks the agent to diagnose a problem; `/ctx-doctor` runs
   `ctx doctor` + `ctx system events --json` and reasons about the
   combined output.

### Edge cases

| Case | Behavior |
|------|----------|
| `event_log: false` or absent | `Append()` is a noop |
| `.context/state/` doesn't exist | Create it (hooks may fire before `ctx init` completes) |
| Concurrent hook writes | JSONL is append-only; OS-level atomic line writes (< 4KB) prevent interleaving |
| Corrupt/truncated line | Reader skips unparseable lines with a warning |
| Log file missing | `ctx system events` prints "No events logged." |
| Rotation race | Check-and-rotate is best-effort; slightly over 1MB is acceptable |
| `ctx doctor` without event log | Runs structural checks, notes event logging is off |

### Error handling

| Error | User-facing message | Recovery |
|-------|---------------------|----------|
| Can't create state dir | Silent (non-fatal, don't break hook) | Skip logging |
| Can't append to file | Silent (non-fatal) | Skip logging |
| Can't parse line in query | Skip line, print warning to stderr | Continue |
| `--session` filter matches nothing | "No events for session <id>." | Exit 0 |

## Interface

### `ctx system events`

Raw event log query. Dumb pipe — filters and outputs, no analysis.

```
ctx system events [flags]
```

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--hook` | `-k` | string | (all) | Filter by hook name |
| `--session` | `-s` | string | (all) | Filter by session ID |
| `--event` | `-e` | string | (all) | Filter by event type (relay, nudge, error) |
| `--last` | `-n` | int | 50 | Show last N events |
| `--json` | `-j` | bool | false | Output raw JSONL (for piping to jq) |
| `--all` | `-a` | bool | false | Include rotated log file |

#### Examples

```bash
# Last 50 events (default)
ctx system events

# Events from a specific session
ctx system events --session eb1dc9cd-0163-4853-89d0-785fbfaae3a6

# Only QA reminder events
ctx system events --hook qa-reminder

# Raw JSONL for jq processing
ctx system events --json | jq '.message'

# How many context-load-gate fires today
ctx system events --hook context-load-gate --json | jq -r '.timestamp' | grep "$(date +%Y-%m-%d)" | wc -l
```

#### Default (human) output format

```
2026-02-27 22:39:31  relay  qa-reminder          QA gate reminder emitted
2026-02-27 22:41:56  relay  qa-reminder          QA gate reminder emitted
2026-02-28 00:48:18  relay  context-load-gate    injected 6 files (~9367 tokens)
```

Columns: timestamp (local TZ), event type, hook name (from detail),
message (truncated to terminal width).

### `ctx doctor`

Structural health check. Runs mechanical checks that don't require
semantic analysis. Think of it as `ctx status` + `ctx drift` +
configuration audit in one pass.

```
ctx doctor [flags]
```

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | false | Machine-readable JSON output |

#### What it checks

| Check | Category | Source |
|-------|----------|--------|
| Context initialized | Structure | `.context/` exists |
| Required files present | Structure | `config.FilesRequired` |
| Drift detected | Quality | `drift.Detect()` |
| Hook config valid | Hooks | `hooks.json` parseable, expected hooks registered |
| Event logging status | Config | `.ctxrc` `event_log` value |
| Webhook configured | Config | `.notify.enc` exists |
| Pending reminders | State | `reminders.json` count |
| Task completion ratio | State | TASKS.md pending vs completed |
| Context token size | Size | `context.EstimateTokens()` |
| Recent event log activity | Events | Last event timestamp (if logging enabled) |

#### Output format (human)

```
ctx doctor
==========

Structure
  ✓ Context initialized (.context/)
  ✓ Required files present (4/4)

Quality
  ⚠ Drift: 2 warnings (stale path in ARCHITECTURE.md, high entry count in LEARNINGS.md)

Hooks
  ✓ hooks.json valid (14 hooks registered)
  ○ Event logging disabled (enable with event_log: true in .ctxrc)

State
  ✓ No pending reminders
  ⚠ Task completion ratio high (18/22 = 82%) — consider archiving

Summary: 2 warnings, 0 errors
```

#### What it does NOT do

- **No event pattern analysis** — that's the skill's job
- **No auto-fixing** — reports findings, doesn't modify anything
- **No external service checks** — doesn't verify webhook endpoints

### Skill: `/ctx-doctor`

The agent-side counterpart. Runs `ctx doctor` for the structural
baseline, then reads event log data and *reasons* about it.

#### Trigger phrases

| User says | Action |
|-----------|--------|
| "why didn't my hook fire?" | Check event log for the hook, check hooks.json |
| "hooks seem broken" | Run doctor + event log query |
| "context seems stale" | Doctor + check last context-load-gate event |
| "too many nudges" | Query event log, identify noisy hooks |
| "what happened last session?" | Events by session ID + `ctx recall` |
| "doctor" / "diagnose" / "troubleshoot" / "health check" | Full triage |

#### Diagnostic playbook (in SKILL.md)

The skill instructs the agent to follow a triage sequence:

```
1. Baseline: run `ctx doctor --json` for structural health
2. Events: run `ctx system events --json --last 100` (if event logging enabled)
3. Correlate: connect findings across both sources
4. Present: structured findings with evidence
5. Suggest: actionable next steps (but don't auto-fix)
```

#### Available data sources

| Source | Command | What it reveals |
|--------|---------|-----------------|
| Structural health | `ctx doctor --json` | All mechanical checks in one call |
| Event log | `ctx system events --json --last 100` | Recent hook activity, gaps, errors |
| Event log (filtered) | `ctx system events --hook <name> --json` | Specific hook behavior |
| Reminders | `ctx remind list` | Pending reminders |
| Hook messages | `ctx system message list` | Custom vs default templates |
| RC config | Read `.ctxrc` | Configuration that affects behavior |

#### Example output from skill

```
## Doctor Report

### Structural health

`ctx doctor` reports 2 warnings, 0 errors:
- Stale path in ARCHITECTURE.md (line 42 references `internal/handler/` which doesn't exist)
- LEARNINGS.md has 35 entries (threshold: 30) — consider consolidating

### Event analysis

Looking at the last 100 events:

1. **qa-reminder fired 12 times in the last hour** — the agent was
   making many small edits without committing. This is normal during
   a refactoring session but worth noting.

2. **context-load-gate hasn't fired since 2026-02-27** — this is
   unexpected. The hook should fire on every PreToolUse. Possible
   causes: hook not registered, or the session hasn't triggered a
   tool use yet.

3. **check-ceremonies fired once, check-persistence never** — the
   persistence nudge may be silenced by a custom message template.

### Suggested actions

- [ ] Fix stale path in ARCHITECTURE.md:42
- [ ] Run `/ctx-consolidate` on LEARNINGS.md
- [ ] Verify context-load-gate is registered: `grep context-load-gate .claude/settings.local.json`
```

#### Graceful degradation

If `event_log` is not enabled, the skill still works but with reduced
capability. It runs `ctx doctor` for structural checks and notes:
"Enable `event_log: true` in `.ctxrc` for hook-level diagnostics."

## Storage

### File: `.context/state/events.jsonl`

One JSON object per line, identical to `notify.Payload`:

```json
{"event":"relay","message":"qa-reminder: QA gate reminder emitted","detail":{"hook":"qa-reminder","variant":"gate"},"session_id":"eb1dc9cd-...","timestamp":"2026-02-27T22:39:31Z","project":"ctx"}
{"event":"relay","message":"context-load-gate: injected 6 files (~9367 tokens)","session_id":"6e011357-...","timestamp":"2026-02-28T00:48:18Z","project":"ctx"}
```

### Rotation

When `events.jsonl` exceeds `EventLogMaxBytes` (default 1MB):

1. Remove `events.1.jsonl` if it exists.
2. Rename `events.jsonl` → `events.1.jsonl`.
3. Create new empty `events.jsonl`.

Total disk budget: ~2MB. For structured log lines averaging ~250
bytes, that's ~8000 events — weeks of heavy usage.

### Gitignore

Event logs are **gitignored** (added to `config.GitignoreEntries`).
They're machine-local diagnostics, not project context. Different
machines produce different event streams.

## Configuration

### `.ctxrc`

```yaml
event_log: true    # default: false
```

Single boolean. No per-hook filtering (use `ctx system events --hook`
at query time). No custom path (always `.context/state/`).

### Constants (`internal/config`)

```go
const (
    FileEventLog     = "events.jsonl"
    FileEventLogPrev = "events.1.jsonl"
    EventLogMaxBytes = 1 << 20  // 1MB
)
```

## Implementation

### New files

```
internal/eventlog/
    eventlog.go       Append(), rotate(), Query()
    eventlog_test.go  Tests

internal/cli/system/
    events.go         eventsCmd(), runEvents()
    events_test.go    Tests

internal/cli/doctor/
    doctor.go         Cmd(), runDoctor(), individual check functions
    doctor_test.go    Tests

internal/assets/claude/skills/ctx-doctor/
    SKILL.md          Diagnostic skill instructions
```

### Modified files

| File | Change |
|------|--------|
| `internal/rc/types.go` | Add `EventLog bool` field |
| `internal/rc/rc.go` | Add `EventLog()` accessor |
| `internal/rc/default.go` | Default `false` |
| `internal/config/dir.go` | Add `FileEventLog`, `FileEventLogPrev`, `EventLogMaxBytes` constants |
| `internal/config/gitignore.go` | Add `state/events.jsonl` and `state/events.1.jsonl` to `GitignoreEntries` |
| `internal/cli/system/system.go` | Register `eventsCmd()` |
| `internal/bootstrap/bootstrap.go` | Register `doctor.Cmd()` |
| System hook files (`check_ceremonies.go`, `check_persistence.go`, `check_context_size.go`, `check_journal.go`, `check_reminders.go`, `check_knowledge.go`, `check_map_staleness.go`, `check_version.go`, `check_resources.go`, `context_load_gate.go`, `post_commit.go`, `qa_reminder.go`) | Add `eventlog.Append()` call alongside `notify.Send()` |

### Key types and functions

```go
// internal/eventlog/eventlog.go

// Append writes a single event to the log file.
// Noop if event logging is disabled in .ctxrc.
//
// Parameters:
//   - event: Event type (e.g., "relay", "nudge")
//   - message: Human-readable description
//   - sessionID: Claude session ID (may be empty)
//   - detail: Optional template reference (may be nil)
func Append(event, message, sessionID string, detail *notify.TemplateRef) {
    if !rc.EventLog() {
        return
    }
    // ... build payload, marshal, append, check rotation
}

// Query reads events from the log, applying filters.
//
// Parameters:
//   - opts: Filter and limit options
//
// Returns:
//   - []notify.Payload: Matching events (newest last)
//   - error: Non-nil if file read fails
func Query(opts QueryOpts) ([]notify.Payload, error) { ... }

// QueryOpts controls event filtering and pagination.
type QueryOpts struct {
    Hook      string // filter by hook name (from detail)
    Session   string // filter by session ID
    Event     string // filter by event type
    Last      int    // return last N events (0 = all)
    IncludeRotated bool // also read events.1.jsonl
}
```

```go
// internal/cli/doctor/doctor.go

// Result represents a single check outcome.
type Result struct {
    Name     string `json:"name"`
    Category string `json:"category"`
    Status   string `json:"status"` // "ok", "warning", "error"
    Message  string `json:"message"`
}

// Report is the complete doctor output.
type Report struct {
    Results  []Result `json:"results"`
    Warnings int      `json:"warnings"`
    Errors   int      `json:"errors"`
}

// Cmd returns the "ctx doctor" command.
func Cmd() *cobra.Command { ... }

// runDoctor executes all checks and prints the report.
func runDoctor(cmd *cobra.Command, jsonOutput bool) error { ... }
```

### Helpers to reuse

- `notify.Payload` — the event struct (already JSON-tagged)
- `notify.TemplateRef` — the detail struct
- `rc.ContextDir()` — base path for state dir
- `config.DirState` — "state" constant
- `drift.Detect()` — existing drift detection
- `context.EstimateTokens()` — token counting
- `filepath.Join()` — path construction

### Call site pattern

Each system hook that currently calls `notify.Send()` gets a
parallel `eventlog.Append()` call with the same arguments:

```go
// existing
notify.Send("relay", rendered, sessionID, detail)

// added
eventlog.Append("relay", rendered, sessionID, detail)
```

The `Append` signature mirrors `Send` intentionally so the call
sites stay mechanical.

## Testing

### `internal/eventlog/eventlog_test.go`

| Test | Scenario |
|------|----------|
| `TestAppend_Disabled` | `event_log: false` → file not created |
| `TestAppend_Basic` | Append event, read back, verify fields |
| `TestAppend_CreatesStateDir` | State dir missing → created automatically |
| `TestAppend_Rotation` | Write > 1MB → rotated, new file started |
| `TestAppend_RotationOverwrite` | Rotation overwrites existing `.1` file |
| `TestQuery_NoFile` | No log file → empty result, no error |
| `TestQuery_FilterHook` | Filter by hook name → only matching events |
| `TestQuery_FilterSession` | Filter by session ID → only matching events |
| `TestQuery_Last` | `--last 5` on 20 events → last 5 returned |
| `TestQuery_IncludeRotated` | `--all` reads both current and rotated |
| `TestQuery_CorruptLine` | Malformed JSON line → skipped with warning |

### `internal/cli/system/events_test.go`

| Test | Scenario |
|------|----------|
| `TestEventsCmd_Default` | Human-readable output, 50-line default |
| `TestEventsCmd_JSON` | `--json` flag → raw JSONL output |
| `TestEventsCmd_NoLog` | No file → "No events logged." |
| `TestEventsCmd_Filters` | Combined `--hook` + `--session` → intersection |

### `internal/cli/doctor/doctor_test.go`

| Test | Scenario |
|------|----------|
| `TestDoctor_Healthy` | All checks pass → summary shows 0 warnings, 0 errors |
| `TestDoctor_NoContext` | No `.context/` → structure error reported |
| `TestDoctor_DriftWarnings` | Drift detected → warnings surfaced |
| `TestDoctor_EventLogOff` | Event logging disabled → info note, not an error |
| `TestDoctor_JSON` | `--json` → valid JSON report |
| `TestDoctor_HighCompletion` | High task completion ratio → warning |

## Documentation

### 1. CLI reference (`docs/cli/system.md`)

Add `ctx system events` section after `ctx system message`. Full
entry with flags table, examples, and human/JSON output format.

### 2. New CLI page (`docs/cli/doctor.md`)

`ctx doctor` is a top-level command (not under `system`), so it gets
its own page. Include:

- Command syntax and `--json` flag
- Check table (what it checks, categories)
- Human and JSON output examples
- "When to use" guidance: `ctx status` for a quick glance,
  `ctx doctor` for a thorough checkup, `/ctx-doctor` when you need
  the agent to reason about what's wrong

### 3. CLI index (`docs/cli/index.md`)

Add row to the commands table:

```markdown
| [`ctx doctor`](doctor.md#ctx-doctor) | Structural health check (hooks, drift, config) |
```

### 4. Configuration docs (`docs/home/configuration.md` or equivalent `.ctxrc` section)

Add `event_log` to the `.ctxrc` reference table:

```markdown
| `event_log` | `bool` | `false` | Enable local hook event logging to `.context/state/events.jsonl` |
```

### 5. Skills reference (`docs/reference/skills.md`)

Add `/ctx-doctor` entry:

```markdown
| `/ctx-doctor` | Troubleshoot ctx behavior. Runs structural health checks, analyzes event log patterns, and presents findings with suggested actions. |
```

Trigger phrases: "diagnose", "troubleshoot", "doctor", "health check",
"why didn't my hook fire?", "hooks seem broken"

### 6. Existing recipe updates

**[Auditing System Hooks](docs/recipes/system-hooks-audit.md):**
Add a section on event logging as a complement to webhook-based auditing.
Mention `ctx system events` as the local alternative to checking Sheets.
Cross-link to the troubleshooting recipe.

**[Detecting and Fixing Drift](docs/recipes/context-health.md):**
Add a note that `ctx doctor` now provides a superset of drift checks
combined with hook and config auditing. Cross-link to the
troubleshooting recipe.

**[Webhook Notifications](docs/recipes/webhook-notifications.md):**
Add a note that event logging provides a local complement to webhooks.
"Don't need a webhook but want diagnostic visibility? Enable
`event_log: true` in `.ctxrc`."

### 7. New recipe: Troubleshooting (`docs/recipes/troubleshooting.md`)

New recipe covering the full diagnostic workflow. Structure:

```markdown
---
title: "Troubleshooting"
icon: lucide/stethoscope
---

## The Problem

Something isn't working: a hook isn't firing, nudges are too noisy,
context seems stale, or the agent isn't following instructions. The
information to diagnose it exists — across status, drift, event logs,
hook config, and session history — but assembling it manually is tedious.

## TL;DR

    ctx doctor                   # structural health check
    ctx system events --last 20  # recent hook activity
    # or ask: "something seems off, can you diagnose?"

## Commands and Skills Used

| Tool                       | Type        | Purpose                              |
|----------------------------|-------------|--------------------------------------|
| `ctx doctor`               | CLI command | Structural health report             |
| `ctx doctor --json`        | CLI command | Machine-readable health report       |
| `ctx system events`        | CLI command | Query local event log                |
| `/ctx-doctor`              | Skill       | Agent-driven diagnosis with analysis |

## The Workflow

### Quick check: `ctx doctor`

Run `ctx doctor` for an instant structural health report...

### Deep dive: `/ctx-doctor`

Ask the agent to diagnose...

### Raw event inspection

For power users: `ctx system events` with filters...

## Common Problems

### "My hook isn't firing"
### "Too many nudges"
### "Context seems stale"
### "The agent isn't following instructions"

## Prerequisites

- Event logging: `event_log: true` in `.ctxrc` (optional but recommended)
- ctx initialized: `ctx init`

## See Also

- [Auditing System Hooks](system-hooks-audit.md)
- [Detecting and Fixing Drift](context-health.md)
- [Webhook Notifications](webhook-notifications.md)
```

### 8. Recipes index (`docs/recipes/index.md`)

Add to the **Maintenance** section:

```markdown
### [Troubleshooting](troubleshooting.md)

Diagnose hook failures, noisy nudges, stale context, and configuration
issues. Start with `ctx doctor` for a structural health check, then
use `/ctx-doctor` for agent-driven analysis of event patterns.

**Uses**: `ctx doctor`, `ctx system events`, `/ctx-doctor`
```

### 9. `zensical.toml`

Add nav entries for new pages:

- `docs/cli/doctor.md` in the CLI nav section
- `docs/recipes/troubleshooting.md` in the Recipes nav section

## Non-goals

- **Real-time streaming** — no `tail -f` mode. Use `tail -f` directly
  on the JSONL file if needed.
- **Aggregation or dashboards** — no counters, no charts, no trends.
  Pipe to `jq` or your tool of choice. Semantic analysis lives in
  the skill, not in Go code.
- **Custom log paths** — always `.context/state/events.jsonl`. One
  knob (`event_log: true/false`), not two.
- **Per-hook log filtering at write time** — all hooks log when
  enabled. Filter at query time.
- **Structured query language** — no SQL, no DSL. Flags + jq.
- **Log shipping** — webhooks already handle that. The local log is
  for local queries.
- **Backward-filling from Sheets** — the external sink is external.
  The local log starts when you enable it.
- **Auto-fixing in `ctx doctor`** — reports findings, never modifies.
  `ctx drift --fix` already handles auto-fixable drift; doctor is
  read-only by design.
- **Programmatic event pattern analysis** — no "alert if X fires
  more than N times." That's the skill's territory. Go code stays
  dumb; the LLM stays smart.

## Open questions

- **Should `ctx init` prompt to enable event logging?** Leaning no —
  it's a power-user feature. Document it in the recipe and let users
  opt in.
