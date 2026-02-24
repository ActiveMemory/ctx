//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"strings"
	"testing"
	"time"
)

func TestGetIndentLevel(t *testing.T) {
	tests := []struct {
		line     string
		expected int
	}{
		{"no indent", 0},
		{" one space", 1},
		{"  two spaces", 2},
		{"\ttab", 1},
		{"    four spaces", 4},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := indentLevel(tt.line)
			if got != tt.expected {
				t.Errorf("indentLevel(%q) = %d, want %d", tt.line, got, tt.expected)
			}
		})
	}
}

func TestParseTaskBlocks_SimpleTask(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [ ] Pending task
- [x] Simple completed task
- [ ] Another pending task

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if !block.IsCompleted {
		t.Error("block should be marked completed")
	}
	if !block.IsArchivable {
		t.Error("simple task should be archivable")
	}
	if len(block.Lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(block.Lines))
	}
	if block.ParentTaskText() != "Simple completed task" {
		t.Errorf("unexpected task text: %q", block.ParentTaskText())
	}
}

func TestParseTaskBlocks_TaskWithMetadata(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] T1.2.8 Bug: archive doesn't handle nested content
  When a parent has indented child lines, only the parent
  is archived, leaving orphaned content.
  #priority:medium
- [ ] Another task

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if !block.IsArchivable {
		t.Error("task with metadata should be archivable")
	}
	if len(block.Lines) != 4 {
		t.Errorf("expected 4 lines (parent + 3 metadata), got %d: %v", len(block.Lines), block.Lines)
	}

	// Verify all lines are captured
	content := block.BlockContent()
	if !strings.Contains(content, "T1.2.8") {
		t.Error("missing task title")
	}
	if !strings.Contains(content, "orphaned content") {
		t.Error("missing metadata line")
	}
	if !strings.Contains(content, "#priority:medium") {
		t.Error("missing priority tag")
	}
}

func TestParseTaskBlocks_TaskWithCompletedSubtasks(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] Implement authentication
  - [x] Add login form
  - [x] Add password validation

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if !block.IsArchivable {
		t.Error("task with all completed subtasks should be archivable")
	}
	if len(block.Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(block.Lines))
	}
}

func TestParseTaskBlocks_TaskWithIncompleteSubtask(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] Implement authentication
  - [x] Add login form
  - [ ] Add OAuth support
- [ ] Other task

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if block.IsArchivable {
		t.Error("task with incomplete subtask should NOT be archivable")
	}
	if !block.IsCompleted {
		t.Error("parent is still completed")
	}
}

func TestParseTaskBlocks_MixedMetadataAndSubtasks(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] T1.2.3 refactor ctx watch
  This task involves sharing validation logic.
  - [x] Extract shared functions
  - [x] Update watch to use them
  #priority:medium

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if !block.IsArchivable {
		t.Error("task with metadata and completed subtasks should be archivable")
	}
	if len(block.Lines) != 5 {
		t.Errorf("expected 5 lines, got %d: %v", len(block.Lines), block.Lines)
	}
}

func TestParseTaskBlocks_SkipsCompletedSection(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] New completed task

## Completed

- [x] Already archived task
  With metadata

## Backlog
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block (only from Next Up), got %d", len(blocks))
	}

	if blocks[0].ParentTaskText() != "New completed task" {
		t.Errorf("wrong task found: %q", blocks[0].ParentTaskText())
	}
}

func TestParseTaskBlocks_MultipleCompletedTasks(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] First completed
  With metadata
- [ ] Pending task
- [x] Second completed

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	if blocks[0].ParentTaskText() != "First completed" {
		t.Errorf("first block wrong: %q", blocks[0].ParentTaskText())
	}
	if len(blocks[0].Lines) != 2 {
		t.Errorf("first block should have 2 lines, got %d", len(blocks[0].Lines))
	}

	if blocks[1].ParentTaskText() != "Second completed" {
		t.Errorf("second block wrong: %q", blocks[1].ParentTaskText())
	}
	if len(blocks[1].Lines) != 1 {
		t.Errorf("second block should have 1 line, got %d", len(blocks[1].Lines))
	}
}

func TestParseTaskBlocks_EmptyLinesInBlock(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] Task with blank line in content
  First line of description

  Second paragraph of description
  Still part of the task
- [ ] Next task

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	// Should include: parent, first line, blank, second paragraph (2 lines)
	if len(block.Lines) < 4 {
		t.Errorf("expected at least 4 lines, got %d: %v", len(block.Lines), block.Lines)
	}
	if !block.IsArchivable {
		t.Error("should be archivable")
	}
}

