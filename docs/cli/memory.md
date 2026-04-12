---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Memory
icon: lucide/brain
---

![ctx](../images/ctx-banner.png)

## `ctx memory`

Bridge Claude Code's auto memory (MEMORY.md) into `.context/`.

Claude Code maintains per-project auto memory at
`~/.claude/projects/<slug>/memory/MEMORY.md`. This command group
discovers that file, mirrors it into `.context/memory/mirror.md`
(git-tracked), and detects drift.

```bash
ctx memory <subcommand>
```

### `ctx memory sync`

Copy MEMORY.md to `.context/memory/mirror.md`. Archives the
previous mirror before overwriting.

```bash
ctx memory sync [flags]
```

**Flags**:

| Flag        | Description                              |
|-------------|------------------------------------------|
| `--dry-run` | Show what would happen without writing   |

**Exit codes**:

| Code | Meaning                                    |
|------|--------------------------------------------|
| 0    | Synced successfully                        |
| 1    | MEMORY.md not found (auto memory inactive) |

**Examples**:

```bash
ctx memory sync
# Archived previous mirror to mirror-2026-03-05-143022.md
# Synced MEMORY.md -> .context/memory/mirror.md
#   Source: ~/.claude/projects/-home-user-project/memory/MEMORY.md
#   Lines: 47 (was 32)
#   New content: 15 lines since last sync

ctx memory sync --dry-run
```

### `ctx memory status`

Show drift, timestamps, line counts, and archive count.

```bash
ctx memory status
```

**Exit codes**:

| Code | Meaning                                       |
|------|-----------------------------------------------|
| 0    | No drift                                      |
| 1    | MEMORY.md not found                           |
| 2    | Drift detected (MEMORY.md changed since sync) |

**Examples**:

```bash
ctx memory status
# Memory Bridge Status
#   Source:      ~/.claude/projects/.../memory/MEMORY.md
#   Mirror:      .context/memory/mirror.md
#   Last sync:   2026-03-05 14:30 (2 hours ago)
#
#   MEMORY.md:  47 lines (modified since last sync)
#   Mirror:     32 lines
#   Drift:      detected (source is newer)
#   Archives:   3 snapshots in .context/memory/archive/
```

### `ctx memory diff`

Show what changed in MEMORY.md since last sync.

```bash
ctx memory diff
```

**Examples**:

```bash
ctx memory diff
# --- .context/memory/mirror.md (mirror)
# +++ ~/.claude/projects/.../memory/MEMORY.md (source)
# +- new learning: memory bridge works
```

No output when files are identical.

### `ctx memory publish`

Push curated `.context/` content into MEMORY.md so the agent
sees it natively.

```bash
ctx memory publish [flags]
```

Content is selected in priority order: pending tasks, recent
decisions (7 days), key conventions, recent learnings (7 days).
Wrapped in `<!-- ctx:published -->` markers. Claude-owned
content outside the markers is preserved.

**Flags**:

| Flag        | Description                              | Default |
|-------------|------------------------------------------|---------|
| `--budget`  | Line budget for published content        | `80`    |
| `--dry-run` | Show what would be published             |         |

**Examples**:

```bash
ctx memory publish --dry-run
# Publishing .context/ -> MEMORY.md...
#   Budget: 80 lines
#   Published block:
#     5 pending tasks (from TASKS.md)
#     3 recent decisions (from DECISIONS.md)
#     5 key conventions (from CONVENTIONS.md)
#   Total: 42 lines (within 80-line budget)
# Dry run - no files written.

ctx memory publish              # Write to MEMORY.md
ctx memory publish --budget 40  # Tighter budget
```

### `ctx memory unpublish`

Remove the ctx-managed marker block from MEMORY.md, preserving
Claude-owned content.

**Examples**:

```bash
ctx memory unpublish
```

**Hook integration**: The `check-memory-drift` hook runs on
every prompt and nudges the agent when MEMORY.md has changed
since last sync. The nudge fires once per session. See
[Memory Bridge](../recipes/memory-bridge.md).

### `ctx memory import`

Classify and promote entries from MEMORY.md into structured
`.context/` files.

```bash
ctx memory import [flags]
```

Each entry is classified by keyword heuristics:

| Keywords                                          | Target         |
|---------------------------------------------------|----------------|
| `always use`, `prefer`, `never use`, `standard`   | CONVENTIONS.md |
| `decided`, `chose`, `trade-off`, `approach`       | DECISIONS.md   |
| `gotcha`, `learned`, `watch out`, `bug`, `caveat` | LEARNINGS.md   |
| `todo`, `need to`, `follow up`                    | TASKS.md       |
| Everything else                                   | Skipped        |

Deduplication prevents re-importing the same entry across runs.

**Flags**:

| Flag        | Description                                    |
|-------------|------------------------------------------------|
| `--dry-run` | Show classification plan without writing       |

**Examples**:

```bash
ctx memory import --dry-run
# Scanning MEMORY.md for new entries...
#   Found 6 entries
#
#   -> "always use ctx from PATH"
#      Classified: CONVENTIONS.md (keywords: always use)
#
#   -> "decided to use heuristic classification over LLM-based"
#      Classified: DECISIONS.md (keywords: decided)
#
# Dry run - would import: 4 entries
# Skipped: 2 entries (session notes/unclassified)

ctx memory import    # Actually write entries to .context/ files
```
