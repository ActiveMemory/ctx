//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/spf13/cobra"
)

func TestCreateArchive(t *testing.T) {
	tmpDir := t.TempDir()

	// Build a small directory tree.
	ctxDir := filepath.Join(tmpDir, ".context")
	if mkErr := os.MkdirAll(ctxDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(ctxDir, "TASKS.md"), []byte("tasks"), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	archive := filepath.Join(tmpDir, "test.tar.gz")
	entries := []archiveEntry{
		{SourcePath: ctxDir, Prefix: ".context"},
	}

	cmd := newTestCmd()
	if archiveErr := createArchive(archive, entries, cmd); archiveErr != nil {
		t.Fatalf("createArchive failed: %v", archiveErr)
	}

	names := tarEntryNames(t, archive)
	if !containsEntry(names, ".context/TASKS.md") {
		t.Errorf("expected .context/TASKS.md in archive, got: %v", names)
	}
}

func TestCreateArchive_MissingDir(t *testing.T) {
	tmpDir := t.TempDir()
	archive := filepath.Join(tmpDir, "test.tar.gz")

	entries := []archiveEntry{
		{SourcePath: filepath.Join(tmpDir, "nonexistent"), Prefix: "nope", Optional: true},
	}

	cmd := newTestCmd()
	if archiveErr := createArchive(archive, entries, cmd); archiveErr != nil {
		t.Fatalf("expected no error for optional missing dir, got: %v", archiveErr)
	}
}

func TestCreateArchive_MissingDirRequired(t *testing.T) {
	tmpDir := t.TempDir()
	archive := filepath.Join(tmpDir, "test.tar.gz")

	entries := []archiveEntry{
		{SourcePath: filepath.Join(tmpDir, "nonexistent"), Prefix: "nope"},
	}

	cmd := newTestCmd()
	archiveErr := createArchive(archive, entries, cmd)
	if archiveErr == nil {
		t.Fatal("expected error for required missing dir")
	}
}

func TestCreateArchive_Exclusions(t *testing.T) {
	tmpDir := t.TempDir()

	// Build tree with excluded subdir.
	ctxDir := filepath.Join(tmpDir, ".context")
	excludedDir := filepath.Join(ctxDir, "journal-site")
	if mkErr := os.MkdirAll(excludedDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(ctxDir, "TASKS.md"), []byte("tasks"), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(excludedDir, "index.html"), []byte("<html>"), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	archive := filepath.Join(tmpDir, "test.tar.gz")
	entries := []archiveEntry{
		{SourcePath: ctxDir, Prefix: ".context", ExcludeDir: "journal-site"},
	}

	cmd := newTestCmd()
	if archiveErr := createArchive(archive, entries, cmd); archiveErr != nil {
		t.Fatalf("createArchive failed: %v", archiveErr)
	}

	names := tarEntryNames(t, archive)
	if !containsEntry(names, ".context/TASKS.md") {
		t.Errorf("expected .context/TASKS.md in archive, got: %v", names)
	}
	for _, name := range names {
		if strings.Contains(name, "journal-site") {
			t.Errorf("excluded dir journal-site found in archive: %s", name)
		}
	}
}

func TestCreateArchive_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()

	bashrc := filepath.Join(tmpDir, ".bashrc")
	if writeErr := os.WriteFile(bashrc, []byte("export PATH"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	archive := filepath.Join(tmpDir, "test.tar.gz")
	entries := []archiveEntry{
		{SourcePath: bashrc, Prefix: ".bashrc"},
	}

	cmd := newTestCmd()
	if archiveErr := createArchive(archive, entries, cmd); archiveErr != nil {
		t.Fatalf("createArchive failed: %v", archiveErr)
	}

	names := tarEntryNames(t, archive)
	if !containsEntry(names, ".bashrc") {
		t.Errorf("expected .bashrc in archive, got: %v", names)
	}
}

func TestParseSMBConfig(t *testing.T) {
	cfg, cfgErr := parseSMBConfig("smb://myhost/myshare", "backups")
	if cfgErr != nil {
		t.Fatalf("unexpected error: %v", cfgErr)
	}
	if cfg.Host != "myhost" {
		t.Errorf("expected host myhost, got %s", cfg.Host)
	}
	if cfg.Share != "myshare" {
		t.Errorf("expected share myshare, got %s", cfg.Share)
	}
	if cfg.Subdir != "backups" {
		t.Errorf("expected subdir backups, got %s", cfg.Subdir)
	}
	if !strings.Contains(cfg.GVFSPath, "server=myhost,share=myshare") {
		t.Errorf("GVFS path missing server/share: %s", cfg.GVFSPath)
	}
}

func TestParseSMBConfig_DefaultSubdir(t *testing.T) {
	cfg, cfgErr := parseSMBConfig("smb://host/share", "")
	if cfgErr != nil {
		t.Fatalf("unexpected error: %v", cfgErr)
	}
	if cfg.Subdir != config.BackupDefaultSubdir {
		t.Errorf("expected default subdir %s, got %s",
			config.BackupDefaultSubdir, cfg.Subdir)
	}
}

func TestParseSMBConfig_InvalidURL(t *testing.T) {
	_, cfgErr := parseSMBConfig("://bad", "")
	if cfgErr == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestParseSMBConfig_MissingShare(t *testing.T) {
	_, cfgErr := parseSMBConfig("smb://host/", "")
	if cfgErr == nil {
		t.Fatal("expected error for missing share")
	}
}

func TestRunBackup_NoSMB(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up fake home.
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(config.EnvBackupSMBURL, "")
	t.Setenv(config.EnvBackupSMBSubdir, "")

	// Create minimal .context and .claude directories.
	origDir, _ := os.Getwd()
	if chErr := os.Chdir(tmpDir); chErr != nil {
		t.Fatal(chErr)
	}
	defer func() { _ = os.Chdir(origDir) }()

	ctxDir := filepath.Join(tmpDir, ".context")
	if mkErr := os.MkdirAll(ctxDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(ctxDir, "TASKS.md"), []byte("tasks"), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	claudeDir := filepath.Join(tmpDir, ".claude")
	if mkErr := os.MkdirAll(claudeDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}

	// Also create ~/.claude for global backup.
	homeClaude := filepath.Join(home, ".claude")
	if mkErr := os.MkdirAll(homeClaude, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(homeClaude, "settings.json"), []byte("{}"), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Create ~/.bashrc.
	if writeErr := os.WriteFile(
		filepath.Join(home, ".bashrc"), []byte("export PATH"), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Build the command.
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	backupC := backupCmd()
	backupC.SetOut(buf)
	backupC.SetArgs([]string{"--scope", "project"})
	if execErr := backupC.Execute(); execErr != nil {
		t.Fatalf("backup execute: %v", execErr)
	}

	out := buf.String()
	if !strings.Contains(out, "project:") {
		t.Errorf("expected 'project:' in output, got: %s", out)
	}

	// Verify marker file was touched.
	markerPath := filepath.Join(home, ".local", "state", config.BackupMarkerFile)
	if _, statErr := os.Stat(markerPath); os.IsNotExist(statErr) {
		t.Error("expected marker file to be created")
	}
}

func TestRunBackup_InvalidScope(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := backupCmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--scope", "invalid"})

	execErr := cmd.Execute()
	if execErr == nil {
		t.Fatal("expected error for invalid scope")
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
	}
	for _, tt := range tests {
		got := formatSize(tt.bytes)
		if got != tt.want {
			t.Errorf("formatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

// tarEntryNames reads a tar.gz file and returns all entry names.
func tarEntryNames(t *testing.T, path string) []string {
	t.Helper()
	f, openErr := os.Open(path)
	if openErr != nil {
		t.Fatal(openErr)
	}
	defer func() { _ = f.Close() }()

	gz, gzErr := gzip.NewReader(f)
	if gzErr != nil {
		t.Fatal(gzErr)
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	var names []string
	for {
		hdr, nextErr := tr.Next()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			t.Fatal(nextErr)
		}
		names = append(names, hdr.Name)
	}
	return names
}

// containsEntry checks if a name appears in the list.
func containsEntry(names []string, target string) bool {
	for _, n := range names {
		if n == target {
			return true
		}
	}
	return false
}
