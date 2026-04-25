//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/resolve"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
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
}

// TestRoot_NoContextDirFlag is the regression guard for the
// removed --context-dir flag (spec:
// specs/single-source-context-anchor.md). Cobra must reject
// the flag with its standard "unknown flag" error.
func TestRoot_NoContextDirFlag(t *testing.T) {
	cmd := RootCmd()
	cmd.SetOut(&discardWriter{})
	cmd.SetErr(&discardWriter{})
	cmd.SetArgs([]string{"--context-dir=/tmp", "status"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for removed --context-dir flag")
	}
	if !strings.Contains(err.Error(), "unknown flag") {
		t.Errorf("error = %q, want cobra unknown-flag error", err.Error())
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
		"setup",
		"learning",
		"task",
		"loop",
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

// TestRootCmdPersistentPreRun_CtxDirEnv: CTX_DIR env declares the
// context directory; non-init annotated dummy bypasses the
// initialized check.
func TestRootCmdPersistentPreRun_CtxDirEnv(t *testing.T) {
	tmp := t.TempDir()
	ctxDir := filepath.Join(tmp, dir.Context)
	if err := os.MkdirAll(ctxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(env.CtxDir, ctxDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	cmd := RootCmd()

	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		Run:         func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	got, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		t.Fatalf("ContextDir: %v", ctxErr)
	}
	if got != ctxDir {
		t.Errorf("ContextDir() = %q, want %q", got, ctxDir)
	}
}

func TestRootCmdPersistentPreRun_DefaultFlags(t *testing.T) {
	// Test PersistentPreRun with default flags.
	// The dummy command carries AnnotationSkipInit, so PersistentPreRunE
	// skips the context-dir declaration gate and returns immediately.
	cmd := RootCmd()

	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		Run:         func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

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

func TestInitGuard_BlocksUninitializedCommand(t *testing.T) {
	tmp := t.TempDir()
	ctxDir := filepath.Join(tmp, dir.Context)
	if err := os.MkdirAll(ctxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(env.CtxDir, ctxDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

	execErr := cmd.Execute()
	if execErr == nil {
		t.Fatal("expected error for uninitialized context directory")
	}
	want := `ctx: not initialized - run "ctx init" first`
	if got := execErr.Error(); got != want {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestInitGuard_AllowsAnnotatedCommand(t *testing.T) {
	tmp := t.TempDir() // empty - not initialized
	ctxDir := filepath.Join(tmp, dir.Context)
	if err := os.MkdirAll(ctxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(env.CtxDir, ctxDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		Run:         func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("annotated command should succeed: %v", execErr)
	}
}

func TestInitGuard_AllowsHiddenCommand(t *testing.T) {
	tmp := t.TempDir() // empty - not initialized
	ctxDir := filepath.Join(tmp, dir.Context)
	if err := os.MkdirAll(ctxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(env.CtxDir, ctxDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:    "dummy",
		Hidden: true,
		Run:    func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

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
	cmd.SetArgs([]string{"group"})

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
	ctxDir := filepath.Join(tmp, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0o700); mkErr != nil {
		t.Fatal(mkErr)
	}

	// Create required context files so Initialized() returns true.
	for _, f := range ctx.FilesRequired {
		path := filepath.Join(ctxDir, f)
		content := []byte("# " + f + "\n")
		if writeErr := os.WriteFile(path, content, 0o600); writeErr != nil {
			t.Fatalf("setup: %v", writeErr)
		}
	}

	t.Setenv(env.CtxDir, ctxDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("initialized command should succeed: %v", execErr)
	}
}

func TestRootCmdToolFlag(t *testing.T) {
	cmd := RootCmd()

	f := cmd.PersistentFlags().Lookup(flag.Tool)
	if f == nil {
		t.Fatal("--tool flag not found")
	}
	if f.DefValue != "" {
		t.Errorf("--tool default = %q, want empty", f.DefValue)
	}
}

func TestResolveTool_FlagOverridesRC(t *testing.T) {
	rc.Reset()
	t.Cleanup(func() { rc.Reset() })

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			tool, err := resolve.Tool(cmd)
			if err != nil {
				return err
			}
			if tool != "cursor" {
				t.Errorf("resolve.Tool() = %q, want %q", tool, "cursor")
			}
			return nil
		},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"--tool", "cursor", "dummy"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestResolveTool_FallsBackToRC(t *testing.T) {
	// When --tool is not set, ResolveTool falls back to rc.Tool().
	// With a fresh rc (no .ctxrc), rc.Tool() returns "" so this
	// should return an error. We test the fallback path indirectly.
	rc.Reset()
	t.Cleanup(func() { rc.Reset() })

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolve.Tool(cmd)
			if err == nil {
				t.Error("resolve.Tool() should return error when no tool is set")
			}
			return nil // swallow error so Execute succeeds
		},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestResolveTool_ErrorMessage(t *testing.T) {
	rc.Reset()
	t.Cleanup(func() { rc.Reset() })

	cmd := RootCmd()
	dummy := &cobra.Command{
		Use:         "dummy",
		Annotations: map[string]string{cli.AnnotationSkipInit: "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolve.Tool(cmd)
			if err == nil {
				t.Fatal("expected error")
			}
			want := "no tool specified: use --tool <tool> or set the tool field in .ctxrc"
			if err.Error() != want {
				t.Errorf("error = %q, want %q", err.Error(), want)
			}
			return nil
		},
	}
	cmd.AddCommand(dummy)
	cmd.SetArgs([]string{"dummy"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}
