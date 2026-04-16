//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package file implements file-level context tracing. It
// shows which context references were active when each
// commit touched a given file.
//
// # Path Argument Parsing
//
// [ParsePathArg] strips an optional :line-range suffix
// from a path argument so git log receives a clean file
// path. Recognized suffixes are numeric ranges like
// ":42" or ":42-60". Non-numeric suffixes such as
// ":latest" are left intact.
//
// # File Tracing
//
// [Trace] runs git log for the given file, limited to
// the last N commits. For each commit it collects
// context refs by combining history entries from
// .context/trace/ with any override annotations. Each
// commit is printed as a single line containing the
// short hash, date, subject, and attached context refs.
//
// The output gives developers a chronological view of
// how context references evolved alongside the file's
// changes.
//
// # Data Flow
//
//  1. CLI parses path argument via ParsePathArg
//  2. Trace calls git log for the clean file path
//  3. For each commit, refs are collected from history
//  4. Each commit line is formatted and printed
package file
