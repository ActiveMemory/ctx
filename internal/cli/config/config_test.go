//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCopyProfile_MissingSource verifies error on nonexistent source file.
func TestCopyProfile_MissingSource(t *testing.T) {
	root := t.TempDir()

	copyErr := copyProfile(root, ".ctxrc.nonexistent")
	if copyErr == nil {
		t.Fatal("expected error for missing source profile")
	}
}

// TestCopyProfile_Success verifies content is copied to .ctxrc.
func TestCopyProfile_Success(t *testing.T) {
	root := t.TempDir()

	srcContent := "# test profile\nnotify:\n  events:\n    - loop\n"
	srcFile := ".ctxrc.test"
	if writeErr := os.WriteFile(
		filepath.Join(root, srcFile), []byte(srcContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	if copyErr := copyProfile(root, srcFile); copyErr != nil {
		t.Fatalf("copyProfile failed: %v", copyErr)
	}

	data, readErr := os.ReadFile(filepath.Join(root, fileCtxRC))
	if readErr != nil {
		t.Fatalf("expected .ctxrc to exist: %v", readErr)
	}

	if string(data) != srcContent {
		t.Errorf("expected .ctxrc content to match source, got: %s", string(data))
	}
}

// TestCmd_HasSubcommands verifies the config command includes expected subcommands.
func TestCmd_HasSubcommands(t *testing.T) {
	cmd := Cmd()

	expected := map[string]bool{
		"switch": false,
		"status": false,
		"schema": false,
	}

	for _, sub := range cmd.Commands() {
		if _, ok := expected[sub.Name()]; ok {
			expected[sub.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("config command should have %q subcommand", name)
		}
	}
}
