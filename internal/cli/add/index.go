//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"regexp"
	"strings"
)

// Index markers used in context files
const (
	IndexStart = "<!-- INDEX:START -->"
	IndexEnd   = "<!-- INDEX:END -->"
)

// IndexEntry represents a parsed entry header from a context file.
//
// Fields:
//   - Timestamp: Full timestamp (YYYY-MM-DD-HHMMSS)
//   - Date: Date only (YYYY-MM-DD)
//   - Title: Entry title
type IndexEntry struct {
	Timestamp string
	Date      string
	Title     string
}

// DecisionEntry is an alias for backward compatibility.
type DecisionEntry = IndexEntry

// entryHeaderRegex matches headers like "## [2026-01-28-051426] Title here"
var entryHeaderRegex = regexp.MustCompile(`## \[(\d{4}-\d{2}-\d{2})-(\d{6})\] (.+)`)

// ParseEntryHeaders extracts all entries from file content.
//
// It scans for headers matching the pattern "## [YYYY-MM-DD-HHMMSS] Title"
// and returns them in the order they appear in the file.
//
// Parameters:
//   - content: The full content of a context file
//
// Returns:
//   - []IndexEntry: Slice of parsed entries (may be empty)
func ParseEntryHeaders(content string) []IndexEntry {
	var entries []IndexEntry

	matches := entryHeaderRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 4 {
			date := match[1]
			time := match[2]
			title := match[3]
			entries = append(entries, IndexEntry{
				Timestamp: date + "-" + time,
				Date:      date,
				Title:     title,
			})
		}
	}

	return entries
}

// ParseDecisionHeaders is an alias for ParseEntryHeaders for backward compatibility.
//
// Parameters:
//   - content: The full content of a context file
//
// Returns:
//   - []DecisionEntry: Slice of parsed entries (may be empty)
func ParseDecisionHeaders(content string) []DecisionEntry {
	return ParseEntryHeaders(content)
}

// GenerateIndexTable creates a markdown table index from entries.
//
// The table has two columns: Date and the specified column header.
// If there are no entries, returns an empty string.
//
// Parameters:
//   - entries: Slice of entries to include
//   - columnHeader: Header for the second column (e.g., "Decision", "Learning")
//
// Returns:
//   - string: Markdown table (without markers) or empty string
func GenerateIndexTable(entries []IndexEntry, columnHeader string) string {
	if len(entries) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("| Date | ")
	sb.WriteString(columnHeader)
	sb.WriteString(" |\n")
	sb.WriteString("|------|")
	sb.WriteString(strings.Repeat("-", len(columnHeader)))
	sb.WriteString("|\n")

	for _, e := range entries {
		// Escape pipe characters in title
		title := strings.ReplaceAll(e.Title, "|", "\\|")
		sb.WriteString("| ")
		sb.WriteString(e.Date)
		sb.WriteString(" | ")
		sb.WriteString(title)
		sb.WriteString(" |\n")
	}

	return sb.String()
}

// GenerateIndex creates a markdown table for decisions (backward compatibility).
//
// Parameters:
//   - entries: Slice of decision entries to include
//
// Returns:
//   - string: Markdown table or empty string if no entries
func GenerateIndex(entries []DecisionEntry) string {
	return GenerateIndexTable(entries, "Decision")
}

// updateFileIndex regenerates the index in file content.
//
// If INDEX:START and INDEX:END markers exist, the content between them
// is replaced. Otherwise, the index is inserted after the specified header.
// If there are no entries, any existing index is removed.
//
// Parameters:
//   - content: The full content of the file
//   - fileHeader: The main header to insert after (e.g., "# Decisions")
//   - columnHeader: Header for the table column (e.g., "Decision")
//
// Returns:
//   - string: Updated content with regenerated index
func updateFileIndex(content, fileHeader, columnHeader string) string {
	entries := ParseEntryHeaders(content)
	indexContent := GenerateIndexTable(entries, columnHeader)

	// Check if markers already exist
	startIdx := strings.Index(content, IndexStart)
	endIdx := strings.Index(content, IndexEnd)

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		// Replace existing index
		if indexContent == "" {
			// No entries - remove index entirely (including markers and surrounding whitespace)
			before := strings.TrimRight(content[:startIdx], "\n")
			after := content[endIdx+len(IndexEnd):]
			after = strings.TrimLeft(after, "\n")
			if after != "" {
				return before + "\n\n" + after
			}
			return before + "\n"
		}
		// Replace content between markers
		before := content[:startIdx+len(IndexStart)]
		after := content[endIdx:]
		return before + "\n" + indexContent + after
	}

	// No existing markers - insert after file header
	if indexContent == "" {
		// No entries, nothing to insert
		return content
	}

	headerIdx := strings.Index(content, fileHeader)
	if headerIdx == -1 {
		// No header found, return unchanged
		return content
	}

	// Find end of header line
	lineEnd := strings.Index(content[headerIdx:], "\n")
	if lineEnd == -1 {
		// Header is at end of file
		return content + "\n\n" + IndexStart + "\n" + indexContent + IndexEnd + "\n"
	}

	insertPoint := headerIdx + lineEnd + 1

	// Build new content with index
	var sb strings.Builder
	sb.WriteString(content[:insertPoint])
	sb.WriteString("\n")
	sb.WriteString(IndexStart)
	sb.WriteString("\n")
	sb.WriteString(indexContent)
	sb.WriteString(IndexEnd)
	sb.WriteString("\n")
	sb.WriteString(content[insertPoint:])

	return sb.String()
}

// UpdateIndex regenerates the decision index in DECISIONS.md content.
//
// Parameters:
//   - content: The full content of DECISIONS.md
//
// Returns:
//   - string: Updated content with regenerated index
func UpdateIndex(content string) string {
	return updateFileIndex(content, "# Decisions", "Decision")
}

// UpdateLearningsIndex regenerates the learning index in LEARNINGS.md content.
//
// Parameters:
//   - content: The full content of LEARNINGS.md
//
// Returns:
//   - string: Updated content with regenerated index
func UpdateLearningsIndex(content string) string {
	return updateFileIndex(content, "# Learnings", "Learning")
}
