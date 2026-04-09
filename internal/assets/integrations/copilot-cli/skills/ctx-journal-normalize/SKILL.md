---
name: ctx-journal-normalize
description: "Normalize journal source markdown for clean rendering. Use after journal site shows rendering issues: fence nesting, metadata formatting, broken lists."
tools: [bash, read, write, edit]
---

Reconstruct journal entries as clean markdown from stripped plain text.

## When to Use

- After `ctx journal site` shows rendering issues
- When journal entries have fence nesting problems
- When metadata blocks render as raw `**Key**: value`
- Before running `ctx-journal-enrich` (clean markdown improves extraction)

## When NOT to Use

- On entries already normalized (check `.state.json`)
- When the site renders correctly
- On non-journal markdown files

## Output Rules

1. **Fences**: Always use backtick fences. Innermost code gets
   3 backticks. Each nesting level adds 1.
2. **Metadata**: `**Key**: value` blocks become collapsed `<details>`.
3. **Tool output**: Collapse into `<details>` when > 10 lines.
4. **Lists**: 2-space indent per level.
5. **No invented content**: Every word in output traces to input.

## Process

1. **Backup first**: copy journal directory to `.bak` sibling
2. Identify files to normalize (skip already-normalized via `.state.json`)
3. Process files turn-by-turn (not whole file at once)
4. Write back the fixed files
5. Mark normalized: `ctx system mark-journal <filename> normalized`
6. Regenerate site: `ctx journal site --build`
7. Report what changed

## Quality Checklist

- [ ] Backup created before modifying
- [ ] Already-normalized files skipped
- [ ] No content was invented or lost
- [ ] State file updated for processed entries
