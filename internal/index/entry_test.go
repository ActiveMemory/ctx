//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"strings"
	"testing"
	"time"
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

func TestEntryBlock_OlderThan(t *testing.T) {
	// An entry from 100 days ago
	oldDate := time.Now().AddDate(0, 0, -100).Format("2006-01-02")
	oldBlock := &EntryBlock{
		Entry: Entry{Date: oldDate},
	}

	if !oldBlock.OlderThan(90) {
		t.Error("100-day-old entry should be older than 90 days")
	}
	if oldBlock.OlderThan(110) {
		t.Error("100-day-old entry should not be older than 110 days")
	}

	// An entry from today
	todayBlock := &EntryBlock{
		Entry: Entry{Date: time.Now().Format("2006-01-02")},
	}

	if todayBlock.OlderThan(1) {
		t.Error("today's entry should not be older than 1 day")
	}
}

func TestEntryBlock_OlderThan_InvalidDate(t *testing.T) {
	eb := &EntryBlock{
		Entry: Entry{Date: "invalid"},
	}
	if eb.OlderThan(1) {
		t.Error("invalid date should return false")
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

func TestRemoveEntryBlocks_Empty(t *testing.T) {
	content := "# Decisions\n\nSome text.\n"
	result := RemoveEntryBlocks(content, nil)
	if result != content {
		t.Errorf("RemoveEntryBlocks with nil blocks should return original content")
	}
}

func TestRemoveEntryBlocks_RemoveOne(t *testing.T) {
	content := `# Decisions

## [2026-01-15-120000] First

Body of first.

## [2026-02-01-090000] Second

Body of second.
`
	blocks := ParseEntryBlocks(content)
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	// Remove the first block
	result := RemoveEntryBlocks(content, blocks[:1])
	if strings.Contains(result, "First") {
		t.Error("removed block should not appear in result")
	}
	if !strings.Contains(result, "Second") {
		t.Error("remaining block should appear in result")
	}
}

func TestRemoveEntryBlocks_RemoveAll(t *testing.T) {
	content := `# Decisions

## [2026-01-15-120000] First

Body of first.

## [2026-02-01-090000] Second

Body of second.
`
	blocks := ParseEntryBlocks(content)

	result := RemoveEntryBlocks(content, blocks)
	if strings.Contains(result, "First") || strings.Contains(result, "Second") {
		t.Error("all blocks should be removed")
	}
	if !strings.Contains(result, "# Decisions") {
		t.Error("file header should be preserved")
	}
}

func TestRemoveEntryBlocks_CleansBlankLines(t *testing.T) {
	content := `# Decisions

## [2026-01-15-120000] First

Body.



## [2026-02-01-090000] Second

Body.
`
	blocks := ParseEntryBlocks(content)
	result := RemoveEntryBlocks(content, blocks[:1])

	// Should not have 3+ consecutive blank lines
	if strings.Contains(result, "\n\n\n\n") {
		t.Error("result should not have 4+ consecutive blank lines")
	}
}
