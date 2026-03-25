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

| Flag                   | Description                                               |
|------------------------|-----------------------------------------------------------|
| `--help`               | Show command help                                         |
| `--version`            | Show version                                              |
| `--context-dir <path>` | Override context directory (default: `.context/`)         |
| `--allow-outside-cwd`  | Allow context directory outside current working directory |

**Initialization required.** Most commands require a `.context/` directory
created by `ctx init`. Running a command without one produces:

```
ctx: not initialized - run "ctx init" first
```

Commands that work before initialization: `ctx init`, `ctx hook`,
`ctx doctor`, and grouping commands that only show help (e.g. `ctx`
with no subcommand, `ctx system`). Hidden hook commands have their
own guards and no-op gracefully.

<!-- drift-check: ctx --help | grep -c '  [a-z]' -->
## Commands

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx init`](init-status.md#ctx-init)         | Initialize `.context/` directory with templates          |
| [`ctx status`](init-status.md#ctx-status)     | Show context summary (files, tokens, drift)              |
| [`ctx agent`](init-status.md#ctx-agent)       | Print token-budgeted context packet for AI consumption   |
| [`ctx load`](init-status.md#ctx-load)         | Output assembled context in read order                   |
| [`ctx add`](context.md#ctx-add)               | Add a task, decision, learning, or convention            |
| [`ctx drift`](context.md#ctx-drift)           | Detect stale paths, secrets, missing files               |
| [`ctx sync`](context.md#ctx-sync)             | Reconcile context with codebase state                    |
| [`ctx compact`](context.md#ctx-compact)       | Archive completed tasks, clean up files                  |
| [`ctx task`](context.md#ctx-task)             | Task completion, archival, and snapshots                 |
| [`ctx permission`](context.md#ctx-permission) | Permission snapshots (golden image)                      |
| [`ctx reindex`](context.md#ctx-reindex)       | Regenerate indices for `DECISIONS.md` and `LEARNINGS.md` |
| [`ctx decision`](context.md#ctx-decision)     | Manage `DECISIONS.md` (reindex)                          |
| [`ctx learning`](context.md#ctx-learning)     | Manage `LEARNINGS.md` (reindex)                          |
| [`ctx recall`](recall.md#ctx-recall)          | Browse and export AI session history                     |
| [`ctx journal`](recall.md#ctx-journal)        | Generate static site from journal entries                |
| [`ctx serve`](recall.md#ctx-serve)            | Serve any zensical directory (default: journal site)     |
| [`ctx watch`](tools.md#ctx-watch)             | Auto-apply context updates from AI output                |
| [`ctx hook`](tools.md#ctx-hook)               | Generate AI tool integration configs                     |
| [`ctx loop`](tools.md#ctx-loop)               | Generate autonomous loop script                          |
| [`ctx memory`](tools.md#ctx-memory)           | Bridge Claude Code auto memory into .context/            |
| [`ctx notify`](tools.md#ctx-notify)           | Send webhook notifications                               |
| [`ctx change`](tools.md#ctx-change)           | Show what changed since last session                     |
| [`ctx dep`](tools.md#ctx-dep)                 | Show package dependency graph                            |
| [`ctx pad`](tools.md#ctx-pad)                 | Encrypted scratchpad for sensitive one-liners            |

| [`ctx remind`](tools.md#ctx-remind)           | Session-scoped reminders that surface at session start   |
| [`ctx completion`](tools.md#ctx-completion)   | Generate shell autocompletion scripts                    |
| [`ctx guide`](tools.md#ctx-guide)             | Quick-reference cheat sheet                              |
| [`ctx why`](tools.md#ctx-why)                 | Read the philosophy behind ctx                           |
| [`ctx site`](tools.md#ctx-site)               | Site management (feed generation)                        |
| [`ctx doctor`](doctor.md#ctx-doctor)          | Structural health check (hooks, drift, config)           |
| [`ctx mcp`](mcp.md#ctx-mcp)                   | MCP server for AI tool integration (stdin/stdout)        |
| [`ctx config`](config.md#ctx-config)          | Manage runtime configuration profiles                    |
| [`ctx system`](system.md#ctx-system)          | System diagnostics and hook commands                     |

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

| Variable                | Description                                         |
|-------------------------|-----------------------------------------------------|
| `CTX_DIR`               | Override default context directory path             |
| `CTX_TOKEN_BUDGET`      | Override default token budget                       |
| `CTX_BACKUP_SMB_URL`    | SMB share URL for backups (e.g. `smb://host/share`) |
| `CTX_BACKUP_SMB_SUBDIR` | Subdirectory on SMB share (default: `ctx-sessions`) |

