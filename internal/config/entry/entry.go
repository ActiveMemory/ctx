//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

// Entry type constants for context updates.
//
// These are the canonical internal representations used in switch statements
// for routing add/update commands to the appropriate handler.
const (
	// Task represents a task entry in TASKS.md.
	Task = "task"
	// Decision represents an architectural decision in DECISIONS.md.
	Decision = "decision"
	// Learning represents a lesson learned in LEARNINGS.md.
	Learning = "learning"
	// Convention represents a code pattern in CONVENTIONS.md.
	Convention = "convention"
	// Complete represents a task completion action (marks the task as done).
	Complete = "complete"
	// Unknown is returned when user input doesn't match any known type.
	Unknown = "unknown"
)
