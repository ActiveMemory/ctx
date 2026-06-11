//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	engine "github.com/ActiveMemory/ctx/internal/dream"
	writeDream "github.com/ActiveMemory/ctx/internal/write/dream"
)

// captureCmd returns a cobra command whose output is captured into buf.
func captureCmd() (*cobra.Command, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	c := &cobra.Command{}
	c.SetOut(buf)
	c.SetErr(buf)
	return c, buf
}

// TestReviewRendersSubstance renders every substance field of each
// pending proposal.
func TestReviewRendersSubstance(t *testing.T) {
	c, buf := captureCmd()
	proposals := []engine.Proposal{
		{
			ID:         "p-1",
			Targets:    []string{"ideas/a.md", "ideas/b.md"},
			Status:     cfgDream.StatusDuplicate,
			Action:     cfgDream.ActionMerge,
			Evidence:   "near-neighbor ideas/b.md (0.9)",
			Confidence: cfgDream.ConfidenceHigh,
			Rationale:  "restates b",
		},
	}

	writeDream.Review(c, proposals)
	out := buf.String()

	for _, want := range []string{
		"p-1", cfgDream.StatusDuplicate, cfgDream.ActionMerge,
		cfgDream.ConfidenceHigh, "ideas/a.md", "ideas/b.md",
		"near-neighbor ideas/b.md (0.9)", "restates b",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("review output missing %q\ngot:\n%s", want, out)
		}
	}
}

// TestReviewNonePrintsMessage prints the no-pending message.
func TestReviewNonePrintsMessage(t *testing.T) {
	c, buf := captureCmd()
	writeDream.ReviewNone(c)
	if strings.TrimSpace(buf.String()) == "" {
		t.Fatal("ReviewNone printed nothing")
	}
}

// TestDigestPrintsCounts prints the source and proposal counts.
func TestDigestPrintsCounts(t *testing.T) {
	c, buf := captureCmd()
	writeDream.Digest(c, 5, 2)
	out := buf.String()
	if !strings.Contains(out, "5") || !strings.Contains(out, "2") {
		t.Fatalf("digest missing counts: %q", out)
	}
}

// TestDispositionGenerativeRoutesToSerendipity points the user at the
// serendipity skill for a generative result.
func TestDispositionGenerativeRoutesToSerendipity(t *testing.T) {
	c, buf := captureCmd()
	writeDream.Disposition(c, "p-9",
		cfgDream.DecisionAccepted,
		engine.ApplyResult{
			Generative: true, Action: cfgDream.ActionPromote,
		},
	)
	out := buf.String()
	if !strings.Contains(out, "p-9") ||
		!strings.Contains(out, "serendipity") {
		t.Fatalf("generative disposition should route to serendipity: %q", out)
	}
}

// TestDispositionMechanicalConfirms confirms an accepted mechanical
// disposition.
func TestDispositionMechanicalConfirms(t *testing.T) {
	c, buf := captureCmd()
	writeDream.Disposition(c, "p-3",
		cfgDream.DecisionAccepted,
		engine.ApplyResult{
			Performed: true, Action: cfgDream.ActionArchive,
		},
	)
	out := buf.String()
	if !strings.Contains(out, "p-3") ||
		!strings.Contains(out, cfgDream.ActionArchive) {
		t.Fatalf("mechanical disposition not confirmed: %q", out)
	}
}

// TestNilCmdNoPanic verifies the helpers no-op on a nil command.
func TestNilCmdNoPanic(t *testing.T) {
	writeDream.Nothing(nil)
	writeDream.Locked(nil)
	writeDream.Digest(nil, 1, 1)
	writeDream.Failmark(nil, "x")
	writeDream.ReviewNone(nil)
	writeDream.Review(nil, nil)
	writeDream.Disposition(nil, "x",
		cfgDream.DecisionRejected, engine.ApplyResult{})
}
