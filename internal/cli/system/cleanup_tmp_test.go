//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCleanupTmp_RemovesOldFiles(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	ctxTmp := filepath.Join(tmpDir, "ctx")
	_ = os.MkdirAll(ctxTmp, 0o700)

	// Create a 16-day-old file (should be removed)
	oldFile := filepath.Join(ctxTmp, "old-counter")
	_ = os.WriteFile(oldFile, []byte("42"), 0o600)
	oldTime := time.Now().Add(-16 * 24 * time.Hour)
	_ = os.Chtimes(oldFile, oldTime, oldTime)

	// Create a 14-day-old file (should be kept)
	newFile := filepath.Join(ctxTmp, "new-counter")
	_ = os.WriteFile(newFile, []byte("7"), 0o600)
	newTime := time.Now().Add(-14 * 24 * time.Hour)
	_ = os.Chtimes(newFile, newTime, newTime)

	if err := runCleanupTmp(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Old file should be gone
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("expected 16-day-old file to be removed")
	}

	// New file should remain
	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Error("expected 14-day-old file to be kept")
	}
}

func TestCleanupTmp_NoDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", filepath.Join(tmpDir, "nonexistent"))

	// Should not error even if dir doesn't exist
	if err := runCleanupTmp(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCleanupTmp_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	_ = os.MkdirAll(filepath.Join(tmpDir, "ctx"), 0o700)

	if err := runCleanupTmp(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
