---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: System
icon: lucide/settings
---

![ctx](../images/ctx-banner.png)


### `ctx system`

Hidden parent command that hosts Claude Code hook plumbing and a small
set of session-lifecycle plumbing subcommands used by skills and editor
integrations. The parent is registered without a visible group in
`ctx --help`; run `ctx system --help` to see its subcommands.

```bash
ctx system <subcommand>
```

!!! note "Commands previously under `ctx system`"
    Several user-facing maintenance commands used to live under
    `ctx system` and were promoted to top-level:

    - `ctx system backup` → **`ctx backup`**
    - `ctx system events` → **`ctx hook event`**
    - `ctx system message` → **`ctx hook message`**
    - `ctx system prune` → **`ctx prune`**
    - `ctx system resources` → **`ctx sysinfo`**
    - `ctx system stats` → **`ctx usage`**

    `ctx system bootstrap` remains under `ctx system` as a hidden,
    agent-only command. Update any scripts or personal docs that
    reference the old paths.

## Plumbing subcommands

These are not hook handlers — they're called by skills and editor
integrations during the session lifecycle. Safe to run manually.

#### `ctx system mark-journal`

Update processing state for a journal entry. Records the current date
in `.context/journal/.state.json`. Used by journal skills to record
pipeline progress.

```bash
ctx system mark-journal <filename> <stage>
```

**Stages**: `exported`, `enriched`, `normalized`, `fences_verified`

| Flag      | Description                           |
|-----------|---------------------------------------|
| `--check` | Check if stage is set (exit 1 if not) |

**Example**:

```bash
ctx system mark-journal 2026-01-21-session-abc12345.md enriched
ctx system mark-journal 2026-01-21-session-abc12345.md normalized
ctx system mark-journal --check 2026-01-21-session-abc12345.md fences_verified
```

#### `ctx system mark-wrapped-up`

Suppress context checkpoint nudges after a wrap-up ceremony. Writes a
marker file that `check-context-size` checks before emitting checkpoint
boxes. The marker expires after 2 hours.

Called automatically by `/ctx-wrap-up` after persisting context
(*not intended for direct use*).

```bash
ctx system mark-wrapped-up
```

No flags, no arguments. Idempotent: running it again updates the
marker timestamp.

#### `ctx system pause` / `ctx system resume`

Session-scoped hook suppression. `ctx system pause` writes a marker
file that causes hook plumbing to no-op for the current session;
`ctx system resume` removes it. These are the hook-plumbing
counterparts to the `ctx hook pause` / `ctx hook resume` commands
(which call them internally).

Read the session ID from stdin JSON (same as hooks) or pass
`--session-id`.

#### `ctx system session-event`

Records a session lifecycle event (start or end) to the event log.
Called by editor integrations when a workspace is opened or closed.

```bash
ctx system session-event --type start --caller vscode
ctx system session-event --type end --caller vscode
```

## Hook subcommands

Hidden Claude Code hook handlers implementing the hook contract: read
JSON from stdin, perform logic, emit output on stdout, exit 0. Block
commands output JSON with a `decision` field.

UserPromptSubmit hooks: `context-load-gate`, `check-context-size`,
`check-persistence`, `check-ceremony`, `check-journal`, `check-version`,
`check-resource`, `check-knowledge`, `check-map-staleness`,
`check-memory-drift`, `check-reminder`, `check-freshness`,
`check-hub-sync`, `check-backup-age`, `check-skill-discovery`,
`heartbeat`.

PreToolUse hooks: `block-non-path-ctx`, `block-dangerous-command`,
`qa-reminder`, `specs-nudge`.

PostToolUse hooks: `post-commit`, `check-task-completion`.

See [AI Tools](../operations/integrations.md#plugin-hooks) for
registration details and the Claude Code plugin integration.
