//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package index builds the **session-ID-to-filename map** that
// every other journal subcommand uses to find a journal entry
// when given only a session ID.
//
// The map matters because users (and skills) routinely refer
// to a session by its ID, a short alphanumeric tag like
// `abc123`, but on disk the journal entry filename is keyed
// by date and slug. The mapping has to be built on demand
// from the entry frontmatter; it cannot be derived from the
// filename alone.
//
// # The Surface
//
//   - **[Session](dir)**: walks the journal directory,
//     reads the YAML frontmatter of every `*.md` entry,
//     extracts each `session_id`, and returns a
//     `map[sessionID]filename`. Entries without a
//     `session_id` field are silently skipped.
//   - **[ExtractSessionID](path)**: reads one file and
//     returns its `session_id` (empty string if not
//     present, error if the file cannot be read or the
//     frontmatter cannot be parsed).
//   - **[LookupSessionFile](dir, sessionID)**: convenience
//     wrapper: calls [Session] and returns the matching
//     filename, or empty string if not found.
//
// # Performance
//
// [Session] reads the frontmatter only, not the full
// body, so the cost scales with `O(N)` files but with
// a small per-file constant. For a journal with a few
// hundred entries, the build typically completes well
// under 100 ms. Callers that need many lookups in a
// row should call [Session] once and cache the map
// rather than calling [LookupSessionFile] repeatedly.
//
// # Concurrency
//
// All functions are stateless. Concurrent callers
// against the same directory each pay the full read
// cost; no module-level cache is implemented because
// the journal directory mutates between sessions and
// stale-cache bugs are worse than the perf cost.
package index
