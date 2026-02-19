//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// testCmd creates a cobra command that captures output.
func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(buf)
	return cmd
}

func setupKnowledgeTest(t *testing.T) (string, func()) {
	t.Helper()
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	rc.Reset()

	ctxDir := filepath.Join(tempDir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0750)

	cleanup := func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	}
	return ctxDir, cleanup
}

func makeEntry(date, title, body string) string {
	return fmt.Sprintf("## [%s-120000] %s\n\n%s", date, title, body)
}

func TestArchiveKnowledgeFile_NoFile(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()
	defer rc.Reset()

	// No .context dir at all
	_ = os.MkdirAll(filepath.Join(tempDir, config.DirContext), 0750)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	_, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 5, false, false,
	)
	if err == nil {
		t.Fatal("expected error when file does not exist")
	}
}

func TestArchiveKnowledgeFile_NoOldEntries(t *testing.T) {
	ctxDir, cleanup := setupKnowledgeTest(t)
	defer cleanup()

	// All entries are recent
	today := time.Now().Format("2006-01-02")
	content := "# Decisions\n\n" + makeEntry(today, "Recent decision", "Body.\n")
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision), []byte(content), 0600)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	count, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 5, false, false,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
}

func TestArchiveKnowledgeFile_DryRun(t *testing.T) {
	ctxDir, cleanup := setupKnowledgeTest(t)
	defer cleanup()

	oldDate := time.Now().AddDate(0, 0, -100).Format("2006-01-02")
	content := "# Decisions\n\n" +
		makeEntry(oldDate, "Old decision", "Body.\n") + "\n" +
		makeEntry(time.Now().Format("2006-01-02"), "Recent decision", "Body.\n")
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision), []byte(content), 0600)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	count, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 1, false, true,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("dry run should return 0, got %d", count)
	}
	if !strings.Contains(buf.String(), "Dry run") {
		t.Error("dry run output should contain 'Dry run'")
	}

	// Verify file not modified
	after, _ := os.ReadFile(filepath.Join(ctxDir, config.FileDecision))
	if string(after) != content {
		t.Error("file should not be modified in dry run")
	}
}

func TestArchiveKnowledgeFile_ArchivesOld(t *testing.T) {
	ctxDir, cleanup := setupKnowledgeTest(t)
	defer cleanup()

	oldDate := time.Now().AddDate(0, 0, -100).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	content := "# Decisions\n\n" +
		makeEntry(oldDate, "Old decision", "Body.\n") + "\n" +
		makeEntry(today, "Recent decision", "Body.\n")
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision), []byte(content), 0600)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	count, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 1, false, false,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}

	// Verify old entry removed from source
	after, _ := os.ReadFile(filepath.Join(ctxDir, config.FileDecision))
	if strings.Contains(string(after), "Old decision") {
		t.Error("archived entry should be removed from source")
	}
	if !strings.Contains(string(after), "Recent decision") {
		t.Error("recent entry should remain in source")
	}

	// Verify archive file created
	archiveDir := filepath.Join(ctxDir, config.DirArchive)
	entries, _ := os.ReadDir(archiveDir)
	if len(entries) == 0 {
		t.Fatal("expected archive file to be created")
	}
	archiveContent, _ := os.ReadFile(filepath.Join(archiveDir, entries[0].Name()))
	if !strings.Contains(string(archiveContent), "Old decision") {
		t.Error("archive should contain the old entry")
	}
}

