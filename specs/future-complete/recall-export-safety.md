# Recall Export Safety: Safe Defaults, Locks, and Ergonomics

## Status

**Ready for implementation.** (Revised 2026-02-21: added safe-by-default
export and confirmation prompt — see Phase 1.)

## Context

`ctx recall export` regenerates journal markdown from JSONL session data.
The conversation body is **always** regenerated — manual edits are lost.
`--force` additionally discards enriched YAML frontmatter. The docs say
"you can edit these files" without warning about this.

### Current behavior (the problem)

| Command                          | Body               | Frontmatter | Confirmation |
|----------------------------------|---------------------|-------------|--------------|
| `export --all`                   | Regenerated (destructive) | Preserved   | None         |
| `export --all --skip-existing`   | Untouched           | Untouched   | None         |
| `export --all --force`           | Regenerated         | Discarded   | None         |

**Two fundamental issues:**

1. **Default is destructive** — `export --all` silently regenerates every
   existing journal body. Users won't RTFM, and silently overwriting 160+
   journal bodies is not a sane default.
2. **No confirmation** — destructive operations proceed without showing what
   will happen or asking for consent.

### What users need

1. **Safe-by-default export** — `--all` should only export *new* sessions
2. **Explicit opt-in for regeneration** — a `--regenerate` flag for re-exporting
3. **Confirmation before destructive ops** — show summary, ask `proceed? [y/N]`
4. **Lock protection** for curated entries
5. **Clearer flag names** — `--keep-frontmatter` instead of `--force`
6. **Better ergonomics** — bare command prints help, `--dry-run` previews

## Phase 1: Safe-by-Default Export

**Files**: `internal/cli/recall/cmd.go`, `run.go`

This is the core behavioral change. `export --all` becomes safe by default.

### 1A: New default — export new sessions only

- Change `runRecallExport` so `--all` (without `--regenerate`) skips files
  that already exist on disk. This makes `--skip-existing` the implicit
  default when using `--all`.
- Deprecate `--skip-existing` via `cmd.Flags().MarkDeprecated` — it's now
  the default behavior and no longer needed as a flag.
- A single-session export (`export <id>`) always writes (specific intent).

### 1B: `--regenerate` flag for re-exporting existing sessions

- Add `--regenerate` flag (bool, default `false`).
- When set, existing files are regenerated (body rewritten, frontmatter
  preserved unless `--keep-frontmatter=false`).
- `--regenerate` without `--all` is an error — regeneration is a bulk concern.

### 1C: Confirmation prompt before destructive writes

- Before any file I/O, compute the plan: count new, regenerate, locked/skipped.
- If `regenerate > 0` (or `force` / `--keep-frontmatter=false`), print summary
  and prompt:
  ```
  Will export 5 new, regenerate 12 existing, skip 3 locked.
  Proceed? [y/N]
  ```
- `--yes` / `-y` flag to skip confirmation (for scripts and automation).
- New-only exports (no regeneration) proceed without confirmation — they're safe.
- `--dry-run` prints the summary and exits (never prompts).

### 1D: Updated behavior matrix

| Command                                        | Body (new)  | Body (existing) | Frontmatter | Confirmation |
|------------------------------------------------|-------------|-----------------|-------------|--------------|
| `export --all`                                 | Exported    | Untouched       | n/a         | No           |
| `export --all --regenerate`                    | Exported    | Regenerated     | Preserved   | **Yes**      |
| `export --all --regenerate --keep-fm=false`    | Exported    | Regenerated     | Discarded   | **Yes**      |
| `export --all --regenerate --yes`              | Exported    | Regenerated     | Preserved   | No (bypassed)|
| `export --all --dry-run`                       | (counted)   | (counted)       | (counted)   | No (preview) |
| `export <id>`                                  | Exported    | Regenerated     | Preserved   | No           |

**Test**: `CGO_ENABLED=0 go test ./internal/cli/recall/...`

## Phase 2: Lock/Unlock State Layer

**Files**: `internal/journal/state/state.go`, `state_test.go`

- Add `Locked string` field to `FileState` (json `"locked,omitempty"`)
- Add `MarkLocked(filename)`, `ClearLocked(filename)`, `IsLocked(filename)`
- Add `"locked"` case to `Mark()` and `ValidStages`
- Tests: mark/clear/round-trip/no-op on missing entry
- Backward compatible: `omitempty` means existing `.state.json` parses fine

**Test**: `CGO_ENABLED=0 go test ./internal/journal/state/...`

## Phase 3: Lock/Unlock CLI + Export Integration

**Files**: new `internal/cli/recall/lock.go`, `lock_test.go`, `run.go`, `recall.go`

### 3A: Lock/Unlock Commands

- `ctx recall lock <pattern>` and `ctx recall unlock <pattern>`, both with `--all`
- Pattern matching: reuse slug/date/id matching from export (extract shared helper)
- Multi-part: locking base also locks all `-pN` parts
- Frontmatter: on lock, insert `locked: true  # managed by ctx` before closing `---`;
  on unlock, remove it
- `.state.json` is source of truth; frontmatter is for human visibility

### 3B: Export Respects Locks

- In `runRecallExport`, after filename is computed, before any file I/O:
  ```
  if jstate.IsLocked(filename) → skip with log line, increment locked counter
  ```
- Neither `--regenerate` nor `--force` overrides locks (require explicit unlock)
- Add `locked` counter to confirmation summary and final output

