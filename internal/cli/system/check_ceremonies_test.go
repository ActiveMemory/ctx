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

func TestCheckCeremonies_NotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	// No .context/ — should be silent
	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence when not initialized, got: %s", out)
	}
}

func TestCheckCeremonies_DailyThrottle(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	setupJournalDir(t, map[string]string{
		"2026-02-01-session.md": "# Session with no ceremonies",
	})

	// Create the throttle marker (touched today)
	_ = os.MkdirAll(filepath.Join(tmpDir, "ctx"), 0o700)
	touchFile(filepath.Join(tmpDir, "ctx", "ceremony-reminded"))

	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Session") {
		t.Errorf("expected silence due to daily throttle, got: %s", out)
	}
}

func TestCheckCeremonies_NoJournalDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	// No journal directory — should be silent

	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence when no journal dir, got: %s", out)
	}
}

func TestCheckCeremonies_BothUsed(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	setupJournalDir(t, map[string]string{
		"2026-02-01-session.md": "Started with /ctx-remember and ended with /ctx-wrap-up",
	})

	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence when both ceremonies used, got: %s", out)
	}
}

func TestCheckCeremonies_BothMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	setupJournalDir(t, map[string]string{
		"2026-02-01-a.md": "# Did some work",
		"2026-02-02-b.md": "# More work",
		"2026-02-03-c.md": "# Even more work",
	})

	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Session Ceremonies") {
		t.Errorf("expected 'Session Ceremonies' nudge, got: %s", out)
	}
	if !strings.Contains(out, "/ctx-remember") {
		t.Errorf("expected /ctx-remember mention, got: %s", out)
	}
	if !strings.Contains(out, "/ctx-wrap-up") {
		t.Errorf("expected /ctx-wrap-up mention, got: %s", out)
	}
	if !strings.Contains(out, "Context:") {
		t.Errorf("expected context dir footer, got: %s", out)
	}
}

func TestCheckCeremonies_OnlyRememberMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	setupJournalDir(t, map[string]string{
		"2026-02-01-session.md": "Ended session with /ctx-wrap-up and persisted context",
	})

	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Session Start") {
		t.Errorf("expected 'Session Start' nudge, got: %s", out)
	}
	if !strings.Contains(out, "/ctx-remember") {
		t.Errorf("expected /ctx-remember mention, got: %s", out)
	}
}

func TestCheckCeremonies_OnlyWrapUpMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	setupJournalDir(t, map[string]string{
		"2026-02-01-session.md": "Started with /ctx-remember and loaded context",
	})

	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Session End") {
		t.Errorf("expected 'Session End' nudge, got: %s", out)
	}
	if !strings.Contains(out, "/ctx-wrap-up") {
		t.Errorf("expected /ctx-wrap-up mention, got: %s", out)
	}
}

func TestCheckCeremonies_CeremoniesAcrossFiles(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	// /ctx-remember in one file, /ctx-wrap-up in another — both found
	setupJournalDir(t, map[string]string{
		"2026-02-01-morning.md": "Started with /ctx-remember",
		"2026-02-02-evening.md": "Ended with /ctx-wrap-up",
	})

	cmd := newTestCmd()
	if err := runCheckCeremonies(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence when both ceremonies found across files, got: %s", out)
	}
}

func TestRecentJournalFiles_SortOrder(t *testing.T) {
	dir := t.TempDir()

	// Create files with date-prefix names in arbitrary order
	for _, name := range []string{
		"2026-01-15-old.md",
		"2026-02-20-newest.md",
		"2026-02-01-middle.md",
		"2026-01-01-oldest.md",
	} {
		_ = os.WriteFile(filepath.Join(dir, name), []byte("# content"), 0o600)
	}

	files := recentJournalFiles(dir, 3)
	if len(files) != 3 {
		t.Fatalf("recentJournalFiles() returned %d files, want 3", len(files))
	}

	// Should be newest first
	wantOrder := []string{"2026-02-20-newest.md", "2026-02-01-middle.md", "2026-01-15-old.md"}
	for i, want := range wantOrder {
		got := filepath.Base(files[i])
		if got != want {
			t.Errorf("files[%d] = %s, want %s", i, got, want)
		}
	}
}

func TestRecentJournalFiles_SkipsNonMd(t *testing.T) {
	dir := t.TempDir()

	_ = os.WriteFile(filepath.Join(dir, "2026-02-01-session.md"), []byte("# md"), 0o600)
	_ = os.WriteFile(filepath.Join(dir, ".state.json"), []byte("{}"), 0o600)
	_ = os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("text"), 0o600)

	files := recentJournalFiles(dir, 10)
	if len(files) != 1 {
		t.Errorf("recentJournalFiles() returned %d files, want 1", len(files))
	}
}

func TestRecentJournalFiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	files := recentJournalFiles(dir, 3)
	if len(files) != 0 {
		t.Errorf("recentJournalFiles() returned %d files, want 0", len(files))
	}
}

func TestRecentJournalFiles_NonexistentDir(t *testing.T) {
	files := recentJournalFiles("/nonexistent/path", 3)
	if files != nil {
		t.Errorf("recentJournalFiles() returned %v, want nil", files)
	}
}

func TestScanJournalsForCeremonies(t *testing.T) {
	tests := []struct {
		name         string
		contents     []string
		wantRemember bool
		wantWrapup   bool
	}{
		{
			name:         "both present",
			contents:     []string{"Used /ctx-remember at start", "Used /ctx-wrap-up at end"},
			wantRemember: true,
			wantWrapup:   true,
		},
		{
			name:         "neither present",
			contents:     []string{"Just did some work", "No ceremonies here"},
			wantRemember: false,
			wantWrapup:   false,
		},
		{
			name:         "only remember",
			contents:     []string{"Ran ctx-remember to load context"},
			wantRemember: true,
			wantWrapup:   false,
		},
		{
			name:         "only wrapup",
			contents:     []string{"Finished with ctx-wrap-up"},
			wantRemember: false,
			wantWrapup:   true,
		},
		{
			name:         "without slash prefix",
			contents:     []string{"ctx-remember and ctx-wrap-up both work"},
			wantRemember: true,
			wantWrapup:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			var files []string
			for i, content := range tt.contents {
				path := filepath.Join(dir, "file"+string(rune('0'+i))+".md")
				_ = os.WriteFile(path, []byte(content), 0o600)
				files = append(files, path)
			}

			remember, wrapup := scanJournalsForCeremonies(files)
			if remember != tt.wantRemember {
				t.Errorf("remember = %v, want %v", remember, tt.wantRemember)
			}
			if wrapup != tt.wantWrapup {
				t.Errorf("wrapup = %v, want %v", wrapup, tt.wantWrapup)
			}
		})
	}
}

// setupJournalDir creates the journal directory with the given files.
func setupJournalDir(t *testing.T, files map[string]string) {
	t.Helper()
	dir := resolvedJournalDir()
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
