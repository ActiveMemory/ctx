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
)

func TestContextLoadGate_NotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	// No .context/ — should be silent
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-no-init"}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence when not initialized, got: %s", out)
	}
}

func TestContextLoadGate_EmptySessionID(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence with empty session_id, got: %s", out)
	}
}

func TestContextLoadGate_FirstToolUse(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-first-tool"}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "STOP. You must read these files in order before proceeding") {
		t.Errorf("expected context-load directive, got: %s", out)
	}
	if !strings.Contains(out, "additionalContext") {
		t.Errorf("expected JSON HookResponse, got: %s", out)
	}
	if !strings.Contains(out, "Do not assess relevance") {
		t.Errorf("expected no-relevance-filtering language, got: %s", out)
	}
	if !strings.Contains(out, "Context Loaded") {
		t.Errorf("expected unconditional checkpoint block, got: %s", out)
	}
	if !strings.Contains(out, "MANDATORY") {
		t.Errorf("expected mandatory language, got: %s", out)
	}
	// Verify all files from FileReadOrder are listed, except GLOSSARY
	for _, f := range config.FileReadOrder {
		if f == config.FileGlossary {
			if strings.Contains(out, f) {
				t.Errorf("GLOSSARY should be excluded from gate, but found in: %s", out)
			}
			continue
		}
		if !strings.Contains(out, f) {
			t.Errorf("expected file %q in directive, got: %s", f, out)
		}
	}
}

func TestContextLoadGate_SecondToolUse_Silent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// First tool use — emits directive and creates marker
	cmd1 := newTestCmd()
	stdin1 := createTempStdin(t, `{"session_id":"test-second-tool"}`)
	if err := runContextLoadGate(cmd1, stdin1); err != nil {
		t.Fatalf("first tool use: unexpected error: %v", err)
	}
	if cmdOutput(cmd1) == "" {
		t.Fatal("first tool use: expected directive output")
	}

	// Second tool use — marker exists, should be silent
	cmd2 := newTestCmd()
	stdin2 := createTempStdin(t, `{"session_id":"test-second-tool"}`)
	if err := runContextLoadGate(cmd2, stdin2); err != nil {
		t.Fatalf("second tool use: unexpected error: %v", err)
	}

	out := cmdOutput(cmd2)
	if out != "" {
		t.Errorf("expected silence on second tool use, got: %s", out)
	}
}

func TestContextLoadGate_DifferentSessions(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// Session A — emits directive
	cmdA := newTestCmd()
	stdinA := createTempStdin(t, `{"session_id":"session-a"}`)
	if err := runContextLoadGate(cmdA, stdinA); err != nil {
		t.Fatalf("session-a: unexpected error: %v", err)
	}
	if cmdOutput(cmdA) == "" {
		t.Fatal("session-a: expected directive output")
	}

	// Verify marker exists for session-a
	marker := filepath.Join(tmpDir, "ctx", "ctx-loaded-session-a")
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("expected marker for session-a, got error: %v", err)
	}

	// Session B — different session_id, should also emit directive
	cmdB := newTestCmd()
	stdinB := createTempStdin(t, `{"session_id":"session-b"}`)
	if err := runContextLoadGate(cmdB, stdinB); err != nil {
		t.Fatalf("session-b: unexpected error: %v", err)
	}
	outB := cmdOutput(cmdB)
	if !strings.Contains(outB, "STOP. You must read these files in order before proceeding") {
		t.Errorf("session-b: expected context-load directive, got: %s", outB)
	}

	// Session A again — should be silent
	cmdA2 := newTestCmd()
	stdinA2 := createTempStdin(t, `{"session_id":"session-a"}`)
	if err := runContextLoadGate(cmdA2, stdinA2); err != nil {
		t.Fatalf("session-a repeat: unexpected error: %v", err)
	}
	if cmdOutput(cmdA2) != "" {
		t.Errorf("session-a repeat: expected silence, got: %s", cmdOutput(cmdA2))
	}
}
