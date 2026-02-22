---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Configuration
icon: lucide/settings
---

![ctx](images/ctx-banner.png)

## Configuration

ctx uses three layers of configuration. Each layer overrides the one below it:

1. **CLI flags** — per-invocation overrides (highest priority)
2. **Environment variables** — shell or CI/CD overrides
3. **`.ctxrc` file** — project-level defaults (YAML)
4. **Built-in defaults** — hardcoded fallbacks (lowest priority)

All settings are optional. If nothing is configured, ctx works out of the box
with sensible defaults.

---

## The `.ctxrc` File

The `.ctxrc` file is an optional YAML file placed in the **project root**
(next to your `.context/` directory). It lets you set project-level defaults
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

ctx looks for `.ctxrc` in the current working directory when any command runs.
There is no global or user-level config file — configuration is always
per-project.

### Full Reference

A commented `.ctxrc` showing all options and their defaults:

```yaml
# .ctxrc — ctx runtime configuration
# https://ctx.ist/configuration/
#
# All settings are optional. Missing values use defaults.
# Priority: CLI flags > environment variables > .ctxrc > defaults

# context_dir: .context
# token_budget: 8000
# auto_archive: true
# archive_after_days: 7
# scratchpad_encrypt: true
# allow_outside_cwd: false
# entry_count_learnings: 30
# entry_count_decisions: 20
# convention_line_count: 200
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

### Option Reference

| Option                  | Type       | Default        | Description                                             |
|-------------------------|------------|----------------|---------------------------------------------------------|
| `context_dir`           | `string`   | `.context`     | Context directory name (relative to project root)       |
| `token_budget`          | `int`      | `8000`         | Default token budget for `ctx agent` and `ctx load`     |
| `auto_archive`          | `bool`     | `true`         | Auto-archive completed tasks during `ctx compact`       |
| `archive_after_days`    | `int`      | `7`            | Days before completed tasks are archived                |
| `scratchpad_encrypt`    | `bool`     | `true`         | Encrypt scratchpad with AES-256-GCM                     |
| `allow_outside_cwd`     | `bool`     | `false`        | Allow context directory outside the current working directory |
| `entry_count_learnings` | `int`      | `30`           | Drift warning when LEARNINGS.md exceeds this entry count (0 = disable) |
| `entry_count_decisions` | `int`      | `20`           | Drift warning when DECISIONS.md exceeds this entry count (0 = disable) |
| `convention_line_count` | `int`      | `200`          | Drift warning when CONVENTIONS.md exceeds this line count (0 = disable) |
| `priority_order`        | `[]string` | *(see below)*  | Custom file loading priority for context assembly       |

**Default priority order** (used when `priority_order` is not set):

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

## Environment Variables

Environment variables override `.ctxrc` values but are overridden by CLI flags.

| Variable           | Description                              | Equivalent `.ctxrc` key |
|--------------------|------------------------------------------|-------------------------|
| `CTX_DIR`          | Override the context directory path       | `context_dir`           |
| `CTX_TOKEN_BUDGET` | Override the default token budget         | `token_budget`          |
| `NO_COLOR`         | Disable colored output when set (any value) | *(CLI flag only)*    |

### Examples

```bash
# Use a shared context directory
CTX_DIR=/shared/team-context ctx status

# Increase token budget for a single run
CTX_TOKEN_BUDGET=16000 ctx agent

# Disable color in CI pipelines
NO_COLOR=1 ctx drift --json
```

!!! note "`NO_COLOR` follows the [no-color.org](https://no-color.org) convention"
    Setting `NO_COLOR` to any non-empty value disables colored terminal output.
    This is equivalent to the `--no-color` CLI flag.

---

## CLI Global Flags

CLI flags have the highest priority and override both environment variables and
`.ctxrc` settings. These flags are available on every `ctx` command.

| Flag                   | Description                                              |
|------------------------|----------------------------------------------------------|
| `--context-dir <path>` | Override context directory (default: `.context/`)        |
| `--allow-outside-cwd`  | Allow context directory outside current working directory |
| `--no-color`           | Disable colored output                                   |
| `--version`            | Show version and exit                                    |
| `--help`               | Show command help and exit                               |

### Examples

```bash
# Point to a different context directory
ctx status --context-dir /path/to/shared/.context

# Allow external context directory (skips boundary check)
ctx status --context-dir /mnt/nas/project-context --allow-outside-cwd
```

---

## Priority Order

When the same setting is configured in multiple layers, the highest-priority
layer wins:

```
CLI flags  >  Environment variables  >  .ctxrc  >  Built-in defaults
(highest)                                           (lowest)
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

### External context directory

Store context outside the project tree (useful for monorepos or shared context):

```yaml
# .ctxrc
context_dir: /home/team/shared-context
allow_outside_cwd: true
```

### Custom token budget

Increase the token budget for projects with large context:

```yaml
# .ctxrc
token_budget: 16000
```

This affects the default budget for `ctx agent` and `ctx load`. You can still
override per-invocation with `ctx agent --budget 4000`.

### Disabled scratchpad encryption

Turn off encryption for the scratchpad (useful in ephemeral environments
where key management is unnecessary):

```yaml
# .ctxrc
scratchpad_encrypt: false
```

!!! warning "Unencrypted scratchpads store secrets in plaintext"
    Only disable encryption if you understand the security implications.
    The scratchpad may contain sensitive data such as API keys, database
    URLs, or deployment credentials.

### Custom priority order

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

### Adjusted drift thresholds

Raise or lower the entry-count thresholds that trigger drift warnings:

```yaml
# .ctxrc
entry_count_learnings: 50   # warn above 50 learnings (default: 30)
entry_count_decisions: 10   # warn above 10 decisions (default: 20)
convention_line_count: 300  # warn above 300 lines (default: 200)
```

Set any threshold to `0` to disable that specific check.

---

**See also**: [CLI Reference](cli-reference.md) | [Context Files](context-files.md) | [Scratchpad](scratchpad.md)
