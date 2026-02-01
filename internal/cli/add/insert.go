//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// insertAfterHeader finds a header line and inserts content after it.
//
// Skips blank lines and ctx markers between the header and insertion point.
// Falls back to appending at the end if the header is not found.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Header line to find (e.g., "# Learnings")
//
// Returns:
//   - []byte: Modified content with entry inserted
func insertAfterHeader(content, entry, header string) []byte {
	hasHeader, idx := contains(content, header)
	if !hasHeader {
		return appendAtEnd(content, entry)
	}

	hasNewLine, lineEnd := containsNewLine(content[idx:])
	if !hasNewLine {
		// Header exists but no newline after (the file ends with a header line)
		return appendAtEnd(content, entry)
	}

	insertPoint := idx + lineEnd
	insertPoint = skipNewline(content, insertPoint)

	// Skip blank lines and ctx markers
	for insertPoint < len(content) {
		if n := skipNewline(content, insertPoint); n > insertPoint {
			insertPoint = n
			continue
		}

		// No context marker: we found the insertion point.
		if !startsWithCtxMarker(content[insertPoint:]) {
			break
		}

		// Skip past the closing marker
		hasCommentEnd, endIdx := containsEndComment(content[insertPoint:])
		if !hasCommentEnd {
			break
		}

		insertPoint += endIdx + len(config.CommentClose)
		insertPoint = skipWhitespace(content, insertPoint)
	}

	return []byte(content[:insertPoint] + entry)
}

// appendAtEnd appends an entry at the end of content.
//
// Ensures proper newline separation between existing content and the new entry.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to append
//
// Returns:
//   - []byte: Content with entry appended
func appendAtEnd(content, entry string) []byte {
	if !endsWithNewline(content) {
		content += config.NewlineLF
	}
	return []byte(content + config.NewlineLF + entry)
}

// insertTask inserts a task entry after a section header in TASKS.md.
//
// Finds the target section (e.g., "## Next Up") and inserts the task
// immediately after the header line. Falls back to appending at the end
// if the section is not found.
//
// Parameters:
//   - existingStr: Existing file content
//   - entry: Formatted task entry to insert
//   - headerSection: Target section name (e.g., "next", "backlog")
//
// Returns:
//   - []byte: Modified content with task inserted
func insertTask(existingStr, entry, headerSection string) []byte {
	targetSectionHeader := normalizeTargetSection(headerSection)

	// Find the section and insert after it
	containsSectionHeader, idx := contains(existingStr, targetSectionHeader)
	if !containsSectionHeader {
		// Section not found: Append at the end.
		if !endsWithNewline(existingStr) {
			existingStr += config.NewlineLF
		}
		return []byte(existingStr + config.NewlineLF + entry)
	}

	// Found section header. Find the end of the section header line
	hasNewLine, lineEnd := containsNewLine(existingStr[idx:])
	if hasNewLine {
		insertPoint := idx + lineEnd
		insertPoint = skipNewline(existingStr, insertPoint)
		return []byte(existingStr[:insertPoint] + config.NewlineLF +
			entry + existingStr[insertPoint:])
	}

	// If no newline; append to the end followed by a newline.
	return []byte(existingStr + config.NewlineLF + entry)
}

// insertDecision inserts a decision entry before existing entries.
//
// Finds the first "## [" marker and inserts before it, maintaining
// reverse-chronological order. Falls back to insertAfterHeader if no entries
// exist.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Header line to insert after (e.g., "# Decisions")
//
// Returns:
//   - []byte: Modified content with entry inserted
func insertDecision(content, entry, header string) []byte {
	// Find the first entry marker "## [" (timestamp-prefixed sections)
	entryIdx := strings.Index(content, "## [")
	if entryIdx != -1 {
		// Insert before the first entry, with a separator after
		return []byte(
			content[:entryIdx] + entry +
				config.NewlineLF + config.Separator +
				config.NewlineLF + config.NewlineLF +
				content[entryIdx:],
		)
	}

	// No existing entries - find the header and insert after it
	return insertAfterHeader(content, entry, header)
}

// insertLearning inserts a learning entry before existing entries.
//
// Finds the first "## [" marker and inserts before it, maintaining
// reverse-chronological order. Falls back to insertAfterHeader if no entries
// exist.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//
// Returns:
//   - []byte: Modified content with entry inserted
func insertLearning(content, entry string) []byte {
	// Find the first entry marker "## [" (timestamp-prefixed sections)
	entryIdx := strings.Index(content, config.HeadingLearningStart)
	if entryIdx != -1 {
		return []byte(
			content[:entryIdx] + entry + config.NewlineLF +
				config.Separator + config.NewlineLF + config.NewlineLF +
				content[entryIdx:],
		)
	}

	// No existing entries - find the header and insert after it
	return insertAfterHeader(content, entry, config.HeadingLearnings)
}
