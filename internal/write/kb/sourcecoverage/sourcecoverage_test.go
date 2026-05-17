//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage_test

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errKbSC "github.com/ActiveMemory/ctx/internal/err/kb/sourcecoverage"
	sc "github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage"
)

func TestValidTransition_AllowedSubset(t *testing.T) {
	cases := []struct {
		from, to string
		want     bool
	}{
		{cfgKB.StateDiscovered, cfgKB.StateAdmitted, true},
		{cfgKB.StateDiscovered, cfgKB.StateSkipped, true},
		{cfgKB.StateAdmitted, cfgKB.StateHighlightsExtracted, true},
		{cfgKB.StateAdmitted, cfgKB.StateComprehensive, true},
		{cfgKB.StateHighlightsExtracted, cfgKB.StateComprehensive, true},
		{cfgKB.StateTopicPageDrafted, cfgKB.StateComprehensive, true},
		// Same-state always allowed (idempotent touch).
		{cfgKB.StateAdmitted, cfgKB.StateAdmitted, true},
		// Illegal: backwards
		{cfgKB.StateComprehensive, cfgKB.StateHighlightsExtracted, false},
		// Illegal: skipped → anything
		{cfgKB.StateSkipped, cfgKB.StateAdmitted, false},
		// Illegal: discovered → comprehensive (must pass through admitted)
		{cfgKB.StateDiscovered, cfgKB.StateComprehensive, false},
	}
	for _, c := range cases {
		if got := sc.ValidTransition(c.from, c.to); got != c.want {
			t.Errorf("%s → %s: want %v; got %v", c.from, c.to, c.want, got)
		}
	}
}

func TestAdvance_NewSourceAtDiscovered(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source-coverage.md")
	row := sc.Row{
		Source: "CURSOR-HOOKS",
		Topic:  "cursor-hooks",
		State:  cfgKB.StateDiscovered,
	}
	if err := sc.Advance(path, row); err != nil {
		t.Fatalf("Advance: %v", err)
	}
	rows, err := sc.Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("count: want 1; got %d", len(rows))
	}
	if rows[0].Source != "CURSOR-HOOKS" {
		t.Errorf("Source: got %q", rows[0].Source)
	}
	if rows[0].State != cfgKB.StateDiscovered {
		t.Errorf("State: got %q", rows[0].State)
	}
}

func TestAdvance_NewSourceAtAdmittedOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source-coverage.md")
	row := sc.Row{
		Source: "FOO",
		Topic:  "n/a",
		State:  cfgKB.StateAdmitted,
	}
	if err := sc.Advance(path, row); err != nil {
		t.Fatalf("Advance: %v", err)
	}
}

func TestAdvance_NewSourceAtInvalidStateRejected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source-coverage.md")
	row := sc.Row{
		Source: "FOO",
		Topic:  "n/a",
		State:  cfgKB.StateComprehensive,
	}
	err := sc.Advance(path, row)
	if !errors.Is(err, errKbSC.ErrUnknownSource) {
		t.Fatalf("want ErrUnknownSource; got %v", err)
	}
}

func TestAdvance_IllegalTransitionRejected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source-coverage.md")
	// First, admit it.
	if err := sc.Advance(path, sc.Row{
		Source: "BAR", Topic: "n/a", State: cfgKB.StateAdmitted,
	}); err != nil {
		t.Fatal(err)
	}
	// Try to demote from admitted → discovered (illegal).
	err := sc.Advance(path, sc.Row{
		Source: "BAR", Topic: "n/a", State: cfgKB.StateDiscovered,
	})
	if !errors.Is(err, errKbSC.ErrIllegalTransition) {
		t.Fatalf("want ErrIllegalTransition; got %v", err)
	}
}

func TestAdvance_LegalTransitionAdvancesRow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source-coverage.md")
	for _, state := range []string{
		cfgKB.StateAdmitted,
		cfgKB.StateHighlightsExtracted,
		cfgKB.StateTopicPageDrafted,
		cfgKB.StateComprehensive,
	} {
		if err := sc.Advance(path, sc.Row{
			Source:  "BAZ",
			Topic:   "baz-topic",
			State:   state,
			Updated: time.Now(),
		}); err != nil {
			t.Fatalf("advance to %s: %v", state, err)
		}
	}
	rows, err := sc.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row; got %d", len(rows))
	}
	if rows[0].State != cfgKB.StateComprehensive {
		t.Errorf("final state: want %q; got %q",
			cfgKB.StateComprehensive, rows[0].State)
	}
}

func TestRead_NonExistentFileReturnsEmpty(t *testing.T) {
	rows, err := sc.Read(filepath.Join(t.TempDir(), "nope.md"))
	if err != nil {
		t.Fatalf("want nil; got %v", err)
	}
	if rows != nil {
		t.Errorf("want nil slice; got %v", rows)
	}
}
