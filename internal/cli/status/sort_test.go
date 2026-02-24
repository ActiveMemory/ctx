//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/context"
)

func TestGetRecentFiles(t *testing.T) {
	now := time.Now()
	files := []context.FileInfo{
		{Name: "old.md", ModTime: now.Add(-3 * time.Hour)},
		{Name: "newest.md", ModTime: now},
		{Name: "mid.md", ModTime: now.Add(-1 * time.Hour)},
	}

	tests := []struct {
		name  string
		n     int
		want  int
		first string
	}{
		{"n less than len", 2, 2, "newest.md"},
		{"n equals len", 3, 3, "newest.md"},
		{"n greater than len", 5, 3, "newest.md"},
		{"n is zero", 0, 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRecentFiles(files, tt.n)
			if len(got) != tt.want {
				t.Errorf("getRecentFiles(n=%d) returned %d files, want %d", tt.n, len(got), tt.want)
			}
			if tt.first != "" && len(got) > 0 && got[0].Name != tt.first {
				t.Errorf("first file = %q, want %q", got[0].Name, tt.first)
			}
		})
	}

	t.Run("empty input", func(t *testing.T) {
		got := getRecentFiles(nil, 5)
		if len(got) != 0 {
			t.Errorf("getRecentFiles(nil) returned %d files, want 0", len(got))
		}
	})
}

func TestSortFilesByPriority(t *testing.T) {
	files := []context.FileInfo{
		{Name: "TASKS.md"},
		{Name: "CONVENTIONS.md"},
		{Name: "CONSTITUTION.md"},
	}

	sortFilesByPriority(files)

	// CONSTITUTION should come before TASKS
	constitutionIdx := -1
	tasksIdx := -1
	for i, f := range files {
		if f.Name == "CONSTITUTION.md" {
			constitutionIdx = i
		}
		if f.Name == "TASKS.md" {
			tasksIdx = i
		}
	}
	if constitutionIdx >= tasksIdx {
		t.Errorf("CONSTITUTION.md (idx=%d) should sort before TASKS.md (idx=%d)", constitutionIdx, tasksIdx)
	}
}
