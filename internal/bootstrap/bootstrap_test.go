//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
)

// discardWriter silences command output in tests.
type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

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
	contextDirFlag := cmd.PersistentFlags().Lookup(flag.ContextDir)
	if contextDirFlag == nil {
		t.Error("--context-dir flag not found")
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
		"agent",
		"drift",
		"sync",
		"compact",
		"decision",
		"watch",
		"hook",
		"learning",
		"task",
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

	flag := cmd.PersistentFlags().Lookup(flag.AllowOutsideCwd)
	if flag == nil {
		t.Fatal("--allow-outside-cwd flag not found")
	}
	if flag.DefValue != "false" {
		t.Errorf("--allow-outside-cwd default = %q, want %q", flag.DefValue, "false")
	}
}

func TestRootCmdPersistentPreRun_ContextDir(t *testing.T) {
	cmd := RootCmd()

	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		Run:         func(cmd *cobra.Command, args []string) {},
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
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		Run:         func(cmd *cobra.Command, args []string) {},
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
// returns an error when --context-dir is outside cwd and --allow-outside-cwd
// is not set.
func TestRootCmdPersistentPreRun_BoundaryViolation(t *testing.T) {
	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		Run:         func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--context-dir", "/etc/not-inside-cwd", "dummy"})

	execErr := cmd.Execute()
	if execErr == nil {
		t.Fatal("expected error from boundary violation")
	}
}

func TestInitGuard_BlocksUninitializedCommand(t *testing.T) {
	tmp := t.TempDir()

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--context-dir", tmp, "--allow-outside-cwd", "dummy"})

	execErr := cmd.Execute()
	if execErr == nil {
		t.Fatal("expected error for uninitialized context directory")
	}
	if got := execErr.Error(); got != `ctx: not initialized - run "ctx init" first` {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestInitGuard_AllowsAnnotatedCommand(t *testing.T) {
	tmp := t.TempDir() // empty - not initialized

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		Run:         func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--context-dir", tmp, "--allow-outside-cwd", "dummy"})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("annotated command should succeed: %v", execErr)
	}
}

func TestInitGuard_AllowsHiddenCommand(t *testing.T) {
	tmp := t.TempDir() // empty - not initialized

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:    "dummy",
		Hidden: true,
		Run:    func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--context-dir", tmp, "--allow-outside-cwd", "dummy"})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("hidden command should succeed: %v", execErr)
	}
}

func TestInitGuard_AllowsGroupingCommand(t *testing.T) {
	cmd := RootCmd()
	// Grouping command: no Run or RunE - just shows help.
	group := &cobra.Command{
		Use:   "group",
		Short: "A grouping command",
	}
	cmd.AddCommand(group)
	cmd.SetArgs([]string{"--allow-outside-cwd", "group"})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("grouping command should succeed: %v", execErr)
	}
}

func TestInitGuard_AllowsCompletionSubcommand(t *testing.T) {
	tmp := t.TempDir() // empty - not initialized
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	rc.Reset()
	t.Cleanup(func() { rc.Reset() })

	cmd := RootCmd()
	Initialize(cmd)

	// "completion bash" is added by cobra during Execute; simulate by
	// running the full command.
	cmd.SetArgs([]string{"completion", "bash"})
	cmd.SetOut(&discardWriter{})
	cmd.SetErr(&discardWriter{})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("completion bash should succeed without init: %v", execErr)
	}
}

func TestInitGuard_AllowsInitializedCommand(t *testing.T) {
	tmp := t.TempDir()

	// Create required context files so Initialized() returns true.
	for _, f := range ctx.FilesRequired {
		path := filepath.Join(tmp, f)
		if writeErr := os.WriteFile(path, []byte("# "+f+"\n"), 0o600); writeErr != nil {
			t.Fatalf("setup: %v", writeErr)
		}
	}

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--context-dir", tmp, "--allow-outside-cwd", "dummy"})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("initialized command should succeed: %v", execErr)
	}
}
