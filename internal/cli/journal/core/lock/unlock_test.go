//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lock_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/journal"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/lock"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

func TestRunLockUnlock_LockSingle(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	// Create a journal file with frontmatter.
	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), fs.PermFile,
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
	cmd := journal.Cmd()
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
	if !strings.Contains(string(data), lock.LockedFrontmatterLine) {
		t.Error("frontmatter should contain locked line")
	}
}

func TestRunLockUnlock_UnlockSingle(t *testing.T) {
	dir := t.TempDir()
	journalDir := filepath.Join(dir, ".context", "journal")
	if err := os.MkdirAll(journalDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\n" +
		lock.LockedFrontmatterLine + "\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), fs.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	// Pre-set locked state.
	jstate := &state.State{
		Version: state.CurrentVersion,
		Entries: map[string]state.File{
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

	cmd := journal.Cmd()
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
	if err := os.MkdirAll(journalDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	files := []string{"a.md", "b.md", "c.md"}
	for _, f := range files {
		content := "---\ndate: \"2026-01-21\"\n---\n\n# " + f + "\n"
		if err := os.WriteFile(
			filepath.Join(journalDir, f), []byte(content), fs.PermFile,
		); err != nil {
			t.Fatal(err)
		}
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := journal.Cmd()
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
	if err := os.MkdirAll(journalDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	filename := "2026-01-21-test-abc12345.md"
	content := "---\ndate: \"2026-01-21\"\n" +
		lock.LockedFrontmatterLine + "\n---\n\n# Test\n"
	if err := os.WriteFile(
		filepath.Join(journalDir, filename), []byte(content), fs.PermFile,
	); err != nil {
		t.Fatal(err)
	}

	// Pre-set locked state.
	jstate := &state.State{
		Version: state.CurrentVersion,
		Entries: map[string]state.File{
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

	cmd := journal.Cmd()
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
	cmd := journal.Cmd()
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
	cmd := journal.Cmd()
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
	if err := os.MkdirAll(journalDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	base := "2026-01-21-test-abc12345.md"
	part2 := "2026-01-21-test-abc12345-p2.md"
	for _, f := range []string{base, part2} {
		content := "---\ndate: \"2026-01-21\"\n---\n\n# " + f + "\n"
		if err := os.WriteFile(
			filepath.Join(journalDir, f), []byte(content), fs.PermFile,
		); err != nil {
			t.Fatal(err)
		}
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := journal.Cmd()
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
		if !strings.Contains(string(data), lock.LockedFrontmatterLine) {
			t.Errorf("%s frontmatter should contain locked line", f)
		}
	}
}
