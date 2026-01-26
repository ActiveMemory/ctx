//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInitCommand tests the init command creates the .context directory.
func TestInitCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save and restore working directory
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Run the init command
	cmd := Cmd()
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Check that .context directory was created
	ctxDir := filepath.Join(tmpDir, ".context")
	info, err := os.Stat(ctxDir)
	if err != nil {
		t.Fatalf(".context directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal(".context should be a directory")
	}

	// Check that required files exist
	requiredFiles := []string{
		"CONSTITUTION.md",
		"TASKS.md",
		"DECISIONS.md",
		"CONVENTIONS.md",
		"ARCHITECTURE.md",
	}

	for _, name := range requiredFiles {
		path := filepath.Join(ctxDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("required file %s was not created", name)
		}
	}
}

// TestFindInsertionPoint tests the insertion point logic for merging.
func TestFindInsertionPoint(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantPos  int
		wantDesc string // description of expected position
	}{
		{
			name:     "H1 at start",
			content:  "# My Project\n\nSome content here.",
			wantPos:  14, // after "# My Project\n"
			wantDesc: "after H1",
		},
		{
			name:     "H1 with blank lines after",
			content:  "# My Project\n\n\nSome content here.",
			wantPos:  15, // after "# My Project\n\n"
			wantDesc: "after H1 and blank lines",
		},
		{
			name:     "H2 first",
			content:  "## Section One\n\nContent here.",
			wantPos:  0,
			wantDesc: "at top (H2 is not H1)",
		},
		{
			name:     "No heading",
			content:  "Just some text content.\n\nMore text.",
			wantPos:  0,
			wantDesc: "at top (no heading)",
		},
		{
			name:     "Empty file",
			content:  "",
			wantPos:  0,
			wantDesc: "at top (empty)",
		},
		{
			name:     "Only whitespace",
			content:  "\n\n   \n",
			wantPos:  0,
			wantDesc: "at top (only whitespace)",
		},
		{
			name:     "H1 after blank lines",
			content:  "\n\n# Title\n\nContent",
			wantPos:  11, // after "\n\n# Title\n"
			wantDesc: "after H1 (skipping leading blanks)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findInsertionPoint(tt.content)
			if got != tt.wantPos {
				t.Errorf("findInsertionPoint() = %d, want %d (%s)", got, tt.wantPos, tt.wantDesc)
			}
		})
	}
}

// TestInitMergeInsertsAfterH1 tests that merge inserts ctx content after H1.
func TestInitMergeInsertsAfterH1(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-merge-h1-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create CLAUDE.md with H1 but no ctx markers
	existingContent := `# My Amazing Project

This is the project description.

## Build Instructions

Run make build.
`
	if err := os.WriteFile("CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	// Run init with --merge flag
	initCmd := Cmd()
	initCmd.SetArgs([]string{"--merge"})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Check CLAUDE.md content
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	contentStr := string(content)

	// H1 should still be at the start
	if !strings.HasPrefix(contentStr, "# My Amazing Project") {
		t.Error("H1 heading should remain at the start")
	}

	// ctx content should appear before "## Build Instructions"
	ctxIdx := strings.Index(contentStr, "ctx:context")
	buildIdx := strings.Index(contentStr, "## Build Instructions")
	if ctxIdx == -1 {
		t.Fatal("ctx:context marker not found")
	}
	if buildIdx == -1 {
		t.Fatal("Build Instructions section not found")
	}
	if ctxIdx > buildIdx {
		t.Error("ctx content should appear before Build Instructions, not after")
	}

	// Original content should be preserved
	if !strings.Contains(contentStr, "Run make build") {
		t.Error("original content was lost")
	}
}

// TestInitMergeInsertsAtTopWhenNoH1 tests merge inserts at top without H1.
func TestInitMergeInsertsAtTopWhenNoH1(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-merge-no-h1-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create CLAUDE.md without H1 (starts with H2)
	existingContent := `## Build Instructions

Run make build.

## Testing

Run make test.
`
	if err := os.WriteFile("CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	// Run init with --merge flag
	initCmd := Cmd()
	initCmd.SetArgs([]string{"--merge"})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Check CLAUDE.md content
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	contentStr := string(content)

	// ctx content should appear before "## Build Instructions"
	ctxIdx := strings.Index(contentStr, "ctx:context")
	buildIdx := strings.Index(contentStr, "## Build Instructions")
	if ctxIdx == -1 {
		t.Fatal("ctx:context marker not found")
	}
	if buildIdx == -1 {
		t.Fatal("Build Instructions section not found")
	}
	if ctxIdx > buildIdx {
		t.Error("ctx content should appear at top, before Build Instructions")
	}

	// Original content should be preserved
	if !strings.Contains(contentStr, "Run make test") {
		t.Error("original content was lost")
	}
}

// TestInitWithExistingClaudeMdWithCtxMarker tests init when CLAUDE.md
// already exists with ctx marker.
func TestInitWithExistingClaudeMdWithCtxMarker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-existing-claude-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create existing CLAUDE.md with ctx marker already present
	existingContent := `# My Project

This is my existing CLAUDE.md content.

<!-- ctx:context -->
Old ctx content here
<!-- ctx:end -->

## Custom Section

Some custom content here.
`
	if err := os.WriteFile("CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	// Run init
	initCmd := Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Check that CLAUDE.md was updated
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	// Should still contain ctx marker (updated)
	if !strings.Contains(string(content), "ctx:context") {
		t.Error("CLAUDE.md missing ctx:context marker")
	}

	// Should preserve custom section
	if !strings.Contains(string(content), "Custom Section") {
		t.Error("CLAUDE.md lost custom section")
	}
}
