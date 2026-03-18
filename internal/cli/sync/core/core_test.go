//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/context/load"
)

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

func TestDetectSyncActions_NoActions(t *testing.T) {
	tmpDir := setupSyncDir(t)

	ctx, err := load.Do("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	_ = tmpDir
	actions := DetectSyncActions(ctx)
	// Just verify it runs without error
	_ = actions
}

func TestCheckNewDirectories_ImportantDirs(t *testing.T) {
	tmpDir := setupSyncDir(t)

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	// Create important directories
	for _, d := range []string{"src", "lib", "pkg", "internal", "cmd", "api"} {
		if mkErr := os.Mkdir(filepath.Join(tmpDir, d), 0750); mkErr != nil {
			t.Fatal(mkErr)
		}
	}

	actions := CheckNewDirectories(ctx)
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
	tmpDir := setupSyncDir(t)

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	// Create directories that should be skipped
	for _, d := range []string{".git", "node_modules", "vendor", "dist", "build"} {
		if mkErr := os.Mkdir(filepath.Join(tmpDir, d), 0750); mkErr != nil {
			t.Fatal(mkErr)
		}
	}

	actions := CheckNewDirectories(ctx)
	for _, a := range actions {
		for _, skip := range []string{".git", "node_modules", "vendor", "dist", "build"} {
			if strings.Contains(a.Description, skip) {
				t.Errorf("should skip %q but got action: %s", skip, a.Description)
			}
		}
	}
}

func TestCheckNewDirectories_DocumentedDirsIgnored(t *testing.T) {
	tmpDir := setupSyncDir(t)

	// Write ARCHITECTURE.md that mentions "src"
	archPath := filepath.Join(tmpDir, dir.Context, ctx.Architecture)
	if err := os.WriteFile(archPath, []byte("# Architecture\n\nThe src directory contains...\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	if mkErr := os.Mkdir(filepath.Join(tmpDir, "src"), 0750); mkErr != nil {
		t.Fatal(mkErr)
	}

	actions := CheckNewDirectories(ctx)
	for _, a := range actions {
		if strings.Contains(a.Description, "'src/'") {
			t.Error("documented directory 'src' should not produce an action")
		}
	}
}

func TestCheckPackageFiles_NoPackages(t *testing.T) {
	setupSyncDir(t)

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := CheckPackageFiles(ctx)
	if len(actions) != 0 {
		t.Errorf("expected no actions, got %d", len(actions))
	}
}

func TestCheckPackageFiles_WithPackageFile(t *testing.T) {
	tmpDir := setupSyncDir(t)

	// Remove any existing dependency docs so the check triggers
	archPath := filepath.Join(tmpDir, dir.Context, ctx.Architecture)
	_ = os.WriteFile(archPath, []byte("# Architecture\n\nSimple app.\n"), 0600)
	depsPath := filepath.Join(tmpDir, dir.Context, ctx.Dependency)
	_ = os.Remove(depsPath)

	// Create a package.json
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name":"test"}`), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := CheckPackageFiles(ctx)
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
	tmpDir := setupSyncDir(t)

	// Create a package.json and DEPENDENCIES.md
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name":"test"}`), 0600); err != nil {
		t.Fatal(err)
	}
	depsPath := filepath.Join(tmpDir, dir.Context, ctx.Dependency)
	if err := os.WriteFile(depsPath, []byte("# Dependencies\n\nAll documented.\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := CheckPackageFiles(ctx)
	for _, a := range actions {
		if a.Type == "DEPS" && strings.Contains(a.Description, "package.json") {
			t.Error("should not produce DEPS action when DEPENDENCIES.md exists")
		}
	}
}

func TestCheckConfigFiles_NoConfigs(t *testing.T) {
	tmpDir := setupSyncDir(t)

	// Remove Makefile created by init (it matches the Makefile config pattern)
	_ = os.Remove(filepath.Join(tmpDir, "Makefile"))
	_ = os.Remove(filepath.Join(tmpDir, "Makefile.ctx"))

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := CheckConfigFiles(ctx)
	// With Makefile removed, no config files should match
	if len(actions) != 0 {
		t.Errorf("expected no actions for clean dir, got %d", len(actions))
	}
}

func TestCheckConfigFiles_WithConfigFile(t *testing.T) {
	tmpDir := setupSyncDir(t)

	// Create a tsconfig.json
	if err := os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := CheckConfigFiles(ctx)
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
	tmpDir := setupSyncDir(t)

	// Create tsconfig.json
	if err := os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	// Write CONVENTIONS.md mentioning tsconfig
	convPath := filepath.Join(tmpDir, dir.Context, ctx.Convention)
	if err := os.WriteFile(convPath, []byte("# Conventions\n\ntsconfig.json is configured for strict mode.\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := CheckConfigFiles(ctx)
	for _, a := range actions {
		if a.Type == "CONFIG" && strings.Contains(a.Description, "tsconfig") {
			t.Error("tsconfig should not produce an action when documented in CONVENTIONS.md")
		}
	}
}

func TestCheckPackageFiles_ArchContainsDependencies(t *testing.T) {
	tmpDir := setupSyncDir(t)

	// Create a go.mod
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Write ARCHITECTURE.md that mentions "dependencies"
	archPath := filepath.Join(tmpDir, dir.Context, ctx.Architecture)
	if err := os.WriteFile(archPath, []byte("# Architecture\n\nProject dependencies are managed via go.mod.\n"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := CheckPackageFiles(ctx)
	for _, a := range actions {
		if a.Type == "DEPS" && strings.Contains(a.Description, "go.mod") {
			t.Error("should not produce DEPS action when ARCHITECTURE.md mentions dependencies")
		}
	}
}

func TestAction_Fields(t *testing.T) {
	a := Action{
		Type:        "NEW_DIR",
		File:        ctx.Architecture,
		Description: "test description",
		Suggestion:  "test suggestion",
	}
	if a.Type != "NEW_DIR" || a.File != ctx.Architecture {
		t.Error("action fields should be set correctly")
	}
}

func TestRunSync_ActionWithEmptySuggestion(t *testing.T) {
	tmpDir := setupSyncDir(t)

	// Create important dir to trigger actions
	if err := os.Mkdir(filepath.Join(tmpDir, "services"), 0750); err != nil {
		t.Fatal(err)
	}

	ctx, err := load.Do("")
	if err != nil {
		t.Fatal(err)
	}

	actions := DetectSyncActions(ctx)
	for _, a := range actions {
		// All actions should have a non-empty Description
		if a.Description == "" {
			t.Error("action should have a description")
		}
	}
}
