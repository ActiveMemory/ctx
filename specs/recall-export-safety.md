# Recall Export Safety: Lock, --keep-frontmatter, Ergonomics

## Status

**Ready for implementation.**

## Context

`ctx recall export` regenerates journal markdown from JSONL session data.
The conversation body is **always** regenerated — manual edits are lost.
`--force` additionally discards enriched YAML frontmatter. The docs say
"you can edit these files" without warning about this. Users need:

1. A way to **protect** curated journal entries from re-export
2. A clearer flag name than `--force`
3. Better ergonomics: bare command prints help, `--dry-run` previews changes

## Phase 1: Lock/Unlock State Layer

**Files**: `internal/journal/state/state.go`, `state_test.go`

- Add `Locked string` field to `FileState` (json `"locked,omitempty"`)
- Add `MarkLocked(filename)`, `ClearLocked(filename)`, `IsLocked(filename)`
- Add `"locked"` case to `Mark()` and `ValidStages`
- Tests: mark/clear/round-trip/no-op on missing entry
- Backward compatible: `omitempty` means existing `.state.json` parses fine

**Test**: `CGO_ENABLED=0 go test ./internal/journal/state/...`

## Phase 2: Lock/Unlock CLI + Export Integration

**Files**: new `internal/cli/recall/lock.go`, `lock_test.go`, `run.go`, `recall.go`

### 2A: Lock/Unlock Commands

- `ctx recall lock <pattern>` and `ctx recall unlock <pattern>`, both with `--all`
- Pattern matching: reuse slug/date/id matching from export (extract shared helper)
- Multi-part: locking base also locks all `-pN` parts
- Frontmatter: on lock, insert `locked: true  # managed by ctx` before closing `---`;
  on unlock, remove it
- `.state.json` is source of truth; frontmatter is for human visibility

### 2B: Export Respects Locks

- In `runRecallExport`, after filename is computed, before any file I/O:
  ```
  if jstate.IsLocked(filename) → skip with log line, increment locked counter
  ```
- `--force` does NOT override locks (require explicit unlock)
- Add `locked` counter to summary output

### 2C: Tests

- Lock single session, verify state + frontmatter
- Unlock, verify state + frontmatter cleaned
- Lock with `--all`
- Lock multi-part entry, verify all parts
- Export skips locked files
- Export with `--force` still skips locked files

**Test**: `CGO_ENABLED=0 go test ./internal/cli/recall/...`

## Phase 3: Replace --force with --keep-frontmatter

**Files**: `internal/cli/recall/cmd.go`, `run.go`, `run_test.go`

- Add `--keep-frontmatter` flag (bool, default `true`)
- Keep `--force` as deprecated alias via `cmd.Flags().MarkDeprecated`
- Effective logic: `discardFrontmatter := !keepFrontmatter || force`
- Rename internal `force` param to `discardFrontmatter` for clarity
- Update help text: explain body is always regenerated, only frontmatter preserved
- Tests: verify `--keep-frontmatter=false` behaves like old `--force`

**Test**: `CGO_ENABLED=0 go test ./internal/cli/recall/...`

## Phase 4: Ergonomics + Documentation

### 4A: Bare export prints help

- `runRecallExport`: when `len(args) == 0 && !all` → return `cmd.Help()`

### 4B: --dry-run flag

- Add `--dry-run` flag to export command
- Same logic but skip file writes and state saves
- Output: "Dry run: would export N new, update M existing, skip K locked"

### 4C: Documentation updates

**Replace `--force` → `--keep-frontmatter=false`:**
- `docs/cli-reference.md` (lines 773, 790)
- `docs/session-journal.md` (lines 83, 135-136)
- `docs/recipes/session-archaeology.md` (lines 180-185)
- `internal/assets/claude/skills/ctx-recall/SKILL.md` (line 75)

**Add lock/unlock docs:**
- `docs/cli-reference.md` — new sections after export
- `docs/session-journal.md` — "Protecting Entries" section
- `internal/assets/claude/skills/ctx-recall/SKILL.md` — lock/unlock subcommands

**Add --dry-run docs:**
- `docs/cli-reference.md` — export flags table
- `internal/assets/claude/skills/ctx-recall/SKILL.md`

**Clarify destructive nature:**
- `docs/session-journal.md` — warn body is always regenerated
- `docs/common-workflows.md` (line 97) — add note about body regeneration
- `docs/recipes/publishing.md` — update pipeline description
- `docs/recipes/session-archaeology.md` — update export behavior

**Test**: `CGO_ENABLED=0 go test ./...` + `make audit`

## Key Design Decisions

1. **`.state.json` is source of truth** for locks; frontmatter `locked: true` is for human visibility
2. **`--force` kept as deprecated alias** — Cobra prints warning, doesn't break scripts
3. **Locks are absolute** — `--force`/`--keep-frontmatter=false` does NOT override; explicit unlock required
4. **Bare export → help** instead of error, follows CLI conventions

## Critical Files

| File | Change |
|------|--------|
| `internal/journal/state/state.go` | Add Locked field + methods |
| `internal/cli/recall/cmd.go` | Replace --force, add --dry-run, --keep-frontmatter |
| `internal/cli/recall/run.go` | Lock check, dry-run mode, bare-command help |
| `internal/cli/recall/lock.go` | New: lock/unlock commands |
| `internal/cli/recall/recall.go` | Register lock/unlock subcommands |
| `docs/cli-reference.md` | Lock/unlock sections, flag updates |
| `docs/session-journal.md` | Protecting Entries section, destructive warnings |
| `internal/assets/claude/skills/ctx-recall/SKILL.md` | Flag + subcommand updates |

## Verification

1. `CGO_ENABLED=0 go test ./...` — all tests pass
2. `make audit` — lint, vet, drift, docs all clean
3. Manual: `ctx recall export --all --dry-run` shows summary
4. Manual: `ctx recall lock <entry>` → `ctx recall export --all` skips it
5. Manual: `ctx recall export --all --keep-frontmatter=false` discards frontmatter
6. Manual: `ctx recall export` (bare) prints help
