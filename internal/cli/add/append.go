//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import "strings"

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
	if fileType == "task" || fileType == "tasks" {
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
			lineEnd := strings.Index(existingStr[idx:], "\n")
			if lineEnd != -1 {
				insertPoint := idx + lineEnd + 1
				return []byte(existingStr[:insertPoint] + "\n" +
					entry + existingStr[insertPoint:])
			}
		}
	}

	// For decisions, prepend after the "# Decisions" header for reverse-chronological order
	if fileType == "decision" || fileType == "decisions" {
		return prependAfterHeader(existingStr, entry, "# Decisions")
	}

	// For learnings, prepend after the header section (after the first "---")
	if fileType == "learning" || fileType == "learnings" {
		return prependAfterSeparator(existingStr, entry)
	}

	// Default (conventions): append at the end
	if !strings.HasSuffix(existingStr, "\n") {
		existingStr += "\n"
	}
	return []byte(existingStr + "\n" + entry)
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
		return []byte(content[:entryIdx] + entry + "\n---\n\n" + content[entryIdx:])
	}

	// No existing entries - find header and insert after it
	idx := strings.Index(content, header)
	if idx != -1 {
		lineEnd := strings.Index(content[idx:], "\n")
		if lineEnd != -1 {
			insertPoint := idx + lineEnd + 1
			// Skip blank lines and comments
			for insertPoint < len(content) {
				if content[insertPoint] == '\n' {
					insertPoint++
				} else if insertPoint+4 < len(content) && content[insertPoint:insertPoint+4] == "<!--" {
					// Skip HTML comment
					endComment := strings.Index(content[insertPoint:], "-->")
					if endComment != -1 {
						insertPoint += endComment + 3
						// Skip trailing whitespace after comment
						for insertPoint < len(content) && (content[insertPoint] == '\n' || content[insertPoint] == ' ') {
							insertPoint++
						}
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
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return []byte(content + "\n" + entry)
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
		return []byte(content[:entryIdx] + entry + "\n" + content[entryIdx:])
	}

	// Also check for section-style learnings "## ["
	if entryIdx = strings.Index(content, "## ["); entryIdx != -1 {
		return []byte(content[:entryIdx] + entry + "\n---\n\n" + content[entryIdx:])
	}

	// No existing entries - find header and insert after it
	idx := strings.Index(content, "# Learnings")
	if idx != -1 {
		lineEnd := strings.Index(content[idx:], "\n")
		if lineEnd != -1 {
			insertPoint := idx + lineEnd + 1
			// Skip blank lines and comments
			for insertPoint < len(content) {
				if content[insertPoint] == '\n' {
					insertPoint++
				} else if insertPoint+4 < len(content) && content[insertPoint:insertPoint+4] == "<!--" {
					// Skip HTML comment
					endComment := strings.Index(content[insertPoint:], "-->")
					if endComment != -1 {
						insertPoint += endComment + 3
						// Skip trailing whitespace after comment
						for insertPoint < len(content) && (content[insertPoint] == '\n' || content[insertPoint] == ' ') {
							insertPoint++
						}
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

	// Final fallback: append at the end
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return []byte(content + "\n" + entry)
}
