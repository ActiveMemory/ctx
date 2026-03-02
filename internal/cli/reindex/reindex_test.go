//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

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

	if cmd.Use != "reindex" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "reindex")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}

	if cmd.RunE == nil {
		t.Error("Cmd().RunE is nil")
	}
}

func TestRunReindex_NoFiles(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	cmd := Cmd()

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when context files do not exist")
	}
}

func TestRunReindex_BothFiles(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	ctxDir := filepath.Join(tempDir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0750)

	decisions := `# Decisions

## [2026-01-15-143022] Use YAML for config

**Context:** Need a config format
**Rationale:** YAML is human-readable
**Consequences:** Added yaml dependency
`
	_ = os.WriteFile(
		filepath.Join(ctxDir, config.FileDecision),
		[]byte(decisions), 0600,
	)

	learnings := `# Learnings

## [2026-01-15-150000] Always validate input

**Context:** Found a bug from invalid input
**Lesson:** Validate at boundaries
**Application:** Add validation to all handlers
`
	_ = os.WriteFile(
		filepath.Join(ctxDir, config.FileLearning),
		[]byte(learnings), 0600,
	)

	cmd := Cmd()

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify both files were updated
	updatedDecisions, readErr := os.ReadFile(filepath.Join(ctxDir, config.FileDecision)) //nolint:gosec // test temp path
	if readErr != nil {
		t.Fatalf("failed to read updated DECISIONS.md: %v", readErr)
	}
	if len(updatedDecisions) == 0 {
		t.Error("updated DECISIONS.md is empty")
	}

	updatedLearnings, readErr := os.ReadFile(filepath.Join(ctxDir, config.FileLearning)) //nolint:gosec // test temp path
	if readErr != nil {
		t.Fatalf("failed to read updated LEARNINGS.md: %v", readErr)
	}
	if len(updatedLearnings) == 0 {
		t.Error("updated LEARNINGS.md is empty")
	}
}

func TestRunReindex_DecisionsMissingLearningsPresent(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	ctxDir := filepath.Join(tempDir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0750)

	// Only create LEARNINGS.md, not DECISIONS.md
	_ = os.WriteFile(
		filepath.Join(ctxDir, config.FileLearning),
		[]byte("# Learnings\n"), 0600,
	)

	cmd := Cmd()

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when DECISIONS.md is missing")
	}
}

func TestRunReindex_EmptyFiles(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	ctxDir := filepath.Join(tempDir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0750)

	_ = os.WriteFile(
		filepath.Join(ctxDir, config.FileDecision),
		[]byte("# Decisions\n"), 0600,
	)
	_ = os.WriteFile(
		filepath.Join(ctxDir, config.FileLearning),
		[]byte("# Learnings\n"), 0600,
	)

	cmd := Cmd()

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
