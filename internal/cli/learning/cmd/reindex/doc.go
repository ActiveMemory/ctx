//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex implements the ctx learning reindex subcommand.
//
// [Cmd] builds the cobra.Command. [Run] parses all timestamped
// entry headers in LEARNINGS.md, regenerates the index table at
// the top of the file sorted by date, and writes the result.
package reindex
