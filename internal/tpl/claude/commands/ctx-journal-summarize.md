---
description: "Generate a summary of sessions over a time period"
---

Generate a narrative summary of sessions from a time period.

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
**FAQ**: Common questions asked and their answers

## Process

1. **Find matching journal entries**:
```!
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
- **[Feature]**: Description

## Decisions Made
- [Decision]: [Rationale]

## Learnings
- [Learning]: [Why it matters]

## Unfinished Work
- [Task]: [Status]

## Sessions
| Date | Session | Summary |
|------|---------|---------|
| 01-24 | twinkly-stirring-kettle | Code cleanup |
```

4. **Review with user** - ask about focus, omissions, filename

5. **Write to** `.context/journal/summary-[period].md`

## Alternate Formats

**Feature narrative** (`/ctx-journal-summarize how we built caching`):
```markdown
# How We Built: Caching Layer

## The Journey
[Narrative of how the feature evolved across sessions]

## Key Sessions
1. [2026-01-20] Initial exploration - evaluated Redis vs Memcached
2. [2026-01-22] Implementation - basic cache layer
3. [2026-01-25] Refinement - added TTL and invalidation

## Decisions Along the Way
- Chose Redis (session: 2026-01-20-...)
- Used write-through pattern (session: 2026-01-22-...)

## Lessons Learned
- Cache invalidation is hard (session: 2026-01-25-...)
```

**Decision trail** (`/ctx-journal-summarize decisions`):
```markdown
# Decision Trail: January 2026

| Decision | Session | Rationale |
|----------|---------|-----------|
| Use Redis | 2026-01-20-... | Speed + team experience |
| JWT auth | 2026-01-22-... | Stateless, scales better |
```

## Output

Summaries appear in the journal site alongside session entries:
- `summary-2026-01.md` (monthly)
- `summary-week-04.md` (weekly)
- `summary-auth-feature.md` (topic)
- `narrative-caching.md` (feature narrative)
- `decisions-2026-01.md` (decision trail)

## Example

User: `/ctx-journal-summarize last week`

You:
> Found 5 sessions from Jan 24-31:
> - 2x recall/export feature
> - 2x journal site
> - 1x bug fixes
>
> Focus areas to emphasize? Things to skip?

After input, generate summary and ask:
> Save to `.context/journal/summary-2026-week-05.md`?
