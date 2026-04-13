---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Personal cross-project brain
icon: lucide/brain
---

![ctx](../images/ctx-banner.png)

# Personal cross-project brain

This recipe shows **how one developer uses a `ctx` Hub
across their own projects day-to-day** — the "Story 1"
shape from the [Hub overview](hub-overview.md). You're not
setting up infrastructure for a team; you're making a
lesson you learned last Tuesday in project A automatically
surface when you open project B next Thursday.

**Prerequisites**: a working `ctx` Hub on localhost
(see [Getting Started](hub-getting-started.md) for the
~5-minute setup). This recipe assumes the hub is already
running and you've registered at least two projects.

## The core loop

Every day, the same three verbs matter:

1. **Record** — notice a decision, learning, or
   convention and capture it with `ctx add --share`.
2. **Subscribe** — every project you care about is
   subscribed to the types you want delivered (set once
   with `ctx connection subscribe`).
3. **Load** — your agent picks up shared entries on next
   session start via the auto-sync hook, or explicitly
   via `ctx agent --include-hub`.

That's the whole workflow. The rest of this recipe fills
in the concrete moments where each verb matters.

## A realistic day

You have three projects on your workstation:

- `~/projects/api` — a Go service you're actively
  developing
- `~/projects/cli` — a companion CLI that consumes the
  API
- `~/projects/dotfiles` — your personal conventions and
  cross-project learnings

All three are registered with a single hub running on
`localhost:9900` (started once at boot, or via a systemd
user unit — see [Hub operations](../operations/hub.md)).
All three subscribe to `decision`, `learning`, and
`convention`.

### 09:00 — Start work on `api`

You `cd ~/projects/api` and start a Claude Code session.
Behind the scenes, the plugin's `PreToolUse` hook calls
`ctx agent --budget 8000 --include-hub` before the first
tool call. Agent loads:

- Local `.context/` (TASKS, DECISIONS, LEARNINGS, etc.)
- Foundation steering files (always-inclusion)
- **Everything you've shared from the other two projects**

So the "use UTC timestamps everywhere" decision you
recorded in `dotfiles` last week is already in Claude's
context for this session, without any manual `sync`.

### 10:30 — You discover a gotcha

While debugging, you find that the API's retry loop
silently drops the last error when the transport times
out. This is the kind of thing you'd normally add to
`LEARNINGS.md` in `api/`. But it's useful across every
Go service you'll ever write, not just this one. So:

```bash
ctx add learning --share \
  --context "Go http.Client retries mask the final error" \
  --lesson  "Transport timeouts don't surface as errors when the retry loop re-assigns err without wrapping. Check for context.DeadlineExceeded on the request context instead." \
  --application "Any retry loop over http.Client.Do that uses a per-attempt timeout"
```

The `--share` flag does two things:

1. Writes the learning to `api/.context/LEARNINGS.md`
   locally (as a normal `ctx add learning` would).
2. Publishes the same entry to the ctx Hub, which stores it
   in the append-only JSONL and fans it out to every
   subscribed client.

Within seconds, `cli/.context/hub/learnings.md` and
`dotfiles/.context/hub/learnings.md` both contain a copy
of this learning (the `ctx connection listen` daemon picks
it up from the ctx Hub's Listen stream).

### 12:00 — You switch to `cli`

`cd ~/projects/cli`, open a new session. The agent
packet for `cli` now includes **the learning you just
recorded in `api`**, because `cli` is subscribed to
`learning` and the entry has already been synced into
`cli/.context/hub/learnings.md`.

You don't have to re-explain the retry-loop gotcha.
Claude already sees it.

### 14:00 — You codify a convention

You've been writing error messages in `api` and decided
you want a consistent pattern: lowercase start, no
trailing period, single-sentence. This is a convention,
not a decision — it applies to every Go project you
touch. Record it in `dotfiles` (since that's your
"personal standards" project), and share it:

```bash
cd ~/projects/dotfiles
ctx add convention --share \
  "Error messages: lowercase start, no trailing period, single sentence (follows Go's stdlib style)"
```

The convention lands in `dotfiles/CONVENTIONS.md` locally
and fans out to `api` and `cli` via the hub. The next
Claude Code session in either project gets the
convention injected into the steering-adjacent slot of
the agent packet.

### 16:30 — End of day

You didn't run `ctx connection sync` once. You didn't
`git push` anything between projects. You didn't
remember to tell your agent about the retry-loop gotcha
in the new project. The hub did all of it for you.

## What the workflow actually looks like

Stripped of prose, the day's commands were:

```bash
# Morning: nothing. Agent loads --include-hub automatically.

# Mid-morning: record a learning that should cross projects
ctx add learning --share \
  --context "..." --lesson "..." --application "..."

# Afternoon: codify a convention in the "standards" project
ctx add convention --share "..."

# Evening: nothing. Everything's already propagated.
```

The hub is passive infrastructure. You never talk **to**
it directly — you talk **through** it by using `--share`
on commands you were already running.

## Tips for solo use

**Pick a "standards" project.** One of your projects
should play the role of "canonical source for rules you
want everywhere." Your dotfiles, a personal scratch
repo, or a dedicated `ctx-standards` project all work.
Record cross-cutting conventions there and let the hub
propagate them to everything else.

**Subscribe to `task` only if you want cross-project
todos.** The four subscribable types are `decision`,
`learning`, `convention`, `task`. Tasks are usually
project-local; subscribing makes every hub-shared task
from every project show up in every other project's
agent packet. That's probably not what you want. Skip
`task` in `ctx connection subscribe` unless you have a
specific reason.

**Run the hub as a user-level daemon** so you don't have
to remember to start it. On Linux with systemd:

```ini
# ~/.config/systemd/user/ctx-hub.service
[Unit]
Description=ctx Hub (personal)

[Service]
Type=simple
ExecStart=/usr/local/bin/ctx hub start
Restart=on-failure

[Install]
WantedBy=default.target
```

```bash
systemctl --user enable --now ctx-hub.service
```

**Don't overthink subscription filters.** For personal
use, subscribe every project to all four types at first
(or three, if you skip `task`). Tune later if the
context packets get noisy.

**Local storage is fine; no TLS needed.** The hub runs
on localhost. No one else is on the network. Skip the
TLS setup from the
[Multi-machine recipe](hub-multi-machine.md) — it's
relevant when the hub is on a LAN host serving multiple
workstations, not when it's a personal daemon.

## What this recipe is *not*

**Not a setup guide.** For the one-time hub install and
project registration, use
[Getting Started](hub-getting-started.md).

**Not a team guide.** If you're sharing across humans,
not just across your own projects, read
[Team knowledge bus](hub-team.md) instead — the trust
model and operational concerns are different.

**Not production operations.** For backup, log
rotation, failure recovery, and HA, see
[Hub operations](../operations/hub.md) and
[Hub failure modes](../operations/hub-failure-modes.md).

## See also

- [Hub overview](hub-overview.md) — when to use the Hub
  and when not to.
- [Team knowledge bus](hub-team.md) — the multi-human
  companion recipe.
- [`ctx connect`](../cli/connection.md) — the client-side
  commands used above (`subscribe`, `publish`, `sync`,
  `listen`, `status`).
- [`ctx add`](../cli/context.md) — the `--share` flag
  reference.
- [`ctx hub`](../cli/hub.md) — operator commands for
  starting, stopping, and inspecting the hub.
