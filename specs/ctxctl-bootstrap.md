# ctxctl Bootstrap + Audit-Channel Migration

Stand up the long-planned `cmd/ctxctl` maintainer binary
(Phase BT, planned 2026-03, never built) with the out-of-band
audit channel as its first real inhabitant — and move that
channel out of the shipped `ctx` binary, where it does not
belong.

## Problem

The audit channel (specs/audit-channel.md, Phase 1a, commit
`aefce517`) shipped into the `ctx` binary as `ctx audit`
(list/show/dismiss) and the `ctx system check-audit`
UserPromptSubmit hook. Both are **maintainer-only** tooling
mis-housed in the user-facing binary:

- **The hook is a per-prompt tax on every user.** A
  UserPromptSubmit hook in the shipped `hooks.json` fires on
  *every prompt for every ctx user*, doing filesystem reads
  on an empty `.context/audit/` for a feature they will never
  produce reports for. Pure overhead for zero value.
- **The auditor is ctx-specific.** `_ctx-surface-audit`
  (already correctly relocated to `.claude/skills/_*`) scans
  ctx's own `internal/` layout. An end user auditing their
  web app gets nothing from it; the channel exists to serve
  ctx's own development discipline.
- **`ctx audit` bloats the user command surface** with
  subcommands a user has no producer for.

