//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("Load missing file: %v", err)
	}
	if s.Version != CurrentVersion {
		t.Errorf("Version = %d, want %d", s.Version, CurrentVersion)
	}
	if len(s.Entries) != 0 {
		t.Errorf("Entries should be empty, got %d", len(s.Entries))
	}
}

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()

	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"2026-01-21-test-abc12345.md": {
				Exported: "2026-01-21",
				Enriched: "2026-01-22",
			},
		},
	}

	if err := s.Save(dir); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Version != CurrentVersion {
		t.Errorf("Version = %d, want %d", loaded.Version, CurrentVersion)
	}

	fs, ok := loaded.Entries["2026-01-21-test-abc12345.md"]
	if !ok {
		t.Fatal("entry not found after round-trip")
	}
	if fs.Exported != "2026-01-21" {
		t.Errorf("Exported = %q, want %q", fs.Exported, "2026-01-21")
	}
	if fs.Enriched != "2026-01-22" {
		t.Errorf("Enriched = %q, want %q", fs.Enriched, "2026-01-22")
	}
}

func TestCountUnenriched(t *testing.T) {
	dir := t.TempDir()

	// Create some .md files
	for _, name := range []string{"a.md", "b.md", "c.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), config.PermFile); err != nil {
			t.Fatal(err)
		}
	}
	// Create a non-md file that should be ignored
	if err := os.WriteFile(filepath.Join(dir, "state.json"), []byte("{}"), config.PermFile); err != nil {
		t.Fatal(err)
	}

	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"a.md": {Enriched: "2026-01-21"},
		},
	}

	count := s.CountUnenriched(dir)
	if count != 2 {
		t.Errorf("CountUnenriched = %d, want 2", count)
	}
}

func TestRename(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"old-name.md": {
				Exported:       "2026-01-21",
				Enriched:       "2026-01-22",
				Normalized:     "2026-01-23",
				FencesVerified: "2026-01-24",
			},
		},
	}

	s.Rename("old-name.md", "new-name.md")

	if _, ok := s.Entries["old-name.md"]; ok {
		t.Error("old entry should be deleted")
	}

	fs, ok := s.Entries["new-name.md"]
	if !ok {
		t.Fatal("new entry not found")
	}
	if fs.Exported != "2026-01-21" {
		t.Errorf("Exported = %q, want preserved value", fs.Exported)
	}
	if fs.FencesVerified != "2026-01-24" {
		t.Errorf("FencesVerified = %q, want preserved value", fs.FencesVerified)
	}
}

func TestRename_NoOp(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: make(map[string]FileState),
	}

	// Should not panic or create entries
	s.Rename("nonexistent.md", "new.md")
	if len(s.Entries) != 0 {
		t.Error("Rename of nonexistent should be no-op")
	}
}

func TestQueryHelpers(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"full.md": {
				Exported:       "2026-01-21",
				Enriched:       "2026-01-22",
				Normalized:     "2026-01-23",
				FencesVerified: "2026-01-24",
			},
			"partial.md": {
				Exported: "2026-01-21",
			},
		},
	}

	if !s.IsExported("full.md") {
		t.Error("full.md should be exported")
	}
	if !s.IsEnriched("full.md") {
		t.Error("full.md should be enriched")
	}
	if !s.IsNormalized("full.md") {
		t.Error("full.md should be normalized")
	}
	if !s.IsFencesVerified("full.md") {
		t.Error("full.md should have fences verified")
	}

	if !s.IsExported("partial.md") {
		t.Error("partial.md should be exported")
	}
	if s.IsEnriched("partial.md") {
		t.Error("partial.md should not be enriched")
	}

	if s.IsExported("missing.md") {
		t.Error("missing.md should not be exported")
	}
}

