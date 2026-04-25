# Spec: Explicit Context Directory (No More Walk-Up)

> **Status: SUPERSEDED by `specs/single-source-context-anchor.md`
> (2026-04-24).** The core thesis (no walk-up, explicit
> declaration) stands and is carried forward by the superseding
> spec. Everything below this header is historical record of
> the first iteration and **should not be used as an
> implementation reference**. The body repeatedly treats
> `--context-dir` as live; it is not.
>
> Concrete deltas in the superseding spec — read that one for
> the current design:
>
> - `--context-dir` flag removed entirely. The single declaration
>   channel is `CTX_DIR` (set by `ctx activate` or per-tool hook
>   injection).
> - Basename of `CTX_DIR` (and any override) is locked to
>   `.context` via a Go-side guard; non-canonical values are
>   rejected with a typed error on first use.
> - Per-tool anchor knowledge (e.g. `CLAUDE_PROJECT_DIR`) stays
>   at the shell layer in `hooks.json`, not in Go. The Go
>   resolver knows only `CTX_DIR`.
> - Hook injection hardened: `CTX_DIR="${CLAUDE_PROJECT_DIR:?msg}/.context"`
>   to fail loud when the anchor is empty or unset.
> - Two partial-coverage safety hooks
>   (`block-dangerous-command`, `block-hack-scripts.sh`) removed
>   in favor of native Claude Code `permissions.deny` rules.
> - New `ctx system check-anchor-drift` hook catches stale
>   cross-session `CTX_DIR` bleed.
>
> In particular, the following sections below are **stale**:
> the Approach section's two-step `--context-dir → CTX_DIR`
> priority; error-message examples suggesting
> `ctx --context-dir <path>`; hook template lines executing
> `ctx --context-dir "$CTX_DIR" …`; implementation notes
> implying flag support. Refer to the superseding spec for the
> current shape.

## Problem

`ctx` infers which `.context/` directory to use by walking up from CWD
with a git-root boundary check (`internal/rc/walk.go`). `.ctxrc` is
read from CWD only, with no walk-up. These inference paths silently
pick the wrong project in several real setups:

- **Nested legitimate projects** (e.g. `WORKSPACE/` meta-project
  containing `WORKSPACE/ctx/`). A sub-agent `cd`'d into the inner
  project writes to `WORKSPACE/ctx/.context/` when the session's
  intent was `WORKSPACE/.context/`.
- **Rogue `.context/` directories** created by agents along the walk
  path get adopted silently.
- **Submodules** create a `.git` that isn't the project's real root,
  breaking the boundary check.
- **Subdir `.ctxrc` invisibility**: a hook firing from a
  subdirectory resolves `.context/` correctly via walk-up but never
  sees the project-root `.ctxrc`, silently dropping webhook routing,
  event filters, and other behavior-altering directives.
- **Sub-agent fragmentation**: parent and sub-agent may resolve
  different `.context/` directories from the same session. Tasks,
  decisions, learnings, and webhook events split across projects.

The root cause is that `ctx` tries to *infer* project boundaries from
filesystem shape. This is unanswerable from the filesystem alone —
the information lives in the dispatching session's intent, not on
disk. No in-process signal distinguishes parent agent, sub-agent, or
human caller. See `specs/context-resolution-analysis.md` for the
full analysis.

Upstream fixes (Anthropic issue #26429 —
`CLAUDE_PROJECT_DIR` not propagated into sub-agent env) would help
but are not under our control.

## Approach

**Stop inferring. Require the caller to declare.**

Every `ctx` subcommand except a deliberately small allowlist (v1
scope) refuses to run without an explicit context directory. Resolution becomes a two-step
lookup with no walk, no heuristic, no filesystem guessing:

1. `--context-dir <path>` flag — highest priority, absolute path
   (absolutized if relative).
2. `CTX_DIR` environment variable — session-level convenience.
3. **Error out** with a kind, actionable message.

