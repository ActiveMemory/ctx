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

Complete reference for all `ctx` commands, grouped by function.

## Global Options

All commands support these flags:

| Flag                   | Description                                               |
|------------------------|-----------------------------------------------------------|
| `--help`               | Show command help                                         |
| `--version`            | Show version                                              |
| `--context-dir <path>` | Override context directory (default: `.context/`)         |
| `--allow-outside-cwd`  | Allow context directory outside current working directory |
| `--tool <name>`        | Override active AI tool identifier (e.g. `kiro`, `cursor`) |

**Initialization required.** Most commands require a `.context/` directory
created by `ctx init`. Running a command without one produces:

```
ctx: not initialized - run "ctx init" first
```

Commands that work before initialization: `ctx init`, `ctx setup`,
`ctx doctor`, and grouping commands that only show help.

<!-- drift-check: ctx --help | grep -c '  [a-z]' -->
## Getting Started

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx init`](init-status.md#ctx-init)         | Initialize `.context/` directory with templates          |
| [`ctx status`](init-status.md#ctx-status)     | Show context summary (files, tokens, drift)              |
| [`ctx guide`](guide.md#ctx-guide)             | Quick-reference cheat sheet                              |

## Context

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx add`](context.md#ctx-add)               | Add a task, decision, learning, or convention            |
| [`ctx load`](init-status.md#ctx-load)         | Output assembled context in read order                   |
| [`ctx agent`](init-status.md#ctx-agent)       | Print token-budgeted context packet for AI consumption   |
| [`ctx skill`](skill.md#ctx-skill)             | Manage reusable instruction bundles                      |
| [`ctx sync`](context.md#ctx-sync)             | Reconcile context with codebase state                    |
| [`ctx drift`](context.md#ctx-drift)           | Detect stale paths, secrets, missing files               |
| [`ctx compact`](context.md#ctx-compact)       | Archive completed tasks, clean up files                  |
| [`ctx fmt`](context.md#ctx-fmt)               | Format context files to 80-char line width               |

## Artifacts

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx decision`](context.md#ctx-decision)     | Manage `DECISIONS.md` (reindex)                          |
| [`ctx learning`](context.md#ctx-learning)     | Manage `LEARNINGS.md` (reindex)                          |
| [`ctx task`](context.md#ctx-task)             | Task completion, archival, and snapshots                 |
| [`ctx reindex`](context.md#ctx-reindex)       | Regenerate indices for `DECISIONS.md` and `LEARNINGS.md` |

## Sessions

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx journal`](journal.md#ctx-journal)       | Browse, import, enrich, and lock session history         |
| [`ctx memory`](memory.md#ctx-memory)          | Bridge Claude Code auto memory into `.context/`          |
| [`ctx remind`](remind.md#ctx-remind)          | Session-scoped reminders that surface at session start   |
| [`ctx pad`](pad.md#ctx-pad)                   | Encrypted scratchpad for sensitive one-liners            |

## Integrations

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx setup`](setup.md#ctx-setup)             | Generate AI tool integration configs                     |
| [`ctx steering`](steering.md#ctx-steering)    | Manage steering files (behavioral rules for AI tools)    |
| [`ctx trigger`](trigger.md#ctx-trigger)       | Manage lifecycle triggers (scripts for automation)       |
| [`ctx serve`](serve.md#ctx-serve)             | Serve a static site locally via zensical                 |
| [`ctx hub`](hub.md#ctx-hub)                   | Operate a ctx Hub server or cluster                      |
| [`ctx connection`](connection.md#ctx-connection) | Connect to a ctx Hub                                  |
| [`ctx mcp`](mcp.md#ctx-mcp)                  | MCP server for AI tool integration (stdin/stdout)        |
| [`ctx watch`](watch.md#ctx-watch)             | Auto-apply context updates from AI output                |
| [`ctx loop`](loop.md#ctx-loop)                | Generate autonomous loop script                          |

## Diagnostics

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx doctor`](doctor.md#ctx-doctor)          | Structural health check (hooks, drift, config)           |
| [`ctx change`](change.md#ctx-change)          | Show what changed since last session                     |
| [`ctx why`](why.md#ctx-why)                   | Read the philosophy behind `ctx`                         |
| [`ctx trace`](trace.md#ctx-trace)             | Show context behind git commits                          |
| [`ctx sysinfo`](sysinfo.md#ctx-sysinfo)      | Show system resource usage (memory, swap, disk, load)    |
| [`ctx usage`](usage.md#ctx-usage)             | Show session token usage stats                           |

## Runtime

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx config`](config.md#ctx-config)          | Manage runtime configuration profiles                    |
| [`ctx permission`](context.md#ctx-permission) | Permission snapshots (golden image)                      |
| [`ctx hook`](hook.md#ctx-hook)                | Hook message, notification, and lifecycle controls       |
| [`ctx backup`](backup.md#ctx-backup)          | Back up context and Claude data to tar.gz / SMB          |
| [`ctx prune`](prune.md#ctx-prune)             | Clean stale per-session state files                      |

## Shell

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx completion`](completion.md#ctx-completion) | Generate shell autocompletion scripts                 |

## Hidden

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx site`](site.md#ctx-site)                | Site management (feed generation)                        |
| [`ctx system`](system.md#ctx-system)          | Hook plumbing and agent-only commands (not user-facing)  |

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
| `CTX_SESSION_ID`        | Active AI session ID (used by `ctx trace` for context linking) |

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
companion_check: true        # Check companion tools at session start
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
tool: ""                     # Active AI tool: claude, cursor, cline, kiro, codex
steering:                    # Steering layer configuration
  dir: .context/steering     # Steering files directory
  default_inclusion: manual  # Default inclusion mode (always, auto, manual)
  default_tools: []          # Default tool filter for new steering files
hooks:                       # Hook system configuration
  dir: .context/hooks        # Hook scripts directory
  timeout: 10                # Per-hook execution timeout in seconds
  enabled: true              # Whether hook execution is enabled
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
| `companion_check`       | `bool`     | `true`         | Check companion tool availability (Gemini Search, GitNexus) during `/ctx-remember`                             |
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
| `tool`                  | `string`   | *(empty)*      | Active AI tool identifier (`claude`, `cursor`, `cline`, `kiro`, `codex`)                                       |
| `steering.dir`          | `string`   | `.context/steering` | Steering files directory                                                                                  |
| `steering.default_inclusion` | `string` | `manual`    | Default inclusion mode for new steering files (`always`, `auto`, `manual`)                                     |
| `steering.default_tools` | `[]string` | *(all)*       | Default tool filter for new steering files (empty = all tools)                                                 |
| `hooks.dir`             | `string`   | `.context/hooks` | Hook scripts directory                                                                                      |
| `hooks.timeout`         | `int`      | `10`           | Per-hook execution timeout in seconds                                                                          |
| `hooks.enabled`         | `bool`     | `true`         | Whether hook execution is enabled                                                                              |

**Priority order:** CLI flags > Environment variables > `.ctxrc` > Defaults

All settings are optional. Missing values use defaults.
