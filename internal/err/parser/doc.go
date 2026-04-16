//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parser defines the typed error constructors
// for session transcript parsing. These errors fire
// when scanning, opening, or parsing AI tool session
// files (Claude Code JSONL, Aider markdown, etc.)
// and when validating frontmatter delimiters.
//
// # Domain
//
// Errors fall into three categories:
//
//   - **Frontmatter**: a session file is missing
//     its opening or closing --- delimiter.
//     Constructors: [MissingOpenDelim],
//     [MissingCloseDelim].
//   - **File IO**: a session file could not be
//     read, opened, scanned, or walked.
//     Constructors: [ReadFile], [OpenFile],
//     [ScanFile], [WalkDir].
//   - **Parse**: no parser matches a file, JSON
//     unmarshaling failed, or a per-file parse
//     error occurred. Constructors: [NoMatch],
//     [Unmarshal], [FileError], [ParseFile].
//
// # Wrapping Strategy
//
// IO and parse constructors wrap their cause with
// fmt.Errorf %w so callers can inspect the
// underlying error. Frontmatter constructors
// return plain errors. All user-facing text is
// resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package parser
