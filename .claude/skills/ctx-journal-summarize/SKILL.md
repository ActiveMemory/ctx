---
name: ctx-journal-summarize
description: "Summarize sessions over time. Use when needing a narrative of what happened over a period or across a feature."
---

Generate a narrative summary of sessions from a time period.

## When to Use

- When needing to understand what happened over a period
- When documenting the evolution of a feature
- When creating a decision trail

## Input

Time range (flexible):
- `last week` / `this week`
- `last month` / `January` / `2026-01`
- `2026-01-20 to 2026-01-27`
- `recent` (last 7 days)

Or topic/feature-based:
- `auth` (sessions mentioning authentication)
- `refactor` (sessions tagged as refactor)
- `how we built caching` (feature narrative across sessions)

## Summary Types

**Timeline summary** (default): What happened in a period
**Feature narrative**: "How we built X" from multiple sessions
**Decision trail**: Link decisions to the sessions that made them

## Process

1. **Find matching journal entries**:
```bash
ls .context/journal/*.md | grep "2026-01"
```

2. **Read each entry** - understand what was worked on, accomplished, decided

3. **Generate summary** with structure:

```markdown
# Summary: [Period]

## Overview
[2-3 sentences: focus and accomplishments]

## Work Completed
- **[Feature]**: Description

## Decisions Made
- [Decision]: [Rationale]

## Learnings
- [Learning]: [Why it matters]

## Sessions
| Date | Session | Summary |
|------|---------|---------|
| 01-24 | twinkly-stirring-kettle | Code cleanup |
```

4. **Review with user** - ask about focus, omissions, filename

5. **Write to** `.context/journal/summary-[period].md`
