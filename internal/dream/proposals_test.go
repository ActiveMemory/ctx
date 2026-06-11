//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/dream"
)

// writeRun writes proposals into a per-run dreams/<ts>/ directory and
// returns the dreams dir and run dir.
func writeRun(
	t *testing.T, proposals []dream.Proposal,
) (dreamsDir, runDir string) {
	t.Helper()
	dreamsDir = filepath.Join(t.TempDir(), "dreams")
	stamp := time.Now().UTC().Format(cfgDream.RunTimeLayout)
	runDir = filepath.Join(dreamsDir, stamp)
	payload, marshalErr := json.Marshal(proposals)
	if marshalErr != nil {
		t.Fatalf("marshal proposals: %v", marshalErr)
	}
	if wrErr := writeFixture(
		filepath.Join(runDir, cfgDream.FileProposals), string(payload),
	); wrErr != nil {
		t.Fatalf("write proposals: %v", wrErr)
	}
	return dreamsDir, runDir
}

// sampleProposals returns two well-formed proposals for list tests.
func sampleProposals() []dream.Proposal {
	return []dream.Proposal{
		{
			ID:         "p-1",
			Targets:    []string{"ideas/a.md"},
			Status:     cfgDream.StatusMeritorious,
			Action:     cfgDream.ActionKeep,
			Evidence:   "near-neighbor ideas/b.md (0.4)",
			Confidence: cfgDream.ConfidenceMed,
			Rationale:  "still live",
		},
		{
			ID:         "p-2",
			Targets:    []string{"ideas/c.md"},
			Status:     cfgDream.StatusImplemented,
			Action:     cfgDream.ActionArchive,
			Evidence:   "commit abc123",
			Confidence: cfgDream.ConfidenceHigh,
			Rationale:  "shipped",
		},
	}
}

// TestLatestRunDirPicksNewest returns the lexically greatest run dir
// and skips notebook files.
func TestLatestRunDirPicksNewest(t *testing.T) {
	dreamsDir, runDir := writeRun(t, sampleProposals())

	got, err := dream.LatestRunDir(dreamsDir)
	if err != nil {
		t.Fatalf("LatestRunDir: %v", err)
	}
	if got != runDir {
		t.Fatalf("LatestRunDir = %q, want %q", got, runDir)
	}
}

// TestLatestRunDirEmpty returns empty (not an error) when no runs exist.
func TestLatestRunDirEmpty(t *testing.T) {
	dreamsDir := filepath.Join(t.TempDir(), "dreams")
	if wrErr := writeFixture(
		filepath.Join(dreamsDir, cfgDream.FileLedger), "",
	); wrErr != nil {
		t.Fatalf("seed ledger: %v", wrErr)
	}

	got, err := dream.LatestRunDir(dreamsDir)
	if err != nil {
		t.Fatalf("LatestRunDir: %v", err)
	}
	if got != "" {
		t.Fatalf("LatestRunDir = %q, want empty", got)
	}
}

// TestReadAndFindProposal reads proposals back and finds one by id.
func TestReadAndFindProposal(t *testing.T) {
	_, runDir := writeRun(t, sampleProposals())

	proposals, err := dream.ReadProposals(runDir)
	if err != nil {
		t.Fatalf("ReadProposals: %v", err)
	}
	if len(proposals) != 2 {
		t.Fatalf("ReadProposals len = %d, want 2", len(proposals))
	}
	p, findErr := dream.FindProposal(runDir, "p-2")
	if findErr != nil {
		t.Fatalf("FindProposal: %v", findErr)
	}
	if p.Action != cfgDream.ActionArchive {
		t.Fatalf("FindProposal p-2 action = %q, want archive", p.Action)
	}
	if _, missErr := dream.FindProposal(runDir, "nope"); missErr == nil {
		t.Fatal("FindProposal must error on a missing id")
	}
}

// TestPendingProposalsDedupAgainstSeen filters out proposals already
// recorded in the ledger.
func TestPendingProposalsDedupAgainstSeen(t *testing.T) {
	dreamsDir, runDir := writeRun(t, sampleProposals())

	if appendErr := dream.AppendLedger(dreamsDir, dream.LedgerEntry{
		ProposalID: "p-1",
		Decision:   cfgDream.DecisionRejected,
		Action:     cfgDream.ActionKeep,
		At:         time.Now().UTC(),
	}); appendErr != nil {
		t.Fatalf("AppendLedger: %v", appendErr)
	}

	proposals, _ := dream.ReadProposals(runDir)
	ledger, _ := dream.ReadLedger(dreamsDir)
	pending := dream.PendingProposals(proposals, ledger)

	if len(pending) != 1 {
		t.Fatalf("pending len = %d, want 1", len(pending))
	}
	if pending[0].ID != "p-2" {
		t.Fatalf("pending[0] = %q, want p-2 (p-1 was seen)", pending[0].ID)
	}
}
