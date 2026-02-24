//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package index provides index generation and parsing for context files.
package index

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/ActiveMemory/ctx/internal/config"
)

// Entry represents a parsed entry header from a context file.
//
// Fields:
//   - Timestamp: Full timestamp (YYYY-MM-DD-HHMMSS)
//   - Date: Date only (YYYY-MM-DD)
//   - Title: Entry title
type Entry struct {
	Timestamp string
	Date      string
	Title     string
}

// ParseHeaders extracts all entries from file content.
//
// It scans for headers matching the pattern "## [YYYY-MM-DD-HHMMSS] Title"
// and returns them in the order they appear in the file.
//
// Parameters:
//   - content: The full content of a context file
//
// Returns:
//   - []Entry: Slice of parsed entries (it may be empty)
func ParseHeaders(content string) []Entry {
	var entries []Entry

	matches := config.RegExEntryHeader.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 4 {
			date := match[1]
			time := match[2]
			title := match[3]
			entries = append(entries, Entry{
				Timestamp: date + "-" + time,
				Date:      date,
				Title:     title,
			})
		}
	}

	return entries
}

// GenerateTable creates a Markdown table index from entries.
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
func GenerateTable(entries []Entry, columnHeader string) string {
	if len(entries) == 0 {
		return ""
	}

	nl := config.NewlineLF
	var sb strings.Builder
	sb.WriteString("| Date | ")
	sb.WriteString(columnHeader)
	sb.WriteString(" |" + nl)
	sb.WriteString("|------|")
	sb.WriteString(strings.Repeat("-", len(columnHeader)))
	sb.WriteString("|" + nl)

	for _, e := range entries {
		// Escape pipe characters in title
		title := strings.ReplaceAll(e.Title, "|", "\\|")
		sb.WriteString("| ")
		sb.WriteString(e.Date)
		sb.WriteString(" | ")
		sb.WriteString(title)
		sb.WriteString(" |" + nl)
	}

	return sb.String()
}

// Update regenerates the index in file content.
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
func Update(content, fileHeader, columnHeader string) string {
	entries := ParseHeaders(content)
	indexContent := GenerateTable(entries, columnHeader)
	nl := config.NewlineLF

	// Check if markers already exist
	startIdx := strings.Index(content, config.IndexStart)
	endIdx := strings.Index(content, config.IndexEnd)

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		// Replace the existing index
		if indexContent == "" {
			// No entries - remove index entirely (including markers
			// and surrounding whitespace)
			before := strings.TrimRight(content[:startIdx], nl)
			after := content[endIdx+len(config.IndexEnd):]
			after = strings.TrimLeft(after, nl)
			if after != "" {
				return before + nl + nl + after
			}
			return before + nl
		}
		// Replace content between markers
		before := content[:startIdx+len(config.IndexStart)]
		after := content[endIdx:]
		return before + nl + indexContent + after
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
	lineEnd := strings.Index(content[headerIdx:], nl)
	if lineEnd == -1 {
		// Header is at the end of the file
		return content + nl + nl +
			config.IndexStart + nl + indexContent +
			config.IndexEnd + nl
	}

	insertPoint := headerIdx + lineEnd + 1

	// Build new content with the index
	var sb strings.Builder
	sb.WriteString(content[:insertPoint])
	sb.WriteString(nl)
	sb.WriteString(config.IndexStart)
	sb.WriteString(nl)
	sb.WriteString(indexContent)
	sb.WriteString(config.IndexEnd)
	sb.WriteString(nl)
	sb.WriteString(content[insertPoint:])

	return sb.String()
}

// UpdateDecisions regenerates the decision index in DECISIONS.md content.
//
// Parameters:
//   - content: The full content of DECISIONS.md
//
// Returns:
//   - string: Updated content with regenerated index
func UpdateDecisions(content string) string {
	return Update(content, config.HeadingDecisions, config.ColumnDecision)
}

// UpdateLearnings regenerates the learning index in LEARNINGS.md content.
//
// Parameters:
//   - content: The full content of LEARNINGS.md
//
// Returns:
//   - string: Updated content with regenerated index
func UpdateLearnings(content string) string {
	return Update(content, config.HeadingLearnings, config.ColumnLearning)
}

// ReindexFile reads a context file, regenerates its index, and writes it back.
//
// This is a convenience function that handles the common reindex workflow:
// check the file exists, read content, apply update function, write back,
// report.
//
// Note: This function uses io.Writer instead of *cobra.Command to keep the
// index package decoupled from CLI concerns. Callers pass cmd.OutOrStdout()
// which writes to the same destination as cmd.Printf.
//
// Parameters:
//   - w: Writer for status output (typically cmd.OutOrStdout())
//   - filePath: Full path to the context file
//   - fileName: Display name for error messages (e.g., "DECISIONS.md")
//   - updateFunc: Function to regenerate the index (e.g., UpdateDecisions)
//   - entryType: Plural noun for the status message (e.g., "decisions")
//
// Returns:
//   - error: Non-nil if file operations fail
func ReindexFile(
	w io.Writer, filePath, fileName string,
	updateFunc func(string) string,
	entryType string,
) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("%s not found. Run 'ctx init' first", fileName)
	}

	content, err := os.ReadFile(filePath) //nolint:gosec // G304: filePath is constructed from known config paths
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	updated := updateFunc(string(content))

	if err := os.WriteFile(filePath, []byte(updated), config.PermFile); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	entries := ParseHeaders(string(content))
	green := color.New(color.FgGreen).SprintFunc()
	if len(entries) == 0 {
		_, err := fmt.Fprintf(
			w, "%s Index cleared (no %s found)\n", green("✓"), entryType)
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintf(
			w,
			"%s Index regenerated with %d entries\n", green("✓"), len(entries),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
