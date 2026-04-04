---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Configuration
icon: lucide/settings
---

![ctx](../images/ctx-banner.png)

## Configuration

ctx uses three layers of configuration. Each layer overrides the one below it:

1. **CLI flags**: Per-invocation overrides (*highest priority*)
2. **Environment variables**: Shell or CI/CD overrides
3. **The `.ctxrc` file**: Project-level defaults (*YAML*)
4. **Built-in defaults**: Hardcoded fallbacks (*lowest priority*)

All settings are optional: If nothing is configured, `ctx` works out of the box
with sensible defaults.

---

## The `.ctxrc` File

The `.ctxrc` file is an optional YAML file placed in the **project root**
(*next to your `.context/` directory*). It lets you set project-level defaults
that apply to every `ctx` command.

### Location

```
my-project/
├── .ctxrc              ← configuration file
├── .context/
│   ├── TASKS.md
│   ├── DECISIONS.md
│   └── ...
└── src/
```

`ctx` looks for `.ctxrc` in the current working directory when any command runs.
There is no global or user-level config file: Configuration is always
per-project.

!!! note "Contributors: Dev Configuration Profile"
    The ctx repo ships two `.ctxrc` source profiles (`.ctxrc.base` and
    `.ctxrc.dev`). The working copy is gitignored and swapped between them
    via `ctx config switch dev` / `ctx config switch base`.
    See [Contributing: Configuration Profiles](contributing.md#configuration-profiles).

!!! tip "Using a Different .context Directory"
    The default `.context/` directory can be changed per-project via the
    `context_dir` key in `.ctxrc`, the `CTX_DIR` environment variable, or the
    `--context-dir` CLI flag.

    See [Environment Variables](#environment-variables)
    and [CLI Global Flags](#cli-global-flags) below for details.

<!-- drift-check: diff <(grep 'yaml:' internal/rc/types.go | grep -oP '"[a-z_]+"' | tr -d '"' | sort -u | grep -v 'desc\|events\|path\|review_url\|profile\|key_path') <(sed -n '/^# \.ctxrc:/,/^```$/p' docs/home/configuration.md | grep -oP '^# ([a-z_]+):' | sed 's/^# //;s/://' | sort -u) -->
### Full Reference

A commented `.ctxrc` showing all options and their defaults:

```yaml
# .ctxrc: ctx runtime configuration
# https://ctx.ist/configuration/
#
# All settings are optional. Missing values use defaults.
# Priority: CLI flags > environment variables > .ctxrc > defaults
#
# context_dir: .context
# token_budget: 8000
# auto_archive: true
# archive_after_days: 7
# scratchpad_encrypt: true
# allow_outside_cwd: false
# event_log: false
# entry_count_learnings: 30
# entry_count_decisions: 20
# convention_line_count: 200
# injection_token_warn: 15000
# context_window: 200000      # auto-detected for Claude Code; override for other tools
# billing_token_warn: 0       # one-shot warning at this token count (0 = disabled)
#
# stale_age_days: 30      # days before drift flags a context file as stale (0 = disabled)
# key_rotation_days: 90
# task_nudge_interval: 5   # Edit/Write calls between task completion nudges
#
# notify:               # requires: ctx notify setup
#   events:             # required: no events sent unless listed
#     - loop
#     - nudge
#     - relay
#
# tool: ""              # Active AI tool: claude, cursor, cline, kiro, codex
#
# steering:             # Steering layer configuration
#   dir: .context/steering
#   default_inclusion: manual
#   default_tools: []
#
# hooks:                # Hook system configuration
#   dir: .context/hooks
#   timeout: 10
#   enabled: true
#
# priority_order:
#   - CONSTITUTION.md
#   - TASKS.md
#   - CONVENTIONS.md
#   - ARCHITECTURE.md
#   - DECISIONS.md
#   - LEARNINGS.md
#   - GLOSSARY.md
#   - AGENT_PLAYBOOK.md
```

<!-- drift-check: diff <(grep 'yaml:' internal/rc/types.go | grep -oP '"[a-z_]+"' | tr -d '"' | sort -u | grep -v 'desc\|events\|path\|review_url\|profile\|key_path') <(sed -n '/Option Reference/,/^\*\*Default/p' docs/home/configuration.md | grep -oP '`([a-z_.]+)`' | tr -d '`' | sed 's/notify\.events/notify/' | sort -u | grep -v 'string\|int\|bool\|\[\]') -->
### Option Reference

| Option                  | Type       | Default       | Description                                                                                                                               |
|-------------------------|------------|---------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| `context_dir`           | `string`   | `.context`    | Context directory name (relative to project root)                                                                                         |
| `token_budget`          | `int`      | `8000`        | Default token budget for `ctx agent` and `ctx load`                                                                                       |
| `auto_archive`          | `bool`     | `true`        | Auto-archive completed tasks during `ctx compact`                                                                                         |
| `archive_after_days`    | `int`      | `7`           | Days before completed tasks are archived                                                                                                  |
| `scratchpad_encrypt`    | `bool`     | `true`        | Encrypt scratchpad with AES-256-GCM                                                                                                       |
| `allow_outside_cwd`     | `bool`     | `false`       | Allow context directory outside the current working directory                                                                             |
| `event_log`             | `bool`     | `false`       | Enable local hook event logging to `.context/state/events.jsonl`                                                                          |
| `entry_count_learnings` | `int`      | `30`          | Drift warning when `LEARNINGS.md` exceeds this entry count (0 = disable)                                                                  |
| `entry_count_decisions` | `int`      | `20`          | Drift warning when `DECISIONS.md` exceeds this entry count (0 = disable)                                                                  |
| `convention_line_count` | `int`      | `200`         | Drift warning when `CONVENTIONS.md` exceeds this line count (0 = disable)                                                                 |
| `injection_token_warn`  | `int`      | `15000`       | Warn when auto-injected context exceeds this token count (0 = disable)                                                                    |
| `context_window`        | `int`      | `200000`      | Context window size in tokens. Auto-detected for Claude Code (200k/1M); override for other AI tools                                       |
| `billing_token_warn`    | `int`      | `0` *(off)*   | One-shot warning when session tokens exceed this threshold (0 = disabled). For plans where tokens beyond an included allowance cost extra |
| `stale_age_days`        | `int`      | `30`          | Days before `ctx drift` flags a context file as stale (0 = disable)                                                                       |
| `key_rotation_days`     | `int`      | `90`          | Days before encryption key rotation nudge                                                                                                 |
| `task_nudge_interval`   | `int`      | `5`           | Edit/Write calls between task completion nudges                                                                                           |
| `notify.events`         | `[]string` | *(all)*       | Event filter for webhook notifications (empty = all)                                                                                      |
| `priority_order`        | `[]string` | *(see below)* | Custom file loading priority for context assembly                                                                                         |
| `tool`                  | `string`   | *(empty)*     | Active AI tool identifier (`claude`, `cursor`, `cline`, `kiro`, `codex`). Used by steering sync and hook dispatch                         |
| `steering.dir`          | `string`   | `.context/steering` | Steering files directory                                                                                                             |
| `steering.default_inclusion` | `string` | `manual` | Default inclusion mode for new steering files (`always`, `auto`, `manual`)                                                                |
| `steering.default_tools` | `[]string` | *(all)*  | Default tool filter for new steering files (empty = all tools)                                                                            |
| `hooks.dir`             | `string`   | `.context/hooks` | Hook scripts directory                                                                                                                |
| `hooks.timeout`         | `int`      | `10`          | Per-hook execution timeout in seconds                                                                                                     |
| `hooks.enabled`         | `bool`     | `true`        | Whether hook execution is enabled                                                                                                         |

**Default priority order** (*used when `priority_order` is not set*):

1. `CONSTITUTION.md`
2. `TASKS.md`
3. `CONVENTIONS.md`
4. `ARCHITECTURE.md`
5. `DECISIONS.md`
6. `LEARNINGS.md`
7. `GLOSSARY.md`
8. `AGENT_PLAYBOOK.md`

See [Context Files](context-files.md#read-order-rationale) for the rationale
behind this ordering.

---

<!-- drift-check: diff <(grep -oP 'os\.Getenv\("[A-Z_]+"\)' internal/rc/rc.go | grep -oP '"[A-Z_]+"' | tr -d '"' | sort) <(sed -n '/Environment Variables/,/^---$/p' docs/home/configuration.md | grep -oP '`CTX_[A-Z_]+`' | tr -d '`' | sort -u) -->
## Environment Variables

Environment variables override `.ctxrc` values but are overridden by CLI flags.

| Variable           | Description                                 | Equivalent `.ctxrc` key |
|--------------------|---------------------------------------------|-------------------------|
| `CTX_DIR`          | Override the context directory path         | `context_dir`           |
| `CTX_TOKEN_BUDGET` | Override the default token budget           | `token_budget`          |

### Examples

```bash
# Use a shared context directory
CTX_DIR=/shared/team-context ctx status

# Increase token budget for a single run
CTX_TOKEN_BUDGET=16000 ctx agent
```

---

<!-- drift-check: diff <(ctx --help 2>&1 | grep -oP '^\s+--[a-z-]+' | sed 's/^\s*//' | sort) <(sed -n '/CLI Global Flags/,/^---$/p' docs/home/configuration.md | grep -oP '`(--[a-z-]+)`' | tr -d '`' | sort -u) -->
## CLI Global Flags

CLI flags have the highest priority and override both environment variables and
`.ctxrc` settings. These flags are available on every `ctx` command.

| Flag                   | Description                                               |
|------------------------|-----------------------------------------------------------|
| `--context-dir <path>` | Override context directory (default: `.context/`)         |
| `--allow-outside-cwd`  | Allow context directory outside current working directory |
| `--tool <name>`        | Override active AI tool identifier (e.g. `kiro`, `cursor`) |
| `--version`            | Show version and exit                                     |
| `--help`               | Show command help and exit                                |

### Examples

```bash
# Point to a different context directory:
ctx status --context-dir /path/to/shared/.context

# Allow external context directory (skips boundary check):
ctx status --context-dir /mnt/nas/project-context --allow-outside-cwd
```

---

## Priority Order

When the same setting is configured in multiple layers, the highest-priority
layer wins:

```
CLI flags  >  Environment variables  >  .ctxrc  >  Built-in defaults
(highest)                                          (lowest)
```

**Example resolution for `context_dir`:**

| Layer              | Value              | Wins? |
|--------------------|--------------------|-------|
| `--context-dir`    | `/tmp/ctx`         | Yes   |
| `CTX_DIR`          | `/shared/context`  | No    |
| `.ctxrc`           | `.my-context`      | No    |
| Default            | `.context`         | No    |

The CLI flag `/tmp/ctx` is used because it has the highest priority.

If the CLI flag were absent, `CTX_DIR=/shared/context` would win. If neither
the flag nor the env var were set, the `.ctxrc` value `.my-context` would be
used. With nothing configured, the default `.context` applies.

---

## Examples

### External `.context` Directory

Store context outside the project tree (*useful for monorepos or shared context*):

```yaml
# .ctxrc
context_dir: /home/team/shared-context
allow_outside_cwd: true
```

### Custom Token Budget

Increase the token budget for projects with large context:

```yaml
# .ctxrc
token_budget: 16000
```

This affects the default budget for `ctx agent` and `ctx load`. You can still
override per-invocation with `ctx agent --budget 4000`.

### Disabled Scratchpad Encryption

Turn off encryption for the scratchpad (*useful in ephemeral environments
where key management is unnecessary*):

```yaml
# .ctxrc
scratchpad_encrypt: false
```

!!! danger "Unencrypted Scratchpads Store Secrets in Plaintext"
    Only disable encryption if you understand the security implications.

    The scratchpad may contain sensitive data such as API keys, database
    URLs, or deployment credentials.

### Custom Priority Order

Reorder context files to prioritize architecture over conventions:

```yaml
# .ctxrc
priority_order:
  - CONSTITUTION.md
  - TASKS.md
  - ARCHITECTURE.md
  - DECISIONS.md
  - CONVENTIONS.md
  - LEARNINGS.md
  - GLOSSARY.md
  - AGENT_PLAYBOOK.md
```

Files not listed in `priority_order` receive the lowest priority (100).
The order affects `ctx agent`, `ctx load`, and drift's file-priority
calculations.

### Billing Token Threshold

Get a one-shot warning when your session crosses a token threshold where
extra charges begin (*e.g., Claude Pro includes 200k tokens; beyond that
costs extra*):

```yaml
# .ctxrc
billing_token_warn: 180000   # warn before hitting the 200k paid boundary
```

The warning fires once per session the first time token usage exceeds
the threshold. Set to `0` (or omit) to disable.

### Adjusted Drift Thresholds

Raise or lower the entry-count thresholds that trigger drift warnings:

```yaml
# .ctxrc
entry_count_learnings: 50   # warn above 50 learnings (default: 30)
entry_count_decisions: 10   # warn above 10 decisions (default: 20)
convention_line_count: 300  # warn above 300 lines (default: 200)
```

Set any threshold to `0` to disable that specific check.

### Webhook Notifications

Get notified when loops complete, hooks fire, or agents reach milestones:

```bash
# Configure the webhook URL (encrypted, safe to commit)
ctx notify setup

# Test delivery
ctx notify test
```

Filter which events reach your webhook:

```yaml
# .ctxrc
notify:
  events:
    - loop      # loop completion/max-iteration
    - nudge     # VERBATIM relay hooks fired
    # - relay   # all hook output (verbose, for debugging)
    # - heartbeat  # every-prompt session-alive signal
```

Notifications are **opt-in**: No events are sent unless explicitly listed.

See [Webhook Notifications](../recipes/webhook-notifications.md) for a
step-by-step recipe.

---

## Hook Message Overrides

Hook messages control what text hooks emit when they fire. Each message
can be overridden per-project by placing a text file at the matching
path under `.context/`:

```
.context/hooks/messages/{hook}/{variant}.txt
```

The override takes priority over the embedded default compiled into the
`ctx` binary. An empty file silences the message while preserving the
hook's logic (counting, state tracking, cooldowns).

Use `ctx system message` to discover and manage overrides:

```bash
ctx system message list                      # see all messages
ctx system message show qa-reminder gate     # view the current template
ctx system message edit qa-reminder gate     # copy default for editing
ctx system message reset qa-reminder gate    # revert to default
```

See [Customizing Hook Messages](../recipes/customizing-hook-messages.md)
for detailed examples including Python, JavaScript, and silence
configurations.

---

## Agent Bootstrapping

AI agents need to know the resolved context directory at session start.
The `ctx system bootstrap` command prints the context path, file list, and
operating rules in both text and JSON formats:

```bash
ctx system bootstrap          # text output for agents
ctx system bootstrap -q       # just the context directory path
ctx system bootstrap --json   # structured output for automation
```

The `CLAUDE.md` template instructs the agent to run this as its first action.
Every nudge (*context checkpoint, persistence reminder, etc.*) also includes a
`Context: <dir>` footer that re-anchors the agent to the correct directory
throughout the session.

This replaces the previous approach of hardcoding `.context/` paths in agent
instructions. 

See [CLI Reference: bootstrap](../cli/system.md#ctx-system-bootstrap)
for full details.

---

**See also**: [CLI Reference](../cli/index.md) | [Context Files](context-files.md) | [Scratchpad](../reference/scratchpad.md)
