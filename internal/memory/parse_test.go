//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"testing"
)

func TestParseEntries_Headers(t *testing.T) {
	content := `# Memory

## Session 2026-03-05: Memory Bridge

Worked on the memory bridge foundation.

## Key decisions

Decided to use heuristic classification.
`
	entries := ParseEntries(content)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d: %+v", len(entries), entries)
	}
	if entries[0].Kind != EntryHeader {
		t.Errorf("entry 0: expected EntryHeader, got %d", entries[0].Kind)
	}
	if entries[0].StartLine != 3 {
		t.Errorf("entry 0: expected StartLine 3, got %d", entries[0].StartLine)
	}
	if entries[1].Kind != EntryHeader {
		t.Errorf("entry 1: expected EntryHeader, got %d", entries[1].Kind)
	}
}

func TestParseEntries_Paragraphs(t *testing.T) {
	content := `# Memory

First paragraph about something.

Second paragraph about something else.
`
	entries := ParseEntries(content)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d: %+v", len(entries), entries)
	}
	if entries[0].Kind != EntryParagraph {
		t.Errorf("entry 0: expected EntryParagraph, got %d", entries[0].Kind)
	}
	if entries[1].Kind != EntryParagraph {
		t.Errorf("entry 1: expected EntryParagraph, got %d", entries[1].Kind)
	}
}

func TestParseEntries_ListItems(t *testing.T) {
	content := `# Memory

- always use bun for this project
- prefer filepath.Join over string concat
- never use global state
`
	entries := ParseEntries(content)
	if len(entries) != 3 {
		t.Fatalf("expected 3 list entries (one per item), got %d: %+v", len(entries), entries)
	}
	for i, e := range entries {
		if e.Kind != EntryList {
			t.Errorf("entry %d: expected EntryList, got %d", i, e.Kind)
		}
	}
	if entries[0].StartLine != 3 {
		t.Errorf("expected StartLine 3, got %d", entries[0].StartLine)
	}
}

func TestParseEntries_Mixed(t *testing.T) {
	content := `# Auto Memory

## Session notes

Worked on import pipeline.

- always use ctx from PATH
- decided to use SQLite over Postgres

### Key learnings

Learned that golangci-lint v2 ignores inline nolint.

Some standalone paragraph.
`
	entries := ParseEntries(content)
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	// Header "## Session notes" absorbs its body paragraph
	if entries[0].Kind != EntryHeader {
		t.Errorf("entry 0: expected EntryHeader, got %d", entries[0].Kind)
	}

	// Two separate list items
	if entries[1].Kind != EntryList {
		t.Errorf("entry 1: expected EntryList, got %d", entries[1].Kind)
	}
	if entries[2].Kind != EntryList {
		t.Errorf("entry 2: expected EntryList, got %d", entries[2].Kind)
	}

	// Sub-header "### Key learnings" absorbs all following paragraphs until next header
	if entries[3].Kind != EntryHeader {
		t.Errorf("entry 3: expected EntryHeader, got %d", entries[3].Kind)
	}
}

func TestParseEntries_Empty(t *testing.T) {
	entries := ParseEntries("")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", len(entries))
	}

	entries = ParseEntries("   \n\n  ")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for whitespace-only input, got %d", len(entries))
	}
}

func TestParseEntries_TopLevelHeadingSkipped(t *testing.T) {
	content := `# Memory

Just a paragraph.
`
	entries := ParseEntries(content)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (top-level heading skipped), got %d", len(entries))
	}
	if entries[0].Kind != EntryParagraph {
		t.Errorf("expected EntryParagraph, got %d", entries[0].Kind)
	}
}

func TestParseEntries_IndividualListItems(t *testing.T) {
	content := `# Memory

- first item

- second item after blank line
`
	entries := ParseEntries(content)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries (blank line separates lists), got %d: %+v", len(entries), entries)
	}
}
