//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task provides task item parsing and matching.
//
// This package handles the domain logic for task items, independent of
// their Markdown representation.
package task

import "github.com/ActiveMemory/ctx/internal/config"

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
	MatchFull    = iota // Full match
	MatchIndent         // Leading whitespace
	MatchState          // "x" or " " or ""
	MatchContent        // Task text
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
	return match[MatchState] == config.MarkTaskComplete
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
	return match[MatchState] != config.MarkTaskComplete
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
//   - string: Task content (empty if the match is invalid)
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
