---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: CLI namespace cleanup — singularization + promotion out of system/
status: accepted
date: 2026-04-11
owner: jose
scope: four-commit namespace cleanup, breaking
---

# Spec — CLI namespace cleanup

## Problem

Two independent but related inconsistencies in the `ctx` CLI surface:

1. **Plural command names.** The 2026-03-18 convention decision
   ("Singular command names for all CLI entities") is violated in seven
   places:
   - Four hidden hook plumbing commands under `ctx system`:
     `block-dangerous-commands`, `check-ceremonies`, `check-reminders`,
     `check-resources`.
   - Two visible commands under `ctx system`: `events`, `resources`.
   - One visible `ctx pad` subcommand: `tags`.
   - One borderline: `ctx system stats` — kept plural, semantically
     acceptable as an abbreviation of "statistics".
2. **User-facing commands parked under `ctx system`.** The `system`
   namespace was originally the home for all internal plumbing — Claude
   Code hook implementations, session lifecycle checks, and so on. Over
   time, user-facing maintenance commands accreted under the same
   namespace (`backup`, `bootstrap`, `events`, `message`, `prune`,
   `resources`, `stats`). These are routinely documented in recipes and
   CLAUDE.md templates as user-facing, but they live under a parent
   whose other 27 children are hidden hook plumbing. That's a design
   smell: the namespace mixes two different trust/visibility tiers.

The rule: if a command is not leveraged by hooks (not invoked by any
hook code path as a subprocess) and is user-facing (documented, runnable
by humans), it does not belong under `ctx system`. The `system`
namespace is reserved for plumbing that agents and hooks call, never
humans.

## Scope — 12 changes, 4 commits

### Commit 1 — Hidden plumbing plural renames

Pure internal rename of four hidden hook subcommands:

| Old                         | New                        |
|-----------------------------|----------------------------|
| `block-dangerous-commands`  | `block-dangerous-command`  |
| `check-ceremonies`          | `check-ceremony`           |
| `check-reminders`           | `check-reminder`           |
| `check-resources`           | `check-resource`           |

Hidden from users. Referenced by hook config in
`.claude/hooks/` (if any) and by YAML asset keys. No deprecation shim —
hook configs are plugin-owned and regenerate on plugin update.

### Commit 2 — `ctx pad tag` rename

`ctx pad tags` → `ctx pad tag`. Visible user-facing. One subcommand of
an otherwise-singular namespace. No deprecation shim per the user's
pure-rename call; release notes flag the breaking change.

### Commit 3 — Promote seven commands out of `system/`

| Old                      | New              | Visible? | Notes                                    |
|--------------------------|------------------|----------|------------------------------------------|
| `ctx system backup`      | `ctx backup`     | yes      | no rename, just parent change            |
| `ctx system bootstrap`   | `ctx bootstrap`  | yes      | no rename, just parent change            |
| `ctx system events`      | `ctx event`      | yes      | promote + singularize                    |
| `ctx system message *`   | `ctx message *`  | yes      | promote, keep all four sub-sub commands  |
| `ctx system prune`       | `ctx prune`      | yes      | promote                                   |
| `ctx system resources`   | `ctx resource`   | yes      | promote + singularize                    |
| `ctx system stats`       | `ctx stats`      | yes      | promote; `stats` stays plural as abbreviation |

After commit 3:

- `ctx system` parent becomes `Hidden: true`. Only the 27 hidden hook
  plumbing subcommands remain underneath. Users who need to inspect
  plumbing can still run `ctx system --help` directly.
- The existing `internal/cli/system/cmd/{backup,bootstrap,events,message,prune,resources,stats}/`
  packages are moved to `internal/cli/{backup,bootstrap,event,message,prune,resource,stats}/`.
  Plural directory names become singular where the command was singularized.
- `internal/config/embed/cmd/system.go` loses 7 Use constants and gains
  in-place updates for the 4 hidden plurals. New Use constants live in
  new files per-command: `internal/config/embed/cmd/backup.go`,
  `bootstrap.go`, `event.go`, `message.go`, `prune.go`, `resource.go`,
  `stats.go`.
