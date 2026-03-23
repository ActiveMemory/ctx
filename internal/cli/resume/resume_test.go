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
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupStateDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("CTX_DIR", tmpDir)
	rc.Reset()
	if mkErr := os.MkdirAll(filepath.Join(tmpDir, dir.State), 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	return tmpDir
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

	// Pause first — creates the marker file.
	nudge.Pause(sessionID)

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
