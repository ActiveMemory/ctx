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
	"time"
)

func TestRecord(t *testing.T) {
	stateDir := t.TempDir()

	beforeRecord := time.Now().UTC().Truncate(time.Second)
	if err := Record("T-1", stateDir); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	entries, readErr := ReadPending(stateDir)
	if readErr != nil {
		t.Fatalf("ReadPending() error: %v", readErr)
	}
	if len(entries) != 1 {
		t.Fatalf("ReadPending() returned %d entries, want 1", len(entries))
	}
	if entries[0].Ref != "T-1" {
		t.Errorf("Ref = %q, want %q", entries[0].Ref, "T-1")
	}
	if entries[0].Timestamp.Before(beforeRecord) {
		t.Errorf("Timestamp %v is before record time %v", entries[0].Timestamp, beforeRecord)
	}
}

func TestRecordMultiple(t *testing.T) {
	stateDir := t.TempDir()

	refs := []string{"T-1", "D-2", "L-3"}
	for _, ref := range refs {
		if err := Record(ref, stateDir); err != nil {
			t.Fatalf("Record(%q) error: %v", ref, err)
		}
	}

	pendingPath := filepath.Join(stateDir, pendingFile)
	data, readErr := os.ReadFile(pendingPath) //nolint:gosec // test file
	if readErr != nil {
		t.Fatalf("ReadFile() error: %v", readErr)
	}

	// Count newlines — each JSONL entry ends with one.
	lineCount := 0
	for _, b := range data {
		if b == '\n' {
			lineCount++
		}
	}
	if lineCount != 3 {
		t.Errorf("got %d lines in file, want 3", lineCount)
	}
}

func TestReadPending(t *testing.T) {
	stateDir := t.TempDir()

	if err := Record("D-1", stateDir); err != nil {
		t.Fatalf("Record(D-1) error: %v", err)
	}
	if err := Record("L-5", stateDir); err != nil {
		t.Fatalf("Record(L-5) error: %v", err)
	}

	entries, readErr := ReadPending(stateDir)
	if readErr != nil {
		t.Fatalf("ReadPending() error: %v", readErr)
	}
	if len(entries) != 2 {
		t.Fatalf("ReadPending() returned %d entries, want 2", len(entries))
	}
	if entries[0].Ref != "D-1" {
		t.Errorf("entries[0].Ref = %q, want %q", entries[0].Ref, "D-1")
	}
	if entries[1].Ref != "L-5" {
		t.Errorf("entries[1].Ref = %q, want %q", entries[1].Ref, "L-5")
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("entries[0].Timestamp is zero")
	}
	if entries[1].Timestamp.IsZero() {
		t.Error("entries[1].Timestamp is zero")
	}
}

func TestReadPendingEmpty(t *testing.T) {
	// Use a directory that does not contain the pending file.
	stateDir := t.TempDir()

	entries, readErr := ReadPending(stateDir)
	if readErr != nil {
		t.Fatalf("ReadPending() on missing file returned error: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf("ReadPending() returned %d entries, want 0", len(entries))
	}
}

func TestTruncatePending(t *testing.T) {
	stateDir := t.TempDir()

	if err := Record("T-2", stateDir); err != nil {
		t.Fatalf("Record() error: %v", err)
	}
	if err := Record("T-3", stateDir); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	if err := TruncatePending(stateDir); err != nil {
		t.Fatalf("TruncatePending() error: %v", err)
	}

	entries, readErr := ReadPending(stateDir)
	if readErr != nil {
		t.Fatalf("ReadPending() after truncate error: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf("ReadPending() after truncate returned %d entries, want 0", len(entries))
	}
}
