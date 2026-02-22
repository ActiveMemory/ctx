//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package learnings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "learnings" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "learnings")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}
}

func TestCmd_HasReindexSubcommand(t *testing.T) {
	cmd := Cmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "reindex" {
			found = true
			if sub.Short == "" {
				t.Error("reindex subcommand has empty Short description")
			}
			if sub.RunE == nil {
				t.Error("reindex subcommand has no RunE function")
			}
			break
		}
	}

	if !found {
		t.Error("reindex subcommand not found")
	}
}

func TestRunReindex_NoFile(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	cmd := Cmd()
	cmd.SetArgs([]string{"reindex"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when LEARNINGS.md does not exist")
	}
}

func TestRunReindex_WithFile(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	// Create the context directory and LEARNINGS.md file
	ctxDir := filepath.Join(tempDir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0750)

	content := `# Learnings

## 2026-01-15 — Always validate input

**Context:** Found a bug from invalid input
**Lesson:** Validate at boundaries
**Application:** Add validation to all handlers
`
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileLearning), []byte(content), 0600)

	cmd := Cmd()
	cmd.SetArgs([]string{"reindex"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the file was updated
	updated, err := os.ReadFile(filepath.Join(ctxDir, config.FileLearning)) //nolint:gosec // test temp path
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}
	if len(updated) == 0 {
		t.Error("updated file is empty")
	}
}

func TestRunReindex_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	// Create the context directory and empty LEARNINGS.md
	ctxDir := filepath.Join(tempDir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0750)
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileLearning), []byte("# Learnings\n"), 0600)

	cmd := Cmd()
	cmd.SetArgs([]string{"reindex"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
