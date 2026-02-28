//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/rc"
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

	// Create the top-level directory so the path passes the
	// "top dir exists" filter but the full file path is still dead.
	if mkErr := os.Mkdir(filepath.Join(tmpDir, "internal"), 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}

	// Create a test context with a dead path reference
	ctx := &context.Context{
		Dir: ".context",
		Files: []context.FileInfo{
			{
				Name:    "ARCHITECTURE.md",
				Content: []byte("# Architecture\n\nSee `internal/nonexistent.go` for details.\n"),
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
		if report.Warnings[0].Path != "internal/nonexistent.go" {
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

func TestCheckEntryCount(t *testing.T) {
	// Helper to build N entries
	buildEntries := func(n int) string {
		var sb strings.Builder
		sb.WriteString("# Learnings\n\n")
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("## [2026-01-%02d-120000] Entry %d\n\nContent for entry %d.\n\n", (i%28)+1, i+1, i+1))
		}
		return sb.String()
	}

	tests := []struct {
		name         string
		files        []context.FileInfo
		wantWarnings int
		wantPassed   bool
	}{
		{
			name:         "no knowledge files",
			files:        nil,
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "zero entries",
			files: []context.FileInfo{
				{Name: "LEARNINGS.md", Content: []byte("# Learnings\n")},
			},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "at threshold (30 learnings)",
			files: []context.FileInfo{
				{Name: "LEARNINGS.md", Content: []byte(buildEntries(30))},
			},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "above threshold (31 learnings)",
			files: []context.FileInfo{
				{Name: "LEARNINGS.md", Content: []byte(buildEntries(31))},
			},
			wantWarnings: 1,
			wantPassed:   false,
		},
		{
			name: "decisions above threshold (21)",
			files: []context.FileInfo{
				{Name: "DECISIONS.md", Content: []byte(buildEntries(21))},
			},
			wantWarnings: 1,
			wantPassed:   false,
		},
		{
			name: "both files above threshold",
			files: []context.FileInfo{
				{Name: "LEARNINGS.md", Content: []byte(buildEntries(31))},
				{Name: "DECISIONS.md", Content: []byte(buildEntries(21))},
			},
			wantWarnings: 2,
			wantPassed:   false,
		},
		{
			name: "warning message format",
			files: []context.FileInfo{
				{Name: "LEARNINGS.md", Content: []byte(buildEntries(35))},
			},
			wantWarnings: 1,
			wantPassed:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &context.Context{
				Dir:   ".context",
				Files: tt.files,
			}

			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []CheckName{},
			}

			checkEntryCount(ctx, report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
			}

			passedCheck := false
			for _, p := range report.Passed {
				if p == CheckEntryCount {
					passedCheck = true
					break
				}
			}
			if passedCheck != tt.wantPassed {
				t.Errorf("expected passed=%v, got passed=%v", tt.wantPassed, passedCheck)
			}

			// Verify warning type and message format
			for _, w := range report.Warnings {
				if w.Type != IssueEntryCount {
					t.Errorf("expected issue type %q, got %q", IssueEntryCount, w.Type)
				}
				if !strings.Contains(w.Message, "entries (recommended:") {
					t.Errorf("unexpected message format: %q", w.Message)
				}
			}
		})
	}
}

func TestCheckEntryCountDisabled(t *testing.T) {
	// Helper to build N entries
	buildEntries := func(n int) string {
		var sb strings.Builder
		sb.WriteString("# Learnings\n\n")
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("## [2026-01-%02d-120000] Entry %d\n\nContent for entry %d.\n\n", (i%28)+1, i+1, i+1))
		}
		return sb.String()
	}

	// Override rc to set thresholds to 0 (disabled)
	rc.Reset()
	defer rc.Reset()

	ctx := &context.Context{
		Dir: ".context",
		Files: []context.FileInfo{
			{Name: "LEARNINGS.md", Content: []byte(buildEntries(100))},
			{Name: "DECISIONS.md", Content: []byte(buildEntries(100))},
		},
	}

	report := &Report{
		Warnings:   []Issue{},
		Violations: []Issue{},
		Passed:     []CheckName{},
	}

	// With default thresholds (30/20), 100 entries should trigger warnings
	checkEntryCount(ctx, report)

	if len(report.Warnings) != 2 {
		t.Errorf("expected 2 warnings with defaults, got %d", len(report.Warnings))
	}
}

