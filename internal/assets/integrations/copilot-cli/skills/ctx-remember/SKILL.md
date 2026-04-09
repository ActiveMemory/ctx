---
name: ctx-remember
description: "Recall project context and present structured readback. Use when the user asks 'do you remember?', at session start, or when context seems lost."
tools: [bash, read]
---

Recall project context and present a structured readback as if
remembering, not searching.

## Before Recalling

Check that the context directory exists. If it does not, tell the
user: "No context directory found. Run `ctx init` to set up context
tracking, then there will be something to remember."

## When to Use

- The user asks "do you remember?", "what were we working on?",
  or any memory-related question
- At the start of a session when context is not yet loaded
- When context seems lost or stale mid-session

## When NOT to Use

- Context was already loaded this session: don't re-fetch
- Mid-session when actively working and context is fresh
- When asking about a specific past session by name: use
  `ctx-recall` instead

## Process

Do all of this **silently** — no narration of the steps:

1. **Load context packet**:
   ```bash
   ctx agent --budget 4000
   ```
2. **Read the files** listed in the packet's "Read These Files"
   section (TASKS.md, DECISIONS.md, LEARNINGS.md, etc.)
3. **List recent sessions**:
   ```bash
   ctx recall list --limit 3
   ```
4. **Present the structured readback**

## Readback Format

**Last session**: Topic, date, and what was accomplished.

**Active work**: Pending and in-progress tasks from TASKS.md.

**Recent context**: 1-2 recent decisions or learnings.

**Next step**: Suggest what to work on next or ask for direction.

## Readback Rules

- Open directly with the readback: not "I don't have memory"
- Skip preamble like "Let me check": go straight to readback
- Present findings as recall, not discovery
- Be honest about the mechanism only if explicitly asked

## Quality Checklist

- [ ] Context packet was loaded
- [ ] Files from the read order were actually read
- [ ] Structured readback has all four sections
- [ ] No narration of the discovery process
- [ ] Readback feels like recall, not a file system tour
