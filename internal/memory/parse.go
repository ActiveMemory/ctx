//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// EntryKind identifies how an entry was delimited in MEMORY.md.
type EntryKind int

const (
	// EntryHeader is a Markdown heading (## or ###).
	EntryHeader EntryKind = iota
	// EntryParagraph is a blank-line-separated paragraph.
	EntryParagraph
	// EntryList is one or more consecutive list items.
	EntryList
)

// Entry is a discrete block parsed from MEMORY.md.
type Entry struct {
	Text      string    // Raw text of the entry (trimmed)
	StartLine int       // 1-based line number where the entry begins
	Kind      EntryKind // How the entry was delimited
}

// ParseEntries splits MEMORY.md content into discrete entries.
//
// Entry boundaries:
//   - Markdown headers (## or ###) start a new entry
//   - Blank lines separate paragraphs into distinct entries
//   - Consecutive list items (- or *) are grouped into a single entry
//
// The top-level heading (# Title) is skipped as it's structural, not content.
func ParseEntries(content string) []Entry {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	lines := strings.Split(content, config.NewlineLF)
	var entries []Entry
	var current []string
	var currentKind EntryKind
	currentStart := 0
	inEntry := false

	flush := func() {
		text := strings.TrimSpace(strings.Join(current, config.NewlineLF))
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

		// Skip top-level heading
		if strings.HasPrefix(trimmed, "# ") && !strings.HasPrefix(trimmed, "## ") {
			if inEntry {
				flush()
			}
			continue
		}

		// Section header (## or ###) starts a new entry
		if strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "### ") {
			if inEntry {
				flush()
			}
			currentStart = lineNum
			currentKind = EntryHeader
			current = []string{line}
			inEntry = true
			continue
		}

		// Blank line
		if trimmed == "" {
			if inEntry && currentKind != EntryHeader {
				flush()
			}
			continue
		}

		// List item — each top-level item is a separate entry for classification
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if inEntry {
				flush()
			}
			currentStart = lineNum
			currentKind = EntryList
			current = []string{line}
			inEntry = true
			continue
		}

		// Regular text — part of a paragraph or continuation of a header block
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
