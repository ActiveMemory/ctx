//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupStateDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	ctxDir := filepath.Join(root, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	// Seed required files so state.Dir() considers the project
	// initialized. See specs/state-dir-no-mkdir-when-uninitialized.md.
	for _, f := range cfgCtx.FilesRequired {
		if wrErr := os.WriteFile(
			filepath.Join(ctxDir, f), []byte("# stub"), 0o600,
		); wrErr != nil {
			t.Fatalf("seed required file %s: %v", f, wrErr)
		}
	}
	t.Chdir(root)
	rc.Reset()
	stateDir := filepath.Join(ctxDir, dir.State)
	if mkErr := os.MkdirAll(stateDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	return ctxDir
}

func TestCmd_WithSessionIDFlag(t *testing.T) {
	setupStateDir(t)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--session-id", "test-xyz"})

	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	got := buf.String()
	want := "resumed for session test-xyz"
	if !strings.Contains(got, want) {
		t.Errorf("output = %q, want it to contain %q", got, want)
	}
}

func TestCmd_PauseResume_Roundtrip(t *testing.T) {
	tmpDir := setupStateDir(t)
	sessionID := "test-roundtrip"

	// Pause first - creates the marker file.
	if pauseErr := nudge.Pause(sessionID); pauseErr != nil {
		t.Fatalf("nudge.Pause() error = %v", pauseErr)
	}

	markerPath := filepath.Join(tmpDir, dir.State, "ctx-paused-"+sessionID)
	if _, statErr := os.Stat(markerPath); statErr != nil {
		t.Fatalf("pause marker should exist after Pause(): %v", statErr)
	}

	// Resume via the command.
	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--session-id", sessionID})

	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	// Verify marker is removed.
	if _, statErr := os.Stat(markerPath); !os.IsNotExist(statErr) {
		t.Error("pause marker should be removed after resume")
	}

	got := buf.String()
	want := "resumed for session test-roundtrip"
	if !strings.Contains(got, want) {
		t.Errorf("output = %q, want it to contain %q", got, want)
	}
}
