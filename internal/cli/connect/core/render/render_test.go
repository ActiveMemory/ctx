//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestWriteEntries_CreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	ctxDir := filepath.Join(tmpDir, ".context")
	if mkErr := os.MkdirAll(ctxDir, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}

	origDir, _ := os.Getwd()
	if chErr := os.Chdir(tmpDir); chErr != nil {
		t.Fatal(chErr)
	}
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	entries := []hub.EntryMsg{
		{
			Type:      "decision",
			Content:   "Use UTC timestamps",
			Origin:    "alpha",
			Timestamp: 1710422400,
			Sequence:  1,
		},
		{
			Type:      "learning",
			Content:   "Avoid mocks in integration tests",
			Origin:    "beta",
			Timestamp: 1710422401,
			Sequence:  2,
		},
	}

	if writeErr := WriteEntries(entries); writeErr != nil {
		t.Fatalf("WriteEntries: %v", writeErr)
	}

	// Check decisions file.
	decPath := filepath.Join(
		ctxDir, "shared", "decisions.md",
	)
	decData, readErr := os.ReadFile(decPath)
	if readErr != nil {
		t.Fatalf("read decisions: %v", readErr)
	}
	decStr := string(decData)
	if !strings.Contains(decStr, "Use UTC timestamps") {
		t.Error("decisions.md missing content")
	}
	if !strings.Contains(decStr, "**Origin**: alpha") {
		t.Error("decisions.md missing origin tag")
	}

	// Check learnings file.
	learnPath := filepath.Join(
		ctxDir, "shared", "learnings.md",
	)
	learnData, learnErr := os.ReadFile(learnPath)
	if learnErr != nil {
		t.Fatalf("read learnings: %v", learnErr)
	}
	if !strings.Contains(
		string(learnData), "Avoid mocks",
	) {
		t.Error("learnings.md missing content")
	}
}

func TestWriteEntries_AppendsToExisting(t *testing.T) {
	tmpDir := t.TempDir()
	ctxDir := filepath.Join(tmpDir, ".context")
	sharedDir := filepath.Join(ctxDir, "shared")
	if mkErr := os.MkdirAll(sharedDir, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}

	origDir, _ := os.Getwd()
	if chErr := os.Chdir(tmpDir); chErr != nil {
		t.Fatal(chErr)
	}
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	// Pre-populate a file.
	existing := "## Existing content\n\n"
	decPath := filepath.Join(sharedDir, "decisions.md")
	if writeErr := os.WriteFile(
		decPath, []byte(existing), 0644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	entries := []hub.EntryMsg{
		{
			Type:      "decision",
			Content:   "New decision",
			Origin:    "proj",
			Timestamp: 1710422400,
			Sequence:  1,
		},
	}

	if writeErr := WriteEntries(entries); writeErr != nil {
		t.Fatal(writeErr)
	}

	data, _ := os.ReadFile(decPath)
	content := string(data)
	if !strings.Contains(content, "Existing content") {
		t.Error("existing content was overwritten")
	}
	if !strings.Contains(content, "New decision") {
		t.Error("new entry was not appended")
	}
}

func TestTypedFileName(t *testing.T) {
	tests := []struct {
		entryType string
		want      string
	}{
		{"decision", "decisions.md"},
		{"learning", "learnings.md"},
		{"convention", "conventions.md"},
	}
	for _, tt := range tests {
		got := typedFileName(tt.entryType)
		if got != tt.want {
			t.Errorf(
				"typedFileName(%q) = %q, want %q",
				tt.entryType, got, tt.want,
			)
		}
	}
}
