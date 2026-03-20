//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/add/cmd/root"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/write/add"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	entrytype "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/entry"
)

// ---------------------------------------------------------------------------
// err.go coverage
// ---------------------------------------------------------------------------

func TestErrNoContent(t *testing.T) {
	err := add.ErrNoContent()
	if err == nil || err.Error() != "no content provided" {
		t.Errorf("ErrNoContent() = %v, want 'no content provided'", err)
	}
}

func TestErrNoContentProvided(t *testing.T) {
	for _, fType := range []string{entrytype.Decision, entrytype.Task, entrytype.Learning, entrytype.Convention, entrytype.Unknown} {
		t.Run(fType, func(t *testing.T) {
			err := add.ErrNoContentProvided(fType, core.ExamplesForType(fType))
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			msg := err.Error()
			if !strings.Contains(msg, "no content provided") {
				t.Errorf("error should contain 'no content provided', got: %s", msg)
			}
			if !strings.Contains(msg, fType) {
				t.Errorf("error should contain type %q, got: %s", fType, msg)
			}
		})
	}
}

func TestErrFileRead(t *testing.T) {
	err := add.ErrFileRead("/some/path", os.ErrNotExist)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "/some/path") {
		t.Errorf("error should contain path, got: %s", err.Error())
	}
}

func TestErrFileWrite(t *testing.T) {
	err := add.ErrFileWriteAdd("/some/path", os.ErrPermission)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "/some/path") {
		t.Errorf("error should contain path, got: %s", err.Error())
	}
}

func TestErrStdinRead(t *testing.T) {
	err := add.ErrStdinRead(os.ErrClosed)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "stdin") {
		t.Errorf("error should mention stdin, got: %s", err.Error())
	}
}

func TestErrIndexUpdate(t *testing.T) {
	err := add.ErrIndexUpdate("/some/file", os.ErrPermission)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "index") {
		t.Errorf("error should mention index, got: %s", err.Error())
	}
}

func TestErrUnknownType(t *testing.T) {
	err := add.ErrUnknownType("foobar")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "foobar") {
		t.Errorf("error should contain the type, got: %s", msg)
	}
	if !strings.Contains(msg, "Valid types") {
		t.Errorf("error should list valid types, got: %s", msg)
	}
}

func TestErrFileNotFound(t *testing.T) {
	err := add.ErrFileNotFound("/missing/file")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "/missing/file") {
		t.Errorf("error should contain path, got: %s", msg)
	}
	if !strings.Contains(msg, "ctx init") {
		t.Errorf("error should suggest 'ctx init', got: %s", msg)
	}
}

func TestErrMissingFields(t *testing.T) {
	err := add.ErrMissingFields("decision", []string{"context", "rationale"})
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "decision") {
		t.Errorf("error should contain entry type, got: %s", msg)
	}
	if !strings.Contains(msg, "context") || !strings.Contains(msg, "rationale") {
		t.Errorf("error should list missing fields, got: %s", msg)
	}
}

// ---------------------------------------------------------------------------
// example.go coverage
// ---------------------------------------------------------------------------