<!-- drift-check: diff <(grep 'yaml:' internal/rc/types.go | grep -oP '"[a-z_]+"' | tr -d '"' | sort) <(sed -n '/Configuration File/,/^##[^#]/p' docs/cli/index.md | grep -oP '`[a-z_]+`' | tr -d '`' | sort -u) -->
## Configuration File

Optional `.ctxrc` (*YAML format*) at project root:

```yaml
# .ctxrc
context_dir: .context        # Context directory name
token_budget: 8000           # Default token budget
priority_order:              # File loading priority
  - TASKS.md
  - DECISIONS.md
  - CONVENTIONS.md
auto_archive: true           # Auto-archive old items
archive_after_days: 7        # Days before archiving tasks
scratchpad_encrypt: true     # Encrypt scratchpad (default: true)
allow_outside_cwd: false     # Skip boundary check (default: false)
event_log: false             # Enable local hook event logging
entry_count_learnings: 30    # Drift warning threshold (0 = disable)
entry_count_decisions: 20    # Drift warning threshold (0 = disable)
convention_line_count: 200   # Line count warning for CONVENTIONS.md (0 = disable)
injection_token_warn: 15000  # Oversize injection warning (0 = disable)
context_window: 200000       # Auto-detected for Claude Code; override for other tools
billing_token_warn: 0        # One-shot billing warning at this token count (0 = disabled)
key_rotation_days: 90        # Days before key rotation nudge
session_prefixes:            # Recognized session header prefixes (extend for i18n)
  - "Session:"               # English (default)
  # - "Oturum:"              # Turkish (add as needed)
  # - "セッション:"             # Japanese (add as needed)
freshness_files:             # Files with technology-dependent constants (opt-in)
  - path: config/thresholds.yaml
    desc: Model token limits and batch sizes
    review_url: https://docs.example.com/limits  # Optional
notify:                      # Webhook notification settings
  events:                    # Required: only listed events fire
    - loop
    - nudge
    - relay
    # - heartbeat            # Every-prompt session-alive signal
```

| Field                   | Type       | Default        | Description                                                                                                    |
|-------------------------|------------|----------------|----------------------------------------------------------------------------------------------------------------|
| `context_dir`           | `string`   | `.context`     | Context directory name (relative to project root)                                                              |
| `token_budget`          | `int`      | `8000`         | Default token budget for `ctx agent`                                                                           |
| `priority_order`        | `[]string` | *(all files)*  | File loading priority for context packets                                                                      |
| `auto_archive`          | `bool`     | `true`         | Auto-archive completed tasks                                                                                   |
| `archive_after_days`    | `int`      | `7`            | Days before completed tasks are archived                                                                       |
| `scratchpad_encrypt`    | `bool`     | `true`         | Encrypt scratchpad with AES-256-GCM                                                                            |
| `allow_outside_cwd`     | `bool`     | `false`        | Skip boundary check for external context dirs                                                                  |
| `event_log`             | `bool`     | `false`        | Enable local hook event logging to `.context/state/events.jsonl`                                               |
| `entry_count_learnings` | `int`      | `30`           | Drift warning when `LEARNINGS.md` exceeds this count                                                           |
| `entry_count_decisions` | `int`      | `20`           | Drift warning when `DECISIONS.md` exceeds this count                                                           |
| `convention_line_count` | `int`      | `200`          | Line count warning for `CONVENTIONS.md`                                                                        |
| `injection_token_warn`  | `int`      | `15000`        | Warn when auto-injected context exceeds this token count (0 = disable)                                         |
| `context_window`        | `int`      | `200000`       | Context window size in tokens. Auto-detected for Claude Code (200k/1M); override for other AI tools            |
| `billing_token_warn`    | `int`      | `0` *(off)*    | One-shot warning when session tokens exceed this threshold (0 = disabled)                                      |
| `key_rotation_days`     | `int`      | `90`           | Days before encryption key rotation nudge                                                                      |
| `session_prefixes`      | `[]string` | `["Session:"]` | Recognized Markdown session header prefixes. Extend to parse sessions written in other languages               |
| `freshness_files`       | `[]object` | *(none)*       | Files to track for staleness (path, desc, optional review_url). Hook warns after 6 months without modification |
| `notify.events`         | `[]string` | *(all)*        | Event filter for webhook notifications (empty = all)                                                           |

**Priority order:** CLI flags > Environment variables > `.ctxrc` > Defaults

All settings are optional. Missing values use defaults.
