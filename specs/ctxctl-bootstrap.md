# ctxctl Bootstrap + Audit-Channel Migration

Stand up the long-planned `ctxctl` maintainer binary at
`tools/ctxctl` — its own Go module (Phase BT, planned 2026-03,
never built) — with the out-of-band audit channel as its first
real inhabitant, and move that channel out of the shipped `ctx`
binary, where it does not belong.

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

### Module structure: separate module at `tools/ctxctl`

`ctxctl` is its **own Go module** at `tools/ctxctl` (module
path `github.com/ActiveMemory/ctx/tools/ctxctl`), NOT
`cmd/ctxctl` in the same module. This **reverses** the earlier
same-module decision (handover 2026-05-26); see DECISIONS.md
2026-05-27.

The earlier decision rested on a false premise — that a
separate go.mod cannot import the parent module's `internal/`
packages. **Empirically disproved** this session: a nested
module whose path is lexically under
`github.com/ActiveMemory/ctx` builds clean while importing
`…/ctx/internal/…`; only a non-nested ("outsider") module path
is rejected. Go's `internal` visibility is import-path-lexical,
not module-scoped.

With that blocker gone, the deciding axis is **blast radius**,
not binary size:

- **Hard boundary, the right direction.** `ctx`'s `go.mod`
  does not `require` `tools/ctxctl`, so `ctx` *literally
  cannot* import ctxctl — enforced by the module graph, not a
  test. ctxctl breaking can never break ctx. The
  one-directional `ctxctl → ctx` coupling is fine: ctxctl is
  disposable maintainer tooling.
- **Reuse, don't duplicate.** ctxctl imports ctx's shared
  `internal/` foundations (`rc`, `assets/read/desc`,
  `config/*`, `io`, `cli/system/core/nudge`, …) in place via
  the nested-module internal allowance. Full self-containment
  (copying those ~20 foundation packages into ctxctl) is
  rejected as a DRY catastrophe — a worse broken window than
  the one this spec fixes.
- **Local dev via `go.work`.** A repo-root `go.work` (`use .`
  and `use ./tools/ctxctl`), committed, wires the workspace so
  every maintainer gets cross-module build/nav/refactor with
  zero setup. Safe for releases: `cmd/ctx` never imports
  ctxctl, so workspace mode cannot pull ctxctl into the shipped
  `ctx` binary.

### Relocate the audit logic, strip the `ctx` wiring

Two moves:

**1. Relocate audit-specific packages → `internal/ctxctl/`.**
The six audit-channel trees (33 files) move out of their
current `internal/cli/...`, `internal/config/...`,
`internal/err/...`, `internal/write/...` homes into an
`internal/ctxctl/` subtree, physically signalling
"ctxctl-only." Shared foundations (`rc`, `desc`, `config/*`
primitives, `nudge`, `io`, …) stay put and are imported by both
binaries. `git mv` + import-path rewrite; a guard test asserts
`cmd/ctx` never imports `internal/ctxctl`.

**2. Strip the `ctx`-side wiring:**

- Remove `audit.Cmd` registration from
  `internal/bootstrap/group.go` (the `ctx audit` command).
- Remove `checkaudit.Cmd()` registration from
  `internal/cli/system/system.go` (the `ctx system
  check-audit` hook).
- Remove the `check-audit` line from the shipped
  `internal/assets/claude/hooks/hooks.json` — **this resolves
  the deliberately-dirty edit in the working tree.** (Left
  dirty as a forcing function; this spec is where the trail
  ends.)
- **Delete** the orphaned `ctx`-side audit descriptors
  outright (not relocate): `commands.yaml` `audit.*` /
  `system.checkaudit`, the `examples.yaml` / `flags.yaml`
  counterparts, `internal/config/embed/text/check_audit.go`,
  the `UseAudit*`/`UseSystemCheckAudit` constants in
  `internal/config/embed/cmd`. ctxctl does **not** reuse these
  — it owns its text as English Go constants (see below), so
  the i18n keys leave the codebase entirely.

### Re-expose under ctxctl

The `tools/ctxctl` module's `main` wires:

- `ctxctl audit list|show|dismiss` — same behavior, importing
  the relocated `internal/ctxctl/...` packages.
- `ctxctl audit-relay` — the hook entry, reusing the relocated
  render logic and ctx's `cli/system/core/nudge` box renderer.
  Single `audit-relay` verb rather than `system check-audit`,
  since ctxctl has no `system` hook-plumbing namespace.

User-facing text (CLI output, relay-box copy) is supplied as
plain English Go constants in `tools/ctxctl`, passed into the
(text-free) `internal/ctxctl` logic. The relocated logic holds
no hardcoded user copy and makes no `desc.Text` calls — that
machinery stays in `ctx` for `ctx`'s own output.

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
# Maintainer binary — built to dist/ and installed to PATH
# (symmetric with ctx) so every repo copy / worktree shares
# one binary and the repo root stays clean:
make reinstall-ctxctl         # -> /usr/local/bin/ctxctl

