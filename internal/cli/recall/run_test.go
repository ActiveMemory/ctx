//   /    ctx:                         https://ctx.ist
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

	// --regenerate without --all should error
	cmd3 := Cmd()
	buf3 := new(bytes.Buffer)
	cmd3.SetOut(buf3)
	cmd3.SetErr(buf3)
	cmd3.SetArgs([]string{"export", "--regenerate", "some-session"})
	err3 := cmd3.Execute()
	if err3 == nil {
		t.Fatal("expected error with --regenerate without --all")
	}
	if !strings.Contains(err3.Error(), "--regenerate requires --all") {
		t.Errorf("unexpected error: %v", err3)
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
	cmd.SetArgs([]string{"export", "--all", "--all-projects", "--regenerate", "--yes"})

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

	// Re-export with --regenerate (safe default skips existing; we need regenerate to trigger re-export)
	exportHelper(t, tmpDir, "--regenerate", "--yes")

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
	exportHelper(t, tmpDir, "--force", "--yes")

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
	exportHelper(t, tmpDir, "--force", "--yes")

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

	// Re-export with --skip-existing (deprecated but still accepted)
	exportHelper(t, tmpDir, "--skip-existing")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != customContent {
		t.Errorf("--skip-existing should leave file unchanged\ngot:  %q\nwant: %q", string(data), customContent)
	}
}

func TestRunRecallExport_AllSkipsExistingByDefault(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-safeskip")
	createTestSessionJSONL(t, projDir, "sess-safe-010", "safe-skip", "/home/test/safeskip")

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

	// Overwrite file body with custom content
	customContent := "my custom notes - safe default\n"
	if err := os.WriteFile(path, []byte(customContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-export with --all (no --regenerate) — should skip existing
	exportHelper(t, tmpDir)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != customContent {
		t.Errorf("--all should skip existing by default\ngot:  %q\nwant: %q", string(data), customContent)
	}
}

func TestRunRecallExport_RegenerateReExports(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-regenproj")
	createTestSessionJSONL(t, projDir, "sess-regen-011", "regen-test", "/home/test/regenproj")

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

	// Overwrite body
	if err := os.WriteFile(path, []byte("overwritten\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-export with --regenerate --yes
	exportHelper(t, tmpDir, "--regenerate", "--yes")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Error("--regenerate should regenerate file content")
	}
}

func TestRunRecallExport_RegenerateRequiresAll(t *testing.T) {
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"export", "--regenerate", "some-session"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error with --regenerate without --all")
	}
	if !strings.Contains(err.Error(), "--regenerate requires --all") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunRecallExport_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-dryproj")
	createTestSessionJSONL(t, projDir, "sess-dry-012", "dry-run-test", "/home/test/dryproj")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
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
	cmd.SetArgs([]string{"export", "--all", "--all-projects", "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Would") {
		t.Errorf("--dry-run should print 'Would' summary, got:\n%s", output)
	}

	// Verify no files were written
	journalDir := filepath.Join(contextDir, "journal")
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		// Directory may not have any .md files, that's fine
		return
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			t.Errorf("--dry-run should not write files, found: %s", e.Name())
		}
	}
}

func TestRunRecallExport_DryRunRegenerate(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-dryregen")
	createTestSessionJSONL(t, projDir, "sess-dryregen-013", "dryregen-test", "/home/test/dryregen")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export to create the file
	exportHelper(t, tmpDir)

	// Dry-run with --regenerate
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"export", "--all", "--all-projects", "--regenerate", "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Would") {
		t.Errorf("--dry-run should print 'Would' summary, got:\n%s", output)
	}
	if !strings.Contains(output, "regenerate") {
		t.Errorf("--dry-run --regenerate should mention regenerate in summary, got:\n%s", output)
	}
}

func TestRunRecallExport_BareExportPrintsHelp(t *testing.T) {
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"export"})

	// Bare export should print help, not error
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("bare export should not error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Export AI sessions") {
		t.Errorf("bare export should print help text, got:\n%s", output)
	}
}

