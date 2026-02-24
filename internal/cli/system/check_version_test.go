//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMajorMinor(t *testing.T) {
	tests := []struct {
		ver       string
		wantMajor int
		wantMinor int
		wantOK    bool
	}{
		{"0.6.0", 0, 6, true},
		{"1.2.3", 1, 2, true},
		{"10.20.30", 10, 20, true},
		{"0.6", 0, 6, true},
		{"dev", 0, 0, false},
		{"", 0, 0, false},
		{"abc.def.ghi", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.ver, func(t *testing.T) {
			major, minor, ok := parseMajorMinor(tt.ver)
			if ok != tt.wantOK {
				t.Errorf("parseMajorMinor(%q) ok = %v, want %v", tt.ver, ok, tt.wantOK)
			}
			if ok && (major != tt.wantMajor || minor != tt.wantMinor) {
				t.Errorf("parseMajorMinor(%q) = (%d, %d), want (%d, %d)",
					tt.ver, major, minor, tt.wantMajor, tt.wantMinor)
			}
		})
	}
}

func TestCheckVersion_NotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Change to empty temp dir (no .context/)
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()

	cmd := newTestCmd()
	if err := runCheckVersion(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Version Mismatch") {
		t.Errorf("expected silence when not initialized, got: %s", out)
	}
}

func TestCheckVersion_DailyThrottle(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Create the throttle marker (touched today)
	_ = os.MkdirAll(filepath.Join(tmpDir, "ctx"), 0o700)
	touchFile(filepath.Join(tmpDir, "ctx", "version-checked"))

	cmd := newTestCmd()
	if err := runCheckVersion(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Version Mismatch") {
		t.Errorf("expected silence due to daily throttle, got: %s", out)
	}
}

func TestCheckVersion_SilentOnMatch(t *testing.T) {
	// When binary == plugin version, no output.
	// In test context, bootstrap.version is "dev" which skips the check.
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	if err := runCheckVersion(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Version Mismatch") {
		t.Errorf("expected silence for dev build, got: %s", out)
	}
}
