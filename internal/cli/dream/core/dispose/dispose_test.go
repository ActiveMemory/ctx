//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dispose_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/dream/core/dispose"
	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	engine "github.com/ActiveMemory/ctx/internal/dream"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// setupProject builds a project root with .context, a git repo that
// gitignores dreams/ and ideas/, an ideas/note.md source, and a dream
// run dir holding the given proposals. It chdir's into the root so the
// cwd-anchored resolver finds it.
func setupProject(
	t *testing.T, proposals []engine.Proposal,
) (root, dreamsDir string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not on PATH: %v", err)
	}
	root = t.TempDir()
	mustMkdir(t, filepath.Join(root, cfgDir.Context))
	gitInit(t, root)
	mustWrite(t, filepath.Join(root, cfgDir.Ideas, "note.md"), "idea\n")

	dreamsDir = filepath.Join(root, cfgDir.Dreams)
	stamp := time.Now().UTC().Format(cfgDream.RunTimeLayout)
	runDir := filepath.Join(dreamsDir, stamp)
	payload, _ := json.Marshal(proposals)
	mustWrite(
		t, filepath.Join(runDir, cfgDream.FileProposals), string(payload),
	)

	t.Chdir(root)
	rc.Reset()
	t.Cleanup(rc.Reset)
	return root, dreamsDir
}

// TestAcceptArchiveThroughCLI accepts an archive proposal end-to-end:
// the idea moves to ideas/done/ and the ledger records it.
func TestAcceptArchiveThroughCLI(t *testing.T) {
	root, dreamsDir := setupProject(t, []engine.Proposal{
		{
			ID:         "d-1",
			Targets:    []string{filepath.Join(cfgDir.Ideas, "note.md")},
			Status:     cfgDream.StatusImplemented,
			Action:     cfgDream.ActionArchive,
			Confidence: cfgDream.ConfidenceHigh,
		},
	})

	c := newCmd()
	if err := dispose.Accept(c, "d-1", ""); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if _, statErr := os.Stat(
		filepath.Join(root, cfgDir.Ideas, cfgDir.Done, "note.md"),
	); statErr != nil {
		t.Fatalf("idea should be archived: %v", statErr)
	}
	if !ledgerHas(t, dreamsDir, "d-1", cfgDream.DecisionAccepted) {
		t.Fatal("ledger missing accepted archive")
	}
}

// TestRejectThroughCLI records a rejection with no mutation.
func TestRejectThroughCLI(t *testing.T) {
	root, dreamsDir := setupProject(t, []engine.Proposal{
		{
			ID:         "d-2",
			Targets:    []string{filepath.Join(cfgDir.Ideas, "note.md")},
			Status:     cfgDream.StatusSidenote,
			Action:     cfgDream.ActionArchive,
			Confidence: cfgDream.ConfidenceLow,
		},
	})

	c := newCmd()
	if err := dispose.Reject(c, "d-2", "no"); err != nil {
		t.Fatalf("Reject: %v", err)
	}
	if _, statErr := os.Stat(
		filepath.Join(root, cfgDir.Ideas, "note.md"),
	); statErr != nil {
		t.Fatal("reject must not move the source")
	}
	if !ledgerHas(t, dreamsDir, "d-2", cfgDream.DecisionRejected) {
		t.Fatal("ledger missing rejection")
	}
}

// TestAmendThroughCLI applies a different action and records amended.
func TestAmendThroughCLI(t *testing.T) {
	_, dreamsDir := setupProject(t, []engine.Proposal{
		{
			ID:         "d-3",
			Targets:    []string{filepath.Join(cfgDir.Ideas, "note.md")},
			Status:     cfgDream.StatusMeritorious,
			Action:     cfgDream.ActionArchive,
			Confidence: cfgDream.ConfidenceMed,
		},
	})

	c := newCmd()
	if err := dispose.Amend(
		c, "d-3", cfgDream.ActionKeep, "",
	); err != nil {
		t.Fatalf("Amend: %v", err)
	}
	if !ledgerHas(t, dreamsDir, "d-3", cfgDream.DecisionAmended) {
		t.Fatal("ledger missing amended decision")
	}
}

// TestAcceptMissingProposalErrors errors when the id is unknown.
func TestAcceptMissingProposalErrors(t *testing.T) {
	setupProject(t, []engine.Proposal{})
	c := newCmd()
	if err := dispose.Accept(c, "ghost", ""); err == nil {
		t.Fatal("Accept must error on a missing proposal id")
	}
}

// --- helpers ---

func newCmd() *cobra.Command {
	c := &cobra.Command{}
	c.SetOut(&bytes.Buffer{})
	c.SetErr(&bytes.Buffer{})
	return c
}

func ledgerHas(
	t *testing.T, dreamsDir, id string, decision cfgDream.Decision,
) bool {
	t.Helper()
	entries, readErr := engine.ReadLedger(dreamsDir)
	if readErr != nil {
		t.Fatalf("read ledger: %v", readErr)
	}
	for _, e := range entries {
		if e.ProposalID == id && e.Decision == decision {
			return true
		}
	}
	return false
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir parent %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func gitInit(t *testing.T, root string) {
	t.Helper()
	run := func(args ...string) {
		//nolint:gosec // test fixture, hardcoded args
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test")
	mustWrite(t, filepath.Join(root, ".gitignore"), "dreams/\nideas/\n")
	run("add", ".gitignore")
	run("commit", "-q", "-m", "init")
}
