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

| Flag            | Description                                                |
|-----------------|------------------------------------------------------------|
| `--help`        | Show command help                                          |
| `--version`     | Show version                                               |
| `--tool <name>` | Override active AI tool identifier (e.g. `kiro`, `cursor`) |

**Tell ctx which `.context/` to use.** ctx does not search the
filesystem for `.context/`: you have to declare it. Three ways:

- `eval "$(ctx activate)"` (recommended): binds `CTX_DIR` for the
  current shell.
- `export CTX_DIR=/abs/path/to/.context` directly, then run any
  `ctx` command.
- `CTX_DIR=/abs/path/to/.context ctx <command>` inline, for a
  one-shot or CI step.

`CTX_DIR` must be an absolute path with `.context` as its basename.
Relative paths and other names are rejected on first use; the
basename guard catches the common footgun
(`export CTX_DIR=$(pwd)`) before stray writes can leak to the
project root.

If you forget, commands fail fast with a linkable
`Error: no context directory specified` pointing at
[Activating a Context Directory](../recipes/activating-context.md).
A handful of commands run without a declaration because they don't
need a project: `ctx init`, `ctx activate`, `ctx deactivate`,
`ctx version`, `ctx help`, `ctx system bootstrap`, `ctx doctor`,
`ctx guide`, `ctx why`, `ctx config switch/status`, and
`ctx hub *`.

**Initialization required.** Once declared, the target must already
have been initialized by `ctx init` (otherwise commands return
`ctx: not initialized`).

<!-- drift-check: ctx --help | grep -c '  [a-z]' -->
## Getting Started

