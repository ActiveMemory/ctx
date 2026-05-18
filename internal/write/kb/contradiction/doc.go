//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package contradiction appends rows to the contradictions
// artifact at `.context/kb/contradictions.md` for the ctx
// knowledge-base editorial pipeline (Phase KB).
//
// Each row records two or more `EV-###` rows that disagree,
// the demotion applied under the demotion policy, and a
// resolution status. Resolved contradictions stay in the file
// as audit trail.
//
// IDs are zero-padded `C-###` allocated monotonically from the
// existing file's high-water mark. The writer scans the
// artifact for the highest existing `C-NNN`, increments, and
// formats; when the file does not exist the first row gets
// `C-001`.
//
// The writer is append-only; when the artifact does not yet
// exist, [Append] initialises it with the schema's table
// header. The schema lives at
// `internal/assets/kb/templates/ingest/schemas/contradictions.md`.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     [cfgKB.Contradictions].
package contradiction
