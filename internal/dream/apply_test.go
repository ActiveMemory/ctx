//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/dream"
)

// fixtureRepo creates a git repo with dreams/ and ideas/ gitignored, an
// ideas/note.md source, and a dreams/ notebook dir. It returns the root
// and dreams dir so the appliers' guards pass.
func fixtureRepo(t *testing.T) (root, dreamsDir string) {
	t.Helper()
	root = t.TempDir()
	gitInit(t, root)
	dreamsDir = filepath.Join(root, cfgDir.Dreams)
	if mkErr := os.MkdirAll(dreamsDir, 0o755); mkErr != nil {
		t.Fatalf("mkdir dreams: %v", mkErr)
	}
	if wrErr := writeFixture(
		filepath.Join(root, "ideas", "note.md"), "an idea\n",
	); wrErr != nil {
		t.Fatalf("write idea: %v", wrErr)
	}
	return root, dreamsDir
}

// ledgerHas reports whether the ledger records a decision for id with
// the given decision and action.
func ledgerHas(
	t *testing.T, dreamsDir, id string,
	decision cfgDream.Decision, action cfgDream.ProposalAction,
) bool {
	t.Helper()
	entries, readErr := dream.ReadLedger(dreamsDir)
	if readErr != nil {
		t.Fatalf("read ledger: %v", readErr)
	}
	for _, e := range entries {
		if e.ProposalID == id &&
			e.Decision == decision && e.Action == action {
			return true
		}
	}
	return false
}

// TestAcceptArchiveMoves accepts an archive proposal: the idea moves to
// ideas/done/ and the ledger records an accepted archive.
func TestAcceptArchiveMoves(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	p := dream.Proposal{
		ID:         "d-001",
		Targets:    []string{"ideas/note.md"},
		Status:     cfgDream.StatusImplemented,
		Action:     cfgDream.ActionArchive,
		Confidence: cfgDream.ConfidenceHigh,
	}

	res, err := dream.Accept(root, dreamsDir, p, "")
	if err != nil {
		t.Fatalf("Accept archive: %v", err)
	}
	if !res.Performed || res.Generative {
		t.Fatalf("archive must be performed mechanically: %+v", res)
	}
	if _, statErr := os.Stat(
		filepath.Join(root, "ideas", "note.md"),
	); !os.IsNotExist(statErr) {
		t.Fatal("source should have been moved out of ideas/")
	}
	if _, statErr := os.Stat(
		filepath.Join(root, "ideas", "done", "note.md"),
	); statErr != nil {
		t.Fatalf("source should be in ideas/done/: %v", statErr)
	}
	if !ledgerHas(
		t, dreamsDir, "d-001",
		cfgDream.DecisionAccepted, cfgDream.ActionArchive,
	) {
		t.Fatal("ledger missing accepted archive")
	}
}

// TestAcceptMarkBlogTagsInPlace accepts a mark-blog proposal: the idea
// stays in place, gains the blog marker, and the ledger records it.
func TestAcceptMarkBlogTagsInPlace(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	p := dream.Proposal{
		ID:         "d-002",
		Targets:    []string{"ideas/note.md"},
		Status:     cfgDream.StatusBlogCandidate,
		Action:     cfgDream.ActionMarkBlog,
		Confidence: cfgDream.ConfidenceMed,
	}

	if _, err := dream.Accept(root, dreamsDir, p, ""); err != nil {
		t.Fatalf("Accept mark-blog: %v", err)
	}
	data, readErr := os.ReadFile(
		filepath.Join(root, "ideas", "note.md"),
	)
	if readErr != nil {
		t.Fatalf("read tagged idea: %v", readErr)
	}
	if !strings.Contains(string(data), cfgDream.BlogMarker) {
		t.Fatal("mark-blog should append the blog marker in place")
	}
	if !ledgerHas(
		t, dreamsDir, "d-002",
		cfgDream.DecisionAccepted, cfgDream.ActionMarkBlog,
	) {
		t.Fatal("ledger missing accepted mark-blog")
	}
}

// TestAcceptKeepNoMutation accepts a keep proposal: no file changes,
// ledger records it.
func TestAcceptKeepNoMutation(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	p := dream.Proposal{
		ID:         "d-003",
		Targets:    []string{"ideas/note.md"},
		Status:     cfgDream.StatusMeritorious,
		Action:     cfgDream.ActionKeep,
		Confidence: cfgDream.ConfidenceLow,
	}

	if _, err := dream.Accept(root, dreamsDir, p, ""); err != nil {
		t.Fatalf("Accept keep: %v", err)
	}
	if _, statErr := os.Stat(
		filepath.Join(root, "ideas", "note.md"),
	); statErr != nil {
		t.Fatalf("keep must not move the source: %v", statErr)
	}
	if !ledgerHas(
		t, dreamsDir, "d-003",
		cfgDream.DecisionAccepted, cfgDream.ActionKeep,
	) {
		t.Fatal("ledger missing accepted keep")
	}
}

