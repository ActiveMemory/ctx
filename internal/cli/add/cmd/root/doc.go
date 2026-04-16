//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements **`ctx add`** — the command
// that adds a new entry (task / decision / learning /
// convention) to the corresponding `.context/` file with
// validated provenance, canonical formatting, and an
// auto-updated index table.
//
// # Public Surface
//
//   - **[Cmd]** — cobra command with the type
//     selector (`-t task|decision|learning|convention`)
//     plus type-specific flags (`--priority`,
//     `--rationale`, `--consequence`, `--lesson`,
//     `--branch`, `--commit`, `--session-id`,
//     `--from-file`, `--application`, etc.).
//   - **[Run]** — validates the supplied flags
//     against the type's required-fields list,
//     extracts content from positional args or
//     `--from-file`, formats the entry via the
//     `core/format` siblings, and inserts it via
//     [internal/cli/add/core/insert].
//
// # Validation Boundaries
//
// All hard checks (required fields, secret patterns,
// length limits, provenance requirements per
// `.ctxrc`) live in [internal/entry] so the rules
// are identical regardless of caller (CLI here, MCP
// `ctx_add` tool elsewhere).
//
// # Concurrency
//
// Single-process, sequential.
package root
