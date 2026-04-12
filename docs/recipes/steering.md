---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Writing Steering Files
icon: lucide/compass
---

![ctx](../images/ctx-banner.png)

# Writing Steering Files

Steering files tell your AI assistant **how to behave**, not
what was decided or how the codebase is written. This recipe
walks through writing a steering file from scratch, validating
which prompts will trigger it, and syncing it out to your
configured AI tools.

!!! tip "Before you start"
    If you're unsure whether a rule belongs in
    `steering/`, `DECISIONS.md`, or `CONVENTIONS.md`, read the
    "Steering vs decisions vs conventions" admonition on the
    [`ctx steering` reference page](../cli/steering.md). The
    short version: if the rule is "the AI should always do X
    when asked about Y," that's steering. Otherwise it's
    probably a decision or convention.

## Start here — customize the foundation files

**`ctx init` scaffolds four foundation steering files** for
you the first time you initialize a project:

| File                                  | Purpose                                     |
|---------------------------------------|---------------------------------------------|
| `.context/steering/product.md`        | Product context, goals, target users        |
| `.context/steering/tech.md`           | Tech stack, constraints, key dependencies   |
| `.context/steering/structure.md`      | Directory layout, naming conventions        |
| `.context/steering/workflow.md`       | Branch strategy, commit rules, pre-commit   |

Each file opens with an **inline HTML comment** that
explains the three inclusion modes, what `priority` means,
and the `tools` scope. The comment is invisible in
rendered markdown but visible when you edit the file.
Delete it once the file is yours.

All four default to `inclusion: always` and `priority: 10`
— they fire on **every** AI tool call until you customize
them. If you're reading this recipe and haven't touched
them yet, **open each one now and replace the placeholder
bullet list with actual rules for your project**. That's
the highest-leverage five minutes you can spend in a new
`ctx` setup.

What to fill in, by file:

**`product.md`** — The elevator pitch plus hard scope:

- One-sentence product description.
- Primary users and their top job-to-be-done.
- Two or three "this is explicitly out of scope" items
  so the AI doesn't wander.

**`tech.md`** — Technology and constraints:

- Languages and versions (`Go 1.22`, `Node 20`, etc.).
- Frameworks and key libraries.
- Runtime and deployment target.
- Hard constraints: "no CGO", "no network at test time",
  "no external DB for unit tests". These are the things
  that burn agents when they don't know them.

**`structure.md`** — Layout and naming:

- Top-level directories and their purpose.
- Where new files should go (and where they should NOT).
- Naming conventions for packages, files, types.

**`workflow.md`** — Process rules:

- Branch strategy (main-only, trunk-based, feature
  branches).
- Commit message format, signed-off-by requirement.
- Pre-commit and pre-push checks.
- Review expectations.

After editing, the next AI tool call in Claude Code will
pick up the new rules automatically via the plugin's
`PreToolUse` hook — no sync step, no restart. Other tools
(Cursor, Cline, Kiro) need `ctx steering sync` to export
into their native format.

!!! note "Prefer a bare `.context/steering/` directory?"
    Re-run `ctx init --no-steering-init` and delete the
    scaffolded files. `ctx init` leaves existing files
    alone, so the flag is only needed if you want to opt
    out of the initial scaffold.

The rest of this recipe walks through creating an
**additional**, scenario-specific steering file beyond the
four foundation defaults.

## Scenario

You're working on a project with a strict input-validation
policy: every new API handler must validate request bodies
before touching the database. You want the AI to flag this
concern automatically whenever it's asked to write an HTTP
handler, without you having to remind it every session.

!!! warning "Claude Code users: pick `always`, not `auto`"
    This walkthrough uses `inclusion: auto` because the
    scenario is a scoped rule that matches a specific kind
    of prompt. That works natively on **Cursor, Cline, and
    Kiro** (they resolve the `description` keyword match
    themselves).

    On **Claude Code**, `auto` does **not** fire through
    the plugin's `PreToolUse` hook — the hook passes an
    empty prompt to `ctx agent`, so only `always` files
    match. Claude can still reach an `auto` file by
    calling the `ctx_steering_get` MCP tool, but that
    requires Claude to decide to call it; there's no
    automatic injection.

    **If Claude Code is your tool**, set `inclusion:
    always` in Step 2 instead of `auto`. The rule will
    fire on every tool call regardless of topic. You may
    want to narrow the rule body so the extra tokens per
    turn aren't wasted on unrelated work.

    See the [`ctx steering` reference](../cli/steering.md)
    "Prefer `inclusion: always` for Claude Code" section
    for the full trade-off.

## Step 1 — scaffold the file

```bash
ctx steering add api-validation
```

That creates `.context/steering/api-validation.md` with default
frontmatter:

```yaml
---
name: api-validation
description:
inclusion: manual
tools: []
priority: 50
---
```

The defaults are deliberately conservative: `inclusion: manual`
means the file won't be applied until you opt in, which keeps
the rules out of the prompt until you've reviewed them.

## Step 2 — fill in the rule

Open the file and write the rule body plus a focused
description. The description is what `inclusion: auto` matches
against later.

```markdown
---
name: api-validation
description: HTTP handler input validation and request parsing
inclusion: auto
tools: []
priority: 20
---

# API request validation

Every new HTTP handler MUST:

1. Parse request bodies into typed structs, never `map[string]any`.
2. Validate required fields before any database call.
3. Return 400 with a machine-readable error for validation failures.
4. Use `context.Context` from the request for all downstream calls.

Prefer existing validation helpers in `internal/validate/`
rather than inline checks.
```

Notes on the choices:

- **`inclusion: auto`** — this rule should fire automatically
  on HTTP-handler-shaped prompts, not always.
