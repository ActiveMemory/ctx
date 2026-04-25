//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupStateDir(t *testing.T) string {
	t.Helper()
	ctxDir := filepath.Join(t.TempDir(), dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	t.Setenv("CTX_DIR", ctxDir)
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
	cmd.SetArgs([]string{"--session-id", "test-abc"})

	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	got := buf.String()
	want := "paused for session test-abc"
	if !strings.Contains(got, want) {
		t.Errorf("output = %q, want it to contain %q", got, want)
	}
}

func TestCmd_DefaultsToEmptySessionID(t *testing.T) {
	setupStateDir(t)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	got := buf.String()
	want := "paused for session"
	if !strings.Contains(got, want) {
		t.Errorf("output = %q, want it to contain %q", got, want)
	}
}
