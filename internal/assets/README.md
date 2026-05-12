![ctx](../../assets/ctx-banner.png)

## `internal/assets/`

The embedded asset tree for the ctx Go binary.

Everything under this directory is compiled into the binary via
`//go:embed` (see `embed.go`) and shipped as raw bytes inside
`ctx`. None of these files execute *inside* ctx itself; they are
written to the user's filesystem at `ctx setup` time, where a
*consumer tool* (Claude Code, OpenCode, Copilot CLI, …) loads and
runs them in its own runtime.

If you are looking for a Go-doc-targeted summary, see `doc.go`.
This README is the longer answer: why the tree looks like it
does, what the contract is, and how to add to it without
breaking the contract.

---

## Why Does the Non-Go Code Lives Under `internal/`?

`internal/` in Go convention means "*private to this module: no
external import*" (*enforced by the Go compiler*). It does **not**
mean "Go source only." What lives here is private *build-time
input*: bytes that the `ctx` build process consumes to produce the
release artifact.

The reason these bytes are TypeScript, Bash, PowerShell, JSON,
YAML, and Markdown (*instead of being fetched at runtime or
distributed as a separate package*) is the single-binary
distribution model:

* `ctx` ships as **one statically-linked Go binary**, no runtime
  dependency tree, no package manager, no network fetch on install.
* Integrations with external tools (*Claude Code plugins,
  OpenCode plugins, Copilot CLI hooks*) require *files those
  tools can load*. Those files have to exist somewhere
  ahead of time.
* `//go:embed` makes the Go binary that "somewhere": at compile
  time the build reads each listed file and stores its bytes in
  the `embed.FS` exported by `package assets`. At install time
  (`ctx setup ...`), ctx reads those bytes back out of itself and
  writes them to the user's filesystem at the location the
  consumer tool expects.

A concrete trace, for the OpenCode plugin:

```
internal/assets/integrations/opencode/plugin/index.ts   (source)
              │
              │  build time: //go:embed in embed.go
              ▼
      ctx binary embeds raw bytes
              │
              │  ctx setup opencode → deployPlugin()
              │  (see internal/cli/setup/core/opencode/plugin.go)
              ▼
   ~/your-project/.opencode/plugins/ctx.ts                (deployed)
              │
              │  OpenCode (Bun runtime) auto-loads
              ▼
              executes inside OpenCode
```

The same shape applies to Copilot CLI scripts, Claude Code skill
markdowns, and every other artifact in this tree: ctx is the
*carrier*, not the *executor*.

---

## The Embed Contract

A file belongs under `internal/assets/` if and only if:

1. It is shipped to users **as bytes**, exactly as committed.
2. A consumer (the ctx binary itself, or an external tool ctx
   installs assets into) needs those bytes available with no
   additional fetch or build step.
3. It is referenced by a `//go:embed` directive in `embed.go`.

If a file is meant to be compiled, generated, fetched, linted,
type-checked, or transformed before reaching a user, it does
**not** belong here. More precisely: only its post-transformation
output does. The directory is a *payload manifest*, not a workspace.

### Hard Go Constraint

`//go:embed` paths are relative to the source file containing
the directive, and cannot reference parents (*`../integrations` 
is a compile error*). The practical consequence is that 
the embed root and the assets must be in the same directory 
subtree. Moving assets out of this tree without also moving 
(*or duplicating*) the `embed.go` declaration will break the build.

---

## Directory map

| Path                                         | Language(s)            | Consumer                      | Deployed to                               |
|----------------------------------------------|------------------------|-------------------------------|-------------------------------------------|
| `claude/CLAUDE.md`                           | Markdown               | Claude Code plugin host       | user project root                         |
| `claude/.claude-plugin/plugin.json`          | JSON                   | Claude Code                   | plugin manifest                           |
| `claude/skills/*/SKILL.md`                   | Markdown + frontmatter | Claude Code skills            | skill registry                            |
| `claude/skills/*/references/*.md`            | Markdown               | Claude Code skill body        | referenced from SKILL.md                  |
| `claude/hooks/hooks.json`                    | JSON                   | Claude Code                   | user-level hooks config                   |
| `context/*.md`                               | Markdown templates     | ctx itself (`ctx init`)       | `.context/` in user project               |
| `entry-templates/*.md`                       | Markdown               | ctx (`ctx decision-add` etc.) | new entries appended to `.context/` files |
| `project/*`                                  | Mixed                  | ctx (`ctx init`)              | project-root files (e.g. Makefile.ctx)    |
| `schema/*.json`                              | JSON Schema            | `.ctxrc` validation           | validated in-memory; not deployed         |
| `why/*.md`                                   | Markdown               | ctx (`ctx why …`)             | rendered to stdout; not deployed          |
| `permissions/*.txt`                          | Text                   | ctx permission lookups        | rendered in-process                       |
| `commands/*.yaml`, `commands/text/*.yaml`    | YAML                   | ctx command/flag descriptions | rendered in-process                       |
| `hooks/messages/*/*.txt`                     | Plain text             | ctx hooks                     | rendered to stdout/stderr in-process      |
| `hooks/messages/registry.yaml`               | YAML                   | ctx hook router               | parsed in-process                         |
| `hooks/trace/*.sh`                           | Bash                   | git tracing                   | written to `.git/hooks/`                  |
| `integrations/agents.md`                     | Markdown               | ctx (`ctx setup` flows)       | written to consumer-tool paths            |
| `integrations/copilot/*.md`                  | Markdown               | GitHub Copilot                | repo instructions                         |
| `integrations/copilot-cli/*.{json,md}`       | JSON + Markdown        | Copilot CLI                   | hook config + instructions                |
| `integrations/copilot-cli/scripts/*.sh`      | Bash                   | Copilot CLI (POSIX shells)    | hook scripts                              |
| `integrations/copilot-cli/scripts/*.ps1`     | PowerShell             | Copilot CLI (Windows)         | hook scripts                              |
| `integrations/copilot-cli/skills/*/SKILL.md` | Markdown + frontmatter | Copilot CLI skills            | skill registry                            |
| `integrations/opencode/plugin/index.ts`      | TypeScript             | OpenCode (Bun)                | `.opencode/plugins/ctx.ts`                |
| `integrations/opencode/skills/*/SKILL.md`    | Markdown + frontmatter | OpenCode skills               | skill registry                            |

