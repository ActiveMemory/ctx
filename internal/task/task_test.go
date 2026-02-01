//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestCompleted(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		want  bool
	}{
		{
			name: "completed task",
			line: "- [x] Do something",
			want: true,
		},
		{
			name: "pending task with space",
			line: "- [ ] Do something",
			want: false,
		},
		{
			name: "pending task empty checkbox",
			line: "- [] Do something",
			want: false,
		},
		{
			name: "indented completed task",
			line: "  - [x] Subtask done",
			want: true,
		},
		{
			name: "indented pending task",
			line: "  - [ ] Subtask pending",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := config.RegExTask.FindStringSubmatch(tt.line)
			if match == nil {
				t.Fatalf("line did not match task pattern: %q", tt.line)
			}
			got := Completed(match)
			if got != tt.want {
				t.Errorf("Completed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompleted_InvalidMatch(t *testing.T) {
	// Test with nil/short match slices
	if Completed(nil) {
		t.Error("Completed(nil) should return false")
	}
	if Completed([]string{}) {
		t.Error("Completed([]) should return false")
	}
	if Completed([]string{"full", "indent"}) {
		t.Error("Completed() with short slice should return false")
	}
}

func TestIsPending(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{
			name: "pending task with space",
			line: "- [ ] Do something",
			want: true,
		},
		{
			name: "pending task empty checkbox",
			line: "- [] Do something",
			want: true,
		},
		{
			name: "completed task",
			line: "- [x] Done task",
			want: false,
		},
		{
			name: "indented pending",
			line: "    - [ ] Deep subtask",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := config.RegExTask.FindStringSubmatch(tt.line)
			if match == nil {
				t.Fatalf("line did not match task pattern: %q", tt.line)
			}
			got := IsPending(match)
			if got != tt.want {
				t.Errorf("IsPending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPending_InvalidMatch(t *testing.T) {
	if IsPending(nil) {
		t.Error("IsPending(nil) should return false")
	}
	if IsPending([]string{"full", "indent"}) {
		t.Error("IsPending() with short slice should return false")
	}
}

func TestIndent(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "no indent",
			line: "- [ ] Top level task",
			want: "",
		},
		{
			name: "two space indent",
			line: "  - [ ] Subtask",
			want: "  ",
		},
		{
			name: "four space indent",
			line: "    - [x] Deep subtask",
			want: "    ",
		},
		{
			name: "tab indent",
			line: "\t- [ ] Tab indented",
			want: "\t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := config.RegExTask.FindStringSubmatch(tt.line)
			if match == nil {
				t.Fatalf("line did not match task pattern: %q", tt.line)
			}
			got := Indent(match)
			if got != tt.want {
				t.Errorf("Indent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIndent_InvalidMatch(t *testing.T) {
	if Indent(nil) != "" {
		t.Error("Indent(nil) should return empty string")
	}
	if Indent([]string{}) != "" {
		t.Error("Indent([]) should return empty string")
	}
}

func TestContent(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "simple task",
			line: "- [ ] Implement feature",
			want: "Implement feature",
		},
		{
			name: "task with tags",
			line: "- [ ] Fix bug #added:2026-01-15-120000",
			want: "Fix bug #added:2026-01-15-120000",
		},
		{
			name: "completed task",
			line: "- [x] Done task #done:2026-01-15-130000",
			want: "Done task #done:2026-01-15-130000",
		},
		{
			name: "task with special characters",
			line: "- [ ] Handle `error` in foo()",
			want: "Handle `error` in foo()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := config.RegExTask.FindStringSubmatch(tt.line)
			if match == nil {
				t.Fatalf("line did not match task pattern: %q", tt.line)
			}
			got := Content(match)
			if got != tt.want {
				t.Errorf("Content() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContent_InvalidMatch(t *testing.T) {
	if Content(nil) != "" {
		t.Error("Content(nil) should return empty string")
	}
	if Content([]string{"full", "indent", "state"}) != "" {
		t.Error("Content() with short slice should return empty string")
	}
}

func TestIsSubTask(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{
			name: "top level task",
			line: "- [ ] Top level",
			want: false,
		},
		{
			name: "single space - not subtask",
			line: " - [ ] One space",
			want: false,
		},
		{
			name: "two space subtask",
			line: "  - [ ] Subtask",
			want: true,
		},
		{
			name: "four space subtask",
			line: "    - [x] Deep subtask",
			want: true,
		},
		{
			name: "tab subtask",
			line: "\t\t- [ ] Tab indented",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := config.RegExTask.FindStringSubmatch(tt.line)
			if match == nil {
				t.Fatalf("line did not match task pattern: %q", tt.line)
			}
			got := IsSubTask(match)
			if got != tt.want {
				t.Errorf("IsSubTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchConstants(t *testing.T) {
	// Verify match indices work correctly
	line := "  - [x] Task content here"
	match := config.RegExTask.FindStringSubmatch(line)
	if match == nil {
		t.Fatal("line did not match task pattern")
	}

	if match[MatchFull] != line {
		t.Errorf("MatchFull = %q, want %q", match[MatchFull], line)
	}
	if match[MatchIndent] != "  " {
		t.Errorf("MatchIndent = %q, want %q", match[MatchIndent], "  ")
	}
	if match[MatchState] != "x" {
		t.Errorf("MatchState = %q, want %q", match[MatchState], "x")
	}
	if match[MatchContent] != "Task content here" {
		t.Errorf("MatchContent = %q, want %q", match[MatchContent], "Task content here")
	}
}
