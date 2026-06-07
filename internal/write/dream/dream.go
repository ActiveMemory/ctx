//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
	engine "github.com/ActiveMemory/ctx/internal/dream"
)

// Nothing prints the empty-delta message and is the no-work exit.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
func Nothing(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteDreamNothing))
}

// Locked prints the lock-held message for the exit-0 path.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
func Locked(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteDreamLocked))
}

// Digest prints the post-pass counts digest.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
//   - sources: number of sources processed this pass.
//   - proposals: number of valid proposals the executor wrote.
func Digest(cmd *cobra.Command, sources, proposals int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteDreamDigest), sources, proposals,
	))
}

// Failmark prints the fail-loud failmark-written notice.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
//   - path: the failmark file path.
func Failmark(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteDreamFailmark), path,
	))
}

// ReviewNone prints the no-pending-proposals review message.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
func ReviewNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteDreamReviewNone))
}

// Review renders the pending proposals substance-forward: a header
// with the count, then each proposal's id/status/action/confidence,
// targets, evidence, and rationale.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
//   - proposals: the pending proposals to render.
func Review(cmd *cobra.Command, proposals []engine.Proposal) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteDreamReviewHeader), len(proposals),
	))
	for _, p := range proposals {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamReviewID),
			p.ID, p.Status, p.Action, p.Confidence,
		))
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamReviewTargets),
			strings.Join(p.Targets, cfgToken.CommaSpace),
		))
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamReviewEvidence), p.Evidence,
		))
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamReviewRationale), p.Rationale,
		))
	}
}

// Disposition prints the confirmation for an applied decision. A
// generative result routes the user to /ctx-serendipity; mechanical
// results confirm the action; a rejection confirms the rejection.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
//   - id: the proposal ID.
//   - decision: the recorded review decision.
//   - res: the apply result describing how the action dispatched.
func Disposition(
	cmd *cobra.Command, id string,
	decision cfgDream.Decision, res engine.ApplyResult,
) {
	if cmd == nil {
		return
	}
	if res.Generative {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamGenerative), id, res.Action,
		))
		return
	}
	switch decision {
	case cfgDream.DecisionRejected:
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamRejected), id,
		))
	case cfgDream.DecisionAmended:
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamAmended), id, res.Action,
		))
	default:
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDreamAccepted), id, res.Action,
		))
	}
}