The `read/` subtree under this directory is **not** an embedded
asset: It is Go code, the typed accessor layer over `FS`. See
`doc.go` for the accessor package overview.

---

## Quality Gates

The current automated coverage (see `embed_test.go` plus the
sibling `read/*/...test.go` files):

* **Presence**: every directory the binary depends on is listed
  by name; missing required files fail the test.
* **Format**: `plugin.json` parses as JSON; `registry.yaml` and
  `.ctxrc` schema parse as YAML/JSON Schema.
* **Schema integrity**: `TestSchemaCoversCtxRC` asserts a
  bidirectional match between `.ctxrc` schema properties and the
  Go struct that consumes them. Drift in either direction fails CI.
* **Spot-content**: targeted substring checks on a handful of
  representative files (e.g. CLAUDE.md contains "Context",
  ctx-history SKILL.md contains "history").
* **Frontmatter shape**: one skill's frontmatter prefix is
  asserted; full validation is not yet generalised.

Anything added to this tree inherits the same exposure: bytes
ship, problems surface at the consumer. Treat new embedded
assets accordingly: add a presence test at minimum, and
prefer a format/parse test where the artifact has any
structure.

---

## Adding a New Embedded Asset

1. **Place the file** under the appropriate subdirectory. If
   the subdirectory does not yet exist, prefer extending an
   existing topic over creating a new top-level folder.
2. **Add an `//go:embed` directive** in `embed.go`. Use the
   most specific glob that captures what you need; avoid
   `**` patterns that may accidentally sweep in new files
   later.
3. **Add a typed accessor** under `read/<domain>/` if callers
   should not need to know the embed path. The package-by-
   domain split keeps callers decoupled from the directory
   layout.
4. **Add a presence test** in `embed_test.go` (or the relevant
   `read/<domain>/..._test.go`). At minimum: assert the file
   reads back non-empty. For structured artifacts (JSON, YAML,
   frontmatter), parse it.
5. **Update the directory map** in this README so the next
   contributor can find your asset without `grep`.
6. **Run `make build && make test`** to confirm the embed
   directive matches an existing file on disk (mismatch is a
   compile error) and the asset is reachable.

---

## What Does **Not** Belong Under `internal/assets/`

* **Go source** that isn't an accessor for `FS`: put it where
  its package belongs.
* **Generated documentation**, transient build artifacts, and
  caches have no business in source control here.
* **Runtime configuration** read from the user's environment
  (the user's `.ctxrc`, secrets, keys). User-owned state lives
  outside the binary.
* **Dev tooling for the embedded assets themselves**
  (`package.json`, `tsconfig.json`, lockfiles, linter
  configs). These are *about* the assets, not part of the
  payload, and would either bloat the embed or pollute the
  contract. Keep them in a sibling tooling directory, with
  tsconfig/lint configs that *reference* this tree via
  relative paths.
* **Anything fetched or generated at install time.** If it
  isn't available at `go build`, it doesn't belong in
  `embed.FS`.
* **Separately-published deliverables.** ctx also ships a
  VS Code extension at `editors/vscode/`. It is *not* embedded
  into the ctx binary: it is built and published independently
  to the VS Code Marketplace under publisher `activememory`
  (see `editors/vscode/README.md`, section "Release"). That
  artifact has its own version, its own toolchain, and its own
  CI gate (`vscode-extension` job in `.github/workflows/ci.yml`).
  Anything that ships via a third-party marketplace, package
  registry, or other out-of-band channel belongs next to *its*
  pipeline, not here.

---

## Embedded vs. Separately-Published: At a Glance

ctx ships two distinct kinds of artifact, and the rules around
them are not the same:

| Dimension | Embedded assets (this tree) | Separately-published (e.g. `editors/vscode/`) |
|---|---|---|
| Carrier | the ctx Go binary | `.vsix` to VS Code Marketplace |
| Build pipeline | `go build` | `npm ci` + `esbuild` + `vsce package` |
| Release pipeline | `release.yml` (`./hack/build-all.sh`) | manual `vsce publish` |
| Version | pinned to the ctx release that compiled them | independent `version` in `package.json` |
| Update reaches user via | ctx binary upgrade | VS Code extension update |
| CI gate today | `typecheck-opencode-plugin` (embedded TS only) | `vscode-extension` (build + production tsc) |
| Lives in repo at | `internal/assets/...` | `editors/<editor-name>/...` |

If you are adding a new harness, decide which model it follows
*before* placing files. Embedded harnesses are simpler to ship
(one binary, no extra publish step) but every byte they carry
becomes part of every ctx release. Separately-published
harnesses have their own release cadence and surface, at the
cost of a second pipeline to maintain.

---

## See Also

* `doc.go`: Go-doc package summary.
* `embed.go`: the single source of truth for what is embedded.
* `embed_test.go`: current presence/format gates.
* `read/`: typed accessors grouped by domain.
* `internal/cli/setup/core/*/`: the `ctx setup` deployers that
  read from `FS` and write to user disk.
