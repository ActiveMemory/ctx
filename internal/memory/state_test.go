//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestStateRoundtrip(t *testing.T) {
	contextDir := t.TempDir()
	stateDir := filepath.Join(contextDir, config.DirState)
	if mkErr := os.MkdirAll(stateDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	var s State
	s.MarkSynced()

	if saveErr := SaveState(contextDir, s); saveErr != nil {
		t.Fatalf("SaveState: %v", saveErr)
	}

	loaded, loadErr := LoadState(contextDir)
	if loadErr != nil {
		t.Fatalf("LoadState: %v", loadErr)
	}

	if loaded.LastSync == nil {
		t.Fatal("expected LastSync to be set after roundtrip")
	}
	if !loaded.LastSync.Equal(*s.LastSync) {
		t.Errorf("LastSync mismatch: got %v, want %v", *loaded.LastSync, *s.LastSync)
	}
	if loaded.LastImport != nil {
		t.Error("expected LastImport to be nil")
	}
}

func TestLoadState_MissingFile(t *testing.T) {
	contextDir := t.TempDir()

	s, loadErr := LoadState(contextDir)
	if loadErr != nil {
		t.Fatalf("expected nil error for missing file, got: %v", loadErr)
	}
	if s.LastSync != nil {
		t.Error("expected nil LastSync for missing file")
	}
	if s.ImportedHashes == nil {
		t.Error("expected non-nil ImportedHashes slice")
	}
}

func TestEntryHash_Deterministic(t *testing.T) {
	h1 := EntryHash("always use bun for this project")
	h2 := EntryHash("always use bun for this project")
	if h1 != h2 {
		t.Errorf("same input should produce same hash: %q vs %q", h1, h2)
	}
	h3 := EntryHash("different text")
	if h1 == h3 {
		t.Error("different input should produce different hash")
	}
}

func TestDedup_ImportedCheck(t *testing.T) {
	var s State
	s.ImportedHashes = []string{}

	hash := EntryHash("some entry")
	if s.Imported(hash) {
		t.Error("expected not imported before marking")
	}

	s.MarkImported(hash, "learning")
	if !s.Imported(hash) {
		t.Error("expected imported after marking")
	}
}

func TestDedup_MarkImportedFormat(t *testing.T) {
	var s State
	s.ImportedHashes = []string{}

	hash := EntryHash("test entry")
	s.MarkImported(hash, "decision")

	if len(s.ImportedHashes) != 1 {
		t.Fatalf("expected 1 hash, got %d", len(s.ImportedHashes))
	}

	// Format should be "hash:target:date"
	entry := s.ImportedHashes[0]
	if len(entry) < len(hash)+2 {
		t.Errorf("stored entry too short: %q", entry)
	}
}

func TestLoadState_CorruptJSON(t *testing.T) {
	contextDir := t.TempDir()
	stateDir := filepath.Join(contextDir, config.DirState)
	if mkErr := os.MkdirAll(stateDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	path := filepath.Join(stateDir, config.FileMemoryState)
	if writeErr := os.WriteFile(path, []byte("{corrupt"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	_, loadErr := LoadState(contextDir)
	if loadErr == nil {
		t.Fatal("expected error for corrupt JSON, got nil")
	}
}
