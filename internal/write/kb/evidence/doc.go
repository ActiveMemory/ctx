//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package evidence appends EV-### rows to
// `.context/kb/evidence-index.md` for the ctx knowledge-base
// editorial pipeline (Phase KB).
//
// The package enforces two load-bearing invariants from the
// spec and from `.context/ingest/KB-RULES.md`:
//
//  1. **No renumber, no delete.** Once an EV-### is minted, its
//     ID is stable. Callers may demote an existing row's
//     confidence band in-place but the ID and claim text stay.
//     Renumbering breaks every existing citation on every
//     topic page.
//
//  2. **Sequential ID allocation.** [Append] reads the
//     existing file, finds the highest EV-### number, and
//     increments. IDs are zero-padded to three digits
//     (`EV-012`, not `EV-12`).
//
// The `evidence-only` tag (an additive entry in the Tags
// column) signals a row was minted in evidence-only mode and
// must be re-read against its source before a topic-page pass
// promotes it into prose. See
// `internal/assets/kb/templates/ingest/schemas/evidence-index.md`
// for the schema template.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     the EV prefix + digit-width constants and the
//     `evidence-only` tag literal.
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/path]
//     resolves the evidence-index file path.
package evidence
