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
)

func TestValidatePath_Valid(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	typeDir := filepath.Join(hooksDir, "pre-tool-use")

	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	script := filepath.Join(typeDir, "check.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := ValidatePath(hooksDir, script); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidatePath_Symlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	typeDir := filepath.Join(hooksDir, "pre-tool-use")

	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a real file and a symlink to it.
	real := filepath.Join(dir, "real.sh")
	if err := os.WriteFile(real, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	link := filepath.Join(typeDir, "link.sh")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}

	err := ValidatePath(hooksDir, link)
	if err == nil {
		t.Fatal("expected symlink error, got nil")
	}

	if got := err.Error(); !contains(got, "symlink") {
		t.Fatalf("expected symlink error, got %q", got)
	}
}

func TestValidatePath_BoundaryEscape(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	outsideDir := filepath.Join(dir, "outside")

	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outsideDir, 0o755); err != nil {
		t.Fatal(err)
	}

	script := filepath.Join(outsideDir, "evil.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	err := ValidatePath(hooksDir, script)
	if err == nil {
		t.Fatal("expected boundary error, got nil")
	}

	if got := err.Error(); !contains(got, "escapes") {
		t.Fatalf("expected boundary error, got %q", got)
	}
}

func TestValidatePath_NotExecutable(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	typeDir := filepath.Join(hooksDir, "session-start")

	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	script := filepath.Join(typeDir, "noexec.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidatePath(hooksDir, script)
	if err == nil {
		t.Fatal("expected not-executable error, got nil")
	}

	if got := err.Error(); !contains(got, "not executable") {
		t.Fatalf("expected not-executable error, got %q", got)
	}
}

func TestValidatePath_NonExistent(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := ValidatePath(hooksDir, filepath.Join(hooksDir, "missing.sh"))
	if err == nil {
		t.Fatal("expected error for non-existent path, got nil")
	}
}

// contains is a small helper to check substring presence.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
