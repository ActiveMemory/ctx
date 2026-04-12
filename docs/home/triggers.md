---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Lifecycle triggers
icon: lucide/zap
---

![ctx](../images/ctx-banner.png)

## Lifecycle Triggers

Some things can't be expressed as a rule you want the AI to
follow. Sometimes you want something to **happen**: block a
dangerous tool call, inject today's standup notes into the
next session, log every file save to a journal. That's what
**triggers** are for.

A trigger is an executable shell script that `ctx` runs at a
specific **lifecycle event**: the start of a session, before
a tool call, when a file is saved, and so on. Triggers read a
JSON payload from stdin, do whatever they need, and write a
JSON response on stdout. They can **allow**, **block**, or
**inject context** into the pipeline depending on the event
type.

## Trigger Types

| Type            | Fires when                          | Use case                               |
|-----------------|-------------------------------------|----------------------------------------|
| `session-start` | A new AI session begins             | Inject rotating context, standup notes |
| `session-end`   | An AI session ends                  | Persist summaries, send notifications  |
| `pre-tool-use`  | Before a tool call executes         | Block, gate, or audit                  |
| `post-tool-use` | After a tool call completes         | Log, react, post-process               |
| `file-save`     | A file is saved                     | Lint on save, update indices           |
| `context-add`   | A new entry is added to `.context/` | Cross-link, notify, enrich             |

## Triggers are arbitrary code — treat them like pre-commit hooks

!!! warning "Only enable scripts you've read and understand"
    A trigger is a shell script with the executable bit set.
    It runs with the same privileges as your AI tool and
    receives JSON input on stdin. A malicious or buggy
    trigger can block tool calls, corrupt context files, or
    exfiltrate data.

    `ctx trigger add` intentionally creates new scripts
    **disabled** (no executable bit). You must
    `ctx trigger enable <name>` after reviewing the contents.
    That's not a suggestion — it's the security model.

## Three hook-like layers in ctx

Triggers are one of **three** distinct hook-like concepts in
ctx. The names are similar but the owners and use cases are
not:

| Layer                  | Owned by    | Where they live                         | When to use                                |
|------------------------|-------------|-----------------------------------------|--------------------------------------------|
| **`ctx trigger`**      | You         | `.context/hooks/<type>/*.sh`            | Project-specific automation, any AI tool   |
| **`ctx system` hooks** | ctx itself  | built-in, wired into tool configs       | Built-in nudges (you don't author these)   |
| **Claude Code hooks**  | Claude Code | `.claude/settings.local.json`           | Claude-Code-only tool-specific integration |

This page is about the first category. The other two run
automatically and are invisible to you.

## Triggers vs steering — same problem, different shape

Triggers are the imperative counterpart to
[**steering files**](steering.md). Steering expresses
*persistent rules* the AI reads before each prompt; triggers
express *side effects* that run on lifecycle events. They're
complementary, not competing:

- Want the AI to *remember* something? → Steering.
- Want a script to *run* when something happens? → Trigger.

Most projects use both.

## Where to go next

- **[Authoring Lifecycle Triggers](../recipes/triggers.md)**
  — walkthrough with security guidance: scaffold, test,
  enable, iterate.
- **[`ctx trigger` reference](../cli/trigger.md)** — command
  reference, trigger type table, input/output contract.
- **[Steering files](steering.md)** — the declarative
  counterpart to triggers.
