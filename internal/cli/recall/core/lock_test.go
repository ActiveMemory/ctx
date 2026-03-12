//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/fs"
)

func TestMultipartBase(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "no multipart suffix",
			filename: "2026-01-21-slug-abc12345.md",
			want:     "2026-01-21-slug-abc12345.md",
		},
		{
			name:     "part 2",
			filename: "2026-01-21-slug-abc12345-p2.md",
			want:     "2026-01-21-slug-abc12345.md",
		},
		{
			name:     "part 10",
			filename: "2026-01-21-slug-abc12345-p10.md",
			want:     "2026-01-21-slug-abc12345.md",
		},
		{
			name:     "not a part suffix",
			filename: "2026-01-21-slug-pickup.md",
			want:     "2026-01-21-slug-pickup.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MultipartBase(tt.filename)
			if got != tt.want {
				t.Errorf("MultipartBase(%q) = %q, want %q",
					tt.filename, got, tt.want)
			}
		})
	}
}

func TestMatchJournalFiles_All(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.md", "b.md", "c.md", "state.json"} {
		if writeErr := os.WriteFile(
			filepath.Join(dir, name), []byte("x"), fs.PermFile,
		); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	files, matchErr := MatchJournalFiles(dir, nil, true)
	if matchErr != nil {
		t.Fatalf("MatchJournalFiles: %v", matchErr)
	}
	if len(files) != 3 {
		t.Errorf("expected 3 .md files, got %d: %v", len(files), files)
	}
}

func TestMatchJournalFiles_Pattern(t *testing.T) {
	dir := t.TempDir()
	names := []string{
		"2026-01-21-hello-abc12345.md",
		"2026-01-22-goodbye-def67890.md",
	}
	for _, name := range names {
		if writeErr := os.WriteFile(
			filepath.Join(dir, name), []byte("x"), fs.PermFile,
		); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	files, matchErr := MatchJournalFiles(dir, []string{"abc12345"}, false)
	if matchErr != nil {
		t.Fatalf("MatchJournalFiles: %v", matchErr)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 match, got %d: %v", len(files), files)
	}
	if len(files) > 0 && files[0] != names[0] {
		t.Errorf("expected %q, got %q", names[0], files[0])
	}
}

func TestMatchJournalFiles_MultipartExpands(t *testing.T) {
	dir := t.TempDir()
	names := []string{
		"2026-01-21-hello-abc12345.md",
		"2026-01-21-hello-abc12345-p2.md",
		"2026-01-21-hello-abc12345-p3.md",
		"2026-01-22-other-def67890.md",
	}
	for _, name := range names {
		if writeErr := os.WriteFile(
			filepath.Join(dir, name), []byte("x"), fs.PermFile,
		); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	files, matchErr := MatchJournalFiles(dir, []string{"abc12345"}, false)
	if matchErr != nil {
		t.Fatalf("MatchJournalFiles: %v", matchErr)
	}
	if len(files) != 3 {
		t.Errorf("expected 3 matches (base + 2 parts), got %d: %v",
			len(files), files)
	}
}

func TestMatchJournalFiles_MissingDir(t *testing.T) {
	files, matchErr := MatchJournalFiles("/nonexistent/path", nil, true)
	if matchErr != nil {
		t.Fatalf("expected nil error for missing dir, got: %v", matchErr)
	}
	if len(files) != 0 {
		t.Errorf("expected no files, got %d", len(files))
	}
}

func TestUpdateLockFrontmatter_Lock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := "---\ndate: \"2026-01-21\"\ntitle: \"Test\"\n---\n\n# Body\n"
	if writeErr := os.WriteFile(path, []byte(content), fs.PermFile); writeErr != nil {
		t.Fatal(writeErr)
	}

	UpdateLockFrontmatter(path, true)

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(data), LockedFrontmatterLine) {
		t.Error("lock should insert locked line into frontmatter")
	}
	if !strings.Contains(string(data), "# Body") {
		t.Error("body content should be preserved")
	}
}

func TestUpdateLockFrontmatter_Unlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := "---\ndate: \"2026-01-21\"\n" +
		LockedFrontmatterLine + "\ntitle: \"Test\"\n---\n\n# Body\n"
	if writeErr := os.WriteFile(path, []byte(content), fs.PermFile); writeErr != nil {
		t.Fatal(writeErr)
	}

	UpdateLockFrontmatter(path, false)

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if strings.Contains(string(data), "locked:") {
		t.Error("unlock should remove locked line from frontmatter")
	}
	if !strings.Contains(string(data), "# Body") {
		t.Error("body content should be preserved")
	}
}

func TestUpdateLockFrontmatter_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := "# No frontmatter here\n\nJust a body.\n"
	if writeErr := os.WriteFile(path, []byte(content), fs.PermFile); writeErr != nil {
		t.Fatal(writeErr)
	}

	UpdateLockFrontmatter(path, true)

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != content {
		t.Error("file without frontmatter should be unchanged")
	}
}

func TestUpdateLockFrontmatter_IdempotentLock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := "---\ndate: \"2026-01-21\"\n" +
		LockedFrontmatterLine + "\n---\n\n# Body\n"
	if writeErr := os.WriteFile(path, []byte(content), fs.PermFile); writeErr != nil {
		t.Fatal(writeErr)
	}

	UpdateLockFrontmatter(path, true)

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatal(readErr)
	}
	// Should not duplicate the locked line.
	count := strings.Count(string(data), "locked:")
	if count != 1 {
		t.Errorf("expected 1 locked line, got %d", count)
	}
}

func TestFrontmatterHasLocked(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "locked true",
			content: "---\ndate: \"2026-01-21\"\nlocked: true\n---\n\n# Body\n",
			want:    true,
		},
		{
			name:    "locked true with managed comment",
			content: "---\ndate: \"2026-01-21\"\nlocked: true  # managed by ctx\n---\n\n# Body\n",
			want:    true,
		},
		{
			name:    "locked false",
			content: "---\ndate: \"2026-01-21\"\nlocked: false\n---\n\n# Body\n",
			want:    false,
		},
		{
			name:    "no locked field",
			content: "---\ndate: \"2026-01-21\"\ntitle: \"Test\"\n---\n\n# Body\n",
			want:    false,
		},
		{
			name:    "no frontmatter",
			content: "# No frontmatter here\n\nJust a body.\n",
			want:    false,
		},
		{
			name:    "empty file",
			content: "",
			want:    false,
		},
		{
			name:    "locked with extra whitespace",
			content: "---\ndate: \"2026-01-21\"\n  locked:   true  \n---\n\n# Body\n",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.md")
			if writeErr := os.WriteFile(path, []byte(tt.content), fs.PermFile); writeErr != nil {
				t.Fatal(writeErr)
			}

			got := FrontmatterHasLocked(path)
			if got != tt.want {
				t.Errorf("FrontmatterHasLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrontmatterHasLocked_MissingFile(t *testing.T) {
	got := FrontmatterHasLocked("/nonexistent/path/test.md")
	if got {
		t.Error("missing file should return false")
	}
}
