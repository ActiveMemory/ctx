# go.work.sum Completeness Gate

The repo tracks `go.work` and `go.work.sum`. Nothing kept
`go.work.sum` complete, so it went stale after Go dependency bumps
and dirtied contributors' trees.

## Problem

dependabot updates each module's `go.mod` and `go.sum` but never
`go.work.sum`. Every Go dependency PR it opens (#173, #174, #175 in
2026-09/10) leaves the workspace sum short of checksums that
workspace-mode commands need.

`go build` and `go vet` don't notice. `go mod download`,
`go list -m all`, `go mod verify`, and `go mod tidy` do: they add the
missing lines to `go.work.sum` on the spot. gopls and IDE module
loaders run the same commands. So after a dependabot merge, the next
contributor finds an unexplained `go.work.sum` diff in their tree (66
lines after the grpc 1.84.0 and raft-boltdb 2.4.2 merges) and has to
work out whether it's
theirs to commit.

History shows the cost: repeated standalone "chore: refresh
go.work.sum" commits, each made after someone tripped over the drift.

## Gate

`hack/check-go-work-sum.sh`, exposed as `make check-go-work-sum`:

1. Snapshot `go.work.sum`.
2. Run `go mod download`. Of the commands above, it produces the
   largest set of additions (a superset of `go list -m all` and
   `go mod verify`), and running it again adds nothing.
3. If the file changed, print the missing line count and the first
   few lines, leave the refreshed file in place, and exit 1.

The comparison is against the working copy (snapshot-then-regenerate,
like `check-steering`), so an uncommitted refresh passes. Unlike the
skill-sync checks, a failure doesn't restore the snapshot, because
the regenerated file is exactly what needs committing.

## Wiring

- `make audit` runs it right after `check-go-version`.
- The CI `lint` job runs it right after `make check-go-version`.

Consequence: a dependabot Go bump now fails CI until `go.work.sum`
is refreshed on that PR. That is the intent: the drift shows up on
the PR that causes it, not later in a contributor's checkout. To fix
one, run `make check-go-work-sum` on the PR branch and commit the
result (or supersede the PR with a branch that includes it).

## Non-Goals

- Automating the refresh on dependabot PRs. A workflow that pushes
  to dependabot branches needs write credentials and produces
  commits that need their own DCO sign-off. Revisit if refreshing
  by hand gets tedious.
- `tools/ctxctl/go.sum`. Module-level sums are maintained by
  dependabot and by `go mod tidy` in that module; this gate covers
  only the workspace file.

## See Also

- `specs/go-version-sync.md`: the sibling gate for the toolchain
  version pins, wired at the same points.
