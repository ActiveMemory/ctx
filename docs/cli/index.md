---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: CLI
icon: lucide/terminal
---

![ctx](../images/ctx-banner.png)

## `ctx` CLI

This is a complete reference for all `ctx` commands.

## Global Options

All commands support these flags:

| Flag                   | Description                                              |
|------------------------|----------------------------------------------------------|
| `--help`               | Show command help                                        |
| `--version`            | Show version                                             |
| `--context-dir <path>` | Override context directory (default: `.context/`)        |
| `--no-color`           | Disable colored output                                   |
| `--allow-outside-cwd`  | Allow context directory outside current working directory |

> The `NO_COLOR=1` environment variable also disables colored output.

<!-- drift-check: ctx --help | grep -c '  [a-z]' -->
## Commands

| Command                           | Description                                               |
|-----------------------------------|-----------------------------------------------------------|
| [`ctx init`](init-status.md#ctx-init)           | Initialize `.context/` directory with templates           |
| [`ctx status`](init-status.md#ctx-status)       | Show context summary (files, tokens, drift)               |
| [`ctx agent`](init-status.md#ctx-agent)         | Print token-budgeted context packet for AI consumption    |
| [`ctx load`](init-status.md#ctx-load)           | Output assembled context in read order                    |
| [`ctx add`](context.md#ctx-add)             | Add a task, decision, learning, or convention             |
| [`ctx complete`](context.md#ctx-complete)   | Mark a task as done                                       |
| [`ctx drift`](context.md#ctx-drift)         | Detect stale paths, secrets, missing files                |
| [`ctx sync`](context.md#ctx-sync)           | Reconcile context with codebase state                     |
| [`ctx compact`](context.md#ctx-compact)     | Archive completed tasks, clean up files                   |
| [`ctx tasks`](context.md#ctx-tasks)         | Task archival and snapshots                               |
| [`ctx permissions`](context.md#ctx-permissions) | Permission snapshots (golden image)                   |
| [`ctx decisions`](context.md#ctx-decisions) | Manage `DECISIONS.md` (reindex)                           |
| [`ctx learnings`](context.md#ctx-learnings) | Manage `LEARNINGS.md` (reindex)                           |
| [`ctx recall`](recall.md#ctx-recall)       | Browse and export AI session history                      |
| [`ctx journal`](recall.md#ctx-journal)     | Generate static site from journal entries                 |
| [`ctx serve`](recall.md#ctx-serve)         | Serve any zensical directory (default: journal site)      |
| [`ctx watch`](tools.md#ctx-watch)          | Auto-apply context updates from AI output                 |
| [`ctx hook`](tools.md#ctx-hook)            | Generate AI tool integration configs                      |
| [`ctx loop`](tools.md#ctx-loop)            | Generate autonomous loop script                           |
| [`ctx notify`](tools.md#ctx-notify)        | Send webhook notifications                                |
| [`ctx pad`](tools.md#ctx-pad)              | Encrypted scratchpad for sensitive one-liners             |
| [`ctx remind`](tools.md#ctx-remind)        | Session-scoped reminders that surface at session start    |
| [`ctx completion`](tools.md#ctx-completion) | Generate shell autocompletion scripts                   |
| [`ctx why`](tools.md#ctx-why)              | Read the philosophy behind ctx                           |
| [`ctx system`](system.md#ctx-system)       | System diagnostics and hook commands                     |

---

## Exit Codes

| Code | Meaning                                |
|------|----------------------------------------|
| 0    | Success                                |
| 1    | General error / warnings (e.g. drift)  |
| 2    | Context not found                      |
| 3    | Violations found (e.g. drift)          |
| 4    | File operation error                   |

## Environment Variables

| Variable           | Description                             |
|--------------------|-----------------------------------------|
| `CTX_DIR`          | Override default context directory path |
| `CTX_TOKEN_BUDGET` | Override default token budget           |
| `NO_COLOR`         | Disable colored output when set         |

## Configuration File

Optional `.ctxrc` (*YAML format*) at project root:

```yaml
# .ctxrc
context_dir: .context                # Context directory name
token_budget: 8000                   # Default token budget
priority_order:                      # File loading priority
  - TASKS.md
  - DECISIONS.md
  - CONVENTIONS.md
auto_archive: true                   # Auto-archive old items
archive_after_days: 7                # Days before archiving tasks
scratchpad_encrypt: true             # Encrypt scratchpad (default: true)
allow_outside_cwd: false             # Skip boundary check (default: false)
entry_count_learnings: 30            # Drift warning threshold (0 = disable)
entry_count_decisions: 20            # Drift warning threshold (0 = disable)
convention_line_count: 200           # Line count warning for CONVENTIONS.md (0 = disable)
notify:                              # Webhook notification settings
  events:                            # Required — only listed events fire
    - loop
    - nudge
    - relay
  key_rotation_days: 90              # Days before key rotation nudge
```

| Field                           | Type       | Default      | Description                                          |
|---------------------------------|------------|--------------|------------------------------------------------------|
| `context_dir`                   | `string`   | `.context`   | Context directory name (relative to project root)    |
| `token_budget`                  | `int`      | `8000`       | Default token budget for `ctx agent`                 |
| `priority_order`                | `[]string` | *(all files)* | File loading priority for context packets           |
| `auto_archive`                  | `bool`     | `false`      | Auto-archive completed tasks                         |
| `archive_after_days`            | `int`      | `7`          | Days before completed tasks are archived             |
| `scratchpad_encrypt`            | `bool`     | `true`       | Encrypt scratchpad with AES-256-GCM                  |
| `allow_outside_cwd`             | `bool`     | `false`      | Skip boundary check for external context dirs        |
| `entry_count_learnings`         | `int`      | `30`         | Drift warning when `LEARNINGS.md` exceeds this count   |
| `entry_count_decisions`         | `int`      | `20`         | Drift warning when `DECISIONS.md` exceeds this count   |
| `convention_line_count`         | `int`      | `200`        | Line count warning for `CONVENTIONS.md`                |
| `notify.events`                 | `[]string` | *(all)*      | Event filter for webhook notifications (empty = all) |
| `notify.key_rotation_days`      | `int`      | `90`         | Days before encryption key rotation nudge            |

**Priority order:** CLI flags > Environment variables > `.ctxrc` > Defaults

All settings are optional. Missing values use defaults.
