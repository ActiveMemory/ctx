//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for `ctx kb` CLI output strings.
const (
	// DescKeyWriteKbCliFindingLine is the findings-log line
	// format (timestamp + text).
	DescKeyWriteKbCliFindingLine = "write.kb-cli.finding-line"
	// DescKeyWriteKbCliReindexed names how many topics were
	// folded into the managed block and the path rewritten.
	DescKeyWriteKbCliReindexed = "write.kb-cli.reindexed"
	// DescKeyWriteKbCliScaffolded names the path of a newly
	// scaffolded topic-page file.
	DescKeyWriteKbCliScaffolded = "write.kb-cli.scaffolded"
	// DescKeyWriteKbCliAppendedTo names the destination of a
	// successful note append.
	DescKeyWriteKbCliAppendedTo = "write.kb-cli.appended-to"
	// DescKeyWriteKbCliAskDrivenHint announces the canonical
	// /ctx-kb-ask skill invocation.
	DescKeyWriteKbCliAskDrivenHint = "write.kb-cli.ask-driven-hint"
	// DescKeyWriteKbCliAskInvokeFormat carries the inline
	// invocation example.
	DescKeyWriteKbCliAskInvokeFormat = "write.kb-cli.ask-invoke-format"
	// DescKeyWriteKbCliAskContractPointer points at the ask
	// contract source-of-truth.
	DescKeyWriteKbCliAskContractPointer = "write.kb-cli.ask-contract-pointer"
	// DescKeyWriteKbCliGroundDrivenHint announces the canonical
	// /ctx-kb-ground skill invocation.
	DescKeyWriteKbCliGroundDrivenHint = "write.kb-cli.ground-driven-hint"
	// DescKeyWriteKbCliGroundContractPointer points at the
	// ground contract source-of-truth.
	DescKeyWriteKbCliGroundContractPointer = "write.kb-cli.ground-contract-pointer"
	// DescKeyWriteKbCliIngestDrivenHint announces the canonical
	// /ctx-kb-ingest skill invocation.
	DescKeyWriteKbCliIngestDrivenHint = "write.kb-cli.ingest-driven-hint"
	// DescKeyWriteKbCliIngestInvokeFormat carries the inline
	// ingest invocation example.
	DescKeyWriteKbCliIngestInvokeFormat = "write.kb-cli.ingest-invoke-format"
	// DescKeyWriteKbCliIngestFallbackHint points at the
	// hand-fallback prompt.
	DescKeyWriteKbCliIngestFallbackHint = "write.kb-cli.ingest-fallback-hint"
	// DescKeyWriteKbCliSiteReviewDrivenHint announces the
	// canonical /ctx-kb-site-review skill invocation.
	DescKeyWriteKbCliSiteReviewDrivenHint = "write.kb-cli.site-review-driven-hint"
	// DescKeyWriteKbCliSiteReviewContractPointer points at the
	// site-review contract source-of-truth.
	DescKeyWriteKbCliSiteReviewContractPointer = "write.kb-cli.site-review-contract-pointer"
)
