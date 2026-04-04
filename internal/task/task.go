//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task provides task item parsing and matching.
//
// This package handles the domain logic for task items, independent of
// their Markdown representation.
package task

import (
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/marker"
)

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
	// matchFull is the index of the full regex match.
	matchFull = iota
	// MatchIndent is the index of leading whitespace.
	MatchIndent
	// MatchState is the index of the checkbox state.
	MatchState
	// MatchContent is the index of the task text.
	MatchContent
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
	return match[MatchState] == marker.MarkTaskComplete
}

// Pending reports whether a match represents a pending task.
//
// Parameters:
//   - match: Result from ItemPattern.FindStringSubmatch
//
// Returns:
//   - bool: True if the checkbox is unchecked ([ ])
func Pending(match []string) bool {
	if len(match) <= MatchState {
		return false
	}
	return match[MatchState] != marker.MarkTaskComplete
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

// Sub reports whether a match represents a subtask (indented).
//
// Parameters:
//   - match: Result from ItemPattern.FindStringSubmatch
//
// Returns:
//   - bool: True if indent is 2+ spaces
func Sub(match []string) bool {
	return len(Indent(match)) >= archive.SubTaskMinIndent
}