- The top-level CLI root gains 7 new `AddCommand` calls for the
  promoted packages.

### Commit 4 — Docs, recipes, skills, CLAUDE.md sweep

All text surfaces carrying the old command names:

- `internal/assets/commands/commands.yaml` — Short/Long descriptions.
- `internal/assets/commands/examples.yaml` — Example entries. Keys
  rename from `system.backup` to `backup`, etc.
- `internal/assets/commands/text/hooks.yaml` — hook message text
  (`"Run: ctx system backup"` → `"Run: ctx backup"`).
- `internal/assets/commands/text/errors.yaml` — two refs to
  `ctx system message list/reset` update to `ctx message list/reset`.
- `internal/assets/claude/CLAUDE.md` — project CLAUDE.md template.
  `ctx system bootstrap` → `ctx bootstrap` in 2 places.
- `internal/assets/claude/skills/ctx-doctor/SKILL.md` — 6 references.
- `internal/assets/claude/skills/ctx-journal-enrich/SKILL.md` — 1 ref.
- `internal/assets/claude/skills/ctx-journal-enrich-all/SKILL.md` — 1 ref.
- `internal/assets/claude/skills/ctx-pause/SKILL.md` — check for refs.
- `docs/cli/system.md` — shrink to plumbing-only.
- New pages: `docs/cli/backup.md`, `docs/cli/bootstrap.md`,
  `docs/cli/event.md`, `docs/cli/message.md`, `docs/cli/prune.md`,
  `docs/cli/resource.md`, `docs/cli/stats.md`.
- `docs/cli/index.md` — nav updates.
- `docs/recipes/customizing-hook-messages.md` — rename all `ctx system
  message *` references.
- `docs/recipes/*` — sweep for any other references.
- `docs/home/*` — sweep for any other references.
- `zensical.toml` — nav for the 7 new CLI pages.
- Project root `CLAUDE.md` — if it references `ctx system bootstrap`,
  update.

## Backwards compatibility

**None.** Per user direction, this is a pure rename with no alias shim.
Release notes call out the rename list explicitly. Users who update
across this change must re-run `ctx init --force` to regenerate their
project CLAUDE.md, and update any hand-written references to the old
commands in their own docs.

Rationale: the alias shim approach was considered and rejected. Shims
add package-level plumbing (deprecation-stub files under `internal/cli/system/cmd/*/`)
that must be tracked and removed in a future release, and the half-state
between "both work" and "only new works" is the worst of both worlds for
users who may be silently running the deprecated path indefinitely. A
hard rename forces the migration.

## Not in scope

- The 27 hidden hook plumbing subcommands (except the 4 plural renames
  in commit 1). Their packages, file locations, and semantics stay.
- The `ctx system` parent itself stays; it just becomes hidden after
  commit 3.
- The internal `internal/write/system/` package taxonomy. That was
  addressed by the 2026-04-XX "Eliminate write/system/" spec in
  `specs/released/v0.8.0/write-system-taxonomy.md`; this spec is about
  the command surface only.
- Any further namespace moves (e.g. grouping `backup`/`bootstrap`/
  `prune` under a `maintenance` parent). If that ever happens, it is a
  separate spec.

## Validation

Each commit passes the hard gate independently:

- `go build ./...` — clean on Linux.
- `GOOS=windows GOARCH=amd64 go build ./...` — clean cross-compile.
- `make lint` — 0 issues.
- `make test` — 0 failures.

Regression coverage: existing tests that invoke commands via their Use
strings update to the new names in the same commit as the Use change.
No new test infrastructure required.

## Out-of-scope follow-ups tracked separately

- Document the new top-level commands in operator-facing recipes that
  don't yet exist (e.g. a "daily maintenance" recipe that chains
  `ctx backup && ctx prune`).
- Consider whether `ctx message` should eventually move under a `ctx
  hook message` umbrella if a `ctx hook` parent is ever introduced for
  hook-related user operations. Not in scope today.
