//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package index generates and maintains the **quick-reference
// index tables** at the top of `DECISIONS.md` and
// `LEARNINGS.md`: the markdown tables wrapped in
// `<!-- INDEX:START -->` / `<!-- INDEX:END -->` markers that
// list every entry by ID, date, and title in chronological
// order.
//
// The index is the affordance that lets `ctx agent` send a
// **token-cheap** version of `DECISIONS.md` / `LEARNINGS.md`
// to the AI: instead of injecting the full prose for hundreds
// of entries, it injects only the index table. The agent
// scans the table, decides which entries it needs, and asks
// for those by ID.
//
// # The Index Format
//
// Each index row mirrors one entry block in the source file:
//
//	| ID | Date | Title |
//	|----|------|-------|
//	| L-43 | 2026-04-12 | Lock acquisition order in fanout |
//
// Entry blocks in the source follow a strict shape:
//
//	## [YYYY-MM-DD-HHMMSS] Title text here
//
// [ParseHeaders] extracts the date + title pair from each
// `## [...]` header. [ParseEntryBlocks] returns full block
// metadata (start/end line, ID, date, title, body) so
// callers can grep, render, or rewrite individual entries.
//
// # Updating in Place
//
// [GenerateTable] turns a parsed entry list into the full
// markdown index (table header + rows).
// [Update](path, newTable) finds the marker pair in the
// existing file and replaces only the content between them,
// leaving the rest of the file untouched. If the markers are
// missing, [Update] inserts them under the H1 heading so the
// next run becomes idempotent. [UpdateDecisions] and the
// matching [UpdateLearnings] are convenience wrappers that
// know the canonical file paths.
//
// # Supersession
//
// An entry can be marked **superseded** by a later one
// (a body line starting with `**Status**: Superseded by
// L-99`). The parser tags such entries so renderers can
// gray-out / sort the index accordingly.
//
// # Concurrency
//
// The package is filesystem-IO at the boundary, pure data
// in the middle. Callers serialize updates externally
// (typically by holding the `.context/` directory
// implicitly through process-level execution).
package index
