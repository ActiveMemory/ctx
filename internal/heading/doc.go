//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package heading parses Markdown headings out of ctx knowledge files.
//
// It is read-only: the package extracts structure, it never rewrites a
// file. A table of contents is projected on demand by `ctx index <file>`
// (and by `ctx agent`), never stored in the file as an index block.
//
// # Two recognizers
//
// [ParseHeaders] and [ParseEntryBlocks] are entry-specific: they match the
// timestamped shape
//
//	## [YYYY-MM-DD-HHMMSS] Title text here
//
// used by DECISIONS.md / LEARNINGS.md. [ParseHeaders] returns the date +
// title pair per entry; [ParseEntryBlocks] returns full block metadata
// (start/end line, header, body) so callers — notably `ctx agent` scoring —
// can grep, render, or score individual entries.
//
// [Headings] is generic: it projects any ATX heading (`##`, `###`, …) up to
// a caller-supplied depth, code-fence-aware. This is what lets one projector
// serve timestamped entry files and untimestamped ones alike (TASKS.md
// `## Phase …`, CONVENTIONS.md), and it backs the `ctx index` command.
//
// # Supersession
//
// An entry can be marked superseded by a later one (a body line starting
// with the strikethrough `~~Superseded` prefix). [EntryBlock.IsSuperseded]
// reports this so renderers can gray-out or sort accordingly.
//
// # Concurrency
//
// Pure data: the package reads content strings and returns values. Callers
// own all file IO.
package heading
