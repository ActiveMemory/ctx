//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// TestDriftCommand tests the drift command.
func TestDriftCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-test-*")
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

	// Run drift - just verify it runs without error
	driftCmd := Cmd()
	driftCmd.SetArgs([]string{})

	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift command failed: %v", err)
	}
}

// TestDriftJSONOutput tests the drift command with JSON output.
func TestDriftJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-json-test-*")
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

	// Test drift with JSON output
	driftCmd := Cmd()
	driftCmd.SetArgs([]string{"--json"})
	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift --json failed: %v", err)
	}
}

// TestRunDrift_NoContext tests drift when no .context/ exists.
func TestRunDrift_NoContext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-nocontext-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{})

	runErr := cmd.Execute()
	if runErr == nil {
		t.Fatal("expected error when no .context/ exists")
	}
	if !strings.Contains(runErr.Error(), "not initialized") {
		t.Errorf("unexpected error: %v", runErr)
	}
}

// helper: set up a temp directory as working dir and init context
func setupContextDir(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "cli-drift-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	rc.Reset()

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	initCmd.SilenceUsage = true
	initCmd.SilenceErrors = true
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	return tmpDir, func() {
		_ = os.Chdir(origDir)
		_ = os.RemoveAll(tmpDir)
		rc.Reset()
	}
}

func TestRunDrift_WithFix(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Write TASKS.md with completed tasks to trigger staleness fix
	tasksPath := filepath.Join(tmpDir, dir.Context, ctx.Task)
	tasksContent := "# Tasks\n\n## In Progress\n\n" +
		"- [ ] Do something\n\n## Completed\n\n" +
		"- [x] Done thing 1\n- [x] Done thing 2\n" +
		"- [x] Done thing 3\n- [x] Done thing 4\n" +
		"- [x] Done thing 5\n- [x] Done thing 6\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--fix"})

	// This may or may not error depending on whether the stale detection
	// actually triggers - just test it doesn't panic
	_ = cmd.Execute()
}

func TestRunDrift_JSONWithViolations(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Create a file that looks like it has secrets to trigger a violation
	constPath := filepath.Join(tmpDir, dir.Context, "CONSTITUTION.md")
	constContent := "# Constitution\n\n- NEVER commit secrets\n"
	if err := os.WriteFile(constPath, []byte(constContent), 0600); err != nil {
		t.Fatalf("failed to write CONSTITUTION.md: %v", err)
	}

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--json"})

	// The command may succeed or fail depending on what drift.Detect finds
	_ = cmd.Execute()
}

func TestRunDrift_FixWithStaleness(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Create TASKS.md with many completed tasks to trigger staleness
	tasksPath := filepath.Join(tmpDir, dir.Context, ctx.Task)
	var sb strings.Builder
	sb.WriteString("# Tasks\n\n## In Progress\n\n" +
		"- [ ] Active task\n\n## Completed\n\n")
	for i := 0; i < 10; i++ {
		fmt.Fprintf(&sb, "- [x] Completed task %d\n", i)
	}
	if err := os.WriteFile(tasksPath, []byte(sb.String()), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd := Cmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--fix"})

	// Execute - may or may not error depending on other drift checks
	_ = cmd.Execute()
}

func TestRunDrift_FixTriggersRecheck(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Remove a required file so fixMissingFile gets called and succeeds
	constPath := filepath.Join(tmpDir, dir.Context, "CONSTITUTION.md")
	_ = os.Remove(constPath)

	// Use Cmd directly with captured output
	cmd := Cmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--fix"})

	_ = cmd.Execute()

	out := buf.String()
	// The fix should have recreated CONSTITUTION.md and triggered re-check
	if !strings.Contains(out, "Applying fixes") {
		t.Errorf("expected 'Applying fixes' in output, got: %s", out)
	}
}

func TestRunDrift_GenericError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-generr-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	// Create .context as a file, not a directory.
	fakePath := filepath.Join(tmpDir, dir.Context)
	if err := os.WriteFile(fakePath, []byte("not a dir"), 0600); err != nil {
		t.Fatalf("failed to create fake .context: %v", err)
	}

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{})

	runErr := cmd.Execute()
	if runErr == nil {
		t.Fatal("expected error when .context is a file")
	}
}
