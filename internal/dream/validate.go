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

// SourceStatusKnown reports whether s is a recognized SourceStatus
// lifecycle state. Exported so state-record producers can validate a
// status before persisting it.
//
// Parameters:
//   - s: the source status value to check
//
// Returns:
//   - bool: true when s is one of cfgDream.KnownSourceStatuses
func SourceStatusKnown(s cfgDream.SourceStatus) bool {
	for _, known := range cfgDream.KnownSourceStatuses {
		if s == known {
			return true
		}
	}
	return false
}

// DecisionKnown reports whether d is a recognized review Decision.
// Exported so ledger producers can validate a decision before
// recording it.
//
// Parameters:
//   - d: the decision value to check
//
// Returns:
//   - bool: true when d is one of cfgDream.KnownDecisions
func DecisionKnown(d cfgDream.Decision) bool {
	for _, known := range cfgDream.KnownDecisions {
		if d == known {
			return true
		}
	}
	return false
}

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
