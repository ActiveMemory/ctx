---
name: ctx-blog
description: "Generate blog post draft. Use when documenting project progress, sharing learnings, or writing about development experience."
---

Generate a blog post draft from recent project activity.

## Before Writing

Two questions — if any answer is "no", reconsider:

1. **"Is there a narrative arc?"** → A blog post needs a story (problem →
   approach → outcome), not just a list of changes
2. **"Would someone outside the project learn something?"** → If the
   insight is only useful internally, use LEARNINGS.md instead

## When to Use

- When documenting significant project progress
- When sharing learnings publicly
- When the user wants to write about the development experience

## When NOT to Use

- For internal-only notes (use session saves or LEARNINGS.md)
- When the work is still in progress with no clear insight yet
- For changelogs (use `/ctx-blog-changelog` instead)

## Input

The user may specify:
- A time range: `last week`, `since Monday`, `January`
- A topic focus: `the refactoring`, `new features`, `lessons learned`
- Or just run it to analyze recent activity

## Sources to Analyze

Gather context from multiple sources:

```bash
# Recent commits
git log --oneline -30

# Recent decisions
ctx status --verbose  # or read DECISIONS.md directly

# Recent learnings
ctx status --verbose  # or read LEARNINGS.md directly

# Recent tasks completed
ctx status  # shows active and completed task counts

# Journal entries (if available)
ctx recall list --limit 10
```

## Blog Post Structure

### Frontmatter

```yaml
---
title: "Descriptive Title: What This Post Is About"
date: YYYY-MM-DD
author: [Ask user]
topics:
  - topic-one
  - topic-two
  - topic-three
---
```

### Body

```markdown
# Title

![ctx](../images/ctx-banner.png)

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
```

## Style Guidelines

- **Personal voice**: Use "I", "we", share the journey
- **Show don't tell**: Include actual code, commits, quotes
- **Tables for comparisons**: Before/after, patterns found
- **Honest about failures**: Include what went wrong and why
- **Concrete examples**: Reference specific files, commits, decisions
- **No em-dashes**: Use `:`, `;`, or restructure the sentence instead
- **Straight quotes only**: Use "dumb quotes" (`"`, `'`), never
  typographic/curly quotes
- **80-character line width**: Wrap prose at ~80 characters; exceptions
  for tables, code blocks, and URLs

## Process

1. Gather sources (git, decisions, learnings, journals)
2. Identify the narrative arc (what's the story?)
3. Draft outline for user approval
4. Write full draft
5. Ask for revisions
6. Save to `docs/blog/YYYY-MM-DD-slug.md`
7. **Update `docs/blog/index.md`** — add entry at the top following the
   existing pattern:

```markdown
### [Post Title](YYYY-MM-DD-slug.md)

*Author / Date*

2-3 sentence blurb.

**Topics**: topic-one, topic-two, topic-three

---
```

## Example Invocations

```
/ctx-blog about the cooldown feature we just built
/ctx-blog last week's refactoring work
/ctx-blog lessons learned from hook design
```
