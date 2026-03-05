---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Bridging Claude Code Auto Memory"
icon: lucide/brain
---

![ctx](../images/ctx-banner.png)

## The Problem

Claude Code maintains per-project auto memory at
`~/.claude/projects/<slug>/memory/MEMORY.md`. This file is:

- **Outside the repo** — not version-controlled, not portable
- **Machine-specific** — tied to one `~/.claude/` directory
- **Invisible to ctx** — context loading and hooks don't read it

Meanwhile, ctx maintains structured context files (DECISIONS.md,
LEARNINGS.md, CONVENTIONS.md) that are git-tracked, portable, and
token-budgeted — but Claude Code doesn't automatically write to them.

The two systems hold complementary knowledge with no bridge between them.

## TL;DR

```bash
ctx memory sync          # Mirror MEMORY.md into .context/memory/mirror.md
ctx memory status        # Check for drift
ctx memory diff          # See what changed since last sync
```

The `check-memory-drift` hook nudges automatically when MEMORY.md
changes — you don't need to remember to sync manually.

## Commands and Skills Used

| Tool                              | Type        | Purpose                                          |
|-----------------------------------|-------------|--------------------------------------------------|
| `ctx memory sync`                 | CLI command | Copy MEMORY.md to mirror, archive previous       |
| `ctx memory status`               | CLI command | Show drift, timestamps, line counts              |
| `ctx memory diff`                 | CLI command | Show changes since last sync                     |
| `ctx memory import`               | CLI command | Classify and promote entries to .context/ files  |
| `ctx memory publish`              | CLI command | Push curated .context/ content to MEMORY.md      |
| `ctx memory unpublish`            | CLI command | Remove published block from MEMORY.md            |
| `ctx system check-memory-drift`   | Hook        | Nudge when MEMORY.md has changed (once/session)  |

## How It Works

### Discovery

Claude Code encodes project paths as directory names under
`~/.claude/projects/`. The encoding replaces `/` with `-` and
prefixes with `-`:

```
/home/jose/WORKSPACE/ctx  →  ~/.claude/projects/-home-jose-WORKSPACE-ctx/
```

`ctx memory` uses this encoding to locate MEMORY.md automatically
from your project root — no configuration needed.

### Mirroring

When you run `ctx memory sync`:

1. The previous mirror is archived to `.context/memory/archive/mirror-<timestamp>.md`
2. MEMORY.md is copied to `.context/memory/mirror.md`
3. Sync state is updated in `.context/state/memory-import.json`

The mirror is git-tracked, so it travels with the project. Archives
provide a fallback for projects that don't use git.

### Drift Detection

The `check-memory-drift` hook compares MEMORY.md's modification time
against the mirror. When drift is detected, the agent sees:

```
┌─ Memory Drift ────────────────────────────────────────────────
│ MEMORY.md has changed since last sync.
│ Run: ctx memory sync
│ Context: .context
└────────────────────────────────────────────────────────────────
```

The nudge fires once per session to avoid noise.

## Typical Workflow

### At Session Start

If the hook fires a drift nudge, sync before diving into work:

```bash
ctx memory diff     # Review what changed
ctx memory sync     # Mirror the changes
```

### Periodic Check

```bash
ctx memory status
# Memory Bridge Status
#   Source:      ~/.claude/projects/.../memory/MEMORY.md
#   Mirror:      .context/memory/mirror.md
#   Last sync:   2026-03-05 14:30 (2 hours ago)
#
#   MEMORY.md:  47 lines
#   Mirror:     32 lines
#   Drift:      detected (source is newer)
#   Archives:   3 snapshots in .context/memory/archive/
```

### Dry Run

Preview what sync would do without writing:

```bash
ctx memory sync --dry-run
```

## Storage Layout

```
.context/
├── memory/
│   ├── mirror.md                          # Raw copy of MEMORY.md (git-tracked)
│   └── archive/
│       ├── mirror-2026-03-05-143022.md    # Timestamped pre-sync snapshots
│       └── mirror-2026-03-04-220015.md
├── state/
│   └── memory-import.json                 # Sync tracking state
```

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| Auto memory not active | `sync` exits 1 with message. `status` reports "not active". Hook skips silently. |
| First sync (no mirror) | Creates mirror without archiving. |
| MEMORY.md is empty | Syncs to empty mirror (valid). |
| Not initialized | Init guard rejects (same as all ctx commands). |

## Importing Entries

Once you've synced, you can classify and promote entries into structured
`.context/` files:

```bash
ctx memory import --dry-run    # Preview classification
ctx memory import              # Actually promote entries
```

Each entry is classified by keyword heuristics:

| Keywords | Target |
|----------|--------|
| `always use`, `prefer`, `never use`, `standard` | CONVENTIONS.md |
| `decided`, `chose`, `trade-off`, `approach` | DECISIONS.md |
| `gotcha`, `learned`, `watch out`, `bug`, `caveat` | LEARNINGS.md |
| `todo`, `need to`, `follow up` | TASKS.md |
| Everything else | Skipped |

Entries that don't match any pattern are skipped — they stay in the mirror
for manual review. Deduplication (hash-based) prevents re-importing the
same entry on subsequent runs.

!!! tip "Review Before Importing"
    Use `--dry-run` first. The heuristic classifier is deliberately simple —
    it may misclassify ambiguous entries. Review the plan, then import.

### Full Workflow

```bash
ctx memory sync                # 1. Mirror MEMORY.md
ctx memory import --dry-run    # 2. Preview what would be imported
ctx memory import              # 3. Promote entries to .context/ files
```

## Publishing Context to MEMORY.md

Push curated `.context/` content back into MEMORY.md so Claude Code sees
structured project context on session start — without needing hooks.

```bash
ctx memory publish --dry-run    # Preview what would be published
ctx memory publish              # Write to MEMORY.md
ctx memory publish --budget 40  # Tighter line budget
```

Published content is wrapped in markers:

```markdown
<!-- ctx:published -->
# Project Context (managed by ctx)

## Pending Tasks
- [ ] Implement feature X
...
<!-- ctx:end -->
```

**Rules:**

- ctx owns everything **between** the markers
- Claude owns everything **outside** the markers
- `ctx memory import` reads only outside the markers
- `ctx memory publish` replaces only inside the markers

To remove the published block entirely:

```bash
ctx memory unpublish
```

!!! tip "Publish at Wrap-Up, Not on Commit"
    The best time to publish is during session wrap-up, after persisting
    decisions and learnings. Never auto-publish — give yourself a chance
    to review what's going into MEMORY.md.

### Full Bidirectional Workflow

```bash
ctx memory sync                 # 1. Mirror MEMORY.md
ctx memory import --dry-run     # 2. Check what Claude wrote
ctx memory import               # 3. Promote entries to .context/
ctx memory publish --dry-run    # 4. Check what would be published
ctx memory publish              # 5. Push context to MEMORY.md
```
