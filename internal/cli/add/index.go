//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

// This file contains backward-compatible aliases for index operations
// that delegate to the internal/index package.

import (
	"github.com/ActiveMemory/ctx/internal/index"
)

// IndexEntry represents a parsed entry header from a context file.
//
// This is an alias for index.Entry for backward compatibility.
type IndexEntry = index.Entry

// DecisionEntry is an alias for IndexEntry for backward compatibility.
type DecisionEntry = index.Entry

// ParseEntryHeaders extracts all entries from file content.
//
// Delegates to index.ParseHeaders.
//
// Parameters:
//   - content: The full content of a context file
//
// Returns:
//   - []IndexEntry: Slice of parsed entries (may be empty)
func ParseEntryHeaders(content string) []IndexEntry {
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
//   - []DecisionEntry: Slice of parsed entries (it may be empty)
func ParseDecisionHeaders(content string) []DecisionEntry {
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
func GenerateIndexTable(entries []IndexEntry, columnHeader string) string {
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
func GenerateIndex(entries []DecisionEntry) string {
	return index.GenerateTable(entries, "Decision")
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
