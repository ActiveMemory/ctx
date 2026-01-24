//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/agent"
	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/cli/complete"
	"github.com/ActiveMemory/ctx/internal/cli/drift"
	"github.com/ActiveMemory/ctx/internal/cli/hook"
	"github.com/ActiveMemory/ctx/internal/cli/init"
	"github.com/ActiveMemory/ctx/internal/cli/load"
	"github.com/ActiveMemory/ctx/internal/cli/loop"
	"github.com/ActiveMemory/ctx/internal/cli/session"
	"github.com/ActiveMemory/ctx/internal/cli/status"
	"github.com/ActiveMemory/ctx/internal/cli/sync"
	"github.com/ActiveMemory/ctx/internal/cli/task"
)

// TestInitCommand tests the init command creates the .context directory
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
	cmd := init.InitCmd()
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

// TestStatusCommand tests the status command
func TestStatusCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-status-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Then status - just verify it runs without error
	// Output goes to stdout, not cmd.Out()
	statusCmd := status.StatusCmd()
	statusCmd.SetArgs([]string{})

	if err := statusCmd.Execute(); err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

// TestAddCommand tests the add command
func TestAddCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test adding a task
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Test task for integration"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task command failed: %v", err)
	}

	// Verify the task was added
	tasksPath := filepath.Join(tmpDir, ".context", "TASKS.md")
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read TASKS.md: %v", err)
	}

	if !strings.Contains(string(content), "Test task for integration") {
		t.Errorf("task was not added to TASKS.md")
	}
}

// TestCompleteCommand tests the complete command
func TestCompleteCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-complete-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add a task
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Task to complete"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task command failed: %v", err)
	}

	// Complete the task
	completeCmd := complete.CompleteCmd()
	completeCmd.SetArgs([]string{"Task to complete"})
	if err := completeCmd.Execute(); err != nil {
		t.Fatalf("complete command failed: %v", err)
	}

	// Verify the task was completed
	tasksPath := filepath.Join(tmpDir, ".context", "TASKS.md")
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read TASKS.md: %v", err)
	}

	if !strings.Contains(string(content), "- [x]") {
		t.Errorf("task was not marked as complete")
	}
}

// TestDriftCommand tests the drift command
func TestDriftCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Run drift - just verify it runs without error
	driftCmd := drift.DriftCmd()
	driftCmd.SetArgs([]string{})

	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift command failed: %v", err)
	}
}

// TestLoadCommand tests the load command
func TestLoadCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-load-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Run load - just verify it runs without error
	loadCmd := load.LoadCmd()
	loadCmd.SetArgs([]string{})

	if err := loadCmd.Execute(); err != nil {
		t.Fatalf("load command failed: %v", err)
	}
}

// TestAgentCommand tests the agent command
func TestAgentCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-agent-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Run agent - just verify it runs without error
	agentCmd := agent.Cmd()
	agentCmd.SetArgs([]string{})

	if err := agentCmd.Execute(); err != nil {
		t.Fatalf("agent command failed: %v", err)
	}
}

// TestHookCommand tests the hook command
func TestHookCommand(t *testing.T) {
	tests := []struct {
		tool     string
		contains string
	}{
		{"claude-code", "Claude Code Integration"},
		{"cursor", "Cursor IDE Integration"},
		{"aider", "Aider Integration"},
		{"copilot", "GitHub Copilot Integration"},
		{"windsurf", "Windsurf Integration"},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			hookCmd := hook.HookCmd()
			hookCmd.SetArgs([]string{tt.tool})

			if err := hookCmd.Execute(); err != nil {
				t.Fatalf("hook %s command failed: %v", tt.tool, err)
			}
		})
	}
}

// TestHookCommandUnknownTool tests hook command with unknown tool
func TestHookCommandUnknownTool(t *testing.T) {
	hookCmd := hook.HookCmd()
	hookCmd.SetArgs([]string{"unknown-tool"})

	err := hookCmd.Execute()
	if err == nil {
		t.Error("hook command should fail for unknown tool")
	}
}

