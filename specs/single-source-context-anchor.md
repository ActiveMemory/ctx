---
title: Single-Source Context Anchor
status: proposed
date: 2026-04-24
owner: jose
scope: architectural — resolver, hooks, settings, deletions
supersedes:
  - specs/rc-contextdir-upward-walk.md
  - specs/explicit-context-dir.md
related:
  - specs/hook-guard-uninitialized.md
  - specs/deprecate-ctx-backup.md
---

# Spec: Single-Source Context Anchor

## Resume here (for fresh session implementers)

Already landed in the working tree before implementation
starts:

- This spec (`specs/single-source-context-anchor.md`).
- Supersession preambles added to
  `specs/rc-contextdir-upward-walk.md` and
  `specs/explicit-context-dir.md`. The bodies are intentionally
  unchanged; the preambles tell readers the bodies are
  historical.
- `internal/assets/claude/skills/ctx-plan/SKILL.md` —
  the adversarial-interview skill that produced this spec.
  **Not** an implementation artifact of Phase SC; leave it.
- `.context/TASKS.md` Phase SC ("Single-Source Context
  Anchor") — full task list, each entry linked back to this
  spec. Implementation work goes against those tasks.

Recommended implementation order (dependencies respected):

1. **Section A** — resolver split (`ContextDir` shape guard
   + `RequireContextDir` boundary check) and new typed
   errors. Foundational; everything downstream depends on it.
2. **Section B** — flag and `OverrideContextDir` removal.
   Pairs with A; no other section can land cleanly until the
   override path is gone.
3. **Section C / C-bis** — hook injection hardening.
   Independent of A/B but trivial; do it next so smoke can
   exercise it.
4. **Section F** — `check-anchor-drift` hook (Go subcommand
   + registration with `CTX_DIR_INHERITED` prepend).
5. **Section G** — `ctx activate` simplification (drop
   explicit-path mode + stale-replacement comment).
6. **Section H** — hub commands `AnnotationSkipInit`. Can
   land at any point but cheapest after A is in.
7. **Section D / E** — settings.local.json sweep + deletions
   (`block-dangerous-command`, `block-hack-scripts.sh`).
   Last because deletions are easiest to do in bulk after the
   replacement deny rules are agreed and grep-verified.
8. **Tests** alongside each chunk per the test plan.
9. **Documentation sweep** per the grep target in the
   "Documentation impact" section.
10. **Pre-commit smoke** (the human-as-pushbutton checklist
    near the end of this spec) — required before final
    commit.

One bulk PR per the branch strategy. Run `make lint && make test`
before committing.

## Problem

The `feat/explicit-context-dir` branch correctly diagnosed silent walk-up
resolution as the root cause of context-dir bugs (stray writes, wrong
project picked, sub-agent fragmentation), and replaced it with an
explicit-declaration model: `--context-dir` flag or `CTX_DIR` env, no
inference, error otherwise.

In practice, three failure modes survived the rewrite:

1. **Override-points-at-root.** `--context-dir=.` (or `=$PWD`,
   `=$(pwd)`) resolves verbatim, then `ctx init` deposits
   `TASKS.md`, `DECISIONS.md`, `AGENT_PLAYBOOK.md`, etc. directly
   into the project root. The basename was never validated.
2. **Per-tool env injection lives in shell glue but is invisible to
   Go.** `internal/assets/claude/hooks/hooks.json` injects
   `CTX_DIR="$CLAUDE_PROJECT_DIR/.context"` inline on every hook
   line. When `CLAUDE_PROJECT_DIR` is unset/empty the expansion
   silently becomes `CTX_DIR=/.context` (an absolute path at the
   filesystem root). Loud failure was the intent; silent escape was
   the outcome.
3. **Bare `ctx system` calls in `.claude/settings.local.json`** had
   no env plumbing under the strict-declaration model and were
   silently failing on every prompt — the user only noticed because
   one of them references a removed command (`check-backup-age`).

These are the same class of bug the branch set out to kill: silent
inference, with no loud failure when the inference produces nonsense.

## Bet

**One source of truth, one resolver, one anchor name. Per-tool
glue stays at the shell layer where it already belongs.**

Concretely:

- The Go resolver knows about exactly one variable: `CTX_DIR`. No
  `CLAUDE_PROJECT_DIR`, no `CURSOR_PROJECT_ROOT`, no walk-up, no
  filesystem guessing. (`ctx activate` retains its explicit
  candidate-scanning behavior — that's user-invoked discovery,
  not silent resolution; it bails loud when ambiguous. The
  no-walk-up rule applies to the resolver used by operating
  commands, not to the activate scanner.)
- Per-tool hook integrations inject `CTX_DIR` from whichever
  variable that tool exports for the project root, using a hardened
  bash idiom (`${VAR:?msg}`) that fails loud on empty/unset.
- `CTX_DIR` is locked to a `.context` basename in Go. Anything
  else is an explicit error before the value is used. The
  common footgun `export CTX_DIR=$(pwd)` (project root instead
  of `.context` subdir) was the trigger for the
  context-files-leak-to-root bug that started this spec;
  the basename guard catches it on first use rather than on
  50th write.
- The `--context-dir` flag is removed. The single declaration
  channel is `CTX_DIR` (set by `ctx activate` interactively, or by
  per-tool hook injection in subprocess contexts).
- A new sanity hook runs at `UserPromptSubmit` (every prompt)
  and compares `CTX_DIR` to `$CLAUDE_PROJECT_DIR/.context`
  (or the equivalent for the active tool). It emits a verbatim
  warning banner only when it observes drift; in normal
  sessions it is silent. The warning path is idempotent and
  cheap, so the every-prompt cadence is a non-issue. Catches
  stale cross-session bleeds.
- Two project-local hardening hooks
  (`block-dangerous-command`, `block-hack-scripts.sh`) that were
  partial coverage at best are deleted; replaced with native
  Claude Code `permissions.deny` rules, which are the right layer
  for this kind of gate.

### Channels for declaring `CTX_DIR`

The model in three labels:

- **`ctx activate`** — interactive project-local discovery.
- **manual `export CTX_DIR=…`** — advanced / manual explicit binding.
- **hook injection** — tool-local subprocess binding.

The only durable declaration is `CTX_DIR`.

`ctx activate` is not a second declaration channel; it is an
interactive helper that discovers exactly one visible
`.context/` from cwd and emits an `export CTX_DIR=…` statement.
Bails loud on zero or many candidates rather than guessing.

Manual `export CTX_DIR=…` is allowed but receives only
resolver-time validation.

Hook injection is subprocess-local glue. It does not change
the user's shell-level `CTX_DIR`; it anchors hook execution
to the tool-provided project root.

The `check-anchor-drift` hook catches the one residual case:
where the user's shell-level `CTX_DIR` (set by activate or
manual export) diverges from the Claude-injected anchor.

