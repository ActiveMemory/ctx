//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// createTestSessionJSONL writes a minimal valid JSONL file for testing.
func createTestSessionJSONL(t *testing.T, dir, sessionID, slug, cwd string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	line1 := fmt.Sprintf(
		`{"uuid":"u1","sessionId":"%s","slug":"%s","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"%s","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"hello from test"}]}}`,
		sessionID, slug, cwd,
	)
	line2 := fmt.Sprintf(
		`{"uuid":"u2","parentUuid":"u1","sessionId":"%s","slug":"%s","type":"assistant","timestamp":"2026-01-20T10:00:30Z","cwd":"%s","version":"2.1.0","message":{"model":"claude-test","role":"assistant","content":[{"type":"text","text":"hi back"}],"usage":{"input_tokens":100,"output_tokens":50}}}`,
		sessionID, slug, cwd,
	)
	content := line1 + "\n" + line2 + "\n"
	file := filepath.Join(dir, sessionID+".jsonl")
	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatalf("write %s: %v", file, err)
	}
}

func init() {
	// Disable color output in all tests to avoid ANSI codes in assertions.
	color.NoColor = true
}

func TestRunRecallExport_ArgValidation(t *testing.T) {
	// --all with a session ID should error
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"export", "--all", "some-session"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error with --all and session ID")
	}
	if !strings.Contains(err.Error(), "cannot use --all with a session ID") {
		t.Errorf("unexpected error: %v", err)
	}

	// Neither --all nor session ID should error
	cmd2 := Cmd()
	buf2 := new(bytes.Buffer)
	cmd2.SetOut(buf2)
	cmd2.SetErr(buf2)
	cmd2.SetArgs([]string{"export"})
	err2 := cmd2.Execute()
	if err2 == nil {
		t.Fatal("expected error with neither --all nor session ID")
	}
	if !strings.Contains(err2.Error(), "please provide a session ID or use --all") {
		t.Errorf("unexpected error: %v", err2)
	}
}

func TestRunRecallList_NoSessions(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Create the expected directory structure (empty)
	claudeDir := filepath.Join(tmpDir, ".claude", "projects")
	if err := os.MkdirAll(claudeDir, 0750); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"list", "--all-projects"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No sessions found") {
		t.Errorf("expected 'No sessions found' message, got:\n%s", output)
	}
}

func TestRunRecallList_WithSessions(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Create session fixture
	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-myproject")
	createTestSessionJSONL(t, projDir, "sess-list-123", "listing-test-session", "/home/test/myproject")

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"list", "--all-projects"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "listing-test-session") {
		t.Errorf("expected slug in output, got:\n%s", output)
	}
}

func TestRunRecallShow_Latest(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-showproj")
	createTestSessionJSONL(t, projDir, "sess-show-456", "show-test-session", "/home/test/showproj")

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"show", "--latest", "--all-projects"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Verify metadata appears
	if !strings.Contains(output, "show-test-session") {
		t.Errorf("expected slug in output, got:\n%s", output)
	}
	if !strings.Contains(output, "sess-show-456") {
		t.Errorf("expected session ID in output, got:\n%s", output)
	}
}

func TestRunRecallShow_BySlug(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-slugproj")
	createTestSessionJSONL(t, projDir, "sess-slug-789", "unique-slug-name", "/home/test/slugproj")

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"show", "unique-slug", "--all-projects"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "unique-slug-name") {
		t.Errorf("expected slug in output, got:\n%s", output)
	}
}

func TestRunRecallExport_SingleSession(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-expproj")
	createTestSessionJSONL(t, projDir, "sess-exp-aaa", "export-session", "/home/test/expproj")

	// Create .context directory for journal output
	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// We need to be in a directory that has .context/ for the export
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"export", "export-session", "--all-projects"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Exported") || !strings.Contains(output, "session") {
		t.Errorf("expected export confirmation, got:\n%s", output)
	}

	// Verify journal file was created
	journalDir := filepath.Join(contextDir, "journal")
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("read journal dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected at least one journal file")
	}

	// Verify content of exported file
	for _, e := range entries {
		if strings.Contains(e.Name(), "export-session") {
			content, err := os.ReadFile(filepath.Join(journalDir, e.Name()))
			if err != nil {
				t.Fatalf("read journal file: %v", err)
			}
			if !strings.Contains(string(content), "export-session") {
				t.Error("journal file missing session slug")
			}
			if !strings.Contains(string(content), "hello from test") {
				t.Error("journal file missing user message")
			}
			return
		}
	}
	t.Error("no journal file found matching export-session slug")
}
