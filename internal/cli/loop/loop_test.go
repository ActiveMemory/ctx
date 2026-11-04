//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestLoopCommand tests the loop command.
func TestLoopCommand(t *testing.T) {
	lookup.Init()
	tmpDir, err := os.MkdirTemp("", "cli-loop-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Create the default prompt file at .context/loop.md
	if err := os.MkdirAll(".context", 0750); err != nil {
		t.Fatalf("failed to create context dir: %v", err)
	}
	loopPath := ".context/loop.md"
	if err := os.WriteFile(loopPath, []byte("# Test Prompt\n"), 0600); err != nil {
		t.Fatalf("failed to create loop.md: %v", err)
	}

	// Test loop command
	loopCmd := Cmd()
	loopCmd.SetArgs([]string{"--tool", "generic"})
	if err := loopCmd.Execute(); err != nil {
		t.Fatalf("loop command failed: %v", err)
	}

	// Verify loop.sh was created
	if _, err := os.Stat("loop.sh"); os.IsNotExist(err) {
		t.Error("loop.sh was not created")
	}

	// Test invalid regime rejection
	invalidRegimeCmd := Cmd()
	invalidRegimeCmd.SetArgs([]string{"--regime", "bogus"})
	if err := invalidRegimeCmd.Execute(); err == nil || !strings.Contains(err.Error(), "invalid loop regime") {
		t.Fatalf("expected invalid loop regime error, got %v", err)
	}

	// Test negative --max-turns rejection
	negTurnsCmd := Cmd()
	negTurnsCmd.SetArgs([]string{"--max-turns", "-1"})
	if err := negTurnsCmd.Execute(); err == nil || !strings.Contains(err.Error(), "--max-turns must be non-negative") {
		t.Fatalf("expected negative max-turns error, got %v", err)
	}
}
