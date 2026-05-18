//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sourcecoverage maintains the source-coverage ledger
// at `.context/kb/source-coverage.md` for the ctx
// knowledge-base editorial pipeline (Phase KB).
//
// The ledger is a state machine over every source the kb has
// touched. Each pass that touches a source MUST advance its
// row honestly before writing the closeout. "Lying to the
// ledger" (recording `comprehensive` while real coverage is
// partial) is a hard anti-pattern named in the spec and in
// `.context/ingest/KB-RULES.md`.
//
// # State machine
//
//	discovered  → admitted | skipped
//	admitted    → highlights-extracted | partially-ingested
//	            | topic-page-drafted | comprehensive
//	highlights-extracted → partially-ingested
//	            | topic-page-drafted | comprehensive
//	partially-ingested → topic-page-drafted | comprehensive
//	topic-page-drafted → comprehensive
//	comprehensive → (terminal until source updates)
//	superseded  → (terminal)
//	skipped     → (terminal until scope changes)
//
// [Advance] refuses transitions that are not in the allow-list.
// The doctor advisory cross-checks the ledger against file
// existence + last-modified time vs. each row's Updated cell
// to catch ledger drift independently of write-time
// validation.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     the state-name constants (`StateAdmitted`, etc.).
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/path]
//     resolves the ledger file path.
package sourcecoverage
