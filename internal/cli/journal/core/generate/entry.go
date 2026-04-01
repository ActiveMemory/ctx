//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package generate

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/format"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// formatIndexEntry formats a single entry for the index.
//
// Parameters:
//   - e: Journal entry to format
//   - nl: Newline string
//
// Returns:
//   - string: Formatted line
//     (e.g., "- 14:30 [title](link.md) (project) `1.2KB`")
func formatIndexEntry(e entity.JournalEntry, nl string) string {
	link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)

	timeStr := ""
	if e.Time != "" && len(e.Time) >= journal.TimePrefixLen {
		timeStr = e.Time[:journal.TimePrefixLen] + " "
	}

	project := ""
	if e.Project != "" {
		project = fmt.Sprintf(desc.Text(text.DescKeyJournalProjectLabel), e.Project)
	}

	size := format.Size(e.Size)

	line := fmt.Sprintf(
		tpl.JournalIndexEntry+nl, timeStr, e.Title, link, project, size,
	)
	if e.Summary != "" {
		line += fmt.Sprintf(tpl.JournalIndexSummary+nl, e.Summary)
	}
	return line
}
