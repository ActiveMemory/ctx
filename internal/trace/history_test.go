//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"testing"
	"time"
)

func TestWriteHistory(t *testing.T) {
	traceDir := t.TempDir()

	entry := HistoryEntry{
		Commit:  "abc1234",
		Refs:    []string{"T-1", "D-2"},
		Message: "feat: add thing",
	}

	beforeWrite := time.Now().UTC().Truncate(time.Second)
	if err := WriteHistory(entry, traceDir); err != nil {
		t.Fatalf("WriteHistory() error: %v", err)
	}

	entries, readErr := ReadHistory(traceDir)
	if readErr != nil {
		t.Fatalf("ReadHistory() error: %v", readErr)
	}
	if len(entries) != 1 {
		t.Fatalf("ReadHistory() returned %d entries, want 1", len(entries))
	}

	got := entries[0]
	if got.Commit != entry.Commit {
		t.Errorf("Commit = %q, want %q", got.Commit, entry.Commit)
	}
	if got.Message != entry.Message {
		t.Errorf("Message = %q, want %q", got.Message, entry.Message)
	}
	if len(got.Refs) != 2 {
		t.Errorf("Refs len = %d, want 2", len(got.Refs))
	}
	if got.Timestamp.Before(beforeWrite) {
		t.Errorf("Timestamp %v is before write time %v", got.Timestamp, beforeWrite)
	}
}

func TestReadHistoryEmpty(t *testing.T) {
	traceDir := t.TempDir()

	entries, readErr := ReadHistory(traceDir)
	if readErr != nil {
		t.Fatalf("ReadHistory() on missing file returned error: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf("ReadHistory() returned %d entries, want 0", len(entries))
	}
}

func TestReadHistoryForCommit(t *testing.T) {
	traceDir := t.TempDir()

	e1 := HistoryEntry{
		Commit:  "deadbeef1234",
		Refs:    []string{"T-1"},
		Message: "first commit",
	}
	e2 := HistoryEntry{
		Commit:  "cafebabe5678",
		Refs:    []string{"D-3"},
		Message: "second commit",
	}

	if err := WriteHistory(e1, traceDir); err != nil {
		t.Fatalf("WriteHistory(e1) error: %v", err)
	}
	if err := WriteHistory(e2, traceDir); err != nil {
		t.Fatalf("WriteHistory(e2) error: %v", err)
	}

	// Find by full hash.
	got, ok := ReadHistoryForCommit("deadbeef1234", traceDir)
	if !ok {
		t.Fatal("ReadHistoryForCommit(full hash) returned false, want true")
	}
	if got.Commit != e1.Commit {
		t.Errorf("Commit = %q, want %q", got.Commit, e1.Commit)
	}

	// Find by short hash (prefix of stored hash).
	got, ok = ReadHistoryForCommit("deadbeef", traceDir)
	if !ok {
		t.Fatal("ReadHistoryForCommit(short hash) returned false, want true")
	}
	if got.Commit != e1.Commit {
		t.Errorf("Commit = %q, want %q", got.Commit, e1.Commit)
	}

	// Missing hash returns false.
	_, ok = ReadHistoryForCommit("0000000", traceDir)
	if ok {
		t.Error("ReadHistoryForCommit(missing) returned true, want false")
	}
}

func TestWriteOverride(t *testing.T) {
	traceDir := t.TempDir()

	entry := OverrideEntry{
		Commit: "abc1234",
		Refs:   []string{"L-7", "T-2"},
	}

	beforeWrite := time.Now().UTC().Truncate(time.Second)
	if err := WriteOverride(entry, traceDir); err != nil {
		t.Fatalf("WriteOverride() error: %v", err)
	}

	entries, readErr := ReadOverrides(traceDir)
	if readErr != nil {
		t.Fatalf("ReadOverrides() error: %v", readErr)
	}
	if len(entries) != 1 {
		t.Fatalf("ReadOverrides() returned %d entries, want 1", len(entries))
	}

	got := entries[0]
	if got.Commit != entry.Commit {
		t.Errorf("Commit = %q, want %q", got.Commit, entry.Commit)
	}
	if len(got.Refs) != 2 {
		t.Errorf("Refs len = %d, want 2", len(got.Refs))
	}
	if got.Timestamp.Before(beforeWrite) {
		t.Errorf("Timestamp %v is before write time %v", got.Timestamp, beforeWrite)
	}
}

func TestReadOverridesEmpty(t *testing.T) {
	traceDir := t.TempDir()

	entries, readErr := ReadOverrides(traceDir)
	if readErr != nil {
		t.Fatalf("ReadOverrides() on missing file returned error: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf("ReadOverrides() returned %d entries, want 0", len(entries))
	}
}

func TestReadOverridesForCommit(t *testing.T) {
	traceDir := t.TempDir()

	// Two overrides for the same commit, one for a different commit.
	o1 := OverrideEntry{
		Commit: "deadbeef1234",
		Refs:   []string{"T-1", "D-2"},
	}
	o2 := OverrideEntry{
		Commit: "deadbeef1234",
		Refs:   []string{"L-5"},
	}
	o3 := OverrideEntry{
		Commit: "cafebabe5678",
		Refs:   []string{"T-9"},
	}

	for _, o := range []OverrideEntry{o1, o2, o3} {
		if err := WriteOverride(o, traceDir); err != nil {
			t.Fatalf("WriteOverride() error: %v", err)
		}
	}

	refs := ReadOverridesForCommit("deadbeef1234", traceDir)
	// o1 has 2 refs, o2 has 1 ref — total 3 for deadbeef1234.
	if len(refs) != 3 {
		t.Errorf("ReadOverridesForCommit() returned %d refs, want 3", len(refs))
	}

	// Different commit should return its own refs.
	otherRefs := ReadOverridesForCommit("cafebabe5678", traceDir)
	if len(otherRefs) != 1 {
		t.Errorf("ReadOverridesForCommit(other) returned %d refs, want 1", len(otherRefs))
	}

	// Non-existent commit returns empty slice.
	noRefs := ReadOverridesForCommit("0000000", traceDir)
	if len(noRefs) != 0 {
		t.Errorf("ReadOverridesForCommit(missing) returned %d refs, want 0", len(noRefs))
	}
}