(`--context-dir` and `ctx activate <path>` were considered
and removed; see "Rejected alternatives" below.)

### What we police, what we don't

The "user declares, we read" framing is convenient shorthand
but understates what the basename guard does. The honest
statement of the contract:

- **Basename: policed.** Any declared value whose
  `filepath.Base` is not `.context` is rejected by
  `rc.ContextDir()` at first use. `ctx activate` only emits
  values discovered by scan, which are already `.context`
  candidates by construction (`rc.ScanCandidates` returns
  only `.context`-basename dirs), so no separate activate-time
  guard is needed. Users cannot declare a `.ctx`, `mycontext`,
  or project-root path as their context dir, even if they
  want to.
- **Location: not policed.** Any directory anywhere with
  basename `.context` is accepted. Shared hubs, non-project
  paths, symlinked dirs are all fine. The user owns *where*.
- **Content: not policed at the declaration layer.** `ctx init`
  writes the canonical file set; subsequent operations trust
  that the directory is well-formed. (Per-file validation lives
  in individual commands, not in the resolver.)

This is policing by shape, not by origin — we're not asking
"who declared this" or "is it the same project," we're asking
"does the basename prove the user meant `.context`." Small
enforcement, clearly scoped, loud on violation.

## Rejected alternatives (with why)

- **Walk-up on `.context` directory name.** The original cause of
  the >1-week design rabbit hole. Edge cases: nested legitimate
  projects (e.g. `WORKSPACE/.context` *and*
  `WORKSPACE/ctx/.context` simultaneously), submodules, rogue
  `.context/` dirs adopted silently, sub-agent fragmentation. No
  filesystem shape is sufficient to disambiguate intent.
- **Walk-up on opt-in init marker (`.ctxrc` or `.ctx-root`)
  sentinel.** Better than naked walk-up, but `ctx init` is callable
  by any process in the session including agents. A rogue agent
  running `ctx init` in a subdirectory creates a "valid" project
  root and walk-up adopts it. Marker can't distinguish user intent
  from agent intent.
- **Bake `CLAUDE_PROJECT_DIR` knowledge into the Go resolver.**
  Couples the binary to one tool vendor's env naming. Both
  options (Go vs. shell) require a ctx release to pick up a
  rename — the hooks file is embedded and deployed to
  `~/.claude/plugins/…`, so neither path is release-free. The
  real difference is **locality and cross-tool scale**:
  shell-layer injection keeps tool-specific knowledge in one
  greppable file per tool (future Cursor/Cline/JetBrains/Codex
  integrations get their own `hooks.json`-equivalent,
  self-contained); Go-layer injection would force a
  switch-per-tool resolver that grows with every new
  integration. Shell wins on the refactor arithmetic, not on
  release ergonomics.
- **Keep the inline `CTX_DIR="$CLAUDE_PROJECT_DIR/.context"` as is.**
  Silently produces `CTX_DIR=/.context` when the anchor is empty.
  Hardening to `${CLAUDE_PROJECT_DIR:?...}` is a one-character fix
  that turns silent escape into loud failure. Keep the pattern,
  fix the operator.
- **Keep `--context-dir` with a basename guard.** Once the basename
  guard exists, the flag can only ever produce values that
  `CTX_DIR` could already produce — it's redundant *and* an extra
  silent surface. Removing it tightens the contract to one channel.
- **Keep `ctx activate <path>` (explicit-path mode) as a
  validation-providing helper.** Initially kept on the theory
  that hub consumers would need to activate `.context` dirs
  outside the project tree. Investigation showed `ctx hub start`
  uses `~/.ctx/hub-data/` and never reads `.context/`, so hub
  client/server scenarios activate from the project root like
  everyone else. Remaining use cases (scripts saving one `cd`,
  occasional niche convenience) didn't outweigh the
  one-canonical-path simplification. Removed in section G;
  activate becomes args-free.
- **Keep `block-dangerous-command` and `block-hack-scripts.sh`.**
  Partial coverage that creates a false sense of safety. Doesn't
  catch bare leading `sudo`, `eval`, `sh -c`, shell aliases,
  `/sbin/`, `/opt/`, custom PATH dirs. Native Claude Code
  `permissions.deny` rules cover the same surface with no regex
  maintenance, no Go-vs-bash asymmetry, and no typo'd-and-silently-dead
  references.

## Top three failure modes (and mitigations)

1. **Stale `CTX_DIR` across sessions.** User ran `eval $(ctx activate)`
   in project A yesterday, opens a new shell today in project B,
   runs `claude` without re-activating. CTX_DIR points at A;
   `CLAUDE_PROJECT_DIR/.context` is B. Silent divergence.
   *Mitigation:* the new `UserPromptSubmit` sanity hook compares the
   two and emits a loud verbatim banner when they disagree.
2. **User forgets `ctx activate`.** Direct CLI calls in their shell
   error loud (correct behavior). Hooks still work because they
   inject `CTX_DIR` independently. The only impact is CLI ergonomics,
   not silent corruption.
   *Mitigation:* docs + a one-time prompt from `ctx init` to add
   `eval "$(ctx activate)"` to shellrc, with a non-mandatory hint.
