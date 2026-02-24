//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupMapStalenessTest(t *testing.T, tracking *mapTrackingInfo) string {
	t.Helper()
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	_ = os.Chdir(dir)

	rc.Reset()
	t.Cleanup(rc.Reset)

	ctxDir := filepath.Join(dir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0o750)

	// Create required files so isInitialized() returns true
	for _, f := range config.FilesRequired {
		_ = os.WriteFile(filepath.Join(ctxDir, f), []byte("# "+f+"\n"), 0o600)
	}

	if tracking != nil {
		data, _ := json.Marshal(tracking)
		_ = os.WriteFile(filepath.Join(ctxDir, config.FileMapTracking), data, 0o600)
	}

	// Point XDG_RUNTIME_DIR to a temp location so marker files don't interfere
	tmpState := filepath.Join(dir, "state")
	_ = os.MkdirAll(tmpState, 0o700)
	t.Setenv("XDG_RUNTIME_DIR", tmpState)

	return dir
}

func TestCheckMapStaleness_NoTrackingFile(t *testing.T) {
	setupMapStalenessTest(t, nil)

	cmd := checkMapStalenessCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected silent when no tracking file, got:\n%s", buf.String())
	}
}

func TestCheckMapStaleness_OptedOut(t *testing.T) {
	old := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	setupMapStalenessTest(t, &mapTrackingInfo{OptedOut: true, LastRun: old})

	cmd := checkMapStalenessCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected silent when opted out, got:\n%s", buf.String())
	}
}

func TestCheckMapStaleness_Fresh(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	setupMapStalenessTest(t, &mapTrackingInfo{LastRun: today})

	cmd := checkMapStalenessCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected silent when fresh, got:\n%s", buf.String())
	}
}

func TestCheckMapStaleness_Uninitialized(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	_ = os.Chdir(dir)

	rc.Reset()
	t.Cleanup(rc.Reset)

	cmd := checkMapStalenessCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected silent when uninitialized, got:\n%s", buf.String())
	}
}

func TestCheckMapStaleness_DailyThrottle(t *testing.T) {
	// Use an old date that's >30 days ago so the staleness condition is met
	old := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	dir := setupMapStalenessTest(t, &mapTrackingInfo{LastRun: old})

	// Initialize a git repo with a commit touching internal/ so countModuleCommits > 0
	initGitWithInternalCommit(t, dir, old)

	cmd1 := checkMapStalenessCmd()
	var buf1 bytes.Buffer
	cmd1.SetOut(&buf1)
	_ = cmd1.RunE(cmd1, nil)

	if buf1.Len() == 0 {
		t.Fatal("first call should produce output")
	}

	// Second call same day should be silent (throttled)
	cmd2 := checkMapStalenessCmd()
	var buf2 bytes.Buffer
	cmd2.SetOut(&buf2)
	_ = cmd2.RunE(cmd2, nil)

	if buf2.Len() != 0 {
		t.Errorf("second call should be silent (throttled), got:\n%s", buf2.String())
	}
}

func TestCheckMapStaleness_StaleWithCommits(t *testing.T) {
	old := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	dir := setupMapStalenessTest(t, &mapTrackingInfo{LastRun: old})

	initGitWithInternalCommit(t, dir, old)

	cmd := checkMapStalenessCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "VERBATIM") {
		t.Error("expected VERBATIM relay header")
	}
	if !strings.Contains(out, "Architecture Map Stale") {
		t.Error("expected 'Architecture Map Stale' in output")
	}
	if !strings.Contains(out, "/ctx-map") {
		t.Error("expected /ctx-map suggestion in output")
	}
}

// initGitWithInternalCommit creates a git repo with a commit touching internal/.
func initGitWithInternalCommit(t *testing.T, dir, since string) {
	t.Helper()

	internalDir := filepath.Join(dir, "internal", "testpkg")
	_ = os.MkdirAll(internalDir, 0o750)
	_ = os.WriteFile(filepath.Join(internalDir, "test.go"), []byte("package testpkg\n"), 0o600)

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "add", "."},
		{"git", "commit", "-m", "init", "--no-gpg-sign"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // test-only: args are hardcoded literals above
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git command %v failed: %v\n%s", args, err, out)
		}
	}
}