func TestRunRecallExport_SingleSessionAlwaysWrites(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-singleproj")
	createTestSessionJSONL(t, projDir, "sess-single-014", "single-write", "/home/test/singleproj")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export via single session
	cmd1 := Cmd()
	buf1 := new(bytes.Buffer)
	cmd1.SetOut(buf1)
	cmd1.SetErr(buf1)
	cmd1.SetArgs([]string{"export", "single-write", "--all-projects"})
	if err := cmd1.Execute(); err != nil {
		t.Fatalf("first export: %v", err)
	}

	// Find the exported file
	journalDir := filepath.Join(contextDir, "journal")
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("read journal dir: %v", err)
	}
	var mdFile string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			mdFile = e.Name()
			break
		}
	}
	if mdFile == "" {
		t.Fatal("no .md file found after first export")
	}
	path := filepath.Join(journalDir, mdFile)

	// Overwrite with custom content
	if writeErr := os.WriteFile(path, []byte("custom\n"), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-export same session by ID — should always regenerate without prompting
	cmd2 := Cmd()
	buf2 := new(bytes.Buffer)
	cmd2.SetOut(buf2)
	cmd2.SetErr(buf2)
	cmd2.SetArgs([]string{"export", "single-write", "--all-projects"})
	if execErr := cmd2.Execute(); execErr != nil {
		t.Fatalf("second export: %v", execErr)
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Error("single-session export should always regenerate content")
	}
}

