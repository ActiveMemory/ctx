---
title: CWD-Anchored Context
status: proposed
date: 2026-05-20
owner: jose
scope: architectural — resolver, activate, deactivate, hooks, gitmeta, specs, docs
supersedes-section-of:
  - specs/single-source-context-anchor.md (the `CTX_DIR` declaration model
    and the `ctx activate` carve-out; the basename guard and the
    no-walk-up rule survive in stronger form)
  - specs/activate-strict-cwd.md (entirely; this spec subsumes it
    by removing `ctx activate` rather than refining its resolver)
related:
  - specs/require-git.md
brief-source: |-
  Debated in session d4cb8647 (2026-05-20). The conversation
  followed `/ctx-plan` shape informally: agent proposed walk-up
  to `.git/`, user countered with strict-CWD, agent pushed back
  citing subdir convenience, user countered with the
  gitmeta-collapse insight and the zensical / Claude Code
  precedent (tools that anchor to a specific config location).
  Convergence: CWD-anchored is strictly simpler across the
  codebase, the convenience tax is acceptable for a tool that
  is closer to claude-code than to git in usage shape.
---

# Spec: CWD-Anchored Context

## Problem

`ctx` currently resolves the target `.context/` through two
channels: an explicit `CTX_DIR` env var (validated for basename,
abs path, declared loud) and `ctx activate`'s upward scan from
CWD (since `specs/activate-strict-cwd.md`: a strict $PWD/.context
check). The two-channel model leaves several recurring frictions:

1. `eval $(ctx activate)` is a per-shell ceremony every user
   learns once, forgets, re-learns, then asks "why do I need
   this?" The eval-around-export pattern is opaque to anyone who
   has not specifically read the `ctx activate` docs.
2. Cross-shell state: `CTX_DIR` set in one shell does not carry
   to another, so the user re-runs `eval` per terminal tab.
3. Cross-project state: `CTX_DIR` set for project A persists
   when the user `cd`s into project B. Mitigated by the env-vs-CWD
   mismatch guard added in `ctx init` (`internal/cli/initialize/core/envmatch/`,
   landed earlier this session), but the mitigation is symptomatic;
   the underlying friction is that the env var is a parallel
   declaration channel to CWD.
4. Walk-up infrastructure exists in two places: `rc.ScanCandidates`
   (no longer called by `ctx activate` after strict-CWD) and
   `gitmeta.RequireGitTree` (Phase RG, used for handover provenance).
   Both walks are bounded by `.git/`, but they exist because
   commands could be invoked outside the project root.

The unifying observation: `.context/` is mandated to be sibling
of `.git/` by `CONSTITUTION.md`'s `require-git` rule. If we also
mandate that `ctx` runs at the project root (where both `.context/`
and `.git/` are), every resolver collapses to a single `os.Stat`
and the entire declaration / walk / activate machinery becomes
unnecessary.

## Approach

**`ctx` is anchored to CWD. `$PWD/.context/` is the context;
`$PWD/.git/` is the git tree. No env var, no walk, no activate.**

Mental model: ctx anchors to its working directory the way
zensical anchors to `zensical.toml`, helm to `Chart.yaml`,
terraform to `.tf` files, Claude Code to `$CLAUDE_PROJECT_DIR`.
Tools that manage local project state are anchored to a config
file or directory; ctx's anchor is the `.context/` directory
itself.

### Why CWD-anchored, not walk-to-`.git/`

Walk-to-`.git/` was the agent's first counter-proposal and is
strictly more convenient for subdir work (no `cd` required).
The user countered with three arguments that decided it:

1. **The walk infrastructure is what we are trying to delete.**
   If gitmeta still walks for `.git/`, we keep the walk code; if
   rc walks too, we maintain two implementations. CWD-anchored
   collapses both into `os.Stat`.
2. **Agents pay no convenience tax.** Most ctx invocations come
   from AI agents (Claude Code, etc.) where the agent runs
   `git rev-parse --show-toplevel` mechanically without
   annoyance. Humans who care about subdir ergonomics keep a
   terminal tab open at the project root or alias `cd $repo`.
3. **The mental model is already paid for** by Claude Code,
   which has the same anchor-to-project-root shape. Users of
   ctx are also users of an AI coding tool with the same
   discipline.

The subdir-convenience tax exists, but it is a fixed cost paid
per-shell, not per-invocation: once you are at the project root
(typical session shape), every subsequent ctx call works.

## Behavior

### Happy Path

