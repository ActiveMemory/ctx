//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// declareCtxDir creates a `.context/` under a fresh tempDir, chdirs
// into the tempDir so that `$PWD/.context/` resolves there, and
// optionally seeds the required files so [state.Initialized] will
// return true. Returns the absolute path of the .context dir.
func declareCtxDir(t *testing.T, initialized bool) string {
	t.Helper()
	tempDir := t.TempDir()
	ctxDir := filepath.Join(tempDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0o700); mkErr != nil {
		t.Fatalf("mkdir .context: %v", mkErr)
	}
	if initialized {
		for _, f := range cfgCtx.FilesRequired {
			path := filepath.Join(ctxDir, f)
			if wrErr := os.WriteFile(path, []byte("# stub"), 0o600); wrErr != nil {
				t.Fatalf("seed required file %s: %v", f, wrErr)
			}
		}
	}
	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
	return ctxDir
}

// TestDir_Initialized verifies the happy path: in an initialized
// project, Dir returns the state path and creates the directory.
func TestDir_Initialized(t *testing.T) {
	ctxDir := declareCtxDir(t, true)

	got, err := state.Dir()
	if err != nil {
		t.Fatalf("Dir() error = %v, want nil", err)
	}
	want := filepath.Join(ctxDir, dir.State)
	if got != want {
		t.Errorf("Dir() = %q, want %q", got, want)
	}
	if _, statErr := os.Stat(got); statErr != nil {
		t.Errorf("state dir was not created: %v", statErr)
	}
}

// TestDir_Uninitialized is the regression test for the cross-IDE
// hook leak: when CTX_DIR is declared but the project is not
// initialized, Dir must return ErrNotInitialized and must NOT
// create the state directory. This is the structural invariant
// established in specs/state-dir-no-mkdir-when-uninitialized.md.
func TestDir_Uninitialized(t *testing.T) {
	ctxDir := declareCtxDir(t, false)
	stateDir := filepath.Join(ctxDir, dir.State)

	got, err := state.Dir()
	if err == nil {
		t.Fatal("Dir() error = nil, want ErrNotInitialized")
	}
	if !errors.Is(err, errCtx.ErrNotInitialized) {
		t.Errorf("Dir() error = %v, want errors.Is(.., ErrNotInitialized)", err)
	}
	if got != "" {
		t.Errorf("Dir() path = %q, want empty string on error", got)
	}
	if _, statErr := os.Stat(stateDir); !os.IsNotExist(statErr) {
		t.Errorf("state/ was created in uninitialized project: stat err = %v (want IsNotExist)", statErr)
	}
}

// TestDir_UninitializedDoesNotCreateContextDir is the strongest
// form of the invariant: if .context/ itself does not exist on
// disk, Dir must surface ErrNoCtxHere (cwd-anchored variant)
// and must NOT materialize .context/ or .context/state/. This is
// the observed Cursor leak shape: opening a non-ctx workspace and
// submitting one prompt must leave the filesystem unchanged.
func TestDir_UninitializedDoesNotCreateContextDir(t *testing.T) {
	tempDir := t.TempDir()
	ctxDir := filepath.Join(tempDir, dir.Context) // does not exist on disk
	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	_, err := state.Dir()
	if !errors.Is(err, errCtx.ErrNoCtxHere) {
		t.Fatalf("Dir() error = %v, want errors.Is(.., ErrNoCtxHere)", err)
	}
	if _, statErr := os.Stat(ctxDir); !os.IsNotExist(statErr) {
		t.Errorf(".context/ was materialized: stat err = %v (want IsNotExist)", statErr)
	}
	stateDir := filepath.Join(ctxDir, dir.State)
	if _, statErr := os.Stat(stateDir); !os.IsNotExist(statErr) {
		t.Errorf(".context/state/ was materialized: stat err = %v (want IsNotExist)", statErr)
	}
}

// TestDir_NoContextHere: cwd lacks `.context/` → ErrNoCtxHere
// propagates out of Dir.
func TestDir_NoContextHere(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	_, err := state.Dir()
	if !errors.Is(err, errCtx.ErrNoCtxHere) {
		t.Errorf("Dir() error = %v, want errors.Is(.., ErrNoCtxHere)", err)
	}
}

// TestDir_Override verifies the test-only override bypasses both
// the resolver and the Initialized gate. Tests that explicitly
// opt into a state dir must continue to work without faking the
// initialized state.
func TestDir_Override(t *testing.T) {
	override := t.TempDir()
	state.SetDirForTest(override)
	t.Cleanup(func() { state.SetDirForTest("") })

	got, err := state.Dir()
	if err != nil {
		t.Fatalf("Dir() with override error = %v, want nil", err)
	}
	if got != override {
		t.Errorf("Dir() with override = %q, want %q", got, override)
	}
}

// TestInitialized_Uninitialized confirms the helper agrees with
// Dir's gate: an uninitialized project reports false.
func TestInitialized_Uninitialized(t *testing.T) {
	declareCtxDir(t, false)
	got, err := state.Initialized()
	if err != nil {
		t.Fatalf("Initialized() error = %v, want nil", err)
	}
	if got {
		t.Error("Initialized() = true, want false")
	}
}

// TestInitialized_Initialized confirms the helper agrees with
// Dir's gate: an initialized project reports true.
func TestInitialized_Initialized(t *testing.T) {
	declareCtxDir(t, true)
	got, err := state.Initialized()
	if err != nil {
		t.Fatalf("Initialized() error = %v, want nil", err)
	}
	if !got {
		t.Error("Initialized() = false, want true")
	}
}
