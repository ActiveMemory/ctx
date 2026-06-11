---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Dream
icon: lucide/moon
---

![ctx](../images/ctx-banner.png)

## `ctx dream`

Run a disciplined, out-of-band **dream** pass over the gitignored
`ideas/` folder: classify each idea against the codebase and specs, and
emit gated, provenance-bearing disposition **proposals** into the
`dreams/` notebook for human review. The dream only ever proposes — it
never writes canonical memory and never acts on a proposal.

The dream is **opt-in and off by default**. Nothing runs until you set
`dream.enabled: true` in `.ctxrc`. See the
[Run the Dream](../recipes/run-the-dream.md) recipe for the full setup
(cron, guard hook, review), and the
[executor contract](../reference/dream-executor-contract.md) to run it
under a non-Claude-Code harness.

Invoked with no subcommand, it runs one bounded pass: it gates on the
idea delta and the quiet window, takes an exclusive lock, invokes the
configured executor (default `claude -p` with the `ctx-dream` skill), and
fails **loud** (writing `dreams/.failed`) if the executor is missing or
errors — it never silently no-ops.

```bash
ctx dream [flags]
ctx dream <subcommand>
```

**Flags**:

| Flag       | Description                                                       |
|------------|-------------------------------------------------------------------|
| `--mode`   | Pass mode (`discipline`; default from `.ctxrc dream.mode`)        |
| `--max`    | Max `ideas/` files processed this pass (default `dream.max`)      |
| `--budget` | Step/token budget for the pass (default `dream.budget`)           |
| `--force`  | Bypass the trigger gate (opt-in + cadence + quiet window)         |

**Examples**:

```bash
ctx dream
ctx dream --max 20 --force
```

### `ctx dream review`

List the pending proposals from the latest pass — those not yet decided
in the ledger — rendered substance-forward (summary, status, action,
evidence, confidence, rationale). This is the read side of the
`/ctx-serendipity` garden walk.

```bash
ctx dream review
```

### `ctx dream accept <id>`

Accept a proposal's recommended action. Mechanical actions (`archive`,
`mark-blog`, `keep`) apply immediately with both guards enforced and a
ledger entry recorded; generative actions (`promote`, `merge`) record
accepted intent and are completed from the full source via
`/ctx-serendipity`.

**Arguments**:

- `id`: the proposal ID (from `ctx dream review`)

**Flags**:

| Flag     | Description                          |
|----------|--------------------------------------|
| `--note` | Optional human note recorded in the ledger |

**Examples**:

```bash
ctx dream accept a1b2c3
ctx dream accept a1b2c3 --note "good catch"
```

### `ctx dream reject <id>`

Record a rejection. No mutation occurs; the proposal is not re-surfaced
unless its source idea changes (dedup-against-seen).

**Arguments**:

- `id`: the proposal ID

**Flags**:

| Flag     | Description                          |
|----------|--------------------------------------|
| `--note` | Optional human note recorded in the ledger |

**Examples**:

```bash
ctx dream reject a1b2c3
ctx dream reject a1b2c3 --note "still relevant"
```

### `ctx dream amend <id> --action <action>`

Apply a different action than the one proposed, recording the decision as
amended (original provenance preserved).

**Arguments**:

- `id`: the proposal ID

**Flags**:

| Flag       | Description                                                  |
|------------|-------------------------------------------------------------|
| `--action` | The action to apply instead (`archive`/`merge`/`promote`/`mark-blog`/`keep`) |
| `--note`   | Optional human note recorded in the ledger                  |

**Examples**:

```bash
ctx dream amend a1b2c3 --action keep
ctx dream amend a1b2c3 --action archive --note "superseded"
```

**See also**: [Run the Dream recipe](../recipes/run-the-dream.md) ·
[Executor contract](../reference/dream-executor-contract.md).
