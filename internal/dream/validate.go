//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"fmt"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
)

// ProposalValid validates that a proposal carries known enum values AND
// the provenance a gated proposal requires: a non-empty target and
// non-empty evidence. It is the schema gate the review and ledger build
// on — an unrecognized field or stripped provenance is rejected before
// the proposal is surfaced or applied. Rejecting evidence-less proposals
// is the spec's "no evidence is not surfaced" rule and the front line
// against corrupted artifacts whose citations have been lost.
//
// Parameters:
//   - p: the proposal to validate
//
// Returns:
//   - error: nil when every field is a known value and provenance is
//     present; otherwise an err/dream.InvalidProposal naming the first
//     offending field
func ProposalValid(p Proposal) error {
	if !statusKnown(p.Status) {
		return errDream.InvalidProposal(p.ID, fmt.Sprintf(
			cfgDream.ReasonUnknownValue,
			cfgDream.FieldStatus, p.Status,
		))
	}
	if !actionKnown(p.Action) {
		return errDream.InvalidProposal(p.ID, fmt.Sprintf(
			cfgDream.ReasonUnknownValue,
			cfgDream.FieldAction, p.Action,
		))
	}
	if !confidenceKnown(p.Confidence) {
		return errDream.InvalidProposal(p.ID, fmt.Sprintf(
			cfgDream.ReasonUnknownValue,
			cfgDream.FieldConfidence, p.Confidence,
		))
	}
	if len(p.Targets) == 0 {
		return errDream.InvalidProposal(p.ID, fmt.Sprintf(
			cfgDream.ReasonMissing, cfgDream.FieldTargets,
		))
	}
	if p.Evidence == "" {
		return errDream.InvalidProposal(p.ID, fmt.Sprintf(
			cfgDream.ReasonMissing, cfgDream.FieldEvidence,
		))
	}
	return nil
}
