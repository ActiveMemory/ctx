//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package learning implements the **`ctx learning`**
// command group for managing `LEARNINGS.md`, currently
// just the `add` subcommand that appends new entries.
//
// `LEARNINGS.md` is the project's running record of
// gotchas, "gotcha" notes, and hard-won lessons. A
// quick-reference table of contents is projected on demand
// by `ctx index LEARNINGS.md`, so `ctx agent` can scan
// available learnings without reading the full prose —
// computed, never stored in the file.
//
// # Subcommands
//
//   - **`ctx learning add [content]`**: appends a new
//     learning entry with structured context, lesson, and
//     application fields plus required provenance metadata
//     (session-id, branch, commit). Implementation in
//     [internal/cli/learning/cmd/add] delegates to the
//     shared add core.
//
// # Shared Add Core
//
// The cmd/add subcommand is a thin adapter; the
// validation, content extraction, formatting, and
// insertion pipeline lives in [internal/cli/add/core]
// (used by every noun-first add command). The
// `_ctx-learning-add` skill wraps `ctx learning add`
// with a guided prompt.
//
// # Concurrency
//
// Stateless. The CLI command runs once and exits.
package learning
