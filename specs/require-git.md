# Require Git as Architectural Precondition

`ctx` already needs git to work properly; this phase enforces
it. `ctx init` and every `ctx` subcommand (except read-only
diagnostics) refuse to operate without `.git/`.

## Problem

ctx's persistent-memory promise is dishonest without an undo
layer. LLM agents are not trustworthy stewards of files; without
git, a deletion is permanent. With git, `git reflog` recovers it.

ctx today silently degrades when git is absent:

- Provenance flags (`--commit`, `--branch`) carry the
  `commit:none` sentinel. Doctor surfaces this as "make a first
  commit so future captures pin to a real SHA" — essentially
  saying "git is required for this to work properly" without
  enforcing it.
- Every git-touching code path (`internal/gitmeta/*`,
  `internal/cli/wrapup/*`, capture commands, doctor advisories)
  carries a "what if no git" branch. That branch is dead code in
  practice — N≈0 ctx projects have run without git in the
  user's experience — and pays testing + maintenance cost for no
  benefit.
- Lifting the editorial pipeline (see
  `specs/kb-editorial-pipeline.md`) inherits an explicit
  git-required assumption: closeout frontmatter has `sha:` /
  `branch:`; evidence-index SHA-pins in-repo citations; handover
  Provenance reads `git HEAD`. Without git, the editorial
  pipeline silently degrades the same way canonical capture
  does.

We should have done this on day zero.

## Approach

Single guard at root command PersistentPreRunE:
`gitmeta.RequireGitTree(projectRoot)` returns nil when `.git/`
is a directory or a worktree pointer file; returns a typed error
otherwise. Opt-out list contains only read-only / help-shaped
commands (`--help`, `--version`, `ctx system bootstrap`, plus
any audit-confirmed others).

Trust boundary: ctx never takes side actions in your repo, so
**no auto-`git init`**. The user runs `git init` first, then
`ctx init`. The error message does the work.

Constitutional amendment: add a "Git is required" rule to
`.context/CONSTITUTION.md` Process Invariants. DECISIONS.md
entry names the breaking change and rationale. Old `commit:none`
rows in existing `DECISIONS.md` / `LEARNINGS.md` / `TASKS.md`
remain valid as historical record; doctor stops counting them as
new advisories because the scenario is unreachable going
forward.

## Behavior

### Happy Path

1. User: `git init` in fresh project root.
2. User: `ctx init`. Init proceeds; lays down `.context/`.
3. All subsequent `ctx` commands work as today.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `ctx init` in dir without `.git/` | Refuse: `ctx init: .git/ not found at <path>; run \`git init\` first, then re-run.` Non-zero exit. |
| Any `ctx <cmd>` without `.git/` (existing user upgraded) | Refuse with cmd name in error: `ctx <cmd>: .git/ not found at <path>; ctx requires a git repository (run \`git init\` to fix).` Non-zero exit. |
| `ctx --help` / `ctx --version` / `ctx system bootstrap` | Allowed without git (read-only / diagnostic). |
| `.git/` is corrupt or unreadable | Refuse with the underlying git/os error wrapped. |
| Detached HEAD | Allowed. Provenance records `branch=detached` per existing convention. |
| Worktree | Allowed. Worktrees have `.git` as a file pointing at the main checkout; the guard accepts file-or-directory. |
| Submodule | Allowed if the submodule itself has a reachable `.git` (file or dir); refuse otherwise. |
| Bare repo | Refuse. ctx assumes a working tree; bare repos have no `.git/` in the working sense. |

### Validation Rules

`RequireGitTree(projectRoot string) error` returns nil when
`<projectRoot>/.git` exists as either:

- a directory (regular repo), OR
- a regular file (worktree pointer per git convention).

Returns `*MissingGitError` when `.git` is absent; returns a
wrapped error for other stat failures.

Wired into the root command PersistentPreRunE. Opt-out list is
explicit and audited; new commands are git-required by default.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `.git/` missing on `ctx init` | `ctx init: .git/ not found at <path>; run \`git init\` first, then re-run.` | `git init` |
| `.git/` missing on other commands | `ctx <cmd>: .git/ not found at <path>; ctx requires a git repository (run \`git init\` to fix).` | `git init` |
| `.git/` corrupt | Wrapped underlying error verbatim | Repair git per its normal guidance |

## Interface

No new commands. Behavior change is at the root-command
PersistentPreRunE; existing commands inherit the precondition
transparently.

## Implementation

### Files to Create / Modify