3. **Tool renames its anchor variable.** Anthropic ships a release
   that renames `CLAUDE_PROJECT_DIR` to `CLAUDE_WORKSPACE_DIR`.
   *Mitigation:* one-line edit to
   `internal/assets/claude/hooks/hooks.json` in the ctx
   repository, followed by a ctx release cycle and user
   re-install / `ctx init --force` to redeploy. This is *not*
   release-free — the file is embedded in the ctx binary and
   deployed to `~/.claude/plugins/…/hooks.json` at init time.
   The shell-layer win is **locality and visibility** (one
   file, 14 identical lines, greppable; vs. buried in a Go
   resolver function that future tool integrations would each
   have to branch on), not release-skipping.

## Concrete diff

### A. Go-side resolver

The resolver splits responsibilities cleanly between two
functions:

- **`ContextDir()`** — declaration *shape* validator. Reads env,
  checks the value is set / absolute / canonically-named, returns
  the cleaned absolute path. **No filesystem syscalls.** Used by
  diagnostics that must observe declared state without erroring
  on broken state (e.g. `check-anchor-drift`).
- **`RequireContextDir()`** — operating-command boundary.
  Calls `ContextDir()`, then `os.Stat`s the result and verifies
  it is an existing directory. Renders tailored friendly errors.
  Used by every command path *not* gated by `PersistentPreRunE`
  (MCP handlers, hook operating commands, async work) and as
  the entry point in `PersistentPreRunE` itself.

**Convention rule** the spec locks in: *operating callers use
`RequireContextDir()`; only diagnostic / exempt callers may use
raw `ContextDir()`.* Without this rule, operating callers would
receive shape-valid but non-existent paths and discover the
problem via confusing downstream errors (`open .../TASKS.md: no
such file or directory`) instead of the friendly tailored
error.

`internal/rc/rc.go` `ContextDir()` pseudocode:

```
ContextDir():
  raw := os.Getenv("CTX_DIR")
  if raw == "":                           return ErrDirNotDeclared
  if !filepath.IsAbs(raw):                return ErrRelativeNotAllowed
  abs := filepath.Clean(raw)
  if filepath.Base(abs) != dir.Context:   return ErrNonCanonicalBasename
  return abs
```

`internal/rc/require.go` `RequireContextDir()` pseudocode
(extending the existing helper):

```
RequireContextDir():
  abs, err := ContextDir()
  if err != nil:
    return "", errCtx.NotDeclared(ScanCandidates(cwd))   // existing tailored error
  info, statErr := os.Stat(abs)
  if statErr != nil:
    if os.IsNotExist(statErr): return "", ErrContextDirNotFound(abs)
    return "", ErrContextDirStat(abs, statErr)
  if !info.IsDir():
    return "", ErrContextDirNotADirectory(abs)
  return abs, nil
```

Resolver rejection reasons: unset/empty, relative, non-canonical
basename. Require-layer rejection adds: path does not exist,
stat fails, or path is not a directory. One success path.

Notes:

- **No `filepath.Abs`; absolute-only is a hardline.** `Abs`
  *would* absolutize a relative input via the process cwd —
  exactly the silent cwd-dependency this branch exists to
  eliminate. Use `filepath.Clean` unconditionally to normalize
  separators, dots, and trailing slashes, but require the input
  itself to be absolute. If a user (or agent) exports
  `CTX_DIR=.context`, fail loud with `ErrRelativeNotAllowed`
  rather than silently joining against cwd. Without this rule,
  rejecting non-`.context` basenames is theater: we'd merely
  have moved the cwd-dependency from walk-up resolution into
  env-var interpretation, and `CTX_DIR=.context` would resolve
  differently in every cwd while still passing the basename
  guard. Absolute-only is what makes the rest of the contract
  actually enforce a single deterministic location.
- **No symlink resolution.** `Clean` doesn't follow symlinks.
  A symlink named `.context` pointing at `mycustom` is
  accepted because the user *declared* `.context` — the basename
  guard validates the declared name, not the underlying target
  name. This is the user-friendly read; if a future spec wants
  to lock down resolved targets too, that's a separate
  decision.
- **`OverrideContextDir()` is deleted.** Not "made private,"
  not "test-only" — gone. The `--context-dir` flag was the
  only legitimate caller; with the flag removed there's nothing
  for the function to feed. The `rcOverrideDir` field, the
  `rcMu` mutex protecting it, and the `Reset()` function's
  override-clearing line all go with it.

New typed errors in `internal/err/context/context.go`:

- `ErrRelativeNotAllowed` — value is non-empty and not absolute.
- `ErrNonCanonicalBasename` — value's basename is not `.context`.
- `ErrContextDirNotFound` — value is shape-valid but the
  directory does not exist on disk.
- `ErrContextDirNotADirectory` — value points at a path that
  exists but is a file (or symlink target is a file, etc.).
- `ErrContextDirStat` — `os.Stat` failed for a reason other
  than not-exist (permission denied, I/O error). Carries the
  underlying error.

`ErrDirNotDeclared` already exists; reuse for unset/empty.

