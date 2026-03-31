# Recall Import Safety: Safe Defaults, Locks, and Ergonomics

## Status

**Ready for implementation.** (Revised 2026-02-21: added safe-by-default
import and confirmation prompt — see Phase 1.)

## Context

`ctx recall import` regenerates journal markdown from JSONL session data.
The conversation body is **always** regenerated — manual edits are lost.
`--force` additionally discards enriched YAML frontmatter. The docs say
"you can edit these files" without warning about this.

### Current behavior (the problem)

| Command                          | Body               | Frontmatter | Confirmation |
|----------------------------------|---------------------|-------------|--------------|
| `import --all`                   | Regenerated (destructive) | Preserved   | None         |
| `import --all --skip-existing`   | Untouched           | Untouched   | None         |
| `import --all --force`           | Regenerated         | Discarded   | None         |

**Two fundamental issues:**

1. **Default is destructive** — `import --all` silently regenerates every
   existing journal body. Users won't RTFM, and silently overwriting 160+
   journal bodies is not a sane default.
2. **No confirmation** — destructive operations proceed without showing what
   will happen or asking for consent.

### What users need

1. **Safe-by-default import** — `--all` should only import *new* sessions
2. **Explicit opt-in for regeneration** — a `--regenerate` flag for re-importing
3. **Confirmation before destructive ops** — show summary, ask `proceed? [y/N]`
4. **Lock protection** for curated entries
5. **Clearer flag names** — `--keep-frontmatter` instead of `--force`
6. **Better ergonomics** — bare command prints help, `--dry-run` previews

## Phase 1: Safe-by-Default Import

**Files**: `internal/cli/recall/cmd.go`, `run.go`

This is the core behavioral change. `import --all` becomes safe by default.

### 1A: New default — import new sessions only

- Change `runRecallImport` so `--all` (without `--regenerate`) skips files
  that already exist on disk. This makes `--skip-existing` the implicit
  default when using `--all`.
- Deprecate `--skip-existing` via `cmd.Flags().MarkDeprecated` — it's now
  the default behavior and no longer needed as a flag.
- A single-session import (`import <id>`) always writes (specific intent).

### 1B: `--regenerate` flag for re-importing existing sessions

- Add `--regenerate` flag (bool, default `false`).
- When set, existing files are regenerated (body rewritten, frontmatter
  preserved unless `--keep-frontmatter=false`).
- `--regenerate` without `--all` is an error — regeneration is a bulk concern.

### 1C: Confirmation prompt before destructive writes

- Before any file I/O, compute the plan: count new, regenerate, locked/skipped.
- If `regenerate > 0` (or `force` / `--keep-frontmatter=false`), print summary
  and prompt:
  ```
  Will import 5 new, regenerate 12 existing, skip 3 locked.
  Proceed? [y/N]
  ```
- `--yes` / `-y` flag to skip confirmation (for scripts and automation).
- New-only imports (no regeneration) proceed without confirmation — they're safe.
- `--dry-run` prints the summary and exits (never prompts).

### 1D: Updated behavior matrix

| Command                                        | Body (new)  | Body (existing) | Frontmatter | Confirmation |
|------------------------------------------------|-------------|-----------------|-------------|--------------|
| `import --all`                                 | Imported    | Untouched       | n/a         | No           |
| `import --all --regenerate`                    | Imported    | Regenerated     | Preserved   | **Yes**      |
| `import --all --regenerate --keep-fm=false`    | Imported    | Regenerated     | Discarded   | **Yes**      |
| `import --all --regenerate --yes`              | Imported    | Regenerated     | Preserved   | No (bypassed)|
| `import --all --dry-run`                       | (counted)   | (counted)       | (counted)   | No (preview) |
| `import <id>`                                  | Imported    | Regenerated     | Preserved   | No           |

**Test**: `CGO_ENABLED=0 go test ./internal/cli/recall/...`

## Phase 2: Lock/Unlock State Layer

**Files**: `internal/journal/state/state.go`, `state_test.go`

- Add `Locked string` field to `FileState` (json `"locked,omitempty"`)
- Add `MarkLocked(filename)`, `ClearLocked(filename)`, `IsLocked(filename)`
- Add `"locked"` case to `Mark()` and `ValidStages`
- Tests: mark/clear/round-trip/no-op on missing entry
- Backward compatible: `omitempty` means existing `.state.json` parses fine

