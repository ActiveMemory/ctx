//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

// KnownStatuses enumerates every valid ProposalStatus. It is the
// single source of truth a validator checks an observed classification
// against, and it keeps every status constant referenced from
// production code.
var KnownStatuses = []ProposalStatus{
	StatusImplemented,
	StatusDuplicate,
	StatusMeritorious,
	StatusSidenote,
	StatusBlogCandidate,
}

// KnownActions enumerates every valid ProposalAction. ActionPromote is
// the only action that sanctions a write into specs/, so callers also
// match against it directly; listing it here keeps every action
// constant referenced from production code.
var KnownActions = []ProposalAction{
	ActionArchive,
	ActionMerge,
	ActionPromote,
	ActionMarkBlog,
	ActionKeep,
}

// KnownConfidences enumerates every valid Confidence level.
var KnownConfidences = []Confidence{
	ConfidenceHigh,
	ConfidenceMed,
	ConfidenceLow,
}

// KnownSourceStatuses enumerates every valid SourceStatus lifecycle
// state a per-source record may hold.
var KnownSourceStatuses = []SourceStatus{
	SourceActive,
	SourceArchived,
	SourcePromoted,
	SourceMerged,
}

// KnownDecisions enumerates every valid review Decision recorded in the
// ledger, including rejections and skips so dedup-against-seen has the
// full vocabulary.
var KnownDecisions = []Decision{
	DecisionAccepted,
	DecisionRejected,
	DecisionAmended,
	DecisionSkipped,
}
