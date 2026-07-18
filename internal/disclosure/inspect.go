//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"github.com/ActiveMemory/ctx/internal/heading"
)

// StagedEntries returns the un-digested entries in a root's staging zone,
// in file order. Conventions have no "## [" entries, so this is empty for
// them (their digestion is a later milestone).
//
// Parameters:
//   - root: the parsed root
//
// Returns:
//   - []StagedEntry: the staged entries, nil when staging holds none
func StagedEntries(root Root) []StagedEntry {
	blocks := heading.ParseEntryBlocks(root.Staging)
	if len(blocks) == 0 {
		return nil
	}
	entries := make([]StagedEntry, 0, len(blocks))
	for _, b := range blocks {
		entries = append(entries, StagedEntry{
			Timestamp: b.Entry.Timestamp,
			Title:     b.Entry.Title,
		})
	}
	return entries
}

// Inspect parses content as a root of kind k and returns the read-only
// view the dry-run pass consumes: kind name, staged entries, and current
// themes. Like Parse, it is total.
//
// Parameters:
//   - content: the full root file content
//   - k: which canonical file this is
//
// Returns:
//   - Inspection: the structured read-only view
func Inspect(content string, k Kind) Inspection {
	root := Parse(content, k)
	return Inspection{
		Kind:    k.String(),
		Staging: StagedEntries(root),
		Themes:  root.Themes,
	}
}