// TestSanitizeFilename tests the sanitizeFilename helper function.
func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple topic", "simple-topic"},
		{"Uppercase Topic", "uppercase-topic"},
		{"topic with   multiple   spaces", "topic-with-multiple-spaces"},
		{"special!@#$%chars", "special-chars"},
		{"already-valid", "already-valid"},
		{"", "session"},
		{"   ", "session"},
		{"---", "session"},
		{"a very long topic name that exceeds the maximum allowed length of fifty characters", "a-very-long-topic-name-that-exceeds-the-maximum-al"},
		{"trailing---", "trailing"},
		{"---leading", "leading"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := session.sanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestTruncate tests the truncate helper function.
func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a longer string", 10, "this is..."},
		{"", 10, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := session.truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestTruncateString tests the truncateString helper function from compact.go.
func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a longer string", 10, "this is..."},
		{"", 10, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := compact.truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestParseIndex tests the parseIndex helper function.
func TestParseIndex(t *testing.T) {
	tests := []struct {
		input     string
		expected  int
		expectErr bool
	}{
		{"1", 1, false},
		{"10", 10, false},
		{"100", 100, false},
		{"0", 0, true},  // index must be positive
		{"-1", 0, true}, // index must be positive
		{"abc", 0, true},
		{"", 0, true},
		{"1.5", 1, false}, // Sscanf stops at the decimal
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := session.parseIndex(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("parseIndex(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("parseIndex(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("parseIndex(%q) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestRemoveEmptySections tests the removeEmptySections helper function.
func TestRemoveEmptySections(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		removed  int
	}{
		{
			name:     "no empty sections",
			input:    "# Title\n\n## Section\n\nContent here\n",
			expected: "# Title\n\n## Section\n\nContent here\n",
			removed:  0,
		},
		{
			name:     "single empty section",
			input:    "# Title\n\n## Empty\n\n## HasContent\n\nSome content\n",
			expected: "# Title\n\n## HasContent\n\nSome content\n",
			removed:  1,
		},
		{
			name:     "multiple empty sections",
			input:    "# Title\n\n## Empty1\n\n## Empty2\n\n## HasContent\n\nContent\n",
			expected: "# Title\n\n## HasContent\n\nContent\n",
			removed:  2,
		},
		{
			name:     "empty section at end",
			input:    "# Title\n\n## Content\n\nText\n\n## EmptyAtEnd\n",
			expected: "# Title\n\n## Content\n\nText\n",
			removed:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, count := compact.removeEmptySections(tt.input)
			if count != tt.removed {
				t.Errorf("removeEmptySections() removed %d sections, want %d", count, tt.removed)
			}
			if result != tt.expected {
				t.Errorf("removeEmptySections() result:\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

// TestSeparateTasks tests the separateTasks helper function.
func TestSeparateTasks(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedCompleted int
		expectedPending   int
	}{
		{
			name:              "mixed tasks",
			input:             "# Tasks\n\n### Phase 1\n- [x] Done task\n- [ ] Pending task\n",
			expectedCompleted: 1,
			expectedPending:   1,
		},
		{
			name:              "all completed",
			input:             "# Tasks\n\n- [x] Task 1\n- [x] Task 2\n",
			expectedCompleted: 2,
			expectedPending:   0,
		},
		{
			name:              "all pending",
			input:             "# Tasks\n\n- [ ] Task 1\n- [ ] Task 2\n",
			expectedCompleted: 0,
			expectedPending:   2,
		},
		{
			name:              "no tasks",
			input:             "# Tasks\n\nNo tasks here.\n",
			expectedCompleted: 0,
			expectedPending:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, stats := task.separateTasks(tt.input)
			if stats.completed != tt.expectedCompleted {
				t.Errorf("separateTasks() completed = %d, want %d", stats.completed, tt.expectedCompleted)
			}
			if stats.pending != tt.expectedPending {
				t.Errorf("separateTasks() pending = %d, want %d", stats.pending, tt.expectedPending)
			}
		})
	}
}

// TestCleanInsight tests the cleanInsight helper function.
func TestCleanInsight(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple text", "simple text"},
		{"  trimmed  ", "trimmed"},
		{"ends with period.", "ends with period"},
		{"ends with comma,", "ends with comma"},
		{"ends with multiple...", "ends with multiple"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := session.cleanInsight(tt.input)
			if result != tt.expected {
				t.Errorf("cleanInsight(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractTextContent tests the extractTextContent helper function.
func TestExtractTextContent(t *testing.T) {
	tests := []struct {
		name     string
		entry    session.transcriptEntry
		expected []string
	}{
		{
			name: "string content",
			entry: session.transcriptEntry{
				Message: session.transcriptMsg{
					Content: "Hello world",
				},
			},
			expected: []string{"Hello world"},
		},
		{
			name: "array content with text",
			entry: session.transcriptEntry{
				Message: session.transcriptMsg{
					Content: []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "First text",
						},
						map[string]interface{}{
							"type": "text",
							"text": "Second text",
						},
					},
				},
			},
			expected: []string{"First text", "Second text"},
		},
		{
			name: "array content with thinking",
			entry: session.transcriptEntry{
				Message: session.transcriptMsg{
					Content: []interface{}{
						map[string]interface{}{
							"type":     "thinking",
							"thinking": "Some thinking",
						},
					},
				},
			},
			expected: []string{"Some thinking"},
		},
		{
			name: "empty content",
			entry: session.transcriptEntry{
				Message: session.transcriptMsg{
					Content: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := session.extractTextContent(tt.entry)
			if len(result) != len(tt.expected) {
				t.Errorf("extractTextContent() returned %d items, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("extractTextContent()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestParseSessionFile tests the parseSessionFile helper function.
func TestParseSessionFile(t *testing.T) {
	// Create a temp session file
	tmpDir, err := os.MkdirTemp("", "session-parse-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sessionContent := `# Session: Test Topic

**Date**: 2025-01-21
**Type**: feature

---

## Summary

This is the summary of the session.

## Tasks
- Task 1
- Task 2
`
	sessionPath := filepath.Join(tmpDir, "test-session.md")
	if err := os.WriteFile(sessionPath, []byte(sessionContent), 0644); err != nil {
		t.Fatalf("failed to write test session: %v", err)
	}

	info, err := session.parseSessionFile(sessionPath)
	if err != nil {
		t.Fatalf("parseSessionFile() error: %v", err)
	}

	if info.Topic != "Test Topic" {
		t.Errorf("parseSessionFile() topic = %q, want %q", info.Topic, "Test Topic")
	}
	if info.Date != "2025-01-21" {
		t.Errorf("parseSessionFile() date = %q, want %q", info.Date, "2025-01-21")
	}
	if info.Type != "feature" {
		t.Errorf("parseSessionFile() type = %q, want %q", info.Type, "feature")
	}
	if info.Summary != "This is the summary of the session." {
		t.Errorf("parseSessionFile() summary = %q, want %q", info.Summary, "This is the summary of the session.")
	}
}

// TestSessionCommands tests the session subcommands.
func TestSessionCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test session save
	t.Run("session save", func(t *testing.T) {
		sessionCmd := session.SessionCmd()
		sessionCmd.SetArgs([]string{"save", "test-topic"})
		if err := sessionCmd.Execute(); err != nil {
			t.Fatalf("session save failed: %v", err)
		}

		// Verify session file was created
		entries, err := os.ReadDir(".context/sessions")
		if err != nil {
			t.Fatalf("failed to read sessions dir: %v", err)
		}
		found := false
		for _, e := range entries {
			if strings.Contains(e.Name(), "test-topic") {
				found = true
				break
			}
		}
		if !found {
			t.Error("session file was not created")
		}
	})

	// Test session list
	t.Run("session list", func(t *testing.T) {
		sessionCmd := session.SessionCmd()
		sessionCmd.SetArgs([]string{"list"})
		if err := sessionCmd.Execute(); err != nil {
			t.Fatalf("session list failed: %v", err)
		}
	})

	// Test session load
	t.Run("session load", func(t *testing.T) {
		sessionCmd := session.SessionCmd()
		sessionCmd.SetArgs([]string{"load", "1"})
		if err := sessionCmd.Execute(); err != nil {
			t.Fatalf("session load failed: %v", err)
		}
	})
}

// TestCompactCommand tests the compact command.
func TestCompactCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-compact-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Run compact with --no-auto-save to skip pre-compact session
	compactCmd := compact.Cmd()
	compactCmd.SetArgs([]string{"--no-auto-save"})
	if err := compactCmd.Execute(); err != nil {
		t.Fatalf("compact failed: %v", err)
	}
}

// TestTasksCommands tests the tasks subcommands.
func TestTasksCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-tasks-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add some tasks
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Test task 1"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task failed: %v", err)
	}

	// Test tasks snapshot
	t.Run("tasks snapshot", func(t *testing.T) {
		tasksCmd := task.TasksCmd()
		tasksCmd.SetArgs([]string{"snapshot", "test-snapshot"})
		if err := tasksCmd.Execute(); err != nil {
			t.Fatalf("tasks snapshot failed: %v", err)
		}

		// Verify snapshot was created
		entries, err := os.ReadDir(".context/archive")
		if err != nil {
			t.Fatalf("failed to read archive dir: %v", err)
		}
		found := false
		for _, e := range entries {
			if strings.Contains(e.Name(), "test-snapshot") {
				found = true
				break
			}
		}
		if !found {
			t.Error("snapshot file was not created")
		}
	})

	// Test tasks archive (dry-run)
	t.Run("tasks archive dry-run", func(t *testing.T) {
		tasksCmd := task.TasksCmd()
		tasksCmd.SetArgs([]string{"archive", "--dry-run"})
		if err := tasksCmd.Execute(); err != nil {
			t.Fatalf("tasks archive failed: %v", err)
		}
	})
}

// TestLoopCommand tests the loop command.
func TestLoopCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-loop-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a PROMPT.md file
	if err := os.WriteFile("PROMPT.md", []byte("# Test Prompt\n"), 0644); err != nil {
		t.Fatalf("failed to create PROMPT.md: %v", err)
	}

	// Test loop command
	loopCmd := loop.LoopCmd()
	loopCmd.SetArgs([]string{"--tool", "generic"})
	if err := loopCmd.Execute(); err != nil {
		t.Fatalf("loop command failed: %v", err)
	}

	// Verify loop.sh was created
	if _, err := os.Stat("loop.sh"); os.IsNotExist(err) {
		t.Error("loop.sh was not created")
	}
}

// TestSyncCommand tests the sync command.
func TestSyncCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-sync-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test sync command
	syncCmd := sync.SyncCmd()
	syncCmd.SetArgs([]string{})
	if err := syncCmd.Execute(); err != nil {
		t.Fatalf("sync command failed: %v", err)
	}
}

// TestAddDecisionAndLearning tests adding decisions and learnings.
func TestAddDecisionAndLearning(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-dl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test adding a decision
	t.Run("add decision", func(t *testing.T) {
		addCmd := add.Cmd()
		addCmd.SetArgs([]string{"decision", "Use PostgreSQL for database"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add decision failed: %v", err)
		}

		content, err := os.ReadFile(".context/DECISIONS.md")
		if err != nil {
			t.Fatalf("failed to read DECISIONS.md: %v", err)
		}
		if !strings.Contains(string(content), "Use PostgreSQL for database") {
			t.Error("decision was not added to DECISIONS.md")
		}
	})

	// Test adding a learning
	t.Run("add learning", func(t *testing.T) {
		addCmd := add.Cmd()
		addCmd.SetArgs([]string{"learning", "Always check for nil before dereferencing"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add learning failed: %v", err)
		}

		content, err := os.ReadFile(".context/LEARNINGS.md")
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}
		if !strings.Contains(string(content), "Always check for nil before dereferencing") {
			t.Error("learning was not added to LEARNINGS.md")
		}
	})

	// Test adding a convention
	t.Run("add convention", func(t *testing.T) {
		addCmd := add.Cmd()
		addCmd.SetArgs([]string{"convention", "Use camelCase for variable names"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add convention failed: %v", err)
		}

		content, err := os.ReadFile(".context/CONVENTIONS.md")
		if err != nil {
			t.Fatalf("failed to read CONVENTIONS.md: %v", err)
		}
		if !strings.Contains(string(content), "Use camelCase for variable names") {
			t.Error("convention was not added to CONVENTIONS.md")
		}
	})
}

// TestAddFromFile tests adding content from a file.
func TestAddFromFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-file-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Create a file with content
	contentFile := filepath.Join(tmpDir, "learning-content.md")
	if err := os.WriteFile(contentFile, []byte("Content from file test"), 0644); err != nil {
		t.Fatalf("failed to create content file: %v", err)
	}

	// Test adding from file
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"learning", "--file", contentFile})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add from file failed: %v", err)
	}

	content, err := os.ReadFile(".context/LEARNINGS.md")
	if err != nil {
		t.Fatalf("failed to read LEARNINGS.md: %v", err)
	}
	if !strings.Contains(string(content), "Content from file test") {
		t.Error("content from file was not added to LEARNINGS.md")
	}
}

// TestAgentJSONOutput tests the agent command with JSON output.
func TestAgentJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-agent-json-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test agent with JSON output
	agentCmd := agent.Cmd()
	agentCmd.SetArgs([]string{"--format", "json"})
	if err := agentCmd.Execute(); err != nil {
		t.Fatalf("agent --format json failed: %v", err)
	}
}

// TestDriftJSONOutput tests the drift command with JSON output.
func TestDriftJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-json-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test drift with JSON output
	driftCmd := drift.DriftCmd()
	driftCmd.SetArgs([]string{"--json"})
	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift --json failed: %v", err)
	}
}

// TestLoadRawOutput tests the load command with raw output.
func TestLoadRawOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-load-raw-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test load with raw output
	loadCmd := load.LoadCmd()
	loadCmd.SetArgs([]string{"--raw"})
	if err := loadCmd.Execute(); err != nil {
		t.Fatalf("load --raw failed: %v", err)
	}
}

// TestStatusJSONOutput tests the status command with JSON output.
func TestStatusJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-status-json-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test status with JSON output
	statusCmd := status.StatusCmd()
	statusCmd.SetArgs([]string{"--json"})
	if err := statusCmd.Execute(); err != nil {
		t.Fatalf("status --json failed: %v", err)
	}
}

// TestSessionParse tests the session parse command with a jsonl file.
func TestSessionParse(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-parse-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a test jsonl file
	jsonlContent := `{"type":"user","message":{"role":"user","content":"Hello"},"timestamp":"2025-01-21T10:00:00Z"}
{"type":"assistant","message":{"role":"assistant","content":"Hi there!"},"timestamp":"2025-01-21T10:00:05Z"}
`
	jsonlPath := filepath.Join(tmpDir, "test-transcript.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	// Test session parse
	sessionCmd := session.SessionCmd()
	sessionCmd.SetArgs([]string{"parse", jsonlPath})
	if err := sessionCmd.Execute(); err != nil {
		t.Fatalf("session parse failed: %v", err)
	}
}

// TestSessionParseWithExtract tests the session parse command with --extract flag.
func TestSessionParseWithExtract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-parse-extract-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a test jsonl file with content that should trigger extraction
	jsonlContent := `{"type":"assistant","message":{"role":"assistant","content":"We decided to use PostgreSQL for the database. I learned that connection pooling is important."},"timestamp":"2025-01-21T10:00:00Z"}
`
	jsonlPath := filepath.Join(tmpDir, "test-transcript.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	// Test session parse with --extract
	sessionCmd := session.SessionCmd()
	sessionCmd.SetArgs([]string{"parse", jsonlPath, "--extract"})
	if err := sessionCmd.Execute(); err != nil {
		t.Fatalf("session parse --extract failed: %v", err)
	}
}

// TestSessionParseWithOutput tests the session parse command with --output flag.
func TestSessionParseWithOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-session-parse-output-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a test jsonl file
	jsonlContent := `{"type":"user","message":{"role":"user","content":"Hello"},"timestamp":"2025-01-21T10:00:00Z"}
`
	jsonlPath := filepath.Join(tmpDir, "test-transcript.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.md")

	// Test session parse with --output
	sessionCmd := session.SessionCmd()
	sessionCmd.SetArgs([]string{"parse", jsonlPath, "--output", outputPath})
	if err := sessionCmd.Execute(); err != nil {
		t.Fatalf("session parse --output failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

// TestFormatTranscriptEntry tests the formatTranscriptEntry helper function.
func TestFormatTranscriptEntry(t *testing.T) {
	tests := []struct {
		name     string
		entry    session.transcriptEntry
		contains []string
	}{
		{
			name: "simple text content",
			entry: session.transcriptEntry{
				Type:      "assistant",
				Timestamp: "2025-01-21T10:00:00Z",
				Message: session.transcriptMsg{
					Role:    "assistant",
					Content: "Hello world",
				},
			},
			contains: []string{"## Assistant", "Hello world"},
		},
		{
			name: "array content with text",
			entry: session.transcriptEntry{
				Type:      "assistant",
				Timestamp: "2025-01-21T10:00:00Z",
				Message: session.transcriptMsg{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "Some text content",
						},
					},
				},
			},
			contains: []string{"## Assistant", "Some text content"},
		},
		{
			name: "array content with thinking",
			entry: session.transcriptEntry{
				Type:      "assistant",
				Timestamp: "2025-01-21T10:00:00Z",
				Message: session.transcriptMsg{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{
							"type":     "thinking",
							"thinking": "Thinking about this...",
						},
					},
				},
			},
			contains: []string{"## Assistant", "Thinking"},
		},
		{
			name: "array content with tool_use",
			entry: session.transcriptEntry{
				Type:      "assistant",
				Timestamp: "2025-01-21T10:00:00Z",
				Message: session.transcriptMsg{
					Role: "assistant",
					Content: []interface{}{
						map[string]interface{}{
							"type": "tool_use",
							"name": "Read",
							"input": map[string]interface{}{
								"file_path": "/some/path",
							},
						},
					},
				},
			},
			contains: []string{"## Assistant", "Tool: Read"},
		},
		{
			name: "array content with tool_result",
			entry: session.transcriptEntry{
				Type:      "assistant",
				Timestamp: "2025-01-21T10:00:00Z",
				Message: session.transcriptMsg{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{
							"type":    "tool_result",
							"content": "Tool result content here",
						},
					},
				},
			},
			contains: []string{"## User", "Tool Result"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := session.formatTranscriptEntry(tt.entry)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("formatTranscriptEntry() result missing %q, got:\n%s", want, result)
				}
			}
		})
	}
}

// TestParseJsonlTranscript tests the parseJsonlTranscript function.
func TestParseJsonlTranscript(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parse-jsonl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test jsonl file
	jsonlContent := `{"type":"user","message":{"role":"user","content":"Hello"},"timestamp":"2025-01-21T10:00:00Z"}
{"type":"assistant","message":{"role":"assistant","content":"Hi!"},"timestamp":"2025-01-21T10:00:05Z"}
{"type":"system","message":{"role":"system","content":"ignored"}}
`
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	result, err := session.parseJsonlTranscript(jsonlPath)
	if err != nil {
		t.Fatalf("parseJsonlTranscript() error: %v", err)
	}

	// Check expected content
	if !strings.Contains(result, "Conversation Transcript") {
		t.Error("result missing header")
	}
	if !strings.Contains(result, "Hello") {
		t.Error("result missing user message")
	}
	if !strings.Contains(result, "Hi!") {
		t.Error("result missing assistant message")
	}
	if !strings.Contains(result, "Total messages: 2") {
		t.Error("result missing message count")
	}
}

// TestExtractInsights tests the extractInsights function.
func TestExtractInsights(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "extract-insights-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test jsonl file with decision and learning content
	jsonlContent := `{"type":"assistant","message":{"role":"assistant","content":"We decided to use Redis for caching because of its speed."},"timestamp":"2025-01-21T10:00:00Z"}
{"type":"assistant","message":{"role":"assistant","content":"I learned that connection pooling is essential for database performance."},"timestamp":"2025-01-21T10:00:05Z"}
{"type":"assistant","message":{"role":"assistant","content":"Gotcha: Always close file handles to avoid resource leaks."},"timestamp":"2025-01-21T10:00:10Z"}
`
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("failed to create jsonl file: %v", err)
	}

	decisions, learnings, err := session.extractInsights(jsonlPath)
	if err != nil {
		t.Fatalf("extractInsights() error: %v", err)
	}

	if len(decisions) == 0 {
		t.Error("extractInsights() found no decisions")
	}
	if len(learnings) == 0 {
		t.Error("extractInsights() found no learnings")
	}
}

// TestCompactWithTasks tests the compact command with actual completed tasks.
func TestCompactWithTasks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-compact-tasks-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// First init
	initCmd := init.InitCmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add and complete a task
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Task to complete"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("add task failed: %v", err)
	}

	completeCmd := complete.CompleteCmd()
	completeCmd.SetArgs([]string{"Task to complete"})
	if err := completeCmd.Execute(); err != nil {
		t.Fatalf("complete task failed: %v", err)
	}

	// Run compact without auto-save
	compactCmd := compact.Cmd()
	compactCmd.SetArgs([]string{"--no-auto-save"})
	if err := compactCmd.Execute(); err != nil {
		t.Fatalf("compact failed: %v", err)
	}
}

// TestInitWithExistingClaudeMdWithCtxMarker tests init when CLAUDE.md already exists with ctx marker.
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
	initCmd := init.InitCmd()
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

// TestFindSessionFile tests the findSessionFile helper with various query types.
func TestFindSessionFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-find-session-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create sessions directory and files
	sessionsDir := filepath.Join(tmpDir, ".context", "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create test session files
	sessions := []string{
		"2025-01-20-100000-first-session.md",
		"2025-01-21-100000-second-session.md",
		"2025-01-22-100000-unique-topic.md",
	}
	for _, name := range sessions {
		path := filepath.Join(sessionsDir, name)
		if err := os.WriteFile(path, []byte("# Test Session\n"), 0644); err != nil {
			t.Fatalf("failed to create session file: %v", err)
		}
	}

	tests := []struct {
		name      string
		query     string
		expectErr bool
	}{
		{"partial match unique", "unique-topic", false},
		{"partial match second", "second-session", false},
		{"index 1", "1", false},
		{"index 2", "2", false},
		{"out of range", "10", true},
		{"no match", "nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := session.findSessionFile(tt.query)
			if tt.expectErr {
				if err == nil {
					t.Errorf("findSessionFile(%q) expected error, got nil", tt.query)
				}
			} else {
				if err != nil {
					t.Errorf("findSessionFile(%q) unexpected error: %v", tt.query, err)
				}
			}
		})
	}
}

// TestBinaryIntegration is an integration test that builds and runs the actual binary
func TestBinaryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir, err := os.MkdirTemp("", "cli-binary-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Build the binary
	binaryPath := filepath.Join(tmpDir, "ctx")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/ctx")
	buildCmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	// Get the project root (go up from internal/cli)
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to get project root: %v", err)
	}
	buildCmd.Dir = projectRoot

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, output)
	}

	// Create a test directory
	testDir := filepath.Join(tmpDir, "test-project")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	// Subtest: ctx init creates expected files
	t.Run("init creates expected files", func(t *testing.T) {
		initCmd := exec.Command(binaryPath, "init")
		initCmd.Dir = testDir
		if output, err := initCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx init failed: %v\n%s", err, output)
		}

		// Check .context directory exists
		ctxDir := filepath.Join(testDir, ".context")
		if _, err := os.Stat(ctxDir); os.IsNotExist(err) {
			t.Fatal(".context directory was not created")
		}

		// Check required files exist
		requiredFiles := []string{
			"CONSTITUTION.md",
			"TASKS.md",
			"DECISIONS.md",
			"LEARNINGS.md",
			"CONVENTIONS.md",
			"ARCHITECTURE.md",
		}
		for _, name := range requiredFiles {
			path := filepath.Join(ctxDir, name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("required file %s was not created", name)
			}
		}
	})

	// Subtest: ctx status returns valid status (not just help text)
	t.Run("status returns valid status", func(t *testing.T) {
		statusCmd := exec.Command(binaryPath, "status")
		statusCmd.Dir = testDir
		output, err := statusCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("ctx status failed: %v\n%s", err, output)
		}

		outputStr := string(output)
		// Verify it's actual status output, not help text
		if strings.Contains(outputStr, "Usage:") || strings.Contains(outputStr, "Available Commands:") {
			t.Error("ctx status returned help text instead of status")
		}
		// Check for expected status output markers
		if !strings.Contains(outputStr, "Context Status") && !strings.Contains(outputStr, "Context Directory") {
			t.Errorf("ctx status did not return expected status output, got:\n%s", outputStr)
		}
	})

	// Subtest: ctx add learning modifies LEARNINGS.md
	t.Run("add learning modifies LEARNINGS.md", func(t *testing.T) {
		addCmd := exec.Command(binaryPath, "add", "learning", "Test learning from integration test")
		addCmd.Dir = testDir
		if output, err := addCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx add learning failed: %v\n%s", err, output)
		}

		// Verify learning was added
		learningsPath := filepath.Join(testDir, ".context", "LEARNINGS.md")
		content, err := os.ReadFile(learningsPath)
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}
		if !strings.Contains(string(content), "Test learning from integration test") {
			t.Error("learning was not added to LEARNINGS.md")
		}
	})

	// Subtest: ctx session save creates session file
	t.Run("session save creates session file", func(t *testing.T) {
		saveCmd := exec.Command(binaryPath, "session", "save")
		saveCmd.Dir = testDir
		if output, err := saveCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx session save failed: %v\n%s", err, output)
		}

		// Check that sessions directory exists and has at least one file
		sessionsDir := filepath.Join(testDir, ".context", "sessions")
		entries, err := os.ReadDir(sessionsDir)
		if err != nil {
			t.Fatalf("failed to read sessions directory: %v", err)
		}
		if len(entries) == 0 {
			t.Error("no session file was created")
		}
	})

	// Subtest: ctx agent returns context packet
	t.Run("agent returns context packet", func(t *testing.T) {
		agentCmd := exec.Command(binaryPath, "agent")
		agentCmd.Dir = testDir
		output, err := agentCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("ctx agent failed: %v\n%s", err, output)
		}

		outputStr := string(output)
		// Verify it's actual agent output, not help text
		if strings.Contains(outputStr, "Usage:") || strings.Contains(outputStr, "Available Commands:") {
			t.Error("ctx agent returned help text instead of context packet")
		}
		// Check for expected context packet markers
		if !strings.Contains(outputStr, "CONSTITUTION") && !strings.Contains(outputStr, "TASKS") {
			t.Errorf("ctx agent did not return expected context packet, got:\n%s", outputStr)
		}
	})

	// Subtest: ctx drift runs without error
	t.Run("drift runs without error", func(t *testing.T) {
		driftCmd := exec.Command(binaryPath, "drift")
		driftCmd.Dir = testDir
		if output, err := driftCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx drift failed: %v\n%s", err, output)
		}
	})

	// Subtest: verify all subcommands execute (not falling through to root help)
	t.Run("subcommands execute without falling through to root help", func(t *testing.T) {
		// Commands that should produce output without "Available Commands:"
		// (which would indicate they fell through to root help)
		subcommands := []struct {
			args     []string
			checkFor string // expected output marker
		}{
			{[]string{"status"}, "Context"},
			{[]string{"agent"}, "Context Packet"},
			{[]string{"drift"}, "Drift"},
			{[]string{"load"}, ""},                 // load outputs context, varies by content
			{[]string{"hook", "cursor"}, "Cursor"}, // hook outputs integration instructions
		}

		for _, tc := range subcommands {
			t.Run(strings.Join(tc.args, "_"), func(t *testing.T) {
				cmd := exec.Command(binaryPath, tc.args...)
				cmd.Dir = testDir
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("ctx %s failed: %v\n%s", strings.Join(tc.args, " "), err, output)
				}

				outputStr := string(output)
				// Critical check: should NOT contain root help indicators
				if strings.Contains(outputStr, "Available Commands:") {
					t.Errorf("ctx %s fell through to root help:\n%s", strings.Join(tc.args, " "), outputStr)
				}
				// If we have an expected marker, check for it
				if tc.checkFor != "" && !strings.Contains(outputStr, tc.checkFor) {
					t.Errorf("ctx %s missing expected output %q:\n%s", strings.Join(tc.args, " "), tc.checkFor, outputStr)
				}
			})
		}
	})
}
