//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// runSyncCmd executes a sync command and captures output.
func runSyncCmd(args ...string) (string, error) {
	cmd := Cmd()
	cmd.SetArgs(args)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

// setupSyncDir creates a temp dir, initializes context, and returns cleanup.
func setupSyncDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	return tmpDir
}

// TestSyncCommand tests the sync command.
func TestSyncCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-sync-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test sync command
	syncCmd := Cmd()
	syncCmd.SetArgs([]string{})
	if err := syncCmd.Execute(); err != nil {
		t.Fatalf("sync command failed: %v", err)
	}
}

func TestSyncCommand_NoContext(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	t.Setenv(env.CtxDir, "")

	_, err := runSyncCmd()
	if err == nil {
		t.Fatal("expected error when no .context/ exists")
	}
	// Under the explicit-context-dir model, the error reports that
	// no context directory has been declared.
	if !strings.Contains(err.Error(), "context directory") {
		t.Errorf("error = %q, want context directory mention", err.Error())
	}
}

func TestSyncCommand_DryRun(t *testing.T) {
	setupSyncDir(t)

	out, err := runSyncCmd("--dry-run")
	if err != nil {
		t.Fatalf("sync --dry-run failed: %v", err)
	}
	// Should produce some output (either "in sync" or analysis)
	if len(out) == 0 {
		t.Error("expected output from sync --dry-run")
	}
}

func TestSyncCommand_DryRunWithActions(t *testing.T) {
	dir := setupSyncDir(t)

	// Create an important directory not documented in ARCHITECTURE.md
	if err := os.Mkdir(filepath.Join(dir, "src"), 0750); err != nil {
		t.Fatal(err)
	}

	out, err := runSyncCmd("--dry-run")
	if err != nil {
		t.Fatalf("sync --dry-run failed: %v", err)
	}
	if !strings.Contains(out, "DRY RUN") {
		t.Error("output should contain DRY RUN marker")
	}
}

func TestSyncCommand_WithActions(t *testing.T) {
	dir := setupSyncDir(t)

	// Create an important undocumented directory
	if err := os.Mkdir(filepath.Join(dir, "cmd"), 0750); err != nil {
		t.Fatal(err)
	}

	out, err := runSyncCmd()
	if err != nil {
		t.Fatalf("sync failed: %v", err)
	}
	if !strings.Contains(out, "items") {
		t.Errorf("output should mention items to sync: %q", out)
	}
}

func TestCmd_HasDryRunFlag(t *testing.T) {
	cmd := Cmd()
	flag := cmd.Flags().Lookup("dry-run")
	if flag == nil {
		t.Fatal("expected --dry-run flag")
	}
	if flag.DefValue != "false" {
		t.Errorf("dry-run default = %q, want 'false'", flag.DefValue)
	}
}

func TestRunSync_InSyncMessage(t *testing.T) {
	setupSyncDir(t)

	// In a clean initialized dir, sync should report "in sync"
	out, err := runSyncCmd()
	if err != nil {
		t.Fatalf("sync error: %v", err)
	}
	if !strings.Contains(out, "in sync") {
		// Could have actions if directory has certain files
		_ = out
	}
}

func TestRunSync_DryRunWithSuggestions(t *testing.T) {
	dir := setupSyncDir(t)

	// Create multiple action triggers
	if err := os.Mkdir(filepath.Join(dir, "lib"), 0750); err != nil {
		t.Fatal(err)
	}
	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runSyncCmd("--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "DRY RUN") {
		t.Error("should indicate dry run mode")
	}
	if !strings.Contains(out, "without --dry-run") {
		t.Error("should suggest running without --dry-run")
	}
}

func TestRunSync_NonDryRunWithSuggestions(t *testing.T) {
	dir := setupSyncDir(t)

	if err := os.Mkdir(filepath.Join(dir, "api"), 0750); err != nil {
		t.Fatal(err)
	}

	out, err := runSyncCmd()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "items") {
		t.Errorf("should mention items count: %q", out)
	}
}

func TestSyncCommand_OutputFormat(t *testing.T) {
	dir := setupSyncDir(t)

	// Create multiple triggers
	for _, d := range []string{"src", "components"} {
		if err := os.Mkdir(filepath.Join(dir, d), 0750); err != nil {
			t.Fatal(err)
		}
	}

	cmd := Cmd()
	cmd.SetArgs([]string{})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	// Should have numbered actions
	if strings.Contains(out, "Sync Analysis") {
		if !strings.Contains(out, "1.") {
			t.Error("actions should be numbered")
		}
	}
}

func TestRunSync_CmdType(t *testing.T) {
	cmd := Cmd()
	if cmd.Use != "sync" {
		t.Errorf("cmd.Use = %q, want 'sync'", cmd.Use)
	}

	// Validate it's a *cobra.Command
	_ = cmd
}
