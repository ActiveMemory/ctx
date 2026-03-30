//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package complete provides terminal output for the task completion
// command (ctx complete).
//
// The single exported function [Completed] prints a confirmation
// message when a task checkbox is toggled from [ ] to [x] in
// TASKS.md.
//
// Example:
//
//	write.Completed(cmd, "Implement session cooldown")
package complete