`PersistentPreRunE` in `internal/bootstrap/cmd.go`: simplifies.
The current sequence is `RequireContextDir() + validate.Initialized()`.
Under the tightened contract, `RequireContextDir()` covers the
shape/existence/type checks; `validate.Initialized()` continues
to run as the *stronger* semantic check ("the directory looks
like a real ctx project," i.e., contains canonical files) and
returns its own error (`errInit.NotInitialized()`) on miss. The
two-step gate stays — they answer different questions — but
the existence-check responsibility moves inward into the
resolver helper rather than living in two places.

### B. Bootstrap / flag wiring

`internal/bootstrap/cmd.go`:
- Remove the `--context-dir` persistent flag and its
  `OverrideContextDir(contextDir)` call in `PersistentPreRunE`.
- Remove `flagbind.PersistentStringFlag(c, &contextDir, flag.ContextDir, …)`.
- Drop the `var contextDir string` declaration.

`internal/config/flag/flag.go` and
`internal/config/embed/flag/flag.go`:
- Remove `flag.ContextDir` and `embedFlag.DescKeyContextDir`
  constants.

`internal/cli/initialize/cmd/root/run.go`:
- The cwd-fallback at `:98-108` stays; init is the bootstrap
  exemption that creates a context dir before any declaration
  exists. The basename guard does not apply at init time
  because init *creates* the canonical-named directory.

### C. Hook injection hardening

`internal/assets/claude/hooks/hooks.json`:
- All operating hook commands use the hardened `CTX_DIR`
  prefix:
  `CTX_DIR="${CLAUDE_PROJECT_DIR:?CLAUDE_PROJECT_DIR unset; cannot anchor ctx}/.context"`.
  This replaces today's
  `CTX_DIR="$CLAUDE_PROJECT_DIR/.context"` (one operator
  change per line, no structural change).
- **Exception:** `check-anchor-drift` (added in section F) uses
  the same hardened prefix but **prepends**
  `CTX_DIR_INHERITED="${CTX_DIR:-}"` so the diagnostic can
  observe the parent shell's pre-injection `CTX_DIR`. Do *not*
  do a blind sed replacement across the file — the drift hook's
  prepend must survive. See section C-bis for the concrete
  shape and section F for the rationale.
- Shell semantics for the hardened prefix: `${VAR:?msg}` exits
  non-zero with `msg` to stderr when `VAR` is unset or empty,
  so any failure to populate `CLAUDE_PROJECT_DIR` produces a
  loud failure rather than a silent escape to `/.context`.

### C-bis. Concrete shape: the new `hooks.json`

Before / after at the representative line level:

```diff
- "command": "CTX_DIR=\"$CLAUDE_PROJECT_DIR/.context\" ctx system context-load-gate"
+ "command": "CTX_DIR=\"${CLAUDE_PROJECT_DIR:?CLAUDE_PROJECT_DIR unset; cannot anchor ctx}/.context\" ctx system context-load-gate"
```

Structurally the file keeps the same three top-level groups
(`PreToolUse`, `PostToolUse`, `UserPromptSubmit`). The only
changes are:

1. Every `command` string gets the hardened prefix.
2. One new entry added under `UserPromptSubmit` for the
   drift-check hook.

Full shape of the new file (elided command bodies, `…/.context`
stands for the hardened prefix for brevity):

```json
{
  "hooks": {
    "PreToolUse": [
      {"matcher": ".*",           "hooks": [{"type": "command", "command": "…/.context\" ctx system context-load-gate"}]},
      {"matcher": "Bash",         "hooks": [{"type": "command", "command": "…/.context\" ctx system block-non-path-ctx"}]},
      {"matcher": "Bash",         "hooks": [{"type": "command", "command": "…/.context\" ctx system qa-reminder"}]},
      {"matcher": "EnterPlanMode","hooks": [{"type": "command", "command": "…/.context\" ctx system specs-nudge"}]},
      {"matcher": ".*",           "hooks": [{"type": "command", "command": "…/.context\" ctx agent --budget 8000 2>/dev/null || true"}]}
    ],
    "PostToolUse": [
      {"matcher": "Bash",  "hooks": [{"type": "command", "command": "…/.context\" ctx system post-commit"}]},
      {"matcher": "Edit",  "hooks": [{"type": "command", "command": "…/.context\" ctx system check-task-completion"}]},
      {"matcher": "Write", "hooks": [{"type": "command", "command": "…/.context\" ctx system check-task-completion"}]}
    ],
    "UserPromptSubmit": [
      {"hooks": [
        {"type": "command", "command": "CTX_DIR_INHERITED=\"${CTX_DIR:-}\" …/.context\" ctx system check-anchor-drift"},
        {"type": "command", "command": "…/.context\" ctx system check-context-size"},
        {"type": "command", "command": "…/.context\" ctx system check-ceremony"},
        {"type": "command", "command": "…/.context\" ctx system check-persistence"},
        {"type": "command", "command": "…/.context\" ctx system check-journal"},
        {"type": "command", "command": "…/.context\" ctx system check-reminder"},
        {"type": "command", "command": "…/.context\" ctx system check-version"},
        {"type": "command", "command": "…/.context\" ctx system check-resource"},
        {"type": "command", "command": "…/.context\" ctx system check-knowledge"},
        {"type": "command", "command": "…/.context\" ctx system check-map-staleness"},
        {"type": "command", "command": "…/.context\" ctx system check-memory-drift"},
        {"type": "command", "command": "…/.context\" ctx system check-freshness"},
        {"type": "command", "command": "…/.context\" ctx system check-skill-discovery"},
        {"type": "command", "command": "…/.context\" ctx system heartbeat"}
      ]}
    ]
  }
}
```

Where `…/.context\"` expands to the hardened prefix:

```
"CTX_DIR=\"${CLAUDE_PROJECT_DIR:?CLAUDE_PROJECT_DIR unset; cannot anchor ctx}/.context\"
```

Semantics:

- **On a well-configured Claude Code session** (Claude sets
  `CLAUDE_PROJECT_DIR` automatically): every hook invocation
  exports `CTX_DIR=<project_root>/.context` and runs `ctx
  system <name>` against it. Identical to today's behavior when
  everything is working, one character different in the source.
- **On a broken/empty anchor** (imagined regression where
  Claude Code stops setting `CLAUDE_PROJECT_DIR`, or a
  non-Claude-Code runtime mistakenly loads this hooks file):
  the `:?` operator fires, bash exits non-zero with the
  configured message to stderr, `ctx system …` never runs.
  Loud failure visible to the user, instead of silent
  `/.context` escape to the filesystem root.
- **`check-anchor-drift` placement**: first entry under
  `UserPromptSubmit` so any drift banner lands at the top of
  the prompt-submit reminder block (most visible on the first
  prompt of a session, but the hook fires on every prompt).
  Same relay shape as `check-context-size`; no-op when no
  drift. **This one entry prepends
  `CTX_DIR_INHERITED="${CTX_DIR:-}"` before the standard
  injection** — see section F for why; without that capture,
  the standard injection overwrites the inherited `CTX_DIR`
  and the drift comparison becomes tautologically equal,
  leaving the hook a permanent no-op.
- **The rest**: no structural change. No hook added, removed,
  or reordered beyond `check-anchor-drift`. No new matcher
  patterns.

### D. Settings.local.json sweep

`.claude/settings.local.json`:
- Delete the `ctx system block-dangerous-commands` line (note: the
  current name is plural, the binary subcommand is singular —
  this hook is currently dead from the typo alone).
- Delete the `ctx system check-backup-age` line (binary command
  removed; only the dangling reference remains).
- Delete the `bash .claude/hooks/block-hack-scripts.sh` line.
- Audit `permissions.allow` for stragglers: remove absolute-path
  one-off allows from prior refactors (`git -C /Users/.../mv …`),
  remove `Bash(./ctx *)` and `Bash(./dist/ctx *)` (violate the
  PATH-only playbook rule), keep only the canonical `Bash(ctx
  …:*)` patterns.
- Audit `permissions.deny` for completeness; add the replacement
  rules below.

`permissions.deny` replacement set (canonical patterns):

```json
"deny": [
  "Bash(*sudo*)",
  "Bash(git push:*)",
  "Bash(*~/.local/bin*)",
  "Bash(install *)",
  "Bash(* /usr/local/bin*)",
  "Bash(* /usr/bin*)",
  "Bash(*hack/*.sh*)"
]
```

Captures the same intent as the deleted hooks with no Go code,
no bash, no regex maintenance.

### E. Deletions

- `internal/cli/system/cmd/block_dangerous_command/` — entire
  package (`cmd.go`, `run.go`, `doc.go`, tests).
- `internal/config/regex/cmd.go` — entries `MidSudo`, `GitPush`,
  `CpMvToBin`, `InstallToLocalBin`. Keep `GitCommit`, `GitAmend`,
  `TaskRef` (used elsewhere).
- `internal/config/regex/cmd_test.go` — corresponding tests.
- `internal/assets/hooks/messages/block-dangerous-command/`
  and `…-commands/` — embedded message templates (the latter is
  already staged for deletion; verify both gone).
- `.claude/hooks/block-hack-scripts.sh` — bash file.
- All `block-dangerous-command*` references from `internal/assets/claude/hooks/hooks.json`
  if any (verify by grep).
- All `check-backup-age` and `block-hack-scripts` references in
  docs (`docs/cli/system.md`, recipes, etc.).

### F. New: stale-anchor sanity hook

New Go subcommand: `ctx system check-anchor-drift`
(`internal/cli/system/cmd/check_anchor_drift/`).

**The diagnostic problem:** every hook line carries a
`CTX_DIR="${CLAUDE_PROJECT_DIR:?…}/.context"` inline assignment.
That assignment *overwrites* the parent shell's `CTX_DIR` for
the hook subprocess. Operating hooks need this — they must
write to the right `.context/` regardless of what the user's
shell happened to export. But for a hook whose entire job is
to *compare* the inherited `CTX_DIR` against the Claude-injected
anchor, the assignment is a tautology generator: the hook
would always see `CTX_DIR == CLAUDE_PROJECT_DIR/.context`
because the line itself just made it so.

**Solution:** capture the inherited `CTX_DIR` into a sibling
variable *before* the standard injection, so the hook can read
both. The hook line:

```sh
CTX_DIR_INHERITED="${CTX_DIR:-}" \
CTX_DIR="${CLAUDE_PROJECT_DIR:?CLAUDE_PROJECT_DIR unset; cannot anchor ctx}/.context" \
ctx system check-anchor-drift
```

Bash evaluates environment-variable assignments on a command
line left-to-right *before* invoking the command, so
`CTX_DIR_INHERITED` snapshots the parent's value (empty if
unset) before the standard `CTX_DIR` assignment runs. The
`:?` guard still fires if `CLAUDE_PROJECT_DIR` is unset,
keeping fail-loud behavior consistent with every other hook.