- **`priority: 20`** — lower than the default, so this rule
  appears near the top of the prompt alongside other
  high-priority rules.
- **Description** is keyword-rich: "HTTP handler input
  validation and request parsing" — the `auto` matcher scores
  prompts against these words.

## Step 3 — preview which prompts match

Before committing the file, validate your description catches
the prompts you care about:

```bash
ctx steering preview "add an endpoint for updating user email"
```

Expected output:

```
Steering files matching prompt "add an endpoint for updating user email":
  api-validation       inclusion=auto     priority=20  tools=all
```

Good — the prompt matches. Try a negative case:

```bash
ctx steering preview "fix a bug in the JSON renderer"
```

Expected: empty match (or whatever else is currently `auto`).
If `api-validation` incorrectly fires for unrelated prompts,
tighten the description. If it misses prompts it should catch,
add more keywords.

## Step 4 — list to confirm metadata

```bash
ctx steering list
```

Should show `api-validation` alongside any other files,
with its inclusion mode and priority. If the list is wrong,
check the frontmatter for typos.

## Step 5 — get the rules in front of the AI

**Steering files are authored once in `.context/steering/`,
but how they reach the AI depends on which tool you use.**
There are two delivery mechanisms:

### Path A — native-rules tools (Cursor, Cline, Kiro)

These tools read a specific directory for rules. `ctx
steering sync` exports your files into that directory with
tool-specific frontmatter:

```bash
ctx steering sync
```

Depending on the active tool in `.ctxrc` or `--tool`:

| Tool   | Target             |
|--------|--------------------|
| Cursor | `.cursor/rules/`   |
| Cline  | `.clinerules/`     |
| Kiro   | `.kiro/steering/`  |

The sync is idempotent — unchanged files are skipped. Run
it whenever you edit a steering file.

### Path B — Claude Code and Codex (hook + MCP)

Claude Code and Codex have **no native rules primitive**,
so `ctx steering sync` is a **no-op** for them — it
deliberately skips both. Instead, steering reaches these
tools through two non-sync channels:

1. **`PreToolUse` hook** (automatic). The `ctx setup
   claude-code` plugin installs a hook that runs
   `ctx agent --budget 8000` before each tool call. `ctx
   agent` loads your steering files, filters them against
   the active prompt, and includes matching bodies as
   Tier 6 of the context packet. The packet gets injected
   into Claude's context automatically.

2. **`ctx_steering_get` MCP tool** (on-demand). Claude can
   call this MCP tool mid-task to fetch matching steering
   files for a specific prompt. Automatic activation comes
   from Claude's judgment, not a hook.

Both channels activate when you run:

```bash
ctx setup claude-code --write
```

That installs the plugin, wires the hook, and registers the
MCP server. After that, steering files you edit are picked
up on the next tool call — no sync step needed.

!!! tip "Running `ctx steering sync` with Claude Code"
    It won't error — it will simply report that Claude and
    Codex aren't sync targets and skip them. If Claude Code
    is your only tool, you never need to run `sync`. If you
    use both Claude Code **and** (say) Cursor, run `sync`
    to keep Cursor up to date; the Claude pipeline takes
    care of itself via the hook.

## Step 6 — verify the AI sees it

Open your AI tool and ask it something the rule should fire
on:

> "Add a POST /users endpoint that accepts email and name."

If the rule is working, the AI's first response should
mention input validation, typed structs, and the
`internal/validate/` package — because that's what the
steering file told it to do.

If nothing happens, the fix depends on which path you're on:

**Path A — Cursor/Cline/Kiro**:

1. Re-run `ctx steering preview` with the literal prompt to
   confirm the match.
2. Run `ctx steering list` and verify `inclusion` is `auto`,
   not `manual`.
3. Check the tool's own config directory (e.g.
   `.cursor/rules/`) — the file should be there after
   `ctx steering sync`.

**Path B — Claude Code**:

1. Re-run `ctx steering preview` with the literal prompt to
   confirm the match.
2. Verify the plugin is installed: `cat .claude/hooks.json`
   should include `ctx agent --budget 8000` under
   `PreToolUse`. If not, re-run `ctx setup claude-code --write`.
3. Run `ctx agent --budget 8000` manually and grep the
   output for your rule body. If it's there, the data is
   fine; if it's missing, the `inclusion` mode or
   `description` is at fault.
4. As a last resort, ask Claude directly: "Call the
   `ctx_steering_get` MCP tool with my prompt and show me
   the result." If the MCP tool returns your rule, Claude
   has access but isn't pulling it into the initial
   context packet — tighten the description keywords.

## Common mistakes

**Too-generic descriptions.** `description: general coding`
will match almost every prompt and flood the context window.
Keep descriptions specific to the scenario the rule applies to.

**Overlapping rules.** If two steering files match the same
prompt and contradict each other, the result is confusing.
Use `priority` to resolve, but better: merge the files or
narrow the descriptions so they don't overlap.

**Putting decisions in steering.** "We decided to use
PostgreSQL" is a decision, not a rule for the AI to follow on
every prompt. Record decisions with `ctx add decision`, not
`ctx steering add`.

**Committing `inclusion: always` without thinking.** Rules
marked `always` fire on every prompt, consuming tier-6 budget
permanently. Only use `always` for true invariants (security,
safety, licensing). Everything else should be `auto` or
`manual`.

## See also

- [`ctx steering` reference](../cli/steering.md) — full
  command, flag, and frontmatter reference.
- [`ctx setup`](../cli/setup.md) — configure which tools the
  steering sync writes to.
- [Authoring triggers](triggers.md) — if you want
  script-based automation, not rule-based prompt injection.
