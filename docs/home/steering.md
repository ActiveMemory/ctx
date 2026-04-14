---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Steering Files
icon: lucide/compass
---

![ctx](../images/ctx-banner.png)

## Steering Files

`ctx` projects talk to AI assistants through several layers
— context files, decisions, conventions, the agent context
packet — but none of those can tell the assistant *how to
behave* when a specific kind of prompt arrives. That's what
**steering files** are for.

A steering file is a small markdown document with YAML
frontmatter that says: "when the user asks about X, prepend
these rules to the prompt." `ctx` manages those files in
`.context/steering/`, decides which ones match each prompt,
and syncs them out to each AI tool's native config (Claude
Code, Cursor, Kiro, Cline) so the rules actually land in the
prompt pipeline.

## Not the Same as Decisions or Conventions

The three look similar on disk but serve different purposes:

| Kind                                             | Purpose                                    |
|--------------------------------------------------|--------------------------------------------|
| [Decisions](context-files.md) (`DECISIONS.md`)   | *What* was chosen and *why*                |
| Conventions (`CONVENTIONS.md`)                   | *How* the codebase is written              |
| **Steering** (`.context/steering/*.md`)          | **How the AI should behave** on matching prompts |

If you find yourself writing "the AI should always do X when
asked about Y," that belongs in steering, not decisions.

## Your First Steering Files

**`ctx init` scaffolds four foundation steering files** in
`.context/steering/` so you start with something to edit
rather than an empty directory:

| File            | What to fill in                                    |
|-----------------|-----------------------------------------------------|
| `product.md`    | What the project is, who it's for, what's out of scope |
| `tech.md`       | Languages, frameworks, runtime, hard constraints   |
| `structure.md`  | Directory layout, where new files go, naming rules |
| `workflow.md`   | Branch strategy, commit conventions, pre-commit checks |

Each file starts with an inline HTML comment explaining the
three inclusion modes, priority semantics, and tool scoping.
The comment is invisible in rendered markdown but visible
when you open the file to edit it — it's self-documenting
scaffolding, not forever guidance. Delete the comment once
you've customized the file.

Default settings for foundation files:

- `inclusion: always` — fires on every AI tool call
- `priority: 10` — injected near the top of the prompt
- `tools: []` — applies to every configured AI tool

**You should open each of these files and replace the
placeholder content with your project's actual rules.**
Re-running `ctx init` is safe: existing files are left
alone, so your edits survive. Use `ctx init --no-steering-init`
to opt out of the scaffold entirely.

## Inclusion Modes

Each steering file declares an inclusion mode in its
frontmatter:

| Mode      | When the file is included                    |
|-----------|-----------------------------------------------|
| `always`  | Every prompt, unconditionally                 |
| `auto`    | When the prompt keywords match the file's description |
| `manual`  | Only when the user explicitly names the file  |

**Which mode to pick depends on the AI tool you use**,
because the two tool families consume steering very
differently.

**Claude Code and Codex** — prefer `inclusion: always`
for rules that must fire reliably. These tools have two
delivery channels:

1. **The plugin's `PreToolUse` hook** runs `ctx agent`
   with an **empty prompt**, so only `always` files match
   and get injected automatically on every tool call.
2. **The `ctx_steering_get` MCP tool**, registered
   automatically when the ctx plugin is installed. Claude
   can call this tool mid-task to fetch `auto` or
   `manual` files matching a specific prompt. Verify
   with `claude mcp list` — look for `ctx: ✓ Connected`.

Use `always` for invariants and anything that **must**
fire every session. Use `auto` for situational rules
where "Claude fetches this when the prompt is relevant"
is the right behavior — those still land, just on
Claude's judgment. Use `manual` for reference libraries
you'll name explicitly.

**Cursor, Cline, Kiro** — `auto` is the natural default.
These tools read `.cursor/rules/`, `.clinerules/`, or
`.kiro/steering/` natively and resolve the description
match on their own, so `auto` files fire when the prompt
matches. `manual` files load on explicit invocation.
`always` still works but consumes context budget on
every turn.

**Mixed setups** — if a rule must fire on Claude Code,
pick `always`, even if it's overkill for your Cursor
setup. The context budget cost is small; the alternative
(silently not firing) is worse.

## Two Families of AI Tools, Two Delivery Paths

Not every AI tool consumes steering the same way. ctx
handles two tool families differently, and it's worth
knowing which family your editor is in before you wonder
why a rule isn't firing.

**Native-rules tools** — **Cursor**, **Cline**, **Kiro** —
have a built-in rules primitive. They read a specific
directory (`.cursor/rules/`, `.clinerules/`,
`.kiro/steering/`) and apply the rules they find there.
ctx handles these via `ctx steering sync`, which exports
your files into the tool-native format. Run `sync`
whenever you edit a steering file.

**Hook + MCP tools** — **Claude Code**, **Codex** — have
no native rules primitive, so `ctx steering sync` is a
**no-op** for them. Instead, ctx delivers steering through
two non-sync channels:

1. **Automatic injection via a `PreToolUse` hook**. The
   `ctx setup claude-code` plugin wires a hook that runs
   `ctx agent --budget 8000` before each tool call.
   `ctx agent` loads your steering files, filters them by
   the active prompt, and includes matching bodies in the
   context packet it prints. Claude Code feeds that output
   back into its context. Every tool call, automatically.
2. **On-demand via the `ctx_steering_get` MCP tool**. The
   ctx MCP server exposes a tool Claude can call mid-task
   to fetch matching steering files for a specific prompt.
   Claude decides when to call it — it's not automatic.

Both channels activate when you run
`ctx setup claude-code --write`. After that, steering just
works for Claude Code.

**Practical takeaway**:

- Using Cursor/Cline/Kiro only? → Run `ctx steering sync`
  after edits.
- Using Claude Code or Codex only? → Never run `sync`; the
  hook+MCP pipeline handles it.
- Using both? → Run `sync` for the native-rules tools; the
  hook+MCP pipeline covers Claude Code automatically.

## Two Shapes of Automation: Rules and Scripts

Steering is one of **two** hook-like layers ctx provides for
customizing AI behavior. They're complementary:

- **Steering** — *persistent rules* that get prepended to
  prompts. Declarative, text-only, scored by match.
- **[Triggers](triggers.md)** — *executable shell scripts*
  that fire at lifecycle events. Imperative, runs arbitrary
  code, gated by exit codes.

Pick steering when you want "always remind the AI of X."
Pick triggers when you want "do Y when event Z happens."
They can coexist — many projects use both.

## Where to Go Next

- **[Writing Steering Files](../recipes/steering.md)** —
  a six-step walkthrough: scaffold, write the rule, preview
  matches, list, get-rules-in-front-of-the-AI (two paths
  depending on tool family), verify.
- **[`ctx steering` reference](../cli/steering.md)** — full
  command, flag, and frontmatter reference; includes the
  per-tool delivery-mechanism table and a dedicated section
  on how Claude Code and Codex consume steering.
- **[`ctx setup`](../cli/setup.md)** — configure which AI
  tools receive steering. For Cursor/Cline/Kiro this is
  about sync targets; for Claude Code/Codex it installs
  the plugin that wires the `PreToolUse` hook and MCP
  server.
- **[Lifecycle Triggers](triggers.md)** — the imperative
  companion to steering files.