Implementation:

- Reads `CTX_DIR_INHERITED` (the parent-shell value) and
  `CTX_DIR` (the post-injection canonical value) directly via
  `os.Getenv`, not through `rc.ContextDir()`. This is a
  diagnostic, not an operating command — it must accept any
  observed value (including unset, including non-canonical)
  so it can describe reality, not impose policy.
- Behavior:
  - `CTX_DIR_INHERITED` empty → silent (user hasn't activated;
    no shell-level declaration exists; hooks still work via
    standard injection on every other hook line).
  - `CTX_DIR_INHERITED` non-empty and equal to `CTX_DIR` after
    `filepath.Clean` on both → silent (correctly anchored).
  - `CTX_DIR_INHERITED` non-empty and unequal to `CTX_DIR` →
    emit a verbatim warning banner via the existing TUI relay
    (same shape as `check-context-size`). Banner names both
    values so the user can see which project's `.context`
    their CLI / `!`-pragma is writing to vs. which project
    Claude Code is in.

Registered as the **first** entry in `UserPromptSubmit` in
`internal/assets/claude/hooks/hooks.json`. Fires on every
prompt; the banner (when emitted) lands at the top of the
prompt-submit reminder block — most visible on the first
prompt of a session, but never gated to it.

`hooks.json` registration shape (note the prepended capture
variable — the difference from every other hook is load-bearing):

```json
"UserPromptSubmit": [
  {"hooks": [
    {"type": "command",
     "command": "CTX_DIR_INHERITED=\"${CTX_DIR:-}\" CTX_DIR=\"${CLAUDE_PROJECT_DIR:?CLAUDE_PROJECT_DIR unset; cannot anchor ctx}/.context\" ctx system check-anchor-drift"},
    {"type": "command",
     "command": "CTX_DIR=\"${CLAUDE_PROJECT_DIR:?CLAUDE_PROJECT_DIR unset; cannot anchor ctx}/.context\" ctx system check-context-size"},
    …rest of the check-* hooks with the standard prefix…
  ]}
]
```

Why this doesn't introduce a new failure mode:
`check-anchor-drift` doesn't *do* anything that requires
`CTX_DIR_INHERITED` to be valid. It reads two env vars and
compares them. It never opens the context directory, never
writes state, never resolves config. `CTX_DIR_INHERITED` lives
only in this hook's command-line scope; it isn't exported
anywhere downstream and no other ctx code reads it.

### G. `ctx activate` simplification + stale-replacement comment

**Drop the explicit-path mode entirely.** `ctx activate
/abs/path/to/.context` is removed. The only invocation form is
`ctx activate` (no args, cwd-scan). The hub-user use case that
nominally justified the path arg turns out not to need it (see
section H — hub server stores at `~/.ctx/hub-data/`,
independent of `.context/`); the remaining justifications
(scripts, niche convenience) don't outweigh the simplification.

