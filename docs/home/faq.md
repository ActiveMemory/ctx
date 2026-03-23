---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: FAQ
icon: lucide/help-circle
---

![ctx](../images/ctx-banner.png)

## Why Markdown?

Markdown is human-readable, version-controllable, and tool-agnostic.
Every AI model can parse it natively. Every developer can read it in a
terminal, a browser, or a code review. There's no schema to learn, no
binary format to decode, no vendor lock-in. You can inspect your
context with `cat`, diff it with `git diff`, and review it in a PR.

## Does ctx work offline?

Yes. ctx is completely local. It reads and writes files on disk,
generates context packets from local state, and requires no network
access. The only feature that touches the network is the optional
[webhook notifications](../recipes/webhook-notifications.md) hook,
which you have to explicitly configure.

## What gets committed to git?

The `.context/` directory: yes, commit it. That's the whole point.
Team members and AI agents read the same context files.

What **not** to commit:

- **`.ctx.key`**: your encryption key. Stored at `~/.ctx/.ctx.key`,
  never in the repo. `ctx init` handles this automatically.
- **`journal/`** and **`logs/`**: generated data, potentially large.
  `ctx init` adds these to `.gitignore`.
- **`scratchpad.enc`**: your choice. It's encrypted, so it's safe to
  commit if you want shared scratchpad state. See
  [Scratchpad](../reference/scratchpad.md) for details.

## How big should my token budget be?

The default is 8000 tokens, which works well for most projects.
Configure it via `.ctxrc` or the `CTX_TOKEN_BUDGET` environment
variable:

```bash
# In .ctxrc
token_budget = 12000

# Or as an environment variable
export CTX_TOKEN_BUDGET=12000

# Or per-invocation
ctx agent --budget 4000
```

Higher budgets include more context but cost more tokens per request.
Lower budgets force sharper prioritization: ctx drops lower-priority
content first, so CONSTITUTION and TASKS always make the cut.

See [Configuration](configuration.md) for all available settings.

## Why not a database?

Files are inspectable, diffable, and reviewable in pull requests.
You can `grep` them, `cat` them, pipe them through `jq` or `awk`.
They work with every version control system and every text editor.

A database would add a dependency, require migrations, and make
context opaque. The design bet is that context should be as visible
and portable as the code it describes.

## Does it work with tools other than Claude Code?

Yes. `ctx agent` outputs a context packet that any AI tool can
consume: paste it into ChatGPT, Cursor, Copilot, Aider, or anything
else that accepts text input.

Claude Code gets first-class integration via the ctx plugin (hooks,
skills, automatic context loading). VS Code Copilot Chat has a
dedicated ctx extension. Other tools integrate via generated
instruction files or manual pasting.

See [Integrations](../operations/integrations.md) for tool-specific
setup, including the [multi-tool recipe](../recipes/multi-tool-setup.md).

## Can I use ctx on an existing project?

Yes. Run `ctx init` in any repo and it creates `.context/` with
template files. Start recording decisions, tasks, and conventions as
you work. Context grows naturally; you don't need to backfill
everything on day one.

See [Getting Started](getting-started.md) for the full setup flow,
or [Joining a ctx Project](joining-a-project.md) if someone else
already initialized it.

## What happens when context files get too big?

Token budgeting handles this automatically. `ctx agent` prioritizes
content by file priority (CONSTITUTION first, GLOSSARY last) and
trims lower-priority entries when the budget is tight.

For manual maintenance, `ctx compact` archives completed tasks and
old entries, keeping active context lean. You can also run
`ctx task archive` to move completed tasks out of TASKS.md.

The goal is to keep context files focused on **current** state.
Historical entries belong in git history or the archive.

## Is .context/ meant to be shared?

Yes. Commit it to your repo. Every team member and every AI agent
reads the same files. That's the mechanism for shared memory:
decisions made in one session are visible in the next, regardless
of who (or what) starts it.

The only per-user state is the encryption key (`~/.ctx/.ctx.key`)
and the optional scratchpad. Everything else is team-shared by
design.

----

**Related**:

* [Getting Started](getting-started.md) - installation and first setup
* [Configuration](configuration.md) - `.ctxrc`, environment variables, and defaults
* [Context Files](context-files.md) - what each file does and how to use it
