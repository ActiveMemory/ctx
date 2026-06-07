//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package review_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/dream/core/review"
	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	engine "github.com/ActiveMemory/ctx/internal/dream"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// seedProject lays out a project root with .context and a dream run
// dir holding proposals, then chdir's in for the cwd-anchored resolver.
func seedProject(
	t *testing.T, proposals []engine.Proposal,
) (dreamsDir string) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(
		filepath.Join(root, cfgDir.Context), 0o755,
	); err != nil {
		t.Fatalf("mkdir .context: %v", err)
	}
	dreamsDir = filepath.Join(root, cfgDir.Dreams)
	stamp := time.Now().UTC().Format(cfgDream.RunTimeLayout)
	runDir := filepath.Join(dreamsDir, stamp)
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run: %v", err)
	}
	payload, _ := json.Marshal(proposals)
	if err := os.WriteFile(
		filepath.Join(runDir, cfgDream.FileProposals), payload, 0o600,
	); err != nil {
		t.Fatalf("write proposals: %v", err)
	}
	t.Chdir(root)
	rc.Reset()
	t.Cleanup(rc.Reset)
	return dreamsDir
}

func captureCmd() (*cobra.Command, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	c := &cobra.Command{}
	c.SetOut(buf)
	c.SetErr(buf)
	return c, buf
}

// TestReviewListsPending lists the pending proposal and shows its id.
func TestReviewListsPending(t *testing.T) {
	seedProject(t, []engine.Proposal{
		{
			ID:         "p-1",
			Targets:    []string{"ideas/a.md"},
			Status:     cfgDream.StatusMeritorious,
			Action:     cfgDream.ActionKeep,
			Evidence:   "live",
			Confidence: cfgDream.ConfidenceMed,
			Rationale:  "keep it",
		},
	})

	c, buf := captureCmd()
	if err := review.Run(c); err != nil {
		t.Fatalf("review.Run: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("p-1")) {
		t.Fatalf("review output missing p-1:\n%s", buf.String())
	}
}

// TestReviewFiltersSeen drops a proposal already recorded in the
// ledger (dedup-against-seen).
func TestReviewFiltersSeen(t *testing.T) {
	dreamsDir := seedProject(t, []engine.Proposal{
		{
			ID:         "p-1",
			Targets:    []string{"ideas/a.md"},
			Status:     cfgDream.StatusMeritorious,
			Action:     cfgDream.ActionKeep,
			Evidence:   "live",
			Confidence: cfgDream.ConfidenceMed,
			Rationale:  "keep it",
		},
	})
	if err := engine.AppendLedger(dreamsDir, engine.LedgerEntry{
		ProposalID: "p-1",
		Decision:   cfgDream.DecisionRejected,
		Action:     cfgDream.ActionKeep,
		At:         time.Now().UTC(),
	}); err != nil {
		t.Fatalf("AppendLedger: %v", err)
	}

	c, buf := captureCmd()
	if err := review.Run(c); err != nil {
		t.Fatalf("review.Run: %v", err)
	}
	if bytes.Contains(buf.Bytes(), []byte("p-1")) {
		t.Fatalf("seen proposal should be filtered:\n%s", buf.String())
	}
}

// TestReviewNoRuns prints the no-pending message when no run exists.
func TestReviewNoRuns(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(
		filepath.Join(root, cfgDir.Context), 0o755,
	); err != nil {
		t.Fatalf("mkdir .context: %v", err)
	}
	t.Chdir(root)
	rc.Reset()
	t.Cleanup(rc.Reset)

	c, buf := captureCmd()
	if err := review.Run(c); err != nil {
		t.Fatalf("review.Run: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("review with no runs should print the no-pending message")
	}
}
