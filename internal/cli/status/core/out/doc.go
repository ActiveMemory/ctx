//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package out renders context status as JSON or
// formatted text output.
//
// # JSON Output
//
// [PersistStatusJSON] encodes the context as an Output
// struct containing directory path, file count, total
// tokens, total size, and per-file status entries.
// When verbose mode is enabled, each file entry
// includes a content preview. The output is written
// as indented JSON to the command's output stream.
//
// # Text Output
//
// [PersistStatusText] renders a human-readable report.
// It prints a header with directory path and totals,
// then lists each file sorted by read priority. Files
// show an indicator icon (empty or OK), token count,
// size, and summary. In verbose mode, content previews
// are appended. A recent activity section shows the
// most recently modified files with relative
// timestamps.
//
// # Types
//
// [Output] models the top-level JSON structure with
// aggregate statistics and a file list. [FileStatus]
// models a single file's metadata including name,
// tokens, size, emptiness flag, summary, modification
// time, and optional preview lines.
package out
