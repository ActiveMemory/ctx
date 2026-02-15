//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestGenerateHomeMOC(t *testing.T) {
	entries := []journalEntry{
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

	got := generateHomeMOC(entries, true, true, true)

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
	entries := []journalEntry{
		{Filename: "entry.md", Title: "Test"},
	}

	got := generateHomeMOC(entries, false, false, false)

	if strings.Contains(got, "[[_Topics") {
		t.Error("should not have topics link when hasTopics=false")
	}
}

func TestGenerateObsidianTopicsMOC(t *testing.T) {
	topics := []topicData{
		{
			Name:    "caching",
			Popular: true,
			Entries: []journalEntry{
				{Filename: "a.md", Title: "A"},
				{Filename: "b.md", Title: "B"},
			},
		},
		{
			Name:    "auth",
			Popular: false,
			Entries: []journalEntry{
				{Filename: "c.md", Title: "C"},
			},
		},
	}

	got := generateObsidianTopicsMOC(topics)

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
	entry := journalEntry{
		Filename: "2026-02-14-main.md",
		Title:    "Main Entry",
		Type:     "feature",
		Topics:   []string{"caching", "auth"},
	}

	topicIndex := map[string][]journalEntry{
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

	got := generateRelatedFooter(entry, topicIndex, 5)

	if !strings.Contains(got, config.ObsidianRelatedHeading) {
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
	// The main entry should not link to itself
	if strings.Contains(got, "[[2026-02-14-main|Main Entry]]") {
		t.Error("entry should not link to itself")
	}
}

func TestGenerateRelatedFooterEmpty(t *testing.T) {
	entry := journalEntry{
		Filename: "entry.md",
		Title:    "No Metadata",
	}

	got := generateRelatedFooter(entry, nil, 5)
	if got != "" {
		t.Errorf("expected empty footer for entry without metadata, got: %q", got)
	}
}

func TestCollectRelated(t *testing.T) {
	main := journalEntry{
		Filename: "main.md",
		Title:    "Main",
		Topics:   []string{"a", "b"},
	}

	topicIndex := map[string][]journalEntry{
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

	related := collectRelated(main, topicIndex, 10)

	if len(related) != 3 {
		t.Fatalf("expected 3 related entries, got %d", len(related))
	}

	// shared-ab should be first (score 2: shared topics a + b)
	if related[0].Filename != "shared-ab.md" {
		t.Errorf("expected shared-ab first (highest score), got %s",
			related[0].Filename)
	}
}

func TestCollectRelatedMaxLimit(t *testing.T) {
	main := journalEntry{
		Filename: "main.md",
		Topics:   []string{"a"},
	}

	topicIndex := map[string][]journalEntry{
		"a": {
			main,
			{Filename: "1.md", Title: "1"},
			{Filename: "2.md", Title: "2"},
			{Filename: "3.md", Title: "3"},
		},
	}

	related := collectRelated(main, topicIndex, 2)
	if len(related) != 2 {
		t.Errorf("expected 2 entries (maxRelated=2), got %d", len(related))
	}
}

func TestFilterFunctions(t *testing.T) {
	entries := []journalEntry{
		{Filename: "regular.md", Title: "Regular", Type: "feature",
			Topics: []string{"a"}, KeyFiles: []string{"b.go"}},
		{Filename: "suggestion.md", Suggestive: true},
		{Filename: "multipart-p2.md", Title: "Part 2"},
		{Filename: "no-meta.md", Title: "No Meta"},
	}

	regular := filterRegularEntries(entries)
	if len(regular) != 2 {
		t.Errorf("filterRegularEntries: expected 2, got %d", len(regular))
	}

	withTopics := filterEntriesWithTopics(entries)
	if len(withTopics) != 1 {
		t.Errorf("filterEntriesWithTopics: expected 1, got %d", len(withTopics))
	}

	withFiles := filterEntriesWithKeyFiles(entries)
	if len(withFiles) != 1 {
		t.Errorf("filterEntriesWithKeyFiles: expected 1, got %d", len(withFiles))
	}

	withType := filterEntriesWithType(entries)
	if len(withType) != 1 {
		t.Errorf("filterEntriesWithType: expected 1, got %d", len(withType))
	}
}

func TestBuildTopicLookup(t *testing.T) {
	entries := []journalEntry{
		{Filename: "a.md", Topics: []string{"go", "testing"}},
		{Filename: "b.md", Topics: []string{"go", "cli"}},
	}

	lookup := buildTopicLookup(entries)

	if len(lookup["go"]) != 2 {
		t.Errorf("expected 2 entries for 'go', got %d", len(lookup["go"]))
	}
	if len(lookup["testing"]) != 1 {
		t.Errorf("expected 1 entry for 'testing', got %d", len(lookup["testing"]))
	}
}

func TestRunJournalObsidianIntegration(t *testing.T) {
	// Create a temporary journal directory with test entries
	tmpDir := t.TempDir()
	journalDir := filepath.Join(tmpDir, config.DirContext, config.DirJournal)
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	// Write test entries with frontmatter
	entry1 := `---
title: "Feature: Add caching"
date: 2026-02-14
type: feature
outcome: completed
topics:
  - caching
  - performance
key_files:
  - internal/cache/store.go
---

# Feature: Add caching

**Time**: 14:30:00
**Project**: ctx

## Summary

Added a caching layer.
`
	entry2 := `---
title: "Fix: Cache invalidation"
date: 2026-02-13
type: bugfix
outcome: completed
topics:
  - caching
  - debugging
---

# Fix: Cache invalidation

**Time**: 10:00:00
**Project**: ctx

## Summary

Fixed cache invalidation bug.
`
	entry3 := `# No frontmatter session

**Time**: 09:00:00
**Project**: ctx

Just a plain session without enrichment.
`

	entries := map[string]string{
		"2026-02-14-add-caching-abc12345.md":       entry1,
		"2026-02-13-fix-cache-def67890.md":         entry2,
		"2026-02-12-plain-session-ghi11111.md":     entry3,
	}

	for name, content := range entries {
		path := filepath.Join(journalDir, name)
		if err := os.WriteFile(path, []byte(content), config.PermFile); err != nil {
			t.Fatal(err)
		}
	}

	// Run the vault generation
	outputDir := filepath.Join(tmpDir, "vault-output")

	// Use a real cobra.Command with captured output
	cmd := &cobra.Command{}
	cmd.SetOut(&strings.Builder{})
	cmd.SetErr(&strings.Builder{})

	err := buildObsidianVault(cmd, journalDir, outputDir)
	if err != nil {
		t.Fatalf("runJournalObsidian failed: %v", err)
	}

	// Verify vault structure
	assertFileExists(t, filepath.Join(outputDir, config.ObsidianConfigDir, config.ObsidianAppConfigFile))
	assertFileExists(t, filepath.Join(outputDir, config.ObsidianHomeMOC))
	assertFileExists(t, filepath.Join(outputDir, config.FilenameReadme))

	// Verify entries were written
	assertFileExists(t, filepath.Join(outputDir, config.ObsidianDirEntries, "2026-02-14-add-caching-abc12345.md"))
	assertFileExists(t, filepath.Join(outputDir, config.ObsidianDirEntries, "2026-02-13-fix-cache-def67890.md"))

	// Verify .obsidian/app.json content
	appConfig, readErr := os.ReadFile(filepath.Join(
		outputDir, config.ObsidianConfigDir, config.ObsidianAppConfigFile)) //nolint:gosec // test file path
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(appConfig), `"useMarkdownLinks": false`) {
		t.Error("app.json missing useMarkdownLinks setting")
	}

	// Verify Home.md contains wikilinks
	home, readErr := os.ReadFile(filepath.Join(outputDir, config.ObsidianHomeMOC)) //nolint:gosec // test file path
	if readErr != nil {
		t.Fatal(readErr)
	}
	homeStr := string(home)
	if !strings.Contains(homeStr, "[[") {
		t.Error("Home.md should contain wikilinks")
	}

	// Verify entry has transformed frontmatter (topics → tags)
	entry1Out, readErr := os.ReadFile(filepath.Join(
		outputDir, config.ObsidianDirEntries, "2026-02-14-add-caching-abc12345.md")) //nolint:gosec // test file path
	if readErr != nil {
		t.Fatal(readErr)
	}
	entry1Str := string(entry1Out)
	if strings.Contains(entry1Str, "\ntopics:") {
		t.Error("entry should have 'tags:' not 'topics:' in frontmatter")
	}
	if !strings.Contains(entry1Str, "tags:") {
		t.Error("entry missing 'tags:' in transformed frontmatter")
	}
	if !strings.Contains(entry1Str, "source_file:") {
		t.Error("entry missing 'source_file:' in transformed frontmatter")
	}

	// Verify entry has related footer
	if !strings.Contains(entry1Str, config.ObsidianRelatedHeading) {
		t.Error("entry missing related sessions footer")
	}

	// Verify topic MOC was created (caching has 2 entries = popular)
	assertFileExists(t, filepath.Join(outputDir, config.ObsidianTopicsMOC))
	topicsMOC, readErr := os.ReadFile(filepath.Join(outputDir, config.ObsidianTopicsMOC))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(topicsMOC), "[[caching]]") {
		t.Error("topics MOC missing caching wikilink")
	}

	// Verify popular topic page was created
	assertFileExists(t, filepath.Join(
		outputDir, config.JournalDirTopics, "caching.md"))
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}
