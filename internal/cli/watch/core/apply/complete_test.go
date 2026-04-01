//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply

import (
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestCompleteTaskNoMatch tests complete with no matching task.
func TestCompleteTaskNoMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-nomatch-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Try to complete a non-existent task
	err = completeTask("nonexistent task query")
	if err == nil {
		t.Error("expected error for non-matching task")
	}
	if !strings.Contains(err.Error(), "no task matching") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompleteTask_Empty(t *testing.T) {
	err := completeTask("")
	if err == nil {
		t.Fatal("expected error for empty query")
	}
	if !strings.Contains(err.Error(), "no task specified") {
		t.Errorf("error = %q, want 'no task specified'", err.Error())
	}
}
