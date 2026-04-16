//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package learning implements the **`ctx learning`**
// command group for managing `LEARNINGS.md` — currently
// just the `reindex` subcommand that regenerates the
// quick-reference index table at the top of the file.
//
// `LEARNINGS.md` is the project's running record of
// gotchas, "gotcha" notes, and hard-won lessons. The
// quick-reference index lets `ctx agent` inject a
// token-cheap **table of contents** instead of the full
// prose, so the AI can scan available learnings and
// request the ones it needs by ID.
//
// # Subcommands
//
//   - **`ctx learning reindex`** — rebuilds the index
//     table by parsing every entry header in
//     `LEARNINGS.md` and emitting a fresh
//     chronologically-sorted table between the
//     `<!-- INDEX:START -->` / `<!-- INDEX:END -->`
//     markers. Idempotent. See
//     [internal/cli/learning/cmd/reindex] for the
//     implementation.
//
// # Adding Entries
//
// New learnings are added through `ctx add learning`
// (the `add` family lives in [internal/cli/add]); this
// package currently only owns the index-maintenance
// side. The `_ctx-learning-add` skill wraps the add
// flow with a guided prompt.
//
// # Concurrency
//
// Stateless. The CLI command runs once and exits.
package learning
