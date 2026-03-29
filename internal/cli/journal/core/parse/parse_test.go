//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseJournalEntry(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-01-21-test-slug-abc12345.md"
	content := `# Test Session

**Time**: 14:30:00
**Project**: my-project

Some content here.
`
	path := filepath.Join(tmpDir, filename)
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatalf("failed to write temp file: %v", writeErr)
	}

	entry := JournalEntry(path, filename)

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
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatalf("failed to write temp file: %v", writeErr)
	}

	entry := JournalEntry(path, filename)

	if !entry.Suggestive {
		t.Error("Suggestive should be true for suggestion mode sessions")
	}
}

func TestParseJournalEntry_MissingFile(t *testing.T) {
	entry := JournalEntry("/nonexistent/path.md", "2026-01-21-test.md")

	if entry.Title != "2026-01-21-test" {
		t.Errorf("Title = %q, want %q", entry.Title, "2026-01-21-test")
	}
}

func TestParseJournalEntry_WithFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-02-14-feature-abc12345.md"
	content := "---\ntitle: \"Feature: Add caching\"\n" +
		"date: 2026-02-14\ntime: \"14:30:00\"\n" +
		"project: ctx\nsession_id: sess-abc123\n" +
		"model: opus\ntokens_in: 1000\n" +
		"tokens_out: 2000\ntype: feature\n" +
		"outcome: completed\ntopics:\n" +
		"  - caching\n  - performance\n" +
		"key_files:\n" +
		"  - internal/cache/store.go\n" +
		"summary: Added a caching layer\n" +
		"---\n\n# Feature: Add caching\n\nContent.\n"
	path := filepath.Join(tmpDir, filename)
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	entry := JournalEntry(path, filename)

	if entry.Title != "Feature: Add caching" {
		t.Errorf("Title = %q, want %q", entry.Title, "Feature: Add caching")
	}
	if entry.Time != "14:30:00" {
		t.Errorf("Time = %q, want %q", entry.Time, "14:30:00")
	}
	if entry.Project != "ctx" {
		t.Errorf("Project = %q, want %q", entry.Project, "ctx")
	}
	if entry.SessionID != "sess-abc123" {
		t.Errorf("SessionID = %q, want %q", entry.SessionID, "sess-abc123")
	}
	if entry.Model != "opus" {
		t.Errorf("Model = %q, want %q", entry.Model, "opus")
	}
	if entry.TokensIn != 1000 {
		t.Errorf("TokensIn = %d, want 1000", entry.TokensIn)
	}
	if entry.TokensOut != 2000 {
		t.Errorf("TokensOut = %d, want 2000", entry.TokensOut)
	}
	if entry.Type != "feature" {
		t.Errorf("Type = %q, want %q", entry.Type, "feature")
	}
	if entry.Outcome != "completed" {
		t.Errorf("Outcome = %q, want %q", entry.Outcome, "completed")
	}
	if len(entry.Topics) != 2 {
		t.Errorf("Topics len = %d, want 2", len(entry.Topics))
	}
	if len(entry.KeyFiles) != 1 {
		t.Errorf("KeyFiles len = %d, want 1", len(entry.KeyFiles))
	}
	if entry.Summary != "Added a caching layer" {
		t.Errorf("Summary = %q, want %q", entry.Summary, "Added a caching layer")
	}
}

func TestParseJournalEntry_SessionID(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-02-20-with-session-abc12345.md"
	content := "---\ntitle: Session With ID\n" +
		"date: 2026-02-20\n" +
		"session_id: 01abc-def-456\n" +
		"---\n\n# Session With ID\n"
	path := filepath.Join(tmpDir, filename)
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	entry := JournalEntry(path, filename)

	if entry.SessionID != "01abc-def-456" {
		t.Errorf("SessionID = %q, want %q", entry.SessionID, "01abc-def-456")
	}
}

func TestParseJournalEntry_NoSessionID(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-02-20-no-session-abc12345.md"
	content := "---\ntitle: No Session ID\n" +
		"date: 2026-02-20\n---\n\n# No Session ID\n"
	path := filepath.Join(tmpDir, filename)
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	entry := JournalEntry(path, filename)

	if entry.SessionID != "" {
		t.Errorf("SessionID = %q, want empty", entry.SessionID)
	}
}

func TestParseJournalEntry_TitleSanitization(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantTitle string
	}{
		{
			name:      "angle brackets escaped",
			title:     "Fix <details> rendering",
			wantTitle: "Fix &lt;details&gt; rendering",
		},
		{
			name:      "backticks stripped",
			title:     "Update `config.go` settings",
			wantTitle: "Update config.go settings",
		},
		{
			name:      "hash stripped",
			title:     "Issue #42 fix",
			wantTitle: "Issue 42 fix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filename := "2026-03-01-test-abc12345.md"
			content := "---\ntitle: \"" + tt.title +
				"\"\ndate: 2026-03-01\n---\n\n# " +
				tt.title + "\n"
			path := filepath.Join(tmpDir, filename)
			if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
				t.Fatal(writeErr)
			}

			entry := JournalEntry(path, filename)
			if entry.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", entry.Title, tt.wantTitle)
			}
		})
	}
}
