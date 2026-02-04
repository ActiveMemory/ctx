---
name: ctx-blog-changelog
description: "Generate themed blog post from commits. Use when writing about changes between releases or documenting a development arc."
---

Generate a blog post about changes since a specific commit, with a given theme.

## When to Use

- When documenting changes between releases
- When writing about a development arc or theme
- When the user wants to explain "what changed and why"

## Input

Required:
- **Commit hash**: Starting point (e.g., `040ce99`, `HEAD~50`, `v0.1.0`)
- **Theme**: The narrative angle (e.g., "human-assisted refactoring", "the recall system")

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

## Blog Structure

```markdown
---
title: "[Theme]: [Specific Angle]"
date: YYYY-MM-DD
author: [Ask user]
---

# [Title]

> [Hook related to theme]

## The Starting Point
[State of codebase at <commit>, what prompted the change]

## The Journey
[Narrative of changes, organized by theme not chronology]

## Before and After
[Comparison table or code diff showing improvement]

## Key Commits
| Commit | Change |
|--------|--------|
| abc123 | Description |

## Lessons Learned
[Insights from this work]

## What's Next
[Future work enabled by these changes]
```
