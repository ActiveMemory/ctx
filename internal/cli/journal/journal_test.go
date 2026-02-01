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
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "journal" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "journal")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}
}

func TestCmd_HasSiteSubcommand(t *testing.T) {
	cmd := Cmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "site" {
			found = true
			if sub.Short == "" {
				t.Error("site subcommand has empty Short description")
			}
			if sub.RunE == nil {
				t.Error("site subcommand has no RunE function")
			}

			// Check flags
			outputFlag := sub.Flags().Lookup("output")
			if outputFlag == nil {
				t.Error("site subcommand missing --output flag")
			}

			buildFlag := sub.Flags().Lookup("build")
			if buildFlag == nil {
				t.Error("site subcommand missing --build flag")
			}

			serveFlag := sub.Flags().Lookup("serve")
			if serveFlag == nil {
				t.Error("site subcommand missing --serve flag")
			}

			break
		}
	}

	if !found {
		t.Error("site subcommand not found")
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0B"},
		{100, "100B"},
		{1023, "1023B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{10240, "10.0KB"},
		{1048576, "1.0MB"},
		{1572864, "1.5MB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatSize(tt.bytes)
			if got != tt.want {
				t.Errorf("formatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestParseJournalEntry(t *testing.T) {
	// Create a temp file with journal content
	tmpDir := t.TempDir()
	filename := "2026-01-21-test-slug-abc12345.md"
	content := `# Test Session

**Time**: 14:30:00
**Project**: my-project

Some content here.
`
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	entry := parseJournalEntry(path, filename)

	if entry.Filename != filename {
		t.Errorf("Filename = %q, want %q", entry.Filename, filename)
	}

	if entry.Date != "2026-01-21" {
		t.Errorf("Date = %q, want %q", entry.Date, "2026-01-21")
	}

	if entry.Title != "Test Session" {
		t.Errorf("Title = %q, want %q", entry.Title, "Test Session")
	}

	if entry.Time != "14:30:00" {
		t.Errorf("Time = %q, want %q", entry.Time, "14:30:00")
	}

	if entry.Project != "my-project" {
		t.Errorf("Project = %q, want %q", entry.Project, "my-project")
	}

	if entry.Size != int64(len(content)) {
		t.Errorf("Size = %d, want %d", entry.Size, len(content))
	}
}

func TestParseJournalEntry_SuggestionMode(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-01-21-suggestion-abc12345.md"
	content := `# Suggestion

[SUGGESTION MODE: some suggestion]

Content here.
`
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	entry := parseJournalEntry(path, filename)

	if !entry.IsSuggestion {
		t.Error("IsSuggestion should be true for suggestion mode sessions")
	}
}

func TestParseJournalEntry_MissingFile(t *testing.T) {
	entry := parseJournalEntry("/nonexistent/path.md", "2026-01-21-test.md")

	// Should use filename as title fallback
	if entry.Title != "2026-01-21-test" {
		t.Errorf("Title = %q, want %q", entry.Title, "2026-01-21-test")
	}
}

func TestGenerateIndex(t *testing.T) {
	entries := []journalEntry{
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
			Filename:     "2026-01-19-suggestion-ghi11111.md",
			Title:        "Suggestion",
			Date:         "2026-01-19",
			Time:         "09:00:00",
			IsSuggestion: true,
			Size:         512,
		},
	}

	index := generateIndex(entries)

	// Should have header
	if !strings.Contains(index, "# Session Journal") {
		t.Error("index missing header")
	}

	// Should have session count
	if !strings.Contains(index, "**Sessions**: 2") {
		t.Error("index missing session count")
	}

	// Should have suggestions count
	if !strings.Contains(index, "**Suggestions**: 1") {
		t.Error("index missing suggestions count")
	}

	// Should have month headers
	if !strings.Contains(index, "## 2026-01") {
		t.Error("index missing month header")
	}

	// Should have entry links
	if !strings.Contains(index, "[Session One]") {
		t.Error("index missing session one link")
	}

	// Should have suggestions section
	if !strings.Contains(index, "## Suggestions") {
		t.Error("index missing suggestions section")
	}
}

func TestFormatIndexEntry(t *testing.T) {
	entry := journalEntry{
		Filename: "2026-01-21-test-abc12345.md",
		Title:    "Test Session",
		Date:     "2026-01-21",
		Time:     "14:30:00",
		Project:  "my-project",
		Size:     1536,
	}

	result := formatIndexEntry(entry, "\n")

	// Should have time prefix
	if !strings.Contains(result, "14:30") {
		t.Error("entry missing time prefix")
	}

	// Should have title link
	if !strings.Contains(result, "[Test Session]") {
		t.Error("entry missing title")
	}

	// Should have link to md file
	if !strings.Contains(result, "(2026-01-21-test-abc12345.md)") {
		t.Error("entry missing link")
	}

	// Should have project
	if !strings.Contains(result, "(my-project)") {
		t.Error("entry missing project")
	}

	// Should have size
	if !strings.Contains(result, "1.5KB") {
		t.Error("entry missing size")
	}
}

func TestGenerateZensicalToml(t *testing.T) {
	entries := []journalEntry{
		{
			Filename: "2026-01-21-test.md",
			Title:    "Test Session",
		},
	}

	toml := generateZensicalToml(entries)

	// Should have project section
	if !strings.Contains(toml, "[project]") {
		t.Error("toml missing [project] section")
	}

	// Should have site name
	if !strings.Contains(toml, `site_name = "Session Journal"`) {
		t.Error("toml missing site_name")
	}

	// Should have nav
	if !strings.Contains(toml, "nav = [") {
		t.Error("toml missing nav")
	}

	// Should have theme section
	if !strings.Contains(toml, "[project.theme]") {
		t.Error("toml missing theme section")
	}
}
