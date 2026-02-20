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

	"github.com/ActiveMemory/ctx/internal/journal/state"
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

	// Verify content of exported file.
	// Filename is now title-based (derived from FirstUserMsg "hello from test").
	for _, e := range entries {
		if strings.Contains(e.Name(), "hello-from-test") {
			content, err := os.ReadFile(filepath.Join(journalDir, e.Name())) //nolint:gosec // test temp path
			if err != nil {
				t.Fatalf("read journal file: %v", err)
			}
			if !strings.Contains(string(content), "hello from test") {
				t.Error("journal file missing user message")
			}
			if !strings.Contains(string(content), "session_id:") {
				t.Error("journal file missing session_id in frontmatter")
			}
			return
		}
	}
	t.Errorf("no journal file found matching hello-from-test slug, got: %v", func() []string {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		return names
	}())
}

func TestRunRecallExport_DedupRenamesOldFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-dedupproj")
	sessionID := "dedup123-full-uuid-value"
	createTestSessionJSONL(t, projDir, sessionID, "random-slug", "/home/test/dedupproj")

	// Create .context directory
	contextDir := filepath.Join(tmpDir, ".context")
	journalDir := filepath.Join(contextDir, "journal")
	if err := os.MkdirAll(journalDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Pre-create a legacy file with the old slug-based name (no session_id).
	// The short ID is the first 8 chars of the session ID: "dedup123".
	oldFilename := "2026-01-20-random-slug-dedup123.md"
	oldContent := "---\ndate: \"2026-01-20\"\n---\n\n# random-slug\n\nOld content\n"
	if err := os.WriteFile(filepath.Join(journalDir, oldFilename), []byte(oldContent), 0600); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"export", "--all", "--all-projects"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("read journal dir: %v", err)
	}

	// Should have exactly 1 file (renamed, not duplicated).
	mdFiles := 0
	var fileNames []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			mdFiles++
			fileNames = append(fileNames, e.Name())
		}
	}
	if mdFiles != 1 {
		t.Errorf("expected 1 journal file (deduped), got %d: %v", mdFiles, fileNames)
	}

	// The old file should be gone.
	if _, err := os.Stat(filepath.Join(journalDir, oldFilename)); err == nil {
		t.Error("old file should have been renamed")
	}

	// The new file should have the title-based slug.
	found := false
	for _, name := range fileNames {
		if strings.Contains(name, "hello-from-test") && strings.Contains(name, "dedup123") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected title-based filename with short ID, got: %v", fileNames)
	}
}

// exportHelper runs "recall export --all --all-projects" in a temp dir and
// returns the journal directory and the name of the first exported .md file.
func exportHelper(t *testing.T, tmpDir string, extraArgs ...string) (journalDir string, mdFile string) {
	t.Helper()

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	args := append([]string{"export", "--all", "--all-projects"}, extraArgs...)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("export: %v\noutput: %s", err, buf.String())
	}

	journalDir = filepath.Join(tmpDir, ".context", "journal")
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("read journal dir: %v", err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			return journalDir, e.Name()
		}
	}
	t.Fatal("no .md file found after export")
	return "", ""
}

func TestRunRecallExport_PreservesFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-fmproj")
	createTestSessionJSONL(t, projDir, "sess-fm-001", "fm-preserve", "/home/test/fmproj")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Read the original frontmatter to get the generated title
	origData, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read original: %v", err)
	}
	origTitle := extractFrontmatterField(string(origData), "title")

	// Inject enriched frontmatter — keep the same title to avoid rename
	enrichedFM := fmt.Sprintf("---\ndate: \"2026-01-20\"\ntitle: %q\nsummary: \"A curated summary\"\ntags:\n  - enriched\n---\n", origTitle)
	body := "# hello from test\n\nBody content\n"
	if writeErr := os.WriteFile(path, []byte(enrichedFM+"\n"+body), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-export (default mode — should preserve frontmatter)
	exportHelper(t, tmpDir)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "A curated summary") {
		t.Error("enriched frontmatter summary should be preserved on re-export")
	}
	if !strings.Contains(content, "enriched") {
		t.Error("enriched frontmatter tags should be preserved on re-export")
	}
}

func TestRunRecallExport_ForceDiscardsFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-forceproj")
	createTestSessionJSONL(t, projDir, "sess-force-002", "force-discard", "/home/test/forceproj")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Read the original frontmatter to get the generated title
	origData, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read original: %v", err)
	}
	origTitle := extractFrontmatterField(string(origData), "title")

	// Inject enriched frontmatter — keep the same title to avoid rename
	enrichedFM := fmt.Sprintf("---\ndate: \"2026-01-20\"\ntitle: %q\nsummary: \"A curated summary\"\ntags:\n  - enriched\n---\n", origTitle)
	body := "# hello from test\n\nBody content\n"
	if writeErr := os.WriteFile(path, []byte(enrichedFM+"\n"+body), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-export with --force — should discard enriched frontmatter
	exportHelper(t, tmpDir, "--force")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	content := string(data)

	if strings.Contains(content, "A curated summary") {
		t.Error("--force should discard enriched frontmatter summary")
	}
	if strings.Contains(content, "tags:") {
		t.Error("--force should discard enriched frontmatter tags")
	}
	// File should still have session content
	if !strings.Contains(content, "session_id:") {
		t.Error("re-exported file should contain session_id in generated frontmatter")
	}
}

func TestRunRecallExport_ForceResetsEnrichmentState(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-stateproj")
	createTestSessionJSONL(t, projDir, "sess-state-003", "state-reset", "/home/test/stateproj")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export
	journalDir, mdFile := exportHelper(t, tmpDir)

	// Manually mark the file as enriched in state
	jstate, err := state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	jstate.MarkEnriched(mdFile)
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		t.Fatalf("save state: %v", saveErr)
	}

	// Verify it's marked enriched
	jstate, _ = state.Load(journalDir)
	if !jstate.IsEnriched(mdFile) {
		t.Fatal("file should be marked enriched before --force re-export")
	}

	// Re-export with --force
	exportHelper(t, tmpDir, "--force")

	// Load state again and verify enriched was cleared
	jstate, err = state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state after force: %v", err)
	}
	if jstate.IsEnriched(mdFile) {
		t.Error("--force re-export should clear enriched state")
	}
	// Exported state should still be set
	if !jstate.IsExported(mdFile) {
		t.Error("file should still be marked exported after --force re-export")
	}
}

func TestRunRecallExport_SkipExistingLeavesFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-skipproj")
	createTestSessionJSONL(t, projDir, "sess-skip-004", "skip-existing", "/home/test/skipproj")

	contextDir := filepath.Join(tmpDir, ".context")
	journalDir := filepath.Join(contextDir, "journal")
	if err := os.MkdirAll(journalDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export to discover the filename
	_, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Overwrite the file with custom content
	customContent := "my custom notes - do not overwrite\n"
	if err := os.WriteFile(path, []byte(customContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-export with --skip-existing
	exportHelper(t, tmpDir, "--skip-existing")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != customContent {
		t.Errorf("--skip-existing should leave file unchanged\ngot:  %q\nwant: %q", string(data), customContent)
	}
}
