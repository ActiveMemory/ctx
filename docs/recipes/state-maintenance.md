---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "State Directory Maintenance"
icon: lucide/folder-cog
---

![ctx](../images/ctx-banner.png)

## The Problem

Every session creates tombstone files in `.context/state/` - small markers
that suppress repeat hook nudges ("already checked context size", "already
sent persistence reminder"). Over days and weeks, these accumulate into
hundreds of files from long-dead sessions.

The files are harmless individually, but the clutter makes it harder to
reason about state, and stale global tombstones can suppress nudges across
sessions entirely.

## TL;DR

```bash
ctx prune --dry-run     # preview what would be removed
ctx prune               # prune files older than 7 days
ctx prune --days 1      # more aggressive: keep only today
```

## Commands Used

| Tool               | Type    | Purpose                                   |
|--------------------|---------|-------------------------------------------|
| `ctx prune` | Command | Remove old per-session state files        |
| `ctx status`       | Command | Quick health overview including state dir |

## Understanding State Files

State files fall into two categories:

**Session-scoped** (contain a UUID in the filename): Created per-session to
suppress repeat nudges. Safe to prune once the session ends. Examples:

```
context-check-11e94c1d-1639-4c04-bf77-63dcf1f50ec7
heartbeat-11e94c1d-1639-4c04-bf77-63dcf1f50ec7
persistence-nudge-11e94c1d-1639-4c04-bf77-63dcf1f50ec7
```

**Global** (no UUID): Persist across sessions. `ctx prune` preserves
these automatically. Some are legitimate state (`events.jsonl`,
`memory-import.json`); others may be stale tombstones that need manual
review.

## The Workflow

### Step 1: Preview

Always dry-run first to see what would be removed:

```bash
ctx prune --dry-run
```

The output shows each file, its age, and a summary:

```text
  would prune: context-check-abc123... (age: 3d)
  would prune: heartbeat-abc123... (age: 3d)

Dry run - would prune 150 files (skip 70 recent, preserve 14 global)
```

### Step 2: Prune

Choose an age threshold. The default is 7 days:

```bash
ctx prune               # older than 7 days
ctx prune --days 3      # older than 3 days
ctx prune --days 1      # older than 1 day (aggressive)
```

### Step 3: Review Global Files

After pruning, check what `prune` preserved:

```bash
ls .context/state/ | grep -v '[0-9a-f]\{8\}-[0-9a-f]\{4\}'
```

Legitimate global files (keep):

- `events.jsonl` - event log
- `memory-import.json` - import tracking state

Stale global tombstones (safe to delete):

- Files like `backup-reminded`, `ceremony-reminded`, `version-checked`
  with no session UUID are one-shot markers. If they are from a previous
  session, they are stale and can be removed manually.

```bash
rm .context/state/backup-reminded .context/state/ceremony-reminded
```

### Step 4: Verify

```bash
ls .context/state/ | wc -l    # should be manageable
```

## When to Prune

- **Weekly**: `ctx prune` with default 7-day threshold
- **After heavy parallel work**: Multiple concurrent sessions create
  many tombstones. Prune with `--days 1` afterward.
- **When state directory exceeds ~100 files**: A sign that pruning
  hasn't run recently

## Tips

**Pruning active sessions is safe but noisy**: If you prune a file belonging
to a still-running session, the corresponding hook will re-fire its nudge
on the next prompt. Minor UX annoyance, not data loss.

**No context files are stored in state**: The state directory contains only
tombstones, counters, and diagnostic data. Nothing in `.context/state/`
affects your decisions, learnings, tasks, or conventions.

**Test artifacts sneak in**: Files like `context-check-statstest` or
`heartbeat-unknown` are artifacts from development or testing. They lack
UUIDs so `prune` preserves them. Delete manually.

## See Also

* [Detecting and Fixing Drift](context-health.md): broader context
  maintenance including drift detection and archival
* [Troubleshooting](troubleshooting.md): diagnostic workflow using
  `ctx doctor` and event logs
* [CLI Reference: system](../cli/system.md): full flag documentation
  for `ctx prune` and related commands
