//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errKbEvidence "github.com/ActiveMemory/ctx/internal/err/kb/evidence"
	"github.com/ActiveMemory/ctx/internal/write/kb/evidence"
)

func TestAppend_AllocatesNextSequentialID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "evidence-index.md")
	row, err := evidence.Append(path, evidence.Row{
		Claim:      "first claim",
		SourceID:   "SRC-A",
		Locator:    "L:1",
		Confidence: cfgKB.ConfidenceMedium,
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if row.ID != "EV-001" {
		t.Errorf("first ID: want EV-001; got %s", row.ID)
	}

	row2, err := evidence.Append(path, evidence.Row{
		Claim:      "second claim",
		SourceID:   "SRC-B",
		Locator:    "ts=00:01:23",
		Confidence: cfgKB.ConfidenceLow,
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if row2.ID != "EV-002" {
		t.Errorf("second ID: want EV-002; got %s", row2.ID)
	}
}

func TestAppend_RespectsExplicitID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "evidence-index.md")
	if _, err := evidence.Append(path, evidence.Row{
		Claim:      "first",
		Confidence: cfgKB.ConfidenceHigh,
	}); err != nil {
		t.Fatal(err)
	}
	// Append an explicit ID well above the auto-incremented one.
	got, err := evidence.Append(path, evidence.Row{
		ID:         "EV-099",
		Claim:      "explicit",
		Confidence: cfgKB.ConfidenceHigh,
	})
	if err != nil {
		t.Fatalf("Append explicit: %v", err)
	}
	if got.ID != "EV-099" {
		t.Errorf("explicit ID preserved: got %s", got.ID)
	}
	// Next auto-allocated ID should be EV-100 (after the max).
	next, err := evidence.Append(path, evidence.Row{
		Claim:      "after gap",
		Confidence: cfgKB.ConfidenceMedium,
	})
	if err != nil {
		t.Fatal(err)
	}
	if next.ID != "EV-100" {
		t.Errorf("max+1 after gap: want EV-100; got %s", next.ID)
	}
}

func TestAppend_RejectsDuplicateID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "evidence-index.md")
	if _, err := evidence.Append(path, evidence.Row{
		Claim:      "first",
		Confidence: cfgKB.ConfidenceHigh,
	}); err != nil {
		t.Fatal(err)
	}
	_, err := evidence.Append(path, evidence.Row{
		ID:         "EV-001",
		Claim:      "dup",
		Confidence: cfgKB.ConfidenceHigh,
	})
	if !errors.Is(err, errKbEvidence.ErrDuplicateID) {
		t.Fatalf("want ErrDuplicateID; got %v", err)
	}
}

func TestAppend_RejectsInvalidBand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "evidence-index.md")
	_, err := evidence.Append(path, evidence.Row{
		Claim:      "x",
		Confidence: "unsure",
	})
	if !errors.Is(err, errKbEvidence.ErrInvalidBand) {
		t.Fatalf("want ErrInvalidBand; got %v", err)
	}
}

func TestAppend_HeaderWrittenOnFirstAppend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "evidence-index.md")
	if _, err := evidence.Append(path, evidence.Row{
		Claim:      "x",
		Confidence: cfgKB.ConfidenceLow,
	}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "# Evidence index") {
		t.Error("header missing")
	}
	if !strings.Contains(string(raw), "| ID |") {
		t.Error("table header missing")
	}
	if !strings.Contains(string(raw), "EV-001") {
		t.Error("first row missing")
	}
}

func TestAppend_EvidenceOnlyTagPreservedVerbatim(t *testing.T) {
	path := filepath.Join(t.TempDir(), "evidence-index.md")
	row, err := evidence.Append(path, evidence.Row{
		Claim:      "from evidence-only mode",
		Confidence: cfgKB.ConfidenceLow,
		Tags:       []string{"cursor", "hooks", cfgKB.EvidenceOnlyTag},
	})
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), cfgKB.EvidenceOnlyTag) {
		t.Errorf("evidence-only tag missing from file")
	}
	if len(row.Tags) != 3 {
		t.Errorf("tags lost in round-trip: %v", row.Tags)
	}
}
