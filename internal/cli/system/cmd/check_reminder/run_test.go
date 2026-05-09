//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_reminder_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_reminder"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// TestRun_NoLeakInUninitializedProject reproduces the cross-IDE
// hook leak (specs/state-dir-no-mkdir-when-uninitialized.md) at
// the entry point that originally triggered it: check-reminder
// was called from Cursor's imported Claude hooks with CTX_DIR
// pointing at a non-ctx workspace, and the call chain
// (Preamble → nudge.Paused → PauseMarkerPath → state.Dir) leaked
// `.context/state/` (mode 0750) into the workspace.
//
// Acceptance: after running check-reminder against a directory
// that has CTX_DIR set but no ctx init, neither `.context/` nor
// `.context/state/` exists on disk. Hook exit must be non-error
// (hooks never fail the parent operation).
func TestRun_NoLeakInUninitializedProject(t *testing.T) {
	tempDir := t.TempDir()
	ctxDir := filepath.Join(tempDir, dir.Context) // not on disk
	stateDir := filepath.Join(ctxDir, dir.State)

	t.Setenv(env.CtxDir, ctxDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	// Feed a minimal valid hook envelope on stdin via a pipe so
	// Run reads from a real *os.File (its signature is exact).
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("os.Pipe: %v", pipeErr)
	}
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = io.Copy(w, bytes.NewReader([]byte(
			`{"session_id":"00000000-0000-0000-0000-000000000000"}`,
		)))
	}()
	t.Cleanup(func() { _ = r.Close() })

	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	if err := check_reminder.Run(cmd, r); err != nil {
		t.Fatalf("Run() error = %v, want nil (hooks must never fail)", err)
	}

	if _, statErr := os.Stat(ctxDir); !os.IsNotExist(statErr) {
		t.Errorf(".context/ leaked into uninitialized project: stat err = %v (want IsNotExist)", statErr)
	}
	if _, statErr := os.Stat(stateDir); !os.IsNotExist(statErr) {
		t.Errorf(".context/state/ leaked into uninitialized project: stat err = %v (want IsNotExist)", statErr)
	}
}
