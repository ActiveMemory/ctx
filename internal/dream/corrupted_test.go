//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"os"
	"path/filepath"
	"testing"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// corruptedFixture is a regression corpus modeled on the corrupted-
// artifact appendix of arXiv 2605.12978 ("Useful Memories Become Faulty
// When Continuously Updated by LLMs"): one well-formed proposal among
// several corrupted ones (an unknown status enum, stripped evidence, a
// dropped target). It guards the review/dedup gate against silently
// admitting a faulty artifact.
const corruptedFixture = "testdata/corrupted-2605.12978.json"

// TestCorruptedArtifactGate feeds the corrupted corpus through the real
// reader and the schema gate and asserts that only the one well-formed,
// provenance-bearing proposal survives — the corrupted entries are
// rejected, not silently admitted.
func TestCorruptedArtifactGate(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	runDir := t.TempDir()

	data, readErr := os.ReadFile(corruptedFixture)
	if readErr != nil {
		t.Fatalf("read fixture: %v", readErr)
	}
	if writeErr := ctxIo.SafeWriteFileAtomic(
		filepath.Join(runDir, cfgDream.FileProposals), data, cfgFs.PermSecret,
	); writeErr != nil {
		t.Fatalf("stage proposals: %v", writeErr)
	}

	proposals, err := dream.ReadProposals(runDir)
	if err != nil {
		t.Fatalf("ReadProposals on a well-formed array: %v", err)
	}
	if len(proposals) != 4 {
		t.Fatalf("expected 4 parsed proposals, got %d", len(proposals))
	}

	var survived []string
	for _, p := range proposals {
		if dream.ProposalValid(p) == nil {
			survived = append(survived, p.ID)
		}
	}
	if len(survived) != 1 || survived[0] != "clean-1" {
		t.Fatalf(
			"gate should pass only the provenance-bearing proposal; survived=%v",
			survived,
		)
	}
}

// TestCorruptedArtifactMalformedJSON asserts a structurally corrupt
// proposals file surfaces as an error, not a panic or a silent empty
// result that would let a pass look successful while reading garbage.
func TestCorruptedArtifactMalformedJSON(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	runDir := t.TempDir()

	if writeErr := ctxIo.SafeWriteFileAtomic(
		filepath.Join(runDir, cfgDream.FileProposals),
		[]byte("{ this is not valid json"), cfgFs.PermSecret,
	); writeErr != nil {
		t.Fatalf("stage malformed proposals: %v", writeErr)
	}

	if _, err := dream.ReadProposals(runDir); err == nil {
		t.Fatal("expected an error on malformed proposals.json, got nil")
	}
}

// TestCorruptedArtifactDedup asserts the dedup-against-seen gate drops a
// proposal already recorded in the ledger — so a corrupted artifact that
// was already decided cannot re-surface on a later pass.
func TestCorruptedArtifactDedup(t *testing.T) {
	proposals := []dream.Proposal{
		{ID: "clean-1"},
		{ID: "already-decided"},
	}
	ledger := []dream.LedgerEntry{
		{ProposalID: "already-decided", Decision: cfgDream.DecisionRejected},
	}

	pending := dream.PendingProposals(proposals, ledger)
	if len(pending) != 1 || pending[0].ID != "clean-1" {
		t.Fatalf(
			"dedup-against-seen should drop the already-decided artifact; pending=%v",
			pending,
		)
	}
}
