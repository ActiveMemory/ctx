//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// This file contains backward-compatible aliases for index operations
// that delegate to the internal/index package.

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/index"
)

// ParseEntryHeaders extracts all entries from file content.
//
// Delegates to index.ParseHeaders.
//
// Parameters:
//   - content: The full content of a context file
//
// Returns:
//   - []index.Entry: Slice of parsed entries (may be empty)
func ParseEntryHeaders(content string) []index.Entry {
	return index.ParseHeaders(content)
}

// ParseDecisionHeaders extracts all entries from file content.
//
// This is an alias for ParseEntryHeaders for backward compatibility.
//
// Parameters:
//   - content: The full content of a context file
//
// Returns:
//   - []index.Entry: Slice of parsed entries (it may be empty)
func ParseDecisionHeaders(content string) []index.Entry {
	return index.ParseHeaders(content)
}

// GenerateIndexTable creates a Markdown table index from entries.
//
// Delegates to index.GenerateTable.
//
// Parameters:
//   - entries: Slice of entries to include
//   - columnHeader: Header for the second column (e.g., "Decision", "Learning")
//
// Returns:
//   - string: Markdown table (without markers) or empty string
func GenerateIndexTable(entries []index.Entry, columnHeader string) string {
	return index.GenerateTable(entries, columnHeader)
}

// GenerateIndex creates a Markdown table for decisions.
//
// This is a convenience wrapper for backward compatibility.
//
// Parameters:
//   - entries: Slice of decision entries to include
//
// Returns:
//   - string: Markdown table or empty string if no entries
func GenerateIndex(entries []index.Entry) string {
	return index.GenerateTable(entries, desc.TextDesc(text.DescKeyColumnDecision))
}

// UpdateIndex regenerates the decision index in DECISIONS.md content.
//
// Delegates to index.UpdateDecisions.
//
// Parameters:
//   - content: The full content of DECISIONS.md
//
// Returns:
//   - string: Updated content with regenerated index
func UpdateIndex(content string) string {
	return index.UpdateDecisions(content)
}

// UpdateLearningsIndex regenerates the learning index in LEARNINGS.md content.
//
// Delegates to index.UpdateLearnings.
//
// Parameters:
//   - content: The full content of LEARNINGS.md
//
// Returns:
//   - string: Updated content with regenerated index
func UpdateLearningsIndex(content string) string {
	return index.UpdateLearnings(content)
}
