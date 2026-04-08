---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Connect
icon: lucide/link
---

## `ctx connect`

Connect to a shared context hub for cross-project knowledge sharing.
Projects publish decisions, learnings, and conventions to a centralized
hub; other projects receive them alongside local context.

### `ctx connect register`

One-time registration with a shared hub. Requires the hub address and
admin token (printed by `ctx serve --shared` on first run).

```bash
ctx connect register localhost:9900 --token ctx_adm_7f3a...
```

On success, stores an encrypted connection config in
`.context/.connect.enc` for future RPCs.

### `ctx connect subscribe`

Set which entry types to receive from the hub. Only matching types
are returned by sync and listen.

```bash
ctx connect subscribe decision learning
ctx connect subscribe decision learning convention
```

### `ctx connect sync`

Pull matching entries from the hub and write them to
`.context/shared/` as markdown files with origin tags and date
headers. Tracks last-seen sequence for incremental sync.

```bash
ctx connect sync
```

### `ctx connect publish`

Push entries to the hub. Specify type and content as arguments.

```bash
ctx connect publish decision "Use UTC timestamps everywhere"
```

### `ctx connect listen`

Stream new entries from the hub in real-time. Writes to
`.context/shared/` as entries arrive. Press Ctrl-C to stop.

```bash
ctx connect listen
```

### `ctx connect status`

Show hub connection state and entry statistics.

```bash
ctx connect status
```

## Shared files

Entries from the hub are stored in `.context/shared/`:

```
.context/shared/
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
ctx agent --include-shared
```

Shared entries are included as Tier 8 in the budget-aware
assembly, scored by recency and type relevance.
