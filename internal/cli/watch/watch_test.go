//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// TestApplyUpdate tests the applyUpdate function routing.
func TestApplyUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-apply-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	tests := []struct {
		name        string
		update      ContextUpdate
		checkFile   string
		checkFor    string
		expectError bool
	}{
		{
			name:      "task update",
			update:    ContextUpdate{Type: config.EntryTask, Content: "Test task from watch"},
			checkFile: config.FileTask,
			checkFor:  "Test task from watch",
		},
		{
			name: "decision update",
			update: ContextUpdate{
				Type:         config.EntryDecision,
				Content:      "Test decision from watch",
				Context:      "Testing watch functionality",
				Rationale:    "Need to verify watch applies decisions",
				Consequences: "Decision will appear in DECISIONS.md",
			},
			checkFile: config.FileDecision,
			checkFor:  "Test decision from watch",
		},
		{
			name: "learning update",
			update: ContextUpdate{
				Type:        config.EntryLearning,
				Content:     "Test learning from watch",
				Context:     "Testing watch functionality",
				Lesson:      "Watch can add learnings",
				Application: "Use structured attributes in context-update tags",
			},
			checkFile: config.FileLearning,
			checkFor:  "Test learning from watch",
		},
		{
			name:        "decision without required fields",
			update:      ContextUpdate{Type: config.EntryDecision, Content: "Missing fields"},
			expectError: true,
		},
		{
			name:        "learning without required fields",
			update:      ContextUpdate{Type: config.EntryLearning, Content: "Missing fields"},
			expectError: true,
		},
		{
			name:      "convention update",
			update:    ContextUpdate{Type: config.EntryConvention, Content: "Test convention from watch"},
			checkFile: config.FileConvention,
			checkFor:  "Test convention from watch",
		},
		{
			name:        "unknown type",
			update:      ContextUpdate{Type: "invalid", Content: "Should fail"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := applyUpdate(tt.update)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("applyUpdate failed: %v", err)
			}

			// Verify content was added
			filePath := filepath.Join(rc.ContextDir(), tt.checkFile)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tt.checkFile, err)
			}
			if !strings.Contains(string(content), tt.checkFor) {
				t.Errorf("expected %s to contain %q", tt.checkFile, tt.checkFor)
			}
		})
	}
}

// TestApplyCompleteUpdate tests the complete update type.
func TestApplyCompleteUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-complete-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add a task to complete
	tasksPath := filepath.Join(rc.ContextDir(), config.FileTask)
	tasksContent := `# Tasks

## Next Up

- [ ] Implement authentication
- [ ] Write tests
`
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0644); err != nil {
		t.Fatalf("failed to write tasks: %v", err)
	}

	// Complete the task
	update := ContextUpdate{Type: config.EntryComplete, Content: "authentication"}
	if err := applyUpdate(update); err != nil {
		t.Fatalf("applyUpdate failed: %v", err)
	}

	// Verify task was marked complete
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read tasks: %v", err)
	}
	if !strings.Contains(string(content), "- [x] Implement authentication") {
		t.Error("task was not marked complete")
	}
	if !strings.Contains(string(content), "- [ ] Write tests") {
		t.Error("other task should remain unchecked")
	}
}

// TestBuildWatchSession tests session snapshot generation.
func TestBuildWatchSession(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 30, 0, 0, time.UTC)
	updates := []ContextUpdate{
		{Type: config.EntryTask, Content: "Task 1"},
		{Type: config.EntryTask, Content: "Task 2"},
		{Type: config.EntryDecision, Content: "Decision 1"},
		{Type: config.EntryLearning, Content: "Learning 1"},
	}

	result := buildWatchSession(timestamp, updates)

	// Check metadata
	if !strings.Contains(result, "**Date**: 2026-01-15") {
		t.Error("missing date in session")
	}
	if !strings.Contains(result, "**Time**: 14:30:00") {
		t.Error("missing time in session")
	}
	if !strings.Contains(result, "watch-auto-save") {
		t.Error("missing session type")
	}

	// Check updates by type
	if !strings.Contains(result, "### Tasks") {
		t.Error("missing Tasks section")
	}
	if !strings.Contains(result, "- Task 1") {
		t.Error("missing Task 1")
	}
	if !strings.Contains(result, "- Task 2") {
		t.Error("missing Task 2")
	}
	if !strings.Contains(result, "### Decisions") {
		t.Error("missing Decisions section")
	}
	if !strings.Contains(result, "### Learnings") {
		t.Error("missing Learnings section")
	}
}