func TestExamplesForType(t *testing.T) {
	tests := []struct {
		fType    string
		contains string
	}{
		{entrytype.Decision, "ctx add decision"},
		{entrytype.Task, "ctx add task"},
		{entrytype.Learning, "ctx add learning"},
		{entrytype.Convention, "ctx add convention"},
		{entrytype.Unknown, "ctx add <type>"},
	}
	for _, tt := range tests {
		t.Run(tt.fType, func(t *testing.T) {
			result := core.ExamplesForType(tt.fType)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("ExamplesForType(%q) should contain %q, got: %s", tt.fType, tt.contains, result)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// fmt.go coverage - FormatTask with priority
// ---------------------------------------------------------------------------

func TestFormatTaskWithPriority(t *testing.T) {
	result := core.FormatTask("My task", "high")
	if !strings.Contains(result, "#priority:high") {
		t.Errorf("FormatTask with priority should contain '#priority:high', got: %s", result)
	}
	if !strings.Contains(result, "My task") {
		t.Errorf("FormatTask should contain task content, got: %s", result)
	}
	if !strings.Contains(result, "#added:") {
		t.Errorf("FormatTask should contain '#added:' timestamp, got: %s", result)
	}
}

func TestFormatTaskWithoutPriority(t *testing.T) {
	result := core.FormatTask("Simple task", "")
	if strings.Contains(result, "#priority:") {
		t.Errorf("FormatTask without priority should not contain '#priority:', got: %s", result)
	}
	if !strings.Contains(result, "Simple task") {
		t.Errorf("FormatTask should contain task content, got: %s", result)
	}
}

// ---------------------------------------------------------------------------
// pos.go coverage
// ---------------------------------------------------------------------------

func TestSkipNewline(t *testing.T) {
	tests := []struct {
		name string
		s    string
		pos  int
		want int
	}{
		{"LF", "abc\ndef", 3, 4},
		{"CRLF", "abc\r\ndef", 3, 5},
		{"no newline", "abcdef", 3, 3},
		{"at end", "abc", 3, 3},
		{"past end", "abc", 5, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.SkipNewline(tt.s, tt.pos)
			if got != tt.want {
				t.Errorf("SkipNewline(%q, %d) = %d, want %d", tt.s, tt.pos, got, tt.want)
			}
		})
	}
}

func TestSkipWhitespace(t *testing.T) {
	tests := []struct {
		name string
		s    string
		pos  int
		want int
	}{
		{"spaces", "   abc", 0, 3},
		{"tabs", "\t\tabc", 0, 2},
		{"newlines", "\n\nabc", 0, 2},
		{"mixed", " \t\n abc", 0, 4},
		{"crlf", "\r\n\r\nabc", 0, 4},
		{"none", "abc", 0, 0},
		{"at end", "abc   ", 3, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.SkipWhitespace(tt.s, tt.pos)
			if got != tt.want {
				t.Errorf("SkipWhitespace(%q, %d) = %d, want %d", tt.s, tt.pos, got, tt.want)
			}
		})
	}
}

func TestFindNewline(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"LF", "abc\ndef", 3},
		{"CRLF", "abc\r\ndef", 3},
		{"none", "abcdef", -1},
		{"empty", "", -1},
		{"starts with LF", "\nabc", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.FindNewline(tt.s)
			if got != tt.want {
				t.Errorf("FindNewline(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// strings.go coverage - ContainsEndComment
// ---------------------------------------------------------------------------

func TestContainsEndComment(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		found, idx := core.ContainsEndComment("some text --> more")
		if !found {
			t.Error("expected to find comment close marker")
		}
		if idx != 10 {
			t.Errorf("expected index 10, got %d", idx)
		}
	})

	t.Run("not found", func(t *testing.T) {
		found, idx := core.ContainsEndComment("no comment close here")
		if found {
			t.Error("should not find comment close marker")
		}
		if idx != -1 {
			t.Errorf("expected index -1, got %d", idx)
		}
	})
}

// ---------------------------------------------------------------------------
// normalize.go coverage - NormalizeTargetSection both branches
// ---------------------------------------------------------------------------

func TestNormalizeTargetSection(t *testing.T) {
	t.Run("without prefix", func(t *testing.T) {
		result := core.NormalizeTargetSection("Phase 1")
		if result != "## Phase 1" {
			t.Errorf("expected '## Phase 1', got %q", result)
		}
	})

	t.Run("with prefix", func(t *testing.T) {
		result := core.NormalizeTargetSection("## Phase 1")
		if result != "## Phase 1" {
			t.Errorf("expected '## Phase 1', got %q", result)
		}
	})
}

// ---------------------------------------------------------------------------
// insert.go coverage - edge cases
// ---------------------------------------------------------------------------

func TestInsertAfterHeader_NoHeader(t *testing.T) {
	content := "Some content without any matching header\n"
	entry := "- New entry\n"

	result := core.InsertAfterHeader(content, entry, "# Missing Header")
	resultStr := string(result)

	if !strings.Contains(resultStr, "New entry") {
		t.Error("entry should be appended when header not found")
	}
}

func TestInsertAfterHeader_HeaderAtEndOfFile(t *testing.T) {
	// Header exists but no newline after it (file ends with header line)
	content := "# Heading"
	entry := "- New entry\n"

	result := core.InsertAfterHeader(content, entry, "# Heading")
	resultStr := string(result)

	if !strings.Contains(resultStr, "New entry") {
		t.Error("entry should be appended when header has no newline after")
	}
}

func TestInsertAfterHeader_WithCtxMarkers(t *testing.T) {
	content := "# Learnings\n" +
		marker.CtxMarkerStart + "\nsome context\n" + marker.CommentClose + "\n\n" +
		"## [2026-01-01] Existing\n"
	entry := "## [2026-01-02] New\n"

	// The header "# Learnings" is found, then markers are skipped
	result := core.InsertAfterHeader(content, entry, desc.TextDesc(text.DescKeyHeadingLearnings))
	resultStr := string(result)

	if !strings.Contains(resultStr, "New") {
		t.Errorf("entry not found in result: %s", resultStr)
	}
}

func TestInsertAfterHeader_CtxMarkerWithoutClose(t *testing.T) {
	// ctx marker start present but no close marker
	content := "# Learnings\n" + marker.CtxMarkerStart + "\nunclosed marker content\nExisting\n"
	entry := "## New entry\n"

	result := core.InsertAfterHeader(content, entry, desc.TextDesc(text.DescKeyHeadingLearnings))
	resultStr := string(result)

	if !strings.Contains(resultStr, "New entry") {
		t.Errorf("entry not found in result: %s", resultStr)
	}
}

func TestAppendAtEnd_WithNewline(t *testing.T) {
	result := core.AppendAtEnd("content\n", "entry\n")
	resultStr := string(result)
	if !strings.Contains(resultStr, "entry") {
		t.Error("entry should be appended")
	}
}

func TestAppendAtEnd_WithoutNewline(t *testing.T) {
	result := core.AppendAtEnd("content", "entry\n")
	resultStr := string(result)
	if !strings.Contains(resultStr, "entry") {
		t.Error("entry should be appended")
	}
	// content should get a newline added before the entry
	if !strings.Contains(resultStr, "content\n") {
		t.Errorf("content should end with newline, got: %q", resultStr)
	}
}

func TestInsertTask_NoPendingNoNewline(t *testing.T) {
	// No unchecked tasks and no trailing newline
	existing := "# Tasks\n\n- [x] Done task"
	entry := "- [ ] New task\n"

	result := core.InsertTask(entry, existing, "")
	resultStr := string(result)

	if !strings.Contains(resultStr, "New task") {
		t.Errorf("new task not found in result: %s", resultStr)
	}
}

func TestInsertTaskAfterSection_SectionNotFound(t *testing.T) {
	content := "# Tasks\n\n- [x] Done\n"
	entry := "- [ ] New task\n"

	result := core.InsertTaskAfterSection(entry, content, "Missing Section")
	resultStr := string(result)

	if !strings.Contains(resultStr, "New task") {
		t.Error("entry should be appended when section not found")
	}
}

func TestInsertTaskAfterSection_SectionAtEnd(t *testing.T) {
	// Section header at end of file without trailing newline after it
	content := "# Tasks\n\n## Phase 1"
	entry := "- [ ] New task\n"

	result := core.InsertTaskAfterSection(entry, content, "Phase 1")
	resultStr := string(result)

	if !strings.Contains(resultStr, "New task") {
		t.Errorf("entry not found in result: %s", resultStr)
	}
}

func TestInsertTaskAfterSection_ContentNoNewline(t *testing.T) {
	// Section not found and no trailing newline
	content := "# Tasks"
	entry := "- [ ] New task\n"

	result := core.InsertTaskAfterSection(entry, content, "Missing")
	resultStr := string(result)

	if !strings.Contains(resultStr, "New task") {
		t.Error("entry should be appended")
	}
}

// ---------------------------------------------------------------------------
// content.go coverage - ExtractContent
// ---------------------------------------------------------------------------

func TestExtractContent_FromFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "content.txt")
	if err := os.WriteFile(tmpFile, []byte("  file content  "), 0600); err != nil {
		t.Fatal(err)
	}

	content, err := core.ExtractContent([]string{"task"}, core.Config{FromFile: tmpFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "file content" {
		t.Errorf("expected 'file content', got %q", content)
	}
}

func TestExtractContent_FromFileMissing(t *testing.T) {
	_, err := core.ExtractContent([]string{"task"}, core.Config{FromFile: "/nonexistent/file"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestExtractContent_FromArgs(t *testing.T) {
	content, err := core.ExtractContent([]string{"task", "hello", "world"}, core.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "hello world" {
		t.Errorf("expected 'hello world', got %q", content)
	}
}

func TestExtractContent_NoContent(t *testing.T) {
	// Only one arg (the type), no file, and stdin is not a pipe in tests
	_, err := core.ExtractContent([]string{"task"}, core.Config{})
	if err == nil {
		t.Fatal("expected error when no content source")
	}
}

// ---------------------------------------------------------------------------
// run.go coverage - ValidateEntry
// ---------------------------------------------------------------------------

func TestValidateEntry(t *testing.T) {
	t.Run("empty content", func(t *testing.T) {
		err := entry.Validate(entry.Params{Type: "task", Content: ""}, nil)
		if err == nil {
			t.Fatal("expected error for empty content")
		}
		if !strings.Contains(err.Error(), "no content provided") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid task", func(t *testing.T) {
		err := entry.Validate(entry.Params{Type: "task", Content: "Do something"}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid convention", func(t *testing.T) {
		err := entry.Validate(entry.Params{Type: "convention", Content: "Use camelCase"}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("decision missing fields", func(t *testing.T) {
		err := entry.Validate(entry.Params{
			Type:    "decision",
			Content: "Some decision",
		}, nil)
		if err == nil {
			t.Fatal("expected error for missing decision fields")
		}
		msg := err.Error()
		if !strings.Contains(msg, "context") {
			t.Errorf("error should mention missing context: %s", msg)
		}
	})

	t.Run("decision valid", func(t *testing.T) {
		err := entry.Validate(entry.Params{
			Type:        "decision",
			Content:     "Use Go",
			Context:     "Need a language",
			Rationale:   "Go is fast",
			Consequence: "Need training",
		}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("learning missing fields", func(t *testing.T) {
		err := entry.Validate(entry.Params{
			Type:    "learning",
			Content: "Some learning",
		}, nil)
		if err == nil {
			t.Fatal("expected error for missing learning fields")
		}
		msg := err.Error()
		if !strings.Contains(msg, "context") {
			t.Errorf("error should mention missing context: %s", msg)
		}
	})

	t.Run("learning valid", func(t *testing.T) {
		err := entry.Validate(entry.Params{
			Type:        "learning",
			Content:     "Go embed",
			Context:     "Tried embedding",
			Lesson:      "Same dir only",
			Application: "Keep files local",
		}, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// run.go coverage - WriteEntry error paths
// ---------------------------------------------------------------------------

func TestWriteEntry_UnknownType(t *testing.T) {
	err := entry.Write(entry.Params{
		Type:    "foobar",
		Content: "something",
	})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !strings.Contains(err.Error(), "foobar") {
		t.Errorf("error should mention the unknown type, got: %v", err)
	}
}

func TestWriteEntry_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// No .context/ directory, so files won't exist
	err := entry.Write(entry.Params{
		Type:    "task",
		Content: "something",
	})
	if err == nil {
		t.Fatal("expected error for missing context file")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention file not found, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// run.go coverage - Run with unknown type
// ---------------------------------------------------------------------------

func TestRun_UnknownType(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := &cobra.Command{}
	addCmd.SetOut(&strings.Builder{})
	addCmd.SetErr(&strings.Builder{})
	err := root.Run(addCmd, []string{"invalidtype", "Some content"}, core.Config{})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

// ---------------------------------------------------------------------------
// run.go coverage - Run with no content (only type arg, no file/stdin)
// ---------------------------------------------------------------------------

func TestRun_NoContent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := &cobra.Command{}
	addCmd.SetOut(&strings.Builder{})
	addCmd.SetErr(&strings.Builder{})
	err := root.Run(addCmd, []string{"task"}, core.Config{})
	if err == nil {
		t.Fatal("expected error when no content provided")
	}
	if !strings.Contains(err.Error(), "no content provided") {
		t.Errorf("expected 'no content provided' error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// run.go coverage - task with priority via Run
// ---------------------------------------------------------------------------

func TestRun_TaskWithPriority(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := &cobra.Command{}
	addCmd.SetOut(&strings.Builder{})
	addCmd.SetErr(&strings.Builder{})
	err := root.Run(addCmd, []string{"task", "High priority task"}, core.Config{Priority: "high"})
	if err != nil {
		t.Fatalf("Run task with priority failed: %v", err)
	}

	content, readErr := os.ReadFile(".context/TASKS.md")
	if readErr != nil {
		t.Fatalf("failed to read TASKS.md: %v", readErr)
	}
	if !strings.Contains(string(content), "#priority:high") {
		t.Error("task with priority should contain '#priority:high'")
	}
}

// ---------------------------------------------------------------------------
// run.go coverage - task with section
// ---------------------------------------------------------------------------

func TestRun_TaskWithSection(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := &cobra.Command{}
	addCmd.SetOut(&strings.Builder{})
	addCmd.SetErr(&strings.Builder{})
	err := root.Run(addCmd, []string{"task", "Sectioned task"}, core.Config{Section: "Next Up"})
	if err != nil {
		t.Fatalf("Run task with section failed: %v", err)
	}

	content, readErr := os.ReadFile(".context/TASKS.md")
	if readErr != nil {
		t.Fatalf("failed to read TASKS.md: %v", readErr)
	}
	if !strings.Contains(string(content), "Sectioned task") {
		t.Error("task should be added to TASKS.md")
	}
}

// ---------------------------------------------------------------------------
// Predicate coverage (already at 100% but ensure plural forms work)
// ---------------------------------------------------------------------------

func TestPredicates(t *testing.T) {
	// Test plural forms
	if !core.FileTypeIsTask("tasks") {
		t.Error("FileTypeIsTask should accept 'tasks'")
	}
	if !core.FileTypeIsDecision("decisions") {
		t.Error("FileTypeIsDecision should accept 'decisions'")
	}
	if !core.FileTypeIsLearning("learnings") {
		t.Error("FileTypeIsLearning should accept 'learnings'")
	}
	// Test negative cases
	if core.FileTypeIsTask("decision") {
		t.Error("FileTypeIsTask should reject 'decision'")
	}
	if core.FileTypeIsDecision("task") {
		t.Error("FileTypeIsDecision should reject 'task'")
	}
	if core.FileTypeIsLearning("convention") {
		t.Error("FileTypeIsLearning should reject 'convention'")
	}
}

// ---------------------------------------------------------------------------
// strings.go coverage - EndsWithNewline edge cases
// ---------------------------------------------------------------------------

func TestEndsWithNewline(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"LF", "content\n", true},
		{"CRLF", "content\r\n", true},
		{"no newline", "content", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.EndsWithNewline(tt.s)
			if got != tt.want {
				t.Errorf("EndsWithNewline(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		found, idx := core.Contains("hello world", "world")
		if !found || idx != 6 {
			t.Errorf("Contains() = (%v, %d), want (true, 6)", found, idx)
		}
	})
	t.Run("not found", func(t *testing.T) {
		found, idx := core.Contains("hello", "world")
		if found || idx != -1 {
			t.Errorf("Contains() = (%v, %d), want (false, -1)", found, idx)
		}
	})
}

func TestContainsNewLine(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		found, idx := core.ContainsNewLine("abc\ndef")
		if !found || idx != 3 {
			t.Errorf("ContainsNewLine() = (%v, %d), want (true, 3)", found, idx)
		}
	})
	t.Run("not found", func(t *testing.T) {
		found, idx := core.ContainsNewLine("abcdef")
		if found || idx != -1 {
			t.Errorf("ContainsNewLine() = (%v, %d), want (false, -1)", found, idx)
		}
	})
}

func TestStartsWithCtxMarker(t *testing.T) {
	if !core.StartsWithCtxMarker(marker.CtxMarkerStart + " rest") {
		t.Error("should detect CtxMarkerStart")
	}
	if !core.StartsWithCtxMarker(marker.CtxMarkerEnd + " rest") {
		t.Error("should detect CtxMarkerEnd")
	}
	if core.StartsWithCtxMarker("no marker here") {
		t.Error("should not detect marker in plain text")
	}
}
