//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence

import "time"

// Row is one EV-### evidence-index entry. Fields match the
// schema at
// `internal/assets/kb/templates/ingest/schemas/evidence-index.md`.
type Row struct {
	// ID is the EV-### identifier. When empty on Append, the
	// next sequential ID is allocated and populated on return.
	ID string
	// Claim is the assertion this row backs.
	Claim string
	// SourceID is the short-name from `source-map.md` that
	// identifies the source.
	SourceID string
	// Locator pins the claim within the source: line range for
	// files, timestamp for transcripts, anchor for URLs,
	// symbol for code.
	Locator string
	// SHA is the short git SHA when the source is an in-repo
	// file. Empty for out-of-repo citations.
	SHA string
	// Confidence is one of the cfgKB.Confidence* constants.
	Confidence string
	// Tags is the per-row tag list. The `evidence-only` tag
	// marks rows minted in evidence-only mode passes (those
	// rows are review-required before a topic-page pass cites
	// them).
	Tags []string
	// Occurred is the date the cited statement was made (for
	// transcripts / dated docs). Empty when not dated.
	Occurred string
	// Extracted is the timestamp at which this row was minted.
	Extracted time.Time
}