```bash
cd ~/projects/foo                    # has .git/ and .context/
ctx status                           # works
ctx add task "..."                   # works
ctx commit                           # works (git operations also from here)
```

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `$PWD/.context/` absent | Refuse with [`errCtx.NoContextHere`] naming `$PWD`. Suggest `ctx init` or `cd` elsewhere. |
| `$PWD/.git/` absent | Refuse with `gitmeta.ErrMissingGitTree` (already enforced by Phase RG). |
| `$PWD/.context/` exists but is a regular file (not a directory) | Refuse with the same `NoContextHere` error; the basename string is `.context` but the type is wrong. |
| `$PWD/.context/` exists, `$PWD/.git/` absent | The git-required gate fires first (existing behavior). |
| `ctx init` in a fresh `git init`'d directory with no `.context/` | Init creates `.context/` and succeeds. No env-var resolution gate. |
| Hook subprocess with unreliable CWD | Hook script must `cd "${CLAUDE_PROJECT_DIR:?missing}"` (or equivalent for the host tool) before invoking `ctx`. Loud failure on empty `CLAUDE_PROJECT_DIR`. |
| CI replay (`CTX_TASK_COMMIT` / `GITHUB_SHA` set) | Unchanged. These env vars override the resolved git HEAD for handover provenance only; they do not influence context-dir resolution. |

### Validation Rules

- `.context/` must exist as a directory at `$PWD/.context`. No
  symlink resolution, no basename inference, no path
  normalization beyond `filepath.Clean($PWD)`.
- Hook scripts must `cd` before calling `ctx`. The convention is
  documented in the hook templates; the loud-fail expansion
  `${VAR:?msg}` is preserved.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `$PWD/.context/` missing | `ctx: no .context/ at <pwd>. Run \`ctx init\` here, or cd to a project that has one.` | Run `ctx init` or `cd`. |
| `$PWD/.git/` missing | (existing) `ctx requires a git working tree; <pwd> has no .git/. Run \`git init\` first.` | Run `git init`. |
| `$PWD/.context` exists but is a file | `ctx: .context at <pwd> is not a directory.` | `rm .context && ctx init`. |

## Interface

### CLI

No new flags. The following are **removed**:

- `ctx activate` — the entire command and its subtree.
- `ctx deactivate` — the entire command and its subtree.
- `CTX_DIR` env var as a resolution channel (the var name may
  remain reserved to refuse it loud if set; tbd in Open Questions).

The following commands are unchanged in user-visible behavior
but **simplified** internally:

- Every read/write command (`ctx status`, `ctx add`, `ctx
  commit`, etc.) calls a renamed `rc.ContextDir()` that now
  returns `filepath.Join(cwd, ".context")` after a single
  `os.Stat`.
- `ctx init` drops the env-vs-CWD mismatch guard (added earlier
  this session for TASKS.md line 63) because the mismatch class
  can no longer occur.

### Skill

No new skills. The trigger map in CLAUDE.md, AGENT_PLAYBOOK,
and the per-tool skill catalogues drop every reference to
"activate" / "deactivate" / `eval $(ctx activate)`.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/rc/rc.go` | `ContextDir()` returns `filepath.Join(cwd, dir.Context)` after a single `os.Stat`; refuse with the new `errCtx.NoContextHere(cwd)` typed error otherwise. Drop env-var read, drop basename guard, drop `ErrDirNotDeclared` / `ErrRelativeNotAllowed` / `ErrNonCanonicalBasename`. |
| `internal/rc/candidates.go`, `internal/rc/require.go` | Drop `ScanCandidates` and its only remaining caller path; `Require` becomes a thin wrapper around the new `ContextDir` shape. |
| `internal/gitmeta/require.go` | `RequireGitTree` becomes `os.Stat($PWD/.git)`; drop the upward walk. |
| `internal/cli/activate/` | Delete entirely. Update `internal/cli/cli.go` registration. |
| `internal/cli/deactivate/` | Delete entirely. Update registration. |
| `internal/write/activate/` | Delete entirely. |
| `internal/err/activate/` | Delete entirely. |
| `internal/cli/initialize/cmd/root/run.go` | Drop the env-vs-CWD mismatch branch; drop the `rc.ContextDir` fallback dance — `ctx init` always targets `$PWD/.context`. |
| `internal/cli/initialize/core/envmatch/` | Delete entirely (the guard becomes vestigial). |
| `internal/err/initialize/initialize.go` | Drop `ErrEnvCwdMismatch` and `EnvCwdMismatch`. |
| `internal/config/embed/text/initialize.go` | Drop the env-cwd-mismatch text keys. |
| `internal/assets/commands/text/errors.yaml` | Drop `err.init-env-cwd-mismatch` and `err.init.env-cwd-mismatch-msg`. |
| `internal/config/env/...` | Drop `CtxDir` constant (or retain only for the "refuse if set" guard; see Open Questions). |
| `internal/config/file/ignore.go` | Drop `.context/.ctx.key` and friends only if they referenced activate state (most don't). |
| `internal/assets/hooks/` (every hook script) | Replace `CTX_DIR="$CLAUDE_PROJECT_DIR/.context"` injection with `cd "${CLAUDE_PROJECT_DIR:?missing CLAUDE_PROJECT_DIR for ctx hook}" \|\| exit 1` at the top of the script. |
| `internal/assets/integrations/copilot-cli/scripts/` | Same `cd` migration for every `.sh` / `.ps1` script. |
| `internal/assets/claude/CLAUDE.md`, `AGENT_PLAYBOOK.md`, `AGENT_PLAYBOOK_GATE.md` | Drop every "Activation" section, every `eval $(ctx activate)` reference, every "Channels for declaring `CTX_DIR`" table. |
| `docs/recipes/activating-context.md` | Either delete (no longer activation-worthy) or rewrite as "Pinning Your Project". |
| `docs/home/getting-started.md` | Drop the "activate the context" step. |
| `specs/activate-strict-cwd.md` | Mark superseded by this spec (preamble note). |
| `specs/single-source-context-anchor.md` | Mark relevant sections superseded; the basename-guard rationale and the no-walk-up principle survive (the latter in stronger form). |

### Key Functions

```go
// internal/rc/rc.go
func ContextDir() (string, error) {
    cwd, err := os.Getwd()
    if err != nil {
        return "", err
    }
    candidate := filepath.Join(cwd, dir.Context)
    info, statErr := os.Stat(candidate)
    if statErr != nil {
        if errors.Is(statErr, os.ErrNotExist) {
            return "", errCtx.NoContextHere(cwd)
        }
        return "", statErr
    }
    if !info.IsDir() {
        return "", errCtx.NotADirectory(candidate)
    }
    return candidate, nil
}