| File | Change |
|------|--------|
| `internal/gitmeta/require.go` | New: `RequireGitTree(projectRoot string) error` + `MissingGitError` type. |
| `internal/cli/parent/cmd.go` (or wherever the root cmd lives) | Add PersistentPreRunE that calls `RequireGitTree`; opt-out list for `--help`, `--version`, `ctx system bootstrap` (audit other read-only/diagnostic commands during implementation). |
| `internal/cli/initcmd/init.go` | Replace any "if no git, fall back to commit:none" logic with the precondition guard at top of init. |
| `internal/gitmeta/resolvehead.go` | Remove `commit:none` fallback path; HEAD must resolve. Update tests. |
| `internal/cli/doctor/advisory.go` | Remove `commit:none` advisory and counts (state unreachable). |
| `internal/cli/<various>/cmd.go` | Audit any command that has its own git check or `commit:none` fallback; remove. |
| `.context/CONSTITUTION.md` | Add "Git is required" rule under Process Invariants. |
| `.context/DECISIONS.md` | New entry: "Mandate git as architectural precondition." Status: Accepted. Context = LLM-safety + provenance honesty + dead-code elimination. Consequence = breaking change for any pre-existing git-less ctx project (N≈0 in practice). |
| `docs/recipes/bootstrap-a-project.md` | Update getting-started to show `git init` before `ctx init`. |
| `README.md` | Update install / first-run section to mention git as a precondition. |
| `docs/cli/init.md` | Note the precondition explicitly with the documented error wording. |
| `dist/RELEASE_NOTES.md` (or equivalent) | Tag as breaking change with migration note: "Run `git init` in any pre-existing git-less ctx projects before upgrading." |

### Key Functions

```go
// internal/gitmeta/require.go
type MissingGitError struct {
    ProjectRoot string
    Cmd         string // optional; populated by root PreRunE
}

func (e *MissingGitError) Error() string {
    if e.Cmd != "" {
        return fmt.Sprintf("ctx %s: .git/ not found at %s; ctx requires a git repository (run `git init` to fix)", e.Cmd, e.ProjectRoot)
    }
    return fmt.Sprintf(".git/ not found at %s; ctx requires a git repository (run `git init` to fix)", e.ProjectRoot)
}

func RequireGitTree(projectRoot string) error {
    p := filepath.Join(projectRoot, ".git")
    info, err := os.Stat(p)
    if err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            return &MissingGitError{ProjectRoot: projectRoot}
        }
        return fmt.Errorf("stat .git at %s: %w", p, err)
    }
    _ = info // accept either dir (regular repo) or file (worktree pointer)
    return nil
}
```

### Helpers to Reuse

- Existing `internal/gitmeta/` resolution code — keep the read
  paths; drop the unborn-HEAD `commit:none` fallback throughout.

## Configuration

None. The precondition is unconditional; no opt-out flag, no
env var escape hatch.

If a CI test fixture genuinely needs to run without git, the
test sets up a temp git repo in its `TestMain` (the existing
pattern in many ctx tests becomes the universal one).

## Testing

### Unit

- `internal/gitmeta/require_test.go` — present `.git` dir → nil;
  present `.git` file (worktree pointer) → nil; absent → typed
  `*MissingGitError`; corrupt → wrapped error.
- `internal/cli/parent/cmd_test.go` — root PreRunE refuses
  without `.git/`; `--help` / `--version` / `bootstrap` allowed.
- `internal/gitmeta/resolvehead_test.go` — verify `commit:none`
  literal removed; HEAD must resolve.

### Integration

- `internal/cli/initcmd/init_test.go` — `ctx init` in dir
  without `.git/` returns the documented error and exits
  non-zero; in dir with `.git/` proceeds normally.
- Compliance test (audit-pass): no remaining `commit:none`
  literal in `internal/` (catches future regressions).

## Non-Goals

- **No auto-`git init`.** ctx does not modify your filesystem
  outside `.context/`.
- **No interactive prompt.** Refuse-and-exit; the user takes the
  action.
- **No `--allow-no-git` escape hatch.** If you need ctx without
  git, you need a different tool.
- **No grandfathering of existing git-less projects.** Release
  notes flag this as breaking with a one-command migration; the
  affected user-base is N≈0.
- **No retroactive cleanup of `commit:none` rows in existing
  DECISIONS.md / LEARNINGS.md / TASKS.md.** Old entries remain
  valid historical record; doctor stops surfacing them as new
  advisories because the scenario is now unreachable.
- **No bare-repo support.** ctx assumes a working tree.

## Open Questions

1. **Opt-out list completeness.** Beyond `--help` / `--version`
   / `ctx system bootstrap`, are there other commands (e.g.
   `ctx guide`, `ctx --version`-shaped subcommands) that should
   be allowed without git? Audit during implementation by
   walking the command tree; default to "git required" unless a
   command is purely informational.
2. **Mandate timing.** Ship in the next minor (v0.X) or wait
   for v1.0? Lean next minor, with the breaking-change note
   prominent in release notes; v1.0 inherits the precondition.

---

## Source

`/ctx-plan` adversarial-interview pass, session 01d0cf92,
2026-05-09. The git-mandate question landed at the close of the
kb-editorial-pipeline planning round; treated as a separate
constitutional spec rather than folded in because of cross-
cutting impact.

Cross-references:

- `specs/kb-editorial-pipeline.md` — depends on this spec.
  Closeout / handover writers and evidence-index SHA-pinning all
  assume git; once this spec ships, those writers don't need a
  `commit:none` fallback.
- `ideas/001-sibling-project-undercover-analysis.md` — the
  surveyed sibling ships the `git`-required assumption verbatim
  (its context dir and `.git/` live as siblings at the project
  root).
- `ideas/003-editorial-pipeline-debated-brief.md` — handover
  Provenance auto-recorded from `git HEAD`.