### 3C: Tests

- Lock single session, verify state + frontmatter
- Unlock, verify state + frontmatter cleaned
- Lock with `--all`
- Lock multi-part entry, verify all parts
- Export skips locked files (with and without `--regenerate`)
- Export with `--force` still skips locked files

**Test**: `CGO_ENABLED=0 go test ./internal/cli/recall/...`

## Phase 4: Replace --force with --keep-frontmatter

**Files**: `internal/cli/recall/cmd.go`, `run.go`, `run_test.go`

- Add `--keep-frontmatter` flag (bool, default `true`)
- Keep `--force` as deprecated alias via `cmd.Flags().MarkDeprecated`
- Effective logic: `discardFrontmatter := !keepFrontmatter || force`
- Rename internal `force` param to `discardFrontmatter` for clarity
- `--keep-frontmatter=false` implies `--regenerate` (can't discard frontmatter
  without regenerating)
- Update help text: explain body is always regenerated, only frontmatter preserved
- Tests: verify `--keep-frontmatter=false` behaves like old `--force`

**Test**: `CGO_ENABLED=0 go test ./internal/cli/recall/...`

## Phase 5: Ergonomics + Documentation

### 5A: Bare export prints help

- `runRecallExport`: when `len(args) == 0 && !all` → return `cmd.Help()`

### 5B: --dry-run flag

- Add `--dry-run` flag to export command
- Same plan computation but skip file writes, state saves, and confirmation
- Output: "Would export N new, regenerate M existing, skip K locked"

### 5C: Documentation updates

**Replace `--force` → `--keep-frontmatter=false`:**
- `docs/cli-reference.md`
- `docs/session-journal.md`
- `docs/recipes/session-archaeology.md`
- `internal/assets/claude/skills/ctx-recall/SKILL.md`

**Add new flags and behavior:**
- `docs/cli-reference.md` — `--regenerate`, `--yes`, `--dry-run`, `--keep-frontmatter`
- `docs/session-journal.md` — "Safe Export" section explaining new defaults
- `internal/assets/claude/skills/ctx-recall/SKILL.md` — updated flag reference

**Add lock/unlock docs:**
- `docs/cli-reference.md` — new sections after export
- `docs/session-journal.md` — "Protecting Entries" section
- `internal/assets/claude/skills/ctx-recall/SKILL.md` — lock/unlock subcommands

**Clarify destructive nature:**
- `docs/session-journal.md` — warn body is regenerated on `--regenerate`
- `docs/common-workflows.md` — add note about export safety
- `docs/recipes/publishing.md` — update pipeline description
- `docs/recipes/session-archaeology.md` — update export behavior

**Deprecation notes:**
- `--skip-existing` deprecated (now the default)
- `--force` deprecated (use `--keep-frontmatter=false`)

**Test**: `CGO_ENABLED=0 go test ./...` + `make audit`

## Key Design Decisions

1. **Safe-by-default** — `export --all` only exports new sessions; regenerating
   existing entries requires explicit `--regenerate`. Users won't RTFM.
2. **Confirmation for destructive ops** — any command that regenerates existing
   files shows a summary and asks `proceed? [y/N]`. `--yes` bypasses.
3. **Single-session export is always direct** — `export <id>` writes without
   confirmation because targeting a specific session is explicit intent.
4. **`.state.json` is source of truth** for locks; frontmatter `locked: true`
   is for human visibility.
5. **`--force` kept as deprecated alias** — Cobra prints warning, doesn't break
   scripts.
6. **Locks are absolute** — `--regenerate`/`--force`/`--keep-frontmatter=false`
   do NOT override locks; explicit `unlock` required.
7. **Bare export → help** instead of error, follows CLI conventions.
8. **`--keep-frontmatter=false` implies `--regenerate`** — you can't discard
   frontmatter without also regenerating the body.

## Critical Files

| File | Change |
|------|--------|
| `internal/cli/recall/cmd.go` | New flags: --regenerate, --yes, --dry-run, --keep-frontmatter; deprecate --force, --skip-existing |
| `internal/cli/recall/run.go` | Safe default, plan/confirm flow, lock check, dry-run mode, bare help |
| `internal/journal/state/state.go` | Add Locked field + methods |
| `internal/cli/recall/lock.go` | New: lock/unlock commands |
| `internal/cli/recall/recall.go` | Register lock/unlock subcommands |
| `docs/cli-reference.md` | Lock/unlock sections, flag updates, deprecation notes |
| `docs/session-journal.md` | Safe Export + Protecting Entries sections |
| `internal/assets/claude/skills/ctx-recall/SKILL.md` | Flag + subcommand updates |

## Verification

1. `CGO_ENABLED=0 go test ./...` — all tests pass
2. `make audit` — lint, vet, drift, docs all clean
3. Manual: `ctx recall export --all` exports only new sessions
4. Manual: `ctx recall export --all --regenerate` prompts for confirmation
5. Manual: `ctx recall export --all --regenerate --yes` bypasses prompt
6. Manual: `ctx recall export --all --dry-run` shows summary without writing
7. Manual: `ctx recall lock <entry>` → `ctx recall export --all --regenerate` skips it
8. Manual: `ctx recall export --all --regenerate --keep-frontmatter=false` discards frontmatter
9. Manual: `ctx recall export` (bare) prints help
