//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	cfgTrigger "github.com/ActiveMemory/ctx/internal/config/trigger"
)

// TestDiscover_ValidExecutableScripts verifies that Discover returns
// executable scripts grouped by hook type.
// Validates: Requirements 6.1
func TestDiscover_ValidExecutableScripts(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	// Create two hook type directories with executable scripts.
	for _, ht := range []string{"pre-tool-use", "session-start"} {
		typeDir := filepath.Join(hooksDir, ht)
		if err := os.MkdirAll(typeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(typeDir, "check.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	result, err := Discover(hooksDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result[cfgTrigger.PreToolUse]) != 1 {
		t.Fatalf("expected 1 pre-tool-use hook, got %d", len(result[cfgTrigger.PreToolUse]))
	}
	if result[cfgTrigger.PreToolUse][0].Name != "check" {
		t.Errorf("expected name %q, got %q", "check", result[cfgTrigger.PreToolUse][0].Name)
	}
	if !result[cfgTrigger.PreToolUse][0].Enabled {
		t.Error("expected hook to be enabled")
	}

	if len(result[cfgTrigger.SessionStart]) != 1 {
		t.Fatalf("expected 1 session-start hook, got %d", len(result[cfgTrigger.SessionStart]))
	}
}

// TestDiscover_SkipsNonExecutable verifies that Discover includes
// non-executable scripts in results but marks them as disabled.
// Note: Discover calls ValidateHookPath which rejects non-executable
// files, so they are skipped entirely.
// Validates: Requirements 6.2
func TestDiscover_SkipsNonExecutable(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	typeDir := filepath.Join(hooksDir, "pre-tool-use")

	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create one executable and one non-executable script.
	if err := os.WriteFile(filepath.Join(typeDir, "enabled.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(typeDir, "disabled.sh"), []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(hooksDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hooks := result[cfgTrigger.PreToolUse]
	if len(hooks) != 1 {
		t.Fatalf("expected 1 hook (non-executable skipped by validation), got %d", len(hooks))
	}
	if hooks[0].Name != "enabled" {
		t.Errorf("expected surviving hook name %q, got %q", "enabled", hooks[0].Name)
	}
}

// TestDiscover_SkipsSymlinks verifies that Discover skips symlinked scripts.
// Validates: Requirements 6.1, 15.1
func TestDiscover_SkipsSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	typeDir := filepath.Join(hooksDir, "post-tool-use")

	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a real executable script.
	real := filepath.Join(dir, "real.sh")
	if err := os.WriteFile(real, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a symlink inside the hook type directory.
	link := filepath.Join(typeDir, "link.sh")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}

	// Also create a valid executable script to confirm it's still found.
	if err := os.WriteFile(filepath.Join(typeDir, "valid.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(hooksDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hooks := result[cfgTrigger.PostToolUse]
	if len(hooks) != 1 {
		t.Fatalf("expected 1 hook (symlink skipped), got %d", len(hooks))
	}
	if hooks[0].Name != "valid" {
		t.Errorf("expected hook name %q, got %q", "valid", hooks[0].Name)
	}
}

// TestDiscover_MissingHooksDir verifies that Discover returns an empty
// map without error when the hooks directory does not exist.
// Validates: Requirements 6.4
func TestDiscover_MissingHooksDir(t *testing.T) {
	result, err := Discover(filepath.Join(t.TempDir(), "nonexistent"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(result))
	}
}

// TestDiscover_AlphabeticalOrder verifies that hooks within each type
// are sorted alphabetically by name.
// Validates: Requirements 6.3
func TestDiscover_AlphabeticalOrder(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	typeDir := filepath.Join(hooksDir, "file-save")

	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create scripts in non-alphabetical order.
	for _, name := range []string{"charlie.sh", "alpha.sh", "bravo.sh"} {
		if err := os.WriteFile(filepath.Join(typeDir, name), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	result, err := Discover(hooksDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hooks := result[cfgTrigger.FileSave]
	if len(hooks) != 3 {
		t.Fatalf("expected 3 hooks, got %d", len(hooks))
	}

	expected := []string{"alpha", "bravo", "charlie"}
	for i, want := range expected {
		if hooks[i].Name != want {
			t.Errorf("hooks[%d].Name = %q, want %q", i, hooks[i].Name, want)
		}
	}
}

// TestFindByName_Found verifies that FindByName locates a hook by name
// across type directories.
// Validates: Requirements 6.1
func TestFindByName_Found(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	typeDir := filepath.Join(hooksDir, "context-add")

	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(typeDir, "notify.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	info, err := FindByName(hooksDir, "notify")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("expected hook info, got nil")
	}
	if info.Name != "notify" {
		t.Errorf("expected name %q, got %q", "notify", info.Name)
	}
	if info.Type != cfgTrigger.ContextAdd {
		t.Errorf("expected type %q, got %q", cfgTrigger.ContextAdd, info.Type)
	}
}

// TestFindByName_NotFound verifies that FindByName returns nil when
// no hook matches the given name.
// Validates: Requirements 6.1
func TestFindByName_NotFound(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}

	info, err := FindByName(hooksDir, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info != nil {
		t.Fatalf("expected nil, got %+v", info)
	}
}
