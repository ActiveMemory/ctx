//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
)

func TestRegExEntryHeader(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMatch bool
		wantDate  string
		wantTime  string
		wantTitle string
	}{
		{
			name:      "valid entry header",
			input:     "## [2026-01-28-051426] Title here",
			wantMatch: true,
			wantDate:  "2026-01-28",
			wantTime:  "051426",
			wantTitle: "Title here",
		},
		{
			name: "entry with long title",
			input: "## [2026-12-31-235959] " +
				"A much longer title with spaces and stuff",
			wantMatch: true,
			wantDate:  "2026-12-31",
			wantTime:  "235959",
			wantTitle: "A much longer title with spaces and stuff",
		},
		{
			name:      "invalid - missing time",
			input:     "## [2026-01-28] Title",
			wantMatch: false,
		},
		{
			name:      "invalid - wrong format",
			input:     "## Title without timestamp",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := regex.EntryHeader.FindStringSubmatch(tt.input)

			if tt.wantMatch {
				if match == nil {
					t.Errorf("expected match for %q", tt.input)
					return
				}
				if match[1] != tt.wantDate {
					t.Errorf("date = %q, want %q", match[1], tt.wantDate)
				}
				if match[2] != tt.wantTime {
					t.Errorf("time = %q, want %q", match[2], tt.wantTime)
				}
				if match[3] != tt.wantTitle {
					t.Errorf("title = %q, want %q", match[3], tt.wantTitle)
				}
			} else if match != nil {
				t.Errorf("expected no match for %q", tt.input)
			}
		})
	}
}

func TestRegExTask(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantMatch   bool
		wantIndent  string
		wantState   string
		wantContent string
	}{
		{
			name:        "pending task",
			input:       "- [ ] Do something",
			wantMatch:   true,
			wantIndent:  "",
			wantState:   " ",
			wantContent: "Do something",
		},
		{
			name:        "completed task",
			input:       "- [x] Done task",
			wantMatch:   true,
			wantIndent:  "",
			wantState:   "x",
			wantContent: "Done task",
		},
		{
			name:        "indented task",
			input:       "  - [ ] Subtask",
			wantMatch:   true,
			wantIndent:  "  ",
			wantState:   " ",
			wantContent: "Subtask",
		},
		{
			name:        "empty checkbox",
			input:       "- [] Task with empty checkbox",
			wantMatch:   true,
			wantIndent:  "",
			wantState:   "",
			wantContent: "Task with empty checkbox",
		},
		{
			name:        "task with tags",
			input:       "- [ ] Task #added:2026-01-15-120000 #in-progress",
			wantMatch:   true,
			wantIndent:  "",
			wantState:   " ",
			wantContent: "Task #added:2026-01-15-120000 #in-progress",
		},
		{
			name:      "not a task - regular bullet",
			input:     "- Regular bullet point",
			wantMatch: false,
		},
		{
			name:      "not a task - numbered list",
			input:     "1. Numbered item",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := regex.Task.FindStringSubmatch(tt.input)

			if tt.wantMatch {
				if match == nil {
					t.Errorf("expected match for %q", tt.input)
					return
				}
				if match[1] != tt.wantIndent {
					t.Errorf("indent = %q, want %q", match[1], tt.wantIndent)
				}
				if match[2] != tt.wantState {
					t.Errorf("state = %q, want %q", match[2], tt.wantState)
				}
				if match[3] != tt.wantContent {
					t.Errorf("content = %q, want %q", match[3], tt.wantContent)
				}
			} else if match != nil {
				t.Errorf("expected no match for %q", tt.input)
			}
		})
	}
}

func TestRegExTaskMultiline(t *testing.T) {
	input := `# Tasks

## Phase 1

- [x] First task
- [ ] Second task
  - [ ] Subtask A
  - [x] Subtask B

## Phase 2

- [ ] Third task
`

	matches := regex.TaskMultiline.FindAllStringSubmatch(input, -1)

	if len(matches) != 5 {
		t.Errorf("expected 5 matches, got %d", len(matches))
	}

	// Verify first match
	if matches[0][3] != "First task" {
		t.Errorf("first match content = %q, want %q", matches[0][3], "First task")
	}
	if matches[0][2] != "x" {
		t.Errorf("first match state = %q, want %q", matches[0][2], "x")
	}
}

