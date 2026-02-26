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

func TestCheckJournal_NoJournalDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	cmd := newTestCmd()
	if err := runCheckJournal(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Journal Reminder") {
		t.Errorf("expected silence when no journal dir, got: %s", out)
	}
}

func TestCheckJournal_DailyThrottle(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	// Create journal dir and projects dir
	_ = os.MkdirAll(resolvedJournalDir(), 0o750)
	fakeProjectsDir := filepath.Join(tmpDir, "claude-projects")
	_ = os.MkdirAll(fakeProjectsDir, 0o750)
	t.Setenv("HOME", tmpDir)
	_ = os.MkdirAll(filepath.Join(tmpDir, ".claude", "projects"), 0o750)

	// Create the throttle marker (touched today)
	_ = os.MkdirAll(filepath.Join(tmpDir, "ctx"), 0o700)
	touchFile(filepath.Join(tmpDir, "ctx", "journal-reminded"))

	cmd := newTestCmd()
	if err := runCheckJournal(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Journal Reminder") {
		t.Errorf("expected silence due to daily throttle, got: %s", out)
	}
}

func TestCheckJournal_Unenriched(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	// Create journal dir with unenriched entry
	_ = os.MkdirAll(resolvedJournalDir(), 0o750)
	_ = os.WriteFile(filepath.Join(resolvedJournalDir(), "2026-01-01-test.md"),
		[]byte("# No frontmatter here"), 0o600)

	// Create Claude projects dir
	_ = os.MkdirAll(filepath.Join(tmpDir, ".claude", "projects"), 0o750)

	cmd := newTestCmd()
	if err := runCheckJournal(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Journal Reminder") {
		t.Errorf("expected journal reminder, got: %s", out)
	}
	if !strings.Contains(out, "entries need enrichment") {
		t.Errorf("expected unenriched message, got: %s", out)
	}
	if !strings.Contains(out, "Context:") {
		t.Errorf("expected context dir footer, got: %s", out)
	}
}

func TestCheckJournal_BothStages(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)
	// Create old journal entry (unenriched) with old mtime
	_ = os.MkdirAll(resolvedJournalDir(), 0o750)
	jDir := resolvedJournalDir()
	_ = os.WriteFile(filepath.Join(jDir, "2025-01-01-test.md"),
		[]byte("# Old entry"), 0o600)
	oldTime := time.Now().Add(-48 * time.Hour)
	_ = os.Chtimes(filepath.Join(jDir, "2025-01-01-test.md"), oldTime, oldTime)

	// Create newer JSONL file (unexported session)
	projectsDir := filepath.Join(tmpDir, ".claude", "projects", "test")
	_ = os.MkdirAll(projectsDir, 0o750)
	_ = os.WriteFile(filepath.Join(projectsDir, "session.jsonl"),
		[]byte(`{"type":"test"}`), 0o600)

	cmd := newTestCmd()
	if err := runCheckJournal(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "not yet exported") {
		t.Errorf("expected export message, got: %s", out)
	}
	if !strings.Contains(out, "entries need enrichment") {
		t.Errorf("expected enrichment message, got: %s", out)
	}
}

func TestCountUnenriched(t *testing.T) {
	dir := t.TempDir()

	// Enriched file (has state entry)
	_ = os.WriteFile(filepath.Join(dir, "enriched.md"),
		[]byte("---\ntitle: test\n---\n# Content"), 0o600)

	// Unenriched file (no state entry)
	_ = os.WriteFile(filepath.Join(dir, "raw.md"),
		[]byte("# Just content"), 0o600)

	// Non-md file (should be ignored)
	_ = os.WriteFile(filepath.Join(dir, "notes.txt"),
		[]byte("not markdown"), 0o600)

	// Create state file marking enriched.md
	stateJSON := `{"version":1,"entries":{"enriched.md":{"enriched":"2026-01-21"}}}`
	_ = os.WriteFile(filepath.Join(dir, ".state.json"),
		[]byte(stateJSON), 0o600)

	count := countUnenriched(dir)
	if count != 1 {
		t.Errorf("countUnenriched() = %d, want 1", count)
	}
}
