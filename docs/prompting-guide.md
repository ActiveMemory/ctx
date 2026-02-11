---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Prompting Guide
icon: lucide/message-circle
---

![ctx](images/ctx-banner.png)

## Prompting Guide

Effective prompts for working with AI assistants in `ctx`-enabled projects.

!!! tip Help Your AI Sidekick
    AI assistants *may not* automatically read context files.

    **The right prompt triggers the right behavior**. 

This guide documents prompts that reliably produce good results.

---

## Session Start

### "Do you remember?"

Triggers the AI to silently read `TASKS.md`, `DECISIONS.md`,
`LEARNINGS.md`, and check recent history via `ctx recall` before
responding with a **structured readback**:

1. **Last session**: most recent session topic and date
2. **Active work**: pending or in-progress tasks
3. **Recent context**: 1-2 recent decisions or learnings
4. **Next step**: offer to continue or ask what to focus on

Use this at the start of every important session.

```
Do you remember what we were working on?
```

This question **implies** prior context exists. The AI checks files
rather than admitting ignorance. The expected response cites specific
context (session names, task counts, decisions), not vague summaries.

If the AI instead narrates its discovery process ("Let me check if
there are files..."), it has not loaded `CLAUDE.md` or
`AGENT_PLAYBOOK.md` properly.

### "What's the current state?"

Prompts reading of `TASKS.md`, recent sessions, and status overview.

Use this when resuming work after a break.

**Variants**:

* "Where did we leave off?"
* "What's in progress?"
* "Show me the open tasks"

---

## During Work

### "Why doesn't X work?"

This triggers **root cause analysis** rather than surface-level fixes.

Use this when something fails unexpectedly.

Framing as "*why*" encourages investigation before action. The AI will trace 
through code, check configurations, and identify the actual cause.

!!! example "Real Example"
    "Why can't I run /ctx-save?" led to discovering missing permissions
    in settings.local.json bootstrapping—a fix that benefited all users.

### "Is this consistent with our decisions?"

This prompts checking `DECISIONS.md` before implementing.

Use this before making architectural choices.

**Variants**:

* "Check if we've decided on this before"
* "Does this align with our conventions?"

### "What would break if we..."

This triggers **defensive thinking** and **impact analysis**.

Use this before making significant changes.

```
What would break if we change the Settings struct?
```

### "Before you start, read X"

This ensures specific context is loaded before work begins.

Use this when you know the relevant context exists in a specific file.

```
Before you start, check ctx recall for the auth discussion session
```

### Scope Control

Constrain the AI to prevent sprawl. These are some of the most
useful prompts in day-to-day work.

```
Only change files in internal/cli/add/. Nothing else.
```

```
No new files. Modify the existing implementation.
```

```
Keep the public API unchanged. Internal refactor only.
```

Use these when the AI tends to "helpfully" modify adjacent code,
add documentation you didn't ask for, or create new abstractions.

### Course Correction

Steer the AI when it goes off-track. Don't wait for it to finish
a wrong approach.

```
Stop. That's not what I meant. Let me clarify.
```

```
Let's step back. Explain what you're about to do before changing anything.
```

```
Undo that last change and try a different approach.
```

These work because they **interrupt momentum**. Without explicit
course correction, the AI tends to commit harder to a wrong path
rather than reconsidering.

## Reflection and Persistence

### "What did we learn?"

This prompts **reflection** on the session and often triggers adding
learnings to `LEARNINGS.md`.

Use this after completing a task or debugging session.

This is an **explicit reflection prompt**. The AI will summarize insights
and often offer to persist them.

### "Add this as a learning/decision"

This is an **explicit persistence request**.

Use this when you have discovered something worth remembering.

```
Add this as a learning: "JSON marshal escapes angle brackets by default"

# or simply.
Add this as a learning.
# and let the AI autonomously infer and summarize.
```

### "Save context before we end"

This triggers **context persistence** before the session closes.

Use it at the end of the session or before switching topics.

**Variants**:

- "Let's persist what we did"
- "Update the context files"
- `/ctx-reflect` (*Agent Skill in Claude Code*)

---

## Exploration and Research

### "Explore the codebase for X"

This triggers thorough codebase search rather than guessing.

Use this when you need to understand how something works.

This works because "**Explore**" signals that **investigation is needed**, 
not immediate action.

### "How does X work in this codebase?"

This prompts reading actual code rather than explaining general concepts.

Use this to understand the existing implementation.

```
How does session saving work in this codebase?
```

### "Find all places where X"

This triggers a **comprehensive search** across the codebase.

Use this before refactoring or understanding the impact.

---

## Meta and Process

### "What should we document from this?"

This prompts identifying learnings, decisions, and conventions
worth persisting.

Use this after complex discussions or implementations.

### "Is this the right approach?"

This invites the AI to challenge the current direction.

Use this when you want a sanity check.

This works because it allows AI to disagree. 
AIs often default to agreeing; this prompt signals you want an 
**honest assessment**.

### "What am I missing?"

This prompts thinking about edge cases, overlooked requirements,
or unconsidered approaches.

Use this before finalizing a design or implementation.

---

## CLI Commands as Prompts

Asking the AI to run `ctx` commands is itself a prompt. These
load context or trigger specific behaviors:

| Command            | What it does                                     |
|--------------------|--------------------------------------------------|
| "Run `ctx status`" | Shows context summary, file presence, staleness  |
| "Run `ctx agent`"  | Loads token-budgeted context packet              |
| "Run `ctx drift`"  | Detects dead paths, stale files, missing context |

### Agent Skills (Claude Code)

Skills are formalized prompts. Use them by name:

| Skill            | When to use                                    |
|------------------|------------------------------------------------|
| `/ctx-reflect`   | Structured reflection checkpoint               |
| `/ctx-recall`    | Browse session history for past discussions    |
| `/ctx-status`    | Quick context summary                          |
| `/ctx-agent`     | Load full context packet                       |
| `/consolidate`   | Detect and fix code-level drift                |
| `/update-docs`   | Sync docs and conventions after code changes   |
| `/verify`        | Verify before claiming work is complete        |
| `/qa`            | Run QA checks before committing                |

Skills combine a prompt, tool permissions, and domain knowledge
into a single invocation. See [Integrations](integrations.md) for
setup details.

---

## Anti-Patterns

Based on our `ctx` development experience (*i.e., "sipping our own champagne"*)
so far, here are some prompts that tend to produce poor results:

| Prompt                 | Problem                       | Better Alternative                        |
|------------------------|-------------------------------|-------------------------------------------|
| "Fix this"             | Too vague, may patch symptoms | "Why is this failing?"                    |
| "Make it work"         | Encourages quick hacks        | "What's the right way to solve this?"     |
| "Just do it"           | Skips planning                | "Plan this, then implement"               |
| "You should remember"  | Confrontational               | "Do you remember?"                        |
| "Obviously..."         | Discourages questions         | State the requirement directly            |
| "Idiomatic X"          | Triggers language priors      | "Follow project conventions"              |
| "Implement everything" | No phasing, sprawl risk       | Break into tasks, implement one at a time |
| "You should know this" | Assumes context is loaded     | "Before you start, read X"                |

---

## Quick Reference

| Goal            | Prompt                                   |
|-----------------|------------------------------------------|
| Load context    | "Do you remember?"                       |
| Resume work     | "What's the current state?"              |
| Debug           | "Why doesn't X work?"                    |
| Validate        | "Is this consistent with our decisions?" |
| Impact analysis | "What would break if we..."              |
| Reflect         | "What did we learn?"                     |
| Persist         | "Add this as a learning"                 |
| Explore         | "How does X work in this codebase?"      |
| Sanity check    | "Is this the right approach?"            |
| Completeness    | "What am I missing?"                     |
| Constrain scope | "Only change files in X. Nothing else."  |
| Course correct  | "Stop. That's not what I meant."         |
| Check health    | "Run `ctx drift`"                        |
| Save context    | `/ctx-reflect`                           |

---

## Explore → Plan → Implement

For non-trivial work, name the phase you want:

```
Explore src/auth and summarize the current flow.
Then propose a plan. After I approve, implement with tests.
```

This prevents the AI from jumping straight to code. The three phases
map to different modes of thinking:

- **Explore**: read, search, understand — no changes
- **Plan**: propose approach, trade-offs, scope — no changes
- **Implement**: write code, run tests, verify — changes

Small fixes skip straight to implement. Complex or uncertain work
benefits from all three.

---

## Prompts by Task Type

Different tasks need different prompt structures. The pattern:
**symptom + location + verification**.

### Bugfix
```
Users report search returns empty results for queries with hyphens.
Reproduce in src/search/. Write a failing test for "foo-bar",
fix the root cause, run: go test ./internal/search/...
```

### Refactor
```
Inspect src/auth/ and list duplication hotspots.
Propose a refactor plan scoped to one module.
After approval, remove duplication without changing behavior.
Add a test if coverage is missing. Run: make audit
```

### Research
```
Explore the request flow around src/api/.
Summarize likely bottlenecks with evidence.
Propose 2-3 hypotheses. Do not implement yet.
```

### Docs
```
Update docs/cli-reference.md to reflect the new --format flag.
Confirm the flag exists in the code and the example works.
```

Notice each prompt includes **what to verify and how**. Without that,
you get "should work now" instead of evidence.

---

## Writing Tasks as Prompts

Tasks in `TASKS.md` are **indirect prompts** to the AI. How you write them
shapes how the AI approaches the work.

### State the Deliverable, Not Just Steps

Bad task (implementation-focused):
```markdown
- [ ] T1.1.0: Parser system
  - [ ] Define data structures
  - [ ] Implement line parser
  - [ ] Implement session grouper
```

The AI may complete all subtasks but miss the actual goal. What does
"Parser system" deliver to the user?

Good task (deliverable-focused):
```markdown
- [ ] T1.1.0: Parser CLI command
  **Deliverable**: `ctx recall list` command that shows parsed sessions
  - [ ] Define data structures
  - [ ] Implement line parser
  - [ ] Implement session grouper
```

Now the AI knows the subtasks serve a specific user-facing deliverable.

### Use Acceptance Criteria

For complex tasks, add explicit "done when" criteria:

```markdown
- [ ] T2.0: Authentication system
  **Done when**:
  - [ ] User can register with email
  - [ ] User can log in and get a token
  - [ ] Protected routes reject unauthenticated requests
```

This prevents premature "task complete" when only the implementation
details are done but the feature doesn't actually work.

### Subtasks ≠ Parent Task

Completing all subtasks does **not** mean the parent task is complete.

The parent task describes **what** the user gets.
Subtasks describe **how** to build it.

Always re-read the parent task description before marking it complete.
Verify the stated deliverable exists and works.

---

## Contributing

Found a prompt that works well?
[Open an issue](https://github.com/ActiveMemory/ctx/issues) or PR with:

1. The prompt text
2. What behavior it triggers
3. When to use it
4. Why it works (*optional but helpful*)