ctxctl audit                  # list (default)
ctxctl audit list
ctxctl audit show <id>
ctxctl audit dismiss <id>
ctxctl audit dismiss --all
ctxctl audit-relay            # hook entry (reads stdin hook JSON)
```

## Files to Create / Modify

Create:

- `tools/ctxctl/go.mod` — module
  `github.com/ActiveMemory/ctx/tools/ctxctl`; `require` the
  `ctx` module (resolved locally via `go.work`).
- `tools/ctxctl/main.go` + wiring — cobra root registering
  `audit list|show|dismiss` and `audit-relay`, importing the
  relocated `internal/ctxctl/...` packages.
- `tools/ctxctl/<text>.go` — English-string constants for all
  ctxctl user-facing output (CLI + relay-box copy).
- `go.work` (repo root, committed) — `use .` and
  `use ./tools/ctxctl`; plus `go.work.sum`.

Move (relocate with `git mv` + import-path rewrite):

- `internal/cli/audit/**`, `internal/cli/system/cmd/checkaudit/**`,
  `internal/cli/system/core/audit/**`, `internal/config/audit/**`,
  `internal/err/audit/**`, `internal/write/audit/**`
  → under `internal/ctxctl/...`. Make the output functions
  text-free (accept strings as parameters); strip their
  `desc.Text` calls.

Modify (strip audit out of `ctx`):

- `internal/bootstrap/group.go` — drop `audit.Cmd`.
- `internal/cli/system/system.go` — drop `checkaudit.Cmd()`.
- `internal/assets/claude/hooks/hooks.json` — drop
  `check-audit` (resolves the dirty edit).
- `Makefile` — add a `ctxctl` build target (`make ctxctl`).
- `.claude/settings.local.json` (gitignored, local only) —
  wire `ctxctl audit-relay` as a UserPromptSubmit hook.

Delete (orphaned `ctx`-side audit descriptors — ctxctl owns
its own English text, so these leave the codebase):

- `internal/config/embed/text/check_audit.go`; the
  `UseAudit*` / `UseSystemCheckAudit` (+ `DescKey*`) constants
  in `internal/config/embed/cmd`.
- `internal/assets/commands/commands.yaml` `audit.*` /
  `system.checkaudit`; matching `examples.yaml`, `flags.yaml`,
  and `internal/assets/hooks/messages/registry.yaml` entries
  (adjust the `registry_test.go` count).

Keep as-is:

- Shared `internal/` foundations (`rc`, `assets/read/desc`,
  `config/*`, `io`, `cli/system/core/nudge`, …) — imported by
  both binaries.
- `.claude/skills/_ctx-surface-audit/SKILL.md`.
- `docs/operations/runbooks/audit-channel.md` (relocated from
  `docs/recipes/` — maintainer-only, so it lives with the
  contributor runbooks) — re-pass once ctxctl exists (invoked via
  `ctxctl`; hook wired locally, not by `ctx setup`).

## Testing

- `tools/ctxctl` builds (workspace mode); `ctxctl audit
  list/show/dismiss` pass the same behavioral tests as the
  Phase 1a CLI (relocate the existing tests into the module).
- **`ctx` cannot reach audit — two layers:** (a) the module
  graph — `ctx`'s `go.mod` does not `require` `tools/ctxctl`;
  (b) a guard test asserting `cmd/ctx`'s transitive imports
  exclude `internal/ctxctl` (`go list -deps ./cmd/ctx` in
  `internal/compliance`).
- Shipped `hooks.json` does NOT contain `check-audit`
  (compliance assertion).
- `ctxctl audit-relay` renders the verbatim box from a dropped
  report (relocate the Phase 1a hook tests), with the English
  copy supplied by `tools/ctxctl`.

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

1. **Relocate audit logic to `internal/ctxctl/`?** —
   **RESOLVED: yes, move it.** Physical relocation signals
   "ctxctl-only" and pairs with the module boundary; stay-put
   was the lazy option. (DECISIONS.md 2026-05-27.)
2. **Shared message registry?** — **RESOLVED: no.** ctxctl
   owns its user-facing text as plain English Go constants
   under `tools/ctxctl`, outside the YAML localization and the
   `desc`/i18n engine — no French ctxctl. The relay machinery
   (`nudge` box renderer) is still reused; it receives
   already-resolved strings.
3. **Version stamping for ctxctl** — reuse `ctx`'s version
   package (importable via the nested-module allowance) now,
   or defer until the release subcommands land? Leaning
   reuse-now (trivial).
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

**Revised 2026-05-27 (session 96765858):** reversed from
same-module `cmd/ctxctl` to a separate module at
`tools/ctxctl`, after a build test disproved the "separate
go.mod can't import `internal/`" premise. Driver: blast-radius
isolation — a hard module boundary so ctxctl can never break
ctx — plus ctxctl owning its English-string messages (no
i18n). See DECISIONS.md, 2026-05-27.
