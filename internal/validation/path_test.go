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
	if runtime.GOOS != osWindows {
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

func TestValidateBoundary_WindowsCaseInsensitive(t *testing.T) {
	if runtime.GOOS != osWindows {
		t.Skip("Windows-only test")
	}

	// Simulate the VS Code plugin scenario: CWD has a lowercase drive letter
	// but EvalSymlinks resolves to the actual (uppercase) casing.
	// When .context doesn't exist yet (first init), the fallback path
	// preserves the lowercase letter, causing a case mismatch.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Swap the drive letter case to simulate VS Code's fsPath
	if len(cwd) >= 2 && cwd[1] == ':' {
		var swapped string
		if cwd[0] >= 'A' && cwd[0] <= 'Z' {
			swapped = strings.ToLower(cwd[:1]) + cwd[1:]
		} else {
			swapped = strings.ToUpper(cwd[:1]) + cwd[1:]
		}

		origDir, _ := os.Getwd()
		if chErr := os.Chdir(swapped); chErr != nil {
			t.Fatalf("cannot chdir to %s: %v", swapped, chErr)
		}
		defer func() { _ = os.Chdir(origDir) }()

		// Non-existent subdir simulates .context before init
		nonExistent := filepath.Join(swapped, ".nonexistent-ctx-dir")
		if err := ValidateBoundary(nonExistent); err != nil {
			t.Errorf("ValidateBoundary(%q) with swapped drive case should pass, got: %v", nonExistent, err)
		}

		// Also test the default relative path that ctx init uses
		if err := ValidateBoundary(".context"); err != nil {
			t.Errorf("ValidateBoundary(.context) with swapped drive case should pass, got: %v", err)
		}
	}
}