**Test**: `CGO_ENABLED=0 go test ./internal/journal/state/...`

## Phase 3: Lock/Unlock CLI + Import Integration

**Files**: new `internal/cli/recall/lock.go`, `lock_test.go`, `run.go`, `recall.go`

### 3A: Lock/Unlock Commands

- `ctx recall lock <pattern>` and `ctx recall unlock <pattern>`, both with `--all`
- Pattern matching: reuse slug/date/id matching from import (extract shared helper)
- Multi-part: locking base also locks all `-pN` parts
- Frontmatter: on lock, insert `locked: true  # managed by ctx` before closing `---`;
  on unlock, remove it
- `.state.json` is source of truth; frontmatter is for human visibility

### 3B: Import Respects Locks

- In `runRecallImport`, after filename is computed, before any file I/O:
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
- Import skips locked files (with and without `--regenerate`)
- Import with `--force` still skips locked files

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

### 5A: Bare import prints help

- `runRecallImport`: when `len(args) == 0 && !all` → return `cmd.Help()`

### 5B: --dry-run flag

- Add `--dry-run` flag to import command
- Same plan computation but skip file writes, state saves, and confirmation
- Output: "Would import N new, regenerate M existing, skip K locked"

### 5C: Documentation updates

**Replace `--force` → `--keep-frontmatter=false`:**
- `docs/cli-reference.md`
- `docs/session-journal.md`
- `docs/recipes/session-archaeology.md`
- `internal/assets/claude/skills/ctx-recall/SKILL.md`

**Add new flags and behavior:**
- `docs/cli-reference.md` — `--regenerate`, `--yes`, `--dry-run`, `--keep-frontmatter`
- `docs/session-journal.md` — "Safe Import" section explaining new defaults
- `internal/assets/claude/skills/ctx-recall/SKILL.md` — updated flag reference

**Add lock/unlock docs:**
- `docs/cli-reference.md` — new sections after import
- `docs/session-journal.md` — "Protecting Entries" section
- `internal/assets/claude/skills/ctx-recall/SKILL.md` — lock/unlock subcommands

**Clarify destructive nature:**
- `docs/session-journal.md` — warn body is regenerated on `--regenerate`
- `docs/common-workflows.md` — add note about import safety
- `docs/recipes/publishing.md` — update pipeline description
- `docs/recipes/session-archaeology.md` — update import behavior

**Deprecation notes:**
- `--skip-existing` deprecated (now the default)
- `--force` deprecated (use `--keep-frontmatter=false`)

**Test**: `CGO_ENABLED=0 go test ./...` + `make audit`

## Key Design Decisions

1. **Safe-by-default** — `import --all` only imports new sessions; regenerating
   existing entries requires explicit `--regenerate`. Users won't RTFM.
2. **Confirmation for destructive ops** — any command that regenerates existing
   files shows a summary and asks `proceed? [y/N]`. `--yes` bypasses.
3. **Single-session import is always direct** — `import <id>` writes without
   confirmation because targeting a specific session is explicit intent.
4. **`.state.json` is source of truth** for locks; frontmatter `locked: true`
   is for human visibility.
5. **`--force` kept as deprecated alias** — Cobra prints warning, doesn't break
   scripts.
6. **Locks are absolute** — `--regenerate`/`--force`/`--keep-frontmatter=false`
   do NOT override locks; explicit `unlock` required.
7. **Bare import → help** instead of error, follows CLI conventions.
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
| `docs/session-journal.md` | Safe Import + Protecting Entries sections |
| `internal/assets/claude/skills/ctx-recall/SKILL.md` | Flag + subcommand updates |

## Verification

1. `CGO_ENABLED=0 go test ./...` — all tests pass
2. `make audit` — lint, vet, drift, docs all clean
3. Manual: `ctx recall import --all` imports only new sessions
4. Manual: `ctx recall import --all --regenerate` prompts for confirmation
5. Manual: `ctx recall import --all --regenerate --yes` bypasses prompt
6. Manual: `ctx recall import --all --dry-run` shows summary without writing
7. Manual: `ctx recall lock <entry>` → `ctx recall import --all --regenerate` skips it
8. Manual: `ctx recall import --all --regenerate --keep-frontmatter=false` discards frontmatter
9. Manual: `ctx recall import` (bare) prints help
