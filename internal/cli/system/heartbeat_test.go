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

	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestHeartbeat_Silent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"hb-silent"}`)
	if err := runHeartbeat(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("heartbeat must produce no stdout, got: %q", out)
	}
}

func TestHeartbeat_CounterIncrement(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// First call: counter should be 1.
	cmd1 := newTestCmd()
	stdin1 := createTempStdin(t, `{"session_id":"hb-counter"}`)
	if err := runHeartbeat(cmd1, stdin1); err != nil {
		t.Fatalf("call 1: unexpected error: %v", err)
	}

	counterFile := filepath.Join(tmpDir, "ctx", "heartbeat-hb-counter")
	count1 := readCounter(counterFile)
	if count1 != 1 {
		t.Errorf("after call 1: expected counter=1, got %d", count1)
	}

	// Second call: counter should be 2.
	cmd2 := newTestCmd()
	stdin2 := createTempStdin(t, `{"session_id":"hb-counter"}`)
	if err := runHeartbeat(cmd2, stdin2); err != nil {
		t.Fatalf("call 2: unexpected error: %v", err)
	}

	count2 := readCounter(counterFile)
	if count2 != 2 {
		t.Errorf("after call 2: expected counter=2, got %d", count2)
	}

	// Third call: counter should be 3.
	cmd3 := newTestCmd()
	stdin3 := createTempStdin(t, `{"session_id":"hb-counter"}`)
	if err := runHeartbeat(cmd3, stdin3); err != nil {
		t.Fatalf("call 3: unexpected error: %v", err)
	}

	count3 := readCounter(counterFile)
	if count3 != 3 {
		t.Errorf("after call 3: expected counter=3, got %d", count3)
	}
}

func TestHeartbeat_ContextModifiedDetection(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// First call establishes baseline mtime.
	cmd1 := newTestCmd()
	stdin1 := createTempStdin(t, `{"session_id":"hb-mtime"}`)
	if err := runHeartbeat(cmd1, stdin1); err != nil {
		t.Fatalf("call 1: unexpected error: %v", err)
	}

	mtimeFile := filepath.Join(tmpDir, "ctx", "heartbeat-mtime-hb-mtime")
	storedMtime := readMtime(mtimeFile)
	if storedMtime == 0 {
		t.Fatal("mtime file should have a non-zero value after first call")
	}

	// Touch a context file to advance its mtime.
	contextDir := rc.ContextDir()
	touchPath := filepath.Join(contextDir, "TASKS.md")
	future := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(touchPath, future, future); err != nil {
		t.Fatalf("failed to touch context file: %v", err)
	}

	// Second call should detect the modification.
	cmd2 := newTestCmd()
	stdin2 := createTempStdin(t, `{"session_id":"hb-mtime"}`)
	if err := runHeartbeat(cmd2, stdin2); err != nil {
		t.Fatalf("call 2: unexpected error: %v", err)
	}

	newStoredMtime := readMtime(mtimeFile)
	if newStoredMtime <= storedMtime {
		t.Errorf("expected updated mtime after context modification: old=%d new=%d", storedMtime, newStoredMtime)
	}
}

func TestHeartbeat_RespectsNotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Work in a directory with no .context/ — not initialized.
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"hb-noinit"}`)
	if err := runHeartbeat(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected no output when not initialized, got: %q", out)
	}

	// Counter file should not exist.
	counterFile := filepath.Join(tmpDir, "ctx", "heartbeat-hb-noinit")
	if _, statErr := os.Stat(counterFile); !os.IsNotExist(statErr) {
		t.Error("counter file should not be created when not initialized")
	}
}

func TestHeartbeat_RespectsPaused(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// Create pause marker.
	Pause("hb-paused")

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"hb-paused"}`)
	if err := runHeartbeat(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected no output when paused, got: %q", out)
	}

	// Counter file should not exist (heartbeat skipped entirely).
	counterFile := filepath.Join(tmpDir, "ctx", "heartbeat-hb-paused")
	if _, statErr := os.Stat(counterFile); !os.IsNotExist(statErr) {
		t.Error("counter file should not be created when paused")
	}
}

func TestHeartbeat_TokenTelemetry(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// Create a fake JSONL file with usage data.
	sessionID := "hb-token-test"
	projectDir := filepath.Join(tmpDir, ".claude", "projects", "testproj")
	if mkErr := os.MkdirAll(projectDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	jsonlContent := `{"type":"assistant","message":{"model":"claude-sonnet-4-5","role":"assistant","content":"hi","usage":{"input_tokens":50000,"output_tokens":500,"cache_creation_input_tokens":8000,"cache_read_input_tokens":100000}}}` + "\n"
	jsonlPath := filepath.Join(projectDir, sessionID+".jsonl")
	if writeErr := os.WriteFile(jsonlPath, []byte(jsonlContent), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"`+sessionID+`"}`)
	if err := runHeartbeat(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("heartbeat must produce no stdout, got: %q", out)
	}

	// Verify the heartbeat log includes token data.
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, "logs", "heartbeat.log")
	logData, readErr := os.ReadFile(logFile) //nolint:gosec // test path
	if readErr != nil {
		t.Fatalf("failed to read heartbeat log: %v", readErr)
	}
	logStr := string(logData)
	if !strings.Contains(logStr, "tokens=") {
		t.Errorf("heartbeat log missing token data: %s", logStr)
	}
	if !strings.Contains(logStr, "pct=") {
		t.Errorf("heartbeat log missing pct data: %s", logStr)
	}
}

func TestHeartbeat_EmptyStdin(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, "")
	if err := runHeartbeat(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected no output with empty stdin, got: %q", out)
	}

	// Should have used "unknown" as session ID.
	counterFile := filepath.Join(tmpDir, "ctx", "heartbeat-unknown")
	count := readCounter(counterFile)
	if count != 1 {
		t.Errorf("expected counter=1 for fallback session, got %d", count)
	}
}
