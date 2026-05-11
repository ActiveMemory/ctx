---
generated-at: 2026-05-11T00:08:47Z
---
# Handover [2026-05-11-000847] Phase SK shipped; 4 polish items on `internal/validate.BodyFlags` for next session

**Provenance:** commit `971bf767` on branch `feat/skill-surface-polish`

## Summary

Phase SK (Skill Surface Polish) — all 7 tasks landed across
6 commits on branch `feat/skill-surface-polish` (off `main` at
`a44edfe3`):

```
971bf767 refactor(validate): consolidate body-flag helpers into internal/validate
ba2faa54 refactor(validate): pure BodyFlags, no PreRunE decoration
1156e44a docs(blog): align thought-piece bylines        (← unrelated, stacked intentionally)
92507039 refactor(validate): single PreRunE enforcement, no panic, no swallowed errors
55acbd81 feat(skills): /ctx-spec --brief, authority boundaries, plan brief offer
f32c8fd9 feat(validate): require body flags on decision/learning add
```

All signed off, `make lint` 0 issues, `go test ./...` 0 failures,
working tree clean. Branch is local-only (not pushed). Spec is at
`specs/skill-surface-polish.md`; design ref at
`ideas/002-editorial-pipeline-and-skill-rigor.md` §3.

The 5-commit churn on the `validate` package (4 commits between
`f32c8fd9` and `971bf767`) reflects iterative correction during
the session — each commit removed a code-smell the previous one
introduced. Functionally correct now, but the API surface still
has 4 specific polish items the user surfaced at session end.
They are non-blocking (the code works, the tests pass) but
should be addressed before PR review.

## Outstanding polish items (next session)

All in `internal/validate/` and its two call sites
(`internal/cli/decision/cmd/add/cmd.go`,
`internal/cli/learning/cmd/add/cmd.go`):

### 1. `RejectPlaceholder` should be unexported

Currently exported as `validate.RejectPlaceholder`. Only call
site is `BodyFlags` in the same package. Rename to
`rejectPlaceholder`. Update test references.

### 2. Per-file convention: private helper lives alone

Convention check confirmed by `internal/sanitize/truncate.go`
(single unexported `truncate` in its own file). After renaming,
move the helper to `internal/validate/rejectplaceholder.go` and
its tests to `internal/validate/rejectplaceholder_test.go`. The
remaining `bodyflags.go` should hold only `BodyFlags`.

### 3. `BodyFlags` takes too much

```go
// Current
func BodyFlags(c *cobra.Command, flags ...string) error
```

Only uses `c.Flags()`. Change to:

```go
func BodyFlags(flags *pflag.FlagSet, names ...string) error
```

Call site:
```go
c.PreRunE = func(cobraCmd *cobra.Command, _ []string) error {
    return validate.BodyFlags(cobraCmd.Flags(), cFlag.Context, ...)
}
```

Update tests — the fixture no longer needs a full
`cobra.Command`; constructing a `pflag.FlagSet` with two flags
and calling `Parse(args)` is simpler.

### 4. `Placeholders` map shape is confusing

`internal/config/validate/placeholder.go` currently defines:

```go
const (
    PlaceholderTBD = "tbd"
    PlaceholderNA  = "n/a"
    ...
)

var Placeholders = map[string]struct{}{
    PlaceholderTBD: {},  // reads as if the key were the identifier
    PlaceholderNA:  {},
    ...
}
```

In a map literal, `PlaceholderTBD:` evaluates to its const value
`"tbd"` — the map key stored is the string, not the identifier.
The user surfaced this as a real code smell: the surface reads
"enum-keyed map" but it's actually a string-keyed set.

**Resolution**: switch to a slice + linear scan (N=9; cost is
negligible):

```go
var placeholders = []string{
    PlaceholderTBD, PlaceholderNA, PlaceholderNAShort,
    PlaceholderNone, PlaceholderSeeChat, PlaceholderSeeAbove,
    PlaceholderSeeBelow, PlaceholderPending, PlaceholderToBeDone,
}
```

```go
// in rejectPlaceholder
key := strings.ToLower(strings.TrimSpace(value))
for _, p := range cfgValidate.Placeholders {
    if key == p {
        return errCli.FlagPlaceholder(flag, value)
    }
}
```

The slice's contents are the same constants, and the slice name
documents the set's purpose without map-key sleight-of-hand. The
audit's magic-strings check passes because every literal lives
in `const PlaceholderXxx`.

After the rename, `Placeholders` (capital P) should be
`placeholders` (private) since only `internal/validate` reads it.

## Gating

- Each item, when fixed, requires the normal gate: `make lint`
  clean, `go test ./...` clean, working tree clean before any
  commit.
- Sign every commit (`git commit -s`).
- Do not add a `Co-Authored-By:` line to commits.
- Stay on `feat/skill-surface-polish`; the branch is the
  active feature branch and stacking polish on top is intended.
- The 4 items are independent; one focused commit per item is
  fine, or fold into a single `refactor(validate): polish surface`
  commit if the diff stays small.

## Session meta — read before resuming

This session accumulated 5 saved feedback memories in
`~/.claude/projects/-Users-volkan-Desktop-WORKSPACE-ctx/memory/`
(2026-05-10), all from the same root cause: I acted on first
impulse instead of grepping the codebase before scaffolding or
editing. Specifically saved:

- `feedback_no_coauthored_by.md` — strip Claude line only, never
  the whole signoff block
- `feedback_branch_before_commit.md` — branch off main first;
  honour "stacking is intentional"
- `feedback_always_signoff.md` — DCO requires `git commit -s`
- `feedback_no_silent_errors_no_panic.md` — propagate errors, no
  `_ =` or `panic` outside `Must`-prefixed functions
- `feedback_no_silent_decoration.md` — helpers do not wrap
  caller's cobra hooks
- `feedback_grep_before_creating_packages.md` — extend existing
  packages by default

The user observed the cumulative drift mid-session and offered
the handover. Resuming agent: **read those memory files first.**
The 4 polish items above are technically small but every fix
this session has introduced a new issue. Slow down: read the
target file fully, grep for the convention, then edit.

## Next session

1. Pull this branch (`feat/skill-surface-polish @ 971bf767`).
2. Read `internal/validate/bodyflags.go` and
   `internal/config/validate/placeholder.go` end-to-end.
3. Read `internal/sanitize/truncate.go` for the
   single-private-helper-per-file pattern.
4. Apply items 1–4 above. One or four commits, either is fine.
5. After lint+test green and tree clean, hand back for push.

Optional follow-up after the 4 fixes land: open the PR. PR body
draft already prepared mid-session; see commit `55acbd81` body
for the user-facing scope summary.

## Open questions

None. The 4 polish items have unambiguous resolutions above.
The user-visible behaviour (decision/learning reject empty and
placeholder body flags) is correct and tested. Polish is purely
about code shape and API surface.
