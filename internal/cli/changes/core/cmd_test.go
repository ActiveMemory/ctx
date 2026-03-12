//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
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
		if got := HumanAgo(tt.d); got != tt.want {
			t.Errorf("HumanAgo(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestExtractTimestamp(t *testing.T) {
	line := `{"event":"context-load-gate","timestamp":"2026-03-03T08:00:00Z","session":"abc"}`
	ts, ok := ExtractTimestamp(line)
	if !ok {
		t.Fatal("ExtractTimestamp returned false")
	}
	if ts.Year() != 2026 || ts.Month() != 3 || ts.Day() != 3 {
		t.Errorf("unexpected timestamp: %v", ts)
	}

	// No timestamp.
	_, ok = ExtractTimestamp(`{"event":"other"}`)
	if ok {
		t.Error("expected false for line without timestamp")
	}
}

func TestParseSinceFlag(t *testing.T) {
	// Duration.
	ts, label, err := ParseSinceFlag("6h")
	if err != nil {
		t.Fatalf("ParseSinceFlag(6h) error: %v", err)
	}
	if !strings.Contains(label, "hour") {
		t.Errorf("expected label with 'hour', got: %s", label)
	}
	if time.Since(ts) < 5*time.Hour {
		t.Errorf("timestamp too recent: %v", ts)
	}

	// Date.
	ts, label, err = ParseSinceFlag("2026-03-01")
	if err != nil {
		t.Fatalf("ParseSinceFlag(2026-03-01) error: %v", err)
	}
	if label != "since 2026-03-01" {
		t.Errorf("unexpected label: %s", label)
	}
	if ts.Year() != 2026 || ts.Month() != 3 || ts.Day() != 1 {
		t.Errorf("unexpected date: %v", ts)
	}

	// Invalid.
	_, _, err = ParseSinceFlag("not-a-date")
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
		if got := Pluralize(tt.n, tt.unit); got != tt.want {
			t.Errorf("Pluralize(%d, %q) = %q, want %q", tt.n, tt.unit, got, tt.want)
		}
	}
}

func TestUniqueTopDirs(t *testing.T) {
	input := "internal/cli/deps/deps.go\ninternal/cli/changes/changes.go\ndocs/index.md\nREADME.md\n"
	got := UniqueTopDirs(input)
	want := []string{"README.md", "docs", "internal"}
	if len(got) != len(want) {
		t.Fatalf("UniqueTopDirs: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("UniqueTopDirs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestUniqueLines(t *testing.T) {
	input := "Alice\nBob\nAlice\nCharlie\n"
	got := UniqueLines(input)
	want := []string{"Alice", "Bob", "Charlie"}
	if len(got) != len(want) {
		t.Fatalf("UniqueLines: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("UniqueLines[%d] = %q, want %q", i, got[i], want[i])
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

func TestDetectReferenceTime_SinceFlag(t *testing.T) {
	_, label, detectErr := DetectReferenceTime("6h")
	if detectErr != nil {
		t.Fatalf("DetectReferenceTime(6h) error: %v", detectErr)
	}
	if !strings.Contains(label, "hour") {
		t.Errorf("expected label containing 'hour', got: %s", label)
	}
}

func TestDetectReferenceTime_Fallback(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CTX_DIR", tmp)
	rc.Reset()

	stateDir := filepath.Join(tmp, dir.State)
	mkErr := os.MkdirAll(stateDir, 0o755)
	if mkErr != nil {
		t.Fatalf("MkdirAll: %v", mkErr)
	}

	_, label, detectErr := DetectReferenceTime("")
	if detectErr != nil {
		t.Fatalf("DetectReferenceTime fallback error: %v", detectErr)
	}
	if !strings.Contains(label, "24 hour") {
		t.Errorf("expected label containing '24 hour', got: %s", label)
	}
}

func TestDetectReferenceTime_FromMarkers(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CTX_DIR", tmp)
	rc.Reset()

	stateDir := filepath.Join(tmp, dir.State)
	mkErr := os.MkdirAll(stateDir, 0o755)
	if mkErr != nil {
		t.Fatalf("MkdirAll: %v", mkErr)
	}

	// Create two marker files with different mtimes.
	older := filepath.Join(stateDir, "ctx-loaded-aaa")
	newer := filepath.Join(stateDir, "ctx-loaded-bbb")

	writeErr := os.WriteFile(older, []byte(""), 0o644)
	if writeErr != nil {
		t.Fatalf("WriteFile older: %v", writeErr)
	}
	writeErr = os.WriteFile(newer, []byte(""), 0o644)
	if writeErr != nil {
		t.Fatalf("WriteFile newer: %v", writeErr)
	}

	olderTime := time.Now().Add(-2 * time.Hour)
	newerTime := time.Now().Add(-30 * time.Minute)

	chtErr := os.Chtimes(older, olderTime, olderTime)
	if chtErr != nil {
		t.Fatalf("Chtimes older: %v", chtErr)
	}
	chtErr = os.Chtimes(newer, newerTime, newerTime)
	if chtErr != nil {
		t.Fatalf("Chtimes newer: %v", chtErr)
	}

	refTime, _, detectErr := DetectReferenceTime("")
	if detectErr != nil {
		t.Fatalf("DetectReferenceTime from markers error: %v", detectErr)
	}

	// Should return the second most recent (older) marker time.
	diff := refTime.Sub(olderTime)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("expected refTime ~%v, got %v (diff=%v)", olderTime, refTime, diff)
	}
}

func TestFindContextChanges(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CTX_DIR", tmp)
	rc.Reset()

	// Create two .md files with different mtimes.
	recentFile := filepath.Join(tmp, "TASKS.md")
	oldFile := filepath.Join(tmp, "OLD.md")

	writeErr := os.WriteFile(recentFile, []byte("# Tasks"), 0o644)
	if writeErr != nil {
		t.Fatalf("WriteFile recent: %v", writeErr)
	}
	writeErr = os.WriteFile(oldFile, []byte("# Old"), 0o644)
	if writeErr != nil {
		t.Fatalf("WriteFile old: %v", writeErr)
	}

	// Set old file to 48 hours ago.
	oldTime := time.Now().Add(-48 * time.Hour)
	chtErr := os.Chtimes(oldFile, oldTime, oldTime)
	if chtErr != nil {
		t.Fatalf("Chtimes old: %v", chtErr)
	}

	// Reference time between old and recent.
	refTime := time.Now().Add(-24 * time.Hour)
	changes, findErr := FindContextChanges(refTime)
	if findErr != nil {
		t.Fatalf("FindContextChanges error: %v", findErr)
	}

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Name != "TASKS.md" {
		t.Errorf("expected TASKS.md, got %s", changes[0].Name)
	}
}

func TestFindContextChanges_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CTX_DIR", tmp)
	rc.Reset()

	refTime := time.Now().Add(-1 * time.Hour)
	changes, findErr := FindContextChanges(refTime)
	if findErr != nil {
		t.Fatalf("FindContextChanges error: %v", findErr)
	}
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}
