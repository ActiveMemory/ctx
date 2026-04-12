---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: The ctx Hub
icon: lucide/network
---

![ctx](../images/ctx-banner.png)

## The `ctx` Hub

`ctx` projects are normally **independent**: each project has its
own `.context/` directory, its own decisions, its own learnings,
its own journal. That's the right default — most work is
project-local, and mixing context across projects tends to dilute
more than it helps.

But sometimes a decision or a learning **should** cross project
boundaries. A convention you codified in one project deserves to
be visible in another. A gotcha you discovered debugging service
A is the same gotcha waiting for you in service B. The **`ctx`
Hub** is the feature that makes those specific entries travel,
without replicating everything else.

## What the Hub actually is

In one paragraph: the `ctx` Hub is a **fan-out channel** for
four specific kinds of structured entries — `decision`,
`learning`, `convention`, and `task`. You publish an entry with
`ctx add --share` in one project, and it appears in
`.context/hub/` for every other project subscribed to that
type. When you run `ctx agent --include-hub`, those shared
entries become part of your next agent context packet.

That is the **entire** feature. The Hub does **not**:

- Share your session journal (`.context/journal/`). That stays
  local to each project.
- Share your scratchpad (`.context/pad`). Encrypted notes never
  leave the machine that created them.
- Share your `TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`, or
  `CONVENTIONS.md` wholesale. Only entries you explicitly
  `--share` cross the boundary.
- Provide user identity or attribution. The Hub identifies
  **projects**, not people.

If you want "my agent in project B sees everything my agent did
in project A," that's not the Hub. Local session density stays
local.

## Who it's for

Two shapes, same mechanics, different trust models.

### Personal cross-project brain

**One developer, many projects.** You want a learning from
project A to show up when you open project B a week later. You
want a convention you codified in your dotfiles project to be
visible everywhere else on your workstation. Run a Hub on
localhost, register each project, done.

### Small trusted team

**A few teammates on a LAN or a hub.ctx-like self-hosted
server.** You want team conventions to propagate without a
wiki. You want lessons from one on-call engineer's 3 AM
incident to reach everyone else's agent on the next session.
Same mechanics as the personal case, plus TLS in front and a
short security runbook.

The Hub is **not** a multi-tenant public service. It assumes
everyone holding a client token is friendly. Don't stand up
`hub.example.com` for untrusted participants.

## Going further

- **First-time setup:** [Hub: Getting Started](../recipes/hub-getting-started.md) —
  a five-minute walkthrough on localhost.
- **Mental model and user stories:** [Hub Overview](../recipes/hub-overview.md) —
  what flows, what doesn't, and when not to use it.
- **Team / LAN deployment:** [Multi-machine setup](../recipes/hub-multi-machine.md).
- **Redundancy:** [HA cluster](../recipes/hub-cluster.md).
- **Operating a Hub:** [Hub Operations](../operations/hub.md)
  and [Hub Failure Modes](../operations/hub-failure-modes.md).
- **Security posture:** [Hub Security Model](../security/hub.md).
- **Command reference:** [`ctx serve`](../cli/serve.md),
  [`ctx connect`](../cli/connection.md),
  [`ctx hub`](../cli/hub.md).