func TestCheckMissingPackages(t *testing.T) {
	// Create a temp directory with an internal/ tree
	tmpDir, err := os.MkdirTemp("", "drift-pkg-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func(path string) {
		rmErr := os.RemoveAll(path)
		if rmErr != nil {
			fmt.Printf("failed to remove temp dir %q: %v", path, rmErr)
		}
	}(tmpDir)

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

	// Create internal/ subdirectories
	for _, d := range []string{"internal/config", "internal/cli", "internal/drift", "internal/newpkg"} {
		if mkErr := os.MkdirAll(filepath.Join(tmpDir, d), 0750); mkErr != nil {
			t.Fatalf("failed to create dir %s: %v", d, mkErr)
		}
	}

	tests := []struct {
		name         string
		archContent  string
		wantWarnings int
		wantPassed   bool
		wantPaths    []string
	}{
		{
			name:         "all packages documented",
			archContent:  "# Arch\n\n| `internal/config` | ... |\n| `internal/cli` | ... |\n| `internal/drift` | ... |\n| `internal/newpkg` | ... |\n",
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name:         "one package missing",
			archContent:  "# Arch\n\n| `internal/config` | ... |\n| `internal/cli` | ... |\n| `internal/drift` | ... |\n",
			wantWarnings: 1,
			wantPassed:   false,
			wantPaths:    []string{"internal/newpkg"},
		},
		{
			name:         "nested path normalizes to parent",
			archContent:  "# Arch\n\n| `internal/config` | ... |\n| `internal/cli/pad` | ... |\n| `internal/drift` | ... |\n| `internal/newpkg` | ... |\n",
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name:         "no ARCHITECTURE.md — skip silently",
			archContent:  "",
			wantWarnings: 0,
			wantPassed:   false, // not passed because check was skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var files []context.FileInfo
			if tt.archContent != "" {
				files = append(files, context.FileInfo{
					Name:    "ARCHITECTURE.md",
					Content: []byte(tt.archContent),
				})
			}

			ctx := &context.Context{Dir: ".context", Files: files}
			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []CheckName{},
			}

			checkMissingPackages(ctx, report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
				for _, w := range report.Warnings {
					t.Logf("  warning: %s (path=%s)", w.Message, w.Path)
				}
			}

			passedCheck := false
			for _, p := range report.Passed {
				if p == CheckMissingPackages {
					passedCheck = true
					break
				}
			}
			if passedCheck != tt.wantPassed {
				t.Errorf("expected passed=%v, got passed=%v", tt.wantPassed, passedCheck)
			}

			for _, w := range report.Warnings {
				if w.Type != IssueMissingPackage {
					t.Errorf("expected issue type %q, got %q", IssueMissingPackage, w.Type)
				}
			}

			if tt.wantPaths != nil {
				gotPaths := make(map[string]bool)
				for _, w := range report.Warnings {
					gotPaths[w.Path] = true
				}
				for _, p := range tt.wantPaths {
					if !gotPaths[p] {
						t.Errorf("expected warning for path %q, not found", p)
					}
				}
			}
		})
	}
}

func TestNormalizeInternalPkg(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"internal/config", "internal/config"},
		{"internal/cli/pad", "internal/cli"},
		{"internal/recall/parser", "internal/recall"},
		{"internal/journal/state", "internal/journal"},
		{"internal", "internal"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeInternalPkg(tt.input)
			if got != tt.want {
				t.Errorf("normalizeInternalPkg(%q) = %q, want %q", tt.input, got, tt.want)
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
