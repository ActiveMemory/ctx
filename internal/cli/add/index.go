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

// Index markers for DECISIONS.md
const (
	IndexStart = "<!-- INDEX:START -->"
	IndexEnd   = "<!-- INDEX:END -->"
)

// DecisionEntry represents a parsed decision header.
type DecisionEntry struct {
	Timestamp string // Full timestamp: YYYY-MM-DD-HHMMSS
	Date      string // Date only: YYYY-MM-DD
	Title     string // Decision title
}

// decisionHeaderRegex matches decision headers like "## [2026-01-28-051426] Title here"
var decisionHeaderRegex = regexp.MustCompile(`## \[(\d{4}-\d{2}-\d{2})-(\d{6})\] (.+)`)

// ParseDecisionHeaders extracts all decision entries from file content.
//
// It scans for headers matching the pattern "## [YYYY-MM-DD-HHMMSS] Title"
// and returns them in the order they appear in the file.
//
// Parameters:
//   - content: The full content of DECISIONS.md
//
// Returns:
//   - []DecisionEntry: Slice of parsed entries (may be empty)
func ParseDecisionHeaders(content string) []DecisionEntry {
	var entries []DecisionEntry

	matches := decisionHeaderRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 4 {
			date := match[1]
			time := match[2]
			title := match[3]
			entries = append(entries, DecisionEntry{
				Timestamp: date + "-" + time,
				Date:      date,
				Title:     title,
			})
		}
	}

	return entries
}

// GenerateIndex creates a markdown table index from decision entries.
//
// The table has two columns: Date and Decision title.
// If there are no entries, returns an empty string.
//
// Parameters:
//   - entries: Slice of decision entries to include
//
// Returns:
//   - string: Markdown table (without markers) or empty string
func GenerateIndex(entries []DecisionEntry) string {
	if len(entries) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("| Date | Decision |\n")
	sb.WriteString("|------|----------|\n")

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

// UpdateIndex regenerates the decision index in file content.
//
// If INDEX:START and INDEX:END markers exist, the content between them
// is replaced. Otherwise, the index is inserted after "# Decisions".
// If there are no decision entries, any existing index is removed.
//
// Parameters:
//   - content: The full content of DECISIONS.md
//
// Returns:
//   - string: Updated content with regenerated index
func UpdateIndex(content string) string {
	entries := ParseDecisionHeaders(content)
	indexContent := GenerateIndex(entries)

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

	// No existing markers - insert after "# Decisions" header
	if indexContent == "" {
		// No entries, nothing to insert
		return content
	}

	headerIdx := strings.Index(content, "# Decisions")
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
