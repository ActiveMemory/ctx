//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestRunWatch_NoContext(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when no .context/ exists")
	}
	if !strings.Contains(err.Error(), "ctx init") {
		t.Errorf("error = %q, want 'ctx init' suggestion", err.Error())
	}
}

func TestRunWatch_WithLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Create a log file with context-update commands
	logContent := `Some output
<context-update type="task" section="Misc">Task from log file</context-update>
More output
`
	logPath := filepath.Join(tmpDir, "test.log")
	if err := os.WriteFile(logPath, []byte(logContent), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--log", logPath})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("runWatch error: %v", err)
	}

	// Verify task was written
	tasksPath := filepath.Join(rc.ContextDir(), ctx.Task)
	content, err := os.ReadFile(filepath.Clean(tasksPath))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "Task from log file") {
		t.Error("task from log file should be added")
	}
}

func TestRunWatch_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Create a log file with updates
	logContent := `<context-update type="task">Dry run task</context-update>
`
	logPath := filepath.Join(tmpDir, "dry.log")
	if err := os.WriteFile(logPath, []byte(logContent), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--log", logPath, "--dry-run"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("runWatch error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DRY RUN") {
		t.Error("output should indicate dry run mode")
	}
	if !strings.Contains(out, "Would apply") {
		t.Error("output should show what would be applied")
	}
}

func TestRunWatch_InvalidLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--log", "/nonexistent/path/to/log"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent log file")
	}
	if !strings.Contains(err.Error(), "failed to open log file") {
		t.Errorf("error = %q, want 'failed to open log file'", err.Error())
	}
}

func TestCmd_HasFlags(t *testing.T) {
	cmd := Cmd()
	if cmd.Use != "watch" {
		t.Errorf("cmd.Use = %q, want 'watch'", cmd.Use)
	}

	logFlag := cmd.Flags().Lookup("log")
	if logFlag == nil {
		t.Fatal("expected --log flag")
	}

	dryRunFlag := cmd.Flags().Lookup("dry-run")
	if dryRunFlag == nil {
		t.Fatal("expected --dry-run flag")
	}
}
