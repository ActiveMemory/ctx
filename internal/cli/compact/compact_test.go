//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	taskComplete "github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
)

// TestCompactCommand tests the compact command.
func TestCompactCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-compact-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Run compact
	compactCmd := Cmd()
	compactCmd.SetArgs([]string{})
	if err := compactCmd.Execute(); err != nil {
		t.Fatalf("compact failed: %v", err)
	}
}

// TestCompactWithTasks tests the compact command with actual completed tasks.
func TestCompactWithTasks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-compact-tasks-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add and complete a task
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Task to complete"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task failed: %v", err)
	}

	completeCmd := taskComplete.Cmd()
	completeCmd.SetArgs([]string{"Task to complete"})
	if err := completeCmd.Execute(); err != nil {
		t.Fatalf("complete task failed: %v", err)
	}

	// Run compact
	compactCmd := Cmd()
	compactCmd.SetArgs([]string{})
	if err := compactCmd.Execute(); err != nil {
		t.Fatalf("compact failed: %v", err)
	}
}
