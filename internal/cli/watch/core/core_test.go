//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/cli/watch/core/apply"
	"github.com/ActiveMemory/ctx/internal/cli/watch/core/stream"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

// TestApplyUpdate tests the Update function routing.
func TestApplyUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-apply-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
		update      apply.ContextUpdate
		checkFile   string
		checkFor    string
		expectError bool
	}{
		{
			name: "task update",
			update: apply.ContextUpdate{
				Type:    entry.Task,
				Content: "Test task from watch",
			},
			checkFile: ctx.Task,
			checkFor:  "Test task from watch",
		},
		{
			name: "decision update",
			update: apply.ContextUpdate{
				Type:        entry.Decision,
				Content:     "Test decision from watch",
				Context:     "Testing watch functionality",
				Rationale:   "Need to verify watch applies decisions",
				Consequence: "Decision will appear in DECISIONS.md",
			},
			checkFile: ctx.Decision,
			checkFor:  "Test decision from watch",
		},
		{
			name: "learning update",
			update: apply.ContextUpdate{
				Type:        entry.Learning,
				Content:     "Test learning from watch",
				Context:     "Testing watch functionality",
				Lesson:      "Watch can add learnings",
				Application: "Use structured attributes in context-update tags",
			},
			checkFile: ctx.Learning,
			checkFor:  "Test learning from watch",
		},
		{
			name: "decision without required fields",
			update: apply.ContextUpdate{
				Type:    entry.Decision,
				Content: "Missing fields",
			},
			expectError: true,
		},
		{
			name: "learning without required fields",
			update: apply.ContextUpdate{
				Type:    entry.Learning,
				Content: "Missing fields",
			},
			expectError: true,
		},
		{
			name: "convention update",
			update: apply.ContextUpdate{
				Type:    entry.Convention,
				Content: "Test convention from watch",
			},
			checkFile: ctx.Convention,
			checkFor:  "Test convention from watch",
		},
		{
			name:        "unknown type",
			update:      apply.ContextUpdate{Type: "invalid", Content: "Should fail"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apply.Update(tt.update)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}

			// Verify content was added
			filePath := filepath.Join(rc.ContextDir(), tt.checkFile)
			content, err := os.ReadFile(filepath.Clean(filePath))
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
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add a task to complete
	tasksPath := filepath.Join(rc.ContextDir(), ctx.Task)
	tasksContent := `# Tasks

## Next Up

- [ ] Implement authentication
- [ ] Write tests
`
	if writeErr := os.WriteFile(
		tasksPath, []byte(tasksContent), 0600,
	); writeErr != nil {
		t.Fatalf("failed to write tasks: %v", writeErr)
	}

	// Complete the task
	update := apply.ContextUpdate{Type: entry.Complete, Content: "authentication"}
	if err = apply.Update(update); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify task was marked complete
	content, err := os.ReadFile(filepath.Clean(tasksPath))
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

// TestProcessStream tests stream processing applies updates.
func TestProcessStream(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watch-stream-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	input := `Some AI output text
<context-update type="task">Stream test task</context-update>
More output
`
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err = stream.Process(cmd, reader, false)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Verify task was written
	tasksPath := filepath.Join(rc.ContextDir(), ctx.Task)
	content, err := os.ReadFile(filepath.Clean(tasksPath))
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
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	input := "Some AI output\n" +
		`<context-update type="learning"` +
		` context="Debugging hooks"` +
		` lesson="Hooks receive JSON via stdin"` +
		` application="Use jq to parse input"` +
		`>Hook Input Format</context-update>` +
		"\nMore output\n"
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var output bytes.Buffer
	cmd.SetOut(&output)

	err = stream.Process(cmd, reader, false)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Verify learning was written with structured fields
	learningsPath := filepath.Join(rc.ContextDir(), ctx.Learning)
	content, err := os.ReadFile(filepath.Clean(learningsPath))
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
		{
			`<context-update type="learning"` +
				` lesson="the lesson"`,
			"lesson", "the lesson",
		},
		{`<context-update type="learning"`, "missing", ""},
		{
			`<context-update type="decision"` +
				` rationale="why we did it"`,
			"rationale", "why we did it",
		},
	}

	for _, tt := range tests {
		result := stream.ExtractAttribute(tt.tag, tt.attr)
		if result != tt.expected {
			t.Errorf(
				"ExtractAttribute(%q, %q) = %q, want %q",
				tt.tag, tt.attr, result, tt.expected,
			)
		}
	}
}

func TestProcessStream_DryRunMode(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	input := `<context-update type="task">Dry run stream task</context-update>
`
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := stream.Process(cmd, reader, true)
	if err != nil {
		t.Fatalf("Process error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Would apply") {
		t.Errorf("dry run should show 'Would apply', got: %q", out)
	}
	if !strings.Contains(out, "Dry run stream task") {
		t.Errorf("dry run should show task content, got: %q", out)
	}
}

func TestProcessStream_FailedApply(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Decision without required fields should fail
	input := `<context-update type="decision">Bad decision</context-update>
`
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := stream.Process(cmd, reader, false)
	if err != nil {
		t.Fatalf("Process should not return error for failed apply: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Failed to apply") {
		t.Error("output should indicate failed apply")
	}
}

func TestProcessStream_MultipleUpdates(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	input := `<context-update type="task">First task</context-update>
<context-update type="task">Second task</context-update>
<context-update type="convention">Use snake_case</context-update>
`
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := stream.Process(cmd, reader, false)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Count(out, "Applied") < 3 {
		t.Errorf("expected 3 applied updates, got: %q", out)
	}
}

func TestProcessStream_DecisionWithAttributes(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	input := `<context-update type="decision"` +
		` context="Need a DB"` +
		` rationale="PostgreSQL is mature"` +
		` consequence="Team needs PG training"` +
		`>Use PostgreSQL</context-update>` + "\n"
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := stream.Process(cmd, reader, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify decision was written
	decPath := filepath.Join(rc.ContextDir(), ctx.Decision)
	content, err := os.ReadFile(filepath.Clean(decPath))
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)
	if !strings.Contains(contentStr, "Use PostgreSQL") {
		t.Error("decision title should be in file")
	}
	if !strings.Contains(contentStr, "Need a DB") {
		t.Error("context attribute should be in file")
	}
}

func TestProcessStream_NoUpdates(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	input := `Just regular text with no updates.
Another line of normal output.
`
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := stream.Process(cmd, reader, false)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "Applied") {
		t.Error("should have no applied updates for plain text")
	}
}

func TestContextUpdate_Fields(t *testing.T) {
	u := apply.ContextUpdate{
		Type:        "learning",
		Content:     "Title",
		Context:     "ctx",
		Lesson:      "lesson",
		Application: "app",
		Rationale:   "rat",
		Consequence: "cons",
	}
	if u.Type != "learning" || u.Content != "Title" {
		t.Error("ContextUpdate fields should be set correctly")
	}
	if u.Context != "ctx" || u.Lesson != "lesson" || u.Application != "app" {
		t.Error("learning fields should be set correctly")
	}
	if u.Rationale != "rat" || u.Consequence != "cons" {
		t.Error("decision fields should be set correctly")
	}
}

func TestExtractAttribute_Consequence(t *testing.T) {
	tag := `<context-update type="decision" consequence="something changes">`
	result := stream.ExtractAttribute(tag, "consequence")
	if result != "something changes" {
		t.Errorf(
			"ExtractAttribute(consequence) = %q,"+
				" want 'something changes'",
			result,
		)
	}
}

func TestExtractAttribute_Application(t *testing.T) {
	tag := `<context-update type="learning" application="use jq">`
	result := stream.ExtractAttribute(tag, "application")
	if result != "use jq" {
		t.Errorf("ExtractAttribute(application) = %q, want 'use jq'", result)
	}
}

func TestProcessStream_CompleteUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Write a task to complete
	tasksPath := filepath.Join(rc.ContextDir(), ctx.Task)
	tasksContent := "# Tasks\n\n- [ ] Implement login\n- [ ] Write tests\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatal(err)
	}

	input := `<context-update type="complete">login</context-update>
`
	reader := strings.NewReader(input)

	cmd := &cobra.Command{Use: "watch"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := stream.Process(cmd, reader, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the task was completed
	content, err := os.ReadFile(filepath.Clean(tasksPath))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "- [x] Implement login") {
		t.Error("login task should be marked complete")
	}
}
