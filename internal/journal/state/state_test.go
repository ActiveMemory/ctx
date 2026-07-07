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

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
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

	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
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

func TestLoad_V1Tolerant(t *testing.T) {
	dir := t.TempDir()

	// A raw v1 state file: version 1, entries only, no sessions map.
	v1 := `{"version":1,"entries":{"2026-01-21-old.md":{"exported":"2026-01-21"}}}`
	statePath := filepath.Join(dir, journal.File)
	if err := os.WriteFile(statePath, []byte(v1), fs.PermFile); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load v1: %v", err)
	}

	// Version normalises to current on load.
	if loaded.Version != CurrentVersion {
		t.Errorf("Version = %d, want %d", loaded.Version, CurrentVersion)
	}
	// Sessions map is initialised (non-nil) but empty.
	if loaded.Sessions == nil {
		t.Error("Sessions should be initialised, not nil")
	}
	if len(loaded.Sessions) != 0 {
		t.Errorf("Sessions should be empty, got %d", len(loaded.Sessions))
	}
	// Existing entries are preserved untouched.
	if !loaded.Exported("2026-01-21-old.md") {
		t.Error("v1 entry should be preserved after load")
	}
}

func TestSourceRoundTrip(t *testing.T) {
	dir := t.TempDir()

	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
			"2026-01-21-test.md": {
				Exported:   "2026-01-21",
				RenderHash: "deadbeef",
			},
		},
		Sessions: map[string]Source{
			"sess-abc": {
				SourceFile:  "/transcripts/sess-abc.jsonl",
				SourceMtime: 1234567890,
				SourceSize:  4096,
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

	src, ok := loaded.SessionSource("sess-abc")
	if !ok {
		t.Fatal("session not found after round-trip")
	}
	if src.SourceFile != "/transcripts/sess-abc.jsonl" {
		t.Errorf("SourceFile = %q, want preserved", src.SourceFile)
	}
	if src.SourceMtime != 1234567890 {
		t.Errorf("SourceMtime = %d, want 1234567890", src.SourceMtime)
	}
	if src.SourceSize != 4096 {
		t.Errorf("SourceSize = %d, want 4096", src.SourceSize)
	}
	if got := loaded.Entries["2026-01-21-test.md"].RenderHash; got != "deadbeef" {
		t.Errorf("RenderHash = %q, want %q", got, "deadbeef")
	}
}

func TestSessionSourceAndMark(t *testing.T) {
	s := &State{
		Version:  CurrentVersion,
		Entries:  make(map[string]File),
		Sessions: make(map[string]Source),
	}

	if _, ok := s.SessionSource("nope"); ok {
		t.Error("unknown session should report not found")
	}

	s.MarkSource("sess-1", "/t/sess-1.jsonl", 100, 200)
	src, ok := s.SessionSource("sess-1")
	if !ok {
		t.Fatal("session should be found after MarkSource")
	}
	if src.SourceFile != "/t/sess-1.jsonl" || src.SourceMtime != 100 ||
		src.SourceSize != 200 {
		t.Errorf("unexpected source stats: %+v", src)
	}

	// MarkSource overwrites with the latest stats.
	s.MarkSource("sess-1", "/t/sess-1.jsonl", 150, 300)
	src, _ = s.SessionSource("sess-1")
	if src.SourceMtime != 150 || src.SourceSize != 300 {
		t.Errorf("MarkSource should overwrite, got %+v", src)
	}
}

func TestMarkSource_NilMap(t *testing.T) {
	// A State with a nil Sessions map (e.g. hand-constructed) must not
	// panic on MarkSource.
	s := &State{Version: CurrentVersion, Entries: make(map[string]File)}
	s.MarkSource("sess-1", "/t/sess-1.jsonl", 1, 2)
	if _, ok := s.SessionSource("sess-1"); !ok {
		t.Error("MarkSource should initialise a nil Sessions map")
	}
}

func TestCountUnenriched(t *testing.T) {
	dir := t.TempDir()

	// Create some .md files
	for _, name := range []string{"a.md", "b.md", "c.md"} {
		fPath := filepath.Join(dir, name)
		if err := os.WriteFile(fPath, []byte("content"), fs.PermFile); err != nil {
			t.Fatal(err)
		}
	}
	// Create a non-md file that should be ignored
	statePath := filepath.Join(dir, "state.json")
	if err := os.WriteFile(statePath, []byte("{}"), fs.PermFile); err != nil {
		t.Fatal(err)
	}

	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
			"a.md": {Enriched: "2026-01-21"},
		},
	}

	count := s.CountUnenriched(dir)
	if count != 2 {
		t.Errorf("CountUnenriched = %d, want 2", count)
	}
}

