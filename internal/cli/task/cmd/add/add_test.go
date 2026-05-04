//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// TestTaskAdd verifies the noun-first ctx task add subcommand
// writes to TASKS.md when invoked without the deprecated noun
// positional arg.
func TestTaskAdd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-task-add-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{
		"Test task for integration",
		"--section", "Misc",
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
	})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx task add failed: %v", err)
	}

	tasksPath := filepath.Join(tmpDir, ".context", "TASKS.md")
	content, err := os.ReadFile(filepath.Clean(tasksPath))
	if err != nil {
		t.Fatalf("failed to read TASKS.md: %v", err)
	}
	if !strings.Contains(string(content), "Test task for integration") {
		t.Errorf("task was not added to TASKS.md")
	}
}

// TestTaskAddRequiresProvenance verifies that omitting the
// provenance flags produces an error referencing --session-id.
func TestTaskAddRequiresProvenance(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-task-add-prov-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{"Missing provenance", "--section", "Misc"})
	err = addCmd.Execute()
	if err == nil {
		t.Fatal("expected error when adding task without provenance")
	}
	if !strings.Contains(err.Error(), "--session-id") {
		t.Errorf("error should mention --session-id: %v", err)
	}
}
