---
name: structure
description: Project structure and directory conventions
inclusion: always
priority: 10
---

# Project Structure

## Top-level layout

| Path | What it is |
|------|-----------|
| `cmd/ctx/` | Cobra entry point. One main package; thin. |
| `internal/` | Private Go packages (compiler-enforced no-external-import). |
| `editors/<editor>/` | Separately-published editor integrations (currently `editors/vscode/`). NOT embedded. |
| `tools/<tool>/` | Dev tooling for embedded assets, sitting outside the embed tree (currently `tools/typecheck/opencode/`). |
| `docs/` | Source for the docs site at https://ctx.ist. |
| `site/` | Built output of `docs/` via `make site` (zensical). Committed. |
| `specs/` | Feature specs; every commit gets a `Spec: specs/<name>.md` trailer. |
| `.context/` | This project's own ctx context (CONSTITUTION, TASKS, DECISIONS, LEARNINGS, CONVENTIONS, steering, journal). |
| `hack/` | Project shell scripts (release, lint helpers, detectors). |
| `ideas/` | Drafts and unscoped exploration; not authoritative. |

## Inside `internal/`

- Organized by **domain**, one package per concern. The split is
  read/write/config/err/cli/etc., not "by layer."
- `internal/assets/` is the embed payload root. **Everything
  under it is `//go:embed`-ed into the binary.** Read
  `internal/assets/README.md` before adding files there: the
  layout has a contract (embedded vs. separately-published) that
  is easy to violate.
- `internal/cli/<domain>/` mirrors the Cobra command tree. New
  commands land in their domain package, not as siblings of the
  root.

## Where new files go

- **New Go domain logic** → existing `internal/<domain>/` if it
  exists. `ls internal/` and read the candidate's `doc.go`
  before creating a new package; extending the existing package
  is the default.
- **New embedded asset** → under `internal/assets/<domain>/`,
  with a matching `//go:embed` directive added in
  `internal/assets/embed.go`. Add a presence test in
  `embed_test.go` at minimum.
- **Dev tooling for an embedded asset** (linters, type-checkers,
  package.json/tsconfig.json) → `tools/typecheck/<asset>/` or
  similar sibling. Never inside `internal/assets/` itself; the
  embed contract forbids it.
- **New separately-published deliverable** (e.g. a new editor
  extension) → `editors/<editor>/`, with its own pipeline. Not
  under `internal/`.
- **User-facing documentation** → `docs/`, then `make site`.
  Each tool that warrants a guide gets `docs/home/<tool>.md`.

## Where new files do NOT go

- Not in the repo root unless they are project-wide config
  (`Makefile`, `go.mod`, `zensical.toml`, etc.).
- Not in `internal/assets/` if they are not actually embedded.
  Foreign-language source belongs only when `embed.go` references
  it; tooling about embedded assets belongs in `tools/`.
- Not under `internal/` at all if they are deliverables to an
  external channel (marketplace, npm registry, etc.).
