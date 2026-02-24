//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCheckPersistence_Init(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"persist-init"}`)
	if err := runCheckPersistence(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First call should just init state, no output
	out := cmdOutput(cmd)
	if strings.Contains(out, "Persistence Checkpoint") {
		t.Errorf("expected silence on init, got: %s", out)
	}

	// Verify state file was created
	stateFile := filepath.Join(tmpDir, "ctx", "persistence-nudge-persist-init")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		t.Error("state file not created")
	}
}

func TestCheckPersistence_MtimeReset(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// Create state file with old mtime to simulate context modification
	stateFile := filepath.Join(tmpDir, "ctx", "persistence-nudge-persist-mtime")
	_ = os.MkdirAll(filepath.Dir(stateFile), 0o700)
	writePersistenceState(stateFile, persistenceState{
		Count:     20,
		LastNudge: 0,
		LastMtime: time.Now().Unix() - 3600, // 1 hour ago
	})

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"persist-mtime"}`)
	if err := runCheckPersistence(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be silent because mtime changed (context was modified)
	out := cmdOutput(cmd)
	if strings.Contains(out, "Persistence Checkpoint") {
		t.Errorf("expected silence after mtime reset, got: %s", out)
	}
}

func TestCheckPersistence_NudgeAt20(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	futureMtime := time.Now().Unix() + 3600

	stateFile := filepath.Join(tmpDir, "ctx", "persistence-nudge-persist-20")
	_ = os.MkdirAll(filepath.Dir(stateFile), 0o700)
	writePersistenceState(stateFile, persistenceState{
		Count:     19, // will become 20
		LastNudge: 0,
		LastMtime: futureMtime,
	})

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"persist-20"}`)
	if err := runCheckPersistence(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Persistence Checkpoint") {
		t.Errorf("expected nudge at prompt 20, got: %s", out)
	}
	if !strings.Contains(out, "Context:") {
		t.Errorf("expected context dir footer, got: %s", out)
	}
}

func TestCheckPersistence_Every15AfterPrompt25(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	futureMtime := time.Now().Unix() + 3600

	stateFile := filepath.Join(tmpDir, "ctx", "persistence-nudge-persist-40")
	_ = os.MkdirAll(filepath.Dir(stateFile), 0o700)
	writePersistenceState(stateFile, persistenceState{
		Count:     39, // will become 40
		LastNudge: 25, // 40 - 25 = 15 >= 15
		LastMtime: futureMtime,
	})

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"persist-40"}`)
	if err := runCheckPersistence(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Persistence Checkpoint") {
		t.Errorf("expected nudge at prompt 40, got: %s", out)
	}
}