func TestRegExPhase(t *testing.T) {
	tests := []struct {
		input     string
		wantMatch bool
	}{
		{"## Phase 1", true},
		{"### Phase 2: Setup", true},
		{"# Phase", true},
		{"###### Phase 99", true},
		{"Phase 1", false},
		{"##Phase 1", false},
		{"## phase 1", false}, // case sensitive
		{"## Not a phase", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			matched := regex.Phase.MatchString(tt.input)
			if matched != tt.wantMatch {
				t.Errorf(
					"Phase.MatchString(%q) = %v, want %v",
					tt.input, matched, tt.wantMatch,
				)
			}
		})
	}
}

func TestRegExPath(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMatch bool
		wantPath  string
	}{
		{
			name:      "go file",
			input:     "Check `internal/config/config.go` for details",
			wantMatch: true,
			wantPath:  "internal/config/config.go",
		},
		{
			name:      "markdown file",
			input:     "See `docs/README.md`",
			wantMatch: true,
			wantPath:  "docs/README.md",
		},
		{
			name:      "no extension",
			input:     "`Makefile`",
			wantMatch: false,
		},
		{
			name:      "code snippet not path",
			input:     "`fmt.Println`",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := regex.CodeFencePath.FindStringSubmatch(tt.input)

			if tt.wantMatch {
				if match == nil {
					t.Errorf("expected match for %q", tt.input)
					return
				}
				if match[1] != tt.wantPath {
					t.Errorf("path = %q, want %q", match[1], tt.wantPath)
				}
			} else if match != nil {
				t.Errorf("expected no match for %q, got %v", tt.input, match)
			}
		})
	}
}

func TestFileTypeMap(t *testing.T) {
	// Verify CtxFile returns expected mappings
	expected := map[string]string{
		entry.Decision:   ctx.Decision,
		entry.Task:       ctx.Task,
		entry.Learning:   ctx.Learning,
		entry.Convention: ctx.Convention,
	}

	for ent, ctxFile := range expected {
		got, ok := entry.CtxFile(ent)
		if !ok {
			t.Errorf("CtxFile(%q) not found", ent)
		} else if got != ctxFile {
			t.Errorf("CtxFile(%q) = %q, want %q", ent, got, ctxFile)
		}
	}
}

func TestRequiredFiles(t *testing.T) {
	// Verify FilesRequired contains essential files
	required := map[string]bool{
		ctx.Constitution: false,
		ctx.Task:         false,
		ctx.Decision:     false,
	}

	for _, f := range ctx.FilesRequired {
		if _, ok := required[f]; ok {
			required[f] = true
		}
	}

	for f, found := range required {
		if !found {
			t.Errorf("FilesRequired missing %q", f)
		}
	}
}

func TestFileReadOrder(t *testing.T) {
	// Verify ReadOrder has expected files in order
	if len(ctx.ReadOrder) == 0 {
		t.Error("ReadOrder is empty")
	}

	// Constitution should be first (most important)
	if ctx.ReadOrder[0] != ctx.Constitution {
		t.Errorf("ReadOrder[0] = %q, want %q (constitution should be first)",
			ctx.ReadOrder[0], ctx.Constitution)
	}

	// Tasks should be second (what to work on)
	if ctx.ReadOrder[1] != ctx.Task {
		t.Errorf("ReadOrder[1] = %q, want %q (tasks should be second)",
			ctx.ReadOrder[1], ctx.Task)
	}
}

func TestConstants(t *testing.T) {
	// Verify important constants are set correctly
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Context", dir.Context, ".context"},
		{"Claude", dir.Claude, ".claude"},
		{"Task", ctx.Task, "TASKS.md"},
		{"Decision", ctx.Decision, "DECISIONS.md"},
		{"Learning", ctx.Learning, "LEARNINGS.md"},
		{"PrefixTaskUndone", marker.PrefixTaskUndone, "- [ ]"},
		{"PrefixTaskDone", marker.PrefixTaskDone, "- [x]"},
		{"IndexStart", marker.IndexStart, "<!-- INDEX:START -->"},
		{"IndexEnd", marker.IndexEnd, "<!-- INDEX:END -->"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}
