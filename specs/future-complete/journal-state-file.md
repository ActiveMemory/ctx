# Plan: External Journal State File

## Context

Journal processing state is tracked via HTML comment markers embedded in journal files (`<!-- normalized: ... -->`, `<!-- fences-verified: ... -->`). This is fundamentally broken because journal files contain AI conversations that include these exact strings as content, causing false positives. Similarly, `countUnenriched()` scans file content for `---` prefix, but section breaks also use `---`. The user is okay nuking all journals and re-exporting from JSONL.

## Approach

Replace all in-band markers with `.context/journal/.state.json` — a single external file that tracks processing status per journal entry. Already gitignored (inherits from `.context/journal/` in `.gitignore`).

### State file format

```json
{
  "version": 1,
  "entries": {
    "2026-01-21-async-roaming-allen-af7cba21.md": {
      "exported": "2026-01-21",
      "enriched": "2026-01-24",
      "normalized": "2026-01-22",
      "fences_verified": "2026-01-23"
    }
  }
}
```

Date strings (not booleans) — provides audit trail at zero extra cost.

## Implementation

### Phase 1: State package

1. **Create `internal/journal/state/state.go`** — types (`JournalState`, `FileState`), `Load`/`Save` (atomic write via temp+rename), query helpers (`IsEnriched`, `CountUnenriched`), mutation helpers (`MarkExported`, `MarkEnriched`, `Rename`).
2. **Create `internal/journal/state/state_test.go`** — load missing file, round-trip, count, rename.
3. **Add `FileJournalState = ".state.json"`** to `internal/config/file.go`.

### Phase 2: Wire into export

4. **`internal/cli/recall/run.go`** — In `runRecallExport()`, load state at start, call `MarkExported()` after each file write, `Rename()` when slug changes, `Save()` at end.

### Phase 3: Wire into site pipeline

5. **`internal/cli/journal/reduce.go`** — Change `stripFences(content string)` → `stripFences(content string, fencesVerified bool)`. Replace `RegExFencesVerified.MatchString(content)` with `if fencesVerified { return content }`.
6. **`internal/cli/journal/normalize.go`** — Change `normalizeContent(content string)` → `normalizeContent(content string, fencesVerified bool)`, forward to `stripFences`.
7. **`internal/cli/journal/run.go`** — Load state at top of `runJournalSite()`. Pass `state.IsFencesVerified(entry.Filename)` per file.

### Phase 4: Wire into check-journal

8. **`internal/cli/system/check_journal.go`** — Replace `countUnenriched()` body: load state, count `.md` files in directory without `enriched` set.

### Phase 5: Update tests

9. **`internal/cli/journal/journal_test.go`** — Update `stripFences` tests to use boolean param instead of in-content markers.

### Phase 6: Cleanup

10. **`internal/config/regex.go`** — Remove `RegExNormalizedMarker` and `RegExFencesVerified`.

### Phase 7: Update skills

11. **`ctx-journal-normalize/SKILL.md`** — Replace marker-based idempotency with `.state.json` reads/writes.
12. **`ctx-journal-enrich-all/SKILL.md`** — Replace `head -1` check with `.state.json` lookup. Update after enriching.
13. **`ctx-journal-enrich/SKILL.md`** — Replace `grep -rL "^---$"` (line 44) with `.state.json` lookup. Update after enriching.

### Phase 8: CLI helper (optional but recommended)

14. **Add `ctx system mark-journal <filename> <stage>`** subcommand — makes state updates from skills trivial and atomic. Skills call `ctx system mark-journal session.md enriched` instead of editing JSON.

## Files changed

| File | Change |
|------|--------|
| `internal/journal/state/state.go` | **NEW** — state package |
| `internal/journal/state/state_test.go` | **NEW** — tests |
| `internal/config/file.go` | Add constant |
| `internal/config/regex.go` | Remove 2 regex constants |
| `internal/cli/journal/reduce.go` | `stripFences` bool param |
| `internal/cli/journal/normalize.go` | Forward bool param |
| `internal/cli/journal/run.go` | Load state, pass per-file |
| `internal/cli/journal/journal_test.go` | Update tests |
| `internal/cli/recall/run.go` | Mark exported on export |
| `internal/cli/system/check_journal.go` | State-based `countUnenriched` |
| `internal/assets/.../ctx-journal-normalize/SKILL.md` | State file instructions |
| `internal/assets/.../ctx-journal-enrich-all/SKILL.md` | State file instructions |
| `internal/assets/.../ctx-journal-enrich/SKILL.md` | State file instructions |
| `internal/cli/journal/journal.go` | Add `mark` subcommand (optional) |

## Verification

1. `go test ./...` — all packages pass
2. `ctx recall export --all` into empty journal dir → `.state.json` created with exported dates
3. `ctx system check-journal` → correct unenriched count from state file
4. `ctx journal site` → `stripFences` respects `fences_verified` from state
5. Run `/ctx-journal-enrich` on one entry → state file updated
6. Plugin version bump (0.6.2 → 0.6.3) for skill changes
