# Go Toolchain Version Sync

The Go version is written down in eight tracked places. A bump is
a hand-edit of each one, and nothing cross-checked them, so a
hurried bump could leave a stale pin behind and ship it silently.
This spec names the single source of truth, enumerates every
dependent site, and adds a gate that fails when any site disagrees.

## Problem

The Go 1.26 → 1.27 bump (2026-09) touched thirteen files. The
working-tree diff updated the steering source and two of its three
tool-native copies, but `.cursor/rules/tech.mdc` still said
`Go 1.26+`. `make check-steering` would have caught that one, but
only if someone ran `make audit`; CI did not run it. And no check
of any kind tied the remaining sites together:

| Site                              | Kind  | What it pins            |
|-----------------------------------|-------|-------------------------|
| `go.mod`                          | exact | `go 1.27.1` (source)    |
| `go.work`                         | exact | must equal `go.mod`     |
| `tools/ctxctl/go.mod`             | exact | must equal `go.mod`     |
| `.github/workflows/ci.yml`        | floor | two `go-version` pins   |
| `.github/workflows/release.yml`   | floor | two `go-version` pins   |
| `hack/tool-versions.txt`          | floor | `go bin required 1.27`  |
| `.context/steering/tech.md`       | floor | `**Go 1.27+**` prose    |
| `.cursor`, `.clinerules`, `.kiro` | floor | generated from steering |

The only version-consistency check in the repo,
`make check-version-sync`, covers `VERSION` against the plugin
manifests. It knows nothing about the toolchain.

## Source of Truth

The `go` directive in the root `go.mod`. It is the value the Go
toolchain itself enforces, `go mod tidy` maintains it, and every
other site is downstream of it. Its `major.minor` is the *floor*
the other sites must state; its full `major.minor.patch` is what
the workspace and the ctxctl module must repeat exactly.

## Gate

`hack/check-go-version.sh`, exposed as `make check-go-version`:

1. Read the `go` directive from `go.mod`; derive the floor.
2. **Exact:** `go.work` and `tools/ctxctl/go.mod` directives must
   equal the full version.
3. **Floor:** every `go-version:` pin in `.github/workflows/*.yml`,
   the `go bin required <min>` row in `hack/tool-versions.txt`, and
   the first `Go X.Y+` token in `.context/steering/tech.md` must
   equal the floor. Finding zero workflow pins is itself a failure,
   so a renamed key cannot silently blind the check.
4. Print one `FAIL:` line per mismatch naming the site and both
   values; exit with the mismatch count.

Tool-native steering copies are not re-checked here. They are
generated, and `make check-steering` already fails when they
diverge from the source. Checking the source is sufficient.

The script is bash 3.2 and BSD awk/grep clean, per
`specs/hack-script-portability.md`.

## Wiring

- `make audit` runs `check-go-version` right after
  `check-version-sync`.
- The CI `lint` job runs `make check-go-version` **and**
  `make check-steering`. The steering gate existed but was
  local-only; the failure that motivated this spec is exactly the
  kind that reaches a PR from an agent that skipped `make audit`.
- `check-steering` now snapshots the tool-native outputs, regenerates,
  and diffs snapshot against regenerated (the `check-copilot-skills`
  shape). It used to diff against `HEAD`, so it failed on any
  uncommitted steering change even when source and outputs agreed,
  which made `make audit` unpassable in the middle of a bump.

## Bump Procedure

1. Change the `go` directive in `go.mod` (and let `go mod tidy`
   settle `go.sum`).
2. Run `make check-go-version`; fix each `FAIL:` line it prints.
3. Run `make sync-steering` if the steering prose changed.
4. Run `make audit`; commit everything together citing this spec.

## Non-Goals

- Example values in prose, such as the `1.27` row in the
  `specs/check-tools.md` scenario table or the version banner
  sample in a `hack/check-tools.sh` comment. They illustrate
  output shape; a stale example there misleads no tool.
- `golangci-lint` or other action pins in the workflows. Those
  have one site each and nothing to drift against.
- Rewriting the check in Go. The sibling `hack/lint-*.sh` gates
  are shell for the same reason: they are grep-shaped checks over
  text files, and `make audit` already assumes a POSIX shell.

## See Also

- `specs/check-tools.md` — the manifest whose `go` row this gate
  reads.
- `specs/steering-sync-drift-respects-configured-tools.md` — the
  `check-steering` gate this spec promotes into CI.
