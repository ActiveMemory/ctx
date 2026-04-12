---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Overview
icon: lucide/compass
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: Overview

Start here before the other hub recipes. This page answers *what*
the hub is, *who* it's for, *why* you'd run one, and —
equally important — *what it is not*.

## Mental model in one paragraph

The hub is a **fan-out channel for structured knowledge
entries across projects**. When you publish a decision, learning,
convention, or task with `--share`, the hub stores it in an
append-only log and delivers it to every other project
subscribed to that type. The next time your agent loads context
in any of those projects, shared entries can be included in the
context packet alongside local ones.

That's the whole feature. It is a **project-to-project
knowledge bus** for a small, curated set of entry types. It is
**not** a shared memory, a shared journal, or a multi-user
database.

## What flows through the hub

Only four entry types:

| Type         | What it is                                |
|--------------|-------------------------------------------|
| `decision`   | Architectural decisions with rationale    |
| `learning`   | Gotchas, lessons, surprising behaviors    |
| `convention` | Coding patterns and standards             |
| `task`       | Work items worth sharing across projects  |

Each entry is an immutable record with a content blob, the
publishing project's name as `Origin`, a timestamp, and a
hub-assigned sequence number. Once published, entries are
never rewritten.

## What does *not* flow through the hub

This is the part new users get wrong most often:

- **Session journals** (`~/.claude/` logs, `.context/journal/`)
  stay local. The hub does **not** sync your AI session history.
- **Scratchpad** (`.context/pad`) stays local. Encrypted notes
  never leave the machine they were written on.
- **Local context files** as a whole — `TASKS.md`,
  `DECISIONS.md`, `LEARNINGS.md`, `CONVENTIONS.md` — are **not**
  mirrored wholesale. Only entries you explicitly `--share`, or
  publish later with `ctx connection publish`, cross the boundary.
- **Anything under `.context/` that isn't one of the four entry
  types above.** Configuration, state, logs, memory, journal
  metadata — all local.

If you were expecting "now my agent in project B can see
everything my agent did in project A," that's not this feature.
Local session density still lives on the local machine.

## Two user stories

The hub makes sense in two different shapes. Pick the one that
matches your situation — the mechanics are identical but the
trust model and threat surface are very different.

### Story 1: Personal cross-project brain

**One developer, many projects, one hub — usually on localhost.**

You're working across several projects on the same machine (or a
handful of machines you own). You want a lesson learned
debugging project A to show up when you open project B a week
later, without re-discovering it. You want a convention you
codified in one project to be visible as-you-type in another.

**Concrete payoff:**

- `ctx add learning --share "..."` in project A →
  `ctx agent --include-hub` in project B shows that learning
  in the next context packet.
- A decision recorded in your personal "dotfiles" project is
  instantly visible to every other project on your workstation.
- Cross-project conventions (e.g., "use UTC timestamps
  everywhere") live in one place and propagate.

**Trust model:** high — you trust every participant because every
participant is *you*. Run the hub on localhost or on your own
LAN, use the default single-node setup, don't worry about TLS.

**Start here:**
[Getting Started](hub-getting-started.md) for the one-time
setup, then [Personal cross-project brain](hub-personal.md)
for the day-to-day workflow.

### Story 2: Small trusted team

**A few teammates, projects they each own, one hub on a LAN host
they all trust.**

Your team has a handful of services and you want a shared
"things we've learned the hard way" stream. Someone on the
platform team records a convention about timestamp handling;
everyone else's agents see it the next session. An on-call
engineer records a learning from a 3 AM incident; the
rest of the team inherits the lesson without needing to read
the postmortem.

**Concrete payoff:**

- Team conventions propagate without needing a wiki or chat.
- Lessons from one team member become available to everyone
  else's agent context packets automatically.
- Cross-project decisions (shared libraries, deployment
  patterns, naming rules) live in a single log the whole team
  reads.

**Trust model:** the hub assumes **everyone holding a client
token is friendly.** There is no per-user attribution you can
rely on, `Origin` is self-asserted by the publishing client, and
there is no read ACL beyond the subscription filter. Treat the
hub like a team wiki: useful because everyone can write to it,
not because it can prove *who* wrote what.

**Operational shape:** run the hub on a LAN host (or a
three-node HA cluster for redundancy), put TLS in front of it
for anything beyond a home LAN, distribute client tokens over a
trusted channel.

**Start here:**
[Multi-machine setup](hub-multi-machine.md) for the
deployment, [Team knowledge bus](hub-team.md) for the
day-to-day team workflow, then [HA cluster](hub-cluster.md)
if you need redundancy.

## Identity: projects, not users

The hub has **no concept of users.** Its unit of identity is the
*project*. `ctx connection register` binds a hub token to a project
directory, not to a person. Two developers working on the same
project share either:

- **The same `.connect.enc`**, copied between machines over a
  trusted channel, or
- **Different project names** (`alpha@laptop-a`,
  `alpha@laptop-b`), because the hub rejects duplicate
  registrations of the same project name.

Either works; neither gives you per-human attribution. If you
need "who wrote this," the hub is the wrong tool.

## When *not* to use it

- **Solo, single-project work.** Local `.context/` files are
  enough. The hub adds operational surface for no payoff.
- **Untrusted participants.** The hub assumes everyone with a
  client token is friendly. It is not hardened against hostile
  insiders or compromised tokens.
- **Compliance-sensitive environments.** There is no audit
  trail that can prove *who* published what, only *which
  project* published what, and `Origin` is self-asserted.
- **Secrets or PII.** Entry content is stored plaintext on the
  hub and fanned out to every subscribed client. Don't publish
  anything you wouldn't paste in a team chat.
- **Wholesale journal sharing.** See "what does not flow"
  above. If that's what you want, this feature won't provide
  it — talk to us in the issue tracker about what *would*.

## How entries reach your agent

Once a project is registered and subscribed, entries arrive by
three mechanisms:

1. **`ctx connection sync`** — an on-demand pull, replays
   everything new since the last sequence you saw.
2. **`ctx connection listen`** — a long-lived gRPC stream that
   writes new entries to `.context/hub/` as they arrive.
3. **`check-hub-sync` hook** — runs at session start, daily
   throttled, so most users never call `sync` manually.

Once entries exist in `.context/hub/`, `ctx agent
--include-hub` adds a dedicated tier to the budget-aware
context packet, scored by recency and type relevance. That's
the end of the pipeline.

## Where to go next

| If you're…                                        | Read                                             |
|---------------------------------------------------|--------------------------------------------------|
| Trying it for yourself on one machine             | [Getting Started](hub-getting-started.md)        |
| A solo developer using the hub day-to-day         | [Personal cross-project brain](hub-personal.md)  |
| Setting up for a small team on a LAN              | [Multi-machine setup](hub-multi-machine.md)      |
| A small team using the hub day-to-day             | [Team knowledge bus](hub-team.md)                |
| Running redundant nodes                           | [HA cluster](hub-cluster.md)                     |
| Operating a hub in production                     | [Operations](../operations/hub.md)               |
| Assessing the security posture                    | [Security model](../security/hub.md)             |
| Debugging a hub in trouble                        | [Failure modes](../operations/hub-failure-modes.md) |
| Just reading the commands                         | [`ctx connect`](../cli/connection.md), [`ctx serve`](../cli/serve.md), [`ctx hub`](../cli/hub.md) |