func TestRunRecallExport_YesBypasses(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-yesproj")
	createTestSessionJSONL(t, projDir, "sess-yes-015", "yes-bypass", "/home/test/yesproj")

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

	// Overwrite body
	if err := os.WriteFile(path, []byte("overwritten\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-export with --regenerate --yes (no stdin prompt)
	exportHelper(t, tmpDir, "--regenerate", "--yes")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Error("--yes should bypass confirmation and regenerate files")
	}
}

func TestRunRecallExport_LockedSkippedByDefault(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(
		tmpDir, ".claude", "projects", "-home-test-lockskip",
	)
	createTestSessionJSONL(
		t, projDir, "sess-lock-016", "lock-skip", "/home/test/lockskip",
	)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export to create the file.
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Lock the entry in state.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state: %v", loadErr)
	}
	jstate.Mark(mdFile, "locked")
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		t.Fatalf("save state: %v", saveErr)
	}

	// Overwrite with custom content.
	custom := "locked content — do not touch\n"
	if writeErr := os.WriteFile(path, []byte(custom), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-export with --regenerate --yes — locked file should be skipped.
	exportHelper(t, tmpDir, "--regenerate", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	if string(data) != custom {
		t.Error("locked file should not be regenerated")
	}
}

func TestRunRecallExport_LockedSkippedByForce(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(
		tmpDir, ".claude", "projects", "-home-test-lockforce",
	)
	createTestSessionJSONL(
		t, projDir, "sess-lock-017", "lock-force", "/home/test/lockforce",
	)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export.
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Lock the entry.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state: %v", loadErr)
	}
	jstate.Mark(mdFile, "locked")
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		t.Fatalf("save state: %v", saveErr)
	}

	// Overwrite.
	custom := "locked content — force cannot override\n"
	if writeErr := os.WriteFile(path, []byte(custom), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Even --force --yes should not overwrite a locked file.
	exportHelper(t, tmpDir, "--force", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	if string(data) != custom {
		t.Error("locked file should not be overwritten even with --force")
	}
}

func TestRunRecallExport_KeepFrontmatterFalse(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(
		tmpDir, ".claude", "projects", "-home-test-keepfm",
	)
	createTestSessionJSONL(
		t, projDir, "sess-keepfm-018", "keepfm-test", "/home/test/keepfm",
	)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export.
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Read generated title to keep filename stable.
	origData, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read original: %v", readErr)
	}
	origTitle := extractFrontmatterField(string(origData), "title")

	// Inject enriched frontmatter.
	enrichedFM := fmt.Sprintf(
		"---\ndate: \"2026-01-20\"\ntitle: %q\nsummary: \"Curated\"\n"+
			"tags:\n  - enriched\n---\n",
		origTitle,
	)
	body := "# hello from test\n\nBody content\n"
	if writeErr := os.WriteFile(
		path, []byte(enrichedFM+"\n"+body), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-export with --keep-frontmatter=false — discards frontmatter.
	exportHelper(t, tmpDir, "--keep-frontmatter=false", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	content := string(data)

	if strings.Contains(content, "Curated") {
		t.Error("--keep-frontmatter=false should discard enriched summary")
	}
	if strings.Contains(content, "tags:") {
		t.Error("--keep-frontmatter=false should discard enriched tags")
	}
	if !strings.Contains(content, "session_id:") {
		t.Error("regenerated file should contain session_id")
	}
}

func TestRunRecallExport_KeepFrontmatterDefault(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(
		tmpDir, ".claude", "projects", "-home-test-keepdef",
	)
	createTestSessionJSONL(
		t, projDir, "sess-keepdef-019", "keepdef-test", "/home/test/keepdef",
	)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export.
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	origData, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read original: %v", readErr)
	}
	origTitle := extractFrontmatterField(string(origData), "title")

	// Inject enriched frontmatter.
	enrichedFM := fmt.Sprintf(
		"---\ndate: \"2026-01-20\"\ntitle: %q\nsummary: \"Preserved\"\n---\n",
		origTitle,
	)
	body := "# hello from test\n\nBody content\n"
	if writeErr := os.WriteFile(
		path, []byte(enrichedFM+"\n"+body), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-export with --regenerate (--keep-frontmatter defaults to true).
	exportHelper(t, tmpDir, "--regenerate", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	if !strings.Contains(string(data), "Preserved") {
		t.Error("--keep-frontmatter=true (default) should preserve frontmatter")
	}
}

func TestRunRecallExport_DryRunShowsLocked(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(
		tmpDir, ".claude", "projects", "-home-test-drylocked",
	)
	createTestSessionJSONL(
		t, projDir, "sess-drylk-020", "drylk-test", "/home/test/drylocked",
	)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Export, then lock.
	journalDir, mdFile := exportHelper(t, tmpDir)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state: %v", loadErr)
	}
	jstate.Mark(mdFile, "locked")
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		t.Fatalf("save state: %v", saveErr)
	}

	// Dry-run with --regenerate should mention locked.
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"export", "--all", "--all-projects",
		"--regenerate", "--dry-run",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("dry-run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "locked") {
		t.Errorf("dry-run should mention locked entries, got:\n%s", output)
	}
}

func TestDiscardFrontmatter(t *testing.T) {
	tests := []struct {
		name string
		opts exportOpts
		want bool
	}{
		{
			name: "defaults",
			opts: exportOpts{keepFrontmatter: true},
			want: false,
		},
		{
			name: "keep-frontmatter=false",
			opts: exportOpts{keepFrontmatter: false},
			want: true,
		},
		{
			name: "force overrides keep-frontmatter",
			opts: exportOpts{keepFrontmatter: true, force: true},
			want: true,
		},
		{
			name: "both false and force",
			opts: exportOpts{keepFrontmatter: false, force: true},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.discardFrontmatter()
			if got != tt.want {
				t.Errorf("discardFrontmatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunRecallExport_FrontmatterLockedSkipsAndPromotesToState(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(
		tmpDir, ".claude", "projects", "-home-test-fmlock",
	)
	createTestSessionJSONL(
		t, projDir, "sess-fmlock-022", "fmlock-test", "/home/test/fmlock",
	)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export to create the file.
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Verify the file is NOT locked in state.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state: %v", loadErr)
	}
	if jstate.Locked(mdFile) {
		t.Fatal("file should not be locked in state initially")
	}

	// Manually add "locked: true" to frontmatter (simulating user edit).
	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	content := string(data)
	// Insert locked: true into existing frontmatter.
	content = strings.Replace(content, "---\n", "---\nlocked: true\n", 1)
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-export with --regenerate --yes.
	exportHelper(t, tmpDir, "--regenerate", "--yes")

	// File should be unchanged (locked via frontmatter).
	after, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read after: %v", readErr)
	}
	if !strings.Contains(string(after), "locked: true") {
		t.Error("frontmatter-locked file should not be regenerated")
	}

	// State should now reflect the lock (promoted from frontmatter).
	jstate, loadErr = state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state after export: %v", loadErr)
	}
	if !jstate.Locked(mdFile) {
		t.Error("frontmatter lock should be promoted to state file")
	}
}

func TestRunRecallExport_KeepFrontmatterFalseImpliesRegenerate(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(
		tmpDir, ".claude", "projects", "-home-test-implyregen",
	)
	createTestSessionJSONL(
		t, projDir, "sess-implyregen-021", "implyregen-test", "/home/test/implyregen",
	)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First export.
	journalDir, mdFile := exportHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Overwrite body with custom content.
	if err := os.WriteFile(path, []byte("overwritten\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-export with --keep-frontmatter=false (no explicit --regenerate).
	// The implication should cause regeneration.
	exportHelper(t, tmpDir, "--keep-frontmatter=false", "--yes")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Error("--keep-frontmatter=false should imply --regenerate")
	}
}
