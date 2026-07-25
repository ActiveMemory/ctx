//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package insert

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreEntry "github.com/ActiveMemory/ctx/internal/cli/add/core/entry"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// AppendEntry inserts a formatted entry into existing file content.
//
// For tasks, inserts after the target section header. For decisions and
// learnings, inserts before existing entries (reverse-chronological order).
// For conventions, inserts at the top of the target H2 section (section
// is honored here too, and the section is created when absent).
//
// Parameters:
//   - existing: Current file content as bytes
//   - entry: Pre-formatted entry text to insert
//   - fileType: AppendEntry type
//     (e.g., "task", "decision", "learning", "convention")
//   - section: Target section header for tasks; defaults to "## Next Up" if
//     empty; a "## " prefix is added automatically if missing
//
// Returns:
//   - []byte: Modified file content with the entry inserted
func AppendEntry(
	existing []byte, entry string, fileType string, section string,
) []byte {
	existingStr := string(existing)

	switch {
	// For tasks, find the appropriate section
	case coreEntry.FileTypeIsTask(fileType):
		return Task(entry, existingStr, section)
	// Decisions: insert before existing entries for reverse-chronological order
	case coreEntry.FileTypeIsDecision(fileType):
		return Decision(
			existingStr, entry, desc.Text(text.DescKeyHeadingDecisions),
		)
	// Learnings: insert before existing entries for reverse-chronological order
	case coreEntry.FileTypeIsLearning(fileType):
		return Learning(existingStr, entry)
	default:
		// Conventions: bullets live inside an H2 section, so the entry is
		// inserted there rather than appended at EOF (which lands below
		// ## Themes on a folded root).
		return Convention(existingStr, entry, section)
	}
}
