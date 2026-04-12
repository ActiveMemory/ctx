---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Connect
icon: lucide/link
---
![ctx](../images/ctx-banner.png)


## `ctx connect`

Connect a project to a `ctx` Hub for cross-project
knowledge sharing. Projects publish decisions, learnings,
conventions, and tasks to a hub; other subscribed projects receive
them alongside local context.

!!! tip "New to the `ctx` Hub?"
    Start with the
    [`ctx` Hub overview](../recipes/hub-overview.md) for
    the mental model (what the hub is, who it's for, what it is
    **not**), then walk through
    [Getting Started](../recipes/hub-getting-started.md).
    This page is a command reference, not an introduction.

**The unit of identity is a project, not a user.** Registering a
directory with `ctx connection register` binds a per-project client
token in `.context/.connect.enc`. Two developers on the same
project either share that file over a trusted channel, or each
register under a different project name.

**Only structured entries flow through the hub** — `decision`,
`learning`, `convention`, `task`. Session journals, scratchpad
contents, and other local state stay on the machine that created
them.

### `ctx connection register`

One-time registration with a ctx Hub. Requires the ctx Hub address and
admin token (printed by `ctx hub start` on first run).

**Examples**:

```bash
ctx connection register localhost:9900 --token ctx_adm_7f3a...
```

On success, stores an encrypted connection config in
`.context/.connect.enc` for future RPCs.

### `ctx connection subscribe`

Set which entry types to receive from the ctx Hub. Only matching types
are returned by sync and listen.

**Examples**:

```bash
ctx connection subscribe decision learning
ctx connection subscribe decision learning convention
```

### `ctx connection sync`

Pull matching entries from the ctx Hub and write them to
`.context/hub/` as markdown files with origin tags and date
headers. Tracks last-seen sequence for incremental sync.

**Examples**:

```bash
ctx connection sync
```

### `ctx connection publish`

Push entries to the ctx Hub. Specify type and content as arguments.

**Examples**:

```bash
ctx connection publish decision "Use UTC timestamps everywhere"
ctx connection publish learning "Go embed requires files in same package"
```

### `ctx connection listen`

Stream new entries from the ctx Hub in real-time. Writes to
`.context/hub/` as entries arrive. Press Ctrl-C to stop.

**Examples**:

```bash
ctx connection listen
```

### `ctx connection status`

Show ctx Hub connection state and entry statistics.

**Examples**:

```bash
ctx connection status
```

## Automatic sharing

Use `--share` on `ctx add` to write locally AND publish to the ctx Hub:

```bash
ctx add decision "Use UTC" --share \
  --context "Need consistency" \
  --rationale "Avoid timezone bugs" \
  --consequence "UI does conversion"
```

If the hub is unreachable, the local write succeeds and a warning
is printed. The `--share` flag is best-effort — it never blocks
local context updates.

## Auto-sync

Once registered, the `check-hub-sync` hook automatically syncs
new entries from the ctx Hub at the start of each session (daily
throttled). No manual `ctx connection sync` needed.

## Shared files

Entries from the ctx Hub are stored in `.context/hub/`:

```
.context/hub/
  decisions.md      # Shared decisions with origin tags
  learnings.md      # Shared learnings
  conventions.md    # Shared conventions
  .sync-state.json  # Last-seen sequence tracker
```

These files are read-only (managed by sync/listen) and never
mixed with local context files.

## Agent integration

Include shared knowledge in agent context packets:

```bash
ctx agent --include-hub
```

Shared entries are included as Tier 8 in the budget-aware
assembly, scored by recency and type relevance.
