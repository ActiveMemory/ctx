---
description: "Generate a blog post draft from recent activity"
---

Generate a blog post draft from recent project activity.

## Input

The user may specify:
- A time range: `last week`, `since Monday`, `January`
- A topic focus: `the refactoring`, `new features`, `lessons learned`
- Or just run it to analyze recent activity

## Sources to Analyze

Gather context from multiple sources:

```!
# Recent commits
git log --oneline -30

# Recent decisions
head -100 .context/DECISIONS.md

# Recent learnings
head -100 .context/LEARNINGS.md

# Recent tasks completed
grep -A2 "\[x\]" .context/TASKS.md | head -50

# Journal entries (if available)
ls -la .context/journal/*.md 2>/dev/null | tail -10
```

## Blog Post Structure

Follow the style of existing posts in `docs/blog/`:

```markdown
---
title: "Descriptive Title: What This Post Is About"
date: YYYY-MM-DD
author: [Ask user]
---

*Author / Date*

# Title

> Opening hook or question

[Introduction: Set the scene, why this matters]

## Section 1: The Context/Problem
[What situation led to this work]

## Section 2: What We Did
[Narrative of the work, with code examples]

## Section 3: What We Learned
[Key insights, gotchas, patterns discovered]

## Section 4: What's Next
[Future work, open questions]

---

*Footer note about how the post was created*
```

## Style Guidelines

Based on `docs/blog/2026-01-27-building-ctx-using-ctx.md`:

- **Personal voice**: Use "I", "we", share the journey
- **Show don't tell**: Include actual code, commits, quotes from sessions
- **Tables for comparisons**: Before/after, patterns found
- **Callout boxes**: For key insights (`!!! info`, `!!! quote`)
- **Honest about failures**: Include what went wrong and why
- **Concrete examples**: Reference specific files, commits, decisions

## Process

1. Gather sources (git, decisions, learnings, journals)
2. Identify the narrative arc (what's the story?)
3. Draft outline for user approval
4. Write full draft
5. Ask for revisions
6. Save to `docs/blog/YYYY-MM-DD-slug.md`

## Example

User: `/ctx-blog last week`

You:
> Looking at the last week, I see:
> - 15 commits focused on journal/recall features
> - 3 new decisions (project matching, tool output labeling)
> - 5 learnings (session parsing, git remote matching)
>
> I'd frame this as: "Building a Memory Palace: How ctx Learned to Remember Sessions"
>
> Narrative arc:
> 1. The problem: sessions scattered, hard to find
> 2. The solution: recall export + journal site
> 3. The gotchas: timezone bugs, missing sessions
> 4. The insight: git remote as stable project ID
>
> Sound good? Any angle you want to emphasize?
