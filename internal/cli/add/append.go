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

// AppendEntry inserts a formatted entry into existing file content.
//
// For task entries, the function locates the target section header and inserts
// the entry immediately after it. For decisions and learnings, entries are
// prepended (inserted after the header section) for reverse-chronological order.
// For conventions, entries are appended to the end of the file.
//
// Parameters:
//   - existing: Current file content as bytes
//   - entry: Pre-formatted entry text to insert
//   - fileType: Entry type (e.g., "task", "decision", "learning", "convention")
//   - section: Target section header for tasks; defaults to "## Next Up" if
//     empty; a "## " prefix is added automatically if missing
//
// Returns:
//   - []byte: Modified file content with the entry inserted
func AppendEntry(
	existing []byte, entry string, fileType string, section string,
) []byte {
	existingStr := string(existing)

	// For tasks, find the appropriate section
	if fileType == config.UpdateTypeTask || fileType == config.UpdateTypeTasks {
		targetSection := section
		if targetSection == "" {
			targetSection = "## Next Up"
		} else if !strings.HasPrefix(targetSection, "##") {
			targetSection = "## " + targetSection
		}

		// Find the section and insert after it
		idx := strings.Index(existingStr, targetSection)
		if idx != -1 {
			// Find the end of the section header line
			lineEnd := findNewline(existingStr[idx:])
			if lineEnd != -1 {
				insertPoint := idx + lineEnd
				insertPoint = skipNewline(existingStr, insertPoint)
				return []byte(existingStr[:insertPoint] + config.NewlineLF +
					entry + existingStr[insertPoint:])
			}
		}
	}

	// For decisions, prepend after the "# Decisions" header for reverse-chronological order
	if fileType == config.UpdateTypeDecision || fileType == config.UpdateTypeDecisions {
		return prependAfterHeader(existingStr, entry, "# Decisions")
	}

	// For learnings, prepend after the header section (after the first "---")
	if fileType == config.UpdateTypeLearning || fileType == config.UpdateTypeLearnings {
		return prependAfterSeparator(existingStr, entry)
	}

	// Default (conventions): append at the end
	if !endsWithNewline(existingStr) {
		existingStr += config.NewlineLF
	}
	return []byte(existingStr + config.NewlineLF + entry)
}

// prependAfterHeader inserts an entry after a header line.
//
// Used for DECISIONS.md to maintain reverse-chronological order.
// Entries are inserted before any existing entries (identified by "## [").
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Header line to insert after (e.g., "# Decisions")
//
// Returns:
//   - []byte: Modified content with entry inserted
func prependAfterHeader(content, entry, header string) []byte {
	// Find the first entry marker "## [" (timestamp-prefixed sections)
	entryIdx := strings.Index(content, "## [")
	if entryIdx != -1 {
		// Insert before the first entry, with separator after
		return []byte(content[:entryIdx] + entry + config.NewlineLF + "---" + config.NewlineLF + config.NewlineLF + content[entryIdx:])
	}

	// No existing entries - find header and insert after it
	return insertAfterHeader(content, entry, header)
}

// prependAfterSeparator inserts an entry for learnings.
//
// Entries are inserted before any existing entries (identified by "- **[").
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//
// Returns:
//   - []byte: Modified content with entry inserted
func prependAfterSeparator(content, entry string) []byte {
	// Find the first entry marker "- **[" (timestamp-prefixed list items)
	entryIdx := strings.Index(content, "- **[")
	if entryIdx != -1 {
		// Insert before the first entry
		return []byte(content[:entryIdx] + entry + config.NewlineLF + content[entryIdx:])
	}

	// Also check for section-style learnings "## ["
	if entryIdx = strings.Index(content, "## ["); entryIdx != -1 {
		return []byte(content[:entryIdx] + entry + config.NewlineLF + "---" + config.NewlineLF + config.NewlineLF + content[entryIdx:])
	}

	// No existing entries - find header and insert after it
	return insertAfterHeader(content, entry, "# Learnings")
}

// insertAfterHeader finds a header line and inserts content after it.
//
// Skips blank lines and ctx markers between the header and insertion point.
// Falls back to appending at the end if header is not found.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Header line to find (e.g., "# Learnings")
//
// Returns:
//   - []byte: Modified content with entry inserted
func insertAfterHeader(content, entry, header string) []byte {
	idx := strings.Index(content, header)
	if idx != -1 {
		lineEnd := findNewline(content[idx:])
		if lineEnd != -1 {
			insertPoint := idx + lineEnd
			insertPoint = skipNewline(content, insertPoint)
			// Skip blank lines and ctx markers
			for insertPoint < len(content) {
				if n := skipNewline(content, insertPoint); n > insertPoint {
					insertPoint = n
				} else if insertPoint+len(config.CommentOpen) <= len(content) &&
					content[insertPoint:insertPoint+len(config.CommentOpen)] == config.CommentOpen {
					// Skip ctx marker
					endComment := strings.Index(content[insertPoint:], config.CommentClose)
					if endComment != -1 {
						insertPoint += endComment + len(config.CommentClose)
						// Skip trailing whitespace after marker
						insertPoint = skipWhitespace(content, insertPoint)
					} else {
						break
					}
				} else {
					break
				}
			}
			return []byte(content[:insertPoint] + entry)
		}
	}

	// Fallback: append at the end
	if !endsWithNewline(content) {
		content += config.NewlineLF
	}
	return []byte(content + config.NewlineLF + entry)
}

// findNewline returns the index of the first newline (CRLF or LF) in s.
// Returns -1 if no newline is found.
func findNewline(s string) int {
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) && s[i] == '\r' && s[i+1] == '\n' {
			return i
		}
		if s[i] == '\n' {
			return i
		}
	}
	return -1
}

// skipNewline advances pos past a newline (CRLF or LF) if present.
// Returns the new position (unchanged if no newline at pos).
func skipNewline(s string, pos int) int {
	if pos >= len(s) {
		return pos
	}
	if pos+1 < len(s) && s[pos] == '\r' && s[pos+1] == '\n' {
		return pos + 2
	}
	if s[pos] == '\n' {
		return pos + 1
	}
	return pos
}

// skipWhitespace advances pos past any whitespace (space, tab, newline).
func skipWhitespace(s string, pos int) int {
	for pos < len(s) {
		if n := skipNewline(s, pos); n > pos {
			pos = n
		} else if s[pos] == ' ' || s[pos] == '\t' {
			pos++
		} else {
			break
		}
	}
	return pos
}

// endsWithNewline reports whether s ends with a newline (CRLF or LF).
func endsWithNewline(s string) bool {
	return strings.HasSuffix(s, config.NewlineCRLF) || strings.HasSuffix(s, config.NewlineLF)
}
