---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Steering
icon: lucide/compass
---

![ctx](../images/ctx-banner.png)

## `ctx steering`

Manage **steering files**: persistent behavioral rules for AI
coding assistants.

A steering file is a small markdown document with YAML
frontmatter that tells the AI *how to behave* in a specific
context. `ctx steering` keeps those files in
`.context/steering/`, decides which ones apply for a given
prompt, and syncs them out to each AI tool's native format
(Claude Code, Cursor, Kiro, Cline).

```bash
ctx steering <subcommand>
```

!!! tip "Steering vs decisions vs conventions"
    The three look similar on disk but serve different purposes:

    - **Decisions** record *what* was chosen and *why*.
      Consumed mostly by humans (and by the agent via
      `ctx agent`).
    - **Conventions** describe *how the codebase is written*.
      Consumed as reference material.
    - **Steering** tells the AI *how to behave when asked
      about X*. Consumed by the AI tool's prompt injection
      layer, conditionally on prompt match.

    If you find yourself writing "the AI should always do X"
    — that belongs in steering, not decisions.

### Anatomy of a steering file

```yaml
---
name: security
description: Security rules for all code changes
inclusion: always    # always | auto | manual
tools: []            # empty = all tools
priority: 10         # lower = injected first
---

# Security rules

- Validate all user input at system boundaries.
- Never log secrets, tokens, or credentials.
- Prefer constant-time comparison for tokens.
```

**Inclusion modes**:

| Mode      | When it's included                                   |
|-----------|------------------------------------------------------|
| `always`  | Every prompt, unconditionally                        |
| `auto`    | When the prompt matches the `description` keywords   |
| `manual`  | Only when the user names it explicitly               |

**Priority**: lower numbers inject first, so high-priority
rules appear at the top of the prompt. Default is `50`.

**Tools**: an empty list means all configured tools receive
the file; list specific tool names to scope it.

### `ctx steering init`

Create a starter set of steering files in `.context/steering/`
to use as a scaffolding baseline.

**Examples**:

```bash
ctx steering init
```

### `ctx steering add`

Create a new steering file with default frontmatter.

```bash
ctx steering add <name>
```

**Arguments**:

- `name`: Steering file name (without `.md` extension)

**Examples**:

```bash
ctx steering add security
# Created .context/steering/security.md
```

The generated file uses `inclusion: manual` and `priority: 50`
by default. Edit the frontmatter to change behavior.

### `ctx steering list`

List all steering files with their inclusion mode, priority,
and tool scoping.

**Examples**:

```bash
ctx steering list
```

### `ctx steering preview`

Preview which steering files would be included for a given
prompt. Useful for validating `auto`-inclusion descriptions
against realistic prompts.

```bash
ctx steering preview [prompt]
```

**Examples**:

```bash
ctx steering preview "create a REST API endpoint"
# Steering files matching prompt "create a REST API endpoint":
#   api-standards        inclusion=auto     priority=20  tools=all
#   security             inclusion=always   priority=10  tools=all
```

### `ctx steering sync`

Sync steering files to tool-native formats for tools that
have a **built-in rules primitive**. Not every tool needs
this — Claude Code and Codex use a different delivery
mechanism (see below).

**Examples**:

```bash
ctx steering sync
```

**Which tools are sync targets?**

| Tool         | Sync target          | Mechanism                               |
|--------------|----------------------|-----------------------------------------|
| Cursor       | `.cursor/rules/`     | Cursor reads the directory natively     |
| Cline        | `.clinerules/`       | Cline reads the directory natively      |
| Kiro         | `.kiro/steering/`    | Kiro reads the directory natively       |
| Claude Code  | *(no-op)*            | **Delivered via hook + MCP** — see next section |
| Codex        | *(no-op)*            | Same as Claude Code                     |

For the three native-rules tools, `ctx steering sync` writes
each matching steering file to the appropriate directory
with tool-specific frontmatter transforms. Unchanged files
are skipped (idempotent).

### How Claude Code and Codex consume steering

Claude Code has no native "steering files" primitive, so
`ctx steering sync` skips it entirely. Instead, steering
reaches Claude through **two non-sync channels**, both
activated by `ctx setup claude-code` (which installs the
plugin):

**1. Automatic injection via the `PreToolUse` hook.** The
Claude Code plugin wires a `PreToolUse` hook that runs
`ctx agent --budget 8000` before each tool call. `ctx
agent` loads `.context/steering/` and calls
`steering.Filter` with an **empty prompt**, so only files
with `inclusion: always` match. Those files are included
as **Tier 6** of the context packet. The packet is
printed on stdout, which Claude Code injects as
additional context. This fires on every tool call; no
user action.

**2. On-demand MCP tool call (`ctx_steering_get`).** The
ctx plugin ships a `.mcp.json` file that automatically
registers the `ctx` MCP server (`ctx mcp serve`) with
Claude Code on plugin install. Once registered, Claude
can invoke the `ctx_steering_get` tool mid-task to fetch
matching steering files for a specific prompt. This is
the **only** path that resolves `inclusion: auto` and
`inclusion: manual` matches for Claude Code — Claude
passes the prompt to the MCP tool, which runs the
keyword match against each file's description.

**Verify the MCP server is registered**:

```bash
claude mcp list
```

Expected line: `ctx: ctx mcp serve - ✓ Connected`. If
it's missing, reinstall the plugin from Claude Code
(`/plugin` → find `ctx` → uninstall → install again) —
older plugin versions shipped without the `.mcp.json`
file.

!!! warning "Prefer `inclusion: always` for Claude Code"
    Because the PreToolUse hook passes an empty prompt to
    `ctx agent`, only `always` files fire automatically.
    `auto` files require Claude to call the
    `ctx_steering_get` MCP tool on its own; `manual` files
    require an explicit user invocation. For rules that
    should reliably fire on every Claude Code session, use
    `inclusion: always`. Reserve `auto`/`manual` for
    situational libraries where the opt-in cost is
    acceptable and you understand Claude may not pull
    them in without prompting.

    The foundation files scaffolded by `ctx init` already
    default to `inclusion: always` for this reason.

**Practical implications**:

- Running `ctx steering sync` before starting a Claude
  session does **nothing** for Claude's benefit. Skip it.
- `ctx steering preview` still works for validating your
  descriptions — it doesn't depend on sync.
- If Claude Code is your only tool, the `ctx steering`
  commands you care about are `add`, `list`, `preview`,
  `init` — never `sync`.
- If you use both Claude Code **and** (say) Cursor,
  `ctx steering sync` covers Cursor (where `auto` and
  `manual` work natively) while the hook+MCP pipeline
  covers Claude Code. For rules you need to fire
  automatically on both, use `inclusion: always`.

### `ctx agent` integration

When `ctx agent` builds a context packet, steering files are
loaded as Tier 6 of the budget-aware assembly (see
[`ctx agent`](init-status.md#ctx-agent)). Files with
`inclusion: always` are always included; `auto` files are
scored against the current prompt and included in priority
order until the tier budget is exhausted.

### See also

- [`ctx setup`](setup.md) — configure which tools receive
  steering syncs
- [`ctx trigger`](trigger.md) — lifecycle scripts (a different
  hooking concept, see below)
- [Building steering files recipe](../recipes/steering.md)
  — walkthrough from first file to synced output
