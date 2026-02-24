//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"strings"
	"testing"
)

func TestParseEntryBlocks_Empty(t *testing.T) {
	blocks := ParseEntryBlocks("")
	if len(blocks) != 0 {
		t.Errorf("ParseEntryBlocks(\"\") = %d blocks, want 0", len(blocks))
	}
}

func TestParseEntryBlocks_NoEntries(t *testing.T) {
	content := "# Decisions\n\nSome intro text.\n"
	blocks := ParseEntryBlocks(content)
	if len(blocks) != 0 {
		t.Errorf("ParseEntryBlocks() = %d blocks, want 0", len(blocks))
	}
}

func TestParseEntryBlocks_Single(t *testing.T) {
	content := `# Decisions

## [2026-01-15-120000] Use YAML for config

**Context:** Need a config format
**Rationale:** YAML is human-readable
`
	blocks := ParseEntryBlocks(content)
	if len(blocks) != 1 {
		t.Fatalf("ParseEntryBlocks() = %d blocks, want 1", len(blocks))
	}

	b := blocks[0]
	if b.Entry.Date != "2026-01-15" {
		t.Errorf("Date = %q, want %q", b.Entry.Date, "2026-01-15")
	}
	if b.Entry.Title != "Use YAML for config" {
		t.Errorf("Title = %q, want %q", b.Entry.Title, "Use YAML for config")
	}
	if b.Entry.Timestamp != "2026-01-15-120000" {
		t.Errorf("Timestamp = %q, want %q", b.Entry.Timestamp, "2026-01-15-120000")
	}
	if len(b.Lines) != 4 {
		t.Errorf("Lines count = %d, want 4", len(b.Lines))
	}
}

func TestParseEntryBlocks_Multiple(t *testing.T) {
	content := `# Decisions

## [2026-01-15-120000] First decision

Body of first.

## [2026-02-01-090000] Second decision

Body of second.
`
	blocks := ParseEntryBlocks(content)
	if len(blocks) != 2 {
		t.Fatalf("ParseEntryBlocks() = %d blocks, want 2", len(blocks))
	}

	if blocks[0].Entry.Title != "First decision" {
		t.Errorf("blocks[0].Title = %q, want %q", blocks[0].Entry.Title, "First decision")
	}
	if blocks[1].Entry.Title != "Second decision" {
		t.Errorf("blocks[1].Title = %q, want %q", blocks[1].Entry.Title, "Second decision")
	}
}

func TestParseEntryBlocks_IndexMarkers(t *testing.T) {
	content := `# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|----------|
| 2026-01-15 | First |
<!-- INDEX:END -->

## [2026-01-15-120000] First

Body.
`
	blocks := ParseEntryBlocks(content)
	if len(blocks) != 1 {
		t.Fatalf("ParseEntryBlocks() = %d blocks, want 1", len(blocks))
	}
	if blocks[0].Entry.Title != "First" {
		t.Errorf("Title = %q, want %q", blocks[0].Entry.Title, "First")
	}
}

func TestEntryBlock_IsSuperseded(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			name:  "not superseded",
			lines: []string{"## [2026-01-15-120000] Test", "Body text"},
			want:  false,
		},
		{
			name:  "superseded",
			lines: []string{"## [2026-01-15-120000] Test", "~~Superseded by newer decision~~"},
			want:  true,
		},
		{
			name:  "superseded with leading space",
			lines: []string{"## [2026-01-15-120000] Test", "  ~~Superseded by newer~~"},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eb := &EntryBlock{Lines: tt.lines}
			if got := eb.IsSuperseded(); got != tt.want {
				t.Errorf("IsSuperseded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryBlock_BlockContent(t *testing.T) {
	eb := &EntryBlock{
		Lines: []string{
			"## [2026-01-15-120000] Test",
			"",
			"Body text here.",
		},
	}

	content := eb.BlockContent()
	if !strings.Contains(content, "## [2026-01-15-120000] Test") {
		t.Error("BlockContent should contain the header")
	}
	if !strings.Contains(content, "Body text here.") {
		t.Error("BlockContent should contain the body")
	}
}
