//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
)

// setup creates a temp dir with a .context/ directory, sets the RC context
// dir override, and returns the temp dir path.
func setup(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()
	rc.OverrideContextDir(dir.Context)

	ctxDir := filepath.Join(tmpDir, dir.Context)
	if err := os.MkdirAll(ctxDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	return tmpDir
}

// newPromptCmd builds a fresh command with the given args.
func newPromptCmd(args ...string) *cobra.Command {
	cmd := Cmd()
	cmd.SetArgs(args)
	return cmd
}

// runCmd captures cobra output.
func runCmd(cmd *cobra.Command) (string, error) {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

func TestList_Empty(t *testing.T) {
	setup(t)

	out, err := runCmd(newPromptCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No prompts found") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestList_NoDir(t *testing.T) {
	setup(t)

	// Don't create prompts dir - should handle gracefully
	out, err := runCmd(newPromptCmd("list"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No prompts found") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestList_WithPrompts(t *testing.T) {
	tmpDir := setup(t)

	// Create prompts directory with files
	promptDir := filepath.Join(tmpDir, dir.Context, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "review.md"), []byte("# Review"), fs.PermFile); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "debug.md"), []byte("# Debug"), fs.PermFile); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPromptCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "review") {
		t.Errorf("expected 'review' in output, got: %s", out)
	}
	if !strings.Contains(out, "debug") {
		t.Errorf("expected 'debug' in output, got: %s", out)
	}
}

func TestShow(t *testing.T) {
	tmpDir := setup(t)

	promptDir := filepath.Join(tmpDir, dir.Context, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}
	content := "# Code Review\n\nReview this code.\n"
	if err := os.WriteFile(filepath.Join(promptDir, "review.md"), []byte(content), fs.PermFile); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPromptCmd("show", "review"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != content {
		t.Errorf("expected %q, got %q", content, out)
	}
}

func TestShow_Missing(t *testing.T) {
	setup(t)

	_, err := runCmd(newPromptCmd("show", "nonexistent"))
	if err == nil {
		t.Fatal("expected error for missing prompt")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestAdd_FromTemplate(t *testing.T) {
	tmpDir := setup(t)

	// Create prompts dir
	promptDir := filepath.Join(tmpDir, dir.Context, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPromptCmd("add", "code-review"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Created") {
		t.Errorf("expected 'Created' message, got: %s", out)
	}

	// Verify file exists with content
	data, readErr := os.ReadFile(filepath.Join(promptDir, "code-review.md"))
	if readErr != nil {
		t.Fatalf("failed to read created file: %v", readErr)
	}
	if !strings.Contains(string(data), "Review") {
		t.Error("created file does not contain expected content")
	}
}

func TestAdd_FromStdin(t *testing.T) {
	tmpDir := setup(t)

	promptDir := filepath.Join(tmpDir, dir.Context, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}

	cmd := newPromptCmd("add", "custom", "--stdin")
	cmd.SetIn(strings.NewReader("# Custom Prompt\n\nDo the thing.\n"))

	out, err := runCmd(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Created") {
		t.Errorf("expected 'Created' message, got: %s", out)
	}

	data, readErr := os.ReadFile(filepath.Join(promptDir, "custom.md"))
	if readErr != nil {
		t.Fatalf("failed to read created file: %v", readErr)
	}
	if !strings.Contains(string(data), "Custom Prompt") {
		t.Error("created file does not contain stdin content")
	}
}

func TestAdd_AlreadyExists(t *testing.T) {
	tmpDir := setup(t)

	promptDir := filepath.Join(tmpDir, dir.Context, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "review.md"), []byte("existing"), fs.PermFile); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPromptCmd("add", "review"))
	if err == nil {
		t.Fatal("expected error for existing prompt")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestAdd_NoTemplate(t *testing.T) {
	setup(t)

	_, err := runCmd(newPromptCmd("add", "nonexistent-template"))
	if err == nil {
		t.Fatal("expected error for missing embedded template")
	}
	if !strings.Contains(err.Error(), "no embedded template") {
		t.Errorf("expected 'no embedded template' error, got: %v", err)
	}
}

func TestRm(t *testing.T) {
	tmpDir := setup(t)

	promptDir := filepath.Join(tmpDir, dir.Context, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "review.md"), []byte("# Review"), fs.PermFile); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPromptCmd("rm", "review"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Removed") {
		t.Errorf("expected 'Removed' message, got: %s", out)
	}

	// Verify file is gone
	if _, statErr := os.Stat(filepath.Join(promptDir, "review.md")); !os.IsNotExist(statErr) {
		t.Error("file should have been removed")
	}
}

func TestRm_Missing(t *testing.T) {
	setup(t)

	_, err := runCmd(newPromptCmd("rm", "nonexistent"))
	if err == nil {
		t.Fatal("expected error for missing prompt")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}