func TestArchiveKnowledgeFile_KeepRecent(t *testing.T) {
	ctxDir, cleanup := setupKnowledgeTest(t)
	defer cleanup()

	// Create 3 old entries — keepRecent=2 should protect the last 2
	oldDate1 := time.Now().AddDate(0, 0, -120).Format("2006-01-02")
	oldDate2 := time.Now().AddDate(0, 0, -110).Format("2006-01-02")
	oldDate3 := time.Now().AddDate(0, 0, -100).Format("2006-01-02")
	content := "# Decisions\n\n" +
		makeEntry(oldDate1, "Oldest", "Body.\n") + "\n" +
		makeEntry(oldDate2, "Middle", "Body.\n") + "\n" +
		makeEntry(oldDate3, "Newest old", "Body.\n")
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision), []byte(content), 0600)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	count, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 2, false, false,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1 (only oldest should be archived)", count)
	}

	after, _ := os.ReadFile(filepath.Join(ctxDir, config.FileDecision))
	if strings.Contains(string(after), "Oldest") {
		t.Error("oldest entry should be archived")
	}
	if !strings.Contains(string(after), "Middle") {
		t.Error("middle entry should be kept (within keepRecent)")
	}
	if !strings.Contains(string(after), "Newest old") {
		t.Error("newest entry should be kept (within keepRecent)")
	}
}

func TestArchiveKnowledgeFile_ArchiveAll(t *testing.T) {
	ctxDir, cleanup := setupKnowledgeTest(t)
	defer cleanup()

	today := time.Now().Format("2006-01-02")
	content := "# Decisions\n\n" +
		makeEntry(today, "First", "Body.\n") + "\n" +
		makeEntry(today, "Second", "Body.\n") + "\n" +
		makeEntry(today, "Third", "Body.\n")
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision), []byte(content), 0600)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	count, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 1, true, false,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 3 entries, keepRecent=1, so 2 archived
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}

	after, _ := os.ReadFile(filepath.Join(ctxDir, config.FileDecision))
	if !strings.Contains(string(after), "Third") {
		t.Error("most recent entry should be kept")
	}
}

func TestArchiveKnowledgeFile_Superseded(t *testing.T) {
	ctxDir, cleanup := setupKnowledgeTest(t)
	defer cleanup()

	today := time.Now().Format("2006-01-02")
	content := "# Decisions\n\n" +
		makeEntry(today, "Superseded one", "~~Superseded by newer~~\n") + "\n" +
		makeEntry(today, "Current one", "Body.\n") + "\n" +
		makeEntry(today, "Another current", "Body.\n")
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision), []byte(content), 0600)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	count, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 2, false, false,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1 (superseded entry)", count)
	}

	after, _ := os.ReadFile(filepath.Join(ctxDir, config.FileDecision))
	if strings.Contains(string(after), "Superseded one") {
		t.Error("superseded entry should be archived")
	}
}

func TestArchiveKnowledgeFile_AppendSameDay(t *testing.T) {
	ctxDir, cleanup := setupKnowledgeTest(t)
	defer cleanup()

	// Create an existing archive file for today
	archiveDir := filepath.Join(ctxDir, config.DirArchive)
	_ = os.MkdirAll(archiveDir, 0750)
	today := time.Now().Format("2006-01-02")
	archiveFile := filepath.Join(archiveDir, fmt.Sprintf("decisions-%s.md", today))
	_ = os.WriteFile(archiveFile, []byte("# Archived Decisions - "+today+"\n\nExisting content.\n"), 0600)

	oldDate := time.Now().AddDate(0, 0, -100).Format("2006-01-02")
	content := "# Decisions\n\n" +
		makeEntry(oldDate, "Old entry", "Body.\n") + "\n" +
		makeEntry(today, "Recent", "Body.\n")
	_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision), []byte(content), 0600)

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	_, err := ArchiveKnowledgeFile(
		cmd, config.FileDecision, "decisions",
		config.HeadingArchivedDecisions, index.UpdateDecisions,
		90, 1, false, false,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify content appended to existing archive
	archiveContent, _ := os.ReadFile(archiveFile)
	if !strings.Contains(string(archiveContent), "Existing content") {
		t.Error("existing archive content should be preserved")
	}
	if !strings.Contains(string(archiveContent), "Old entry") {
		t.Error("new archived content should be appended")
	}
}
