//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package insert

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/add/core/normalize"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/inspect"
)

// AfterHeader finds a header line and inserts content after it.
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
func AfterHeader(content, entry, header string) []byte {
	hasHeader, idx := inspect.Contains(content, header)
	if !hasHeader {
		return AppendAtEnd(content, entry)
	}

	hasNewLine, lineEnd := inspect.ContainsNewLine(content[idx:])
	if !hasNewLine {
		// Header exists but no newline after (the file ends with a header line)
		return AppendAtEnd(content, entry)
	}

	insertPoint := idx + lineEnd
	insertPoint = inspect.SkipNewline(content, insertPoint)

	// Skip blank lines and any HTML comment blocks (<!-- ... -->).
	// This handles INDEX markers, format-guide comments, and ctx markers alike.
	for insertPoint < len(content) {
		if n := inspect.SkipNewline(content, insertPoint); n > insertPoint {
			insertPoint = n
			continue
		}

		// Not an HTML comment: we found the insertion point.
		if !strings.HasPrefix(content[insertPoint:], marker.CommentOpen) {
			break
		}

		// Skip past the closing --> of this comment block.
		hasCommentEnd, endIdx := inspect.ContainsEndComment(content[insertPoint:])
		if !hasCommentEnd {
			break
		}

		insertPoint += endIdx + len(marker.CommentClose)
		insertPoint = inspect.SkipWhitespace(content, insertPoint)
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
	if !inspect.EndsWithNewline(content) {
		content += token.NewlineLF
	}
	return []byte(content + token.NewlineLF + entry)
}

// Task inserts a task entry into TASKS.md.
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
func Task(entry, existingStr, section string) []byte {
	// Explicit section: honor it.
	if section != "" {
		return TaskAfterSection(entry, existingStr, section)
	}

	// Default: insert before the first unchecked task.
	pendingIdx := strings.Index(existingStr, marker.PrefixTaskUndone)
	if pendingIdx != -1 {
		return []byte(existingStr[:pendingIdx] + entry +
			token.NewlineLF + existingStr[pendingIdx:])
	}

	// No unchecked tasks: append at the end.
	if !inspect.EndsWithNewline(existingStr) {
		existingStr += token.NewlineLF
	}
	return []byte(existingStr + token.NewlineLF + entry)
}

// TaskAfterSection inserts a task after a named section header.
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
//   - []byte: Modified content with the task inserted
func TaskAfterSection(entry, content, section string) []byte {
	header := normalize.TargetSection(section)

	found, idx := inspect.Contains(content, header)
	if !found {
		if !inspect.EndsWithNewline(content) {
			content += token.NewlineLF
		}
		return []byte(content + token.NewlineLF + entry)
	}

	hasNewLine, lineEnd := inspect.ContainsNewLine(content[idx:])
	if hasNewLine {
		insertPoint := idx + lineEnd
		insertPoint = inspect.SkipNewline(content, insertPoint)
		return []byte(content[:insertPoint] + token.NewlineLF +
			entry + content[insertPoint:])
	}

	return []byte(content + token.NewlineLF + entry)
}

// ExistsInsideHTMLComment reports whether the position idx in content falls
// inside an HTML comment block (<!-- ... -->).
//
// Parameters:
//   - content: String to check
//   - idx: Position to test
//
// Returns:
//   - bool: True if idx is between a <!-- and its closing -->
func ExistsInsideHTMLComment(content string, idx int) bool {
	// Find the last <!-- before idx
	openIdx := strings.LastIndex(content[:idx], marker.CommentOpen)
	if openIdx == -1 {
		return false
	}
	// Check whether a --> closes that block before idx
	closeIdx := strings.Index(content[openIdx:], marker.CommentClose)
	if closeIdx == -1 {
		// Unclosed comment - treat as inside
		return true
	}
	// The comment closes at openIdx+closeIdx; if that position is >= idx,
	// the position is still inside the comment.
	return openIdx+closeIdx+len(marker.CommentClose) > idx
}

// Decision inserts a decision entry before existing entries.
//
// Finds the first "## [" marker that is NOT inside an HTML comment block
// and inserts before it, maintaining reverse-chronological order.
// Falls back to AfterHeader if no real entries exist yet.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Header line to insert after (e.g., "# Decisions")
//
// Returns:
//   - []byte: Modified content with entry inserted
func Decision(content, entry, header string) []byte {
	return insertBeforeFirstEntry(content, entry, header)
}

// Learning inserts a learning entry before existing entries.
//
// Finds the first "## [" marker that is NOT inside an HTML comment block
// and inserts before it, maintaining reverse-chronological order.
// Falls back to AfterHeader if no real entries exist yet.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//
// Returns:
//   - []byte: Modified content with entry inserted
func Learning(content, entry string) []byte {
	return insertBeforeFirstEntry(
		content, entry, desc.Text(text.DescKeyHeadingLearnings),
	)
}
