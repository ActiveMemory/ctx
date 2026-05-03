//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// TestLearningAdd verifies the noun-first ctx learning add
// subcommand writes a structured entry to LEARNINGS.md.
func TestLearningAdd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-learning-add-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{
		"Always check for nil before dereferencing",
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
		"--context", "Got a nil pointer panic in production",
		"--lesson", "Always validate pointers before use",
		"--application", "Add nil checks in all pointer-receiving functions",
	})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx learning add failed: %v", err)
	}

	content, err := os.ReadFile(".context/LEARNINGS.md")
	if err != nil {
		t.Fatalf("failed to read LEARNINGS.md: %v", err)
	}
	contentStr := string(content)
	for _, want := range []string{
		"Always check for nil before dereferencing",
		"Got a nil pointer panic in production",
		"Always validate pointers before use",
		"Add nil checks in all pointer-receiving functions",
	} {
		if !strings.Contains(contentStr, want) {
			t.Errorf("expected %q in LEARNINGS.md", want)
		}
	}
}

// TestLearningAddRequiresFlags verifies that omitting
// required flags produces an error.
func TestLearningAddRequiresFlags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-learning-add-req-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{"Incomplete learning"})
	err = addCmd.Execute()
	if err == nil {
		t.Fatal("expected error when adding learning without required flags")
	}
	if !strings.Contains(err.Error(), "--session-id") {
		t.Errorf("error should mention missing --session-id flag: %v", err)
	}
}

// TestLearningAddFromFile verifies reading content from --file.
func TestLearningAddFromFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-learning-add-file-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	contentFile := filepath.Join(tmpDir, "learning-content.md")
	if err = os.WriteFile(
		contentFile, []byte("Content from file test"), 0600,
	); err != nil {
		t.Fatalf("failed to create content file: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{
		"--file", contentFile,
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
		"--context", "Testing file input",
		"--lesson", "File input works",
		"--application", "Use --file for long content",
	})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx learning add --file failed: %v", err)
	}

	content, err := os.ReadFile(".context/LEARNINGS.md")
	if err != nil {
		t.Fatalf("failed to read LEARNINGS.md: %v", err)
	}
	if !strings.Contains(string(content), "Content from file test") {
		t.Error("content from file was not added to LEARNINGS.md")
	}
}