| Command                                             | Description                                              |
|-----------------------------------------------------|----------------------------------------------------------|
| [`ctx init`](init-status.md#ctx-init)               | Initialize `.context/` directory with templates          |
| [`ctx activate`](init-status.md#ctx-activate)       | Emit `export CTX_DIR=...` to bind context for the shell  |
| [`ctx deactivate`](init-status.md#ctx-deactivate)   | Emit `unset CTX_DIR` to clear the binding                |
| [`ctx status`](init-status.md#ctx-status)           | Show context summary (files, tokens, drift)              |
| [`ctx guide`](guide.md#ctx-guide)                   | Quick-reference cheat sheet                              |
| [`ctx why`](why.md#ctx-why)                         | Read the philosophy behind `ctx`                         |

## Context

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx load`](context.md#ctx-load)             | Output assembled context in read order                   |
| [`ctx agent`](context.md#ctx-agent)           | Print token-budgeted context packet for AI consumption   |
| [`ctx sync`](context.md#ctx-sync)             | Reconcile context with codebase state                    |
| [`ctx drift`](context.md#ctx-drift)           | Detect stale paths, secrets, missing files               |
| [`ctx compact`](context.md#ctx-compact)       | Archive completed tasks, clean up files                  |
| [`ctx fmt`](context.md#ctx-fmt)               | Format context files to 80-char line width               |
| [`ctx task`](context.md#ctx-task)             | Add tasks, mark complete, archive, snapshot              |
| [`ctx decision`](context.md#ctx-decision)     | Add decisions and reindex `DECISIONS.md`                 |
| [`ctx learning`](context.md#ctx-learning)     | Add learnings and reindex `LEARNINGS.md`                 |
| [`ctx convention`](context.md#ctx-convention) | Add conventions to `CONVENTIONS.md`                      |
| [`ctx reindex`](context.md#ctx-reindex)       | Regenerate indices for `DECISIONS.md` and `LEARNINGS.md` |
| [`ctx permission`](context.md#ctx-permission) | Permission snapshots (golden image)                      |
| [`ctx change`](change.md#ctx-change)          | Show what changed since last session                     |
| [`ctx memory`](memory.md#ctx-memory)          | Bridge Claude Code auto memory into `.context/`          |
| [`ctx watch`](watch.md#ctx-watch)             | Auto-apply context updates from AI output                |

## Sessions

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx journal`](journal.md#ctx-journal)       | Browse, import, enrich, and lock session history         |
| [`ctx pad`](pad.md#ctx-pad)                   | Encrypted scratchpad for sensitive one-liners            |
| [`ctx remind`](remind.md#ctx-remind)          | Session-scoped reminders that surface at session start   |
| [`ctx hook pause`](pause.md)                  | Pause context hooks for the current session              |
| [`ctx hook resume`](resume.md)                | Resume paused context hooks                              |

## Integrations

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx setup`](setup.md#ctx-setup)             | Generate AI tool integration configs                     |
| [`ctx steering`](steering.md#ctx-steering)    | Manage steering files (behavioral rules for AI tools)    |
| [`ctx trigger`](trigger.md#ctx-trigger)       | Manage lifecycle triggers (scripts for automation)       |
| [`ctx skill`](skill.md#ctx-skill)             | Manage reusable instruction bundles                      |
| [`ctx mcp`](mcp.md#ctx-mcp)                   | MCP server for AI tool integration (stdin/stdout)        |
| [`ctx hook notify`](notify.md)                | Webhook notifications (setup, test, send)                |
| [`ctx loop`](loop.md#ctx-loop)                | Generate autonomous loop script                          |
| [`ctx connection`](connection.md#ctx-connection) | Client-side commands for connecting to a `ctx` Hub    |
| [`ctx hub`](hub.md#ctx-hub)                   | Operate a `ctx` Hub server or cluster                    |
| [`ctx serve`](serve.md#ctx-serve)             | Serve a static site locally via zensical                 |
| [`ctx site`](site.md#ctx-site)                | Site management (feed generation)                        |

## Diagnostics

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx doctor`](doctor.md#ctx-doctor)          | Structural health check (hooks, drift, config)           |
| [`ctx trace`](trace.md#ctx-trace)             | Show context behind git commits                          |
| [`ctx sysinfo`](sysinfo.md#ctx-sysinfo)       | Show system resource usage (memory, swap, disk, load)    |
| [`ctx usage`](usage.md#ctx-usage)             | Show session token usage stats                           |

## Runtime

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx config`](config.md#ctx-config)          | Manage runtime configuration profiles                    |
| [`ctx prune`](prune.md#ctx-prune)             | Clean stale per-session state files                      |
| [`ctx hook`](hook.md#ctx-hook)                | Hook message, notification, and lifecycle controls       |
| [`ctx system`](system.md#ctx-system)          | Hook plumbing and agent-only commands (not user-facing)  |

## Shell

| Command                                       | Description                                              |
|-----------------------------------------------|----------------------------------------------------------|
| [`ctx completion`](completion.md#ctx-completion) | Generate shell autocompletion scripts                 |

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
| `CTX_SESSION_ID`        | Active AI session ID (used by `ctx trace` for context linking) |

<!-- drift-check: diff <(grep 'yaml:' internal/rc/types.go | grep -oP '"[a-z_]+"' | tr -d '"' | sort) <(sed -n '/Configuration File/,/^##[^#]/p' docs/cli/index.md | grep -oP '`[a-z_]+`' | tr -d '`' | sort -u) -->
## Configuration File

Optional `.ctxrc` (*YAML format*) at project root:

```yaml
# .ctxrc
token_budget: 8000           # Default token budget
priority_order:              # File loading priority
  - TASKS.md
  - DECISIONS.md
  - CONVENTIONS.md
auto_archive: true           # Auto-archive old items
archive_after_days: 7        # Days before archiving tasks
scratchpad_encrypt: true     # Encrypt scratchpad (default: true)
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
| `token_budget`          | `int`      | `8000`         | Default token budget for `ctx agent`                                                                           |
| `priority_order`        | `[]string` | *(all files)*  | File loading priority for context packets                                                                      |
| `auto_archive`          | `bool`     | `true`         | Auto-archive completed tasks                                                                                   |
| `archive_after_days`    | `int`      | `7`            | Days before completed tasks are archived                                                                       |
| `scratchpad_encrypt`    | `bool`     | `true`         | Encrypt scratchpad with AES-256-GCM                                                                            |
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
