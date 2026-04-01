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

	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
)

func TestCollectDeduplicates(t *testing.T) {
	contextDir := t.TempDir()

	// Create state directory for pending records.
	stateDir := filepath.Join(contextDir, cfgDir.State)
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}

	// Write TASKS.md with one pending task so WorkingRefs returns "task:1".
	tasks := "# Tasks\n\n- [ ] First pending task\n"
	if err := os.WriteFile(filepath.Join(contextDir, "TASKS.md"), []byte(tasks), 0o600); err != nil {
		t.Fatalf("WriteFile(TASKS.md) error: %v", err)
	}

	// Record "task:1" in pending — it will also appear from WorkingRefs.
	if err := Record("task:1", stateDir); err != nil {
		t.Fatalf("Record(task:1) error: %v", err)
	}
	// Record "decision:5" in pending — appears only once.
	if err := Record("decision:5", stateDir); err != nil {
		t.Fatalf("Record(decision:5) error: %v", err)
	}

	refs := Collect(contextDir)

	seen := map[string]int{}
	for _, r := range refs {
		seen[r]++
	}

	if seen["task:1"] != 1 {
		t.Errorf("task:1 count = %d, want 1 (deduplication failed); refs = %v", seen["task:1"], refs)
	}
	if seen["decision:5"] != 1 {
		t.Errorf("decision:5 count = %d, want 1; refs = %v", seen["decision:5"], refs)
	}
}

func TestCollectEmptyReturnsNil(t *testing.T) {
	contextDir := t.TempDir()

	// Create state directory but no pending file; no TASKS.md; no session env.
	stateDir := filepath.Join(contextDir, cfgDir.State)
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}

	// Write an empty TASKS.md so WorkingRefs finds no pending tasks.
	if err := os.WriteFile(filepath.Join(contextDir, "TASKS.md"), []byte("# Tasks\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(TASKS.md) error: %v", err)
	}

	refs := Collect(contextDir)
	if len(refs) != 0 {
		t.Errorf("Collect() returned %v, want empty", refs)
	}
}

func TestFormatTrailer(t *testing.T) {
	refs := []string{"decision:12", "task:8", "session:abc123"}
	got := FormatTrailer(refs)
	want := "ctx-context: decision:12, task:8, session:abc123"
	if got != want {
		t.Errorf("FormatTrailer() = %q, want %q", got, want)
	}
}

func TestFormatTrailerEmpty(t *testing.T) {
	got := FormatTrailer(nil)
	if got != "" {
		t.Errorf("FormatTrailer(nil) = %q, want empty string", got)
	}
}
