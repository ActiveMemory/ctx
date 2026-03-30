//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgMemory "github.com/ActiveMemory/ctx/internal/config/memory"
)

const fixtureMemory = `# Auto Memory

## Session 2026-03-05: Memory Bridge Design

Worked on the memory bridge foundation today.

- always use ctx from PATH
- decided to use heuristic classification over LLM-based
- learned that symlinks produce different slugs across machines
- need to add integration tests for import

Some generic session notes that should be skipped.
`

func TestIntegration_ParseClassifyPromote(t *testing.T) {
	contextDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Parse
	entries := Entries(fixtureMemory)
	if len(entries) == 0 {
		t.Fatal("expected entries from fixture")
	}

	// Classify and promote
	state := State{ImportedHashes: []string{}}
	var imported, skipped int

	for _, entry := range entries {
		hash := EntryHash(entry.Text)
		if state.Imported(hash) {
			continue
		}

		classification := Classify(entry)
		if classification.Target == cfgMemory.TargetSkip {
			skipped++
			continue
		}

		if promoteErr := Promote(entry, classification); promoteErr != nil {
			t.Fatalf("Promote failed for %q: %v", entry.Text[:40], promoteErr)
		}
		state.MarkImported(hash, classification.Target)
		imported++
	}

	if imported == 0 {
		t.Fatal("expected at least one imported entry")
	}
	if skipped == 0 {
		t.Fatal("expected at least one skipped entry")
	}

	// Verify entries landed in correct files
	convData, _ := os.ReadFile(filepath.Join(contextDir, ctx.Convention))
	if !strings.Contains(string(convData), "ctx from PATH") {
		t.Error("expected convention 'always use ctx from PATH' in CONVENTIONS.md")
	}

	decData, _ := os.ReadFile(filepath.Join(contextDir, ctx.Decision))
	if !strings.Contains(string(decData), "heuristic classification") {
		t.Error("expected decision about classification in DECISIONS.md")
	}

	lrnData, _ := os.ReadFile(filepath.Join(contextDir, ctx.Learning))
	if !strings.Contains(string(lrnData), "symlinks") {
		t.Error("expected learning about symlinks in LEARNINGS.md")
	}

	taskData, _ := os.ReadFile(filepath.Join(contextDir, ctx.Task))
	if !strings.Contains(string(taskData), "integration tests") {
		t.Error("expected task about integration tests in TASKS.md")
	}

	// Verify dedup: re-run should import zero
	var reimported int
	for _, entry := range entries {
		hash := EntryHash(entry.Text)
		if state.Imported(hash) {
			continue
		}
		classification := Classify(entry)
		if classification.Target == cfgMemory.TargetSkip {
			continue
		}
		reimported++
	}
	if reimported != 0 {
		t.Errorf("expected 0 reimports after dedup, got %d", reimported)
	}
}
