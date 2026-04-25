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
	"github.com/ActiveMemory/ctx/internal/cli/task/core/count"
	"github.com/ActiveMemory/ctx/internal/cli/task/core/path"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

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

	testctx.Declare(t, tmpDir)

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Add some tasks
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Test task 1", "--section", "Misc", "--session-id", "test1234", "--branch", "main", "--commit", "abc123"})
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

	testctx.Declare(t, tmpDir)

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
			count := count.Pending(tt.lines)
			if count != tt.expected {
				t.Errorf("Pending() = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestTasksFilePath(t *testing.T) {
	setupTaskDir(t)

	p, err := path.File()
	if err != nil {
		t.Fatalf("File: %v", err)
	}
	if !strings.Contains(p, ctx.Task) {
		t.Errorf("File() = %q, want to contain %q", p, ctx.Task)
	}
}

func TestArchiveDirPath(t *testing.T) {
	setupTaskDir(t)

	p, err := path.ArchiveDir()
	if err != nil {
		t.Fatalf("ArchiveDir: %v", err)
	}
	if !strings.Contains(p, dir.Archive) {
		t.Errorf("ArchiveDir() = %q, want to contain %q", p, dir.Archive)
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
	})

	// Create .context but no TASKS.md
	ctxDir := testctx.Declare(t, tmpDir)
	if err := os.MkdirAll(ctxDir, 0750); err != nil {
		t.Fatal(err)
	}

	_, err := runTaskCmd("snapshot")
	if err == nil {
		t.Fatal("expected error when TASKS.md doesn't exist")
	}
	if !strings.Contains(err.Error(), "TASKS.md not found") {
		t.Errorf("error = %q, want 'TASKS.md not found'", err.Error())
	}
}

func TestSnapshotCommand_DefaultName(t *testing.T) {
	setupTaskDir(t)

	// Add a task so TASKS.md has content
	addCmd := add.Cmd()
	addCmd.SetArgs([]string{"task", "Test task", "--section", "Misc", "--session-id", "test1234", "--branch", "main", "--commit", "abc123"})
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
	entries, err := os.ReadDir(filepath.Join(dir.Context, dir.Archive))
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
	})

	ctxDir := testctx.Declare(t, tmpDir)
	if err := os.MkdirAll(ctxDir, 0750); err != nil {
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
	addCmd.SetArgs([]string{"task", "Pending task", "--section", "Misc", "--session-id", "test1234", "--branch", "main", "--commit", "abc123"})
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
	tasksPath := filepath.Join(dir.Context, ctx.Task)
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
	tasksPath := filepath.Join(dir.Context, ctx.Task)
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
	if cmd.Use != "task" {
		t.Errorf("cmd.Use = %q, want 'task'", cmd.Use)
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

func TestSnapshotCommand_SnapshotContentFormat(t *testing.T) {
	setupTaskDir(t)

	tasksContent := "# Tasks\n\n- [ ] My task\n"
	tasksPath := filepath.Join(dir.Context, ctx.Task)
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := runTaskCmd("snapshot", "my-snap")
	if err != nil {
		t.Fatal(err)
	}

	// Find the snapshot file and verify content
	entries, err := os.ReadDir(filepath.Join(dir.Context, dir.Archive))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), "my-snap") {
			data, err := os.ReadFile(filepath.Join(dir.Context, dir.Archive, e.Name()))
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
