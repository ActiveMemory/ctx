---
title: Scratchpad tag filtering
date: 2026-04-06
status: ready
---

# Scratchpad Tag Filtering

## Problem

Users organically tag pad entries with `#word` tokens (e.g., `#later`,
`#urgent`, `#idea`) and filter them with shell pipelines:

```bash
ctx pad | grep -v "#later"
ctx pad edit 79 --append "#later"
```

This works but is fragile: grep matches line numbers containing the
digit sequence, doesn't understand entry boundaries, and can't negate
cleanly. The pad should support tags natively since users are already
using them.

## Solution

Convention-based tags: a `#word` token anywhere in an entry's text is a
tag. No new storage format, no metadata layer. Entries already containing
`#word` tokens work retroactively.

Three additions:

1. **`--tag` filter flag** on the root `ctx pad` list command
2. **`--tag` filter flag** on `ctx pad show` (filter then pick by index)
3. **`ctx pad tags` subcommand** listing all tags with counts

## Design

### Tag Definition

A tag is a `#` followed by one or more word characters (`[a-zA-Z0-9_-]`),
bounded by whitespace or start/end of string. The regex:

```
(?:^|\s)#([a-zA-Z0-9_-]+)(?:\s|$)
```

Examples:
- `fix the flaky test #later #ci` → tags: `later`, `ci`
- `#urgent deploy hotfix` → tags: `urgent`
- `issue #42 is broken` → tags: `42` (intentional: simple rule, no exceptions)
- `my-app:::base64...` → blob entries: tags extracted from label only

Tags are case-sensitive. `#Later` and `#later` are distinct. This
matches how users naturally type them and avoids surprises.

### Tag Extraction

New package: `internal/cli/pad/core/tag/tag.go`

```go
// Extract returns all unique tags from an entry string.
// For blob entries, only the label portion is scanned.
func Extract(entry string) []string

// Has returns true if the entry contains the given tag.
func Has(entry string, tagName string) bool

// Match returns true if the entry satisfies the filter.
// Supports negation: "~later" matches entries WITHOUT #later.
func Match(entry string, filter string) bool
```

### Root List Command: `--tag` Flag

```bash
ctx pad --tag later        # entries with #later
ctx pad --tag '!later'     # entries WITHOUT #later
ctx pad --tag urgent       # entries with #urgent
```

Behavior:
- Filters the displayed list to matching entries
- Entry numbers remain **original** (not renumbered) so that
  `ctx pad rm 5` still targets the correct entry after filtering
- Multiple `--tag` flags: AND logic (all must match)
- Empty result: print the standard "scratchpad is empty" message

Implementation: add a `[]string` flag `--tag` (short: `-t`) to the
root pad command. After `store.ReadEntries()`, filter with `tag.Match`
before the display loop.

### Show Command: No `--tag` Flag

`ctx pad show N` operates on absolute entry numbers. Adding tag
filtering here would create ambiguity about what N means. Users can
combine `ctx pad --tag X` to find the number, then `ctx pad show N`.

### Tags Subcommand: `ctx pad tags`

```bash
$ ctx pad tags
ci       3
later    7
urgent   2
```

Output: one line per tag, sorted alphabetically, with count.
Tab-separated for pipe friendliness.

When `--json` flag is passed:
```json
[{"tag": "ci", "count": 3}, {"tag": "later", "count": 7}]
```

No tags found: print "no tags found" (not an error).

### Edit `--tag` Flag

```bash
ctx pad edit 42 --tag later        # appends " #later" to entry 42
ctx pad edit 42 --append "text" --tag done   # combinable with other modes
ctx pad edit 42 "new text" --tag urgent      # combinable with replace
```

The `--tag` flag appends ` #tagname` to the entry. It can be used alone
or combined with `--append`, `--prepend`, or positional replace. It
conflicts with `--file`/`--label` (blob content operations).

For blob entries, the tag is appended to the label (same as `--append`
on blobs, which also modifies the label).

### Blob Append/Prepend

Previously, `--append` and `--prepend` were blocked on blob entries.
This restriction is removed: both now modify the **label** portion of
the blob entry while preserving the binary data. This enables tagging
blob entries naturally.

### What This Spec Excludes

- **Tag remove subcommands**: `ctx pad edit N "text without tag"` works.
  Dedicated commands are premature.
- **Tag namespaces or hierarchies**: no `#category:value` support.
- **Colored tag rendering**: cosmetic, not v1.
- **Tag-aware sorting or grouping**: list is position-ordered, period.
- **Cross-pad tag queries**: one scratchpad at a time.

## File Changes

### New Files

| File | Purpose |
|------|---------|
| `internal/cli/pad/core/tag/tag.go` | Extract, Has, Match functions |
| `internal/cli/pad/core/tag/tag_test.go` | Unit tests for tag logic |
| `internal/cli/pad/cmd/tags/cmd.go` | Tags subcommand definition |
| `internal/cli/pad/cmd/tags/run.go` | Tags subcommand logic |
| `internal/config/embed/cmd/pad_tags.go` | Embedded command metadata |
| `internal/config/embed/flag/pad_tags.go` | Embedded flag metadata |
| `internal/config/embed/text/pad_tags.go` | Embedded text strings |
| `internal/write/pad/tags.go` | Output formatting for tags |
| `internal/err/pad/tags.go` | Tag-specific errors (if needed) |

### Modified Files

| File | Change |
|------|--------|
| `internal/cli/pad/pad.go` | Add `--tag` flag, filter loop, register `tags` subcommand |
| `internal/config/flag/flag.go` | Add `Tag` constant |
| `internal/assets/read/desc/*.toml` | Add description keys for new flag/command |
| `docs/reference/scratchpad.md` | Document `--tag` and `tags` subcommand |
| `internal/assets/claude/skills/generated/pad/SKILL.md` | Add tag examples |

### Unchanged

- Storage format (scratchpad.enc / scratchpad.md) — no changes
- `parse.go`, `store.go` — tags are a display/query concern, not storage
- `add.go`, `edit.go` — users manage tags via existing text operations

## Testing

1. **`tag.Extract`**: empty string, no tags, single tag, multiple tags,
   tags at boundaries, blob entries (label-only extraction), hyphenated
   tags, numeric tags
2. **`tag.Match`**: positive match, negative match (`!`), no match,
   blob entries
3. **`--tag` flag integration**: single filter, multiple filters (AND),
   negation, original numbering preserved, empty result
4. **`tags` subcommand**: counts correct, alphabetical sort, no tags,
   JSON output, blob label tags included
