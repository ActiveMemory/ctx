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
	"time"
)

func TestMarkWrappedUp_CreatesMarker(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	cmd := newTestCmd()
	if err := runMarkWrappedUp(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	markerPath := filepath.Join(tmpDir, "ctx", wrappedUpMarker)
	if _, statErr := os.Stat(markerPath); statErr != nil {
		t.Fatalf("marker file not created: %v", statErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "marked wrapped-up") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

func TestMarkWrappedUp_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	cmd1 := newTestCmd()
	if err := runMarkWrappedUp(cmd1); err != nil {
		t.Fatalf("first run: %v", err)
	}

	cmd2 := newTestCmd()
	if err := runMarkWrappedUp(cmd2); err != nil {
		t.Fatalf("second run: %v", err)
	}

	markerPath := filepath.Join(tmpDir, "ctx", wrappedUpMarker)
	if _, statErr := os.Stat(markerPath); statErr != nil {
		t.Fatalf("marker file missing after second run: %v", statErr)
	}
}

func TestWrappedUpRecently_Fresh(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Create a fresh marker.
	markerPath := filepath.Join(tmpDir, "ctx", wrappedUpMarker)
	_ = os.MkdirAll(filepath.Dir(markerPath), 0o700)
	_ = os.WriteFile(markerPath, []byte("wrapped-up"), 0o600)

	if !wrappedUpRecently() {
		t.Error("expected wrappedUpRecently() = true for fresh marker")
	}
}

func TestWrappedUpRecently_Expired(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Create a marker and backdate it beyond the expiry.
	markerPath := filepath.Join(tmpDir, "ctx", wrappedUpMarker)
	_ = os.MkdirAll(filepath.Dir(markerPath), 0o700)
	_ = os.WriteFile(markerPath, []byte("wrapped-up"), 0o600)

	expired := time.Now().Add(-3 * time.Hour)
	_ = os.Chtimes(markerPath, expired, expired)

	if wrappedUpRecently() {
		t.Error("expected wrappedUpRecently() = false for expired marker")
	}
}

func TestWrappedUpRecently_NoMarker(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	if wrappedUpRecently() {
		t.Error("expected wrappedUpRecently() = false when no marker exists")
	}
}
