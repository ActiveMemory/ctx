//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"testing"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/dream"
)

// TestProposalValid accepts a fully-known proposal and rejects ones
// carrying an unknown status, action, or confidence.
func TestProposalValid(t *testing.T) {
	good := dream.Proposal{
		ID:         "p1",
		Targets:    []string{"ideas/a.md"},
		Status:     cfgDream.StatusMeritorious,
		Action:     cfgDream.ActionKeep,
		Evidence:   "spec",
		Confidence: cfgDream.ConfidenceHigh,
	}
	if err := dream.ProposalValid(good); err != nil {
		t.Fatalf("known proposal rejected: %v", err)
	}

	cases := []struct {
		name string
		p    dream.Proposal
	}{
		{
			name: "bad status",
			p: dream.Proposal{
				ID: "p2", Status: "nonsense",
				Action:     cfgDream.ActionKeep,
				Confidence: cfgDream.ConfidenceHigh,
			},
		},
		{
			name: "bad action",
			p: dream.Proposal{
				ID: "p3", Status: cfgDream.StatusSidenote,
				Action: "nuke", Confidence: cfgDream.ConfidenceLow,
			},
		},
		{
			name: "bad confidence",
			p: dream.Proposal{
				ID: "p4", Status: cfgDream.StatusDuplicate,
				Action: cfgDream.ActionMerge, Confidence: "certain",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := dream.ProposalValid(tc.p); err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

// TestSourceStatusAndDecisionKnown exercises the lifecycle and decision
// predicates across known and unknown values.
func TestSourceStatusAndDecisionKnown(t *testing.T) {
	if !dream.SourceStatusKnown(cfgDream.SourceMerged) {
		t.Fatal("SourceMerged must be known")
	}
	if dream.SourceStatusKnown("zombie") {
		t.Fatal("unknown source status must be rejected")
	}
	if !dream.DecisionKnown(cfgDream.DecisionAmended) {
		t.Fatal("DecisionAmended must be known")
	}
	if dream.DecisionKnown("maybe") {
		t.Fatal("unknown decision must be rejected")
	}
}
