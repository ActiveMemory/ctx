//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for `ctx kb` CLI output strings.
const (
	// DescKeyWriteKbFindingLine is the findings-log line
	// format (timestamp + text).
	DescKeyWriteKbFindingLine = "write.kb.finding-line"
	// DescKeyWriteKbReindexed names how many topics were
	// folded into the managed block and the path rewritten.
	DescKeyWriteKbReindexed = "write.kb.reindexed"
	// DescKeyWriteKbScaffolded names the path of a newly
	// scaffolded topic-page file.
	DescKeyWriteKbScaffolded = "write.kb.scaffolded"
	// DescKeyWriteKbAppendedTo names the destination of a
	// successful note append.
	DescKeyWriteKbAppendedTo = "write.kb.appended-to"
	// DescKeyWriteKbAskDrivenHint announces the canonical
	// /ctx-kb-ask skill invocation.
	DescKeyWriteKbAskDrivenHint = "write.kb.ask-driven-hint"
	// DescKeyWriteKbAskInvokeFormat carries the inline
	// invocation example.
	DescKeyWriteKbAskInvokeFormat = "write.kb.ask-invoke-format"
	// DescKeyWriteKbAskContractPointer points at the ask
	// contract source-of-truth.
	DescKeyWriteKbAskContractPointer = "write.kb.ask-contract-pointer"
	// DescKeyWriteKbGroundDrivenHint announces the canonical
	// /ctx-kb-ground skill invocation.
	DescKeyWriteKbGroundDrivenHint = "write.kb.ground-driven-hint"
	// DescKeyWriteKbGroundContractPointer points at the
	// ground contract source-of-truth.
	DescKeyWriteKbGroundContractPointer = "write.kb.ground-contract-pointer"
	// DescKeyWriteKbIngestDrivenHint announces the canonical
	// /ctx-kb-ingest skill invocation.
	DescKeyWriteKbIngestDrivenHint = "write.kb.ingest-driven-hint"
	// DescKeyWriteKbIngestInvokeFormat carries the inline
	// ingest invocation example.
	DescKeyWriteKbIngestInvokeFormat = "write.kb.ingest-invoke-format"
	// DescKeyWriteKbIngestFallbackHint points at the
	// hand-fallback prompt.
	DescKeyWriteKbIngestFallbackHint = "write.kb.ingest-fallback-hint"
	// DescKeyWriteKbSiteReviewDrivenHint announces the
	// canonical /ctx-kb-site-review skill invocation.
	DescKeyWriteKbSiteReviewDrivenHint = "write.kb.site-review-driven-hint"
	// DescKeyWriteKbSiteReviewContractPointer points at the
	// site-review contract source-of-truth.
	DescKeyWriteKbSiteReviewContractPointer = "write.kb.site-review-contract-pointer"
)
