# Export Update Mode

## Problem

T2.1 requests an `--update` mode for `ctx recall export` that preserves
enrichment metadata while regenerating conversation content. Investigation
reveals the current codebase already implements this as the default behavior
— but with a bug and unclear semantics.

## Current State

### Default behavior (no flags)

When re-exporting a session whose journal file already exists:

1. Regenerate conversation content from JSONL source
2. Read existing file and extract YAML frontmatter
3. Merge: existing frontmatter + new content (without its frontmatter)
4. Write merged result

This **already preserves enriched frontmatter** (title, type, outcome,
topics, technologies, summary) during re-export. The enrichment skill
stores all metadata in the YAML frontmatter block, not in markdown
sections, so the current `extractFrontmatter()` + `stripFrontmatter()`
merge captures everything.

### --force flag (broken)

The flag description says "Overwrite existing files completely (discard
frontmatter)" but the code unconditionally preserves frontmatter when the
file exists — the `force` flag only changes the counter (exported vs
updated) and output message. The frontmatter merge at `run.go:179-187`
has no `!force` guard.

### --skip-existing flag

Skips files entirely. No re-export, no update.

## What T2.1 Actually Needs

The original task described three behaviors:

| Intent             | Flag              | Status          |
|--------------------|-------------------|-----------------|
| Preserve + update  | (default)         | Already works   |
| Skip existing      | `--skip-existing` | Already works   |
| Fresh overwrite    | `--force`         | Broken (bug)    |

The `--update` flag as originally conceived is unnecessary — the default
behavior already does what it describes. The real work is:

1. Fix `--force` to actually discard frontmatter
2. Add tests for the merge behavior (currently untested)
3. Clarify documentation and help text

## Design

### Fix --force

Guard the frontmatter preservation behind `!force`:

```go
// Preserve enriched YAML frontmatter from existing file
if fileExists && !force {
    existing, readErr := os.ReadFile(filepath.Clean(path))
    if readErr == nil {
        if fm := extractFrontmatter(string(existing)); fm != "" {
            content = fm + "\n" + stripFrontmatter(content)
        }
    }
}
```

When `--force` is set, the file is overwritten completely — enrichment
metadata is discarded. This matches the flag's documented behavior.

### Reset enrichment state on force

When `--force` overwrites an enriched file, the state entry should clear
the `Enriched` date so the file shows up in `countUnenriched()` again:

```go
if force {
    jstate.ClearEnriched(filename)
}
```

This requires adding a `ClearEnriched(filename)` method to the journal
state package.

### Tests

Add tests for:

1. Default re-export preserves frontmatter
2. `--force` re-export discards frontmatter
3. `--force` re-export resets enrichment state
4. `--skip-existing` leaves file untouched
5. Multipart files: frontmatter preserved per-part
6. Malformed frontmatter: graceful degradation

### Documentation

Update help text in `cmd.go` to clarify:

- Default: "Updates existing files — enriched metadata (title, type,
  topics, summary) is preserved, conversation content is regenerated
  from source."
- `--force`: "Overwrites completely — discards enriched metadata. Files
  will need re-enrichment."

## Tasks

### Phase 2: Export Preservation

- [ ] T2.1.1: Fix `--force` to actually discard frontmatter — add `!force`
      guard around frontmatter preservation in `run.go`. Currently the
      merge runs unconditionally when file exists.
      File: `internal/cli/recall/run.go:179-187`
      #priority:medium

- [ ] T2.1.2: Add `ClearEnriched()` method to journal state — when
      `--force` overwrites an enriched file, clear the enrichment date
      so it appears in `countUnenriched()` again.
      File: `internal/journal/state/state.go`
      #priority:medium

- [ ] T2.1.3: Add tests for export merge behavior — default preserves
      frontmatter, `--force` discards it, `--skip-existing` leaves file
      untouched, multipart preservation, malformed frontmatter graceful
      degradation.
      File: `internal/cli/recall/run_test.go` (or new export_test.go)
      #priority:medium

- [ ] T2.1.4: Update help text and docs — clarify default update behavior
      in `cmd.go` Long description, update `docs/session-journal.md` if
      it documents export flags.
      #priority:low

## What This Does NOT Do

- Does not add a new `--update` flag — the default already does this.
- Does not preserve manual edits to conversation content — the
  conversation section is always regenerated from JSONL source.
- Does not merge enrichment from multiple files — each file's
  frontmatter is self-contained.
