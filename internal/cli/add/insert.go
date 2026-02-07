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

// insertTask inserts a task entry into TASKS.md.
//
// When section is explicitly provided, inserts after that section header.
// When section is empty (default), finds the first unchecked task and
// inserts before it so the new task lands among existing pending work.
// Falls back to appending at the end if neither is found.
//
// Parameters:
//   - entry: Formatted task entry to insert
//   - existingStr: Existing file content
//   - section: Explicit section name, or empty for auto-placement
//
// Returns:
//   - []byte: Modified content with task inserted
func insertTask(entry, existingStr, section string) []byte {
	// Explicit section: honor it.
	if section != "" {
		return insertTaskAfterSection(entry, existingStr, section)
	}

	// Default: insert before the first unchecked task.
	pendingIdx := strings.Index(existingStr, config.PrefixTaskUndone)
	if pendingIdx != -1 {
		return []byte(existingStr[:pendingIdx] + entry +
			config.NewlineLF + existingStr[pendingIdx:])
	}

	// No unchecked tasks: append at the end.
	if !endsWithNewline(existingStr) {
		existingStr += config.NewlineLF
	}
	return []byte(existingStr + config.NewlineLF + entry)
}

// insertTaskAfterSection inserts a task after a named section header.
//
// Normalizes the section name to a markdown heading, finds it in the
// content, and inserts the entry immediately after. Falls back to
// appending at the end if the header is not found.
//
// Parameters:
//   - entry: Formatted task entry to insert
//   - content: Existing file content
//   - section: Section name (e.g., "Phase 1", "Maintenance")
//
// Returns:
//   - []byte: Modified content with task inserted
func insertTaskAfterSection(entry, content, section string) []byte {
	header := normalizeTargetSection(section)

	found, idx := contains(content, header)
	if !found {
		if !endsWithNewline(content) {
			content += config.NewlineLF
		}
		return []byte(content + config.NewlineLF + entry)
	}

	hasNewLine, lineEnd := containsNewLine(content[idx:])
	if hasNewLine {
		insertPoint := idx + lineEnd
		insertPoint = skipNewline(content, insertPoint)
		return []byte(content[:insertPoint] + config.NewlineLF +
			entry + content[insertPoint:])
	}

	return []byte(content + config.NewlineLF + entry)
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