// TestAcceptPromoteIsGenerative accepts a promote: no mutation here, the
// result is generative, and the ledger records the accepted intent.
func TestAcceptPromoteIsGenerative(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	p := dream.Proposal{
		ID:         "d-004",
		Targets:    []string{"ideas/note.md"},
		Status:     cfgDream.StatusMeritorious,
		Action:     cfgDream.ActionPromote,
		Confidence: cfgDream.ConfidenceHigh,
	}

	res, err := dream.Accept(root, dreamsDir, p, "")
	if err != nil {
		t.Fatalf("Accept promote: %v", err)
	}
	if !res.Generative || res.Performed {
		t.Fatalf("promote must be generative, not performed: %+v", res)
	}
	if _, statErr := os.Stat(
		filepath.Join(root, "ideas", "note.md"),
	); statErr != nil {
		t.Fatal("promote must not move the source from the CLI")
	}
	if !ledgerHas(
		t, dreamsDir, "d-004",
		cfgDream.DecisionAccepted, cfgDream.ActionPromote,
	) {
		t.Fatal("ledger missing accepted promote intent")
	}
}

// TestRejectRecordsNoMutation rejects a proposal: no file changes, the
// ledger records a rejection.
func TestRejectRecordsNoMutation(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	p := dream.Proposal{
		ID:         "d-005",
		Targets:    []string{"ideas/note.md"},
		Status:     cfgDream.StatusSidenote,
		Action:     cfgDream.ActionArchive,
		Confidence: cfgDream.ConfidenceLow,
	}

	res, err := dream.Reject(dreamsDir, p, "not now")
	if err != nil {
		t.Fatalf("Reject: %v", err)
	}
	if !res.Performed {
		t.Fatal("reject is a recorded, performed disposition")
	}
	if _, statErr := os.Stat(
		filepath.Join(root, "ideas", "note.md"),
	); statErr != nil {
		t.Fatal("reject must not move the source")
	}
	if !ledgerHas(
		t, dreamsDir, "d-005",
		cfgDream.DecisionRejected, cfgDream.ActionArchive,
	) {
		t.Fatal("ledger missing rejection")
	}
}

// TestAmendAppliesDifferentAction amends an archive proposal to keep:
// the source stays and the ledger records an amended keep.
func TestAmendAppliesDifferentAction(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	p := dream.Proposal{
		ID:         "d-006",
		Targets:    []string{"ideas/note.md"},
		Status:     cfgDream.StatusMeritorious,
		Action:     cfgDream.ActionArchive,
		Confidence: cfgDream.ConfidenceMed,
	}

	res, err := dream.Amend(root, dreamsDir, p, cfgDream.ActionKeep, "")
	if err != nil {
		t.Fatalf("Amend: %v", err)
	}
	if res.Action != cfgDream.ActionKeep {
		t.Fatalf("amended action = %q, want keep", res.Action)
	}
	if _, statErr := os.Stat(
		filepath.Join(root, "ideas", "note.md"),
	); statErr != nil {
		t.Fatal("amend-to-keep must not move the source")
	}
	if !ledgerHas(
		t, dreamsDir, "d-006",
		cfgDream.DecisionAmended, cfgDream.ActionKeep,
	) {
		t.Fatal("ledger missing amended keep")
	}
}

// TestAmendUnknownActionRefused rejects an unrecognized amend action
// before any mutation or ledger write.
func TestAmendUnknownActionRefused(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	p := dream.Proposal{
		ID:         "d-007",
		Targets:    []string{"ideas/note.md"},
		Status:     cfgDream.StatusMeritorious,
		Action:     cfgDream.ActionKeep,
		Confidence: cfgDream.ConfidenceMed,
	}

	if _, err := dream.Amend(
		root, dreamsDir, p, "obliterate", "",
	); err == nil {
		t.Fatal("unknown action must be refused")
	}
	entries, _ := dream.ReadLedger(dreamsDir)
	if len(entries) != 0 {
		t.Fatal("unknown action must not write a ledger entry")
	}
}

// TestAcceptGuardRefusesTrackedTarget refuses an archive whose target
// resolves to a tracked path (the don't-leak guard).
func TestAcceptGuardRefusesTrackedTarget(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)
	// A tracked top-level file: .gitignore is committed by gitInit.
	p := dream.Proposal{
		ID:         "d-008",
		Targets:    []string{".gitignore"},
		Status:     cfgDream.StatusImplemented,
		Action:     cfgDream.ActionMarkBlog,
		Confidence: cfgDream.ConfidenceHigh,
	}

	if _, err := dream.Accept(root, dreamsDir, p, ""); err == nil {
		t.Fatal("write to a tracked path must be refused by the guard")
	}
	entries, _ := dream.ReadLedger(dreamsDir)
	if len(entries) != 0 {
		t.Fatal("a guard refusal must not write a ledger entry")
	}
}

// TestBackupSnapshotsSource backs up a source into dreams/ before a
// destructive mutation, leaving the original intact.
func TestBackupSnapshotsSource(t *testing.T) {
	root, dreamsDir := fixtureRepo(t)

	if err := dream.Backup(root, dreamsDir, "ideas/note.md"); err != nil {
		t.Fatalf("Backup: %v", err)
	}
	if _, statErr := os.Stat(
		filepath.Join(dreamsDir, "note.md"+cfgDream.BackupSuffix),
	); statErr != nil {
		t.Fatalf("backup snapshot missing: %v", statErr)
	}
	if _, statErr := os.Stat(
		filepath.Join(root, "ideas", "note.md"),
	); statErr != nil {
		t.Fatal("backup must not remove the original")
	}
}
