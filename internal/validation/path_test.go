//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validation

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestValidateBoundary(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{"relative inside cwd", ".context", false},
		{"absolute inside cwd", filepath.Join(cwd, ".context"), false},
		{"deeply nested", filepath.Join(cwd, "a", "b", "c"), false},
		{"cwd itself", cwd, false},
		{"dot", ".", false},
		{"escapes cwd", "../../etc", true},
		{"absolute outside cwd", "/tmp/evil", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBoundary(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBoundary(%q) error = %v, wantErr %v", tt.dir, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBoundaryCaseInsensitive(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("case-insensitive path test only applies to Windows")
	}

	// On Windows, EvalSymlinks normalizes casing to the filesystem's
	// canonical form. When .context/ doesn't exist yet the fallback
	// preserves the original cwd casing. The prefix check must be
	// case-insensitive to avoid false "outside cwd" errors.
	tmp := t.TempDir()

	// Change cwd to a case-mangled version of the temp dir.
	// TempDir returns canonical casing; flip it.
	mangled := strings.ToUpper(tmp)
	if mangled == tmp {
		mangled = strings.ToLower(tmp)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(orig) }()

	if err := os.Chdir(mangled); err != nil {
		t.Skipf("cannot chdir to case-mangled path %q: %v", mangled, err)
	}

	// .context doesn't exist — this is the exact scenario that caused the
	// false positive on Windows.
	if err := ValidateBoundary(".context"); err != nil {
		t.Errorf("ValidateBoundary(.context) with case-mangled cwd: %v", err)
	}
}

func TestCheckSymlinks(t *testing.T) {
	t.Run("regular directory passes", func(t *testing.T) {
		dir := t.TempDir()
		// Create a regular file inside.
		if err := os.WriteFile(filepath.Join(dir, "file.md"), []byte("ok"), 0600); err != nil {
			t.Fatal(err)
		}

		if err := CheckSymlinks(dir); err != nil {
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

		err := CheckSymlinks(linkDir)
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

		err := CheckSymlinks(dir)
		if err == nil {
			t.Error("CheckSymlinks with symlinked child: expected error, got nil")
		}
	})

	t.Run("non-existent directory passes", func(t *testing.T) {
		if err := CheckSymlinks("/nonexistent/path"); err != nil {
			t.Errorf("CheckSymlinks on non-existent dir: unexpected error: %v", err)
		}
	})
}
