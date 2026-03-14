//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
)

func TestGenerateHomeMOC(t *testing.T) {
	entries := []JournalEntry{
		{
			Filename: "2026-02-14-session-a.md",
			Title:    "Session A",
			Type:     "feature",
			Outcome:  "completed",
		},
		{
			Filename: "2026-02-13-session-b.md",
			Title:    "Session B",
			Type:     "bugfix",
		},
	}

	got := GenerateHomeMOC(entries, true, true, true)

	if !strings.Contains(got, "# Session Journal") {
		t.Error("missing main heading")
	}
	if !strings.Contains(got, "[[_Topics|Topics]]") {
		t.Error("missing topics MOC link")
	}
	if !strings.Contains(got, "[[_Key Files|Key Files]]") {
		t.Error("missing files MOC link")
	}
	if !strings.Contains(got, "[[_Session Types|Session Types]]") {
		t.Error("missing types MOC link")
	}
	if !strings.Contains(got, "[[2026-02-14-session-a|Session A]]") {
		t.Error("missing entry wikilink")
	}
}

func TestGenerateHomeMOCNoSections(t *testing.T) {
	entries := []JournalEntry{
		{Filename: "entry.md", Title: "Test"},
	}

	got := GenerateHomeMOC(entries, false, false, false)

	if strings.Contains(got, "[[_Topics") {
		t.Error("should not have topics link when hasTopics=false")
	}
}

func TestGenerateObsidianTopicsMOC(t *testing.T) {
	topics := []TopicData{
		{
			Name:    "caching",
			Popular: true,
			Entries: []JournalEntry{
				{Filename: "a.md", Title: "A"},
				{Filename: "b.md", Title: "B"},
			},
		},
		{
			Name:    "auth",
			Popular: false,
			Entries: []JournalEntry{
				{Filename: "c.md", Title: "C"},
			},
		},
	}

	got := GenerateObsidianTopicsMOC(topics)

	if !strings.Contains(got, "[[caching]]") {
		t.Error("missing popular topic wikilink")
	}
	if !strings.Contains(got, "**auth**") {
		t.Error("missing longtail topic")
	}
	if !strings.Contains(got, "[[c|C]]") {
		t.Error("missing longtail entry wikilink")
	}
}

func TestGenerateRelatedFooter(t *testing.T) {
	entry := JournalEntry{
		Filename: "2026-02-14-main.md",
		Title:    "Main Entry",
		Type:     "feature",
		Topics:   []string{"caching", "auth"},
	}

	topicIndex := map[string][]JournalEntry{
		"caching": {
			entry,
			{Filename: "2026-02-13-related.md", Title: "Related Entry"},
		},
		"auth": {
			entry,
			{Filename: "2026-02-13-related.md", Title: "Related Entry"},
			{Filename: "2026-02-12-other.md", Title: "Other Entry"},
		},
	}

	got := GenerateRelatedFooter(entry, topicIndex, 5)

	if !strings.Contains(got, assets.ObsidianRelatedHeading) {
		t.Error("missing related heading")
	}
	if !strings.Contains(got, "[[_Topics|Topics MOC]]") {
		t.Error("missing topics MOC link")
	}
	if !strings.Contains(got, "[[caching]]") {
		t.Error("missing topic link")
	}
	if !strings.Contains(got, "[[feature]]") {
		t.Error("missing type link")
	}
	if !strings.Contains(got, "[[2026-02-13-related|Related Entry]]") {
		t.Error("missing related entry link")
	}
	if strings.Contains(got, "[[2026-02-14-main|Main Entry]]") {
		t.Error("entry should not link to itself")
	}
}

func TestGenerateRelatedFooterEmpty(t *testing.T) {
	entry := JournalEntry{
		Filename: "entry.md",
		Title:    "No Metadata",
	}

	got := GenerateRelatedFooter(entry, nil, 5)
	if got != "" {
		t.Errorf("expected empty footer for entry without metadata, got: %q", got)
	}
}

func TestCollectRelated(t *testing.T) {
	main := JournalEntry{
		Filename: "main.md",
		Title:    "Main",
		Topics:   []string{"a", "b"},
	}

	topicIndex := map[string][]JournalEntry{
		"a": {
			main,
			{Filename: "shared-ab.md", Title: "Shared AB"},
			{Filename: "only-a.md", Title: "Only A"},
		},
		"b": {
			main,
			{Filename: "shared-ab.md", Title: "Shared AB"},
			{Filename: "only-b.md", Title: "Only B"},
		},
	}

	related := CollectRelated(main, topicIndex, 10)

	if len(related) != 3 {
		t.Fatalf("expected 3 related entries, got %d", len(related))
	}

	if related[0].Filename != "shared-ab.md" {
		t.Errorf("expected shared-ab first (highest score), got %s",
			related[0].Filename)
	}
}

func TestCollectRelatedMaxLimit(t *testing.T) {
	main := JournalEntry{
		Filename: "main.md",
		Topics:   []string{"a"},
	}

	topicIndex := map[string][]JournalEntry{
		"a": {
			main,
			{Filename: "1.md", Title: "1"},
			{Filename: "2.md", Title: "2"},
			{Filename: "3.md", Title: "3"},
		},
	}

	related := CollectRelated(main, topicIndex, 2)
	if len(related) != 2 {
		t.Errorf("expected 2 entries (maxRelated=2), got %d", len(related))
	}
}

func TestFilterFunctions(t *testing.T) {
	entries := []JournalEntry{
		{Filename: "regular.md", Title: "Regular", Type: "feature",
			Topics: []string{"a"}, KeyFiles: []string{"b.go"}},
		{Filename: "suggestion.md", Suggestive: true},
		{Filename: "multipart-p2.md", Title: "Part 2"},
		{Filename: "no-meta.md", Title: "No Meta"},
	}

	regular := FilterRegularEntries(entries)
	if len(regular) != 2 {
		t.Errorf("FilterRegularEntries: expected 2, got %d", len(regular))
	}

	withTopics := FilterEntriesWithTopics(entries)
	if len(withTopics) != 1 {
		t.Errorf("FilterEntriesWithTopics: expected 1, got %d", len(withTopics))
	}

	withFiles := FilterEntriesWithKeyFiles(entries)
	if len(withFiles) != 1 {
		t.Errorf("FilterEntriesWithKeyFiles: expected 1, got %d", len(withFiles))
	}

	withType := FilterEntriesWithType(entries)
	if len(withType) != 1 {
		t.Errorf("FilterEntriesWithType: expected 1, got %d", len(withType))
	}
}

func TestBuildTopicLookup(t *testing.T) {
	entries := []JournalEntry{
		{Filename: "a.md", Topics: []string{"go", "testing"}},
		{Filename: "b.md", Topics: []string{"go", "cli"}},
	}

	lookup := BuildTopicLookup(entries)

	if len(lookup["go"]) != 2 {
		t.Errorf("expected 2 entries for 'go', got %d", len(lookup["go"]))
	}
	if len(lookup["testing"]) != 1 {
		t.Errorf("expected 1 entry for 'testing', got %d", len(lookup["testing"]))
	}
}
