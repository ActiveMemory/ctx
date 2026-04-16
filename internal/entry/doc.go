//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry is the **shared write-side API** for
// adding entries to context files (DECISIONS.md,
// LEARNINGS.md, CONVENTIONS.md). Both the CLI add
// commands and the MCP `ctx_add` tool route through here
// so the validation rules and the on-disk format are
// applied uniformly regardless of caller.
//
// # Public Surface
//
//   - **[Validate](params)** — applies all the rules
//     a new entry must satisfy before it can be
//     written: required fields present, body
//     non-empty, identifier-like fields well-formed,
//     provenance flags satisfy the per-project
//     [internal/rc.ProvenanceConfig]. Returns a typed
//     error per failure for actionable messages.
//   - **[Write](params)** — writes the entry: builds
//     the timestamped header, formats the body with
//     the canonical attribute order, appends to the
//     target file, and updates the file's index
//     table via [internal/index].
//   - **[ValidateAndWrite](params)** — convenience
//     wrapper that runs [Validate] and then [Write]
//     when validation passes; this is the function
//     the CLI commands and the MCP handler actually
//     call.
//
// # Validation Rules
//
// Beyond presence checks, validation enforces:
//
//   - **Title length** — fits the index-table
//     column width without truncation.
//   - **Body has at least one substantive line** —
//     not just whitespace or template placeholders.
//   - **Provenance** — `--session-id`, `--branch`,
//     `--commit` are required when
//     `provenance_required` enables them in
//     `.ctxrc`.
//   - **No secrets** — body is scanned against
//     [internal/config/token.SecretPatterns]; a
//     match aborts with a typed error so the user
//     can scrub before retry.
//
// # On-Disk Format
//
// Entries follow the canonical shape:
//
//	## [YYYY-MM-DD-HHMMSS] Title text here
//
//	Body paragraph(s)...
//
//	**Attribute**: value
//	**Attribute**: value
//
// Attributes are emitted in a fixed order so re-runs
// produce stable diffs.
//
// # Concurrency
//
// Single-process write assumption. Concurrent writes
// to the same file would race on the append; ctx
// CLI is single-process by design.
package entry
