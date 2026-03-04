//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package changes

import (
	"strings"
	"testing"
	"time"
)

func TestHumanAgo(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{5 * time.Second, "just now"},
		{30 * time.Second, "just now"},
		{5 * time.Minute, "5 minutes ago"},
		{1 * time.Minute, "1 minute ago"},
		{3 * time.Hour, "3 hours ago"},
		{1 * time.Hour, "1 hour ago"},
		{48 * time.Hour, "2 days ago"},
		{24 * time.Hour, "1 day ago"},
	}
	for _, tt := range tests {
		if got := humanAgo(tt.d); got != tt.want {
			t.Errorf("humanAgo(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestExtractTimestamp(t *testing.T) {
	line := `{"event":"context-load-gate","timestamp":"2026-03-03T08:00:00Z","session":"abc"}`
	ts, ok := extractTimestamp(line)
	if !ok {
		t.Fatal("extractTimestamp returned false")
	}
	if ts.Year() != 2026 || ts.Month() != 3 || ts.Day() != 3 {
		t.Errorf("unexpected timestamp: %v", ts)
	}

	// No timestamp.
	_, ok = extractTimestamp(`{"event":"other"}`)
	if ok {
		t.Error("expected false for line without timestamp")
	}
}

func TestParseSinceFlag(t *testing.T) {
	// Duration.
	ts, label, err := parseSinceFlag("6h")
	if err != nil {
		t.Fatalf("parseSinceFlag(6h) error: %v", err)
	}
	if !strings.Contains(label, "hour") {
		t.Errorf("expected label with 'hour', got: %s", label)
	}
	if time.Since(ts) < 5*time.Hour {
		t.Errorf("timestamp too recent: %v", ts)
	}

	// Date.
	ts, label, err = parseSinceFlag("2026-03-01")
	if err != nil {
		t.Fatalf("parseSinceFlag(2026-03-01) error: %v", err)
	}
	if label != "since 2026-03-01" {
		t.Errorf("unexpected label: %s", label)
	}
	if ts.Year() != 2026 || ts.Month() != 3 || ts.Day() != 1 {
		t.Errorf("unexpected date: %v", ts)
	}

	// Invalid.
	_, _, err = parseSinceFlag("not-a-date")
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		n    int
		unit string
		want string
	}{
		{1, "commit", "1 commit"},
		{5, "commit", "5 commits"},
		{0, "file", "0 files"},
	}
	for _, tt := range tests {
		if got := pluralize(tt.n, tt.unit); got != tt.want {
			t.Errorf("pluralize(%d, %q) = %q, want %q", tt.n, tt.unit, got, tt.want)
		}
	}
}

func TestUniqueTopDirs(t *testing.T) {
	input := "internal/cli/deps/deps.go\ninternal/cli/changes/changes.go\ndocs/index.md\nREADME.md\n"
	got := uniqueTopDirs(input)
	want := []string{"README.md", "docs", "internal"}
	if len(got) != len(want) {
		t.Fatalf("uniqueTopDirs: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("uniqueTopDirs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestUniqueLines(t *testing.T) {
	input := "Alice\nBob\nAlice\nCharlie\n"
	got := uniqueLines(input)
	want := []string{"Alice", "Bob", "Charlie"}
	if len(got) != len(want) {
		t.Fatalf("uniqueLines: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("uniqueLines[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRenderChanges(t *testing.T) {
	ctxChanges := []ContextChange{
		{Name: "TASKS.md", ModTime: time.Date(2026, 3, 3, 14, 30, 0, 0, time.UTC)},
	}
	code := CodeSummary{
		CommitCount: 5,
		LatestMsg:   "Add deps command",
		Dirs:        []string{"internal", "docs"},
		Authors:     []string{"Volkan"},
	}

	out := RenderChanges("6 hours ago", ctxChanges, code)
	if !strings.Contains(out, "## Changes Since Last Session") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "TASKS.md") {
		t.Error("missing context change")
	}
	if !strings.Contains(out, "5 commits") {
		t.Error("missing commit count")
	}
	if !strings.Contains(out, "Add deps command") {
		t.Error("missing latest message")
	}
}

func TestRenderChangesForHook(t *testing.T) {
	ctxChanges := []ContextChange{
		{Name: "TASKS.md", ModTime: time.Now()},
	}
	code := CodeSummary{CommitCount: 3, LatestMsg: "Fix bug"}

	out := RenderChangesForHook("2 hours ago", ctxChanges, code)
	if !strings.Contains(out, "Changes since last session") {
		t.Error("missing hook header")
	}
	if !strings.Contains(out, "TASKS.md") {
		t.Error("missing file name in hook output")
	}

	// Empty case.
	out = RenderChangesForHook("1 hour ago", nil, CodeSummary{})
	if out != "" {
		t.Errorf("expected empty for no changes, got: %q", out)
	}
}

func TestRenderChanges_NoChanges(t *testing.T) {
	out := RenderChanges("1 hour ago", nil, CodeSummary{})
	if !strings.Contains(out, "No changes detected") {
		t.Error("expected 'No changes detected' message")
	}
}
