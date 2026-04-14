---
title: Home
icon: lucide/home
---

![ctx](../images/ctx-banner.png)

* **`ctx` is not a prompt**.
* **`ctx` is version-controlled cognitive state**.

`ctx` is the persistence layer for human-AI reasoning.

*Deterministic. Git-native. Human-readable. Local-first*.

**Start here**.

Learn what `ctx` does, set it up, and run your first session.

!!! warning "Pre-1.0: Moving Fast"
    `ctx` is under active development. This website tracks the
    **development branch**, not the latest release:

    Some features described here may not exist in the binary
    you have installed.

    Expect rough edges.

    If something is missing or broken,
    [open an issue](https://github.com/ActiveMemory/ctx/issues).

---

## Introduction

### [About](about.md)

What `ctx` is, how it works, and why **persistent context
changes** how you work with AI.

### [Is It Right for Me?](is-ctx-right.md)

Good fit, not-so-good fit, and a **5-minute trial**
to find out for yourself.

### [FAQ](faq.md)

Quick answers to the questions newcomers ask most about
**`ctx`**, files, tooling, and trade-offs.

---

## Get Started

### [Getting Started](getting-started.md)

Install the **binary**, set up the **plugin**, and **verify** it works.

### [Your First Session](first-session.md)

**Step-by-step** walkthrough from `ctx init` to verified recall.

### [Common Workflows](common-workflows.md)

Day-to-day commands for **tracking** context, **checking** health,
and browsing **history**.

---

## Concepts

### [Context Files](context-files.md)

What each `.context/` file does. What's their **purpose**.
How do we best **leverage** them.

### [Configuration](configuration.md)

Flexible **configuration**: `.ctxrc`, environment variables, and CLI flags.

### [Hub](hub.md)

A **fan-out channel** for decisions, learnings, conventions, and
tasks that need to cross **project boundaries** — without replicating
everything else.

---

## Working with AI

### [Prompting Guide](prompting-guide.md)

**Effective prompts** for AI sessions with `ctx`.

### [Keeping AI Honest](keeping-ai-honest.md)

AI agents **confabulate**: they invent history, claim familiarity
with decisions never made, and sometimes declare tasks complete
when they aren't. Tools and habits to push back.

### [My AI Keeps Making the Same Mistakes](repeated-mistakes.md)

Stop **rediscovering** the same bugs and dead-ends across sessions.

### [Joining a Project](joining-a-project.md)

You inherited a `.context/` directory. Get **oriented fast**:
priority order, what to read first, how to ramp up.

---

## Customization

### [Steering Files](steering.md)

Tell the assistant **how to behave** when a specific kind
of prompt arrives.

### [Lifecycle Triggers](triggers.md)

Make things **happen** at session boundaries: block dangerous
tool calls, inject standup notes, log file saves.

---

## Community

### [#ctx](community.md)

We are the builders who care about **durable** context.<br />
Join the community. Hang out in IRC. Star `ctx` on GitHub.

### [Contributing](contributing.md)

**Development setup**, project layout, and pull request process.
