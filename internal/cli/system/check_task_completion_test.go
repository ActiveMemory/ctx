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

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestCheckTaskCompletion_SilentBeforeThreshold(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Calls 1-4 should produce no output (default interval = 5).
	for i := 1; i <= 4; i++ {
		cmd := newTestCmd()
		stdin := createTempStdin(t, `{"session_id":"test-silent"}`)
		if runErr := runCheckTaskCompletion(cmd, stdin); runErr != nil {
			t.Fatalf("call %d: unexpected error: %v", i, runErr)
		}
		out := cmdOutput(cmd)
		if strings.Contains(out, "hookSpecificOutput") {
			t.Errorf("call %d: expected silence before threshold, got: %s", i, out)
		}
	}

	// Verify counter was written.
	counterPath := filepath.Join(rc.ContextDir(), config.DirState, "task-nudge-test-silent")
	count := readCounter(counterPath)
	if count != 4 {
		t.Errorf("expected counter=4 after 4 calls, got %d", count)
	}
}

func TestCheckTaskCompletion_NudgeAtThreshold(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Pre-set counter to 4 so next call (5th) triggers nudge.
	counterPath := filepath.Join(rc.ContextDir(), config.DirState, "task-nudge-test-nudge")
	writeCounter(counterPath, 4)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-nudge"}`)
	if runErr := runCheckTaskCompletion(cmd, stdin); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "hookSpecificOutput") {
		t.Errorf("expected JSON hook response at threshold, got: %s", out)
	}
	if !strings.Contains(out, "TASKS.md") {
		t.Errorf("expected TASKS.md mention in nudge, got: %s", out)
	}

	// Counter should be reset to 0.
	count := readCounter(counterPath)
	if count != 0 {
		t.Errorf("expected counter reset to 0 after nudge, got %d", count)
	}
}

func TestCheckTaskCompletion_IntervalZeroDisabled(t *testing.T) {
	origDir, _ := os.Getwd()
	workDir := t.TempDir()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Write .ctxrc with interval=0 to disable.
	if writeErr := os.WriteFile(filepath.Join(workDir, ".ctxrc"),
		[]byte("task_nudge_interval: 0\n"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	setupContextDir(t)

	for i := 1; i <= 10; i++ {
		cmd := newTestCmd()
		stdin := createTempStdin(t, `{"session_id":"test-disabled"}`)
		if runErr := runCheckTaskCompletion(cmd, stdin); runErr != nil {
			t.Fatalf("call %d: unexpected error: %v", i, runErr)
		}
		out := cmdOutput(cmd)
		if strings.Contains(out, "hookSpecificOutput") {
			t.Errorf("call %d: expected silence when disabled, got: %s", i, out)
		}
	}
}

func TestCheckTaskCompletion_PausedSilent(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Pause the session.
	Pause("test-paused")

	// Pre-set counter to threshold-1 so it would fire if not paused.
	counterPath := filepath.Join(rc.ContextDir(), config.DirState, "task-nudge-test-paused")
	writeCounter(counterPath, 4)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-paused"}`)
	if runErr := runCheckTaskCompletion(cmd, stdin); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "hookSpecificOutput") {
		t.Errorf("expected silence when paused, got: %s", out)
	}
}

func TestCheckTaskCompletion_UninitializedSilent(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()
	// No setupContextDir — simulates pre-init state.

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-uninit"}`)
	if runErr := runCheckTaskCompletion(cmd, stdin); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "hookSpecificOutput") {
		t.Errorf("expected silence when uninitialized, got: %s", out)
	}
}

func TestCheckTaskCompletion_CustomInterval(t *testing.T) {
	origDir, _ := os.Getwd()
	workDir := t.TempDir()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Write .ctxrc with custom interval=3.
	if writeErr := os.WriteFile(filepath.Join(workDir, ".ctxrc"),
		[]byte("task_nudge_interval: 3\n"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	setupContextDir(t)

	// Calls 1-2 silent, call 3 fires.
	for i := 1; i <= 2; i++ {
		cmd := newTestCmd()
		stdin := createTempStdin(t, `{"session_id":"test-custom"}`)
		if runErr := runCheckTaskCompletion(cmd, stdin); runErr != nil {
			t.Fatalf("call %d: unexpected error: %v", i, runErr)
		}
		out := cmdOutput(cmd)
		if strings.Contains(out, "hookSpecificOutput") {
			t.Errorf("call %d: expected silence, got: %s", i, out)
		}
	}

	// 3rd call should fire.
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-custom"}`)
	if runErr := runCheckTaskCompletion(cmd, stdin); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}
	out := cmdOutput(cmd)
	if !strings.Contains(out, "hookSpecificOutput") {
		t.Errorf("expected nudge at interval 3, got: %s", out)
	}
}