Ergonomics are preserved via a new `ctx activate` subcommand that
emits `export CTX_DIR=...` for `eval`-style shell integration
(same pattern as `pyvenv activate`, `direnv`, `rbenv`).

`.ctxrc` resolution becomes deterministic: always read from
`dirname(context_dir)/.ctxrc`. No walk-up. Missing file → defaults.

**Configuration belongs to the selected context root, not to the
caller's working directory.** That single shift is the philosophical
core of this change: once the context root is declared, everything
else — config, hooks, state — is anchored relative to it, not to
wherever `ctx` happens to be invoked from.

## The `.context`-is-project-root Contract

**Each project OWNS a `.context/` directory at its project root.
`filepath.Dir(ContextDir())` IS the project root. This is a contract,
not a heuristic.**

Several subsystems read the filesystem relative to the parent of the
context directory:

- `ctx sync` scans `filepath.Dir(ctxDir)` for code directories and
  suggests which to document in ARCHITECTURE.md.
- `ctx drift` scans `filepath.Dir(ctxDir)` for secret-looking files
  (`.env`, `credentials`, `api_key`).
- `ctx check-memory-drift` discovers `MEMORY.md` at
  `filepath.Dir(ctxDir)`.
- `.ctxrc` loads from `filepath.Dir(ctxDir)/.ctxrc`.

If `CTX_DIR` points somewhere decoupled from the codebase (e.g.
`/usr/share/context/my-project/.context` while the code lives at
`~/WORKSPACE/my-project`), every one of the above silently scans the
wrong tree and produces confident-but-wrong results. The parent of
the context directory is not a hint; it is the root.

### Why not support shared `.context` across projects

A shared `.context/` looks attractive but is a false advertisement.
The directory contains:

- `journals/` — per-session, per-developer, append-only; collisions
  lose data.
- `state/` — runtime tombstones keyed by session IDs that are only
  meaningful to one machine.
- `.connect.enc`, `.notify.enc` — per-project secrets (encrypted,
  but still scoped to one project).
- `hub/` — the mechanism that is already designed to be shared;
  sharing the hub itself is a paradox.

Only `CONSTITUTION.md`, `CONVENTIONS.md`, and `ARCHITECTURE.md` want
to be shared across projects, and `ctx hub` serves exactly that use
case at the right granularity: cherry-picked entries, not a bulk
mount. Different developers sharing one `.context/` directly would
need per-developer journal isolation, lock ordering, and CRDT-style
state merging — all to rescue a use case `ctx hub` already solves
better.

The contract:

- One project, one `.context/`, at the project root.
- `filepath.Dir(ContextDir())` is authoritative for "where is the
  code."
- For sharing across projects: `ctx hub`, not a shared `.context/`.
- `CTX_DIR` still points anywhere the user wants; whoever sets
  `CTX_DIR=/somewhere/.context` is declaring `/somewhere/` as their
  project root and owns that choice. We don't police the
  declaration.

### Recommended Project Layout

```
~/WORKSPACE/my-to-do-list
  ├── .git
  ├── .context          ← owned by this project; do not share
  ├── ideas
  │   └── ...
  ├── Makefile
  ├── Makefile.ctx
  └── specs
      └── ...
```

`filepath.Dir(.context)` → `~/WORKSPACE/my-to-do-list` → the code.
Every `ctx` subsystem that needs the project root reads it from
there. No `CTX_PROJECT_DIR` env var. No second source of truth.

Hook scripts generated by `ctx hooks install` honor
`$CLAUDE_PROJECT_DIR` at runtime with a baked-in absolute path as
fallback:

```sh
CTX_DIR="${CLAUDE_PROJECT_DIR:-/absolute/baked/path}/.context"
exec ctx --context-dir "$CTX_DIR" ...
```

## Commands Exempt from the Requirement

These run without `--context-dir` / `CTX_DIR`:

