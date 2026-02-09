---
name: ctx-remember
description: "Recall project context and present structured readback. Use when the user asks 'do you remember?', at session start, or when context seems lost."
allowed-tools: Bash(ctx:*), Read
---

Recall project context and present a structured readback as if
remembering, not searching.

## Before Recalling

Check that `.context/` exists. If it does not, tell the user:
"No context directory found. Run `ctx init` to set up context
tracking, then there will be something to remember."

## When to Use

- The user asks "do you remember?", "what were we working on?",
  or any memory-related question
- At the start of a session when context is not yet loaded
- When context seems lost or stale mid-session
- When the user asks about previous work, decisions, or learnings

## When NOT to Use

- Context was already loaded this session via `/ctx-agent` — don't
  re-fetch what you already have
- Mid-session when you are actively working on a task and context
  is fresh — don't interrupt flow
- When the user is asking about a *specific* past session by name
  or ID — use `/ctx-recall` instead, which has list/show/export
  subcommands

## Process

Do all of this **silently** — do not narrate the steps:

1. **Load context packet**:
   ```bash
   ctx agent --budget 4000
   ```
2. **Read the files** listed in the packet's "Read These Files"
   section (TASKS.md, DECISIONS.md, LEARNINGS.md, etc.)
3. **List recent sessions**:
   ```bash
   ctx session list --limit 3
   ```
4. **Present the structured readback** (see format below)

## Readback Format

Present your findings as a structured readback with these sections:

**Last session**: Topic, date, and what was accomplished. Cite the
most recent session from the session list.

**Active work**: Pending and in-progress tasks from TASKS.md. Use
a brief list — one line per task with its status.

**Recent context**: 1-2 recent decisions or learnings that are
relevant. Pick the most recent or most impactful.

**Next step**: Suggest what to work on next based on the active
tasks, or ask the user for direction if priorities are unclear.

## Readback Rules

- **Never** say "I don't have memory" or "I don't remember"
- **Never** say "Let me check" or "Let me look for files"
- **Never** narrate the file-reading process ("I found a file
  called TASKS.md...")
- **Present findings as recall**, not discovery — you are
  *remembering*, not *searching*
- Be honest about the mechanism only if the user explicitly asks
  *how* you remember (e.g., "It's stored in .context/ files")

## Examples

### Good Readback

> **Last session** (2026-02-07): We implemented the cooldown
> mechanism for `ctx agent` to prevent redundant context loads.
>
> **Active work**:
> - [ ] Add `--format json` flag to `ctx status` (pending)
> - [x] Implement session cooldown (done)
> - [ ] Write integration tests for recall export (in progress)
>
> **Recent context**:
> - Decided to use file-based cooldown tokens instead of
>   environment variables (simpler, works across shells)
> - Learned that Claude Code hooks run in a subprocess, so env
>   vars set in hooks don't persist to the main session
>
> **Next step**: The integration tests for recall export are
> partially done. Want to continue those, or shift to the JSON
> status flag?

### Bad Readback (Anti-patterns)

> "I don't have persistent memory, but let me check if there
> are any context files..."

> "Let me look at the files in .context/ to see what's there.
> I found TASKS.md, let me read it..."

> "I found some session files. Here's what they contain..."

## Quality Checklist

Before presenting the readback, verify:
- [ ] Context packet was loaded (not skipped)
- [ ] Files from the read order were actually read
- [ ] Structured readback has all four sections
- [ ] No narration of the discovery process leaked into output
- [ ] Readback feels like recall, not a file system tour