`internal/cli/activate/cmd/root/cmd.go`:
- Set `Args: cobra.NoArgs` on the cobra command. Any argument
  produces cobra's standard "accepts 0 arg(s), received N"
  error.
- Remove the `[path]` from the `Use` line; update Examples to
  show only the no-arg form.

`internal/cli/activate/core/resolve/`:
- Delete `internal.go`'s `explicit()` and `hasCanonicalFile()`
  functions; delete the corresponding test file.
- Simplify `Selected()` in `resolve.go` to call `scan()`
  directly; remove the `if len(args) == 1` branch and the
  `args` parameter (caller no longer passes args).
- `scan()` and `rc.ScanCandidates` are unchanged — they are
  the canonical resolution path.

`internal/cli/activate/cmd/root/run.go`:
- Update the call site to drop the `args` parameter.
- When emitting `export CTX_DIR=...`, if the parent shell
  already has `CTX_DIR` set to a different value, prepend a
  comment line to the output: `# ctx: replacing stale CTX_DIR=<old>`.
- This surfaces the change in `eval` output where the user
  can see it before the replacement takes effect.

`internal/err/activate/`:
- Delete `errActivate.NotContext`, `errActivate.NotDirectory`,
  and `errActivate.InvalidPath` if they are referenced only by
  the deleted `explicit()` code (verify by grep before
  deleting). `errActivate.NoCandidates` and `errActivate.Ambiguous`
  remain — both are scan-mode errors.

Activate's surface collapses to: *one command, zero arguments,
one resolution path, one emitted line.* Easier to document,
test, and reason about.

### H. Hub commands: exempt from context-init gate

`ctx hub start`, `ctx hub peer`, `ctx hub status`,
`ctx hub stop`, and `ctx hub stepdown` do not read or write
`.context/`. The hub server stores its state at
`~/.ctx/hub-data/` (`internal/cli/hub/core/server/setup.go`'s
`defaultDataDir()`), independent of any project context dir.
Under the current explicit-context-dir model, however, the
global `PersistentPreRunE` in
`internal/bootstrap/cmd.go:99-108` calls
`rc.RequireContextDir()` for any non-`AnnotationSkipInit`
command — which means hub commands currently fail with
`ErrDirNotDeclared` when `CTX_DIR` is unset, even though they
never need it.

This is a no-broken-windows fix that ships with this spec.
Hub-using teams (specifically the AWS/EKS group ramping on
`ctx hub`) hit this on first contact otherwise.

For each of the five hub command files —
`internal/cli/hub/cmd/start/cmd.go`,
`internal/cli/hub/cmd/peer/cmd.go`,
`internal/cli/hub/cmd/status/cmd.go`,
`internal/cli/hub/cmd/stop/cmd.go`,
`internal/cli/hub/cmd/stepdown/cmd.go` — add the annotation:

```go
Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
```

(Same pattern as `internal/cli/activate/cmd/root/cmd.go:39`,
`internal/cli/initialize/cmd/root/cmd.go:52`,
`internal/cli/doctor/cmd/root/cmd.go:32`, etc.)

Verify the parent `ctx hub` grouping command also doesn't
require it (it has no `RunE` / `Run` of its own — Cobra short-
circuits in PreRunE per `bootstrap/cmd.go:87-89` — so the
parent itself is fine). Children inherit the annotation
individually.

### H. MCP server verification

`internal/mcp/server/server.go` and dispatch paths:
- Verify the MCP server's `ContextDir` resolution goes through
  `rc.ContextDir()` (it currently does — `internal/mcp/handler/*.go`
  uses the `d.ContextDir` field set by `internal/mcp/server/server.go:40`).
- Confirm no parallel `os.Getenv("CTX_DIR")` reads bypass the
  resolver. If found, route through `rc.ContextDir()`.
- Add a regression test asserting MCP server fails closed
  when `CTX_DIR` is unset.

## Test plan

### New tests

- `internal/rc/rc_test.go`:
  - `TestContextDir_RejectsUnset` — `CTX_DIR` unset returns
    `ErrDirNotDeclared`.
  - `TestContextDir_RejectsEmpty` — `CTX_DIR=""` returns
    `ErrDirNotDeclared` (treated as unset).
  - `TestContextDir_RejectsRelative_DotContext` —
    `CTX_DIR=.context` returns `ErrRelativeNotAllowed`. Critical
    regression guard against silent cwd-dependency.
  - `TestContextDir_RejectsRelative_DotSlashContext` —
    `CTX_DIR=./.context` returns `ErrRelativeNotAllowed`.
  - `TestContextDir_RejectsRelative_DotDot` —
    `CTX_DIR=../foo/.context` returns
    `ErrRelativeNotAllowed`.
  - `TestContextDir_RejectsNonCanonicalBasename` —
    `CTX_DIR=/tmp/notdotcontext` returns
    `ErrNonCanonicalBasename` with the offending basename in
    the error message.
  - `TestContextDir_RejectsRoot` — `CTX_DIR=/` returns
    `ErrNonCanonicalBasename` (basename is `/`, not
    `.context`).
  - `TestContextDir_AcceptsCanonical` — `CTX_DIR=/tmp/.context`
    returns `/tmp/.context`.
  - `TestContextDir_NormalizesTrailingSlash` —
    `CTX_DIR=/tmp/.context/` returns `/tmp/.context` (Clean
    strips trailing slash; basename guard passes).
  - `TestContextDir_NormalizesDotSegments` —
    `CTX_DIR=/tmp/./.context` returns `/tmp/.context`.
  - `TestContextDir_AcceptsSymlinkNamedDotContext` — value
    points at a symlink whose basename is `.context` even when
    the link target's basename differs; verify guard checks
    declared name, not resolved target.

