//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema is the **CLI-side wrapper** around the
// underlying [internal/journal/schema] validator: it
// resolves which directories to scan based on user
// flags, runs validation across every JSONL session file
// it finds, and persists the resulting drift report
// under `.context/reports/`.
//
// Used by two surfaces:
//
//  1. **`ctx journal schema`** — the standalone
//     drift-check command users run when investigating
//     parser issues.
//  2. **`ctx journal import`** — runs validation
//     **after** an import as a post-flight check;
//     prints a summary if drift is found so users know
//     the next Claude Code release may need a parser
//     update.
//
// # Public Surface
//
//   - **[Run](opts)** — orchestration: resolve scan
//     paths from flags, dispatch validation per file
//     via [internal/journal/schema], aggregate the
//     [Report], optionally write it to
//     `.context/reports/schema-drift-<ts>.md`.
//
// # The Drift Report
//
// Drift is **informational, not fatal** — a session
// with unknown fields still imports cleanly. The
// report exists so maintainers can update
// [internal/journal/parser]'s schema declarations
// when a new Claude Code release adds fields. See
// the [internal/journal/schema] doc.go for the
// upstream semantics.
//
// # Concurrency
//
// Sequential. The validation itself is fast (a few
// milliseconds per JSONL file).
//
// # Related Packages
//
//   - [internal/journal/schema]            — the
//     validator engine.
//   - [internal/cli/journal/cmd/schema]    — the
//     standalone `ctx journal schema` CLI surface.
//   - [internal/cli/journal/cmd/importer]  — runs
//     this package after every import.
package schema