func TestRemoveBlocksFromLines(t *testing.T) {
	lines := []string{
		"# Tasks",
		"",
		"## Next Up",
		"",
		"- [x] First task",
		"  metadata",
		"- [ ] Second task",
		"- [x] Third task",
		"",
		"## Completed",
	}

	blocks := []TaskBlock{
		{StartIndex: 4, EndIndex: 6, Lines: []string{"- [x] First task", "  metadata"}},
		{StartIndex: 7, EndIndex: 8, Lines: []string{"- [x] Third task"}},
	}

	result := RemoveBlocksFromLines(lines, blocks)

	expected := []string{
		"# Tasks",
		"",
		"## Next Up",
		"",
		"- [ ] Second task",
		"",
		"## Completed",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %v", len(expected), len(result), result)
	}

	for i, line := range result {
		if line != expected[i] {
			t.Errorf("line %d: got %q, want %q", i, line, expected[i])
		}
	}
}

func TestBlockContent(t *testing.T) {
	block := TaskBlock{
		Lines: []string{
			"- [x] Parent task",
			"  First child",
			"  Second child",
		},
	}

	expected := "- [x] Parent task\n  First child\n  Second child"
	got := block.BlockContent()

	if got != expected {
		t.Errorf("BlockContent() = %q, want %q", got, expected)
	}
}

func TestParentTaskText(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected string
	}{
		{
			name:     "simple task",
			lines:    []string{"- [x] Simple task"},
			expected: "Simple task",
		},
		{
			name:     "task with metadata",
			lines:    []string{"- [x] Task with stuff #tag", "  metadata"},
			expected: "Task with stuff #tag",
		},
		{
			name:     "empty block",
			lines:    []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := TaskBlock{Lines: tt.lines}
			got := block.ParentTaskText()
			if got != tt.expected {
				t.Errorf("ParentTaskText() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseTaskBlocks_DeeplyNested(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] Top level task
  - [x] Second level
    - [x] Third level
      - [x] Fourth level

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if !block.IsArchivable {
		t.Error("deeply nested complete tasks should be archivable")
	}
	if len(block.Lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(block.Lines))
	}
}

func TestParseTaskBlocks_IncompleteDeepChild(t *testing.T) {
	lines := strings.Split(`# Tasks

## Next Up

- [x] Top level task
  - [x] Second level
    - [ ] Incomplete third level

## Completed
`, "\n")

	blocks := ParseTaskBlocks(lines)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	if block.IsArchivable {
		t.Error("task with incomplete deep child should NOT be archivable")
	}
}

func TestParseDoneTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantNil  bool
		wantYear int
	}{
		{
			name:    "no timestamp",
			line:    "- [x] Simple completed task",
			wantNil: true,
		},
		{
			name:     "with done timestamp",
			line:     "- [x] Task with timestamp #added:2026-01-15-100000 #done:2026-01-20-143000",
			wantNil:  false,
			wantYear: 2026,
		},
		{
			name:     "only done timestamp",
			line:     "- [x] Task #done:2025-12-25-120000",
			wantNil:  false,
			wantYear: 2025,
		},
		{
			name:    "invalid timestamp format",
			line:    "- [x] Task #done:2026-01-20",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDoneTimestamp(tt.line)
			if tt.wantNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Error("expected non-nil time")
				} else if result.Year() != tt.wantYear {
					t.Errorf("year = %d, want %d", result.Year(), tt.wantYear)
				}
			}
		})
	}
}

func TestTaskBlockIsOlderThan(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		doneTime *time.Time
		days     int
		want     bool
	}{
		{
			name:     "nil done time",
			doneTime: nil,
			days:     7,
			want:     false,
		},
		{
			name:     "completed today",
			doneTime: func() *time.Time { t := now; return &t }(),
			days:     7,
			want:     false,
		},
		{
			name:     "completed 10 days ago",
			doneTime: func() *time.Time { t := now.AddDate(0, 0, -10); return &t }(),
			days:     7,
			want:     true,
		},
		{
			name:     "completed exactly 7 days ago",
			doneTime: func() *time.Time { t := now.AddDate(0, 0, -7).Add(-time.Hour); return &t }(),
			days:     7,
			want:     true,
		},
		{
			name:     "completed 5 days ago with 7 day threshold",
			doneTime: func() *time.Time { t := now.AddDate(0, 0, -5); return &t }(),
			days:     7,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := TaskBlock{DoneTime: tt.doneTime}
			got := block.OlderThan(tt.days)
			if got != tt.want {
				t.Errorf("OlderThan(%d) = %v, want %v", tt.days, got, tt.want)
			}
		})
	}
}
