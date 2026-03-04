# Multi-Ecosystem Dependency Graphs

**Status**: Implemented

## Problem

`ctx deps` was hardcoded to Go projects. The output pipeline
(`map[string][]string` → Mermaid/table/JSON) was already
language-agnostic, but the input side only knew about `go list`.

## Solution

Interface-based `GraphBuilder` with a registry. Each ecosystem
implements `Detect()`, `Name()`, and `Build(external bool)`.

Detection walks the registry in order; first match wins.
The `--type` flag overrides auto-detection.

## Ecosystems

| Ecosystem | Manifest | Method | External tool required |
|-----------|----------|--------|-----------------------|
| Go | `go.mod` | `go list -json ./...` | Yes (`go`) |
| Node.js | `package.json` | Parse JSON directly | No |
| Python | `requirements.txt` / `pyproject.toml` | Parse text directly | No |
| Rust | `Cargo.toml` | `cargo metadata` | Yes (`cargo`) |

### Detection order

Go → Node.js → Python → Rust (first match wins).

### Node.js details

- **Single package**: internal graph is empty (no internal deps).
  External graph lists all dependencies/devDependencies.
- **Workspaces**: internal graph shows workspace-to-workspace deps.
  External graph shows all deps per workspace.
- Supports both array and object workspace formats.

### Python details

- Parses `requirements.txt` (version specifiers, extras, comments,
  environment markers) or `pyproject.toml` (PEP 621 inline arrays
  and Poetry-style sections).
- `--external` includes dev dependencies from pyproject.toml.
- No import tracing — shows declared dependencies only.

### Rust details

- Requires `cargo` binary. Fails with clear error if missing.
- Internal graph: workspace member cross-dependencies.
- External graph: all dependencies per workspace member.

## CLI surface

```
ctx deps [--format mermaid|table|json] [--external] [--type go|node|python|rust]
```

## Non-goals

- Import-level tracing for Python (would require AST parsing)
- Lock file parsing (package-lock.json, poetry.lock, Cargo.lock)
- Transitive dependency resolution for Node.js/Python