- `internal/rc/require_test.go`:
  - `TestRequireContextDir_PathDoesNotExist` —
    `CTX_DIR=/nonexistent/.context` returns
    `ErrContextDirNotFound` after passing `ContextDir`'s
    shape checks.
  - `TestRequireContextDir_PathIsAFile` —
    `CTX_DIR=/path/to/.context` where `.context` is a regular
    file (not a directory) returns `ErrContextDirNotADirectory`.
  - `TestRequireContextDir_StatPermissionDenied` —
    `CTX_DIR=/root-only/.context` where the parent is
    unreadable returns `ErrContextDirStat` wrapping the
    underlying error.
  - `TestRequireContextDir_HappyPath` — existing valid
    directory returns the absolute path, nil error.
  - `TestRequireContextDir_DelegatesShapeChecks` — covers all
    three `ContextDir` shape errors flowing through unchanged
    (unset, relative, non-canonical basename). Spot-checks
    that the wrapper doesn't accidentally swallow upstream
    errors.

- `internal/cli/system/cmd/check_anchor_drift/run_test.go`:
  - `TestCheckAnchorDrift_Match` — `CTX_DIR_INHERITED` and
    `CTX_DIR` equal after filepath.Clean: silent.
  - `TestCheckAnchorDrift_Mismatch` — `CTX_DIR_INHERITED=/project-a/.context`,
    `CTX_DIR=/project-b/.context`: emits expected banner naming
    both values.
  - `TestCheckAnchorDrift_InheritedEmpty` — `CTX_DIR_INHERITED`
    unset/empty (user never activated): silent regardless of
    `CTX_DIR`.
  - `TestCheckAnchorDrift_AcceptsNonCanonicalInherited` — when
    `CTX_DIR_INHERITED=/some/random/path` (non-`.context`
    basename) and `CTX_DIR=/project-a/.context`: the hook still
    runs (does not error via the basename guard) and emits the
    drift banner. Verifies the hook bypasses `rc.ContextDir()`
    so diagnostics work even when the env is broken.

- Shell-level (`hack/test-hook-injection.sh` or equivalent under
  `internal/assets/`):
  - `TestHookInjection_EmptyClaudeProjectDir` — invoking the hook
    pattern with `CLAUDE_PROJECT_DIR=""` exits non-zero with the
    `:?` message on stderr.

- `internal/cli/activate/cmd/root/run_test.go`:
  - `TestActivate_RejectsArgs` —
    `ctx activate /any/path` returns cobra's standard
    "accepts 0 arg(s)" error; no scan, no emit.
  - `TestActivate_StaleReplacementComment` — when env has
    `CTX_DIR=/old/.context` and activate resolves
    `/new/.context`, the output contains the
    `# ctx: replacing stale CTX_DIR=/old/.context` comment
    line.

- `internal/cli/hub/cmd/*/cmd_test.go` (start, peer, status,
  stop, stepdown):
  - `TestHub<Cmd>_AnnotationSkipInit` — verify the cobra
    command carries `cli.AnnotationSkipInit: cli.AnnotationTrue`.
    Structural; cheap; one per command (5 total).
  - **Required** integration-style smoke for one hub command
    (e.g. `ctx hub status` since it reads-only and avoids
    daemonization complexity): construct the root command
    tree as `bootstrap.RootCmd()` does in production, set
    `os.Setenv("CTX_DIR", "")` (or `Unsetenv`), execute the
    hub command, and assert the returned error is **not**
    `errCtx.ErrDirNotDeclared` and PreRunE did not short-
    circuit. Verifies the actual gate-bypass behavior, not
    just the annotation. Without this, a future refactor
    that breaks PreRunE's annotation-handling could leave
    the annotation in place while regressing the bypass.

### Regression / removal verification

- `internal/assets/embed_test.go`:
  - Asset directory list no longer contains
    `block-dangerous-command*` message dirs.
  - Asset directory list still contains `ctx-plan` skill (already
    confirmed by the embed test passing on registration).

- `internal/bootstrap/cmd_test.go`:
  - `TestRoot_NoContextDirFlag` — `ctx --context-dir=foo status`
    fails with cobra's "unknown flag" error.

- Build / `go vet` / `go test ./...` after deletions: no orphan
  imports, no unused symbols.

## Migration & breaking changes

This is alpha; users accept breakage as part of the contract.
Communication, not compatibility shims, is the deliverable.

- **`--context-dir` flag removal.** Anything using the flag in
  scripts, CI, aliases breaks with cobra's "unknown flag" error.
  Action: PR description explicitly lists this; release notes
  carry it forward. No deprecation cycle (alpha).
- **Basename pin to `.context`.** Anything relying on a custom
  context dir name breaks with the basename-guard error. Same
  treatment.
- **`ctx activate <path>` removal.** Scripts that pass an
  explicit path get cobra's "accepts 0 arg(s)" error. Action:
  use `cd <project>; eval "$(ctx activate)"` instead, or
  fall back to `export CTX_DIR=<absolute>` directly (skips
  activate-time validation but the resolver guard still
  fires). Release notes call this out next to the flag removal.
- **Hook subprocess hardening.** Existing user-installed
  `hooks.json` files (deployed by older `ctx init`) won't have the
  `:?` operator. They'll continue working as long as
  `CLAUDE_PROJECT_DIR` is set, which it is in current Claude
  Code. The hardening only matters when the variable is missing —
  and that's the case the new operator is meant to surface.
  Action: `ctx init --force` re-deploys the new hook config.
- **Removed hooks.** `block-dangerous-command` and
  `block-hack-scripts.sh` no longer fire. Action: spec ships
  with `permissions.deny` replacement rules; PR description
  highlights the policy continuation; users porting their own
  `settings.local.json` should adopt the new rules.

PR description template:

```
BREAKING CHANGES (alpha):
- `--context-dir` flag removed. Use `CTX_DIR` env (set by
  `ctx activate`) or per-tool hook injection.
- `CTX_DIR` must name an absolute `.context` directory;
  relative paths and non-canonical basenames are rejected.
- `block-dangerous-command` and `block-hack-scripts.sh` hooks
  removed. Replaced with native `permissions.deny` rules in
  `settings.local.json` (template provided).

MIGRATION:
- Run `eval "$(ctx activate)"` once per shell, or add to shellrc.
- Run `ctx init --force` to re-deploy hardened hook config.
- Adopt the new `permissions.deny` block in your
  `.claude/settings.local.json`.
```

## Documentation impact

Sweep target:

