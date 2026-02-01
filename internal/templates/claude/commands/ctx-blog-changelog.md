---
description: "Generate a blog post from a commit range with a theme"
---

Generate a blog post about changes since a specific commit, with a given theme.

## Input

Required:
- **Commit hash**: Starting point (e.g., `040ce99`, `HEAD~50`)
- **Theme**: The narrative angle (e.g., "human-assisted refactoring", "the recall system")

Optional:
- **Reference post**: An existing post to match the style

## Usage Examples

```text
/ctx-blog-changelog 040ce99 "human-assisted refactoring"
/ctx-blog-changelog HEAD~30 "building the journal system" --like 2026-01-27-building-ctx-using-ctx.md
/ctx-blog-changelog v0.1.0 "what's new in v0.2.0"
```

## Process

1. **Analyze the commit range**:
```!
git log --oneline <commit>..HEAD
git diff --stat <commit>..HEAD
git log --format="%s" <commit>..HEAD | head -50
```

2. **Gather supporting context**:
```!
# Decisions made in this period
git log --format="%H %s" <commit>..HEAD | grep -i decision

# Files most changed
git diff --stat <commit>..HEAD | sort -t'|' -k2 -rn | head -20

# Journal entries from this period (if available)
ls .context/journal/*.md
```

3. **If reference post provided**, read it:
```!
cat docs/blog/<reference-post>.md
```

4. **Draft the narrative** following the theme:

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

### [Sub-theme 1]
[Changes, examples, code snippets]

### [Sub-theme 2]
[Changes, examples, code snippets]

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

## Example

User: `/ctx-blog-changelog 040ce99 "human-assisted refactoring"`

You:
> Analyzing 040ce99..HEAD (47 commits)
>
> I see the theme of "human-assisted refactoring" in:
> - Consolidation of magic strings → config constants
> - Package reorganization (init → initialize)
> - Test colocation
> - Convention documentation
>
> Reference: The existing post discusses YOLO vs human-guided modes.
> I'll frame this as a continuation: "The Refactoring Sessions: What Happened After YOLO Mode"
>
> Outline:
> 1. Where YOLO left us (the debt)
> 2. The consolidation pattern (constants, naming)
> 3. Before/after comparisons
> 4. What the AI couldn't see (why human guidance mattered)
>
> Does this angle work? Any commits or changes to highlight?