Meanwhile `ctxctl` — the maintainer/contributor binary
specified in TASKS.md Phase BT (line 1387) — has never been
built. Its planned first inhabitants (build/release script
replacements) were each blocked or deferred ("Rewrite
lint-style scripts in Go as ctxctl subcommands — blocked:
prerequisite ctxctl does not exist yet. Deferred."). The
audit channel is a cleaner, self-contained first inhabitant
that forces the binary into existence.

## The Dividing Line (with one refinement)

From Phase BT:

> `ctx` is the user/agent tool, `ctxctl` is the
> maintainer/contributor tool. If a developer clones the
> repo and needs to build, test, release, or validate —
> that's `ctxctl`. If a user is working in a project and
> needs context — that's `ctx`.

Phase BT also states: *"Anything Claude Code hooks call —
hooks must call `ctx`, not `ctxctl`."* This spec **refines**
that rule, because it was written assuming all hooks are
shipped product hooks:

- **Shipped product hooks** (in
  `internal/assets/claude/hooks/hooks.json`, installed by
  `ctx setup`) call `ctx`. Unchanged.
- **Repo-local dev hooks** (wired in the ctx repository's
  own gitignored `.claude/settings.local.json`, never
  shipped) MAY call `ctxctl`. This is the audit-relay hook's
  home.

The distinction is "does this hook reach end users?" Shipped:
`ctx`. Repo-internal: `ctxctl` is fine.

## Approach

### Module structure: same module, `cmd/ctxctl`

`ctxctl` lives at `cmd/ctxctl` in the **same Go module** as
`ctx` (per Phase BT), NOT as a separate module at
`tools/ctx/ctxctl` with its own go.mod. Decided after
weighing both:

- A **separate go.mod** cannot cleanly import the parent
  module's `internal/` packages. The audit channel is ~25
  files already living under `internal/` (`internal/cli/audit`,
  `internal/config/audit`, `internal/err/audit`,
  `internal/write/audit`, `internal/cli/system/core/audit`,
  `internal/cli/audit/core/parse|store`). A module split
  forces relocating or duplicating all of it.
- **Same module still keeps audit out of the `ctx` binary.**
  Go compiles a package into a binary only if that binary's
  `main` transitively imports it. As long as
  `cmd/ctxctl/main` imports the audit packages and
  `cmd/ctx/main` does not, the `ctx` binary never carries
  them. Binary-level isolation — the actual goal — without a
  module split.
- A separate go.mod's only real win is **dependency
  isolation** (keeping heavy build/release deps out of
  `ctx`'s module graph). The audit channel needs no heavy
  deps (only `yaml`, already a ctx dependency). If a future
  ctxctl subcommand pulls in heavy tooling deps, revisit the
  module question *then* — do not pay the split tax now for
  deps that do not exist.

### Move the audit entry points out of `ctx`

The internal logic packages STAY in `internal/` (reused via
import). What moves is the **wiring** — the cobra
registration and the hook entry point:

- Remove `audit.Cmd` registration from
  `internal/bootstrap/group.go` (the `ctx audit` top-level
  command).
- Remove `checkaudit.Cmd()` registration from
  `internal/cli/system/system.go` (the `ctx system
  check-audit` hook subcommand).
- Remove the `check-audit` line from the shipped
  `internal/assets/claude/hooks/hooks.json` — **this
  resolves the deliberately-dirty edit currently in the
  working tree.** (That edit was left dirty as a forcing
  function; this spec is where the trail ends.)
- Remove now-orphaned `ctx`-side constants/yaml that only
  served the `ctx audit` / `ctx system check-audit` surface
  (e.g. `commands.yaml` `audit.*` and `system.checkaudit`
  entries, `UseAudit*`/`DescKeyAudit*` if not reused by
  ctxctl). Re-add under ctxctl's own descriptor namespace.

### Re-expose under ctxctl

- `ctxctl audit list|show|dismiss` — same behavior, reusing
  `internal/cli/audit/cmd/*` and `internal/cli/audit/core/*`.
- `ctxctl audit-relay` — the hook entry, reusing the render
  logic in `internal/cli/system/core/audit`. (Naming: a
  single `audit-relay` verb rather than `system check-audit`,
  since ctxctl has no `system` hook-plumbing namespace.)

### Wire the repo-local dev hook

The ctx repository's own (gitignored)
`.claude/settings.local.json` gets a UserPromptSubmit entry
calling `ctxctl audit-relay`. This makes the channel live for
ctx's own development. End users never see it because it is
neither in the shipped `hooks.json` nor installed by `ctx
setup`.

## Behavior

### Happy path (maintainer, in the ctx repo)

1. Maintainer lands a feature on a branch.
2. From a second Claude Code session: `/_ctx-surface-audit`
   → writes `.context/audit/surface.md`.
3. Back in the working session, the repo-local UserPromptSubmit
   hook fires `ctxctl audit-relay`, which verbatim-relays the
   report in the standard box.
4. Maintainer addresses findings, runs `ctxctl audit dismiss
   surface`.

### End user (any non-ctx project)

Sees no `ctx audit` command, no `check-audit` hook, no
per-prompt audit tax. `ctxctl` is not installed in their
environment and is not referenced by anything `ctx setup`
writes.

## Interface

```
# Maintainer binary, installed via:
go install github.com/ActiveMemory/ctx/cmd/ctxctl@latest
# or built locally by the Makefile.

ctxctl audit                  # list (default)
ctxctl audit list
ctxctl audit show <id>
ctxctl audit dismiss <id>
ctxctl audit dismiss --all
ctxctl audit-relay            # hook entry (reads stdin hook JSON)
```

## Files to Create / Modify

Create:

- `cmd/ctxctl/main.go` — binary entry, cobra root.
- `cmd/ctxctl/<wiring>` — audit + audit-relay subcommand
  registration importing the existing internal packages.
- `internal/config/embed/cmd/ctxctl.go` (or reuse) —
  ctxctl-namespaced Use/DescKey constants.

Modify (move audit out of `ctx`):

- `internal/bootstrap/group.go` — drop `audit.Cmd`.
- `internal/cli/system/system.go` — drop `checkaudit.Cmd()`.
- `internal/assets/claude/hooks/hooks.json` — drop
  `check-audit` (resolves the dirty edit).
- `internal/assets/commands/commands.yaml`,
  `examples.yaml`, `flags.yaml`,
  `internal/config/embed/cmd/{audit.go,system.go}` — relocate
  the `audit.*` / `system.checkaudit` descriptors to ctxctl's
  namespace; remove the `ctx`-side orphans.
- `internal/assets/hooks/messages/registry.yaml` +
  `registry_test.go` count — the `check-audit` message
  variant: decide whether ctxctl reuses the same message
  registry (likely yes — the nudge/message machinery is in
  `internal/`, shared) or gets its own. If reused, the
  registry entry stays and the count test is unaffected.
- `Makefile` — add a `ctxctl` build target.
- `.claude/settings.local.json` (gitignored, local only) —
  wire `ctxctl audit-relay` as a UserPromptSubmit hook.

Keep as-is:

- All `internal/.../audit` logic packages (reused by ctxctl).
- `.claude/skills/_ctx-surface-audit/SKILL.md`.
- `docs/recipes/audit-channel.md` — but re-pass once ctxctl
  exists: the channel is invoked via `ctxctl`, and the hook
  is wired locally, not by `ctx setup`.

## Testing

- `cmd/ctxctl` builds; `ctxctl audit list/show/dismiss` pass
  the same behavioral tests as the Phase 1a CLI (relocate or
  re-point the existing tests).
- **`ctx` binary excludes audit**: a guard test asserting
  `cmd/ctx`'s transitive imports do NOT include
  `internal/cli/audit` (parse the import graph, or a
  `go list -deps ./cmd/ctx` assertion in
  `internal/compliance`).
- Shipped `hooks.json` does NOT contain `check-audit`
  (compliance assertion).
- `ctxctl audit-relay` renders the verbatim box from a
  dropped report (relocate the Phase 1a hook tests).

## Non-Goals

- **The build/release ctxctl subcommands** (`sync`, `build`,
  `release`, `check`, `tag` from Phase BT). This spec only
  bootstraps the binary + migrates the audit channel. Those
  subcommands are a later phase, now unblocked.
- **Audit channel Phase 2** (auto-dismissal, sibling audit
  skills, stale escalation) — unchanged, still tracked under
  specs/audit-channel.md.
- **Shipping ctxctl to end users.** It is a maintainer tool.
  `ctx init` / `ctx setup` must not reference it.

## Open Questions

1. **Do the audit logic packages move from `internal/cli/...`
   to an `internal/ctxctl/...` subtree** to physically signal
   "ctxctl-only," or stay put with a doc-note? Leaning
   stay-put + doc-note: a physical move is churn and the
   import-graph guard test already enforces the boundary
   mechanically. Decide during implementation.
2. **Does the nudge/message registry stay shared** between
   `ctx` and `ctxctl` (both in `internal/`), or does ctxctl
   get its own? Leaning shared — the relay machinery is
   generic and a second copy is waste.
3. **Version stamping for ctxctl** — reuse `ctx`'s version
   package now, or defer until the release subcommands land?
   Leaning reuse-now (trivial, same module).
4. **`ctxctl` install ergonomics in the Makefile** — a
   `make ctxctl` target plus a dev-setup note; final shape
   deferred to the build-tooling phase.

## Source

User session 2026-05-24, immediately after relocating
`_ctx-surface-audit` out of shipped assets. Realization: the
audit channel itself (CLI + hook), not just the skill, is
maintainer tooling — the check-audit hook would tax every
end user's every prompt for a feature they never use. Rather
than ship a half-baked `ctx system` command nobody uses, this
is the forcing function to finally build the Phase-BT
`ctxctl` binary, with the audit channel as its first
inhabitant. The deliberately-dirty `hooks.json` edit in the
working tree is the burned bridge: the only way forward is
the migration this spec describes.
