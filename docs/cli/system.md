---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: System
icon: lucide/settings
---

### `ctx system`

System diagnostics and hook commands.

```bash
ctx system <subcommand>
```

The parent command shows available subcommands. Hidden plumbing subcommands
(`ctx system mark-journal`) are used by skills and automation. Hidden hook
subcommands (`ctx system check-*`) are used by the Claude Code plugin — see
[AI Tools](../operations/integrations.md#plugin-hooks) for details.

#### `ctx system bootstrap`

Print context location and rules for AI agents. This is the recommended first
command for AI agents to run at session start — it tells them where the context
directory is and how to use it.

```bash
ctx system bootstrap [flags]
```

**Flags**:

| Flag     | Description           |
|----------|-----------------------|
| `--json` | Output in JSON format |

**Text output**:

```
ctx bootstrap
=============

context_dir: .context

Files:
  CONSTITUTION.md, TASKS.md, DECISIONS.md, LEARNINGS.md,
  CONVENTIONS.md, ARCHITECTURE.md, GLOSSARY.md

Rules:
  1. Use context_dir above for ALL file reads/writes
  2. Never say "I don't have memory" — context IS your memory
  3. Read files silently, present as recall (not search)
  4. Persist learnings/decisions before session ends
  5. Run `ctx agent` for content summaries
  6. Run `ctx status` for context health
```

**JSON output**:

```json
{
  "context_dir": ".context",
  "files": ["CONSTITUTION.md", "TASKS.md", ...],
  "rules": ["Use context_dir above for ALL file reads/writes", ...]
}
```

**Examples**:

```bash
ctx system bootstrap                          # Text output
ctx system bootstrap --json                   # JSON output
ctx system bootstrap --json | jq .context_dir # Extract context path
```

**Why it exists**: When users configure an external context directory via
`.ctxrc` (`context_dir: /mnt/nas/.context`), the AI agent needs to know where
context lives. Bootstrap resolves the configured path and communicates it to
the agent at session start. Every nudge also includes a context directory
footer for reinforcement.

#### `ctx system resources`

Show system resource usage with threshold-based alerts.

```bash
ctx system resources [flags]
```

Displays memory, swap, disk, and CPU load with two severity tiers:

| Resource | WARNING | DANGER |
|----------|---------|--------|
| Memory | >= 80% used | >= 90% used |
| Swap | >= 50% used | >= 75% used |
| Disk (cwd) | >= 85% full | >= 95% full |
| Load (1m) | >= 0.8x CPUs | >= 1.5x CPUs |

**Flags**:

| Flag     | Description           |
|----------|-----------------------|
| `--json` | Output in JSON format |

**Examples**:

```bash
ctx system resources               # Text output with status indicators
ctx system resources --json        # Machine-readable JSON
ctx system resources --json | jq '.alerts'   # Extract alerts only
```

**Text output**:

```
System Resources
====================

Memory:    4.2 / 16.0 GB (26%)                     ✓ ok
Swap:      0.0 /  8.0 GB (0%)                      ✓ ok
Disk:    180.2 / 500.0 GB (36%)                     ✓ ok
Load:     0.52 / 0.41 / 0.38  (8 CPUs, ratio 0.07) ✓ ok

All clear — no resource warnings.
```

When resources breach thresholds, alerts are listed below the summary:

```
Alerts:
  ✖ Memory 92% used (14.7 / 16.0 GB)
  ✖ Swap 78% used (6.2 / 8.0 GB)
  ✖ Load 1.56x CPU count
```

**Platform support**: Full metrics on Linux and macOS. Windows shows
disk only; memory and load report as unsupported.

#### `ctx system message`

Manage hook message templates. Hook messages control what text hooks emit.
The hook logic (when to fire, counting, state tracking) is universal; the
messages are opinions that can be customized per-project.

```bash
ctx system message <subcommand>
```

**Subcommands**:

| Subcommand | Args | Flags | Description |
|------------|------|-------|-------------|
| `list` | *(none)* | `--json` | Show all hook messages with category and override status |
| `show` | `<hook> <variant>` | *(none)* | Print the effective message template with source |
| `edit` | `<hook> <variant>` | *(none)* | Copy embedded default to `.context/` for editing |
| `reset` | `<hook> <variant>` | *(none)* | Delete user override, revert to embedded default |

**Examples**:

```bash
ctx system message list                      # Table of all 24 messages
ctx system message list --json               # Machine-readable JSON
ctx system message show qa-reminder gate     # View the QA gate template
ctx system message edit qa-reminder gate     # Copy default to .context/ for editing
ctx system message reset qa-reminder gate    # Delete override, revert to default
```

Override files are placed at `.context/hooks/messages/{hook}/{variant}.txt`.
An empty override file silences the message while preserving the hook's logic.

See the [Customizing Hook Messages](../recipes/customizing-hook-messages.md)
recipe for detailed examples.

#### `ctx system events`

Query the local hook event log. Reads events from `.context/state/events.jsonl`
and outputs them in human-readable or raw JSONL format. Requires `event_log: true`
in `.ctxrc`.

```bash
ctx system events [flags]
```

**Flags**:

| Flag        | Short | Type   | Default | Description                                  |
|-------------|-------|--------|---------|----------------------------------------------|
| `--hook`    | `-k`  | string | *(all)* | Filter by hook name                          |
| `--session` | `-s`  | string | *(all)* | Filter by session ID                         |
| `--event`   | `-e`  | string | *(all)* | Filter by event type (relay, nudge)          |
| `--last`    | `-n`  | int    | `50`    | Show last N events                           |
| `--json`    | `-j`  | bool   | `false` | Output raw JSONL (for piping to `jq`)        |
| `--all`     | `-a`  | bool   | `false` | Include rotated log file (`events.1.jsonl`)  |

All filter flags use intersection (AND) logic.

**Text output**:

```
2026-02-27 22:39:31  relay  qa-reminder          QA gate reminder emitted
2026-02-27 22:41:56  relay  qa-reminder          QA gate reminder emitted
2026-02-28 00:48:18  relay  context-load-gate    injected 6 files (~9367 tokens)
```

Columns: timestamp (*local timezone*), event type, hook name, message (*truncated
to terminal width*).

**JSON output** (`--json`):

Each line is a standalone JSON object identical to the webhook payload format:

```json
{"event":"relay","message":"qa-reminder: QA gate reminder emitted","detail":{"hook":"qa-reminder","variant":"gate"},"session_id":"eb1dc9cd-...","timestamp":"2026-02-27T22:39:31Z","project":"ctx"}
```

**Examples**:

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
ctx system events --hook context-load-gate --json \
  | jq -r '.timestamp' | grep "$(date +%Y-%m-%d)" | wc -l

# Include rotated events
ctx system events --all --last 100
```

**Why it exists**: System hooks fire invisibly. When something goes wrong (*"why
didn't my hook fire?"*), the event log provides a local, queryable record of
what hooks fired, when, and how often. Event logging is opt-in via
`event_log: true` in `.ctxrc` to avoid surprises for existing users.

**See also**: [Troubleshooting](../recipes/troubleshooting.md),
[Auditing System Hooks](../recipes/system-hooks-audit.md),
[`ctx doctor`](doctor.md#ctx-doctor)

---

#### `ctx system mark-journal`

Update processing state for a journal entry. Records the current date
in `.context/journal/.state.json`. Used by journal skills to record
pipeline progress.

```bash
ctx system mark-journal <filename> <stage>
```

**Stages**: `exported`, `enriched`, `normalized`, `fences_verified`

| Flag      | Description                           |
|-----------|---------------------------------------|
| `--check` | Check if stage is set (exit 1 if not) |

**Example**:

```bash
ctx system mark-journal 2026-01-21-session-abc12345.md enriched
ctx system mark-journal 2026-01-21-session-abc12345.md normalized
ctx system mark-journal --check 2026-01-21-session-abc12345.md fences_verified
```
