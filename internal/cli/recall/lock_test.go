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
			got := multipartBase(tt.filename)
			if got != tt.want {
				t.Errorf("multipartBase(%q) = %q, want %q",
					tt.filename, got, tt.want)
			}
		})
	}
}

func TestMatchJournalFiles_All(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.md", "b.md", "c.md", "state.json"} {
		if err := os.WriteFile(
			filepath.Join(dir, name), []byte("x"), config.PermFile,
		); err != nil {
			t.Fatal(err)
		}
	}

	files, err := matchJournalFiles(dir, nil, true)
	if err != nil {
		t.Fatalf("matchJournalFiles: %v", err)
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
		if err := os.WriteFile(
			filepath.Join(dir, name), []byte("x"), config.PermFile,
		); err != nil {
			t.Fatal(err)
		}
	}

	files, err := matchJournalFiles(dir, []string{"abc12345"}, false)
	if err != nil {
		t.Fatalf("matchJournalFiles: %v", err)
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
		if err := os.WriteFile(
			filepath.Join(dir, name), []byte("x"), config.PermFile,
		); err != nil {
			t.Fatal(err)
		}
	}

	files, err := matchJournalFiles(dir, []string{"abc12345"}, false)
	if err != nil {
		t.Fatalf("matchJournalFiles: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("expected 3 matches (base + 2 parts), got %d: %v",
			len(files), files)
	}
}

func TestMatchJournalFiles_MissingDir(t *testing.T) {
	files, err := matchJournalFiles("/nonexistent/path", nil, true)
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected no files, got %d", len(files))
	}
}

func TestUpdateLockFrontmatter_Lock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := "---\ndate: \"2026-01-21\"\ntitle: \"Test\"\n---\n\n# Body\n"
	if err := os.WriteFile(path, []byte(content), config.PermFile); err != nil {
		t.Fatal(err)
	}

	updateLockFrontmatter(path, true)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), lockedFrontmatterLine) {
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
		lockedFrontmatterLine + "\ntitle: \"Test\"\n---\n\n# Body\n"
	if err := os.WriteFile(path, []byte(content), config.PermFile); err != nil {
		t.Fatal(err)
	}

	updateLockFrontmatter(path, false)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
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
	if err := os.WriteFile(path, []byte(content), config.PermFile); err != nil {
		t.Fatal(err)
	}

	updateLockFrontmatter(path, true)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Error("file without frontmatter should be unchanged")
	}
}

func TestUpdateLockFrontmatter_IdempotentLock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := "---\ndate: \"2026-01-21\"\n" +
		lockedFrontmatterLine + "\n---\n\n# Body\n"
	if err := os.WriteFile(path, []byte(content), config.PermFile); err != nil {
		t.Fatal(err)
	}

	updateLockFrontmatter(path, true)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}
	// Should not duplicate the locked line.
	count := strings.Count(string(data), "locked:")
	if count != 1 {
		t.Errorf("expected 1 locked line, got %d", count)
	}
}

func TestRunLockUnlock_LockSingle(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	// Create a journal file with frontmatter.
	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	// Point rc.ContextDir() to our temp dir.
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Lock via CLI.
	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"lock", "abc12345"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("lock: %v\noutput: %s", err, buf.String())
	}

	// Verify state.
	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if !jstate.Locked(filename) {
		t.Error("file should be locked in state")
	}

	// Verify frontmatter.
	data, err := os.ReadFile(filepath.Join(journalDir, filename))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), lockedFrontmatterLine) {
		t.Error("frontmatter should contain locked line")
	}
}

func TestRunLockUnlock_UnlockSingle(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\n" +
		lockedFrontmatterLine + "\n---\n\n# Test\n"
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
	cmd.SetArgs([]string{"unlock", "abc12345"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unlock: %v\noutput: %s", err, buf.String())
	}

	// Verify state.
	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if jstate.Locked(filename) {
		t.Error("file should not be locked after unlock")
	}

	// Verify frontmatter.
	data, readErr := os.ReadFile(filepath.Join(journalDir, filename))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if strings.Contains(string(data), "locked:") {
		t.Error("frontmatter should not contain locked line after unlock")
	}
}

func TestRunLockUnlock_LockAll(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	files := []string{"a.md", "b.md", "c.md"}
	for _, f := range files {
		content := "---\ndate: \"2026-01-21\"\n---\n\n# " + f + "\n"
		if err := os.WriteFile(
			filepath.Join(journalDir, f), []byte(content), config.PermFile,
		); err != nil {
			t.Fatal(err)
		}
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
	cmd.SetArgs([]string{"lock", "--all"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("lock --all: %v\noutput: %s", err, buf.String())
	}

	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	for _, f := range files {
		if !jstate.Locked(f) {
			t.Errorf("%s should be locked", f)
		}
	}
}

func TestRunLockUnlock_AlreadyLocked(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\n" +
		lockedFrontmatterLine + "\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), config.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	// Pre-set locked state.
	jstate := &state.JournalState{
		Version: state.CurrentVersion,
		Entries: map[string]state.FileState{
			filename: {Locked: "2026-01-21"},
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
	cmd.SetArgs([]string{"lock", "abc12345"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("lock: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No changes") {
		t.Errorf("expected 'No changes' for already-locked, got:\n%s", output)
	}
}

func TestRunLockUnlock_NoArgsNoAll(t *testing.T) {
	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"lock"})

	// Should print help (no error).
	if err := cmd.Execute(); err != nil {
		t.Fatalf("bare lock should not error: %v", err)
	}
}

func TestRunLockUnlock_AllWithPattern(t *testing.T) {
	cmd := Cmd()
	buf := new(strings.Builder)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"lock", "--all", "abc12345"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error with --all and a pattern")
	}
	if !strings.Contains(err.Error(), "cannot use --all with a pattern") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunLockUnlock_LockMultipart(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		t.Fatal(err)
	}

	base := "2026-01-21-test-abc12345.md"
	part2 := "2026-01-21-test-abc12345-p2.md"
	for _, f := range []string{base, part2} {
		content := "---\ndate: \"2026-01-21\"\n---\n\n# " + f + "\n"
		if err := os.WriteFile(
			filepath.Join(journalDir, f), []byte(content), config.PermFile,
		); err != nil {
			t.Fatal(err)
		}
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
	cmd.SetArgs([]string{"lock", "abc12345"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("lock: %v\noutput: %s", err, buf.String())
	}

	// Verify both parts locked in state.
	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if !jstate.Locked(base) {
		t.Error("base file should be locked")
	}
	if !jstate.Locked(part2) {
		t.Error("part 2 should be locked")
	}

	// Verify frontmatter on both files.
	for _, f := range []string{base, part2} {
		data, readErr := os.ReadFile(filepath.Join(journalDir, f))
		if readErr != nil {
			t.Fatal(readErr)
		}
		if !strings.Contains(string(data), lockedFrontmatterLine) {
			t.Errorf("%s frontmatter should contain locked line", f)
		}
	}
}
