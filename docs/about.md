---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: About ctx
icon: lucide/info
---

![ctx](images/ctx-banner.png)

## `ctx`

`ctx` (*Context*) is a file-based system that enables AI coding assistants to
persist project knowledge across sessions. Instead of re-explaining your
codebase every time, context files let AI tools remember decisions,
conventions, and learnings:

* A session is interactive.
* `ctx` enables **cognitive continuity**.
* **Cognitive continuity** enables durable, *symbiotic-like* human–AI workflows.

!!! quote "The `ctx` Manifesto"
    **Creation, not code. Context, not prompts. Verification, not vibes.**

    Without durable context, intelligence resets.
    With `ctx`, creation compounds.

    **[Read the Manifesto →](https://ctx.ist/)**

## Community

**[Join the ctx Community →](community.md)**

## Why Do I Keep Re-Explaining My Codebase?!?!

You open a new AI session. The first thing you do is re-explain your project.

Again.

The architecture, the database choice, the naming conventions, the thing you
tried last week that didn't work. You've said all of this before — maybe
yesterday, maybe an hour ago — but the AI doesn't know that.

- You explain the same architecture **every session**
- The AI suggests an approach you already rejected — **again**
- A decision you made three sessions ago gets relitigated from scratch
- You spend more time *setting context* than *building features*

This isn't an AI problem. It's a **context problem**. Without persistent
memory, every session starts at zero.

### Before & After

=== "Without ctx"

    ```text
    Session 12 — Monday morning

    You:  "We use PostgreSQL, not MySQL. I explained this on Thursday."
    AI:   "Got it! Let me adjust the schema for PostgreSQL..."

    You:  "Also, we decided to use JWT for auth, not sessions."
    AI:   "Understood! Here's the updated approach..."

    You:  "And the API uses snake_case, not camelCase."
    AI:   "I'll fix the naming..."

    20 minutes spent re-explaining before any real work starts.
    ```

=== "With ctx"

    ```text
    Session 12 — Monday morning

    You:  "Let's continue the auth implementation."
    AI:   "Based on your context: the project uses PostgreSQL with JWT
           authentication. The API follows snake_case conventions.
           Last session you completed the login endpoint — the
           registration endpoint is next on TASKS.md. Want me to
           start there?"

    0 minutes re-explaining. Straight to building.
    ```

The difference: **ctx gives your AI a memory that persists across sessions**.

### How ctx Solves This

`ctx` creates a `.context/` directory in your project that stores structured
knowledge files:

| File | What It Remembers |
|------|-------------------|
| `TASKS.md` | What you're working on and what's next |
| `DECISIONS.md` | Architectural choices and *why* you made them |
| `LEARNINGS.md` | Gotchas, bugs, things that didn't work |
| `CONVENTIONS.md` | Naming patterns, code style, project rules |
| `CONSTITUTION.md` | Hard rules the AI must never violate |

These files **version with your code** in git. They load automatically at
session start (via hooks in Claude Code, or manually with `ctx agent` for
other tools). The AI reads them, cites them, and builds on them — instead
of asking you to start over.

Context accumulates. Every decision you record, every lesson you capture,
makes the *next* session smarter than the last.

----

**Ready to get started?**

- [Getting Started →](getting-started.md) — full installation and setup
- [Your First Session →](first-session.md) — step-by-step walkthrough from `ctx init` to verified recall
