//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package section

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/entity"
)

func TestBuildTopicIndex(t *testing.T) {
	entries := []entity.JournalEntry{
		{Filename: "a.md", Title: "A", Topics: []string{"go", "testing"}},
		{Filename: "b.md", Title: "B", Topics: []string{"go", "cli"}},
		{Filename: "c.md", Title: "C", Topics: []string{"python"}},
	}

	topics := BuildTopicIndex(entries)

	if len(topics) != 4 {
		t.Fatalf("expected 4 topics, got %d", len(topics))
	}

	if topics[0].Name != "go" {
		t.Errorf("first topic should be 'go' (2 entries), got %q", topics[0].Name)
	}
	if !topics[0].Popular {
		t.Error("'go' should be popular (2+ entries)")
	}

	for _, tp := range topics[1:] {
		if tp.Popular {
			t.Errorf("topic %q should not be popular (1 entry)", tp.Name)
		}
	}
}

func TestGenerateTopicsIndex(t *testing.T) {
	topics := []entity.TopicData{
		{
			Name:    "caching",
			Popular: true,
			Entries: []entity.JournalEntry{
				{Filename: "a.md", Title: "A"},
				{Filename: "b.md", Title: "B"},
			},
		},
		{
			Name:    "auth",
			Popular: false,
			Entries: []entity.JournalEntry{
				{Filename: "c.md", Title: "C"},
			},
		},
	}

	got := GenerateTopicsIndex(topics)

	if !strings.Contains(got, "# Topics") {
		t.Error("missing heading")
	}
	if !strings.Contains(got, "[caching]") {
		t.Error("missing popular topic link")
	}
	if !strings.Contains(got, "auth") {
		t.Error("missing longtail topic")
	}
}

func TestGenerateTopicPage(t *testing.T) {
	topic := entity.TopicData{
		Name: "caching",
		Entries: []entity.JournalEntry{
			{
				Filename: "2026-02-14-a.md", Title: "A",
				Date: "2026-02-14", Time: "10:00:00",
			},
			{
				Filename: "2026-01-20-b.md", Title: "B",
				Date: "2026-01-20", Time: "09:00:00",
			},
		},
	}

	got := GenerateTopicPage(topic)

	if !strings.Contains(got, "# caching") {
		t.Error("missing heading")
	}
	if !strings.Contains(got, "2 sessions") {
		t.Error("missing session count")
	}
	if !strings.Contains(got, "[A]") {
		t.Error("missing entry link")
	}
}

func TestBuildKeyFileIndex(t *testing.T) {
	entries := []entity.JournalEntry{
		{Filename: "a.md", KeyFiles: []string{"cmd/main.go", "internal/config.go"}},
		{Filename: "b.md", KeyFiles: []string{"cmd/main.go"}},
		{Filename: "c.md", KeyFiles: []string{"README.md"}},
	}

	keyFiles := BuildKeyFileIndex(entries)

	if len(keyFiles) != 3 {
		t.Fatalf("expected 3 key files, got %d", len(keyFiles))
	}
	if keyFiles[0].Path != "cmd/main.go" {
		t.Errorf("first key file should be cmd/main.go, got %q", keyFiles[0].Path)
	}
	if !keyFiles[0].Popular {
		t.Error("cmd/main.go should be popular (2+ entries)")
	}
}

func TestGenerateKeyFilesIndex(t *testing.T) {
	keyFiles := []entity.KeyFileData{
		{
			Path:    "cmd/main.go",
			Popular: true,
			Entries: []entity.JournalEntry{
				{Filename: "a.md", Title: "A"},
				{Filename: "b.md", Title: "B"},
			},
		},
		{
			Path:    "README.md",
			Popular: false,
			Entries: []entity.JournalEntry{
				{Filename: "c.md", Title: "C"},
			},
		},
	}

	got := GenerateKeyFilesIndex(keyFiles)

	if !strings.Contains(got, "# Key Files") {
		t.Error("missing heading")
	}
	if !strings.Contains(got, "cmd/main.go") {
		t.Error("missing key file path")
	}
}

func TestGenerateKeyFilePage(t *testing.T) {
	kf := entity.KeyFileData{
		Path: "internal/config.go",
		Entries: []entity.JournalEntry{
			{Filename: "2026-02-14-a.md", Title: "A", Date: "2026-02-14"},
		},
	}

	got := GenerateKeyFilePage(kf)

	if !strings.Contains(got, "internal/config.go") {
		t.Error("missing file path")
	}
	if !strings.Contains(got, "1 sessions") {
		t.Error("missing session count")
	}
}

func TestBuildTypeIndex(t *testing.T) {
	entries := []entity.JournalEntry{
		{Filename: "a.md", Type: "feature"},
		{Filename: "b.md", Type: "feature"},
		{Filename: "c.md", Type: "bugfix"},
	}

	types := BuildTypeIndex(entries)

	if len(types) != 2 {
		t.Fatalf("expected 2 types, got %d", len(types))
	}
	if types[0].Name != "feature" {
		t.Errorf("first type should be 'feature', got %q", types[0].Name)
	}
}

func TestGenerateTypesIndex(t *testing.T) {
	types := []entity.TypeData{
		{
			Name: "feature",
			Entries: []entity.JournalEntry{
				{Filename: "a.md"}, {Filename: "b.md"},
			},
		},
		{Name: "bugfix", Entries: []entity.JournalEntry{{Filename: "c.md"}}},
	}

	got := GenerateTypesIndex(types)

	if !strings.Contains(got, "# Session Types") {
		t.Error("missing heading")
	}
	if !strings.Contains(got, "2 types") {
		t.Error("missing type count")
	}
}

func TestGenerateTypePage(t *testing.T) {
	st := entity.TypeData{
		Name: "feature",
		Entries: []entity.JournalEntry{
			{Filename: "2026-02-14-a.md", Title: "A", Date: "2026-02-14"},
		},
	}

	got := GenerateTypePage(st)

	if !strings.Contains(got, "# feature") {
		t.Error("missing heading")
	}
	if !strings.Contains(got, "feature") {
		t.Error("missing type name")
	}
}
