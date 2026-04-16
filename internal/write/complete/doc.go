//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package complete provides terminal output for the task
// completion command (ctx task complete).
//
// # Exported Functions
//
// [Completed] prints a confirmation message when a task
// checkbox is toggled from [ ] to [x] in TASKS.md. The
// message includes the task description so the user can
// verify which task was marked done.
//
// # Message Categories
//
//   - Info: task completion confirmation with the task
//     description echoed back to the user
//
// # Usage
//
// The calling command identifies the target task by
// index or text match, toggles the checkbox, and then
// calls Completed to confirm the change:
//
//	complete.Completed(cmd, "Implement session cooldown")
package complete