func TestClearEnriched(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"test.md": {
				Exported: "2026-01-21",
				Enriched: "2026-01-22",
			},
		},
	}

	if !s.IsEnriched("test.md") {
		t.Fatal("should be enriched before clear")
	}

	s.ClearEnriched("test.md")

	if s.IsEnriched("test.md") {
		t.Error("should not be enriched after ClearEnriched")
	}
	// Other fields should be untouched
	if !s.IsExported("test.md") {
		t.Error("exported should be preserved after ClearEnriched")
	}
}

func TestClearEnriched_NoOp(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"test.md": {Exported: "2026-01-21"},
		},
	}

	// Should not panic on file that isn't enriched
	s.ClearEnriched("test.md")
	if s.IsEnriched("test.md") {
		t.Error("should remain unenriched")
	}

	// Should not panic on missing entry
	s.ClearEnriched("nonexistent.md")
}

func TestMark(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: make(map[string]FileState),
	}

	if ok := s.Mark("test.md", "exported"); !ok {
		t.Error("Mark exported should succeed")
	}
	if !s.IsExported("test.md") {
		t.Error("test.md should be exported after Mark")
	}

	if ok := s.Mark("test.md", "invalid_stage"); ok {
		t.Error("Mark invalid stage should fail")
	}
}

func TestMark_AllStages(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: make(map[string]FileState),
	}

	for _, stage := range ValidStages {
		if ok := s.Mark("test.md", stage); !ok {
			t.Errorf("Mark %q should succeed", stage)
		}
	}

	fs := s.Entries["test.md"]
	if fs.Exported == "" || fs.Enriched == "" || fs.Normalized == "" ||
		fs.FencesVerified == "" || fs.Locked == "" {
		t.Error("all stages should be set")
	}
}

func TestClear(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"test.md": {
				Exported: "2026-01-21",
				Enriched: "2026-01-22",
				Locked:   "2026-01-23",
			},
		},
	}

	if ok := s.Clear("test.md", "locked"); !ok {
		t.Error("Clear locked should succeed")
	}
	if s.Locked("test.md") {
		t.Error("should not be locked after Clear")
	}
	// Other fields preserved.
	if !s.IsExported("test.md") {
		t.Error("exported should be preserved after Clear locked")
	}
	if !s.IsEnriched("test.md") {
		t.Error("enriched should be preserved after Clear locked")
	}

	if ok := s.Clear("test.md", "invalid"); ok {
		t.Error("Clear invalid stage should fail")
	}
}

func TestClear_AllStages(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: make(map[string]FileState),
	}

	// Set all stages, then clear all.
	for _, stage := range ValidStages {
		s.Mark("test.md", stage)
	}
	for _, stage := range ValidStages {
		if ok := s.Clear("test.md", stage); !ok {
			t.Errorf("Clear %q should succeed", stage)
		}
	}

	fs := s.Entries["test.md"]
	if fs.Exported != "" || fs.Enriched != "" || fs.Normalized != "" ||
		fs.FencesVerified != "" || fs.Locked != "" {
		t.Error("all stages should be empty after Clear")
	}
}

func TestLocked(t *testing.T) {
	s := &JournalState{
		Version: CurrentVersion,
		Entries: make(map[string]FileState),
	}

	if s.Locked("test.md") {
		t.Error("should not be locked initially")
	}

	s.Mark("test.md", "locked")
	if !s.Locked("test.md") {
		t.Error("should be locked after Mark")
	}

	s.Clear("test.md", "locked")
	if s.Locked("test.md") {
		t.Error("should not be locked after Clear")
	}
}

func TestLocked_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	s := &JournalState{
		Version: CurrentVersion,
		Entries: map[string]FileState{
			"locked.md":   {Exported: "2026-01-21", Locked: "2026-01-22"},
			"unlocked.md": {Exported: "2026-01-21"},
		},
	}

	if err := s.Save(dir); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !loaded.Locked("locked.md") {
		t.Error("locked.md should be locked after round-trip")
	}
	if loaded.Locked("unlocked.md") {
		t.Error("unlocked.md should not be locked after round-trip")
	}
}
