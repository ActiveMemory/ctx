//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package question appends rows to the outstanding-questions
// artifact at `.context/kb/outstanding-questions.md` for the
// ctx knowledge-base editorial pipeline (Phase KB).
//
// Each row is an open gap the kb has not yet resolved. An open
// question blocks promotion of any claim that depends on it:
// a topic page with a `Q-###` in its `## Open questions`
// section cannot ship at `confidence: high`.
//
// IDs are zero-padded `Q-###` allocated monotonically from the
// existing file's high-water mark. The writer scans the
// artifact for the highest existing `Q-NNN`, increments, and
// formats; when the file does not exist the first row gets
// `Q-001`.
//
// The writer is append-only; when the artifact does not yet
// exist, [Append] initialises it with the schema's table
// header. The schema lives at
// `internal/assets/kb/templates/ingest/schemas/outstanding-questions.md`.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     [cfgKB.OutstandingQuestions].
package question
