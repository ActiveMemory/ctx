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
// required body flags produces an error. The placeholder
// check fires first (PreRunE runs before cobra's required-
// flag validation), so the message names the first empty
// body flag rather than --session-id.
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
	if !strings.Contains(err.Error(), "--context") {
		t.Errorf("error should mention missing --context flag: %v", err)
	}
}

// TestLearningAddRejectsPlaceholderLesson verifies placeholder
// rejection on --lesson.
func TestLearningAddRejectsPlaceholderLesson(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-learning-add-ph-*")
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
		"Learning body",
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
		"--context", "real context",
		"--lesson", "see chat",
		"--application", "real application",
	})
	err = addCmd.Execute()
	if err == nil {
		t.Fatal("expected placeholder rejection for --lesson=\"see chat\"")
	}
	if !strings.Contains(err.Error(), "lesson") {
		t.Errorf("error should name --lesson: %v", err)
	}
	if !strings.Contains(err.Error(), "placeholder") {
		t.Errorf("error should explain placeholder rejection: %v", err)
	}
}

// TestLearningAddAcceptsSubstringMatchingPlaceholder verifies
// that values containing a placeholder word as a substring are
// NOT rejected (only exact whole-value matches are placeholders).
func TestLearningAddAcceptsSubstringMatchingPlaceholder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-learning-add-substr-*")
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
		"Learning body",
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
		"--context", "we left this as TBD originally then resolved it",
		"--lesson", "real lesson",
		"--application", "real application",
	})
	if err = addCmd.Execute(); err != nil {
		t.Errorf("substring match should be accepted: %v", err)
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

// TestLearningAddFromJSONFile verifies --json-file populates the
// learning's typed fields (context/lesson/application) and provenance.
func TestLearningAddFromJSONFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-learning-add-json-*")
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

	payload := filepath.Join(tmpDir, "learning.json")
	if err = os.WriteFile(payload, []byte(`{
		"title": "Hooks run in a subprocess",
		"context": "env vars set in a hook did not persist",
		"lesson": "hook output is the only channel back to the session",
		"application": "relay via stdout, not environment",
		"provenance": {"session_id": "json1234", "branch": "main", "commit": "abc123"}
	}`), 0o600); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{"--json-file", payload})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx learning add --json-file failed: %v", err)
	}

	content, err := os.ReadFile(".context/LEARNINGS.md")
	if err != nil {
		t.Fatalf("failed to read LEARNINGS.md: %v", err)
	}
	contentStr := string(content)
	for _, want := range []string{
		"Hooks run in a subprocess",
		"env vars set in a hook did not persist",
		"hook output is the only channel back to the session",
		"relay via stdout, not environment",
	} {
		if !strings.Contains(contentStr, want) {
			t.Errorf("expected %q in LEARNINGS.md", want)
		}
	}
}
