---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Dependency graph
icon: lucide/git-fork
---

![ctx](../images/ctx-banner.png)

## `ctx dep`

Generate a dependency graph from source code.

Auto-detects the project ecosystem from manifest files and
outputs a dependency graph in Mermaid, table, or JSON format.

```bash
ctx dep [flags]
```

**Supported ecosystems**:

| Ecosystem | Manifest                                | Method                                 |
|-----------|-----------------------------------------|----------------------------------------|
| Go        | `go.mod`                                | `go list -json ./...`                  |
| Node.js   | `package.json`                          | Parse `package.json` (workspace-aware) |
| Python    | `requirements.txt` or `pyproject.toml`  | Parse manifest directly                |
| Rust      | `Cargo.toml`                            | `cargo metadata`                       |

Detection order: Go, Node.js, Python, Rust. First match wins.

**Flags**:

| Flag         | Description                                     | Default       |
|--------------|-------------------------------------------------|---------------|
| `--format`   | Output format: `mermaid`, `table`, `json`       | `mermaid`     |
| `--external` | Include external/third-party dependencies       | `false`       |
| `--type`     | Force ecosystem: `go`, `node`, `python`, `rust` | auto-detect   |

**Examples**:

```bash
# Auto-detect and show internal deps as Mermaid
ctx dep

# Include external dependencies
ctx dep --external

# Force Node.js detection (useful when multiple manifests exist)
ctx dep --type node

# Machine-readable output
ctx dep --format json

# Table format
ctx dep --format table
```

**Ecosystem notes**:

- **Go**: Uses `go list -json ./...`. Requires `go` in PATH.
- **Node.js**: Parses `package.json` directly (no npm/yarn
  needed). For monorepos with workspaces, shows
  workspace-to-workspace deps (internal) or all deps per
  workspace (external).
- **Python**: Parses `requirements.txt` or `pyproject.toml`
  directly (no pip needed). Shows declared dependencies; does
  not trace imports. With `--external`, includes dev
  dependencies from `pyproject.toml`.
- **Rust**: Requires `cargo` in PATH. Uses `cargo metadata`
  for accurate dependency resolution.
