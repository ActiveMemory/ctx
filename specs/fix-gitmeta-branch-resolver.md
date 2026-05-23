# Fix gitmeta Branch Resolver

`internal/gitmeta.resolveBranchOrDetached` was invoking
`git rev-parse --show-current` instead of
`git branch --show-current`. The handover writer surfaces the
result in frontmatter, so the bug shipped a literal `branch:
--show-current` line into every wrap-up.

## Problem

`--show-current` is a flag of the `git branch` subcommand. It is
**not** recognized by `git rev-parse`. When passed to rev-parse,
git treats it as an unknown revision/object argument and — because
rev-parse's fallback for unknown args is to echo them back as
literal output with exit 0 — the resolver returned the string
`"--show-current"` verbatim. Exit 0 meant the error guard
(`runErr != nil` ⇒ `"detached"`) never tripped.

Symptom (confirmed on `git version 2.50.0`):

```
$ git rev-parse --show-current
--show-current
$ git branch --show-current
main
```

The constant `FlagShowCurrent` was further misclassified in
`internal/config/git/git.go` under a comment group labeled
"Rev-parse flags", which is what made the wrong call site look
right at review time.

The existing test suite (`internal/gitmeta/resolvehead_test.go`)
only exercised paths where git was absent or env overrides took
priority, so the happy-path resolver was never invoked under
test — the regression had no coverage.

## Solution

1. `internal/gitmeta/branch.go` — swap `cfgGit.RevParse` for
   `cfgGit.Branch` in the `execGit.Run` call.
2. `internal/config/git/git.go` — split `FlagShowCurrent` out of
   the "Rev-parse flags" group into a new "Branch subcommand
   flags" group. The misclassification was the latent enabler.
3. `internal/gitmeta/head.go` — corrected the doc comment that
   claimed branch resolution used
   `git rev-parse --abbrev-ref HEAD` (wrong subcommand AND wrong
   flag).
4. `internal/config/git/doc.go` — package-level prose updated for
   the new grouping.
5. `internal/gitmeta/resolvehead_test.go` — added
   `TestResolveHead_RealRepoReturnsBranchName`: initializes a
   real git repo on branch `trunk`, calls `ResolveHead`, and
   asserts the resolved branch (a) equals `"trunk"` and (b)
   contains no `--` flag literal. The dual assertion is
   deliberate — a future regression that returned a different
   wrong flag (`--abbrev-ref`, say) would still be caught by
   (b) even if the test fixture's branch name changed.

## Out of Scope

- Backfilling correct `branch:` values in already-written
  handovers. The bad string is parsed back as a literal by the
  read side; downstream tools that expect `main` will trip, but
  rewriting historical handover frontmatter is a separate cleanup
  task and not blocking.
- Auditing other `cfgGit.*` constants for similar
  misclassifications. The doc.go regrouping fixes the only
  current instance; a broader sweep is unwarranted without
  evidence of further misuse.
