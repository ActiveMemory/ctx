//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package generate

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/entity"
)

func TestGenerateIndex(t *testing.T) {
	entries := []entity.JournalEntry{
		{
			Filename: "2026-01-21-session-one-abc12345.md",
			Title:    "Session One",
			Date:     "2026-01-21",
			Time:     "14:30:00",
			Project:  "project-a",
			Size:     1024,
		},
		{
			Filename: "2026-01-20-session-two-def67890.md",
			Title:    "Session Two",
			Date:     "2026-01-20",
			Time:     "10:00:00",
			Project:  "project-b",
			Size:     2048,
		},
		{
			Filename:   "2026-01-19-suggestion-ghi11111.md",
			Title:      "Suggestion",
			Date:       "2026-01-19",
			Time:       "09:00:00",
			Suggestive: true,
			Size:       512,
		},
	}

	index := Index(entries)

	if !strings.Contains(index, "# Session Journal") {
		t.Error("index missing header")
	}
	if !strings.Contains(index, "**Sessions**: 2") {
		t.Error("index missing session count")
	}
	if !strings.Contains(index, "**Suggestions**: 1") {
		t.Error("index missing suggestions count")
	}
	if !strings.Contains(index, "## 2026-01") {
		t.Error("index missing month header")
	}
	if !strings.Contains(index, "[Session One]") {
		t.Error("index missing session one link")
	}
	if !strings.Contains(index, "## Suggestions") {
		t.Error("index missing suggestions section")
	}
}

func TestInjectSourceLink_WithFrontmatter(t *testing.T) {
	content := "---\ntitle: Test\n---\n\n# Heading\n"
	result := InjectedSourceLink(content, "/home/user/.context/journal/test.md")

	wantLink := "[View source]" +
		"(file:///home/user/.context/journal/test.md)"
	if !strings.Contains(result, wantLink) {
		t.Errorf("missing file:// link:\n%s", result)
	}
	if !strings.Contains(result, ".context/journal/test.md") {
		t.Errorf("missing relative path:\n%s", result)
	}
	if !strings.Contains(result, "# Heading") {
		t.Error("original content missing")
	}
}

func TestInjectSourceLink_NoFrontmatter(t *testing.T) {
	content := "# Heading\n\nSome text.\n"
	result := InjectedSourceLink(content, "/path/to/file.md")

	if !strings.HasPrefix(result, "*[View source](file:///path/to/file.md)") {
		t.Errorf("source link not at top:\n%s", result)
	}
	if !strings.Contains(result, ".context/journal/file.md") {
		t.Errorf("missing relative path:\n%s", result)
	}
	if !strings.Contains(result, "# Heading") {
		t.Error("original content missing")
	}
}

func TestFormatIndexEntry(t *testing.T) {
	tests := []struct {
		name  string
		entry entity.JournalEntry
		check func(t *testing.T, got string)
	}{
		{
			name: "full entry",
			entry: entity.JournalEntry{
				Filename: "2026-01-21-session-abc12345.md",
				Title:    "Test Session",
				Time:     "14:30:00",
				Project:  "ctx",
				Size:     1536,
			},
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "14:30") {
					t.Error("missing time")
				}
				if !strings.Contains(got, "[Test Session]") {
					t.Error("missing title link")
				}
				if !strings.Contains(got, "(ctx)") {
					t.Error("missing project")
				}
				if !strings.Contains(got, "1.5KB") {
					t.Error("missing size")
				}
			},
		},
		{
			name: "no time or project",
			entry: entity.JournalEntry{
				Filename: "2026-01-21-session-abc12345.md",
				Title:    "Minimal",
				Size:     100,
			},
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "[Minimal]") {
					t.Error("missing title")
				}
				if !strings.Contains(got, "100B") {
					t.Error("missing size")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatIndexEntry(tt.entry, "\n")
			tt.check(t, got)
		})
	}
}

func TestGenerateZensicalToml_WithAllNav(t *testing.T) {
	entries := []entity.JournalEntry{
		{Filename: "a.md", Title: "A", Date: "2026-01-01"},
	}
	topics := []entity.TopicData{{Name: "t", Entries: entries}}
	keyFiles := []entity.KeyFileData{{Path: "f.go", Entries: entries}}
	sessionTypes := []entity.TypeData{{Name: "feature", Entries: entries}}

	got := ZensicalToml(entries, topics, keyFiles, sessionTypes)

	if !strings.Contains(got, "Topics") {
		t.Error("missing Topics nav")
	}
	if !strings.Contains(got, "Files") {
		t.Error("missing Files nav")
	}
	if !strings.Contains(got, "Types") {
		t.Error("missing Types nav")
	}
}

func TestGenerateZensicalToml_NoTopics(t *testing.T) {
	entries := []entity.JournalEntry{
		{Filename: "a.md", Title: "A", Date: "2026-01-01"},
	}

	got := ZensicalToml(entries, nil, nil, nil)

	if strings.Contains(got, "Topics") {
		t.Error("should not have Topics nav when empty")
	}
}
