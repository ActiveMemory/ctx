---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Trigger
icon: lucide/zap
---

![ctx](../images/ctx-banner.png)

## `ctx trigger`

Manage **lifecycle triggers**: executable scripts that fire at
specific events during an AI session. Triggers can block tool
calls, inject context, automate reactions — any side effect
you want at session boundaries, tool boundaries, or file-save
events.

```bash
ctx trigger <subcommand>
```

!!! warning "Triggers execute arbitrary scripts"
    A trigger is a shell script with the executable bit set.
    It runs with the same privileges as your AI tool and
    receives JSON input on stdin. Treat triggers like
    pre-commit hooks: only enable scripts you've read and
    understand. A malicious or buggy trigger can block tool
    calls, corrupt context files, or exfiltrate data.

### Where triggers live

Triggers live in `.context/hooks/<trigger-type>/` as executable
scripts. The on-disk directory name is still `hooks/` for
historical reasons even though the command is `ctx trigger`.
Each script:

- Reads a JSON payload from stdin.
- Returns a JSON payload on stdout.
- Returns a non-zero exit code to block or error.

```
.context/
└── hooks/
    ├── session-start/
    │   └── inject-context.sh
    ├── pre-tool-use/
    │   └── block-legacy.sh
    └── post-tool-use/
        └── record-edit.sh
```

### Trigger types

| Type            | Fires when                           |
|-----------------|--------------------------------------|
| `session-start` | An AI session begins                 |
| `session-end`   | An AI session ends                   |
| `pre-tool-use`  | Before an AI tool call is executed   |
| `post-tool-use` | After an AI tool call returns        |
| `file-save`     | When a file is saved                 |
| `context-add`   | When a context entry is added        |

### Input and output contract

Each trigger receives a JSON object on stdin with the event
details. Minimal contract (fields vary by trigger type):

```json
{
  "type": "pre-tool-use",
  "tool": "write_file",
  "path": "src/auth.go",
  "session_id": "abc123-..."
}
```

The trigger may write a JSON object to stdout to influence
behavior. Example for a blocking `pre-tool-use` trigger:

```json
{
  "action": "block",
  "message": "Editing src/auth.go requires approval from #security"
}
```

For non-blocking event loggers, simply read stdin and exit 0
without writing to stdout.

### `ctx trigger add`

Create a new trigger script with a template. The generated
file has a bash shebang, a stdin reader using `jq`, and a
basic JSON output structure.

```bash
ctx trigger add <trigger-type> <name>
```

**Arguments**:

- `trigger-type`: One of `session-start`, `session-end`,
  `pre-tool-use`, `post-tool-use`, `file-save`, `context-add`
- `name`: Script name (without `.sh` extension)

**Examples**:

```bash
ctx trigger add session-start inject-context
# Created .context/hooks/session-start/inject-context.sh

ctx trigger add pre-tool-use block-legacy
# Created .context/hooks/pre-tool-use/block-legacy.sh
```

The generated script is **not** executable by default. Enable
it with `ctx trigger enable` after reviewing the contents.

### `ctx trigger list`

List all discovered triggers, grouped by trigger type, with
their enabled/disabled status.

**Examples**:

```bash
ctx trigger list
```

### `ctx trigger test`

Run all enabled triggers of a given type against a mock
payload. Use `--tool` and `--path` to customize the mock
input for tool-related events.

```bash
ctx trigger test <trigger-type> [flags]
```

**Flags**:

| Flag     | Description                       |
|----------|-----------------------------------|
| `--tool` | Tool name to put in mock input    |
| `--path` | File path to put in mock input    |

**Examples**:

```bash
ctx trigger test session-start
ctx trigger test pre-tool-use --tool write_file --path src/main.go
```

### `ctx trigger enable`

Enable a trigger by setting its executable permission bit.
Searches every trigger-type directory for a script matching
`<name>`.

```bash
ctx trigger enable <name>
```

**Examples**:

```bash
ctx trigger enable inject-context
# Enabled .context/hooks/session-start/inject-context.sh
```

### `ctx trigger disable`

Disable a trigger by clearing its executable permission bit.
Searches every trigger-type directory for a script matching
`<name>`.

```bash
ctx trigger disable <name>
```

**Examples**:

```bash
ctx trigger disable inject-context
# Disabled .context/hooks/session-start/inject-context.sh
```

### Three hooking concepts in ctx — don't confuse them

This is a common source of confusion. `ctx` has three
distinct hook-like layers, and they serve different purposes:

| Layer                 | Owned by    | Where it runs                              | Configured via                         |
|-----------------------|-------------|--------------------------------------------|----------------------------------------|
| **`ctx trigger`**     | You         | `.context/hooks/<type>/*.sh`               | `ctx trigger add/enable`               |
| **`ctx system` hooks**| ctx itself  | built-in, called by ctx's own lifecycle    | internal (see `ctx system --help`)     |
| **Claude Code hooks** | Claude Code | `.claude/settings.local.json`              | edit JSON, or `/ctx-sanitize-permissions` |

Use `ctx trigger` when you want project-specific automation
that your AI tool will run at lifecycle events. Use Claude
Code hooks for tool-specific integrations that don't need to
be portable across tools. `ctx system` hooks are not something
you author — they're the internal nudge machinery that ships
with ctx.

### See also

- [`ctx steering`](steering.md) — persistent AI behavioral
  rules (a different concept; rules vs scripts)
- [Authoring triggers recipe](../recipes/triggers.md) — a
  full walkthrough with security guidance
