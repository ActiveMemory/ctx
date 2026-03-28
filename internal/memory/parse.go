//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Entries splits MEMORY.md content into discrete entries.
//
// Entry boundaries:
//   - Markdown headers (## or ###) start a new entry
//   - Blank lines separate paragraphs into distinct entries
//   - Consecutive list items (- or *) are grouped into a single entry
//
// The top-level heading (# Title) is skipped as it's structural, not content.
//
// Parameters:
//   - content: Raw markdown content to parse into entries
//
// Returns:
//   - []Entry: Parsed entries grouped by headers, paragraphs, and list blocks
func Entries(content string) []Entry {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	lines := strings.Split(content, token.NewlineLF)
	var entries []Entry
	var current []string
	var currentKind EntryKind
	currentStart := 0
	inEntry := false

	flush := func() {
		text := strings.TrimSpace(strings.Join(current, token.NewlineLF))
		if text != "" {
			entries = append(entries, Entry{
				Text:      text,
				StartLine: currentStart,
				Kind:      currentKind,
			})
		}
		current = nil
		inEntry = false
	}

	for i, line := range lines {
		lineNum := i + 1 // 1-based
		trimmed := strings.TrimSpace(line)

		// Section header (H2+) starts a new entry
		if strings.HasPrefix(trimmed, token.HeadingLevelTwoStart) ||
			strings.HasPrefix(trimmed, token.HeadingLevelThreeStart) {
			if inEntry {
				flush()
			}
			currentStart = lineNum
			currentKind = EntryHeader
			current = []string{line}
			inEntry = true
			continue
		}

		// Skip top-level heading (H2+ already handled above)
		if strings.HasPrefix(trimmed, token.HeadingLevelOneStart) {
			if inEntry {
				flush()
			}
			continue
		}

		// Blank line
		if trimmed == "" {
			if inEntry && currentKind != EntryHeader {
				flush()
			}
			continue
		}

		// List item: each top-level item is a separate entry for classification
		if strings.HasPrefix(trimmed, token.PrefixListDash) ||
			strings.HasPrefix(trimmed, token.PrefixListStar) {
			if inEntry {
				flush()
			}
			currentStart = lineNum
			currentKind = EntryList
			current = []string{line}
			inEntry = true
			continue
		}

		// Regular text: part of a paragraph or continuation of a header block
		if inEntry && (currentKind == EntryHeader || currentKind == EntryParagraph) {
			current = append(current, line)
			continue
		}
		if inEntry {
			flush()
		}
		if !inEntry {
			currentStart = lineNum
			currentKind = EntryParagraph
			current = []string{line}
			inEntry = true
		}
	}

	if inEntry {
		flush()
	}

	return entries
}