// TestProcessStream tests stream processing applies updates.
func TestProcessStream(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-stream-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Ensure dry-run is off
	watchDryRun = false
	watchAutoSave = false

	input := `Some AI output text
<context-update type="task">Stream test task</context-update>
More output
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var output bytes.Buffer
	cmd.SetOut(&output)

	err = processStream(cmd, reader)
	if err != nil {
		t.Fatalf("processStream failed: %v", err)
	}

	// Verify task was written
	tasksPath := filepath.Join(rc.ContextDir(), config.FileTask)
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read tasks: %v", err)
	}
	if !strings.Contains(string(content), "Stream test task") {
		t.Error("task should have been added to file")
	}
}

// TestProcessStreamWithAttributes tests parsing of structured attributes.
func TestProcessStreamWithAttributes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-attr-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Ensure dry-run is off
	watchDryRun = false
	watchAutoSave = false

	input := `Some AI output
<context-update type="learning" context="Debugging hooks" lesson="Hooks receive JSON via stdin" application="Use jq to parse input">Hook Input Format</context-update>
More output
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var output bytes.Buffer
	cmd.SetOut(&output)

	err = processStream(cmd, reader)
	if err != nil {
		t.Fatalf("processStream failed: %v", err)
	}

	// Verify learning was written with structured fields
	learningsPath := filepath.Join(rc.ContextDir(), config.FileLearning)
	content, err := os.ReadFile(learningsPath)
	if err != nil {
		t.Fatalf("failed to read learnings: %v", err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "Hook Input Format") {
		t.Error("learning title should be in file")
	}
	if !strings.Contains(contentStr, "Debugging hooks") {
		t.Error("context attribute should be in file")
	}
	if !strings.Contains(contentStr, "Hooks receive JSON via stdin") {
		t.Error("lesson attribute should be in file")
	}
	if !strings.Contains(contentStr, "Use jq to parse input") {
		t.Error("application attribute should be in file")
	}
	// Should NOT contain placeholders since attributes were provided
	if strings.Contains(contentStr, "[Context from watch") {
		t.Error("should not have placeholder when context attribute provided")
	}
}

// TestExtractAttribute tests the attribute extraction helper.
func TestExtractAttribute(t *testing.T) {
	tests := []struct {
		tag      string
		attr     string
		expected string
	}{
		{`<context-update type="learning"`, "type", "learning"},
		{`<context-update type="decision" context="test ctx"`, "context", "test ctx"},
		{`<context-update type="learning" lesson="the lesson"`, "lesson", "the lesson"},
		{`<context-update type="learning"`, "missing", ""},
		{`<context-update type="decision" rationale="why we did it"`, "rationale", "why we did it"},
	}

	for _, tt := range tests {
		result := extractAttribute(tt.tag, tt.attr)
		if result != tt.expected {
			t.Errorf("extractAttribute(%q, %q) = %q, want %q", tt.tag, tt.attr, result, tt.expected)
		}
	}
}

// TestRunCompleteSilentNoMatch tests complete with no matching task.
func TestRunCompleteSilentNoMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-nomatch-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Try to complete a non-existent task
	err = runCompleteSilent([]string{"nonexistent task query"})
	if err == nil {
		t.Error("expected error for non-matching task")
	}
	if !strings.Contains(err.Error(), "no task matching") {
		t.Errorf("unexpected error: %v", err)
	}
}
