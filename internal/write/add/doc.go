//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add provides terminal output for the context
// entry addition commands (ctx task add, ctx decision add,
// ctx learning add, ctx convention add).
//
// # Exported Functions
//
// [Added] prints a confirmation message after an entry
// is appended to a context file. The message includes
// the target filename so the user knows which file was
// modified (e.g. TASKS.md, DECISIONS.md).
//
// [SpecNudge] prints a one-line tip suggesting that
// the user create a feature spec when a task is complex
// enough to benefit from structured planning. This nudge
// is shown conditionally by the calling command.
//
// # Message Categories
//
//   - Info: confirmation that an entry was written
//   - Nudge: optional guidance tip after task creation
//
// # Usage
//
// Both functions accept a *cobra.Command for output
// routing. Messages are loaded from the embedded
// descriptor system and formatted with the target
// filename.
//
//	add.Added(cmd, "TASKS.md")
//	add.SpecNudge(cmd)
package add
