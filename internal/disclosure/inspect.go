//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

// StagedEntries returns the un-digested entries in a root's staging zone,
// in file order, per the root's own kind: timestamped "## [" entries for
// LEARNINGS/DECISIONS, curated "## " prose sections for CONVENTIONS. A
// convention's Timestamp is always empty — its identity is its title.
//
// Parameters:
//   - root: the parsed root
//
// Returns:
//   - []StagedEntry: the staged entries, nil when staging holds none
func StagedEntries(root Root) []StagedEntry {
	blocks := stagedBlocks(root.Staging, root.Kind)
	if len(blocks) == 0 {
		return nil
	}
	entries := make([]StagedEntry, 0, len(blocks))
	for _, b := range blocks {
		entries = append(entries, StagedEntry{
			Timestamp: b.Timestamp,
			Title:     b.Title,
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
