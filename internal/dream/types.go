//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"time"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
)

// Proposal is one atomic, provenance-bearing triage proposal the dream
// emits into the gitignored dreams/ notebook for human review. The dream
// never acts on a proposal; the serendipity gate does.
type Proposal struct {
	// ID is the stable identifier for accept/reject/amend and ledger
	// reference. Stable so v2 supersession is not foreclosed.
	ID string `json:"id"`
	// Targets are the idea file path(s) the proposal concerns (more than
	// one for a merge).
	Targets []string `json:"targets"`
	// Status is the observed classification of the idea.
	Status cfgDream.ProposalStatus `json:"status"`
	// Action is the recommended disposition.
	Action cfgDream.ProposalAction `json:"action"`
	// Evidence is the grounding citation (commit, spec, or near-neighbor
	// + similarity) that justifies the classification. Required: a
	// proposal with no evidence is not surfaced.
	Evidence string `json:"evidence"`
	// Confidence drives attention triage during review.
	Confidence cfgDream.Confidence `json:"confidence"`
	// Rationale is a one-line, human-readable why.
	Rationale string `json:"rationale"`
}

// SourceState is the per-idea record the dream tracks across passes
// (persisted in dreams/state.json). It drives the discipline clock:
// re-triage only when the content hash changes.
type SourceState struct {
	// Path is the idea file path, relative to the project root.
	Path string `json:"path"`
	// Hash is the content hash; an unchanged hash means skip re-triage.
	Hash string `json:"hash"`
	// SummaryRef points at the cached summary for this source; it is
	// regenerated when Hash changes.
	SummaryRef string `json:"summary_ref,omitempty"`
	// LastModified is the source file's last-modified time at last triage.
	LastModified time.Time `json:"last_modified"`
	// LastSurfaced is when this source was last shown in a review round.
	LastSurfaced time.Time `json:"last_surfaced,omitempty"`
	// Merit is the attention-ranking score (0..1) feeding ruthless
	// self-rejection — never an autonomous promote threshold.
	Merit float64 `json:"merit"`
	// Status is the lifecycle state of the idea.
	Status cfgDream.SourceStatus `json:"status"`
	// History is the chronological list of dispositions decided for this
	// source's proposals.
	History []LedgerEntry `json:"history,omitempty"`
}

// GuardDecision is the structured result of a guard check: whether a
// write target is allowed and, when refused, a human-readable reason
// sourced from the dream error/text registry (never an inline English
// literal). It is returned by the executor-agnostic guard logic so a
// Claude Code PreToolUse hook and a raw-API tool executor enforce the
// same invariant the same way.
type GuardDecision struct {
	// Allowed reports whether the write target passed the guard.
	Allowed bool
	// Reason explains a refusal (empty when Allowed is true). The text
	// originates from internal/err/dream, so guards carry no inline
	// English string literals.
	Reason string
}

// LedgerEntry is one disposition recorded in the append-only ledger
// (dreams/ledger.md). Rejections are recorded too, so dedup-against-seen
// keeps decided proposals from re-surfacing unless the source changes.
type LedgerEntry struct {
	// ProposalID links back to the Proposal this disposition resolved.
	ProposalID string `json:"proposal_id"`
	// Decision is the human's review outcome.
	Decision cfgDream.Decision `json:"decision"`
	// Action is the disposition that was applied (may differ from the
	// proposed action when the decision was amended).
	Action cfgDream.ProposalAction `json:"action"`
	// At is when the disposition was recorded.
	At time.Time `json:"at"`
	// Note is an optional human note captured at decision time.
	Note string `json:"note,omitempty"`
}
