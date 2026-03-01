//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"
	"os/exec"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestRootCmd(t *testing.T) {
	cmd := RootCmd()

	if cmd == nil {
		t.Fatal("RootCmd() returned nil")
	}

	if cmd.Use != "ctx" {
		t.Errorf("RootCmd().Use = %q, want %q", cmd.Use, "ctx")
	}

	if cmd.Short == "" {
		t.Error("RootCmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("RootCmd().Long is empty")
	}

	// Check global flags exist
	contextDirFlag := cmd.PersistentFlags().Lookup("context-dir")
	if contextDirFlag == nil {
		t.Error("--context-dir flag not found")
	}

	noColorFlag := cmd.PersistentFlags().Lookup("no-color")
	if noColorFlag == nil {
		t.Error("--no-color flag not found")
	}
}

func TestInitialize(t *testing.T) {
	root := RootCmd()
	cmd := Initialize(root)

	if cmd == nil {
		t.Fatal("Initialize() returned nil")
	}

	// Verify all expected subcommands are registered
	expectedCommands := []string{
		"init",
		"status",
		"load",
		"add",
		"complete",
		"agent",
		"drift",
		"sync",
		"compact",
		"decisions",
		"watch",
		"hook",
		"learnings",
		"tasks",
		"loop",
		"recall",
		"journal",
		"serve",
		"guide",
	}

	commands := make(map[string]bool)
	for _, c := range cmd.Commands() {
		commands[c.Use] = true
		// Handle commands with args in Use (e.g., "serve [directory]")
		for _, exp := range expectedCommands {
			if c.Name() == exp {
				commands[exp] = true
			}
		}
	}

	for _, exp := range expectedCommands {
		if !commands[exp] {
			t.Errorf("missing subcommand: %s", exp)
		}
	}
}

func TestRootCmdVersion(t *testing.T) {
	cmd := RootCmd()

	if cmd.Version == "" {
		t.Error("RootCmd().Version is empty")
	}
}

func TestRootCmdAllowOutsideCwdFlag(t *testing.T) {
	cmd := RootCmd()

	flag := cmd.PersistentFlags().Lookup("allow-outside-cwd")
	if flag == nil {
		t.Fatal("--allow-outside-cwd flag not found")
	}
	if flag.DefValue != "false" {
		t.Errorf("--allow-outside-cwd default = %q, want %q", flag.DefValue, "false")
	}
}

func TestRootCmdPersistentPreRun_NoColor(t *testing.T) {
	cmd := RootCmd()
	// Set --no-color and --allow-outside-cwd so boundary check doesn't fail
	cmd.SetArgs([]string{"--no-color", "--allow-outside-cwd"})

	// Add a dummy subcommand so Execute doesn't just print help
	dummy := &cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--no-color", "--allow-outside-cwd", "dummy"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !color.NoColor {
		t.Error("expected color.NoColor to be true after --no-color")
	}
}

func TestRootCmdPersistentPreRun_ContextDir(t *testing.T) {
	cmd := RootCmd()

	dummy := &cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--context-dir", "/tmp/test-ctx", "--allow-outside-cwd", "dummy"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	got := rc.ContextDir()
	if got != "/tmp/test-ctx" {
		t.Errorf("ContextDir() = %q, want %q", got, "/tmp/test-ctx")
	}
}

func TestRootCmdPersistentPreRun_DefaultFlags(t *testing.T) {
	// Test PersistentPreRun with default flags (no --context-dir, no --no-color)
	// --allow-outside-cwd needed since test cwd may not have .context
	cmd := RootCmd()

	dummy := &cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--allow-outside-cwd", "dummy"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestInitializeReturnsSameCommand(t *testing.T) {
	root := RootCmd()
	result := Initialize(root)
	if result != root {
		t.Error("Initialize() should return the same command pointer")
	}
}

func TestInitializeSubcommandCount(t *testing.T) {
	root := RootCmd()
	Initialize(root)

	// There should be at least 19 subcommands registered
	count := len(root.Commands())
	if count < 19 {
		t.Errorf("Initialize() registered %d subcommands, want at least 19", count)
	}
}

// TestRootCmdPersistentPreRun_BoundaryViolation tests that boundary validation
// causes a non-zero exit when --context-dir is outside cwd and
// --allow-outside-cwd is not set. We use a subprocess because the code
// under test calls os.Exit(1).
func TestRootCmdPersistentPreRun_BoundaryViolation(t *testing.T) {
	if os.Getenv("TEST_BOUNDARY_EXIT") == "1" {
		cmd := RootCmd()
		dummy := &cobra.Command{
			Use: "dummy",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		cmd.AddCommand(dummy)
		cmd.SetArgs([]string{"--context-dir", "/etc/not-inside-cwd", "dummy"})
		_ = cmd.Execute()
		// If we reach here, the boundary check didn't exit
		return
	}

	// Run this test in a subprocess with the env var set
	sub := exec.Command(os.Args[0], "-test.run=^TestRootCmdPersistentPreRun_BoundaryViolation$") //nolint:gosec // test re-exec pattern
	sub.Env = append(os.Environ(), "TEST_BOUNDARY_EXIT=1")
	err := sub.Run()
	if err == nil {
		t.Fatal("expected subprocess to exit with non-zero status")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() == 0 {
		t.Fatal("expected non-zero exit code from boundary violation")
	}
}
