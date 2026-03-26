//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
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

	"github.com/ActiveMemory/ctx/internal/cli/recall/core/extract"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/index"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/entity"
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
}

func TestRunRecallImport_ArgValidation(t *testing.T) {
	// --all with a session ID should error
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"import", "--all", "some-session"})
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
	cmd3.SetArgs([]string{"import", "--regenerate", "some-session"})
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

func TestRunRecallImport_SingleSession(t *testing.T) {
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
	cmd.SetArgs([]string{"import", "export-session", "--all-projects"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Imported") || !strings.Contains(output, "session") {
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
			content, readErr := os.ReadFile(filepath.Join(journalDir, e.Name())) //nolint:gosec // test temp path
			if readErr != nil {
				t.Fatalf("read journal file: %v", readErr)
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

func TestRunRecallImport_DedupRenamesOldFile(t *testing.T) {
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
	cmd.SetArgs([]string{"import", "--all", "--all-projects", "--regenerate", "--yes"})

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
	if _, statErr := os.Stat(filepath.Join(journalDir, oldFilename)); statErr == nil {
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

// importHelper runs "recall import --all --all-projects" in a temp dir and
// returns the journal directory and the name of the first imported .md file.
func importHelper(t *testing.T, tmpDir string, extraArgs ...string) (journalDir string, mdFile string) {
	t.Helper()

	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	args := append([]string{"import", "--all", "--all-projects"}, extraArgs...)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("import: %v\noutput: %s", err, buf.String())
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
	t.Fatal("no .md file found after import")
	return "", ""
}

func TestRunRecallImport_PreservesFrontmatter(t *testing.T) {
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

	// First import
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Read the original frontmatter to get the generated title
	origData, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read original: %v", err)
	}
	origTitle := index.ExtractFrontmatterField(string(origData), "title")

	// Inject enriched frontmatter - keep the same title to avoid rename
	enrichedFM := fmt.Sprintf("---\ndate: \"2026-01-20\"\ntitle: %q\nsummary: \"A curated summary\"\ntags:\n  - enriched\n---\n", origTitle)
	body := "# hello from test\n\nBody content\n"
	if writeErr := os.WriteFile(path, []byte(enrichedFM+"\n"+body), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-import with --regenerate (safe default skips existing; we need regenerate to trigger re-export)
	importHelper(t, tmpDir, "--regenerate", "--yes")

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

func TestRunRecallImport_KeepFrontmatterFalseDiscards(t *testing.T) {
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

	// First import
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Read the original frontmatter to get the generated title
	origData, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read original: %v", err)
	}
	origTitle := index.ExtractFrontmatterField(string(origData), "title")

	// Inject enriched frontmatter - keep the same title to avoid rename
	enrichedFM := fmt.Sprintf("---\ndate: \"2026-01-20\"\ntitle: %q\nsummary: \"A curated summary\"\ntags:\n  - enriched\n---\n", origTitle)
	body := "# hello from test\n\nBody content\n"
	if writeErr := os.WriteFile(path, []byte(enrichedFM+"\n"+body), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-import with --keep-frontmatter=false - should discard enriched frontmatter
	importHelper(t, tmpDir, "--keep-frontmatter=false", "--yes")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	content := string(data)

	if strings.Contains(content, "A curated summary") {
		t.Error("--keep-frontmatter=false should discard enriched frontmatter summary")
	}
	if strings.Contains(content, "tags:") {
		t.Error("--keep-frontmatter=false should discard enriched frontmatter tags")
	}
	// File should still have session content
	if !strings.Contains(content, "session_id:") {
		t.Error("re-exported file should contain session_id in generated frontmatter")
	}
}

func TestRunRecallImport_KeepFrontmatterFalseResetsEnrichmentState(t *testing.T) {
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

	// First import
	journalDir, mdFile := importHelper(t, tmpDir)

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
	if !jstate.Enriched(mdFile) {
		t.Fatal("file should be marked enriched before re-export")
	}

	// Re-import with --keep-frontmatter=false
	importHelper(t, tmpDir, "--keep-frontmatter=false", "--yes")

	// Load state again and verify enriched was cleared
	jstate, err = state.Load(journalDir)
	if err != nil {
		t.Fatalf("load state after re-export: %v", err)
	}
	if jstate.Enriched(mdFile) {
		t.Error("re-export with --keep-frontmatter=false should clear enriched state")
	}
	// Exported state should still be set
	if !jstate.Exported(mdFile) {
		t.Error("file should still be marked exported after re-export")
	}
}

func TestRunRecallImport_AllSkipsExistingByDefault(t *testing.T) {
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

	// First import
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Overwrite file body with custom content
	customContent := "my custom notes - safe default\n"
	if err := os.WriteFile(path, []byte(customContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-import with --all (no --regenerate) - should skip existing
	importHelper(t, tmpDir)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != customContent {
		t.Errorf("--all should skip existing by default\ngot:  %q\nwant: %q", string(data), customContent)
	}
}

func TestRunRecallImport_RegenerateReExports(t *testing.T) {
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

	// First import
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Overwrite body
	if err := os.WriteFile(path, []byte("overwritten\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-import with --regenerate --yes
	importHelper(t, tmpDir, "--regenerate", "--yes")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Error("--regenerate should regenerate file content")
	}
}

func TestRunRecallImport_RegenerateRequiresAll(t *testing.T) {
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"import", "--regenerate", "some-session"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error with --regenerate without --all")
	}
	if !strings.Contains(err.Error(), "--regenerate requires --all") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunRecallImport_DryRun(t *testing.T) {
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
	cmd.SetArgs([]string{"import", "--all", "--all-projects", "--dry-run"})

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

func TestRunRecallImport_DryRunRegenerate(t *testing.T) {
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

	// First import to create the file
	importHelper(t, tmpDir)

	// Dry-run with --regenerate
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"import", "--all", "--all-projects", "--regenerate", "--dry-run"})

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

func TestRunRecallImport_BareExportPrintsHelp(t *testing.T) {
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"import"})

	// Bare export should print help, not error
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("bare export should not error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Import AI sessions") {
		t.Errorf("bare export should print help text, got:\n%s", output)
	}
}

func TestRunRecallImport_SingleSessionAlwaysWrites(t *testing.T) {
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

	// First import via single session
	cmd1 := Cmd()
	buf1 := new(bytes.Buffer)
	cmd1.SetOut(buf1)
	cmd1.SetErr(buf1)
	cmd1.SetArgs([]string{"import", "single-write", "--all-projects"})
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

	// Re-import same session by ID - should always regenerate without prompting
	cmd2 := Cmd()
	buf2 := new(bytes.Buffer)
	cmd2.SetOut(buf2)
	cmd2.SetErr(buf2)
	cmd2.SetArgs([]string{"import", "single-write", "--all-projects"})
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

func TestRunRecallImport_YesBypasses(t *testing.T) {
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

	// First import
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Overwrite body
	if err := os.WriteFile(path, []byte("overwritten\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-import with --regenerate --yes (no stdin prompt)
	importHelper(t, tmpDir, "--regenerate", "--yes")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Error("--yes should bypass confirmation and regenerate files")
	}
}

func TestRunRecallImport_LockedSkippedByDefault(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-lockskip")
	createTestSessionJSONL(t, projDir, "sess-lock-016", "lock-skip", "/home/test/lockskip")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import to create the file.
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Lock the entry in state.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state: %v", loadErr)
	}
	jstate.Mark(mdFile, journal.StageLocked)
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		t.Fatalf("save state: %v", saveErr)
	}

	// Overwrite with custom content.
	custom := "locked content - do not touch\n"
	if writeErr := os.WriteFile(path, []byte(custom), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-import with --regenerate --yes - locked file should be skipped.
	importHelper(t, tmpDir, "--regenerate", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	if string(data) != custom {
		t.Error("locked file should not be regenerated")
	}
}

func TestRunRecallImport_LockedSkippedByKeepFrontmatterFalse(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-lockforce")
	createTestSessionJSONL(t, projDir, "sess-lock-017", "lock-force", "/home/test/lockforce")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import.
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Lock the entry.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state: %v", loadErr)
	}
	jstate.Mark(mdFile, journal.StageLocked)
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		t.Fatalf("save state: %v", saveErr)
	}

	// Overwrite.
	custom := "locked content - cannot override\n"
	if writeErr := os.WriteFile(path, []byte(custom), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Even --keep-frontmatter=false --yes should not overwrite a locked file.
	importHelper(t, tmpDir, "--keep-frontmatter=false", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	if string(data) != custom {
		t.Error("locked file should not be overwritten even with --keep-frontmatter=false")
	}
}

func TestRunRecallImport_KeepFrontmatterFalse(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-keepfm")
	createTestSessionJSONL(t, projDir, "sess-keepfm-018", "keepfm-test", "/home/test/keepfm")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import.
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Read generated title to keep filename stable.
	origData, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read original: %v", readErr)
	}
	origTitle := index.ExtractFrontmatterField(string(origData), "title")

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

	// Re-import with --keep-frontmatter=false - discards frontmatter.
	importHelper(t, tmpDir, "--keep-frontmatter=false", "--yes")

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

func TestRunRecallImport_KeepFrontmatterDefault(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-keepdef")
	createTestSessionJSONL(t, projDir, "sess-keepdef-019", "keepdef-test", "/home/test/keepdef")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import.
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	origData, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read original: %v", readErr)
	}
	origTitle := index.ExtractFrontmatterField(string(origData), "title")

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

	// Re-import with --regenerate (--keep-frontmatter defaults to true).
	importHelper(t, tmpDir, "--regenerate", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	if !strings.Contains(string(data), "Preserved") {
		t.Error("--keep-frontmatter=true (default) should preserve frontmatter")
	}
}

func TestRunRecallImport_DryRunShowsLocked(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-drylocked")
	createTestSessionJSONL(t, projDir, "sess-drylk-020", "drylk-test", "/home/test/drylocked")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Import, then lock.
	journalDir, mdFile := importHelper(t, tmpDir)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		t.Fatalf("load state: %v", loadErr)
	}
	jstate.Mark(mdFile, journal.StageLocked)
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		t.Fatalf("save state: %v", saveErr)
	}

	// Dry-run with --regenerate should mention locked.
	cmd := Cmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"import", "--all", "--all-projects",
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
		opts entity.ImportOpts
		want bool
	}{
		{
			name: "defaults",
			opts: entity.ImportOpts{KeepFrontmatter: true},
			want: false,
		},
		{
			name: "keep-frontmatter=false",
			opts: entity.ImportOpts{KeepFrontmatter: false},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.DiscardFrontmatter()
			if got != tt.want {
				t.Errorf("DiscardFrontmatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunRecallImport_FrontmatterLockedSkipsAndPromotesToState(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-fmlock")
	createTestSessionJSONL(t, projDir, "sess-fmlock-022", "fmlock-test", "/home/test/fmlock")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import to create the file.
	journalDir, mdFile := importHelper(t, tmpDir)
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

	// Re-import with --regenerate --yes.
	importHelper(t, tmpDir, "--regenerate", "--yes")

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

func TestRunRecallImport_KeepFrontmatterFalseImpliesRegenerate(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-implyregen")
	createTestSessionJSONL(t, projDir, "sess-implyregen-021", "implyregen-test", "/home/test/implyregen")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import.
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Overwrite body with custom content.
	if err := os.WriteFile(path, []byte("overwritten\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Re-import with --keep-frontmatter=false (no explicit --regenerate).
	// The implication should cause regeneration.
	importHelper(t, tmpDir, "--keep-frontmatter=false", "--yes")

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Error("--keep-frontmatter=false should imply --regenerate")
	}
}

func TestRunRecallImport_MalformedFrontmatterGracefulDegradation(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-malformed")
	createTestSessionJSONL(t, projDir, "sess-malformed-030", "malformed-fm", "/home/test/malformed")

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import to create the file.
	journalDir, mdFile := importHelper(t, tmpDir)
	path := filepath.Join(journalDir, mdFile)

	// Overwrite with malformed YAML frontmatter (unclosed delimiter, invalid YAML).
	malformedContent := "---\ndate: \"2026-01-20\"\ntitle: \"test\"\nsummary: [invalid yaml\n\n# Body\n\nSome content here\n"
	if writeErr := os.WriteFile(path, []byte(malformedContent), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-import with --regenerate --yes - should not crash.
	importHelper(t, tmpDir, "--regenerate", "--yes")

	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	content := string(data)

	// The file should have valid content (regenerated from session data).
	if !strings.Contains(content, "session_id:") {
		t.Error("regenerated file should contain session_id in frontmatter")
	}
	if !strings.Contains(content, "hello from test") {
		t.Error("regenerated file should contain session content")
	}
}

// createLargeTestSessionJSONL writes a JSONL file with the specified number of
// user/assistant message pairs for testing multipart splitting.
func createLargeTestSessionJSONL(t *testing.T, dir, sessionID, slug, cwd string, pairs int) {
	t.Helper()
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	var lines []string
	for i := 0; i < pairs; i++ {
		userUUID := fmt.Sprintf("u%d", i*2+1)
		assistUUID := fmt.Sprintf("u%d", i*2+2)
		ts := fmt.Sprintf("2026-01-20T10:%02d:%02dZ", i/60, i%60)

		userLine := fmt.Sprintf(
			`{"uuid":%q,"sessionId":%q,"slug":%q,"type":"user","timestamp":%q,"cwd":%q,"version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"message %d from user"}]}}`,
			userUUID, sessionID, slug, ts, cwd, i+1,
		)
		assistLine := fmt.Sprintf(
			`{"uuid":%q,"parentUuid":%q,"sessionId":%q,"slug":%q,"type":"assistant","timestamp":%q,"cwd":%q,"version":"2.1.0","message":{"model":"claude-test","role":"assistant","content":[{"type":"text","text":"reply %d from assistant"}],"usage":{"input_tokens":100,"output_tokens":50}}}`,
			assistUUID, userUUID, sessionID, slug, ts, cwd, i+1,
		)
		lines = append(lines, userLine, assistLine)
	}
	content := strings.Join(lines, "\n") + "\n"
	file := filepath.Join(dir, sessionID+".jsonl")
	if writeErr := os.WriteFile(file, []byte(content), 0600); writeErr != nil {
		t.Fatalf("write %s: %v", file, writeErr)
	}
}

func TestRunRecallImport_MultipartFrontmatterPreservation(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projDir := filepath.Join(tmpDir, ".claude", "projects", "-home-test-multipart")
	// 110 pairs = 220 messages, exceeding config.MaxMessagesPerPart (200) → 2 parts.
	createLargeTestSessionJSONL(t, projDir, "sess-multi-031", "multipart-fm", "/home/test/multipart", 110)

	contextDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First import - should produce 2 files (part 1 and part 2).
	journalDir, _ := importHelper(t, tmpDir)

	entries, readErr := os.ReadDir(journalDir)
	if readErr != nil {
		t.Fatalf("read journal dir: %v", readErr)
	}
	var mdFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			mdFiles = append(mdFiles, e.Name())
		}
	}
	if len(mdFiles) < 2 {
		t.Fatalf("expected at least 2 .md files for multipart export, got %d: %v", len(mdFiles), mdFiles)
	}

	// Inject enriched frontmatter into part 1 (the base file without -p2 suffix).
	part1Path := filepath.Join(journalDir, mdFiles[0])
	origData, readErr := os.ReadFile(filepath.Clean(part1Path))
	if readErr != nil {
		t.Fatalf("read part1: %v", readErr)
	}
	origTitle := index.ExtractFrontmatterField(string(origData), "title")

	enrichedFM := fmt.Sprintf(
		"---\ndate: \"2026-01-20\"\ntitle: %q\nsummary: \"Multipart curated summary\"\ntags:\n  - multipart-enriched\n---\n",
		origTitle,
	)
	body := extract.StripFrontmatter(string(origData))
	if writeErr := os.WriteFile(part1Path, []byte(enrichedFM+"\n"+body), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Re-import with --regenerate --yes.
	importHelper(t, tmpDir, "--regenerate", "--yes")

	// Verify part 1 preserved enriched frontmatter.
	data, readErr := os.ReadFile(filepath.Clean(part1Path))
	if readErr != nil {
		t.Fatalf("read part1 after re-export: %v", readErr)
	}
	content := string(data)

	if !strings.Contains(content, "Multipart curated summary") {
		t.Error("part 1 enriched frontmatter summary should be preserved on re-export")
	}
	if !strings.Contains(content, "multipart-enriched") {
		t.Error("part 1 enriched frontmatter tags should be preserved on re-export")
	}

	// Verify part 2 still exists and has valid content.
	part2Path := filepath.Join(journalDir, mdFiles[1])
	data2, readErr := os.ReadFile(filepath.Clean(part2Path))
	if readErr != nil {
		t.Fatalf("read part2 after re-export: %v", readErr)
	}
	if !strings.Contains(string(data2), "session_id:") {
		t.Error("part 2 should contain session_id in frontmatter")
	}
}
