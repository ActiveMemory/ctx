//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildSessionIndex_WithSessionID(t *testing.T) {
	dir := t.TempDir()

	// File with session_id in frontmatter.
	content := "---\ndate: \"2026-01-15\"\nsession_id: \"abc12345-full-uuid\"\n---\n\n# Test\n"
	if err := os.WriteFile(filepath.Join(dir, "2026-01-15-fix-auth-abc12345.md"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	idx := buildSessionIndex(dir)

	if got, ok := idx["abc12345-full-uuid"]; !ok {
		t.Error("expected session_id key in index")
	} else if got != "2026-01-15-fix-auth-abc12345.md" {
		t.Errorf("index[session_id] = %q, want %q", got, "2026-01-15-fix-auth-abc12345.md")
	}
}

func TestBuildSessionIndex_ShortIDFallback(t *testing.T) {
	dir := t.TempDir()

	// Legacy file without session_id (no frontmatter with session_id).
	content := "---\ndate: \"2026-01-15\"\n---\n\n# old-slug\n"
	if err := os.WriteFile(filepath.Join(dir, "2026-01-15-old-slug-abc12345.md"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	idx := buildSessionIndex(dir)

	if got, ok := idx["abc12345"]; !ok {
		t.Error("expected short ID key in index")
	} else if got != "2026-01-15-old-slug-abc12345.md" {
		t.Errorf("index[shortID] = %q, want %q", got, "2026-01-15-old-slug-abc12345.md")
	}
}

func TestBuildSessionIndex_SkipsMultipartFiles(t *testing.T) {
	dir := t.TempDir()

	// Base file.
	base := "---\ndate: \"2026-01-15\"\n---\n\n# test\n"
	if err := os.WriteFile(filepath.Join(dir, "2026-01-15-test-slug-abc12345.md"), []byte(base), 0600); err != nil {
		t.Fatal(err)
	}
	// Multipart file (-p2).
	part := "---\ndate: \"2026-01-15\"\n---\n\n# test part 2\n"
	if err := os.WriteFile(filepath.Join(dir, "2026-01-15-test-slug-abc12345-p2.md"), []byte(part), 0600); err != nil {
		t.Fatal(err)
	}

	idx := buildSessionIndex(dir)

	// Should only have one entry for the short ID.
	if got, ok := idx["abc12345"]; !ok {
		t.Error("expected short ID key in index")
	} else if got != "2026-01-15-test-slug-abc12345.md" {
		t.Errorf("index[shortID] = %q, want base file", got)
	}
}

func TestBuildSessionIndex_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	idx := buildSessionIndex(dir)
	if len(idx) != 0 {
		t.Errorf("expected empty index, got %d entries", len(idx))
	}
}

func TestBuildSessionIndex_NonexistentDir(t *testing.T) {
	idx := buildSessionIndex("/nonexistent/path/to/journal")
	if len(idx) != 0 {
		t.Errorf("expected empty index for nonexistent dir, got %d entries", len(idx))
	}
}

func TestLookupSessionFile(t *testing.T) {
	idx := map[string]string{
		"abc12345-full-uuid": "2026-01-15-fix-auth-abc12345.md",
		"def67890":           "2026-01-16-old-slug-def67890.md",
	}

	tests := []struct {
		name      string
		sessionID string
		want      string
	}{
		{"full ID match", "abc12345-full-uuid", "2026-01-15-fix-auth-abc12345.md"},
		{"short ID fallback", "def67890-some-longer-id", "2026-01-16-old-slug-def67890.md"},
		{"no match", "xxxxxxxx-unknown", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lookupSessionFile(idx, tt.sessionID)
			if got != tt.want {
				t.Errorf("lookupSessionFile(%q) = %q, want %q", tt.sessionID, got, tt.want)
			}
		})
	}
}

func TestExtractSessionID(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			"with session_id",
			"---\ndate: \"2026-01-15\"\nsession_id: \"abc12345-uuid\"\n---\n\n# Test\n",
			"abc12345-uuid",
		},
		{
			"without session_id",
			"---\ndate: \"2026-01-15\"\n---\n\n# Test\n",
			"",
		},
		{
			"no frontmatter",
			"# Just a heading\nSome text\n",
			"",
		},
		{
			"unquoted session_id",
			"---\nsession_id: abc12345-uuid\n---\n\n# Test\n",
			"abc12345-uuid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSessionID(tt.content)
			if got != tt.want {
				t.Errorf("extractSessionID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractFrontmatterField(t *testing.T) {
	content := "---\ndate: \"2026-01-15\"\ntitle: \"Fix Auth Bug\"\nsession_id: \"abc-123\"\n---\n\n# Test\n"

	tests := []struct {
		field string
		want  string
	}{
		{"title", "Fix Auth Bug"},
		{"session_id", "abc-123"},
		{"date", "2026-01-15"},
		{"missing", ""},
	}
	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			got := extractFrontmatterField(content, tt.field)
			if got != tt.want {
				t.Errorf("extractFrontmatterField(%q) = %q, want %q", tt.field, got, tt.want)
			}
		})
	}
}

func TestRenameJournalFiles(t *testing.T) {
	dir := t.TempDir()

	// Create old base file.
	oldBase := "2026-01-15-old-slug-abc12345"
	if err := os.WriteFile(
		filepath.Join(dir, oldBase+".md"),
		[]byte("# old content"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	newBase := "2026-01-15-fix-auth-abc12345"
	renameJournalFiles(dir, oldBase, newBase, 1)

	// Old file should not exist.
	if _, err := os.Stat(filepath.Join(dir, oldBase+".md")); err == nil {
		t.Error("old file should not exist after rename")
	}
	// New file should exist.
	if _, err := os.Stat(filepath.Join(dir, newBase+".md")); err != nil {
		t.Error("new file should exist after rename")
	}
}

func TestRenameJournalFiles_Multipart(t *testing.T) {
	dir := t.TempDir()

	oldBase := "2026-01-15-old-slug-abc12345"
	newBase := "2026-01-15-fix-auth-abc12345"

	// Create base and part files with nav links.
	baseContent := "# old\n[Next →](" + oldBase + "-p2.md)\n"
	p2Content := "# old p2\n[← Previous](" + oldBase + ".md)\n[Next →](" + oldBase + "-p3.md)\n"
	p3Content := "# old p3\n[← Previous](" + oldBase + "-p2.md)\n"

	for fname, content := range map[string]string{
		oldBase + ".md":    baseContent,
		oldBase + "-p2.md": p2Content,
		oldBase + "-p3.md": p3Content,
	} {
		if err := os.WriteFile(filepath.Join(dir, fname), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}

	renameJournalFiles(dir, oldBase, newBase, 3)

	// Verify all old files are gone.
	for _, suffix := range []string{".md", "-p2.md", "-p3.md"} {
		if _, err := os.Stat(filepath.Join(dir, oldBase+suffix)); err == nil {
			t.Errorf("old file %s should not exist", oldBase+suffix)
		}
	}

	// Verify all new files exist.
	for _, suffix := range []string{".md", "-p2.md", "-p3.md"} {
		if _, err := os.Stat(filepath.Join(dir, newBase+suffix)); err != nil {
			t.Errorf("new file %s should exist", newBase+suffix)
		}
	}

	// Verify nav links were updated.
	data, err := os.ReadFile(filepath.Join(dir, newBase+"-p2.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !contains(content, newBase+".md") {
		t.Error("p2 should link to new base file")
	}
	if !contains(content, newBase+"-p3.md") {
		t.Error("p2 should link to new p3 file")
	}
	if contains(content, oldBase) {
		t.Error("p2 should not contain old base name")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
