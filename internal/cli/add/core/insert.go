//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// InsertAfterHeader finds a header line and inserts content after it.
//
// Skips blank lines and HTML comment blocks (<!-- ... -->) between the header
// and the insertion point, so new entries land after index tables, format
// guides, and other comment-wrapped metadata. Falls back to appending at the
// end if the header is not found.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Header line to find (e.g., "# Learnings")
//
// Returns:
//   - []byte: Modified content with entry inserted
func InsertAfterHeader(content, entry, header string) []byte {
	hasHeader, idx := Contains(content, header)
	if !hasHeader {
		return AppendAtEnd(content, entry)
	}

	hasNewLine, lineEnd := ContainsNewLine(content[idx:])
	if !hasNewLine {
		// Header exists but no newline after (the file ends with a header line)
		return AppendAtEnd(content, entry)
	}

	insertPoint := idx + lineEnd
	insertPoint = SkipNewline(content, insertPoint)

	// Skip blank lines and any HTML comment blocks (<!-- ... -->).
	// This handles INDEX markers, format-guide comments, and ctx markers alike.
	for insertPoint < len(content) {
		if n := SkipNewline(content, insertPoint); n > insertPoint {
			insertPoint = n
			continue
		}

		// Not an HTML comment: we found the insertion point.
		if !strings.HasPrefix(content[insertPoint:], config.CommentOpen) {
			break
		}

		// Skip past the closing --> of this comment block.
		hasCommentEnd, endIdx := ContainsEndComment(content[insertPoint:])
		if !hasCommentEnd {
			break
		}

		insertPoint += endIdx + len(config.CommentClose)
		insertPoint = SkipWhitespace(content, insertPoint)
	}

	return []byte(content[:insertPoint] + entry)
}

// AppendAtEnd appends an entry at the end of content.
//
// Ensures proper newline separation between existing content and the new entry.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to append
//
// Returns:
//   - []byte: Content with entry appended
func AppendAtEnd(content, entry string) []byte {
	if !EndsWithNewline(content) {
		content += config.NewlineLF
	}
	return []byte(content + config.NewlineLF + entry)
}

// InsertTask inserts a task entry into TASKS.md.
//
// When section is explicitly provided, inserts after that section header.
// When section is empty (default), finds the first unchecked task and
// inserts before it, so the new task lands among existing pending work.
// Falls back to appending at the end if neither is found.
//
// Parameters:
//   - entry: Formatted task entry to insert
//   - existingStr: Existing file content
//   - section: Explicit section name, or empty for auto-placement
//
// Returns:
//   - []byte: Modified content with task inserted
func InsertTask(entry, existingStr, section string) []byte {
	// Explicit section: honor it.
	if section != "" {
		return InsertTaskAfterSection(entry, existingStr, section)
	}

	// Default: insert before the first unchecked task.
	pendingIdx := strings.Index(existingStr, config.PrefixTaskUndone)
	if pendingIdx != -1 {
		return []byte(existingStr[:pendingIdx] + entry +
			config.NewlineLF + existingStr[pendingIdx:])
	}

	// No unchecked tasks: append at the end.
	if !EndsWithNewline(existingStr) {
		existingStr += config.NewlineLF
	}
	return []byte(existingStr + config.NewlineLF + entry)
}

// InsertTaskAfterSection inserts a task after a named section header.
//
// Normalizes the section name to a Markdown heading, finds it in the
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
func InsertTaskAfterSection(entry, content, section string) []byte {
	header := NormalizeTargetSection(section)

	found, idx := Contains(content, header)
	if !found {
		if !EndsWithNewline(content) {
			content += config.NewlineLF
		}
		return []byte(content + config.NewlineLF + entry)
	}

	hasNewLine, lineEnd := ContainsNewLine(content[idx:])
	if hasNewLine {
		insertPoint := idx + lineEnd
		insertPoint = SkipNewline(content, insertPoint)
		return []byte(content[:insertPoint] + config.NewlineLF +
			entry + content[insertPoint:])
	}

	return []byte(content + config.NewlineLF + entry)
}

// IsInsideHTMLComment reports whether the position idx in content falls
// inside an HTML comment block (<!-- ... -->).
//
// Parameters:
//   - content: String to check
//   - idx: Position to test
//
// Returns:
//   - bool: True if idx is between a <!-- and its closing -->
func IsInsideHTMLComment(content string, idx int) bool {
	// Find the last <!-- before idx
	openIdx := strings.LastIndex(content[:idx], config.CommentOpen)
	if openIdx == -1 {
		return false
	}
	// Check whether a --> closes that block before idx
	closeIdx := strings.Index(content[openIdx:], config.CommentClose)
	if closeIdx == -1 {
		// Unclosed comment — treat as inside
		return true
	}
	// The comment closes at openIdx+closeIdx; if that position is >= idx,
	// the position is still inside the comment.
	return openIdx+closeIdx+len(config.CommentClose) > idx
}

// InsertDecision inserts a decision entry before existing entries.
//
// Finds the first "## [" marker that is NOT inside an HTML comment block
// and inserts before it, maintaining reverse-chronological order.
// Falls back to InsertAfterHeader if no real entries exist yet.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Header line to insert after (e.g., "# Decisions")
//
// Returns:
//   - []byte: Modified content with entry inserted
func InsertDecision(content, entry, header string) []byte {
	return insertBeforeFirstEntry(content, entry, header)
}

// InsertLearning inserts a learning entry before existing entries.
//
// Finds the first "## [" marker that is NOT inside an HTML comment block
// and inserts before it, maintaining reverse-chronological order.
// Falls back to InsertAfterHeader if no real entries exist yet.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//
// Returns:
//   - []byte: Modified content with entry inserted
func InsertLearning(content, entry string) []byte {
	return insertBeforeFirstEntry(content, entry, config.HeadingLearnings)
}
