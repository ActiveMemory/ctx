---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: About
icon: lucide/info
---

![ctx](../images/ctx-banner.png)

"**Creation, not code**;
**Context, not prompts**;
**Verification, not vibes**."
    
[Read the `ctx` Manifesto →](../index.md)

"**Without durable context, intelligence *resets*;
with `ctx`, creation compounds**."

Without persistent memory, every session starts at *zero*;
`ctx` makes sessions **cumulative**.

[Join the `ctx` Community →](community.md)

## What Is `ctx`?

`ctx` (*Context*) is a file-based system that enables AI coding assistants
to persist project knowledge across sessions.
It lives in a `.context/` directory in your repo.

* A session is *interactive*.
* `ctx` enables **cognitive continuity**.
* **Cognitive continuity** enables durable, *symbiotic-like* human–AI workflows.

!!! tip "Context Files"
    Context files let AI tools **remember** [*decisions*, *conventions*, and *learnings*](context-files.md):

    Context files are **explicit** and **versionable** contracts 
    between you and your agents.


## Why Do I Keep Re-Explaining My Codebase?!?!

You open a new AI session. The first thing you do is re-explain your project.

**Again**.

The architecture, the database choice, the naming conventions, the thing you
tried last week that didn't work. You have said all of this before (*maybe
yesterday, maybe an hour ago*), but the AI doesn't know that.

* You explain the same architecture **every session**;
* The AI suggests an approach you already rejected, **again**;
* A decision you made three sessions ago gets relitigated **from scratch**;
* You spend more time *setting context* than *building features*.

This isn't AI failure: It is **staleness**. 

Without persistent memory, every session starts at zero.

### Day and Night

`ctx` gives your AI a memory that persists across sessions:

=== "Without ctx"

    ```text
    Session 12: Monday morning

    AI:  "I'll scaffold a MySQL schema for this."
    You: "We use PostgreSQL."

    AI:  "For auth, session-based would be simplest."
    You: "We decided on JWT."

    AI:  "Endpoints in camelCase?"
    You: "snake_case! Like the rest of the API!!!!!!"

    You: *scrolls up to copy an existing endpoint*

    20 minutes spent restoring decisions that already exist.
    ```

=== "With ctx"

    ```text
    Session 12: Monday morning

    You:  "Let's continue the auth implementation."
    AI:   "Per DECISIONS.md (2026-02-12): JWT, not sessions: 
           The project uses PostgreSQL with JWT
           authentication. The API follows snake_case conventions.
           Last session you completed the login endpoint — the
           registration endpoint is next on TASKS.md. Want me to
           start there?"

    0 minutes re-explaining. Straight to building.
    ```

### How ctx Solves This

`ctx` creates a `.context/` directory in your project that stores [structured
knowledge files](context-files.md):

| File              | What It Remembers                             |
|-------------------|-----------------------------------------------|
| `TASKS.md`        | What you're working on and what's next        |
| `DECISIONS.md`    | Architectural choices and *why* you made them |
| `LEARNINGS.md`    | Gotchas, bugs, things that didn't work        |
| `CONVENTIONS.md`  | Naming patterns, code style, project rules    |
| `CONSTITUTION.md` | Hard rules the AI must never violate          |

These files can **version with your code** in `git`: 

* They load automatically at the session start 
  (*via hooks in Claude Code, or manually with `ctx agent` for
  other tools*). 
* The AI **reads** them, **cites** them, and **builds on** them, instead
  of asking you to start over. 
      * And when it **acts**, it can point to the exact file and line that 
        **justifies** the choice.

Every decision you record, every lesson you capture,
makes the *next* session **smarter**.

`ctx` **accumulates**.

----

**Connect with ctx**

* [Join the Community →](community.md): ask questions, share workflows, and help shape what comes next
* [Read the Blog →](../../blog/): real-world patterns, ponderings, and lessons learned from building `ctx` using `ctx`

----

**Ready to get started?**

* [Getting Started →](getting-started.md) — full installation and setup
* [Your First Session →](first-session.md) — step-by-step walkthrough from `ctx init` to verified recall
