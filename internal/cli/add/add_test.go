//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
)

// TestAddCommand tests the add command.
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
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test adding a task
	addCmd := Cmd()
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
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test adding a decision
	t.Run("add decision", func(t *testing.T) {
		addCmd := Cmd()
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
		addCmd := Cmd()
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
		addCmd := Cmd()
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

// TestPrependOrder tests that decisions and learnings are prepended (newest first).
func TestPrependOrder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-add-prepend-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

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

	t.Run("decisions are prepended", func(t *testing.T) {
		// Add first decision
		addCmd := Cmd()
		addCmd.SetArgs([]string{"decision", "First decision"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add first decision failed: %v", err)
		}

		// Add second decision
		addCmd = Cmd()
		addCmd.SetArgs([]string{"decision", "Second decision"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add second decision failed: %v", err)
		}

		content, err := os.ReadFile(".context/DECISIONS.md")
		if err != nil {
			t.Fatalf("failed to read DECISIONS.md: %v", err)
		}

		contentStr := string(content)
		firstIdx := strings.Index(contentStr, "First decision")
		secondIdx := strings.Index(contentStr, "Second decision")

		if firstIdx == -1 || secondIdx == -1 {
			t.Fatal("decisions not found in file")
		}
		if secondIdx >= firstIdx {
			t.Errorf("second decision should appear before first (prepended), but first at %d, second at %d", firstIdx, secondIdx)
		}
	})

	t.Run("learnings are prepended", func(t *testing.T) {
		// Add first learning
		addCmd := Cmd()
		addCmd.SetArgs([]string{"learning", "First learning"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add first learning failed: %v", err)
		}

		// Add second learning
		addCmd = Cmd()
		addCmd.SetArgs([]string{"learning", "Second learning"})
		if err := addCmd.Execute(); err != nil {
			t.Fatalf("add second learning failed: %v", err)
		}

		content, err := os.ReadFile(".context/LEARNINGS.md")
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}

		contentStr := string(content)
		firstIdx := strings.Index(contentStr, "First learning")
		secondIdx := strings.Index(contentStr, "Second learning")

		if firstIdx == -1 || secondIdx == -1 {
			t.Fatal("learnings not found in file")
		}
		if secondIdx >= firstIdx {
			t.Errorf("second learning should appear before first (prepended), but first at %d, second at %d", firstIdx, secondIdx)
		}
	})
}

// TestAppendEntry tests the AppendEntry function directly.
func TestAppendEntry(t *testing.T) {
	t.Run("decision prepends after header", func(t *testing.T) {
		// Use timestamp format "## [" to match what FormatDecision produces
		existing := []byte("# Decisions\n\n## [2026-01-01] Old Decision\n\nContent\n")
		entry := "## [2026-01-02] New Decision\n\nNew content\n"

		result := AppendEntry(existing, entry, "decision", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New Decision")
		oldIdx := strings.Index(resultStr, "Old Decision")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatalf("decisions not found in result: %s", resultStr)
		}
		if newIdx >= oldIdx {
			t.Errorf("new decision should appear before old, but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("learning prepends after separator", func(t *testing.T) {
		// Use timestamp format "- **[" to match what FormatLearning produces
		existing := []byte("# Learnings\n\n<!-- comment -->\n\n- **[2026-01-01]** Old Learning\n")
		entry := "- **[2026-01-02]** New Learning\n"

		result := AppendEntry(existing, entry, "learning", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New Learning")
		oldIdx := strings.Index(resultStr, "Old Learning")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatalf("learnings not found in result: %s", resultStr)
		}
		if newIdx >= oldIdx {
			t.Errorf("new learning should appear before old, but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("convention appends at end", func(t *testing.T) {
		existing := []byte("# Conventions\n\n- Old convention\n")
		entry := "- New convention\n"

		result := AppendEntry(existing, entry, "convention", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New convention")
		oldIdx := strings.Index(resultStr, "Old convention")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatal("conventions not found in result")
		}
		if newIdx <= oldIdx {
			t.Errorf("new convention should appear after old (appended), but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("decision on empty file", func(t *testing.T) {
		existing := []byte("# Decisions\n\n<!-- Add decisions here -->\n")
		entry := "## [2026-01-01] First Decision\n\nContent\n"

		result := AppendEntry(existing, entry, "decision", "")

		if !strings.Contains(string(result), "First Decision") {
			t.Errorf("decision not found in result: %s", result)
		}
	})

	t.Run("learning on empty file", func(t *testing.T) {
		existing := []byte("# Learnings\n\n<!-- Add gotchas here -->\n")
		entry := "- **[2026-01-01]** First Learning\n"

		result := AppendEntry(existing, entry, "learning", "")

		if !strings.Contains(string(result), "First Learning") {
			t.Errorf("learning not found in result: %s", result)
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
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
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
	addCmd := Cmd()
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
