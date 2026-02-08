//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import "testing"

func TestUserInputToEntry(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Task variations
		{"task", EntryTask},
		{"tasks", EntryTask},
		{"Task", EntryTask},
		{"TASKS", EntryTask},

		// Decision variations
		{"decision", EntryDecision},
		{"decisions", EntryDecision},
		{"Decision", EntryDecision},
		{"DECISION", EntryDecision},

		// Learning variations
		{"learning", EntryLearning},
		{"learnings", EntryLearning},
		{"Learning", EntryLearning},
		{"LEARNINGS", EntryLearning},

		// Convention variations
		{"convention", EntryConvention},
		{"conventions", EntryConvention},
		{"Convention", EntryConvention},
		{"CONVENTIONS", EntryConvention},

		// Unknown inputs
		{"", EntryUnknown},
		{"unknown", EntryUnknown},
		{"foo", EntryUnknown},
		{"taskss", EntryUnknown},
		{"learn", EntryUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := UserInputToEntry(tt.input)
			if got != tt.want {
				t.Errorf("UserInputToEntry(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRegExFromAttrName(t *testing.T) {
	tests := []struct {
		name      string
		attrName  string
		input     string
		wantMatch bool
		wantValue string
	}{
		{
			name:      "type attribute",
			attrName:  "type",
			input:     `type="task"`,
			wantMatch: true,
			wantValue: "task",
		},
		{
			name:      "context attribute",
			attrName:  "context",
			input:     `context="some context here"`,
			wantMatch: true,
			wantValue: "some context here",
		},
		{
			name:      "attribute in larger string",
			attrName:  "id",
			input:     `<tag id="123" class="foo">`,
			wantMatch: true,
			wantValue: "123",
		},
		{
			name:      "no match",
			attrName:  "missing",
			input:     `type="task"`,
			wantMatch: false,
			wantValue: "",
		},
		{
			name:      "empty value",
			attrName:  "empty",
			input:     `empty=""`,
			wantMatch: true,
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := RegExFromAttrName(tt.attrName)
			match := re.FindStringSubmatch(tt.input)

			if tt.wantMatch {
				if match == nil {
					t.Errorf("expected match for %q in %q", tt.attrName, tt.input)
					return
				}
				if len(match) < 2 {
					t.Errorf("match has no capture group")
					return
				}
				if match[1] != tt.wantValue {
					t.Errorf("got value %q, want %q", match[1], tt.wantValue)
				}
			} else if match != nil {
				t.Errorf("expected no match for %q in %q, got %v", tt.attrName, tt.input, match)
			}
		})
	}
}

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
			name:      "entry with long title",
			input:     "## [2026-12-31-235959] A much longer title with spaces and stuff",
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
			match := RegExEntryHeader.FindStringSubmatch(tt.input)

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
			match := RegExTask.FindStringSubmatch(tt.input)

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

	matches := RegExTaskMultiline.FindAllStringSubmatch(input, -1)

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
			matched := RegExPhase.MatchString(tt.input)
			if matched != tt.wantMatch {
				t.Errorf("RegExPhase.MatchString(%q) = %v, want %v", tt.input, matched, tt.wantMatch)
			}
		})
	}
}

func TestRegExTaskDoneTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMatch bool
		wantTime  string
	}{
		{
			name:      "task with done timestamp",
			input:     "- [x] Task #done:2026-01-15-143022",
			wantMatch: true,
			wantTime:  "2026-01-15-143022",
		},
		{
			name:      "task with multiple tags",
			input:     "- [x] Task #added:2026-01-01-000000 #done:2026-01-15-143022",
			wantMatch: true,
			wantTime:  "2026-01-15-143022",
		},
		{
			name:      "task without done",
			input:     "- [ ] Task #added:2026-01-01-000000",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := RegExTaskDoneTimestamp.FindStringSubmatch(tt.input)

			if tt.wantMatch {
				if match == nil {
					t.Errorf("expected match for %q", tt.input)
					return
				}
				if match[1] != tt.wantTime {
					t.Errorf("timestamp = %q, want %q", match[1], tt.wantTime)
				}
			} else if match != nil {
				t.Errorf("expected no match for %q", tt.input)
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
			match := RegExPath.FindStringSubmatch(tt.input)

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
	// Verify FileType map contains expected mappings
	expected := map[string]string{
		EntryDecision:   FileDecision,
		EntryTask:       FileTask,
		EntryLearning:   FileLearning,
		EntryConvention: FileConvention,
	}

	for entry, file := range expected {
		if FileType[entry] != file {
			t.Errorf("FileType[%q] = %q, want %q", entry, FileType[entry], file)
		}
	}
}

func TestRequiredFiles(t *testing.T) {
	// Verify FilesRequired contains essential files
	required := map[string]bool{
		FileConstitution: false,
		FileTask:         false,
		FileDecision:     false,
	}

	for _, f := range FilesRequired {
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
	// Verify FileReadOrder has expected files in order
	if len(FileReadOrder) == 0 {
		t.Error("FileReadOrder is empty")
	}

	// Constitution should be first (most important)
	if FileReadOrder[0] != FileConstitution {
		t.Errorf("FileReadOrder[0] = %q, want %q (constitution should be first)",
			FileReadOrder[0], FileConstitution)
	}

	// Tasks should be second (what to work on)
	if FileReadOrder[1] != FileTask {
		t.Errorf("FileReadOrder[1] = %q, want %q (tasks should be second)",
			FileReadOrder[1], FileTask)
	}
}

func TestEntryPlural(t *testing.T) {
	tests := []struct {
		entry string
		want  string
	}{
		{EntryTask, "tasks"},
		{EntryDecision, "decisions"},
		{EntryLearning, "learnings"},
		{EntryConvention, "conventions"},
	}

	for _, tt := range tests {
		t.Run(tt.entry, func(t *testing.T) {
			got := EntryPlural[tt.entry]
			if got != tt.want {
				t.Errorf("EntryPlural[%q] = %q, want %q", tt.entry, got, tt.want)
			}
		})
	}
}

func TestDefaultClaudePermissions(t *testing.T) {
	if len(DefaultClaudePermissions) == 0 {
		t.Error("DefaultClaudePermissions should not be empty")
	}

	// Check that essential ctx commands are included
	expected := []string{
		"Bash(ctx status:*)",
		"Bash(ctx agent:*)",
		"Bash(ctx add:*)",
		"Bash(ctx session:*)",
	}

	permSet := make(map[string]bool)
	for _, p := range DefaultClaudePermissions {
		permSet[p] = true
	}

	for _, e := range expected {
		if !permSet[e] {
			t.Errorf("DefaultClaudePermissions missing: %s", e)
		}
	}
}

func TestConstants(t *testing.T) {
	// Verify important constants are set correctly
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"DirContext", DirContext, ".context"},
		{"DirClaude", DirClaude, ".claude"},
		{"FileTask", FileTask, "TASKS.md"},
		{"FileDecision", FileDecision, "DECISIONS.md"},
		{"FileLearning", FileLearning, "LEARNINGS.md"},
		{"PrefixTaskUndone", PrefixTaskUndone, "- [ ]"},
		{"PrefixTaskDone", PrefixTaskDone, "- [x]"},
		{"IndexStart", IndexStart, "<!-- INDEX:START -->"},
		{"IndexEnd", IndexEnd, "<!-- INDEX:END -->"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}