- `ctx init` — creates a new `.context/` at CWD. Takes optional
  `--dir <path>`.
- `ctx activate [path]` — emits shell export. Walks up from CWD if
  no path given (this is the one command where walk-up is tolerable,
  because it runs interactively and crashes loudly on ambiguity).
- `ctx deactivate` — emits `unset CTX_DIR`.
- `ctx version`, `ctx help`, `ctx --help` — trivially exempt.
- `ctx system bootstrap` — runs before context is known; must handle
  both "context declared" and "no context declared" modes
  cleanly.

Everything else (`task`, `decision`, `learning`, `journal`, `status`,
`agent`, `pad`, `remember`, hooks, etc.) errors out without
resolution.

## Behavior Changes

### `internal/rc/walk.go`

**Deleted.** Walk-up logic, `findGitRoot`, and the boundary-check
dance all go away.

### `internal/rc/rc.go:ContextDir()`

New implementation:

```go
func ContextDir() string {
    override := getOverride()
    if override != "" {
        return absolutize(override)
    }
    if envDir := os.Getenv(env.CtxDir); envDir != "" {
        return absolutize(envDir)
    }
    // No inference. Caller declares or we exit.
    exitWithActivationHint()
}
```

### `.ctxrc` Loading (`internal/rc/load.go`)

Read from `filepath.Join(filepath.Dir(ContextDir()), ".ctxrc")`.
Still optional; missing → safe defaults. No walk.

### New Subcommand: `ctx activate`

Interactive shell activator. Emits to stdout:

```
export CTX_DIR=/abs/path/to/.context
```

Resolution rules for activate:

- `ctx activate /abs/path` — emit unconditionally. Validate the
  path exists and looks like a `.context/` directory.
- `ctx activate` (no args) — walk up from CWD collecting every
  `.context/` directory found along the way. Behavior by candidate
  count:
  - **Zero candidates** → error, suggest `ctx init`.
  - **Exactly one candidate** → emit its path.
  - **Two or more candidates** → refuse, print all candidate paths
    with their locations, suggest re-running with an explicit path
    argument.

Activate does **not** classify candidates as "legitimate" or "rogue."
Every `.context/` directory found counts equally. Presence of
multiple is the signal to stop and ask the human — not a cue to try
harder. This keeps activate's resolution mechanism simple enough to
fit in the spec without depending on mechanisms the spec declines to
ship (markers, fingerprints, etc.). Smarter disambiguation is future
work.

**This is the only command allowed to walk.** Every other command
reads `CTX_DIR` or `--context-dir` or errors.

Usage:

```sh
eval "$(ctx activate)"
ctx task add "foo"
```

### New Subcommand: `ctx deactivate`

Emits `unset CTX_DIR`.

### Error Message When Context Dir Is Unset

The message must name the exact subcommand the user was trying to
run and suggest the shortest path to fixing it. When resolution
fails, ctx re-uses the same candidate-counting logic `activate`
would run (read-only; no binding) and tailors the message.

**Zero candidates visible from CWD:**

```
ctx: no context directory specified

No .context/ directory found from this location.
Run: ctx init                # create a new .context/ here
Or:  ctx --context-dir <path> # one-shot override
Or:  export CTX_DIR=<path>    # session-level
```

**Exactly one candidate visible from CWD:**

```
ctx: no context directory specified

From this directory, the likely candidate is /repo/.context
Run: eval "$(ctx activate)"  # bind for this shell
Or:  ctx --context-dir /repo/.context <command>
```

**Two or more candidates visible:**

```
ctx: no context directory specified

Multiple .context/ directories visible from this location:
  /repo/.context
  /repo/sub/.context
Run: ctx activate /abs/path  # choose one explicitly
```

The hint is informational only. It never auto-selects — that
decision stays with the human. The candidate list is a read, not a
bind.

### Hook Script Template

`ctx hooks install` generates scripts of the form:

