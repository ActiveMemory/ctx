//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_parseRef(t *testing.T) {
	tests := []struct {
		input      string
		wantType   string
		wantNumber int
		wantText   string
	}{
		{"decision:12", "decision", 12, ""},
		{"learning:5", "learning", 5, ""},
		{"task:8", "task", 8, ""},
		{"convention:3", "convention", 3, ""},
		{"session:abc123", "session", 0, "abc123"},
		{`"Hotfix note"`, "note", 0, "Hotfix note"},
		{"unknown", "note", 0, "unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			gotType, gotNumber, gotText := parseRef(tc.input)
			if gotType != tc.wantType {
				t.Errorf("parseRef(%q) type = %q, want %q", tc.input, gotType, tc.wantType)
			}
			if gotNumber != tc.wantNumber {
				t.Errorf("parseRef(%q) number = %d, want %d", tc.input, gotNumber, tc.wantNumber)
			}
			if gotText != tc.wantText {
				t.Errorf("parseRef(%q) text = %q, want %q", tc.input, gotText, tc.wantText)
			}
		})
	}
}

func TestResolveDecision(t *testing.T) {
	contextDir := t.TempDir()

	decisions := `# Decisions

## [2026-01-10-120000] Use PostgreSQL for storage

Some rationale here.

## [2026-02-15-140000] Adopt hexagonal architecture

Another rationale.
`
	if err := os.WriteFile(filepath.Join(contextDir, "DECISIONS.md"), []byte(decisions), 0o600); err != nil {
		t.Fatalf("WriteFile(DECISIONS.md) error: %v", err)
	}

	resolved := Resolve("decision:1", contextDir)

	if !resolved.Found {
		t.Fatalf("Resolve(decision:1) Found = false, want true")
	}
	if resolved.Title != "Use PostgreSQL for storage" {
		t.Errorf("Resolve(decision:1) Title = %q, want %q", resolved.Title, "Use PostgreSQL for storage")
	}
	if resolved.Raw != "decision:1" {
		t.Errorf("Resolve(decision:1) Raw = %q, want %q", resolved.Raw, "decision:1")
	}
	if resolved.Type != "decision" {
		t.Errorf("Resolve(decision:1) Type = %q, want %q", resolved.Type, "decision")
	}
}

func TestResolveTask(t *testing.T) {
	contextDir := t.TempDir()

	tasks := `# Tasks

- [ ] First pending task
- [x] Completed task
- [ ] Second pending task
`
	if err := os.WriteFile(filepath.Join(contextDir, "TASKS.md"), []byte(tasks), 0o600); err != nil {
		t.Fatalf("WriteFile(TASKS.md) error: %v", err)
	}

	resolved := Resolve("task:1", contextDir)

	if !resolved.Found {
		t.Fatalf("Resolve(task:1) Found = false, want true")
	}
	if resolved.Title != "First pending task" {
		t.Errorf("Resolve(task:1) Title = %q, want %q", resolved.Title, "First pending task")
	}
}

func TestResolveTaskCompleted(t *testing.T) {
	contextDir := t.TempDir()

	tasks := `# Tasks

- [ ] First pending task
- [x] Completed task
- [ ] Second pending task
`
	if err := os.WriteFile(filepath.Join(contextDir, "TASKS.md"), []byte(tasks), 0o600); err != nil {
		t.Fatalf("WriteFile(TASKS.md) error: %v", err)
	}

	// task:2 should be the completed task (second top-level task overall)
	resolved := Resolve("task:2", contextDir)

	if !resolved.Found {
		t.Fatalf("Resolve(task:2) Found = false, want true")
	}
	if resolved.Title != "Completed task" {
		t.Errorf("Resolve(task:2) Title = %q, want %q", resolved.Title, "Completed task")
	}
	if resolved.Detail != "Status: completed" {
		t.Errorf("Resolve(task:2) Detail = %q, want %q", resolved.Detail, "Status: completed")
	}
}

func TestResolveNotFound(t *testing.T) {
	contextDir := t.TempDir()

	// Empty DECISIONS.md
	if err := os.WriteFile(filepath.Join(contextDir, "DECISIONS.md"), []byte("# Decisions\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(DECISIONS.md) error: %v", err)
	}

	resolved := Resolve("decision:999", contextDir)

	if resolved.Found {
		t.Errorf("Resolve(decision:999) Found = true, want false")
	}
	if resolved.Raw != "decision:999" {
		t.Errorf("Resolve(decision:999) Raw = %q, want %q", resolved.Raw, "decision:999")
	}
}

func TestResolveNote(t *testing.T) {
	contextDir := t.TempDir()

	resolved := Resolve(`"Hotfix for production bug"`, contextDir)

	if !resolved.Found {
		t.Fatalf("Resolve(note) Found = false, want true")
	}
	if resolved.Title != "Hotfix for production bug" {
		t.Errorf("Resolve(note) Title = %q, want %q", resolved.Title, "Hotfix for production bug")
	}
	if resolved.Type != "note" {
		t.Errorf("Resolve(note) Type = %q, want %q", resolved.Type, "note")
	}
}

func TestResolveSession(t *testing.T) {
	contextDir := t.TempDir()

	resolved := Resolve("session:abc123", contextDir)

	if !resolved.Found {
		t.Fatalf("Resolve(session:abc123) Found = false, want true")
	}
	if resolved.Title != "abc123" {
		t.Errorf("Resolve(session:abc123) Title = %q, want %q", resolved.Title, "abc123")
	}
	if resolved.Type != "session" {
		t.Errorf("Resolve(session:abc123) Type = %q, want %q", resolved.Type, "session")
	}
}
