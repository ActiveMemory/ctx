//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "serve [directory]" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "serve [directory]")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}

	if cmd.RunE == nil {
		t.Error("Cmd().RunE is nil")
	}
}

func TestCmd_AcceptsArgs(t *testing.T) {
	cmd := Cmd()

	// Should accept 0 or 1 args
	if err := cmd.Args(cmd, []string{}); err != nil {
		t.Errorf("should accept 0 args: %v", err)
	}

	if err := cmd.Args(cmd, []string{"./docs"}); err != nil {
		t.Errorf("should accept 1 arg: %v", err)
	}

	if err := cmd.Args(cmd, []string{"a", "b"}); err == nil {
		t.Error("should reject 2 args")
	}
}

func TestRunServe_DirNotFound(t *testing.T) {
	err := runServe([]string{"/tmp/nonexistent-dir-ctx-test-xyz"})
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
	if !strings.Contains(err.Error(), "directory not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunServe_NotADirectory(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ctx-serve-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_ = tmpFile.Close()

	serveErr := runServe([]string{tmpFile.Name()})
	if serveErr == nil {
		t.Fatal("expected error for non-directory path")
	}
	if !strings.Contains(serveErr.Error(), "not a directory") {
		t.Errorf("unexpected error: %v", serveErr)
	}
}

func TestRunServe_NoSiteConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ctx-serve-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	serveErr := runServe([]string{tmpDir})
	if serveErr == nil {
		t.Fatal("expected error for missing zensical.toml")
	}
	if !strings.Contains(serveErr.Error(), "no zensical.toml found") {
		t.Errorf("unexpected error: %v", serveErr)
	}
}

func TestRunServe_ZensicalNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ctx-serve-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a zensical.toml so we pass the config check
	tomlPath := filepath.Join(tmpDir, config.FileZensicalToml)
	if err := os.WriteFile(tomlPath, []byte("[site]\n"), 0600); err != nil {
		t.Fatalf("failed to create zensical.toml: %v", err)
	}

	// Ensure zensical is not in PATH
	t.Setenv("PATH", "")

	serveErr := runServe([]string{tmpDir})
	if serveErr == nil {
		t.Fatal("expected error for missing zensical binary")
	}
	if !strings.Contains(serveErr.Error(), "zensical not found") {
		t.Errorf("unexpected error: %v", serveErr)
	}
}

func TestRunServe_DefaultDir(t *testing.T) {
	// When no args are given, runServe uses the default journal-site directory
	// which won't exist in test, so we expect directory not found
	err := runServe([]string{})
	if err == nil {
		t.Fatal("expected error when default dir doesn't exist")
	}
	// Should be either "directory not found" or "not a directory"
	if !strings.Contains(err.Error(), "directory not found") &&
		!strings.Contains(err.Error(), "not a directory") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestErrDirNotFound(t *testing.T) {
	err := errDirNotFound("/some/path")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "directory not found") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "/some/path") {
		t.Errorf("error should contain path, got: %v", err)
	}
}

func TestErrNotDir(t *testing.T) {
	err := errNotDir("/some/file")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "/some/file") {
		t.Errorf("error should contain path, got: %v", err)
	}
}

func TestErrNoSiteConfig(t *testing.T) {
	err := errNoSiteConfig("/some/dir")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "no zensical.toml found") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "/some/dir") {
		t.Errorf("error should contain dir, got: %v", err)
	}
}

func TestRunServe_WithMockZensical(t *testing.T) {
	// Create a temporary directory with a zensical.toml
	tmpDir, err := os.MkdirTemp("", "ctx-serve-mock-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tomlPath := filepath.Join(tmpDir, config.FileZensicalToml)
	if err := os.WriteFile(tomlPath, []byte("[site]\n"), 0600); err != nil {
		t.Fatalf("failed to create zensical.toml: %v", err)
	}

	// Create a fake zensical binary that just exits successfully
	binDir := filepath.Join(tmpDir, "bin")
	if err := os.MkdirAll(binDir, 0750); err != nil {
		t.Fatalf("failed to create bin dir: %v", err)
	}
	fakeZensical := filepath.Join(binDir, "zensical")
	if err := os.WriteFile(fakeZensical, []byte("#!/bin/sh\nexit 0\n"), 0600); err != nil {
		t.Fatalf("failed to create fake zensical: %v", err)
	}
	if err := os.Chmod(fakeZensical, 0755); err != nil { //nolint:gosec // test needs executable
		t.Fatalf("failed to chmod fake zensical: %v", err)
	}

	// Set PATH to include our fake binary
	origPath := os.Getenv("PATH")
	t.Setenv("PATH", binDir+":"+origPath)

	serveErr := runServe([]string{tmpDir})
	if serveErr != nil {
		t.Errorf("unexpected error: %v", serveErr)
	}
}

func TestCmd_RunE(t *testing.T) {
	// Test that Cmd().RunE actually invokes runServe via the command
	cmd := Cmd()
	// Set args to a nonexistent dir so we get a predictable error
	cmd.SetArgs([]string{"/tmp/nonexistent-ctx-test-xyz"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from Execute with nonexistent dir")
	}
	if !strings.Contains(err.Error(), "directory not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestErrZensicalNotFound(t *testing.T) {
	err := errZensicalNotFound()
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "zensical not found") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "pipx install zensical") {
		t.Errorf("error should contain install instructions, got: %v", err)
	}
}
