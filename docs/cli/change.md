---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Change
icon: lucide/diff
---

![ctx](../images/ctx-banner.png)

## `ctx change`

Show what changed in context files and code since your last
session.

Automatically detects the previous session boundary from state
markers or event log. Useful at session start to quickly see
what moved while you were away.

```bash
ctx change [flags]
```

**Flags**:

| Flag      | Description                                             |
|-----------|---------------------------------------------------------|
| `--since` | Time reference: duration (`24h`) or date (`2026-03-01`) |

**Reference time detection** (priority order):

1. `--since` flag (duration, date, or RFC3339 timestamp)
2. `ctx-loaded-*` marker files in `.context/state/` (second most recent)
3. Last `context-load-gate` event from `.context/state/events.jsonl`
4. Fallback: 24 hours ago

**Examples**:

```bash
# Auto-detect last session, show what changed
ctx change

# Changes in the last 48 hours
ctx change --since 48h

# Changes since a specific date
ctx change --since 2026-03-10
```

**Output**:

```
## Changes Since Last Session

**Reference point**: 6 hours ago

### Context File Changes
- `TASKS.md` - modified 2026-03-12 14:30
- `DECISIONS.md` - modified 2026-03-12 09:15

### Code Changes
- **12 commits** since reference point
- **Latest**: Fix journal enrichment ordering
- **Directories touched**: internal, docs, specs
- **Authors**: jose, claude
```

Context file changes are detected by filesystem mtime (works
without git). Code changes use `git log --since` (empty when
not in a git repo).

**See also**: [Reviewing Session Changes](../recipes/session-changes.md).
