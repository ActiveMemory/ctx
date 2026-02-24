//   /    ctx:                         https://ctx.ist
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
	tasksPath := filepath.Join(rc.ContextDir(), config.FileTask)
	tasksContent := `# Tasks

## Next Up

- [ ] Implement authentication
- [ ] Write tests
`
	if writeErr := os.WriteFile(tasksPath, []byte(tasksContent), 0600); writeErr != nil {
		t.Fatalf("failed to write tasks: %v", writeErr)
	}

	// Complete the task
	update := ContextUpdate{Type: config.EntryComplete, Content: "authentication"}
	if err = applyUpdate(update); err != nil {
		t.Fatalf("applyUpdate failed: %v", err)
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

	// Ensure dry-run is off
	watchDryRun = false

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

	// Ensure dry-run is off
	watchDryRun = false

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

	// Try to complete a non-existent task
	err = runCompleteSilent([]string{"nonexistent task query"})
	if err == nil {
		t.Error("expected error for non-matching task")
	}
	if !strings.Contains(err.Error(), "no task matching") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunCompleteSilent_NoArgs(t *testing.T) {
	err := runCompleteSilent([]string{})
	if err == nil {
		t.Fatal("expected error for empty args")
	}
	if !strings.Contains(err.Error(), "no task specified") {
		t.Errorf("error = %q, want 'no task specified'", err.Error())
	}
}

func TestRunWatch_NoContext(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	// Reset package-level vars
	watchLog = ""
	watchDryRun = false

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := runWatch(cmd, nil)
	if err == nil {
		t.Fatal("expected error when no .context/ exists")
	}
	if !strings.Contains(err.Error(), "ctx init") {
		t.Errorf("error = %q, want 'ctx init' suggestion", err.Error())
	}
}

func TestRunWatch_WithLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		watchLog = ""
		watchDryRun = false
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Create a log file with context-update commands
	logContent := `Some output
<context-update type="task">Task from log file</context-update>
More output
`
	logPath := filepath.Join(tmpDir, "test.log")
	if err := os.WriteFile(logPath, []byte(logContent), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--log", logPath})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("runWatch error: %v", err)
	}

	// Verify task was written
	tasksPath := filepath.Join(rc.ContextDir(), config.FileTask)
	content, err := os.ReadFile(filepath.Clean(tasksPath))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "Task from log file") {
		t.Error("task from log file should be added")
	}
}

func TestRunWatch_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		watchLog = ""
		watchDryRun = false
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Create a log file with updates
	logContent := `<context-update type="task">Dry run task</context-update>
`
	logPath := filepath.Join(tmpDir, "dry.log")
	if err := os.WriteFile(logPath, []byte(logContent), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--log", logPath, "--dry-run"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("runWatch error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DRY RUN") {
		t.Error("output should indicate dry run mode")
	}
	if !strings.Contains(out, "Would apply") {
		t.Error("output should show what would be applied")
	}
}

func TestRunWatch_InvalidLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		watchLog = ""
		watchDryRun = false
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--log", "/nonexistent/path/to/log"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent log file")
	}
	if !strings.Contains(err.Error(), "failed to open log file") {
		t.Errorf("error = %q, want 'failed to open log file'", err.Error())
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
		watchDryRun = false
		watchLog = ""
		rc.Reset()
	})

	rc.Reset()

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Create a log file with the dry-run content (use Execute to properly set flags)
	logContent := `<context-update type="task">Dry run stream task</context-update>
`
	logPath := filepath.Join(tmpDir, "drystream.log")
	if err := os.WriteFile(logPath, []byte(logContent), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--log", logPath, "--dry-run"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("processStream error: %v", err)
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
		watchDryRun = false
	})

	// Initialize context
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	watchDryRun = false

	// Decision without required fields should fail
	input := `<context-update type="decision">Bad decision</context-update>
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := processStream(cmd, reader)
	if err != nil {
		t.Fatalf("processStream should not return error for failed apply: %v", err)
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
		watchDryRun = false
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	watchDryRun = false

	input := `<context-update type="task">First task</context-update>
<context-update type="task">Second task</context-update>
<context-update type="convention">Use snake_case</context-update>
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := processStream(cmd, reader)
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
		watchDryRun = false
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	watchDryRun = false

	input := `<context-update type="decision" context="Need a DB" rationale="PostgreSQL is mature" consequences="Team needs PG training">Use PostgreSQL</context-update>
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := processStream(cmd, reader)
	if err != nil {
		t.Fatal(err)
	}

	// Verify decision was written
	decPath := filepath.Join(rc.ContextDir(), config.FileDecision)
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
		watchDryRun = false
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	watchDryRun = false

	input := `Just regular text with no updates.
Another line of normal output.
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := processStream(cmd, reader)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "Applied") {
		t.Error("should have no applied updates for plain text")
	}
}

func TestContextUpdate_Fields(t *testing.T) {
	u := ContextUpdate{
		Type:         "learning",
		Content:      "Title",
		Context:      "ctx",
		Lesson:       "lesson",
		Application:  "app",
		Rationale:    "rat",
		Consequences: "cons",
	}
	if u.Type != "learning" || u.Content != "Title" {
		t.Error("ContextUpdate fields should be set correctly")
	}
	if u.Context != "ctx" || u.Lesson != "lesson" || u.Application != "app" {
		t.Error("learning fields should be set correctly")
	}
	if u.Rationale != "rat" || u.Consequences != "cons" {
		t.Error("decision fields should be set correctly")
	}
}

func TestCmd_HasFlags(t *testing.T) {
	cmd := Cmd()
	if cmd.Use != "watch" {
		t.Errorf("cmd.Use = %q, want 'watch'", cmd.Use)
	}

	logFlag := cmd.Flags().Lookup("log")
	if logFlag == nil {
		t.Fatal("expected --log flag")
	}

	dryRunFlag := cmd.Flags().Lookup("dry-run")
	if dryRunFlag == nil {
		t.Fatal("expected --dry-run flag")
	}
}

func TestExtractAttribute_Consequences(t *testing.T) {
	tag := `<context-update type="decision" consequences="something changes">`
	result := extractAttribute(tag, "consequences")
	if result != "something changes" {
		t.Errorf("extractAttribute(consequences) = %q, want 'something changes'", result)
	}
}

func TestExtractAttribute_Application(t *testing.T) {
	tag := `<context-update type="learning" application="use jq">`
	result := extractAttribute(tag, "application")
	if result != "use jq" {
		t.Errorf("extractAttribute(application) = %q, want 'use jq'", result)
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
		watchDryRun = false
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Write a task to complete
	tasksPath := filepath.Join(rc.ContextDir(), config.FileTask)
	tasksContent := "# Tasks\n\n- [ ] Implement login\n- [ ] Write tests\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatal(err)
	}

	watchDryRun = false

	input := `<context-update type="complete">login</context-update>
`
	reader := strings.NewReader(input)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := processStream(cmd, reader)
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
