//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heading

import (
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// ParseHeaders extracts all timestamped entry headers from file content.
//
// It scans for headers matching the pattern "## [YYYY-MM-DD-HHMMSS] Title"
// and returns them in the order they appear in the file. Unlike Headings,
// which projects any ATX heading, ParseHeaders is entry-specific: it is the
// recognizer that ctx agent uses to enumerate DECISIONS/LEARNINGS entries.
//
// Parameters:
//   - content: The full content of a context file
//
// Returns:
//   - []entity.IndexEntry: Slice of parsed entries (it may be empty)
func ParseHeaders(content string) []entity.IndexEntry {
	var entries []entity.IndexEntry

	matches := regex.EntryHeader.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == regex.EntryHeaderGroups {
			date := match[1]
			time := match[2]
			title := match[3]
			entries = append(entries, entity.IndexEntry{
				Timestamp: date + token.Dash + time,
				Date:      date,
				Title:     title,
			})
		}
	}

	return entries
}
