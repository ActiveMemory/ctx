# Spec: User-Level Directory Relocation (~/.local/ctx → ~/.ctx)

## Problem

The user-level directory for ctx state lives at `~/.local/ctx/` (currently
only `keys/` subdirectory). This follows XDG conventions but doesn't match
the mental model of ctx's primary users — Claude Code users who see
`~/.claude/` as the convention for tool-level config.

Having `~/.local/ctx/` as the user-level path is the odd one out.

## Decision

Relocate the user-level directory from `~/.local/ctx/` to `~/.ctx/`,
matching Claude's `~/.claude` convention.

## Current State

- `~/.local/ctx/keys/<slug>--<sha8>.key` — encryption keys
- Referenced in `internal/config/keypath.go` (`KeyDir()`)
- Docs reference `~/.local/ctx/keys/` in 12+ files (recently updated)

## Changes Required

### Code

1. **`internal/config/keypath.go`**: Change `KeyDir()` to return
   `~/.ctx/keys/` instead of `~/.local/ctx/keys/`
2. **`internal/config/migrate.go`**: Add migration tier that checks
   `~/.local/ctx/keys/` and moves to `~/.ctx/keys/` on first access
   (same copy-then-delete pattern as project-local → user-level migration)
3. **Tests**: Update `keypath_test.go`, `migrate_test.go` for new paths

### Docs (12+ files)

All files referencing `~/.local/ctx/keys/`:

- `docs/recipes/scratchpad-sync.md` (heaviest — scp examples, tips)
- `docs/reference/scratchpad.md`
- `docs/operations/upgrading.md`
- `docs/operations/migration.md`
- `docs/home/first-session.md`
- `internal/cli/notify/setup.go` (help text)
- `internal/cli/pad/doc.go` (package doc)
- `internal/cli/pad/pad.go` (help text)
- `internal/cli/initialize/run.go` (godoc)
- `.context/ARCHITECTURE.md`, `DETAILED_DESIGN.md` (if they mention the path)

### Skills

- `internal/assets/claude/skills/ctx-pad/SKILL.md` (if it mentions key path)

## Migration Strategy

Three-tier resolution (highest priority wins):

1. `.ctxrc` `key_path` override (explicit)
2. `~/.ctx/keys/<slug>.key` (new default)
3. `~/.local/ctx/keys/<slug>.key` (legacy, auto-migrated)
4. `.context/.ctx.key` (oldest legacy, already auto-migrated)

Migration happens in `MigrateKeyFile()` — add a tier between current
user-level and project-local.

## Non-Goals

- Moving `~/.claude/` — that's Anthropic's convention, not ours
- Creating other subdirectories under `~/.ctx/` yet (that comes in the
  state consolidation spec for global state)

## Future: ~/.ctx as Global State Home

Once `~/.ctx/` exists, it becomes the natural home for any future
global (non-project) ctx state. The state consolidation spec may
use `~/.ctx/` for global state that doesn't belong in `.context/state/`.
