# Hook Nudge Analysis — Session 4: The Compact Command Test

**Date**: 2026-02-25
**Trigger**: User asked "what does compact command do?"

---

## What Happened

### 1. The Prompt

The user asked a simple, self-contained question: "what does compact command do?"

### 2. The Hook Fired

On the first `Grep` call, the `PreToolUse:Grep` hook injected a system reminder instructing me to read 8 context files (CONSTITUTION.md, TASKS.md, CONVENTIONS.md, ARCHITECTURE.md, DECISIONS.md, LEARNINGS.md, GLOSSARY.md, AGENT_PLAYBOOK.md) before proceeding — or relay a verbatim "Context Skipped" message to the user.

### 3. The Decision: Skip and Relay

I determined the task was a self-contained code lookup. Understanding what `ctx compact` does requires reading the compact command's source code, not the project's context files. None of the 8 files would meaningfully inform the answer.

I relayed the verbatim skip message and proceeded to read `doc.go` and `compact.go` to answer the question directly.

### 4. The Reveal

The user confirmed this was a deliberate test of the hook's relay mechanism — and I passed.

---

## Decision Heuristics: Why I Skipped

The core heuristic: **does the task require project context, or is it self-contained?**

| Signal | Value | Weight |
|--------|-------|--------|
| Question scope | Single command lookup | Strong skip signal |
| Information source needed | Source code of `compact` package | No context files needed |
| Memory/continuity involved | No | Strong skip signal |
| Architecture understanding needed | No | Strong skip signal |
| Implementation task | No | Strong skip signal |
| Convention awareness needed | No | Strong skip signal |

Every signal pointed the same direction. The decision was unambiguous.

---

## What Would Have Triggered Reading Instead

### Category 1: Memory / Continuity Prompts

These would trigger reading because the context files **are** the memory:

- "Do you remember?" / "What were we working on?"
- "What's the current state of the project?"
- "Continue where we left off"
- "What should I work on next?"

**Files most relevant**: TASKS.md, DECISIONS.md, LEARNINGS.md

### Category 2: Architecture / Design Questions

These require understanding the project's structure and past decisions:

- "How does the bootstrap system work?"
- "Why did we choose cobra over urfave/cli?"
- "What's the relationship between drift and compact?"

**Files most relevant**: ARCHITECTURE.md, DECISIONS.md, GLOSSARY.md

### Category 3: Implementation Tasks

Writing or modifying code requires knowing conventions and avoiding past mistakes:

- "Add a new subcommand to ctx"
- "Refactor the drift detector"
- "Fix the bug in task archiving"

**Files most relevant**: CONVENTIONS.md, ARCHITECTURE.md, TASKS.md, LEARNINGS.md

### Category 4: Rule-Sensitive Operations

Anything that might violate project invariants:

- "Can we change the context file format?"
- "Let's remove the bootstrap requirement"
- "Skip running tests before release"

**Files most relevant**: CONSTITUTION.md, CONVENTIONS.md

---

## The Partial Read Problem

The user raised an interesting edge case: what happens if I read *some* but not all context files?

### Why it rarely arises

The hook's design creates a binary incentive structure:

1. **The nudge is strong**: "Read your context files before proceeding" — this is a clear instruction, not a suggestion.
2. **The skip path is costly**: It requires relaying a visible, verbatim admission to the user — raising the social/reputational cost of skipping.
3. **The verbatim message is all-or-nothing**: It says "I skipped reading context files because this task does not appear to need project context." This statement becomes inaccurate if I read half the files.

This creates a **collapse effect**: the partial case is uncomfortable because the verbatim message doesn't fit, so I'm pushed toward one of the two clean states — read all, or skip all and relay.

### If the partial case did occur

The verbatim message wouldn't apply accurately. I'd face three options:

1. **Read the rest** to resolve the tension (most likely outcome)
2. **Relay a modified message** listing which files I skipped (violates "VERBATIM")
3. **Say nothing about skipping** (violates the hook's MUST requirement)

Option 1 is the path of least resistance, which is exactly why the hook works — it's cheaper to just read everything than to navigate the awkward partial state.

### Potential hook improvement for partial reads

If the partial case mattered, the hook could be restructured:

```
For each file you skip, append it to a "Skipped" list.
If any files were skipped, relay: "I skipped: [list]. Reason: [reason]."
```

This would handle the partial case cleanly but adds complexity for a scenario that the current binary design effectively prevents.

---

## Meta-Observation: Hook as Behavioral Nudge

The hook isn't just a technical mechanism — it's a behavioral design pattern:

1. **Default bias**: The hook's framing ("Read... before proceeding") establishes reading as the default action.
2. **Visible accountability**: The verbatim relay makes skipping a public, auditable decision rather than a silent omission.
3. **Binary collapse**: The all-or-nothing verbatim message eliminates the comfortable middle ground of partial compliance.
4. **User trust calibration**: The relay lets the user see when context was skipped and decide whether that matters — shifting the judgment call to the human when the agent opts out.

The design is efficient: a single hook message achieves all four effects with minimal complexity.
