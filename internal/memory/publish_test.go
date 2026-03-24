//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	cfgMem "github.com/ActiveMemory/ctx/internal/config/memory"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestMergePublished_EmptyFile(t *testing.T) {
	published := "# Project Context (managed by ctx)\n\n## Pending Tasks\n- [ ] task one\n"
	merged, missing := MergePublished("", published)

	if !strings.Contains(merged, marker.PublishMarkerStart) {
		t.Error("expected marker start in output")
	}
	if !strings.Contains(merged, marker.PublishMarkerEnd) {
		t.Error("expected marker end in output")
	}
	if !strings.Contains(merged, "task one") {
		t.Error("expected published content in output")
	}
	if !missing {
		t.Error("expected missing=true for empty file")
	}
}

func TestMergePublished_ReplaceExisting(t *testing.T) {
	existing := "# Auto Memory\n\nClaude notes here.\n\n" +
		marker.PublishMarkerStart + "\nold content\n" + marker.PublishMarkerEnd + "\n\nMore Claude notes.\n"
	published := "# Project Context (managed by ctx)\n\nnew content\n"

	merged, missing := MergePublished(existing, published)

	if missing {
		t.Error("expected missing=false when markers exist")
	}
	if strings.Contains(merged, "old content") {
		t.Error("old content should be replaced")
	}
	if !strings.Contains(merged, "new content") {
		t.Error("new content should be present")
	}
	if !strings.Contains(merged, "Claude notes here") {
		t.Error("Claude content before markers should be preserved")
	}
	if !strings.Contains(merged, "More Claude notes") {
		t.Error("Claude content after markers should be preserved")
	}
}

func TestMergePublished_MarkersStripped(t *testing.T) {
	existing := "# Auto Memory\n\nClaude rewrote everything.\n"
	published := "# Project Context (managed by ctx)\n\nnew block\n"

	merged, missing := MergePublished(existing, published)

	if !missing {
		t.Error("expected missing=true when markers absent")
	}
	if !strings.Contains(merged, "Claude rewrote everything") {
		t.Error("existing content should be preserved")
	}
	if !strings.Contains(merged, "new block") {
		t.Error("published block should be appended")
	}
}

func TestRemovePublished(t *testing.T) {
	content := "# Auto Memory\n\nNotes.\n\n" +
		marker.PublishMarkerStart + "\npublished stuff\n" + marker.PublishMarkerEnd + "\n\nMore notes.\n"

	cleaned, found := RemovePublished(content)

	if !found {
		t.Error("expected found=true")
	}
	if strings.Contains(cleaned, "published stuff") {
		t.Error("published block should be removed")
	}
	if !strings.Contains(cleaned, "Notes.") {
		t.Error("content before markers should remain")
	}
	if !strings.Contains(cleaned, "More notes.") {
		t.Error("content after markers should remain")
	}
}

func TestRemovePublished_NoMarkers(t *testing.T) {
	content := "# Memory\n\nJust notes.\n"
	cleaned, found := RemovePublished(content)

	if found {
		t.Error("expected found=false when no markers")
	}
	if cleaned != content {
		t.Error("content should be unchanged")
	}
}

func TestPublishResult_TrimToBudget(t *testing.T) {
	r := PublishResult{
		Tasks:       []string{"- [ ] t1", "- [ ] t2"},
		Decisions:   []string{"d1", "d2", "d3"},
		Conventions: []string{"- c1", "- c2", "- c3"},
		Learnings:   []string{"l1", "l2", "l3"},
	}

	// Very tight budget: should trim learnings and conventions first
	r.trimToBudget(10)

	if len(r.Learnings) != 0 {
		t.Errorf("expected 0 learnings after tight trim, got %d", len(r.Learnings))
	}
	if len(r.Tasks) != 2 {
		t.Errorf("tasks should not be trimmed, got %d", len(r.Tasks))
	}
}

func TestSelectContent(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	rc.Reset()
	defer func() { _ = os.Chdir(origDir) }()

	contextDir := filepath.Join(workDir, dir.Context)
	if mkErr := os.MkdirAll(contextDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	// Create TASKS.md with pending items
	tasks := "# Tasks\n\n- [x] done task\n- [ ] pending task one\n- [ ] pending task two\n"
	if writeErr := os.WriteFile(filepath.Join(contextDir, ctx.Task), []byte(tasks), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Create DECISIONS.md with a recent entry
	ts := time.Now().Format(cfgTime.TimestampCompact)
	decisions := fmt.Sprintf("# Decisions\n\n## [%s] Use SQLite\n\nContext: testing\n", ts)
	if writeErr := os.WriteFile(filepath.Join(contextDir, ctx.Decision), []byte(decisions), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Create CONVENTIONS.md
	conventions := "# Conventions\n\n- Always use ctx from PATH\n- Prefer filepath.Join\n"
	if writeErr := os.WriteFile(filepath.Join(contextDir, ctx.Convention), []byte(conventions), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Create empty LEARNINGS.md
	if writeErr := os.WriteFile(filepath.Join(contextDir, ctx.Learning), []byte("# Learnings\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	result, selectErr := SelectContent(contextDir, cfgMem.DefaultPublishBudget)
	if selectErr != nil {
		t.Fatalf("SelectContent: %v", selectErr)
	}

	if len(result.Tasks) != 2 {
		t.Errorf("expected 2 pending tasks, got %d", len(result.Tasks))
	}
	if len(result.Decisions) != 1 {
		t.Errorf("expected 1 decision, got %d", len(result.Decisions))
	}
	if len(result.Conventions) != 2 {
		t.Errorf("expected 2 conventions, got %d", len(result.Conventions))
	}

	formatted := result.Format()
	if !strings.Contains(formatted, "pending task one") {
		t.Error("expected pending task in formatted output")
	}
	if !strings.Contains(formatted, "SQLite") {
		t.Error("expected decision in formatted output")
	}
}