// internal/gitmeta/require.go
func RequireGitTree(_ string) error {
    cwd, err := os.Getwd()
    if err != nil {
        return err
    }
    if _, statErr := os.Stat(filepath.Join(cwd, ".git")); statErr != nil {
        return errGitmeta.ErrMissingGitTree
    }
    return nil
}
```

(The `string` parameter is preserved for signature compatibility
during the transition; it is ignored. Callers can drop it on a
follow-up sweep.)

### Helpers to Reuse

- `internal/config/dir.Context` for the literal `.context`
  basename.
- `internal/err/context/...` for the typed-error pattern (new
  `NoContextHere`, `NotADirectory` sentinels follow the existing
  `entity.Sentinel` shape).
- `internal/assets/read/desc.Text` for the user-facing message
  bodies, sourced from `internal/assets/commands/text/errors.yaml`.

## Configuration

- `CTX_DIR` env var: removed from the resolver. Open Question:
  retain a "refuse if set" guard so users with stale shellrc
  exports get a loud error rather than silent misconfiguration?
- `.ctxrc`: unaffected. The rcfile lives at `$PWD/.ctxrc` and
  has always been CWD-anchored.

## Testing

- **Unit**:
  - `rc.ContextDir()` returns the cwd `.context` when present,
    refuses when absent.
  - `gitmeta.RequireGitTree` accepts when `$PWD/.git` present,
    refuses when absent.
- **Integration**:
  - `ctx status` from project root succeeds.
  - `ctx status` from any subdir refuses with the typed error
    naming `$PWD`.
  - `ctx init` in a `git init`'d empty dir creates `.context/`
    and succeeds.
  - `ctx init` in a dir that already has populated `.context/`
    refuses per existing safety (unchanged).
- **Edge cases**:
  - `$PWD/.context` is a regular file → refuse.
  - `$PWD/.context` is a symlink to a real dir → accept
    (filesystem behaves transparently; no special handling).
- **Migration tests**:
  - Hook subprocess with `$CLAUDE_PROJECT_DIR` set: `cd` works,
    ctx finds `.context`.
  - Hook subprocess with `$CLAUDE_PROJECT_DIR` empty / unset:
    fail loud per the `${VAR:?msg}` expansion.

## Non-Goals

- **Backward compatibility for `eval $(ctx activate)`** in
  user shellrc files. The eval becomes a no-op error
  (`unknown command: activate`). Migration note in release
  notes; no shim.
- **Multi-context projects.** Single `.context/` per project,
  same as today.
- **Subdir support via fallback walks.** Deliberately rejected.
  Users who want subdir invocation use `cd $(git rev-parse --show-toplevel)`
  or a shell alias.
- **Removing the basename guard from explicit env declarations.**
  Moot: there are no explicit env declarations under this model.
- **Changing the `.context/` ↔ `.git/` sibling contract.**
  This spec depends on it; do not weaken.

## Open Questions

1. **`CTX_DIR` env-var presence: refuse or ignore?**
   If the user has `export CTX_DIR=/some/path` in shellrc from
   a prior ctx version, should the new ctx (a) silently ignore
   it, (b) print a one-line deprecation notice on every
   invocation, or (c) refuse loud and tell them to remove it
   from shellrc? Leaning toward (b) for one release cycle, then
   (a) — but the choice affects user upgrade friction.
2. **Activate-shaped backward-compat shim?** Should `ctx activate`
   exist as a deprecation message ("you no longer need this;
   delete `eval $(ctx activate)` from your shellrc") for one
   release, or is the hard removal cleaner?
3. **VSCode extension and other editor integrations.** Do any of
   them rely on the env-var declaration? Need to grep
   `editors/vscode/` and confirm before the implementation
   starts.
4. **Multi-step implementation order.** Likely (1) rc + gitmeta
   resolver simplification, (2) init guard removal, (3) hook
   `cd` migration, (4) activate / deactivate deletion, (5) docs
   sweep. Each step CI-green before the next. Confirm before
   implementation kicks off.
