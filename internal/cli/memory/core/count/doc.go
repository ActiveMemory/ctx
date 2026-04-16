//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package count provides line counting utilities for memory
// file analysis.
//
// [FileLines] counts the number of newline characters in
// raw file bytes using [bytes.Count]. This gives a fast
// approximation of line count without parsing the file
// content. The count is newline-based (LF), consistent with
// the project's LF-only convention.
//
// The memory status command uses this to report source and
// mirror line counts side by side. A mismatch between source
// and mirror line counts indicates drift: the external
// memory source (e.g., Claude Code's MEMORY.md) has changed
// since the last sync, or the local mirror has been edited
// independently.
//
// # Design Choice
//
// Counting newlines rather than splitting into lines avoids
// allocating a string slice for what is purely a numeric
// query. For large memory files this keeps the status
// command's memory footprint proportional to file size,
// not line count.
package count