```bash
grep -rln 'CTX_DIR\|--context-dir\|CLAUDE_PROJECT_DIR\|context-dir\|context dir' \
  internal/assets/ docs/ README.md CLAUDE.md AGENT_PLAYBOOK*.md
```

Files expected to need edits:

- `README.md` — overview of context-dir resolution; activate-as-contract.
- `CLAUDE.md` and `internal/assets/claude/CLAUDE.md` — agent guidance
  on session start; remove `--context-dir` mentions.
- `internal/assets/context/AGENT_PLAYBOOK.md` and `…_GATE.md` — agent
  playbook content; align with the new contract.
- `docs/recipes/activating-context.md` — primary activation recipe;
  update examples and rationale.
- `docs/cli/index.md`, `docs/cli/system.md`, `docs/cli/init-status.md`,
  `docs/cli/bootstrap.md` — CLI reference, drop flag.
- `docs/recipes/troubleshooting.md` — env-asymmetry table (hooks vs
  `!`-pragma vs interactive) belongs here so users know about the
  two-tab pattern.
- `docs/recipes/external-context.md`, `docs/recipes/session-lifecycle.md`,
  `docs/recipes/customizing-hook-messages.md`, `docs/recipes/hook-output-patterns.md`,
  `docs/recipes/hook-sequence-diagrams.md` — incidental references.
- `docs/operations/runbooks/architecture-exploration.md` — runbook
  references.
- `internal/assets/claude/skills/*/SKILL.md` (40+ files) —
  full grep pass; most won't match, but any that reference
  `--context-dir`, walk-up, or CTX_DIR semantics need alignment.
- `internal/assets/getting-started.md` (or equivalent quick-start) —
  initial activate guidance.

The skill registered earlier this session (`ctx-plan`) gets two
documentation hooks:

- A short entry in any "skills index" doc that lists
  user-invocable skills.
- A recipe under `docs/recipes/` (suggested name:
  `docs/recipes/scrutinizing-a-plan.md`) explaining when to use
  `/ctx-plan`, the stop conditions, and a worked example. The
  recipe is the on-ramp; the SKILL.md is the working contract.

## Spec coordination & supersession

Header preamble to add to `specs/rc-contextdir-upward-walk.md`:

```markdown
> **Status: SUPERSEDED by `specs/single-source-context-anchor.md`
> (2026-04-24).** The upward-walk approach was rolled back in
> `specs/explicit-context-dir.md`; this spec further refines the
> resolution model. Retained as historical record.
```

Header preamble added to `specs/explicit-context-dir.md`
promotes it from AMENDED → SUPERSEDED, because the body's
prose still treats `--context-dir` as live throughout
(9 references). A SUPERSEDED header with an explicit delta
list keeps the historical record intact while preventing
readers from mistaking the body as current implementation
reference. The delta list names the flag removal, basename
guard, shell-layer injection split, hook hardening, hook
deletions, and the new drift-check hook — each a place where
the old spec's body is actively misleading.

`specs/hook-guard-uninitialized.md`: no preamble; the
`Initialized()` guard work is complementary and unaffected.

`specs/deprecate-ctx-backup.md`: no preamble; backup deprecation
is shipped (binary removed, message templates staged for
deletion). The `check-backup-age` references this spec deletes
are dangling housekeeping, not active spec work.

## Pre-commit smoke test (human-as-pushbutton)

Required before final commit. Each step is a checkbox the user
walks through; agent observes / verifies.

```
□ User: cd into a fresh empty directory (e.g. /tmp/ctx-smoke)
□ User: run `ctx init`; agent verifies .context/ created at
        $PWD/.context with all expected templates
□ User: run `ctx --context-dir=foo status`; agent verifies
        cobra "unknown flag" error
□ User: run `ctx status` with CTX_DIR unset; agent verifies
        ErrDirNotDeclared with friendly hint
□ User: run `eval "$(ctx activate)"`; verify CTX_DIR exported
□ User: run `ctx status`; verify success
□ User: export CTX_DIR=/tmp; run any ctx command; agent
        verifies basename-guard error (basename "tmp" ≠
        ".context")
□ User: open a Claude Code session from the smoke directory
□ User: prompt to verify hooks fire (check-context-size etc.)
□ Agent: verify session-start banner does NOT report drift
        (CTX_DIR matches CLAUDE_PROJECT_DIR/.context)
□ User: in a separate shell, set CTX_DIR=/path/to/some/.context
        and launch a new claude there; verify drift banner
        fires verbatim
□ User: simulate empty CLAUDE_PROJECT_DIR by stripping the
        injection in a copy of hooks.json; verify hooks fail
        loud with the :? message rather than silently
        producing /.context
```

If any step regresses, no commit. Issues filed against this spec.

## Branch strategy

`feat/explicit-context-dir` is retained as the working branch.
50+ files already modified toward the original bet; the parts
that survive (resolver tightening, bootstrap audit, KeyPath
cleanup, init exemption) stay in. The parts that don't
(`--context-dir` flag wiring, `OverrideContextDir`) get reverted
within the same branch as part of this spec's implementation.
Final commit lands as one bulk PR. Long but cohesive — the test
plan above is the safety net that justifies the bulk shape.

## Out of scope

Deferred to future specs:

- `cmd/ctxctl` — maintainer tooling. TASKS.md Phase BT covers
  the full design. Not a prerequisite for this spec; explicitly
  decoupled.
- `block-dangerous-command` regex completeness audit (bare
  `sudo`, `/sbin`, `/opt`, custom PATH dirs). Moot under this
  spec since the hook is being deleted; if a future spec
  reintroduces a deny mechanism, that audit would belong there.
- Cross-tool hook configurations (Cursor, Cline, JetBrains AI,
  Codex). The current spec ships Claude Code integration; other
  tools follow the same per-tool-shell-injection pattern but
  are scoped to their own integration specs.
- Per-project hardening hook policy (whether project-local
  hardening hooks have a documented home pattern). Surfaced by
  this spec's deletion of two such hooks; design discussion
  belongs to a separate "project-local hook conventions" spec.

## Open questions

None remaining. All design decisions settled in the
`/ctx-plan` session that produced this spec
(2026-04-24, session 03823d97).
