//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package complete implements the "ctx task complete"
// cobra subcommand.
//
// This command marks a single task as done in
// TASKS.md by changing its checkbox from "- [ ]" to
// "- [x]". Tasks can be identified by number, partial
// text match, or full text.
//
// # Usage
//
//	ctx task complete <query>
//
// # Arguments
//
// Exactly one positional argument is required:
//
//   - query: a task number (e.g. "3"), a partial text
//     match (e.g. "fix auth"), or the full task text.
//     The query is matched against pending tasks in
//     TASKS.md.
//
// # Behavior
//
// The command delegates to task/core/complete which:
//
//   - Reads TASKS.md and finds the matching pending
//     task.
//   - Rewrites TASKS.md with the matched task's
//     checkbox toggled to [x].
//   - Returns the matched task text and its number.
//
// After completion, the command records a trace
// reference (type "task", value is the task number)
// for commit tracing so the next git commit can be
// linked back to the completed task.
//
// # Output
//
// Prints a confirmation line showing the completed
// task text.
//
// # Delegation
//
// Task matching and file rewriting are in
// task/core/complete. Trace recording uses the
// trace package. Output formatting uses
// write/complete.
package complete
