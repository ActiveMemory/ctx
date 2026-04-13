---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Getting Started
icon: lucide/share-2
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: Getting Started

Stand up a **single-node** `ctx` Hub on localhost, register
two projects, publish a decision from one, and see it appear in the
other — all in under five minutes.

!!! tip "Read this first"
    If you haven't already, skim the
    [`ctx` Hub overview](hub-overview.md). It explains the
    mental model, names the two user stories (personal vs small
    team), and — importantly — lists what the hub **does not do**.
    This recipe assumes you already know you want the feature.

## What you'll get out of this recipe

By the end, you will have:

1. A local hub process running on port `9900`.
2. Two project directories both registered with the ctx Hub.
3. A decision published from project `alpha` that appears
   automatically in project `beta`'s `.context/hub/` and in
   `ctx agent --include-hub` output.

Concretely, the payoff this unlocks: a lesson you record in one
project becomes visible to your agent the next time you open
another project — **without** touching local files in the second
project or opening another editor window.

## What this recipe does *not* cover

- Sharing `.context/journal/`, `.context/pad`, or any other
  local state. The hub only fans out `decision`, `learning`,
  `convention`, and `task` entries. Everything else stays local.
- Multi-user attribution. The hub identifies **projects**, not
  people.
- Running over a LAN — see
  [Multi-machine setup](hub-multi-machine.md).
- Redundancy — see [HA cluster](hub-cluster.md).

## Prerequisites

- `ctx` installed and on `PATH`
- Two project directories, each already initialized with
  `ctx init`

## Step 1 — Start the hub

In a dedicated terminal:

```bash
ctx hub start
```

On first run, the hub generates an **admin token** and prints it to
stdout. Copy it — you'll need it for each project registration:

```
ctx hub listening on :9900
admin token: ctx_adm_7f3a1c2d...
data dir: ~/.ctx/hub-data/
```

The admin token is written to `~/.ctx/hub-data/admin.token` so you
can recover it later. Treat it like a password.

## Step 2 — Register the first project

```bash
cd ~/projects/alpha
ctx connection register localhost:9900 --token ctx_adm_7f3a1c2d...
```

This stores an **encrypted** connection config in
`.context/.connect.enc`. The admin token is exchanged for a
per-project client token; the admin token itself is never persisted
in the project.

## Step 3 — Choose what to receive

```bash
ctx connection subscribe decision learning convention
```

Only the entry types you subscribe to will be delivered by `sync`
and `listen`.

## Step 4 — Publish a decision

Either use `ctx add --share` to write locally *and* push to the ctx Hub:

```bash
ctx add decision "Use UTC timestamps everywhere" --share \
  --context "We had timezone drift between the API and journal" \
  --rationale "Single source of truth avoids conversion bugs" \
  --consequence "The UI does conversion at render time"
```

Or publish an existing entry directly:

```bash
ctx connection publish decision "Use UTC timestamps everywhere"
```

## Step 5 — Register a second project and sync

```bash
cd ~/projects/beta
ctx connection register localhost:9900 --token ctx_adm_7f3a1c2d...
ctx connection subscribe decision learning convention
ctx connection sync
```

The decision from `alpha` now appears in
`~/projects/beta/.context/hub/decisions.md` with an origin tag
and timestamp.

## Step 6 — Watch entries arrive live

Instead of re-running `sync`, stream new entries as they land:

```bash
ctx connection listen
```

Leave this running in a terminal; every `--share` publish from any
registered project will appear in `.context/hub/` immediately.

## Step 7 — Feed shared knowledge into the agent

Once entries exist in `.context/hub/`, include them in the agent
context packet:

```bash
ctx agent --include-hub
```

Shared entries are added as a dedicated tier in the budget-aware
assembly, scored by recency and type relevance.

## Auto-sync on session start

After `register`, the `check-hub-sync` hook pulls new entries at the
start of each session (daily throttled). Most users never need to
call `ctx connection sync` manually.

## Where to go next

- **[Multi-machine hub](hub-multi-machine.md)**: run the hub
  on a LAN host and connect from other workstations.
- **[HA cluster](hub-cluster.md)**: Raft-based leader
  election for high availability.
- **[Hub operations](../operations/hub.md)**: daemon mode,
  backup, log rotation, JSONL store layout.
- **[Hub security model](../security/hub.md)**: token
  lifecycle, encryption at rest, threat model.
- **[`ctx connect` reference](../cli/connection.md)** and
  **[`ctx hub start` reference](../cli/serve.md)**.
