//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for site feed generation.
const (
	DescKeySiteFeedGenerated = "site.feed-generated"
	DescKeySiteFeedSkipped   = "site.feed-skipped"
	DescKeySiteFeedWarnings  = "site.feed-warnings"
	DescKeySiteFeedItem      = "site.feed-item"
)

// DescKeys for site generation skip reasons.
const (
	DescKeySiteSkipCannotRead    = "site.skip-cannot-read"
	DescKeySiteSkipNoFrontmatter = "site.skip-no-frontmatter"
	DescKeySiteSkipMalformed     = "site.skip-malformed"
	DescKeySiteSkipParseError    = "site.skip-parse-error"
	DescKeySiteSkipNotFinalized  = "site.skip-not-finalized"
	DescKeySiteSkipMissingTitle  = "site.skip-missing-title"
	DescKeySiteSkipMissingDate   = "site.skip-missing-date"
	DescKeySiteWarnNoSummary     = "site.warn-no-summary"
)
