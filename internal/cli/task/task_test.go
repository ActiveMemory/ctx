//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

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
			_, _, stats := separateTasks(tt.input)
			if stats.completed != tt.expectedCompleted {
				t.Errorf("separateTasks() completed = %d, want %d", stats.completed, tt.expectedCompleted)
			}
			if stats.pending != tt.expectedPending {
				t.Errorf("separateTasks() pending = %d, want %d", stats.pending, tt.expectedPending)
			}
		})
	}
}

// TestTasksCommands tests the tasks subcommands.
func TestTasksCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-tasks-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
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
		tasksCmd := Cmd()
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
		tasksCmd := Cmd()
		tasksCmd.SetArgs([]string{"archive", "--dry-run"})
		if err := tasksCmd.Execute(); err != nil {
			t.Fatalf("tasks archive failed: %v", err)
		}
	})
}

// setupTaskDir creates a temp dir with initialized context.
func setupTaskDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	return tmpDir
}

// runTaskCmd executes a task command and captures output.
func runTaskCmd(args ...string) (string, error) {
	cmd := Cmd()
	cmd.SetArgs(args)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

func TestSeparateTasks_WithSubtasks(t *testing.T) {
	content := `# Tasks

### Phase 1
- [x] Completed parent
  - [ ] Subtask of completed (should be archived)
  - [x] Done subtask
- [ ] Pending parent
  - [ ] Subtask of pending (should remain)
`

	remaining, archived, stats := separateTasks(content)

	if stats.completed != 1 {
		t.Errorf("completed = %d, want 1", stats.completed)
	}
	if stats.pending != 1 {
		t.Errorf("pending = %d, want 1", stats.pending)
	}

	// Archived should contain the completed parent and its subtasks
	if !strings.Contains(archived, "Completed parent") {
		t.Error("archived should contain completed parent")
	}
	if !strings.Contains(archived, "Subtask of completed") {
		t.Error("archived should contain subtask of completed parent")
	}

	// Remaining should contain the pending parent and its subtask
	if !strings.Contains(remaining, "Pending parent") {
		t.Error("remaining should contain pending parent")
	}
	if !strings.Contains(remaining, "Subtask of pending") {
		t.Error("remaining should contain subtask of pending parent")
	}
}

func TestSeparateTasks_MultiplePhases(t *testing.T) {
	content := `# Tasks

### Phase 1
- [x] Phase 1 done
- [ ] Phase 1 pending

### Phase 2
- [x] Phase 2 done
- [ ] Phase 2 pending
`

	remaining, archived, stats := separateTasks(content)

	if stats.completed != 2 {
		t.Errorf("completed = %d, want 2", stats.completed)
	}
	if stats.pending != 2 {
		t.Errorf("pending = %d, want 2", stats.pending)
	}

	// Each phase header should appear in archived since both have completed tasks
	if !strings.Contains(archived, "Phase 1") {
		t.Error("archived should contain Phase 1 header")
	}
	if !strings.Contains(archived, "Phase 2") {
		t.Error("archived should contain Phase 2 header")
	}

	// Remaining should still have phase headers and pending tasks
	if !strings.Contains(remaining, "Phase 1 pending") {
		t.Error("remaining should contain Phase 1 pending task")
	}
	if !strings.Contains(remaining, "Phase 2 pending") {
		t.Error("remaining should contain Phase 2 pending task")
	}
}

func TestSeparateTasks_PhaseWithNoCompletedTasks(t *testing.T) {
	content := `# Tasks

### Phase 1
- [ ] Only pending

### Phase 2
- [x] Only completed
`

	_, archived, _ := separateTasks(content)

	// Phase 1 should NOT appear in archived (no completed tasks)
	lines := strings.Split(archived, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Phase 1") {
			t.Error("Phase 1 should not be in archived (no completed tasks)")
		}
	}
	if !strings.Contains(archived, "Phase 2") {
		t.Error("Phase 2 should be in archived")
	}
}

