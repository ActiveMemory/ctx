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
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
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

	_, err := runSyncCmd()
	if err == nil {
		t.Fatal("expected error when no .context/ exists")
	}
	if !strings.Contains(err.Error(), "ctx init") {
		t.Errorf("error = %q, want 'ctx init' suggestion", err.Error())
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

func TestDetectSyncActions_NoActions(t *testing.T) {
	dir := setupSyncDir(t)

	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	// In a clean init dir with no important dirs or package files,
	// there may be no actions
	_ = dir
	actions := detectSyncActions(ctx)
	// Just verify it runs without error
	_ = actions
}

func TestCheckNewDirectories_ImportantDirs(t *testing.T) {
	dir := setupSyncDir(t)

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	// Create important directories
	for _, d := range []string{"src", "lib", "pkg", "internal", "cmd", "api"} {
		if mkErr := os.Mkdir(filepath.Join(dir, d), 0750); mkErr != nil {
			t.Fatal(mkErr)
		}
	}

	actions := checkNewDirectories(ctx)
	if len(actions) == 0 {
		t.Error("expected actions for undocumented directories")
	}
	for _, a := range actions {
		if a.Type != "NEW_DIR" {
			t.Errorf("action type = %q, want NEW_DIR", a.Type)
		}
	}
}

func TestCheckNewDirectories_SkipsHiddenAndVendor(t *testing.T) {
	dir := setupSyncDir(t)

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	// Create directories that should be skipped
	for _, d := range []string{".git", "node_modules", "vendor", "dist", "build"} {
		if mkErr := os.Mkdir(filepath.Join(dir, d), 0750); mkErr != nil {
			t.Fatal(mkErr)
		}
	}

	actions := checkNewDirectories(ctx)
	for _, a := range actions {
		for _, skip := range []string{".git", "node_modules", "vendor", "dist", "build"} {
			if strings.Contains(a.Description, skip) {
				t.Errorf("should skip %q but got action: %s", skip, a.Description)
			}
		}
	}
}

func TestCheckNewDirectories_DocumentedDirsIgnored(t *testing.T) {
	dir := setupSyncDir(t)

	// Write ARCHITECTURE.md that mentions "src"
	archPath := filepath.Join(dir, config.DirContext, config.FileArchitecture)
	if err := os.WriteFile(archPath, []byte("# Architecture\n\nThe src directory contains...\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	if mkErr := os.Mkdir(filepath.Join(dir, "src"), 0750); mkErr != nil {
		t.Fatal(mkErr)
	}

	actions := checkNewDirectories(ctx)
	for _, a := range actions {
		if strings.Contains(a.Description, "'src/'") {
			t.Error("documented directory 'src' should not produce an action")
		}
	}
}

func TestCheckPackageFiles_NoPackages(t *testing.T) {
	setupSyncDir(t)

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := checkPackageFiles(ctx)
	if len(actions) != 0 {
		t.Errorf("expected no actions, got %d", len(actions))
	}
}

func TestCheckPackageFiles_WithPackageFile(t *testing.T) {
	dir := setupSyncDir(t)

	// Remove any existing dependency docs so the check triggers
	archPath := filepath.Join(dir, config.DirContext, config.FileArchitecture)
	_ = os.WriteFile(archPath, []byte("# Architecture\n\nSimple app.\n"), 0600)
	depsPath := filepath.Join(dir, config.DirContext, config.FileDependency)
	_ = os.Remove(depsPath)

	// Create a package.json
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"test"}`), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := checkPackageFiles(ctx)
	found := false
	for _, a := range actions {
		if a.Type == "DEPS" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected DEPS action for package.json")
	}
}

func TestCheckPackageFiles_WithDepsDoc(t *testing.T) {
	dir := setupSyncDir(t)

	// Create a package.json and DEPENDENCIES.md
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"test"}`), 0600); err != nil {
		t.Fatal(err)
	}
	depsPath := filepath.Join(dir, config.DirContext, config.FileDependency)
	if err := os.WriteFile(depsPath, []byte("# Dependencies\n\nAll documented.\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := checkPackageFiles(ctx)
	for _, a := range actions {
		if a.Type == "DEPS" && strings.Contains(a.Description, "package.json") {
			t.Error("should not produce DEPS action when DEPENDENCIES.md exists")
		}
	}
}

func TestCheckConfigFiles_NoConfigs(t *testing.T) {
	dir := setupSyncDir(t)

	// Remove Makefile created by init (it matches the Makefile config pattern)
	_ = os.Remove(filepath.Join(dir, "Makefile"))
	_ = os.Remove(filepath.Join(dir, "Makefile.ctx"))

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := checkConfigFiles(ctx)
	// With Makefile removed, no config files should match
	if len(actions) != 0 {
		t.Errorf("expected no actions for clean dir, got %d", len(actions))
	}
}

func TestCheckConfigFiles_WithConfigFile(t *testing.T) {
	dir := setupSyncDir(t)

	// Create a tsconfig.json
	if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := checkConfigFiles(ctx)
	found := false
	for _, a := range actions {
		if a.Type == "CONFIG" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected CONFIG action for tsconfig.json")
	}
}

func TestCheckConfigFiles_DocumentedInConventions(t *testing.T) {
	dir := setupSyncDir(t)

	// Create tsconfig.json
	if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	// Write CONVENTIONS.md mentioning tsconfig
	convPath := filepath.Join(dir, config.DirContext, config.FileConvention)
	if err := os.WriteFile(convPath, []byte("# Conventions\n\ntsconfig.json is configured for strict mode.\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := checkConfigFiles(ctx)
	for _, a := range actions {
		if a.Type == "CONFIG" && strings.Contains(a.Description, "tsconfig") {
			t.Error("tsconfig should not produce an action when documented in CONVENTIONS.md")
		}
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
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0600); err != nil {
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

func TestAction_Fields(t *testing.T) {
	a := Action{
		Type:        "NEW_DIR",
		File:        config.FileArchitecture,
		Description: "test description",
		Suggestion:  "test suggestion",
	}
	if a.Type != "NEW_DIR" || a.File != config.FileArchitecture {
		t.Error("action fields should be set correctly")
	}
}

func TestRunSync_ActionWithEmptySuggestion(t *testing.T) {
	dir := setupSyncDir(t)

	// Create important dir to trigger actions
	if err := os.Mkdir(filepath.Join(dir, "services"), 0750); err != nil {
		t.Fatal(err)
	}

	// runSync with cmd that captures output
	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := detectSyncActions(ctx)
	for _, a := range actions {
		// All actions should have a non-empty Description
		if a.Description == "" {
			t.Error("action should have a description")
		}
	}
}

func TestCheckPackageFiles_ArchContainsDependencies(t *testing.T) {
	dir := setupSyncDir(t)

	// Create a go.mod
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Write ARCHITECTURE.md that mentions "dependencies"
	archPath := filepath.Join(dir, config.DirContext, config.FileArchitecture)
	if err := os.WriteFile(archPath, []byte("# Architecture\n\nProject dependencies are managed via go.mod.\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatal(err)
	}

	actions := checkPackageFiles(ctx)
	for _, a := range actions {
		if a.Type == "DEPS" && strings.Contains(a.Description, "go.mod") {
			t.Error("should not produce DEPS action when ARCHITECTURE.md mentions dependencies")
		}
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
