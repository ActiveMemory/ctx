//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// setupContextDir creates a minimal .context/ for promotion tests.
func setupContextDir(t *testing.T) (string, func()) {
	t.Helper()
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	rc.Reset()

	contextDir := filepath.Join(workDir, config.DirContext)
	if mkErr := os.MkdirAll(contextDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	// Create required context files
	for _, f := range []string{
		config.FileConstitution, config.FileTask, config.FileDecision,
		config.FileLearning, config.FileConvention,
	} {
		content := "# " + strings.TrimSuffix(f, ".md") + "\n\n"
		if writeErr := os.WriteFile(filepath.Join(contextDir, f), []byte(content), 0o644); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	return contextDir, func() { _ = os.Chdir(origDir) }
}

func TestPromote_Convention(t *testing.T) {
	contextDir, cleanup := setupContextDir(t)
	defer cleanup()

	entry := Entry{Text: "always use bun for this project", Kind: EntryList}
	classification := Classification{Target: config.EntryConvention, Keywords: []string{"always use"}}

	if promoteErr := Promote(entry, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, config.FileConvention))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(data), "always use bun") {
		t.Error("expected convention content in CONVENTIONS.md")
	}
}

func TestPromote_Learning(t *testing.T) {
	contextDir, cleanup := setupContextDir(t)
	defer cleanup()

	entry := Entry{Text: "learned that nolint is ignored in v2", Kind: EntryParagraph}
	classification := Classification{Target: config.EntryLearning, Keywords: []string{"learned"}}

	if promoteErr := Promote(entry, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, config.FileLearning))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(data), "nolint is ignored") {
		t.Error("expected learning content in LEARNINGS.md")
	}
}

func TestPromote_Decision(t *testing.T) {
	contextDir, cleanup := setupContextDir(t)
	defer cleanup()

	entry := Entry{Text: "decided to use SQLite over Postgres", Kind: EntryParagraph}
	classification := Classification{Target: config.EntryDecision, Keywords: []string{"decided"}}

	if promoteErr := Promote(entry, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, config.FileDecision))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(data), "SQLite") {
		t.Error("expected decision content in DECISIONS.md")
	}
}

func TestPromote_Task(t *testing.T) {
	contextDir, cleanup := setupContextDir(t)
	defer cleanup()

	entry := Entry{Text: "need to add tests for import", Kind: EntryList}
	classification := Classification{Target: config.EntryTask, Keywords: []string{"need to"}}

	if promoteErr := Promote(entry, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, config.FileTask))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(data), "add tests for import") {
		t.Error("expected task content in TASKS.md")
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"## Session notes", "Session notes"},
		{"### Key decisions", "Key decisions"},
		{"- always use bun", "always use bun"},
		{"* prefer filepath.Join", "prefer filepath.Join"},
		{"plain text", "plain text"},
		{"## Multi\nline entry", "Multi"},
	}
	for _, tt := range tests {
		got := extractTitle(tt.input)
		if got != tt.want {
			t.Errorf("extractTitle(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