func TestSeparateTasks_NonTaskLines(t *testing.T) {
	content := `# Tasks

Some description text.

- [x] Done
- [ ] Pending

More notes.
`

	remaining, _, _ := separateTasks(content)

	if !strings.Contains(remaining, "Some description text.") {
		t.Error("non-task lines should remain")
	}
	if !strings.Contains(remaining, "More notes.") {
		t.Error("trailing non-task lines should remain")
	}
}

func TestCountPendingTasks(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected int
	}{
		{
			name:     "no tasks",
			lines:    []string{"# Tasks", "Some text"},
			expected: 0,
		},
		{
			name:     "only pending",
			lines:    []string{"- [ ] Task 1", "- [ ] Task 2"},
			expected: 2,
		},
		{
			name:     "mixed",
			lines:    []string{"- [x] Done", "- [ ] Pending"},
			expected: 1,
		},
		{
			name:     "subtasks not counted",
			lines:    []string{"- [ ] Parent", "  - [ ] Subtask"},
			expected: 1,
		},
		{
			name:     "all done",
			lines:    []string{"- [x] Done 1", "- [x] Done 2"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := countPendingTasks(tt.lines)
			if count != tt.expected {
				t.Errorf("countPendingTasks() = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestTasksFilePath(t *testing.T) {
	setupTaskDir(t)

	path := tasksFilePath()
	if !strings.Contains(path, config.FileTask) {
		t.Errorf("tasksFilePath() = %q, want to contain %q", path, config.FileTask)
	}
}

func TestArchiveDirPath(t *testing.T) {
	setupTaskDir(t)

	path := archiveDirPath()
	if !strings.Contains(path, config.DirArchive) {
		t.Errorf("archiveDirPath() = %q, want to contain %q", path, config.DirArchive)
	}
}

func TestSnapshotCommand_NoTasks(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	// Create .context but no TASKS.md
	rc.Reset()
	rc.OverrideContextDir(config.DirContext)
	if err := os.MkdirAll(config.DirContext, 0750); err != nil {
		t.Fatal(err)
	}

	_, err := runTaskCmd("snapshot")
	if err == nil {
		t.Fatal("expected error when TASKS.md doesn't exist")
	}
	if !strings.Contains(err.Error(), "no TASKS.md") {
		t.Errorf("error = %q, want 'no TASKS.md'", err.Error())
	}
}

func TestSnapshotCommand_DefaultName(t *testing.T) {
	setupTaskDir(t)

	// Add a task so TASKS.md has content
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Test task"})
	if err := addCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err := runTaskCmd("snapshot")
	if err != nil {
		t.Fatalf("snapshot error: %v", err)
	}
	if !strings.Contains(out, "Snapshot saved") {
		t.Errorf("output = %q, want 'Snapshot saved'", out)
	}

	// Verify file was created with default name
	entries, err := os.ReadDir(filepath.Join(config.DirContext, config.DirArchive))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range entries {
		if strings.Contains(e.Name(), "snapshot") {
			found = true
		}
	}
	if !found {
		t.Error("snapshot file with default name should be created")
	}
}

func TestArchiveCommand_NoTasks(t *testing.T) {
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
	rc.OverrideContextDir(config.DirContext)
	if err := os.MkdirAll(config.DirContext, 0750); err != nil {
		t.Fatal(err)
	}

	_, err := runTaskCmd("archive")
	if err == nil {
		t.Fatal("expected error when TASKS.md doesn't exist")
	}
}

func TestArchiveCommand_NoCompletedTasks(t *testing.T) {
	setupTaskDir(t)

	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Pending task"})
	if err := addCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err := runTaskCmd("archive")
	if err != nil {
		t.Fatalf("archive error: %v", err)
	}
	if !strings.Contains(out, "No completed tasks") {
		t.Errorf("output = %q, want 'No completed tasks'", out)
	}
}

func TestArchiveCommand_WithCompletedTasks(t *testing.T) {
	setupTaskDir(t)

	// Write TASKS.md with completed and pending tasks
	tasksContent := `# Tasks

## Next Up

- [x] Completed task 1
- [ ] Pending task 1
- [x] Completed task 2
`
	tasksPath := filepath.Join(config.DirContext, config.FileTask)
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runTaskCmd("archive")
	if err != nil {
		t.Fatalf("archive error: %v", err)
	}
	if !strings.Contains(out, "Archived") {
		t.Errorf("output = %q, want 'Archived'", out)
	}
	if !strings.Contains(out, "pending tasks remain") {
		t.Errorf("output should mention pending tasks: %q", out)
	}

	// Verify TASKS.md no longer has completed tasks
	data, err := os.ReadFile(tasksPath) //nolint:gosec // test temp path
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "Completed task 1") {
		t.Error("completed task 1 should be removed from TASKS.md")
	}
	if !strings.Contains(string(data), "Pending task 1") {
		t.Error("pending task 1 should remain in TASKS.md")
	}
}

func TestArchiveCommand_DryRunWithCompleted(t *testing.T) {
	setupTaskDir(t)

	tasksContent := `# Tasks

- [x] Done task
- [ ] Not done task
`
	tasksPath := filepath.Join(config.DirContext, config.FileTask)
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runTaskCmd("archive", "--dry-run")
	if err != nil {
		t.Fatalf("archive --dry-run error: %v", err)
	}
	if !strings.Contains(out, "Dry run") {
		t.Error("output should indicate dry run")
	}
	if !strings.Contains(out, "Would archive") {
		t.Error("output should show what would be archived")
	}

	// Verify TASKS.md was NOT modified
	data, err := os.ReadFile(tasksPath) //nolint:gosec // test temp path
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Done task") {
		t.Error("dry run should not modify TASKS.md")
	}
}

func TestCmd_HasSubcommands(t *testing.T) {
	cmd := Cmd()
	if cmd.Use != "tasks" {
		t.Errorf("cmd.Use = %q, want 'tasks'", cmd.Use)
	}

	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	if !names["archive"] {
		t.Error("missing archive subcommand")
	}
	if !names["snapshot"] {
		t.Error("missing snapshot subcommand")
	}
}

func TestArchiveCommand_DryRunFlag(t *testing.T) {
	cmd := Cmd()
	archiveCmd, _, err := cmd.Find([]string{"archive"})
	if err != nil {
		t.Fatal(err)
	}
	flag := archiveCmd.Flags().Lookup("dry-run")
	if flag == nil {
		t.Fatal("archive should have --dry-run flag")
	}
}

func TestSeparateTasks_EmptyContent(t *testing.T) {
	remaining, archived, stats := separateTasks("")
	if stats.completed != 0 || stats.pending != 0 {
		t.Errorf("stats = %+v, want zero for empty content", stats)
	}
	_ = remaining
	_ = archived
}

func TestSnapshotCommand_SnapshotContentFormat(t *testing.T) {
	setupTaskDir(t)

	tasksContent := "# Tasks\n\n- [ ] My task\n"
	tasksPath := filepath.Join(config.DirContext, config.FileTask)
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := runTaskCmd("snapshot", "my-snap")
	if err != nil {
		t.Fatal(err)
	}

	// Find the snapshot file and verify content
	entries, err := os.ReadDir(filepath.Join(config.DirContext, config.DirArchive))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), "my-snap") {
			data, err := os.ReadFile(filepath.Join(config.DirContext, config.DirArchive, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			content := string(data)
			if !strings.Contains(content, "Snapshot") {
				t.Error("snapshot should have header")
			}
			if !strings.Contains(content, "My task") {
				t.Error("snapshot should contain original tasks")
			}
			if !strings.Contains(content, "---") {
				t.Error("snapshot should contain separator")
			}
			return
		}
	}
	t.Error("snapshot file not found")
}
