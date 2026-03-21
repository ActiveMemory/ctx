//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import "strings"

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

// Plural forms used as labels and resource identifiers.
const (
	Decisions   = "decisions"
	Learnings   = "learnings"
	Conventions = "conventions"
	Tasks       = "tasks"
)

// Priority levels for task entries.
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// Priorities lists all valid priority levels for shell completion.
var Priorities = []string{PriorityHigh, PriorityMedium, PriorityLow}

// FromUserInput normalizes user input to a canonical entry type.
//
// Accepts singular and plural forms, case-insensitive.
//
// Parameters:
//   - s: user-supplied type string (e.g. "tasks", "Decision")
//
// Returns:
//   - string: canonical entry constant, or Unknown
func FromUserInput(s string) string {
	switch strings.ToLower(s) {
	case "task", "tasks":
		return Task
	case "decision", "decisions":
		return Decision
	case "learning", "learnings":
		return Learning
	case "convention", "conventions":
		return Convention
	default:
		return Unknown
	}
}
