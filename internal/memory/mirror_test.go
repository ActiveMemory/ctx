//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/memory"
)

func TestSync_FirstRun(t *testing.T) {
	contextDir := t.TempDir()
	sourceDir := t.TempDir()
	sourcePath := filepath.Join(sourceDir, "MEMORY.md")

	content := "# Memory\n\n## Session notes\n- discovered a bug\n"
	if writeErr := os.WriteFile(sourcePath, []byte(content), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	result, syncErr := Sync(contextDir, sourcePath)
	if syncErr != nil {
		t.Fatalf("Sync: %v", syncErr)
	}

	if result.ArchivedTo != "" {
		t.Errorf("expected no archive on first sync, got %q", result.ArchivedTo)
	}
	if result.SourceLines != 4 {
		t.Errorf("SourceLines = %d, want 4", result.SourceLines)
	}

	mirrorPath := filepath.Join(contextDir, dir.Memory, memory.Mirror)
	mirrorData, readErr := os.ReadFile(mirrorPath)
	if readErr != nil {
		t.Fatalf("reading mirror: %v", readErr)
	}
	if string(mirrorData) != content {
		t.Errorf("mirror content mismatch:\ngot:  %q\nwant: %q", string(mirrorData), content)
	}
}

func TestSync_WithArchive(t *testing.T) {
	contextDir := t.TempDir()
	sourceDir := t.TempDir()
	sourcePath := filepath.Join(sourceDir, "MEMORY.md")

	// Create initial mirror
	mirrorDir := filepath.Join(contextDir, dir.Memory)
	if mkErr := os.MkdirAll(mirrorDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	mirrorPath := filepath.Join(mirrorDir, memory.Mirror)
	oldContent := "# Memory v1\n"
	if writeErr := os.WriteFile(mirrorPath, []byte(oldContent), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Write updated source
	newContent := "# Memory v2\n\n## New stuff\n"
	if writeErr := os.WriteFile(sourcePath, []byte(newContent), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	result, syncErr := Sync(contextDir, sourcePath)
	if syncErr != nil {
		t.Fatalf("Sync: %v", syncErr)
	}

	if result.ArchivedTo == "" {
		t.Error("expected archive path on second sync")
	}
	if result.MirrorLines != 1 {
		t.Errorf("MirrorLines = %d, want 1", result.MirrorLines)
	}
	if result.SourceLines != 3 {
		t.Errorf("SourceLines = %d, want 3", result.SourceLines)
	}

	// Verify archive content
	archiveData, readErr := os.ReadFile(result.ArchivedTo)
	if readErr != nil {
		t.Fatalf("reading archive: %v", readErr)
	}
	if string(archiveData) != oldContent {
		t.Errorf("archive content mismatch:\ngot:  %q\nwant: %q", string(archiveData), oldContent)
	}

	// Verify mirror updated
	mirrorData, mirrorReadErr := os.ReadFile(mirrorPath)
	if mirrorReadErr != nil {
		t.Fatalf("reading updated mirror: %v", mirrorReadErr)
	}
	if string(mirrorData) != newContent {
		t.Error("mirror was not updated to new content")
	}
}

func TestDiff_Identical(t *testing.T) {
	contextDir := t.TempDir()
	sourceDir := t.TempDir()

	content := "# Memory\nsame content\n"

	mirrorDir := filepath.Join(contextDir, dir.Memory)
	if mkErr := os.MkdirAll(mirrorDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	mirrorPath := filepath.Join(mirrorDir, memory.Mirror)
	if writeErr := os.WriteFile(mirrorPath, []byte(content), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	sourcePath := filepath.Join(sourceDir, "MEMORY.md")
	if writeErr := os.WriteFile(sourcePath, []byte(content), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	diff, diffErr := Diff(contextDir, sourcePath)
	if diffErr != nil {
		t.Fatalf("Diff: %v", diffErr)
	}
	if diff != "" {
		t.Errorf("expected empty diff for identical files, got:\n%s", diff)
	}
}

func TestDiff_WithChanges(t *testing.T) {
	contextDir := t.TempDir()
	sourceDir := t.TempDir()

	mirrorDir := filepath.Join(contextDir, dir.Memory)
	if mkErr := os.MkdirAll(mirrorDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	mirrorPath := filepath.Join(mirrorDir, memory.Mirror)
	if writeErr := os.WriteFile(mirrorPath, []byte("# Memory\nold line\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	sourcePath := filepath.Join(sourceDir, "MEMORY.md")
	if writeErr := os.WriteFile(sourcePath, []byte("# Memory\nnew line\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	diff, diffErr := Diff(contextDir, sourcePath)
	if diffErr != nil {
		t.Fatalf("Diff: %v", diffErr)
	}
	if !strings.Contains(diff, "-old line") {
		t.Error("diff should contain removed old line")
	}
	if !strings.Contains(diff, "+new line") {
		t.Error("diff should contain added new line")
	}
}

func TestSync_EmptySource(t *testing.T) {
	contextDir := t.TempDir()
	sourceDir := t.TempDir()
	sourcePath := filepath.Join(sourceDir, "MEMORY.md")

	if writeErr := os.WriteFile(sourcePath, []byte(""), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	result, syncErr := Sync(contextDir, sourcePath)
	if syncErr != nil {
		t.Fatalf("Sync: %v", syncErr)
	}
	if result.SourceLines != 0 {
		t.Errorf("SourceLines = %d, want 0 for empty file", result.SourceLines)
	}
}

func TestArchiveCount(t *testing.T) {
	contextDir := t.TempDir()
	archiveDir := filepath.Join(contextDir, dir.MemoryArchive)
	if mkErr := os.MkdirAll(archiveDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	// No archives yet
	if got := ArchiveCount(contextDir); got != 0 {
		t.Errorf("ArchiveCount = %d, want 0", got)
	}

	// Create two archives
	for _, name := range []string{"mirror-2026-03-01-120000.md", "mirror-2026-03-02-120000.md"} {
		if writeErr := os.WriteFile(filepath.Join(archiveDir, name), []byte("x"), 0o644); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	if got := ArchiveCount(contextDir); got != 2 {
		t.Errorf("ArchiveCount = %d, want 2", got)
	}
}
