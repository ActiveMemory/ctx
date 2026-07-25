# layout

## File Organization

- **Public API in main file, private helpers in separate logical files**
  - `loader.go` (exports `Load()`) + `process.go` (unexported helpers)
  - NOT: one file with unexported functions stacked at the bottom
- Reasoning: agent loads only the public API file unless
  it needs implementation detail
- **Name files after what they contain, not their role**
  - `format.go`, `sort.go`, `parse.go` — named by responsibility
  - NOT: `util.go`, `utils.go`, `helper.go`, `common.go` — junk drawer names
  - If a file can't be named without a generic label,
    its contents don't belong together
  - Existing junk drawers should be split as their contents grow

## Maintainer-Only Binaries (Layout and Installation)

Maintainer-only binaries — tooling that must never ship to end
users — live in `tools/<name>/` as separate Go modules. The
module path is lexically nested under the main ctx module
(`github.com/ActiveMemory/ctx/tools/<name>`) so the new module
CAN import the parent's `internal/` packages (Go's
internal-import rule is path-lexical, not module-scoped — see
LEARNINGS.md), reusing `rc`, `desc`, `nudge`, `config`
primitives without duplication.

Build and install:

- Built to `dist/<name>` via `make <name>` (keeps the repo
  root clean).
- PATH-installed to `/usr/local/bin/<name>` via
  `make install-<name>` / `make reinstall-<name>` —
  mirroring ctx's `install` / `reinstall` targets so one
  binary serves every worktree and repo copy.
- The shipped `ctx` binary's `go.mod` must NOT `require` the
  maintainer module, giving a **hard module-graph guarantee**
  that the maintainer code can never leak into `ctx`.

Repo-local hooks calling the maintainer binary live in the
gitignored `.claude/settings.local.json`, **not** in the
shipped `internal/assets/claude/hooks/hooks.json`. The hook
command shape is `cd "$CLAUDE_PROJECT_DIR" && <name>
<subcommand>` (PATH binary, project-root cwd so `.context/`
resolves correctly under cwd-anchoring).

`tools/ctxctl/` is the first inhabitant. Future maintainer
binaries follow the same shape.

- Maintainer-only docs (features that run through the ctxctl binary, never
  shipped to users) live under docs/operations/runbooks/, not docs/recipes/
  (user-facing). contributing.md is the entry point that links to the runbook.