```sh
#!/bin/sh
# Generated by ctx hooks install
CTX_DIR="${CLAUDE_PROJECT_DIR:-/absolute/path/baked/at/install/time}/.context"
exec ctx --context-dir "$CTX_DIR" <hook-command> "$@"
```

Behavior:

- Under Claude Code, `CLAUDE_PROJECT_DIR` is set in the hook
  payload env → that wins, portable through repo moves.
- Under other harnesses or raw invocation, the baked path wins.
- If the repo moves *and* `CLAUDE_PROJECT_DIR` is unset, the baked
  path is stale. Document `ctx hooks reinstall` as the fix.

**Hook failure semantics:**

- If the resolved `CTX_DIR` does not exist or is not a valid
  `.context/` directory, the hook **must fail loudly** (non-zero
  exit + stderr message naming both the resolved path and its
  source — `CLAUDE_PROJECT_DIR` or baked) and emit a
  `ctx hooks reinstall` hint.
- Hooks **never fall back to walk-up.** Walk logic lives in exactly
  one place (`ctx activate`). Silent drift is precisely the bug this
  whole change is fixing — reintroducing it inside hooks would
  undo the main guarantee.
- Failing loudly in hooks is acceptable UX because hooks surface
  errors to the harness's event stream (Claude Code shows them to
  the user). A broken hook is visible; a silently-wrong hook is
  what we're eliminating.

## Non-Goals

- **No deprecation path.** Clean break. Existing users must update
  their workflow after upgrading. Bumping to 0.9 (or equivalent)
  signals the break.
- **No `include:` directive** for `.ctxrc`. If two projects need to
  share webhook config, they duplicate it. A little copy is better
  than a little dependency.
- **No `.context/PROJECT` marker file.** Once inference is gone, the
  primary motivation for markers (distinguishing real projects from
  rogue dirs during walk-up) evaporates. Markers may return later
  for `ctx activate`'s ambiguity detection (see Future Work) but are
  not required for the main change.
- **No observability audit log** (the sidecar `~/.ctx/audit.jsonl`
  idea from the analysis doc). With explicit resolution, fragmentation
  becomes a caller-level concern, not a ctx-level one — if the caller
  passes the wrong path, that's their bug to fix.
- **No change to commands that don't read/write `.context/`.**
  Exempt list stays small and obvious.
- **No automatic migration** of existing `.ctxrc` locations (if any
  users have them at unusual paths). Adjacent-to-`.context/` is the
  only supported location after this change.

## Decisions

Resolved during spec review:

1. **`ctx activate` refuses on ambiguity.** Multiple candidates →
   refuse + print the list. No silent innermost-wins, no rogue
   classification (see `ctx activate` section).
2. **`ctx init` prints the activation hint.** After creating
   `.context/`, output `eval "$(ctx activate /abs/path)"` as the
   next step. Closes the loop for new projects without requiring
   two commands' worth of manual work.
3. **`ctx system bootstrap` does not walk.** Walk logic lives only
   in `ctx activate`. Bootstrap either reports the declared
   `CTX_DIR` / `--context-dir`, or reports "unset" — same as any
   other command.
4. **`CTX_DIR` env var name stays.** Already used across the
   codebase (`internal/config/env/env.go:16`). Reasoning:
   - Alternatives like `CTX_CONTEXT_DIR` stutter — "ctx context"
     is two words for the same thing. `CTX_DIR` is already scoped
     by the `CTX_` prefix; the second token should add meaning,
     not repeat it.
   - This change already breaks every existing user once. Renaming
     the env var at the same time doubles the friction for no
     architectural gain.
   - The reviewer is right that env names calcify after shipping —
     but `CTX_DIR` is already shipped. The decision to make now is
     "keep" rather than "choose fresh."
5. **`ctx activate` supports bash/zsh in v1.** Structure the
   emitter as a shell-keyed dispatch (`emit(shell) → string`) so
   adding fish, nushell, powershell later is a plugin-style
   addition, not a redesign. v1 ships bash/zsh; others on request.
