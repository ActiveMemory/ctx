//   /    ctx:                         https://ctx.ist
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
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

func setupPruneDir(t *testing.T) (string, func()) {
	t.Helper()
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	rc.Reset()

	contextDir := filepath.Join(workDir, config.DirContext)
	stateDir := filepath.Join(contextDir, config.DirState)
	if mkErr := os.MkdirAll(stateDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	// Create required files for isInitialized
	for _, f := range config.FilesRequired {
		if writeErr := os.WriteFile(filepath.Join(contextDir, f), []byte("# "+f+"\n"), 0o644); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	return stateDir, func() { _ = os.Chdir(origDir) }
}

func createStateFile(t *testing.T, dir, name string, age time.Duration) {
	t.Helper()
	path := filepath.Join(dir, name)
	if writeErr := os.WriteFile(path, nil, 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	mtime := time.Now().Add(-age)
	if chtErr := os.Chtimes(path, mtime, mtime); chtErr != nil {
		t.Fatal(chtErr)
	}
}

func TestPrune_AgeBasedDeletion(t *testing.T) {
	stateDir, cleanup := setupPruneDir(t)
	defer cleanup()

	// Old session files (10 days old)
	createStateFile(t, stateDir, "context-check-a1b2c3d4-e5f6-7890-abcd-ef1234567890", 10*24*time.Hour)
	createStateFile(t, stateDir, "heartbeat-a1b2c3d4-e5f6-7890-abcd-ef1234567890", 10*24*time.Hour)

	// Recent session file (1 day old)
	createStateFile(t, stateDir, "context-check-11111111-2222-3333-4444-555555555555", 1*24*time.Hour)

	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	runErr := runPrune(cmd, 7, false)
	if runErr != nil {
		t.Fatalf("runPrune: %v", runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "Pruned 2") {
		t.Errorf("expected 'Pruned 2', got: %s", output)
	}

	// Old files should be gone
	if _, statErr := os.Stat(filepath.Join(stateDir, "context-check-a1b2c3d4-e5f6-7890-abcd-ef1234567890")); !os.IsNotExist(statErr) {
		t.Error("expected old file to be deleted")
	}

	// Recent file should remain
	if _, statErr := os.Stat(filepath.Join(stateDir, "context-check-11111111-2222-3333-4444-555555555555")); statErr != nil {
		t.Error("expected recent file to be preserved")
	}
}

func TestPrune_DryRun(t *testing.T) {
	stateDir, cleanup := setupPruneDir(t)
	defer cleanup()

	createStateFile(t, stateDir, "heartbeat-a1b2c3d4-e5f6-7890-abcd-ef1234567890", 10*24*time.Hour)

	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	runErr := runPrune(cmd, 7, true)
	if runErr != nil {
		t.Fatalf("runPrune: %v", runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "would prune") {
		t.Errorf("expected 'would prune' in dry-run output, got: %s", output)
	}

	// File should still exist
	if _, statErr := os.Stat(filepath.Join(stateDir, "heartbeat-a1b2c3d4-e5f6-7890-abcd-ef1234567890")); statErr != nil {
		t.Error("expected file to be preserved in dry-run")
	}
}

func TestPrune_PreservesGlobalFiles(t *testing.T) {
	stateDir, cleanup := setupPruneDir(t)
	defer cleanup()

	// Global files (no UUID)
	createStateFile(t, stateDir, "events.jsonl", 30*24*time.Hour)
	createStateFile(t, stateDir, "memory-import.json", 30*24*time.Hour)
	createStateFile(t, stateDir, "check-knowledge", 30*24*time.Hour)
	createStateFile(t, stateDir, "version-checked", 30*24*time.Hour)

	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	runErr := runPrune(cmd, 1, false)
	if runErr != nil {
		t.Fatalf("runPrune: %v", runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "preserved 4 global") {
		t.Errorf("expected 'preserved 4 global', got: %s", output)
	}

	// All global files should still exist
	for _, name := range []string{"events.jsonl", "memory-import.json", "check-knowledge", "version-checked"} {
		if _, statErr := os.Stat(filepath.Join(stateDir, name)); statErr != nil {
			t.Errorf("global file %s should be preserved", name)
		}
	}
}
