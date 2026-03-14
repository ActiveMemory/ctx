---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Reviewing Session Changes
---

## What Changed While You Were Away?

Between sessions, teammates commit code, context files get updated,
and decisions pile up. `ctx changes` gives you a single-command
summary of everything that moved since your last session.

## Quick Start

```bash
# Auto-detects your last session and shows what changed
ctx changes

# Check what changed in the last 48 hours
ctx changes --since 48h

# Check since a specific date
ctx changes --since 2026-03-10
```

## How Reference Time Works

`ctx changes` needs a reference point to compare against. It tries
these sources in order:

1. **`--since` flag** — explicit duration (`24h`, `72h`) or date
   (`2026-03-10`, RFC3339 timestamp)
2. **Session markers** — `ctx-loaded-*` files in `.context/state/`;
   picks the second-most-recent (your *previous* session start)
3. **Event log** — last `context-load-gate` event from
   `.context/state/events.jsonl`
4. **Fallback** — 24 hours ago

The marker-based detection means `ctx changes` usually just works
without any flags: it knows when you last loaded context and shows
everything after that.

## What It Reports

### Context file changes

Any `.md` file in `.context/` modified after the reference time:

```
### Context File Changes
- `TASKS.md` — modified 2026-03-11 14:30
- `DECISIONS.md` — modified 2026-03-11 09:15
```

### Code changes

Git activity since the reference time:

```
### Code Changes
- **12 commits** since reference point
- **Latest**: Fix journal enrichment ordering
- **Directories touched**: internal, docs, specs
- **Authors**: jose, claude
```

## Integrating Into Session Start

Pair `ctx changes` with the `/ctx-remember` ceremony for a complete
session-start picture:

```bash
# 1. Load context (this also creates the session marker)
ctx agent --budget 4000

# 2. See what changed since your last session
ctx changes
```

Or script it:

```bash
# .context/hooks/session-start.sh
ctx agent --budget 4000
echo "---"
ctx changes
```

## Team Workflows

When multiple people share a `.context/` directory, `ctx changes`
shows who changed what:

```bash
# After pulling from remote
git pull
ctx changes --since 72h
```

This surfaces context file changes from teammates that you might
otherwise miss in the commit log.

## Tips

- **No changes?** If nothing shows up, the reference time might be
  wrong. Use `--since 48h` to widen the window.
- **Works without git.** Context file changes are detected by
  filesystem mtime, not git. Code changes require git.
- **Hook integration.** The `context-load-gate` hook writes the
  session marker that `ctx changes` uses for auto-detection. If
  you're not using the ctx plugin, markers won't exist and it falls
  back to the event log or 24h window.