6. **Validation for `ctx activate /abs/path` is strict, no
   `--force`.** The path must exist, be a directory, and contain
   at least one canonical file (`CONSTITUTION.md` or `TASKS.md`).
   If any check fails, refuse with a message that names the missing
   piece. No `--force` escape hatch in v1 — adding one multiplies
   the branches on the explicit-path path and eats into the
   simplicity win. If a real use case for forcing appears later,
   add the flag then.

## Open Questions (Remaining)

1. **Subdir writes inside a project:** a human `cd`'d into
   `/repo/src/foo` with `CTX_DIR=/repo/.context` exported via
   activate should work uniformly. Verify no subcommand assumes
   `CWD == project root` beyond what ctx explicitly resolves
   through `ContextDir()`. This is a code-audit task, not a design
   decision — flagging it so implementation doesn't skip it.
(Note: validation rules for `ctx activate /abs/path` were an open
question in the first draft. Resolved in Decisions §6.)

## Implementation Sketch

Rough file-level change set (not final):

- **Delete**: `internal/rc/walk.go`, associated tests.
- **Rewrite**: `internal/rc/rc.go:ContextDir()`, `KeyPath()` and any
  downstream that assumed walk-up behavior.
- **Rewrite**: `internal/rc/load.go` to compute `.ctxrc` path from
  resolved context dir rather than CWD.
- **New**: `internal/cli/activate/` subcommand with resolver that
  walks (only here) and emits shell exports.
- **New**: `internal/cli/deactivate/` subcommand.
- **Update**: hook install template in
  `internal/assets/claude/hooks/` (and equivalents for other
  integrations) to the `${CLAUDE_PROJECT_DIR:-baked}` pattern.
- **Update**: every entry-point command to call a shared
  "require-context-dir-or-exit-kindly" helper. Exempt list lives in
  one place.
- **Update**: error messages, `ctx --help` output, `README.md`,
  `AGENT_PLAYBOOK.md`, `CLAUDE.md` project instructions.
- **Delete**: tests that exercised walk-up scenarios. Replace with
  tests that exercise the exempt-list boundary and the error message.
- **Add**: integration tests for `ctx activate` in a fresh dir,
  in a nested dir, and in a dir with ambiguous candidates.

## Risks

- **User friction on first session after upgrade.** Users running
  `ctx` in a plain shell will hit the error immediately. Mitigation:
  kind error message that points to `ctx activate`. This is the
  primary discomfort; it is accepted.
- **Hook scripts with baked paths go stale after repo moves.**
  Mitigation: `CLAUDE_PROJECT_DIR` fallback at runtime handles the
  common case; `ctx hooks reinstall` handles the edge. Document both.
- **Third-party scripts that invoke `ctx` in CI or cron.** They must
  set `CTX_DIR` explicitly. Documented migration step.
- **Loss of "just works" feel for new users.** Mitigation: `ctx init`
  prints the activate command as a hint; first-run UX closes the loop
  immediately.

## Future Work

- **Smarter disambiguation / candidate classification in `ctx
  activate`.** v1 refuses on any multi-candidate case and lists
  them. A future version may classify candidates (e.g., prefer
  ones backed by a `.context/PROJECT` fingerprint, surface stale
  or orphaned dirs separately) so `activate` can recommend a
  default without losing the "refuse loud" guarantee. Not required
  for the main change.
- **`.context/PROJECT` marker file** — fingerprinted with the
  enclosing repo's git remote — as the machinery behind the
  smarter disambiguation above. Explicitly deferred from v1.
- **Fish / nushell / powershell support in `activate`.** The v1
  shell-keyed emitter is designed to make these a drop-in addition.
- **Per-tool session probe for non-Claude harnesses** (Cursor,
  Cline, Kiro, Codex). Out of scope here; relevant for future
  hub/federation work.
