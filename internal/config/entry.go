//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.

package config

import "strings"

// Entry type constants for context updates.
//
// These are the canonical internal representations used in switch statements
// for routing add/update commands to the appropriate handler.
const (
	// EntryTask represents a task entry in TASKS.md.
	EntryTask = "task"
	// EntryDecision represents an architectural decision in DECISIONS.md.
	EntryDecision = "decision"
	// EntryLearning represents a lesson learned in LEARNINGS.md.
	EntryLearning = "learning"
	// EntryConvention represents a code pattern in CONVENTIONS.md.
	EntryConvention = "convention"
	// EntryComplete represents a task completion action (marks the task as done).
	EntryComplete = "complete"
	// EntryUnknown is returned when user input doesn't match any known type.
	EntryUnknown = "unknown"
)

// UserInputToEntry normalizes user input to a canonical entry type.
//
// Accepts both singular and plural forms (e.g., "task" or "tasks") and
// returns the canonical singular form. Matching is case-insensitive.
// Unknown inputs return EntryUnknown.
//
// Parameters:
//   - s: User-provided entry type string
//
// Returns:
//   - string: Canonical entry type constant (EntryTask, EntryDecision, etc.)
func UserInputToEntry(s string) string {
	switch strings.ToLower(s) {
	case "task", "tasks":
		return EntryTask
	case "decision", "decisions":
		return EntryDecision
	case "learning", "learnings":
		return EntryLearning
	case "convention", "conventions":
		return EntryConvention
	default:
		return EntryUnknown
	}
}
