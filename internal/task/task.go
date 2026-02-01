//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task provides task item parsing and matching.
//
// This package handles the domain logic for task items, independent of
// their markdown representation.
package task

// Match indices for accessing capture groups.
//
// Usage:
//
//	match := task.ItemPattern.FindStringSubmatch(line)
//	if match != nil {
//	    indent := match[task.MatchIndent]
//	    state := match[task.MatchState]
//	    content := match[task.MatchContent]
//	}
const (
	MatchFull    = 0 // Full match
	MatchIndent  = 1 // Leading whitespace
	MatchState   = 2 // "x" or " " or ""
	MatchContent = 3 // Task text
)

// Completed reports whether a match represents a completed task.
//
// Parameters:
//   - match: Result from ItemPattern.FindStringSubmatch
//
// Returns:
//   - bool: True if the checkbox is checked ([x])
func Completed(match []string) bool {
	if len(match) <= MatchState {
		return false
	}
	return match[MatchState] == "x"
}

// IsPending reports whether a match represents a pending task.
//
// Parameters:
//   - match: Result from ItemPattern.FindStringSubmatch
//
// Returns:
//   - bool: True if the checkbox is unchecked ([ ])
func IsPending(match []string) bool {
	if len(match) <= MatchState {
		return false
	}
	return match[MatchState] != "x"
}

// Indent returns the leading whitespace from a match.
//
// Parameters:
//   - match: Result from ItemPattern.FindStringSubmatch
//
// Returns:
//   - string: Indent string (may be empty for top-level tasks)
func Indent(match []string) string {
	if len(match) <= MatchIndent {
		return ""
	}
	return match[MatchIndent]
}

// Content returns the task text from a match.
//
// Parameters:
//   - match: Result from ItemPattern.FindStringSubmatch
//
// Returns:
//   - string: Task content (empty if match is invalid)
func Content(match []string) string {
	if len(match) <= MatchContent {
		return ""
	}
	return match[MatchContent]
}

// IsSubTask reports whether a match represents a subtask (indented).
//
// Parameters:
//   - match: Result from ItemPattern.FindStringSubmatch
//
// Returns:
//   - bool: True if indent is 2+ spaces
func IsSubTask(match []string) bool {
	return len(Indent(match)) >= 2
}