func TestRename(t *testing.T) {
	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
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
	s := &State{
		Version: CurrentVersion,
		Entries: make(map[string]File),
	}

	// Should not panic or create entries
	s.Rename("nonexistent.md", "new.md")
	if len(s.Entries) != 0 {
		t.Error("Rename of nonexistent should be no-op")
	}
}

func TestQueryHelpers(t *testing.T) {
	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
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

	if !s.Exported("full.md") {
		t.Error("full.md should be exported")
	}
	if !s.Enriched("full.md") {
		t.Error("full.md should be enriched")
	}
	if !s.Normalized("full.md") {
		t.Error("full.md should be normalized")
	}
	if !s.FencesVerified("full.md") {
		t.Error("full.md should have fences verified")
	}

	if !s.Exported("partial.md") {
		t.Error("partial.md should be exported")
	}
	if s.Enriched("partial.md") {
		t.Error("partial.md should not be enriched")
	}

	if s.Exported("missing.md") {
		t.Error("missing.md should not be exported")
	}
}

func TestClearEnriched(t *testing.T) {
	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
			"test.md": {
				Exported: "2026-01-21",
				Enriched: "2026-01-22",
			},
		},
	}

	if !s.Enriched("test.md") {
		t.Fatal("should be enriched before clear")
	}

	s.ClearEnriched("test.md")

	if s.Enriched("test.md") {
		t.Error("should not be enriched after ClearEnriched")
	}
	// Other fields should be untouched
	if !s.Exported("test.md") {
		t.Error("exported should be preserved after ClearEnriched")
	}
}

func TestClearEnriched_NoOp(t *testing.T) {
	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
			"test.md": {Exported: "2026-01-21"},
		},
	}

	// Should not panic on file that isn't enriched
	s.ClearEnriched("test.md")
	if s.Enriched("test.md") {
		t.Error("should remain unenriched")
	}

	// Should not panic on missing entry
	s.ClearEnriched("nonexistent.md")
}

func TestMark(t *testing.T) {
	s := &State{
		Version: CurrentVersion,
		Entries: make(map[string]File),
	}

	if ok := s.Mark("test.md", journal.StageExported); !ok {
		t.Error("Mark exported should succeed")
	}
	if !s.Exported("test.md") {
		t.Error("test.md should be exported after Mark")
	}

	if ok := s.Mark("test.md", "invalid_stage"); ok {
		t.Error("Mark invalid stage should fail")
	}
}

func TestMark_AllStages(t *testing.T) {
	s := &State{
		Version: CurrentVersion,
		Entries: make(map[string]File),
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
	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
			"test.md": {
				Exported: "2026-01-21",
				Enriched: "2026-01-22",
				Locked:   "2026-01-23",
			},
		},
	}

	if ok := s.Clear("test.md", journal.StageLocked); !ok {
		t.Error("Clear locked should succeed")
	}
	if s.Locked("test.md") {
		t.Error("should not be locked after Clear")
	}
	// Other fields preserved.
	if !s.Exported("test.md") {
		t.Error("exported should be preserved after Clear locked")
	}
	if !s.Enriched("test.md") {
		t.Error("enriched should be preserved after Clear locked")
	}

	if ok := s.Clear("test.md", "invalid"); ok {
		t.Error("Clear invalid stage should fail")
	}
}

func TestClear_AllStages(t *testing.T) {
	s := &State{
		Version: CurrentVersion,
		Entries: make(map[string]File),
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
	s := &State{
		Version: CurrentVersion,
		Entries: make(map[string]File),
	}

	if s.Locked("test.md") {
		t.Error("should not be locked initially")
	}

	s.Mark("test.md", journal.StageLocked)
	if !s.Locked("test.md") {
		t.Error("should be locked after Mark")
	}

	s.Clear("test.md", journal.StageLocked)
	if s.Locked("test.md") {
		t.Error("should not be locked after Clear")
	}
}

func TestLocked_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	s := &State{
		Version: CurrentVersion,
		Entries: map[string]File{
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
