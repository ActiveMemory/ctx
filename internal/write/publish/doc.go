//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package publish provides terminal output for the
// memory publish command (ctx memory publish).
//
// Publishing compiles selected context files into a
// single MEMORY.md block with a line budget. The
// output functions narrate each stage of that
// process.
//
// # Planning
//
// [Plan] displays the full publish plan: a header,
// source file list, line budget, per-file counts
// (tasks, decisions, conventions, learnings), and
// the total line count versus the budget.
//
// # Execution
//
// [Done] confirms a successful publish with marker
// information. [DryRun] prints a notice that no
// changes were written.
//
// # Unpublish
//
// [NotFound] reports that no published block exists
// in the target file. [Unpublished] confirms that
// the published block was removed.
package publish
