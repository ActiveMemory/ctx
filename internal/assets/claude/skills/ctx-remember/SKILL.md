---
name: ctx-remember
description: "Recall project context and present structured readback. Use when the user asks 'do you remember?', at session start, or when context seems lost."
allowed-tools: Bash(ctx:*), Read
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
- When the user asks about previous work, decisions, or learnings

## When NOT to Use

- Context was already loaded this session via `/ctx-agent`: don't
  re-fetch what you already have
- Mid-session when you are actively working on a task and context
  is fresh: don't interrupt flow
- When the user is asking about a *specific* past session by name
  or ID: use `/ctx-history` instead, which has list/show/export
  subcommands

## Process

Do all of this **silently**: narrating the steps makes the readback
feel like a file search rather than genuine recall:

1. **Load context packet**:
   ```bash
   ctx agent
   ```
2. **Read the files** listed in the packet's "Read These Files"
   section (TASKS.md, DECISIONS.md, LEARNINGS.md, etc.)
3. **List recent sessions**:
   ```bash
   ctx journal source --limit 3
   ```
4. **Present the structured readback** (see format below)

## Readback Format

Present your findings as a structured readback with these sections:

**Last session**: Topic, date, and what was accomplished. Cite the
most recent session from the session list.

**Active work**: Pending and in-progress tasks from TASKS.md. Use
a brief list: one line per task with its status.

**Recent context**: 1-2 recent decisions or learnings that are
relevant. Pick the most recent or most impactful.

**Next step**: Suggest what to work on next based on the active
tasks, or ask the user for direction if priorities are unclear.

## Readback Rules

- Open directly with the readback: instead of "I don't have memory",
  present what you found
- Skip preamble like "Let me check": go straight to the structured
  readback
- Present findings as recall, not discovery: you are *remembering*,
  not *searching*
- Be honest about the mechanism only if the user explicitly asks
  *how* you remember (e.g., "It's stored in context files managed
  by ctx")

## Examples

### Good Readback

> **Last session** (2026-02-07): We implemented the cooldown
> mechanism for `ctx agent` to prevent redundant context loads.
>
> **Active work**:
> - [ ] Add `--format json` flag to `ctx status` (pending)
> - [x] Implement session cooldown (done)
> - [ ] Write integration tests for journal import (in progress)
>
> **Recent context**:
> - Decided to use file-based cooldown tokens instead of
>   environment variables (simpler, works across shells)
> - Learned that Claude Code hooks run in a subprocess, so env
>   vars set in hooks don't persist to the main session
>
> **Next step**: The integration tests for journal import are
> partially done. Want to continue those, or shift to the JSON
> status flag?

### Bad Readback (Anti-patterns)

> "I don't have persistent memory, but let me check if there
> are any context files..."

> "Let me look at the context files to see what's there.
> I found TASKS.md, let me read it..."

> "I found some session files. Here's what they contain..."

## Companion Tool Check

After presenting the readback, check companion tool availability.
Skip this section entirely if `companion_check: false` is set in
`.ctxrc`: check by running `ctx config status` and looking for
the field value.

**Companion tools** enhance ctx skills with web search and code
intelligence. They are optional but recommended:

| Tool          | Purpose                                                | Smoke test                                                           |
|---------------|--------------------------------------------------------|----------------------------------------------------------------------|
| Gemini Search | Grounded web search with citations                     | Call `mcp__gemini-search__search_with_grounding` with a simple query |
| GitNexus      | Code knowledge graph (symbols, blast radius, clusters) | Call `mcp__gitnexus__list_repos`                                     |

**Check procedure:**

1. Attempt each smoke test silently
2. For tools that respond: note as available (no output needed)
3. For tools that fail or are not connected: append a brief note
   after the readback:
   > "Companion tools: Gemini Search is not connected (web search
   > will fall back to built-in). Install via MCP settings if
   > needed."
4. For GitNexus specifically: if it responds but the current repo
   is not indexed or the index is stale, suggest:
   > "GitNexus index is stale: run `npx gitnexus analyze` to
   > rehydrate."

Present companion status as a one-line note after the readback,
not a separate section. If everything is healthy, say nothing.

## Quality Checklist

Before presenting the readback, verify:
- [ ] Context packet was loaded (not skipped)
- [ ] Files from the read order were actually read
- [ ] Structured readback has all four sections
- [ ] No narration of the discovery process leaked into output
- [ ] Readback feels like recall, not a file system tour
- [ ] Companion tool check ran (unless suppressed via .ctxrc)
