//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package decision appends rows to the kb-scoped
// domain-decisions artifact at
// `.context/kb/domain-decisions.md` for the ctx knowledge-base
// editorial pipeline (Phase KB).
//
// This file is DISTINCT from project-level
// `.context/DECISIONS.md`. The two files have different
// schemas, different write authority, and different lifecycles.
// `domain-decisions.md` records kb-scoped positions on the
// subject matter under study, cites `EV-###` rows (not commit
// SHAs), and is written by `/ctx-kb-ingest`.
//
// IDs are zero-padded `DD-###` allocated monotonically from the
// existing file's high-water mark. (The brief drafted these as
// `D-###`; the canonical schema at
// `internal/assets/kb/templates/ingest/schemas/domain-decisions.md`
// pins the prefix to `DD-` to keep the namespace distinct from
// any future project-side `D-###` series, and this writer
// follows the schema.) When the file does not exist the first
// row gets `DD-001`.
//
// The writer is append-only; when the artifact does not yet
// exist, [Append] initialises it with the schema's table
// header.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     [cfgKB.DomainDecisions].
package decision
