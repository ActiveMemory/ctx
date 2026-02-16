//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func cmdOutput(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}

func TestCheckContextSize_SilentEarly(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Change to temp dir so .context/logs don't pollute
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-silent"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Context Checkpoint") {
		t.Errorf("expected silence at prompt 1, got: %s", out)
	}
}

func TestCheckContextSize_CheckpointAt18(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Pre-set counter to 17 so next increment = 18 (18 > 15, 18 is not divisible by 5)
	// Need count 20 for first trigger (20 > 15, 20 % 5 == 0)
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-test-18")
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("19"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-18"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Checkpoint") {
		t.Errorf("expected checkpoint at prompt 20, got: %s", out)
	}
	if !strings.Contains(out, "prompt #20") {
		t.Errorf("expected 'prompt #20' in output, got: %s", out)
	}
}

func TestCheckContextSize_CheckpointAt33(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Pre-set counter to 32 so next = 33 (33 > 30, 33 % 3 == 0)
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-test-33")
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("32"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-33"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Checkpoint") {
		t.Errorf("expected checkpoint at prompt 33, got: %s", out)
	}
}

func TestCheckContextSize_EmptyStdin(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, "")
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not panic or error with empty input
}

// setupContextDir creates a minimal .context/ with essential files so that
// isInitialized() returns true. Must be called after chdir to the work dir.
func setupContextDir(t *testing.T) {
	t.Helper()
	dir := ".context"
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	for _, f := range config.FilesRequired {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("# "+f+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// createTempStdin writes content to a temp file and returns it opened for reading.
func createTempStdin(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "stdin-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	return f
}
