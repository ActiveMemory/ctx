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

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestCheckBackupMarker_NoMarker(t *testing.T) {
	warnings := checkBackupMarker("/nonexistent/path/marker", nil)
	if len(warnings) == 0 {
		t.Fatal("expected warnings for missing marker")
	}
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "No backup marker found") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'No backup marker found' warning, got: %v", warnings)
	}
}

func TestCheckBackupMarker_FreshMarker(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "marker")
	if err := os.WriteFile(marker, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	warnings := checkBackupMarker(marker, nil)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for fresh marker, got: %v", warnings)
	}
}

func TestCheckBackupMarker_StaleMarker(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "marker")
	if err := os.WriteFile(marker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	// Set mtime to 3 days ago
	staleTime := time.Now().Add(-3 * 24 * time.Hour)
	if err := os.Chtimes(marker, staleTime, staleTime); err != nil {
		t.Fatal(err)
	}

	warnings := checkBackupMarker(marker, nil)
	if len(warnings) == 0 {
		t.Fatal("expected warnings for stale marker")
	}
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "days old") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'days old' warning, got: %v", warnings)
	}
}

func TestCheckSMBMount_NoURL(t *testing.T) {
	warnings := checkSMBMountWarnings("", nil)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for empty URL, got: %v", warnings)
	}
}

func TestCheckSMBMount_Unmounted(t *testing.T) {
	// Use a URL that definitely won't have a GVFS mount
	warnings := checkSMBMountWarnings("smb://testhost/testshare", nil)
	if len(warnings) == 0 {
		t.Fatal("expected warnings for unmounted SMB share")
	}
	found := false
	for _, w := range warnings {
		if strings.Contains(w, "not mounted") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'not mounted' warning, got: %v", warnings)
	}
}

func TestCheckBackupAge_DailyThrottle(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	throttleFile := filepath.Join(secureTempDir(), backupThrottleID)
	touchFile(throttleFile)

	cmd := newTestCmd()
	if err := runCheckBackupAge(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Backup Warning") {
		t.Errorf("expected silence on throttled run, got: %s", out)
	}
}

func TestCheckBackupAge_StaleMarkerEmitsWarning(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Create a stale marker at the expected location
	home := t.TempDir()
	t.Setenv("HOME", home)
	stateDir := filepath.Join(home, ".local", "state")
	if err := os.MkdirAll(stateDir, 0o750); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(stateDir, config.BackupMarkerFile)
	if err := os.WriteFile(marker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	staleTime := time.Now().Add(-3 * 24 * time.Hour)
	if err := os.Chtimes(marker, staleTime, staleTime); err != nil {
		t.Fatal(err)
	}

	// Ensure no SMB URL (only age check)
	t.Setenv("CTX_BACKUP_SMB_URL", "")

	cmd := newTestCmd()
	if err := runCheckBackupAge(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Backup Warning") {
		t.Errorf("expected Backup Warning, got: %s", out)
	}
	if !strings.Contains(out, "days old") {
		t.Errorf("expected 'days old' in output, got: %s", out)
	}
}

func TestCheckBackupAge_FreshMarkerSilent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	home := t.TempDir()
	t.Setenv("HOME", home)
	stateDir := filepath.Join(home, ".local", "state")
	if err := os.MkdirAll(stateDir, 0o750); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(stateDir, config.BackupMarkerFile)
	if err := os.WriteFile(marker, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CTX_BACKUP_SMB_URL", "")

	cmd := newTestCmd()
	if err := runCheckBackupAge(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Backup Warning") {
		t.Errorf("expected silence for fresh marker, got: %s", out)
	}
}

func TestCheckBackupAge_MissingMarkerWarns(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Use a home dir with no marker file
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CTX_BACKUP_SMB_URL", "")

	cmd := newTestCmd()
	if err := runCheckBackupAge(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Backup Warning") {
		t.Errorf("expected Backup Warning for missing marker, got: %s", out)
	}
	if !strings.Contains(out, "No backup marker found") {
		t.Errorf("expected 'No backup marker found', got: %s", out)
	}
}
