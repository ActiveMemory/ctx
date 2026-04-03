//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/rc"
)

const (
	devContent  = "profile: dev\nnotify:\n  events:\n    - loop\n"
	baseContent = "profile: base\n# context_dir: .context\n"
)

func chdirWithCleanup(t *testing.T, dir string) {
	t.Helper()
	origDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	rc.Reset()
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})
}

// TestCopy_MissingSource verifies error on nonexistent
// source file.
func TestCopy_MissingSource(t *testing.T) {
	root := t.TempDir()

	copyErr := Copy(root, ".ctxrc.nonexistent")
	if copyErr == nil {
		t.Fatal("expected error for missing source profile")
	}
}

// TestCopy_Success verifies content is copied to .ctxrc.
func TestCopy_Success(t *testing.T) {
	root := t.TempDir()

	srcContent := "# test profile\nnotify:\n  events:\n    - loop\n"
	srcFile := ".ctxrc.test"
	if writeErr := os.WriteFile(
		filepath.Join(root, srcFile),
		[]byte(srcContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	if copyErr := Copy(root, srcFile); copyErr != nil {
		t.Fatalf("Copy failed: %v", copyErr)
	}

	data, readErr := os.ReadFile(
		filepath.Join(root, file.CtxRC),
	)
	if readErr != nil {
		t.Fatalf("expected .ctxrc to exist: %v", readErr)
	}

	if string(data) != srcContent {
		t.Errorf(
			"expected .ctxrc content to match source, got: %s",
			string(data),
		)
	}
}

// TestDetect_Dev verifies detection of the dev profile.
func TestDetect_Dev(t *testing.T) {
	root := t.TempDir()
	if writeErr := os.WriteFile(
		filepath.Join(root, file.CtxRC),
		[]byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	chdirWithCleanup(t, root)

	got := Detect()
	if got != file.ProfileDev {
		t.Errorf("expected dev, got %q", got)
	}
}

// TestDetect_Base verifies detection of the base profile.
func TestDetect_Base(t *testing.T) {
	root := t.TempDir()
	if writeErr := os.WriteFile(
		filepath.Join(root, file.CtxRC),
		[]byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	chdirWithCleanup(t, root)

	got := Detect()
	if got != file.ProfileBase {
		t.Errorf("expected base, got %q", got)
	}
}

// TestDetect_Missing verifies empty string for missing file.
func TestDetect_Missing(t *testing.T) {
	root := t.TempDir()
	chdirWithCleanup(t, root)
	got := Detect()
	if got != "" {
		t.Errorf("expected empty for missing file, got %q", got)
	}
}
