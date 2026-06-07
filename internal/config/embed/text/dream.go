//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for ctx-dream user-facing write output.
const (
	// DescKeyWriteDreamNothing is the text key for the empty-delta
	// "nothing to dream" message.
	DescKeyWriteDreamNothing = "write.dream-nothing"
	// DescKeyWriteDreamLocked is the text key for the lock-held
	// exit-0 message.
	DescKeyWriteDreamLocked = "write.dream-locked"
	// DescKeyWriteDreamDigest is the text key for the post-pass
	// counts digest.
	DescKeyWriteDreamDigest = "write.dream-digest"
	// DescKeyWriteDreamFailmark is the text key for the fail-loud
	// failmark-written message.
	DescKeyWriteDreamFailmark = "write.dream-failmark"
	// DescKeyWriteDreamReviewNone is the text key for the no-pending
	// review message.
	DescKeyWriteDreamReviewNone = "write.dream-review-none"
	// DescKeyWriteDreamReviewHeader is the text key for the pending
	// review header with count.
	DescKeyWriteDreamReviewHeader = "write.dream-review-header"
	// DescKeyWriteDreamReviewID is the text key for a proposal's
	// id/status/action/confidence line.
	DescKeyWriteDreamReviewID = "write.dream-review-id"
	// DescKeyWriteDreamReviewTargets is the text key for a proposal's
	// targets line.
	DescKeyWriteDreamReviewTargets = "write.dream-review-targets"
	// DescKeyWriteDreamReviewEvidence is the text key for a
	// proposal's evidence line.
	DescKeyWriteDreamReviewEvidence = "write.dream-review-evidence"
	// DescKeyWriteDreamReviewRationale is the text key for a
	// proposal's rationale line.
	DescKeyWriteDreamReviewRationale = "write.dream-review-rationale"
	// DescKeyWriteDreamAccepted is the text key for the accepted
	// mechanical disposition confirmation.
	DescKeyWriteDreamAccepted = "write.dream-accepted"
	// DescKeyWriteDreamRejected is the text key for the rejected
	// confirmation.
	DescKeyWriteDreamRejected = "write.dream-rejected"
	// DescKeyWriteDreamAmended is the text key for the amended
	// disposition confirmation.
	DescKeyWriteDreamAmended = "write.dream-amended"
	// DescKeyWriteDreamGenerative is the text key for the accepted
	// generative-intent message routing to /ctx-serendipity.
	DescKeyWriteDreamGenerative = "write.dream-generative"
)
