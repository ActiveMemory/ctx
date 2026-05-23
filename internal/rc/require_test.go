//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
)

// TestRequireContextDir_NoContextHere: cwd has no .context/ →
// ErrNoCtxHere.
func TestRequireContextDir_NoContextHere(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	got, err := RequireContextDir()
	if !errors.Is(err, errCtx.ErrNoCtxHere) {
		t.Errorf("RequireContextDir() err = %v, want ErrNoCtxHere", err)
	}
	if got != "" {
		t.Errorf("RequireContextDir() = %q, want \"\"", got)
	}
}

// TestRequireContextDir_PathIsAFile: $PWD/.context exists as a
// regular file → ErrContextDirNotADirectory.
func TestRequireContextDir_PathIsAFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, dir.Context)
	if err := os.WriteFile(filePath, []byte("not a dir"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	_, err := RequireContextDir()
	if !errors.Is(err, errCtx.ErrContextDirNotADirectory) {
		t.Errorf("RequireContextDir() err = %v, want ErrContextDirNotADirectory",
			err)
	}
}

// TestRequireContextDir_StatPermissionDenied: stat fails for a
// reason other than not-exist → ErrContextDirStat. Skipped on
// platforms where chmod 000 doesn't block stat (Windows) or where
// the test runs as root.
func TestRequireContextDir_StatPermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission semantics differ on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks")
	}
	tempDir := t.TempDir()
	parent := filepath.Join(tempDir, "locked")
	if err := os.MkdirAll(parent, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	target := filepath.Join(parent, dir.Context)
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	// t.Chdir must happen *before* the chmod-0 lockdown: a chmod-0
	// directory cannot be made the cwd. Once cwd is set, the chmod
	// prevents the resolver from stat'ing children.
	t.Chdir(parent)
	if err := os.Chmod(parent, 0); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() {
		// Restore rwx so t.TempDir's recursive cleanup can
		// remove the directory. gosec G302 flags 0o700 as too
		// permissive for files; it is fine for an in-test
		// directory chmod that needs read+write+execute for
		// cleanup to succeed.
		_ = os.Chmod(parent, 0o700) //nolint:gosec // dir needs rwx for cleanup
	})

	Reset()
	t.Cleanup(Reset)

	_, err := RequireContextDir()
	if err == nil {
		t.Fatal("RequireContextDir() err = nil, want non-nil")
	}
	// Either ErrNoCtxHere or ErrContextDirStat depending on
	// the underlying syscall: macOS often returns ENOENT through a
	// chmod-0 parent because lookup short-circuits, while Linux
	// typically surfaces EACCES. Both are acceptable diagnostics for
	// the user.
	if !errors.Is(err, errCtx.ErrContextDirStat) &&
		!errors.Is(err, errCtx.ErrNoCtxHere) {
		t.Errorf(
			"RequireContextDir() err = %v, want ErrContextDirStat or ErrNoCtxHere",
			err)
	}
}

// TestRequireContextDir_HappyPath: cwd has .context/ → returns
// absolute path, nil error.
func TestRequireContextDir_HappyPath(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, dir.Context)
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	got, err := RequireContextDir()
	if err != nil {
		t.Fatalf("RequireContextDir() err = %v, want nil", err)
	}
	gotResolved, _ := filepath.EvalSymlinks(got)
	wantResolved, _ := filepath.EvalSymlinks(target)
	if gotResolved != wantResolved {
		t.Errorf("RequireContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}
