//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

// Mode is an execution mode for a dream pass. v1 builds only
// discipline; creative is sketched and deferred.
type Mode = string

// ProposalStatus is the observed classification the dream assigns to
// an idea during triage.
type ProposalStatus = string

// Proposal status constants — what the dream observed about an idea.
const (
	// StatusImplemented means the idea is already realized in code/specs
	// (the dream cites the commit or spec as evidence).
	StatusImplemented ProposalStatus = "implemented"
	// StatusDuplicate means the idea restates a near-neighbor already
	// captured elsewhere (the dream cites the neighbor).
	StatusDuplicate ProposalStatus = "duplicate"
	// StatusMeritorious means the idea is still live and worth keeping
	// or acting on.
	StatusMeritorious ProposalStatus = "meritorious"
	// StatusSidenote means the idea is a throwaway aside with little
	// standalone merit.
	StatusSidenote ProposalStatus = "sidenote"
	// StatusBlogCandidate means the idea reads as publishable material.
	StatusBlogCandidate ProposalStatus = "blog-candidate"
)

// ProposalAction is the disposition the dream recommends for an idea.
type ProposalAction = string

// Proposal action constants — the recommended disposition. Mechanical
// actions (archive, mark-blog, keep) apply with no LLM cost; generative
// actions (merge, promote) drop to the agent reading the full source.
const (
	// ActionArchive moves the idea to ideas/done/ (reversible by relocation).
	ActionArchive ProposalAction = "archive"
	// ActionMerge folds the idea into another (destructive: backup first).
	ActionMerge ProposalAction = "merge"
	// ActionPromote drafts specs/<name>.md from the full source — the one
	// sanctioned declassification across the don't-leak boundary.
	ActionPromote ProposalAction = "promote"
	// ActionMarkBlog tags the idea as blog material in place.
	ActionMarkBlog ProposalAction = "mark-blog"
	// ActionKeep leaves the idea untouched.
	ActionKeep ProposalAction = "keep"
)

// Confidence is the dream's self-assessed confidence in a proposal,
// driving attention triage during review.
type Confidence = string

// Confidence level constants.
const (
	// ConfidenceHigh is a high-confidence proposal.
	ConfidenceHigh Confidence = "high"
	// ConfidenceMed is a medium-confidence proposal.
	ConfidenceMed Confidence = "med"
	// ConfidenceLow is a low-confidence proposal.
	ConfidenceLow Confidence = "low"
)

// SourceStatus is the lifecycle state of a triaged idea, tracked in the
// per-source state record.
type SourceStatus = string

// Source status constants.
const (
	// SourceActive means the idea is live and eligible for triage.
	SourceActive SourceStatus = "active"
	// SourceArchived means the idea was moved to ideas/done/.
	SourceArchived SourceStatus = "archived"
	// SourcePromoted means the idea was drafted into specs/.
	SourcePromoted SourceStatus = "promoted"
	// SourceMerged means the idea was folded into another.
	SourceMerged SourceStatus = "merged"
)

// Decision is the human's disposition recorded in the ledger during a
// serendipity review.
type Decision = string

// Ledger decision constants — every review outcome, including rejections,
// so dedup-against-seen keeps decided items from re-surfacing.
const (
	// DecisionAccepted means the human accepted the proposed action.
	DecisionAccepted Decision = "accepted"
	// DecisionRejected means the human rejected the proposal; it is not
	// re-surfaced unless the source content changes.
	DecisionRejected Decision = "rejected"
	// DecisionAmended means the human changed the action before applying.
	DecisionAmended Decision = "amended"
	// DecisionSkipped means the human deferred; the proposal may re-surface.
	DecisionSkipped Decision = "skipped"
)

// Notebook file names within the gitignored dreams/ directory.
const (
	// FileState is the per-source state record file under dreams/.
	FileState = "state.json"
	// FileLedger is the append-only disposition ledger under dreams/.
	FileLedger = "ledger.md"
)

// Proposal field labels used in validation diagnostics — the field
// name plus the offending value form an invalid-proposal reason.
const (
	// FieldStatus labels the Status field in a validation reason.
	FieldStatus = "status"
	// FieldAction labels the Action field in a validation reason.
	FieldAction = "action"
	// FieldConfidence labels the Confidence field in a validation
	// reason.
	FieldConfidence = "confidence"
	// FieldEvidence labels the Evidence field in a validation reason.
	FieldEvidence = "evidence"
	// FieldTargets labels the Targets field in a validation reason.
	FieldTargets = "targets"
)

// ReasonUnknownValue is the format for an invalid-proposal reason:
// "<field> has unknown value %q". The field label and the offending
// value are filled by the validator.
const ReasonUnknownValue = "%s has unknown value %q"

// ReasonMissing is the format for an invalid-proposal reason when a
// required, provenance-bearing field is absent: "<field> is required".
// The field label is filled by the validator.
const ReasonMissing = "%s is required"

// JSONIndent is the indent unit used when encoding the state file as
// human-readable JSON (the notebook is meant to be inspectable).
const JSONIndent = "  "
