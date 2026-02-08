//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/context"
)

func TestReportStatus(t *testing.T) {
	tests := []struct {
		name     string
		report   Report
		expected StatusType
	}{
		{
			name:     "no issues",
			report:   Report{},
			expected: StatusOk,
		},
		{
			name: "only warnings",
			report: Report{
				Warnings: []Issue{{File: "test.md", Type: IssueStaleness}},
			},
			expected: StatusWarning,
		},
		{
			name: "only violations",
			report: Report{
				Violations: []Issue{{File: "test.md", Type: IssueSecret}},
			},
			expected: StatusViolation,
		},
		{
			name: "warnings and violations",
			report: Report{
				Warnings:   []Issue{{File: "test.md", Type: IssueStaleness}},
				Violations: []Issue{{File: "test.md", Type: IssueSecret}},
			},
			expected: StatusViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.report.Status()
			if result != tt.expected {
				t.Errorf("Status() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDetect(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "drift-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func(path string) {
		rmErr := os.RemoveAll(path)
		if rmErr != nil {
			fmt.Printf("failed to remove temp dir %q: %v", path, rmErr)
		}
	}(tmpDir)

	// Save and restore the current working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func(dir string) {
		chdirErr := os.Chdir(dir)
		if chdirErr != nil {
			fmt.Printf("failed to chdir: %v", chdirErr)
		}
	}(origDir)

	// Create a .context directory with test files
	ctxDir := filepath.Join(tmpDir, ".context")
	if mkErr := os.Mkdir(ctxDir, 0750); mkErr != nil {
		t.Fatalf("failed to create .context dir: %v", mkErr)
	}

	// Create required files
	files := map[string]string{
		"CONSTITUTION.md": "# Constitution\n\n- [ ] Never break the build\n",
		"TASKS.md":        "# Tasks\n\n- [ ] Do something\n",
		"DECISIONS.md":    "# Decisions\n\n## Decision 1\n\nContent\n",
		"ARCHITECTURE.md": "# Architecture\n\nMain file is `main.go`.\n",
	}

	for name, content := range files {
		path := filepath.Join(ctxDir, name)
		if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
			t.Fatalf("failed to write %s: %v", name, writeErr)
		}
	}

	// Create the main.go file so the path reference check passes
	mainGo := filepath.Join(tmpDir, "main.go")
	if writeErr := os.WriteFile(mainGo, []byte("package main"), 0600); writeErr != nil {
		t.Fatalf("failed to write main.go: %v", writeErr)
	}

	// Load the context
	ctx, err := context.Load(ctxDir)
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	// Run detection
	report := Detect(ctx)

	// Check that no violations exist (no secret files in this test)
	if len(report.Violations) > 0 {
		t.Errorf("expected no violations, got %d", len(report.Violations))
	}

	// Check that passed checks are recorded
	if len(report.Passed) == 0 {
		t.Error("expected at least one passed check")
	}
}

func TestCheckPathReferences(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "drift-path-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func(path string) {
		rmErr := os.RemoveAll(path)
		if rmErr != nil {
			fmt.Printf("failed to remove temp dir %q: %v", path, rmErr)
		}
	}(tmpDir)

	// Save and restore the current working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func(dir string) {
		chdirErr := os.Chdir(dir)
		if chdirErr != nil {
			fmt.Printf("failed to chdir: %v", chdirErr)
		}
	}(origDir)

	// Create a test context with a dead path reference
	ctx := &context.Context{
		Dir: ".context",
		Files: []context.FileInfo{
			{
				Name:    "ARCHITECTURE.md",
				Content: []byte("# Architecture\n\nSee `nonexistent.go` for details.\n"),
			},
		},
	}

	report := &Report{
		Warnings:   []Issue{},
		Violations: []Issue{},
		Passed:     []CheckName{},
	}

	checkPathReferences(ctx, report)

	// Should find the dead path
	if len(report.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(report.Warnings))
	} else {
		if report.Warnings[0].Type != "dead_path" {
			t.Errorf("expected warning type 'dead_path', got %q", report.Warnings[0].Type)
		}
		if report.Warnings[0].Path != "nonexistent.go" {
			t.Errorf("expected path 'nonexistent.go', got %q", report.Warnings[0].Path)
		}
	}
}

func TestCheckStaleness(t *testing.T) {
	tests := []struct {
		name         string
		tasksContent string
		wantWarnings int
	}{
		{
			name:         "few completed tasks",
			tasksContent: "# Tasks\n\n- [x] Done 1\n- [x] Done 2\n- [ ] Todo\n",
			wantWarnings: 0,
		},
		{
			name:         "many completed tasks",
			tasksContent: "# Tasks\n\n- [x] Done 1\n- [x] Done 2\n- [x] Done 3\n- [x] Done 4\n- [x] Done 5\n- [x] Done 6\n- [x] Done 7\n- [x] Done 8\n- [x] Done 9\n- [x] Done 10\n- [x] Done 11\n",
			wantWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &context.Context{
				Dir: ".context",
				Files: []context.FileInfo{
					{
						Name:    "TASKS.md",
						Content: []byte(tt.tasksContent),
					},
				},
			}

			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []CheckName{},
			}

			checkStaleness(ctx, report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
			}
		})
	}
}

func TestCheckRequiredFiles(t *testing.T) {
	tests := []struct {
		name         string
		files        []string
		wantWarnings int
	}{
		{
			name:         "all required files present",
			files:        []string{"CONSTITUTION.md", "TASKS.md", "DECISIONS.md"},
			wantWarnings: 0,
		},
		{
			name:         "missing CONSTITUTION.md",
			files:        []string{"TASKS.md", "DECISIONS.md"},
			wantWarnings: 1,
		},
		{
			name:         "missing all required files",
			files:        []string{"OTHER.md"},
			wantWarnings: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fileInfos []context.FileInfo
			for _, name := range tt.files {
				fileInfos = append(fileInfos, context.FileInfo{Name: name})
			}

			ctx := &context.Context{
				Dir:   ".context",
				Files: fileInfos,
			}

			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []CheckName{},
			}

			checkRequiredFiles(ctx, report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
			}
		})
	}
}

func TestIsTemplateFile(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "empty file",
			content:  []byte{},
			expected: false,
		},
		{
			name:     "regular content",
			content:  []byte("DATABASE_URL=postgres://localhost/db"),
			expected: false,
		},
		{
			name:     "template with YOUR_",
			content:  []byte("API_KEY=YOUR_API_KEY_HERE"),
			expected: true,
		},
		{
			name:     "template with REPLACE_",
			content:  []byte("SECRET=REPLACE_WITH_SECRET"),
			expected: true,
		},
		{
			name:     "template with TODO:",
			content:  []byte("# TODO: Add your config here"),
			expected: true,
		},
		{
			name:     "template with CHANGEME",
			content:  []byte("PASSWORD=CHANGEME"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTemplateFile(tt.content)
			if result != tt.expected {
				t.Errorf("isTemplateFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}
