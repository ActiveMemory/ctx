---
name: ctx-blog-changelog
description: "Generate themed blog post from commits. Use when writing about changes between releases or documenting a development arc."
---

Generate a blog post about changes since a specific commit, with a given theme.

## Before Writing

Two questions; if any answer is "no", reconsider:

1. **"Is there enough change to tell a story?"** → A handful of typo
   fixes doesn't warrant a post
2. **"Is the theme clear?"** → If the commit range covers unrelated
   work, narrow the scope or split into multiple posts

## When to Use

- When documenting changes between releases
- When writing about a development arc or theme
- When the user wants to explain "what changed and why"

## When NOT to Use

- For general project updates without a commit range (use `/ctx-blog`)
- When the changes are minor or routine maintenance
- When there's no unifying theme across the commits

## Input

Required:
- **Commit hash**: Starting point (e.g., `040ce99`, `HEAD~50`, `v0.1.0`)
- **Theme**: The narrative angle (e.g., "human-assisted refactoring",
  "the recall system")

Optional:
- **Reference post**: An existing post to match the style

## Usage Examples

```text
/ctx-blog-changelog 040ce99 "human-assisted refactoring"
/ctx-blog-changelog HEAD~30 "building the journal system"
/ctx-blog-changelog v0.1.0 "what's new in v0.2.0"
```

## Process

1. **Analyze the commit range**:
```bash
git log --oneline <commit>..HEAD
git diff --stat <commit>..HEAD
git log --format="%s" <commit>..HEAD | head -50
```

2. **Gather supporting context**:
```bash
# Files most changed
git diff --stat <commit>..HEAD | sort -t'|' -k2 -rn | head -20

# Journal entries from this period
ls .context/journal/*.md
```

3. **Draft the narrative** following the theme
4. Save to `docs/blog/YYYY-MM-DD-slug.md`
5. **Update `docs/blog/index.md`** with an entry at the top:

```markdown
### [Post Title](YYYY-MM-DD-slug.md)

*Author / Date*

2-3 sentence blurb.

**Topics**: topic-one, topic-two, topic-three

---
```

## Blog Structure

### Frontmatter

```yaml
---
title: "[Theme]: [Specific Angle]"
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
# [Title]

![ctx](../images/ctx-banner.png)

> [Hook related to theme]

## The Starting Point
[State of codebase at <commit>, what prompted the change]

## The Journey
[Narrative of changes, organized by theme not chronology]

## Before and After
[Comparison table or code diff showing improvement]

## Key Commits

| Commit | Change      |
|--------|-------------|
| abc123 | Description |

## Lessons Learned
[Insights from this work]

## What's Next
[Future work enabled by these changes]
```

## Style Guidelines

- **Personal voice**: Use "I", "we", share the journey
- **Show don't tell**: Include actual code, commits, diffs
- **Tables for comparisons**: Before/after, key commits
- **Honest about failures**: Include what went wrong and why
- **Concrete examples**: Reference specific files, commits, decisions
- **No em-dashes**: Use `:`, `;`, or restructure the sentence instead
- **Straight quotes only**: Use "dumb quotes" (`"`, `'`), never
  typographic/curly quotes
- **80-character line width**: Wrap prose at ~80 characters; exceptions
  for tables, code blocks, and URLs
