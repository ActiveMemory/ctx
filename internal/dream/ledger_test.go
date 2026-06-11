//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"path/filepath"
	"testing"
	"time"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/dream"
)

// TestLedgerAppendReadback appends entries and reads them back in order.
func TestLedgerAppendReadback(t *testing.T) {
	dreamsDir := filepath.Join(t.TempDir(), "dreams")
	at := time.Date(2026, 6, 7, 2, 30, 0, 0, time.UTC)

	entries := []dream.LedgerEntry{
		{
			ProposalID: "p1",
			Decision:   cfgDream.DecisionAccepted,
			Action:     cfgDream.ActionArchive,
			At:         at,
		},
		{
			ProposalID: "p2",
			Decision:   cfgDream.DecisionRejected,
			Action:     cfgDream.ActionKeep,
			At:         at,
			Note:       "not relevant anymore",
		},
	}
	for _, e := range entries {
		if err := dream.AppendLedger(dreamsDir, e); err != nil {
			t.Fatalf("AppendLedger: %v", err)
		}
	}

	got, err := dream.ReadLedger(dreamsDir)
	if err != nil {
		t.Fatalf("ReadLedger: %v", err)
	}
	if len(got) != len(entries) {
		t.Fatalf("read %d entries, want %d", len(got), len(entries))
	}
	if got[0].ProposalID != "p1" || got[1].ProposalID != "p2" {
		t.Fatalf("append order not preserved: %+v", got)
	}
	if got[1].Note != "not relevant anymore" {
		t.Fatalf("note not preserved: %q", got[1].Note)
	}
}

// TestLedgerSeenDedup verifies the dedup-against-seen signal: a recorded
// proposal (including a rejection) reports as seen; an unrecorded one
// does not.
func TestLedgerSeenDedup(t *testing.T) {
	dreamsDir := filepath.Join(t.TempDir(), "dreams")
	if err := dream.AppendLedger(dreamsDir, dream.LedgerEntry{
		ProposalID: "rejected-1",
		Decision:   cfgDream.DecisionRejected,
		Action:     cfgDream.ActionKeep,
		At:         time.Now().UTC(),
	}); err != nil {
		t.Fatalf("AppendLedger: %v", err)
	}

	entries, err := dream.ReadLedger(dreamsDir)
	if err != nil {
		t.Fatalf("ReadLedger: %v", err)
	}
	if !dream.Seen(entries, "rejected-1") {
		t.Fatal("rejected proposal must report as seen")
	}
	if dream.Seen(entries, "never-surfaced") {
		t.Fatal("unrecorded proposal must not report as seen")
	}
}

// TestLedgerReadMissing returns an empty slice when no ledger exists.
func TestLedgerReadMissing(t *testing.T) {
	dreamsDir := filepath.Join(t.TempDir(), "dreams")
	got, err := dream.ReadLedger(dreamsDir)
	if err != nil {
		t.Fatalf("ReadLedger: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want empty, got %d", len(got))
	}
}
