---
name: ctx-journal-normalize
description: "Normalize journal source markdown for clean rendering. Use after `ctx journal site` shows rendering issues: fence nesting, metadata formatting, broken lists."
---

Reconstruct journal entries as clean markdown from stripped plain text.

## Before Normalizing

1. **Export first**: run `ctx recall export --all` (or `make journal`)
   so there are journal files to normalize
2. **Check if needed**: if the journal site renders cleanly, skip this

## When to Use

- After `ctx journal site` shows rendering issues
- When journal entries have fence nesting problems (no code highlighting)
- When metadata blocks render as raw `**Key**: value` instead of tables
- Before running `/ctx-journal-enrich` (clean markdown improves extraction)

## When NOT to Use

- On entries already marked with `<!-- normalized -->` (they're done)
- When the site renders correctly (don't fix what isn't broken)
- On non-journal markdown files (this is journal-specific)

## Usage Examples

```text
/ctx-journal-normalize
/ctx-journal-normalize twinkly-stirring-kettle
/ctx-journal-normalize --compact
```

## Architecture

`ctx journal site` strips all code fence markers from site copies,
eliminating nesting conflicts. The result is readable plain text with
structural markers preserved. This skill goes further: it reconstructs
proper markdown in **source files** so the site renders with code
highlighting and proper formatting.

## Input Format

Source journal entries have these structural markers:
- Turn boundaries: `### N. Role (HH:MM:SS)`
- Tool calls: `Tool Name: args` on their own line
- Tool output: block following a Tool Output turn header
- Section breaks: `---`
- Frontmatter: YAML between `---` delimiters at file start

## Output Rules

1. **Fences**: Always use **backtick** fences, never tildes.
   Innermost code gets 3 backticks. Each nesting level adds 1.
   Never nest same-count fences.
2. **Metadata**: `**Key**: value` blocks become collapsed `<details>`
   with `<table>`. Summary from Date/Duration/Turns/Model.
3. **Tool output**: Collapse into
   `<details><summary>N lines</summary>` when > 10 lines.
4. **Lists**: 2-space indent per level. Continuation lines match
   list item indent.
5. **No invented content**: Every word in output must trace to input.
   Structure changes only.

## Modes

**Default (lossless)**: Reformat only. All content preserved.

**Compact** (when user requests `--compact` or "compact"): May
summarize tool outputs > 50 lines; keep first/last 5 lines with
`... (N lines omitted)`. Flag this to user before proceeding.

## Process

1. **Backup first**: `cp -r .context/journal/ .context/journal.bak/`
   - Always back up before modifying; files may contain user edits
   - Tell the user where the backup is
2. Identify files to normalize:
   - If user specifies a file/pattern, use that
   - Otherwise, scan `.context/journal/*.md`
   - **Skip files with** `<!-- normalized: YYYY-MM-DD -->` marker
3. Process turn-by-turn (not whole file at once; large files blow
   context):
   - Fix fence nesting, metadata, lists per output rules
   - Add `<!-- normalized: YYYY-MM-DD -->` after frontmatter
4. Write back the fixed files
5. Regenerate site: `ctx journal site --build`
6. Report what changed and remind user of backup location

## Idempotency

Two markers, two passes:

- `<!-- normalized: YYYY-MM-DD -->`: metadata tables done.
  Skip metadata conversion on re-run.
- `<!-- fences-verified: YYYY-MM-DD -->`: fence reconstruction done.
  `stripFences` in `ctx journal site` only skips files with this
  marker.

Files without `fences-verified` get all fences stripped in site copies
(readable but no code highlighting). Add this marker after verifying
fence nesting is correct.

## Scope

- Operate on **source files** (`.context/journal/`)
- Changes persist; no repeated normalization needed
- Preserve all substantive content; only fix formatting
