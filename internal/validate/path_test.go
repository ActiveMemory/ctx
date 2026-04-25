//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckSymlinks(t *testing.T) {
	t.Run("regular directory passes", func(t *testing.T) {
		dir := t.TempDir()
		// Create a regular file inside.
		fPath := filepath.Join(dir, "file.md")
		if err := os.WriteFile(fPath, []byte("ok"), 0600); err != nil {
			t.Fatal(err)
		}

		if err := Symlinks(dir); err != nil {
			t.Errorf("CheckSymlinks on regular dir: unexpected error: %v", err)
		}
	})

	t.Run("directory that is a symlink fails", func(t *testing.T) {
		tmp := t.TempDir()
		realDir := filepath.Join(tmp, "real")
		if err := os.Mkdir(realDir, 0750); err != nil {
			t.Fatal(err)
		}
		linkDir := filepath.Join(tmp, "link")
		if err := os.Symlink(realDir, linkDir); err != nil {
			t.Fatal(err)
		}

		err := Symlinks(linkDir)
		if err == nil {
			t.Error("CheckSymlinks on symlinked dir: expected error, got nil")
		}
	})

	t.Run("directory containing symlinked file fails", func(t *testing.T) {
		dir := t.TempDir()
		// Create a real file elsewhere and symlink it into the dir.
		realFile := filepath.Join(t.TempDir(), "real.md")
		if err := os.WriteFile(realFile, []byte("secret"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(realFile, filepath.Join(dir, "TASKS.md")); err != nil {
			t.Fatal(err)
		}

		err := Symlinks(dir)
		if err == nil {
			t.Error("CheckSymlinks with symlinked child: expected error, got nil")
		}
	})

	t.Run("non-existent directory passes", func(t *testing.T) {
		if err := Symlinks("/nonexistent/path"); err != nil {
			t.Errorf("CheckSymlinks on non-existent dir: unexpected error: %v", err)
		}
	})
}
