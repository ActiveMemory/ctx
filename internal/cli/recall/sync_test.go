//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

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
			if err := os.WriteFile(path, []byte(tt.content), config.PermFile); err != nil {
				t.Fatal(err)
			}

			got := frontmatterHasLocked(path)
			if got != tt.want {
				t.Errorf("frontmatterHasLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrontmatterHasLocked_MissingFile(t *testing.T) {
	got := frontmatterHasLocked("/nonexistent/path/test.md")
	if got {
		t.Error("missing file should return false")
	}
}

func TestRunSync_LocksFromFrontmatter(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	// File with locked: true in frontmatter but not in state.
	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\nlocked: true  # managed by ctx\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"sync"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("sync: %v\noutput: %s", err, buf.String())
	}

	// Verify state was updated.
	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if !jstate.Locked(filename) {
		t.Error("file should be locked in state after sync")
	}

	output := buf.String()
	if !strings.Contains(output, "(locked)") {
		t.Errorf("expected '(locked)' in output, got:\n%s", output)
	}
}

func TestRunSync_UnlocksFromFrontmatter(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	// File WITHOUT locked in frontmatter, but locked in state.
	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\ntitle: \"Test\"\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	// Pre-set locked state.
	jstate := &state.JournalState{
		Version: state.CurrentVersion,
		Entries: map[string]state.FileState{
			filename: {Exported: "2026-01-21", Locked: "2026-01-21"},
		},
	}
	if err := jstate.Save(journalDir); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"sync"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("sync: %v\noutput: %s", err, buf.String())
	}

	// Verify state was cleared.
	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if jstate.Locked(filename) {
		t.Error("file should not be locked in state after sync")
	}

	output := buf.String()
	if !strings.Contains(output, "(unlocked)") {
		t.Errorf("expected '(unlocked)' in output, got:\n%s", output)
	}
}

func TestRunSync_NoChanges(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	// Locked in both frontmatter and state — already in sync.
	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\nlocked: true\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	jstate := &state.JournalState{
		Version: state.CurrentVersion,
		Entries: map[string]state.FileState{
			filename: {Exported: "2026-01-21", Locked: "2026-01-21"},
		},
	}
	if err := jstate.Save(journalDir); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"sync"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("sync: %v\noutput: %s", err, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "No changes") {
		t.Errorf("expected 'No changes' when already in sync, got:\n%s", output)
	}
}

func TestRunSync_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"sync"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("sync: %v\noutput: %s", err, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "No journal entries found") {
		t.Errorf("expected 'No journal entries found', got:\n%s", output)
	}
}

func TestRunSync_MixedFiles(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	// File A: locked in frontmatter, not in state → should lock.
	fileA := "2026-01-21-test-aaa11111.md"
	contentA := "---\ndate: \"2026-01-21\"\nlocked: true\n---\n\n# A\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, fileA), []byte(contentA), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	// File B: no locked in frontmatter, locked in state → should unlock.
	fileB := "2026-01-22-test-bbb22222.md"
	contentB := "---\ndate: \"2026-01-22\"\n---\n\n# B\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, fileB), []byte(contentB), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	// File C: no frontmatter, not locked → no change.
	fileC := "2026-01-23-test-ccc33333.md"
	contentC := "# C\n\nNo frontmatter.\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, fileC), []byte(contentC), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	jstate := &state.JournalState{
		Version: state.CurrentVersion,
		Entries: map[string]state.FileState{
			fileB: {Exported: "2026-01-22", Locked: "2026-01-22"},
		},
	}
	if err := jstate.Save(journalDir); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"sync"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("sync: %v\noutput: %s", err, buf.String())
	}

	// Reload state and verify.
	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}

	if !jstate.Locked(fileA) {
		t.Error("file A should be locked after sync")
	}
	if jstate.Locked(fileB) {
		t.Error("file B should be unlocked after sync")
	}
	if jstate.Locked(fileC) {
		t.Error("file C should not be locked")
	}

	output := buf.String()
	if !strings.Contains(output, "Locked 1") {
		t.Errorf("expected 'Locked 1' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "Unlocked 1") {
		t.Errorf("expected 'Unlocked 1' in output, got:\n%s", output)
	}
}
