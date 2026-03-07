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

// insertBeforeFirstEntry scans for the first "## [" marker not inside an
// HTML comment and inserts the entry before it. Falls back to
// InsertAfterHeader when no real entries exist yet.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted entry to insert
//   - header: Section header to fall back to
//
// Returns:
//   - []byte: Modified content with entry inserted
func insertBeforeFirstEntry(content, entry, header string) []byte {
	search := content
	offset := 0
	for {
		rel := strings.Index(search, config.HeadingLearningStart)
		if rel == -1 {
			break
		}
		entryIdx := offset + rel
		if !IsInsideHTMLComment(content, entryIdx) {
			return []byte(
				content[:entryIdx] + entry +
					config.NewlineLF + config.Separator +
					config.NewlineLF + config.NewlineLF +
					content[entryIdx:],
			)
		}
		offset = entryIdx + len(config.HeadingLearningStart)
		search = content[offset:]
	}

	return InsertAfterHeader(content, entry, header)
}
