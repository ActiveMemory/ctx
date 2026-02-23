# Plan: `ctx recall sync` — Frontmatter-to-State Lock Sync

## Context

`ctx recall lock` writes lock state to `.state.json` and updates frontmatter
for visibility. But from the user's perspective during journal enrichment, the
natural flow is: edit markdown frontmatter (add `locked: true`), then sweep
all files to propagate that state. The current `lock` command requires naming
files explicitly, which is friction after a batch enrichment pass.

## Approach

Add `ctx recall sync` as a sister command to `lock`/`unlock`. It scans all
journal markdowns, reads their YAML frontmatter for `locked: true`, and syncs
that into `.state.json` — the inverse direction of `lock`.

### Behavior

| Frontmatter     | `.state.json` | Action                     |
|-----------------|---------------|----------------------------|
| `locked: true`  | not locked    | Mark locked in state       |
| no `locked:`    | locked        | Clear lock in state        |
| `locked: true`  | locked        | No change                  |
| no `locked:`    | not locked    | No change                  |

### CLI surface

```
ctx recall sync
```

No flags. No arguments. Scans everything, reports changes.

### Output

```
  ✓ 2026-01-21-session-abc12345.md (locked)
  ✓ 2026-01-22-session-def67890.md (unlocked)

Locked 1 entry(s).
Unlocked 1 entry(s).
```

Or: `No changes — state already matches frontmatter.`

## Implementation

### Phase 1: Core command (done)

1. **`internal/cli/recall/sync.go`** — `recallSyncCmd()`, `runSync()`,
   `frontmatterHasLocked()` helper.

### Phase 2: Wire up

2. **`internal/cli/recall/recall.go`** — Add `cmd.AddCommand(recallSyncCmd())`
   and update help text.

### Phase 3: Tests

3. **`internal/cli/recall/sync_test.go`** — Test cases:
   - Sync locks files with `locked: true` in frontmatter
   - Sync unlocks files missing `locked:` in frontmatter but locked in state
   - No changes when already in sync
   - No journal entries found (empty dir)
   - Files without frontmatter are treated as unlocked
   - `locked: false` in frontmatter treated as unlocked

### Phase 4: Documentation

4. **`docs/reference/cli-reference.md`** — Add `ctx recall sync` section
   after `ctx recall unlock`.
5. **`docs/reference/session-journal.md`** — Mention sync in the enrichment
   workflow section.
6. **`docs/recipes/session-archaeology.md`** — Add sync to the lock/unlock
   discussion.

## Files changed

| File | Change |
|------|--------|
| `internal/cli/recall/sync.go` | **NEW** — sync command |
| `internal/cli/recall/sync_test.go` | **NEW** — tests |
| `internal/cli/recall/recall.go` | Wire up + help text |
| `docs/reference/cli-reference.md` | Add `ctx recall sync` section |
| `docs/reference/session-journal.md` | Mention sync in enrichment workflow |
| `docs/recipes/session-archaeology.md` | Add sync reference |

## Verification

1. `go test ./internal/cli/recall/...` — all tests pass
2. `go test ./...` — full suite green
3. `ctx recall sync` — no-op on project with state already matching
4. Manually add `locked: true` to a frontmatter, run sync, verify state updates
