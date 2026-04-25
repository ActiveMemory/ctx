# internal/config: Constants Package Structure

## Why 60+ Sub-Packages?

This directory contains ~60 sub-packages, each holding constants,
compiled regexes, type definitions, or text keys for a single
domain. This looks unusual. **It's intentional**.

### The problem it solves

A monolithic `config` package creates a false dependency: importing
`config` to use `config.TokenBudget` also imports every regex
pattern, every MCP constant, every entry type, and every CLI flag
name. In Go, the package is the dependency unit: importing one
symbol imports the whole package. A change to any constant in the
package marks every consumer as stale for recompilation and makes
the blast radius of any change the entire codebase.

### The design principle

**One domain, one package.** Each sub-package groups related
constants that change together and are consumed by the same set of
packages. Consumers import exactly what they need:

```go
import "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
import "github.com/ActiveMemory/ctx/internal/config/entry"
import "github.com/ActiveMemory/ctx/internal/config/regex"
```

Not:

```go
import "github.com/ActiveMemory/ctx/internal/config"  // everything
```

### What it costs

- `go list ./internal/config/...` outputs 60+ lines
- IDE autocomplete shows many `config/` entries
- New contributors need this README to find the right package

### What it buys

- Surgical dependency tracking (change `config/mcp/tool` and only
  MCP packages recompile)
- Zero import cycles (all sub-packages are leaves: zero internal
  dependencies)
- Clear ownership (each file belongs to one domain)
- Safe to modify (changing a constant in `config/agent` cannot
  affect `config/mcp`)

## Package Categories

### Flat single-file (pure `const` blocks)

`agent/`, `architecture/`, `bootstrap/`, `box/`, `ceremony/`,
`cli/`, `content/`, `copilot/`, `crypto/`, `ctx/`, `dir/`,
`entry/`, `env/`, `event/`, `flag/`, `fmt/`, `format/`,
`freshness/`, `fs/`, `git/`, `heartbeat/`, `hook/`, `http/`,
`knowledge/`, `load_gate/`, `loop/`, `marker/`, `msg/`, `nudge/`,
`obsidian/`, `pad/`, `project/`, `reminder/`, `rss/`, `runtime/`,
`session/`, `stats/`, `sync/`, `sysinfo/`, `time/`, `token/`,
`trace/`, `version/`, `vscode/`, `warn/`, `watch/`, `why/`,
`wrap/`, `zensical/`

Each contains a `doc.go` and 1-3 files of `const`/`var` definitions.

### Multi-file thematic

- **`regex/`**: 14 files of compiled `regexp.MustCompile()` objects,
  organized by domain (fence, task, entry, markdown, etc.)
- **`file/`**: extensions, ignore patterns, names, limits
- **`dep/`, `doctor/`**: multi-file domain constants

### Hierarchical (nested sub-packages)

- **`embed/`**: user-facing text, organized in 3 tiers:
  - `embed/cmd/`: command Short/Long descriptions (22 files)
  - `embed/flag/`: flag description keys (~10 files)
  - `embed/text/`: output text DescKey constants (~100 files)

- **`mcp/`**: MCP protocol constants, split into 12 sub-packages:
  `cfg/`, `event/`, `field/`, `governance/`, `method/`, `mime/`,
  `notify/`, `prompt/`, `resource/`, `schema/`, `server/`, `tool/`

- **`memory/`**: memory bridge constants

## How To Find the Right Package

**Adding a new constant?**

1. Is it a file name, extension, or path? → `file/` or `dir/`
2. Is it a regex pattern? → `regex/`
3. Is it a CLI flag name? → `flag/`
4. Is it user-facing text? → `embed/text/` (add DescKey + YAML)
5. Is it an MCP protocol value? → `mcp/<sub>/`
6. Is it a time duration, threshold, or limit? → the domain
   package it belongs to (e.g., `agent/` for agent budgets)
7. None of the above? → create a new sub-package named after
   the domain, with a `doc.go` explaining its purpose.

**Looking for an existing constant?**

```bash
# Search by value
grep -r '"the-value"' internal/config/

# Search by name
grep -r 'ConstantName' internal/config/

# List all packages
go list ./internal/config/...
```

## Rules

- **Zero internal dependencies.** Config sub-packages import
  stdlib only (regexp, time, strings, path). Never import another
  internal package.
- **No logic.** One exception: `entry.FromUserInput()` normalizes
  user input to entry type constants. Everything else is pure
  `const`/`var` declarations.
- **Every package has a `doc.go`.** Documents what the package
  provides and what domain it serves.
- **Audit-enforced.** TestDescKeyYAMLLinkage verifies all 879+
  DescKey constants resolve to non-empty YAML values.

## config/ vs entity/ for Types

String-typed enums (`type IssueType string`) and their const
values live in `config/`: the same place all other string
constants live. The type annotation adds compile-time safety but
does not change where the definition belongs.

**When to promote to `entity/`:** When the type grows behavior:
method receivers, interface participation, or business logic. A
type with `func (t IssueType) Severity() int` has outgrown
`config/` and belongs in `entity/`.

| Stage                        | Home               | Example                                           |
|------------------------------|--------------------|---------------------------------------------------|
| Pure value enum              | `config/<domain>/` | `type IssueType string` with const values         |
| Cross-package value enum     | `config/<domain>/` | Same; `config/` is already importable everywhere |
| Type with methods            | `entity/`          | `func (t IssueType) Severity() int`               |
| Type implementing interfaces | `entity/`          | `var _ fmt.Stringer = IssueType("")`              |

The migration path is natural: start in `config/`, promote to
`entity/` when behavior appears. `TestCrossPackageTypes` catches
the cross-package signal that indicates a type may need promotion.
