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
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
)

// TestRequireContextDir_PathDoesNotExist: shape-valid declaration
// pointing at a path that doesn't exist on disk → ErrContextDirNotFound.
func TestRequireContextDir_PathDoesNotExist(t *testing.T) {
	t.Setenv(env.CtxDir, "/nonexistent-test-dir/.context")
	Reset()
	t.Cleanup(Reset)

	got, err := RequireContextDir()
	if !errors.Is(err, errCtx.ErrContextDirNotFound) {
		t.Errorf("RequireContextDir() err = %v, want ErrContextDirNotFound",
			err)
	}
	if got != "" {
		t.Errorf("RequireContextDir() = %q, want \"\"", got)
	}
}

// TestRequireContextDir_PathIsAFile: CTX_DIR points at an existing
// regular file → ErrContextDirNotADirectory.
func TestRequireContextDir_PathIsAFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, dir.Context)
	if err := os.WriteFile(filePath, []byte("not a dir"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Setenv(env.CtxDir, filePath)
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

	t.Setenv(env.CtxDir, target)
	Reset()
	t.Cleanup(Reset)

	_, err := RequireContextDir()
	if err == nil {
		t.Fatal("RequireContextDir() err = nil, want non-nil")
	}
	// Either ErrContextDirNotFound or ErrContextDirStat depending on
	// the underlying syscall: macOS often returns ENOENT through a
	// chmod-0 parent because lookup short-circuits, while Linux
	// typically surfaces EACCES. Both are acceptable diagnostics for
	// the user.
	if !errors.Is(err, errCtx.ErrContextDirStat) &&
		!errors.Is(err, errCtx.ErrContextDirNotFound) {
		t.Errorf(
			"RequireContextDir() err = %v, want ErrContextDirStat or ErrContextDirNotFound",
			err)
	}
}

// TestRequireContextDir_HappyPath: existing dir, canonical name →
// returns absolute path, nil error.
func TestRequireContextDir_HappyPath(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, dir.Context)
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Setenv(env.CtxDir, target)
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

// TestRequireContextDir_DelegatesShapeChecks: ContextDir shape
// errors flow through with their precise meaning preserved. Only
// the truly-unset case gets rewrapped as the tailored
// "no context directory specified" message with candidate hints;
// relative and non-canonical-basename errors propagate unchanged so
// the user sees what's wrong with the value they declared instead
// of "you didn't declare it" when they actually did.
func TestRequireContextDir_DelegatesShapeChecks(t *testing.T) {
	cases := []struct {
		name           string
		val            string
		wantSentinel   error
		wantMsgContain string
	}{
		{
			name:           "unset",
			val:            "",
			wantSentinel:   errCtx.ErrDirNotDeclared,
			wantMsgContain: "no context directory specified",
		},
		{
			name:           "relative",
			val:            "relative-path",
			wantSentinel:   errCtx.ErrRelativeNotAllowed,
			wantMsgContain: "absolute",
		},
		{
			name:           "non-canonical",
			val:            "/tmp/notdotcontext",
			wantSentinel:   errCtx.ErrNonCanonicalBasename,
			wantMsgContain: "notdotcontext",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv(env.CtxDir, c.val)
			Reset()
			t.Cleanup(Reset)

			got, err := RequireContextDir()
			if err == nil {
				t.Fatalf("RequireContextDir() err = nil, want non-nil for %q",
					c.val)
			}
			if got != "" {
				t.Errorf("RequireContextDir() = %q, want \"\"", got)
			}
			// "unset" gets rewrapped into a tailored message that no
			// longer wraps the original sentinel. The other two
			// shape errors propagate the sentinel unchanged.
			if c.name != "unset" && !errors.Is(err, c.wantSentinel) {
				t.Errorf("RequireContextDir() err = %v, want errors.Is matching %v",
					err, c.wantSentinel)
			}
			if msg := err.Error(); msg == "" {
				t.Error("RequireContextDir() returned empty error message")
			}
			if !strings.Contains(err.Error(), c.wantMsgContain) {
				t.Errorf("RequireContextDir() msg = %q; want substring %q",
					err.Error(), c.wantMsgContain)
			}
		})
	}
}
