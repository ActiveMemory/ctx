# Spec: Global Encryption Key (~/.ctx/.ctx.key)

Supersedes: `specs/user-level-dir-relocation.md` (key portions only;
the `~/.local/ctx → ~/.ctx` directory move remains valid context).

## Problem

The current key resolution has three layers and per-project slug
filenames (`~/.local/ctx/keys/<slug>--<sha8>.key`). This is
over-engineered for the common case: one user, one machine, one key.

Pain points:

1. **Unnecessary complexity** — three resolution tiers, slug generation
   with SHA hashing, and migration logic between tiers
2. **Wrong default location** — `~/.local/ctx/` doesn't match Claude
   Code's `~/.claude/` convention; `~/.ctx/` is the natural home
3. **Per-project keys without clear benefit** — the slug convention
   creates one key per project, but users who want isolation can use
   `.ctxrc key_path` or a project-local `.context/.ctx.key`
4. **Tilde not expanded** — `.ctxrc key_path` with `~/...` fails
   because the shell never expands it (separate bugfix)

## Decision

Single global key at `~/.ctx/.ctx.key`. Per-project override via
`.ctxrc key_path` or project-local `.context/.ctx.key`.

## New Resolution Order

Two tiers (highest priority wins):

1. **Project-level override**: `.ctxrc key_path` (explicit) OR
   `$PROJECT/.context/.ctx.key` (file exists)
2. **Global default**: `~/.ctx/.ctx.key`

Drop: per-project slug directory (`~/.local/ctx/keys/`), the
`ProjectKeySlug()` function, `ProjectKeyPath()`, `KeyDir()`.

## Changes Required

### Phase 1: Code — key resolution and migration

| File | Change |
|------|--------|
| `internal/config/keypath.go` | Replace `KeyDir()`, `ProjectKeySlug()`, `ProjectKeyPath()` with `GlobalKeyPath()` returning `~/.ctx/.ctx.key`. Simplify `ResolveKeyPath()` to two-tier: project-local/override → global. Add tilde expansion for override path. |
| `internal/config/keypath_test.go` | Rewrite tests for new two-tier resolution. |
| `internal/config/migrate.go` | Replace `MigrateKeyFile()`: (a) promote any `~/.local/ctx/keys/*.key` to `~/.ctx/.ctx.key` if global doesn't exist yet (warn if multiple distinct keys found), (b) promote project-local `.context/.ctx.key` to global, (c) clean up legacy names (`.context.key`, `.scratchpad.key`). |
| `internal/config/migrate_test.go` | Rewrite for new migration tiers. |
| `internal/rc/rc.go` | `KeyPath()` calls simplified `ResolveKeyPath()`. |
| `internal/rc/types.go` | No structural change; `KeyPathOverride` stays. |
| `internal/config/file.go` | Keep `FileContextKey` (`.ctx.key`) — still used for project-local override detection. |
| `internal/config/dir.go` | Keep `.context/.ctx.key` in gitignore entries. |

### Phase 2: Code — callers (no logic change, just verify)

| File | What to verify |
|------|----------------|
| `internal/cli/initialize/run.go` | `initScratchpad()` uses `rc.KeyPath()` — no change needed, but `os.MkdirAll` target dir changes from `~/.local/ctx/keys/` to `~/.ctx/`. |
| `internal/cli/pad/store.go` | `keyPath()` calls `rc.KeyPath()` — no change. |
| `internal/notify/notify.go` | `LoadWebhook()` / `SaveWebhook()` call `rc.KeyPath()` — no change. |
| `internal/cli/system/check_version.go` | `checkKeyAge()` calls `rc.KeyPath()` — no change. |
| `internal/cli/pad/pad_test.go` | `setupEncrypted()` uses `config.ProjectKeyPath()` — must switch to `config.GlobalKeyPath()`. |
| `internal/cli/initialize/initialize_test.go` | Same — switch from `ProjectKeyPath()`. |

### Phase 3: Documentation

All files referencing `~/.local/ctx/keys/` or slug-based paths:

- `docs/reference/scratchpad.md` — rewrite "Encrypted by Default" and
  "Key Distribution" sections
- `docs/recipes/scratchpad-sync.md` — heaviest; scp examples, tips
- `docs/recipes/webhook-notifications.md`
- `docs/recipes/parallel-worktrees.md`
- `docs/recipes/scratchpad-with-claude.md`
- `docs/operations/upgrading.md`
- `docs/operations/migration.md`
- `docs/home/first-session.md`
- `internal/cli/pad/doc.go` (package doc)
- `internal/cli/pad/pad.go` (help text)
- `internal/cli/notify/setup.go` (help text)
- `.context/ARCHITECTURE.md`, `DETAILED_DESIGN.md`

### Phase 4: Cleanup

- Delete `specs/user-level-dir-relocation.md` (superseded)
- Remove stale test key files from `~/.local/ctx/keys/` (manual or
  migration cleans on first access)
- Record decision in DECISIONS.md

## Migration Strategy

`MigrateKeyFile()` runs on every `rc.KeyPath()` call (existing pattern).

**New migration flow:**

```
1. If ~/.ctx/.ctx.key exists → done (global key in place)
2. If ~/.local/ctx/keys/*.key exists:
   a. If all keys are identical → copy one to ~/.ctx/.ctx.key
   b. If keys differ → copy this project's slug key, warn on stderr
3. If .context/.ctx.key exists → copy to ~/.ctx/.ctx.key
4. Legacy names (.context.key, .scratchpad.key) → rename to .context/.ctx.key, then step 3
5. Clean up: remove project-local .context/.ctx.key after promotion
```

After migration settles (1-2 releases), the legacy `~/.local/ctx/keys/`
tier can be removed entirely.

## Tilde Expansion (Bugfix)

`ResolveKeyPath()` must expand leading `~/` in the override path using
`os.UserHomeDir()`. This is a standalone fix that applies regardless
of the key relocation.

## Non-Goals

- Moving `~/.claude/` — Anthropic's convention
- Global state beyond the key file (future work)
- Removing `.ctxrc key_path` override — still needed for unusual setups
- Encrypting with different keys per project by default

## Rollback

If a user needs per-project keys, they set `key_path` in each
project's `.ctxrc`. The global default is the 99% case.
