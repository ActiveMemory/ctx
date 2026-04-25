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

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// setupContextDir creates a minimal .context/ for promotion tests.
func setupContextDir(t *testing.T) (string, func()) {
	t.Helper()
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	testctx.Declare(t, workDir)

	contextDir := filepath.Join(workDir, dir.Context)
	if mkErr := os.MkdirAll(contextDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	// Create required context files
	for _, f := range []string{
		ctx.Constitution, ctx.Task, ctx.Decision,
		ctx.Learning, ctx.Convention,
	} {
		content := "# " + strings.TrimSuffix(f, ".md") + "\n\n"
		fPath := filepath.Join(contextDir, f)
		if writeErr := os.WriteFile(fPath, []byte(content), 0o644); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	return contextDir, func() { _ = os.Chdir(origDir) }
}

func TestPromote_Convention(t *testing.T) {
	contextDir, cleanup := setupContextDir(t)
	defer cleanup()

	e := Entry{Text: "always use bun for this project", Kind: EntryList}
	classification := Classification{
		Target:   entry.Convention,
		Keywords: []string{"always use"},
	}

	if promoteErr := Promote(e, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, ctx.Convention))
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

	e := Entry{Text: "learned that nolint is ignored in v2", Kind: EntryParagraph}
	classification := Classification{
		Target:   entry.Learning,
		Keywords: []string{"learned"},
	}

	if promoteErr := Promote(e, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, ctx.Learning))
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

	e := Entry{Text: "decided to use SQLite over Postgres", Kind: EntryParagraph}
	classification := Classification{
		Target:   entry.Decision,
		Keywords: []string{"decided"},
	}

	if promoteErr := Promote(e, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, ctx.Decision))
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

	e := Entry{Text: "need to add tests for import", Kind: EntryList}
	classification := Classification{
		Target:   entry.Task,
		Keywords: []string{"need to"},
	}

	if promoteErr := Promote(e, classification); promoteErr != nil {
		t.Fatalf("Promote: %v", promoteErr)
	}

	data, readErr := os.ReadFile(filepath.Join(contextDir, ctx.Task))
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
