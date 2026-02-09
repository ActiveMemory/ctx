---
title: Recipes
icon: lucide/chef-hat
---

![ctx](../images/ctx-banner.png)

Workflow recipes showing how ctx commands and skills work together.
Each recipe solves a specific problem by combining multiple tools
into a step-by-step workflow.

!!! tip "Commands vs. Skills"
    **Commands** (`ctx status`, `ctx add task`) run in your
    terminal.

    **Skills** (`/ctx-save`, `/ctx-next`) run inside
    your AI coding assistant. 

    Recipes combine both.

    Think of commands as **structure** and skills as **behavior**.

!!! info "Proactive Agents"
    These recipes show explicit commands and skills, but agents
    trained on the ctx playbook are **proactive**:

    They offer to: 

    * save learnings after debugging,
    * record decisions after trade-offs, 
    * create follow-up tasks after completing work, 
    * and suggest what to work on next. 

    You often **don’t** need to invoke commands directly: A simple
    "*what have we learned?*" or "*let's wrap up*" triggers the same
    workflows. 
    
    Each recipe includes a **Conversational Approach** section
    showing these natural-language patterns.

---

## Getting Started

### [Setting Up ctx Across AI Tools](multi-tool-setup.md)

Initialize ctx and configure hooks for Claude Code, Cursor,
Aider, Copilot, or Windsurf. Includes shell completion,
watch mode for non-native tools, and verification.

**Uses**: `ctx init`, `ctx hook`, `ctx agent`, `ctx completion`,
`ctx watch`, `ctx session parse`

---

## Daily Workflow

These recipes cover the workflows you’ll use every day when
working with ctx.

### [The Complete Session](session-lifecycle.md)

Walk through a full ctx session from start to finish: loading
context, picking what to work on, committing with context
capture, reflecting, and saving a snapshot.

**Uses**: `ctx status`, `ctx agent`, `ctx session save/list/load`,
`/ctx-remember`, `/ctx-next`, `/ctx-commit`, `/ctx-reflect`,
`/ctx-save`

---

### [Tracking Work Across Sessions](task-management.md)

Add, prioritize, complete, snapshot, and archive tasks. Keep
TASKS.md focused as your project evolves across dozens of
sessions.

**Uses**: `ctx add task`, `ctx complete`, `ctx tasks archive`,
`ctx tasks snapshot`, `/ctx-add-task`, `/ctx-archive`, `/ctx-next`

---

### [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md)

Record architectural decisions with rationale, capture gotchas
and lessons learned, and codify coding conventions so they
survive across sessions and team members.

**Uses**: `ctx add decision`, `ctx add learning`,
`ctx add convention`, `ctx decisions reindex`,
`ctx learnings reindex`, `/ctx-add-decision`,
`/ctx-add-learning`, `/ctx-reflect`

---

## Maintenance

### [Detecting and Fixing Drift](context-health.md)

Keep context files accurate by detecting stale paths,
reconciling with the codebase, compacting bloated files, and
auditing prompt quality.

**Uses**: `ctx drift`, `ctx sync`, `ctx compact`, `ctx status`,
`/ctx-drift`, `/ctx-status`, `/ctx-prompt-audit`

---

## History and Discovery

### [Browsing and Enriching Past Sessions](session-archaeology.md)

Export your AI session history to a browsable journal site.
Normalize rendering, enrich entries with metadata, and search
across months of work.

**Uses**: `ctx recall list/show/export`, `ctx journal site`,
`ctx serve`, `/ctx-recall`, `/ctx-journal-normalize`,
`/ctx-journal-enrich`

---

## Advanced

### [Running an Unattended AI Agent](autonomous-loops.md)

Set up a loop where an AI agent works through tasks overnight
without you at the keyboard, using ctx for persistent memory
between iterations.

This recipe shows how ctx supports long-running agent loops
without losing context or intent.

**Uses**: `ctx init --ralph`, `ctx loop`, `ctx watch`, `ctx load`,
`/ctx-loop`, `/ctx-implement`, `/ctx-context-monitor`

---

### [Turning Activity into Content](publishing.md)

Generate blog posts from project activity, write changelog
posts from commit ranges, and publish a browsable journal
site from your session history.

**Uses**: `ctx journal site`, `ctx serve`, `ctx recall export`,
`/ctx-blog`, `/ctx-blog-changelog`, `/ctx-journal-enrich`,
`/ctx-journal-normalize`
