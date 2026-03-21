//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	"github.com/spf13/cobra"
)

func TestRunJournalObsidianIntegration(t *testing.T) {
	// Create a temporary journal directory with test entries
	tmpDir := t.TempDir()
	journalDir := filepath.Join(tmpDir, dir.Context, dir.Journal)
	if mkErr := os.MkdirAll(journalDir, fs.PermExec); mkErr != nil {
		t.Fatal(mkErr)
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
		"2026-02-14-add-caching-abc12345.md":   entry1,
		"2026-02-13-fix-cache-def67890.md":     entry2,
		"2026-02-12-plain-session-ghi11111.md": entry3,
	}

	for name, content := range entries {
		path := filepath.Join(journalDir, name)
		if writeErr := os.WriteFile(path, []byte(content), fs.PermFile); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	// Run the vault generation
	outputDir := filepath.Join(tmpDir, "vault-output")

	// Use a real cobra.Command with captured output
	cmd := &cobra.Command{}
	cmd.SetOut(&strings.Builder{})
	cmd.SetErr(&strings.Builder{})

	buildErr := BuildObsidianVault(cmd, journalDir, outputDir)
	if buildErr != nil {
		t.Fatalf("BuildObsidianVault failed: %v", buildErr)
	}

	// Verify vault structure
	assertFileExists(t, filepath.Join(outputDir, obsidian.DirConfig, obsidian.AppConfigFile))
	assertFileExists(t, filepath.Join(outputDir, obsidian.MOCHome))
	assertFileExists(t, filepath.Join(outputDir, file.Readme))

	// Verify entries were written
	assertFileExists(t, filepath.Join(outputDir, obsidian.DirEntries, "2026-02-14-add-caching-abc12345.md"))
	assertFileExists(t, filepath.Join(outputDir, obsidian.DirEntries, "2026-02-13-fix-cache-def67890.md"))

	// Verify .obsidian/app.json content
	appConfig, readErr := os.ReadFile(filepath.Join( //nolint:gosec // test file path
		outputDir, obsidian.DirConfig, obsidian.AppConfigFile))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(appConfig), `"useMarkdownLinks": false`) {
		t.Error("app.json missing useMarkdownLinks setting")
	}

	// Verify Home.md contains wikilinks
	home, readErr := os.ReadFile(filepath.Join(outputDir, obsidian.MOCHome)) //nolint:gosec // test file path
	if readErr != nil {
		t.Fatal(readErr)
	}
	homeStr := string(home)
	if !strings.Contains(homeStr, "[[") {
		t.Error("Home.md should contain wikilinks")
	}

	// Verify entry has transformed frontmatter (topics -> tags)
	entry1Out, readErr := os.ReadFile(filepath.Join( //nolint:gosec // test file path
		outputDir, obsidian.DirEntries, "2026-02-14-add-caching-abc12345.md"))
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
	if !strings.Contains(entry1Str, desc.Text(text.DescKeyHeadingObsidianRelated)) {
		t.Error("entry missing related sessions footer")
	}

	// Verify topic MOC was created (caching has 2 entries = popular)
	assertFileExists(t, filepath.Join(outputDir, obsidian.MOCTopics))
	topicsMOC, readErr := os.ReadFile(filepath.Join(outputDir, obsidian.MOCTopics)) //nolint:gosec // test file path
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(topicsMOC), "[[caching]]") {
		t.Error("topics MOC missing caching wikilink")
	}

	// Verify popular topic page was created
	assertFileExists(t, filepath.Join(
		outputDir, dir.JournTopics, "caching.md"))
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Errorf("expected file to exist: %s", path)
	}
}
