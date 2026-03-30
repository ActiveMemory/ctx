//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex implements the ctx decision reindex subcommand.
//
// [Cmd] builds the cobra.Command. [Run] regenerates the index
// table at the top of DECISIONS.md by parsing all entry headers
// and rebuilding the sorted table.
package reindex
