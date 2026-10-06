# CODEOWNERS: Full Ownership for Every Maintainer

The `main` ruleset requires a code-owner approval on every pull
request. The intent is that each of the five maintainers
(@josealekhine, @parlakisik, @CoderMungan, @hamzaerbay, @bilersan)
owns the whole repository, so any one of them can approve. The
CODEOWNERS file did not say that.

## Problem

GitHub applies only the **last** CODEOWNERS pattern that matches a
file. The file had:

- **Five separate `* @user` lines**, one appended per maintainer
  (all on 2026-03-08). Each overrode the one above it, so for most
  of the repo only `* @bilersan`, the last line, counted.
- **Eight path rules naming only @josealekhine** (`/cmd/`,
  `/internal/`, `*.md`, `/docs/`, `/specs/`, `/hack/`, `/.github/`,
  `Makefile`). Being later in the file, they replaced the `*` owners
  for those paths instead of adding to them.

Net effect: @bilersan was sole owner of everything outside those
paths, @josealekhine was sole owner of everything inside them, and
the other three owned nothing. GitHub's CODEOWNERS error check
reported nothing, because the file was syntactically valid.

## Solution

- **One rule:** `* @josealekhine @parlakisik @CoderMungan
  @hamzaerbay @bilersan`. A comment above it explains the
  last-match rule, so the next maintainer added goes onto this line.
- **Path rules removed.** Each one could only narrow ownership,
  which contradicts the policy.
- **Guard:** `internal/compliance/codeowners_test.go` requires
  exactly one non-comment rule, with pattern `*` and at least one
  `@handle`. It doesn't pin the names, so adding a maintainer stays a
  one-line edit, but appending a second line fails `go test`.

All five accounts have write access, which GitHub requires before a
code owner counts.

## Verification

- `TestCodeownersSingleRule` passes on the new file, and fails on
  the old one (13 rules).
- After merge, `gh api repos/ActiveMemory/ctx/codeowners/errors`
  reports no errors.

## Non-Goals

- Per-area ownership. If areas get dedicated owners later, the
  policy changes and so must this guard; every rule must then repeat
  all the owners it intends.
