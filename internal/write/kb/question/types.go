//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package question

// Row is one outstanding-questions entry. Fields match the
// schema at
// `internal/assets/kb/templates/ingest/schemas/outstanding-questions.md`,
// rendered as a single markdown table row. The ID is allocated
// by [Append]; callers do not supply it.
type Row struct {
	// Question is the one-sentence open question, phrased as
	// an interrogative.
	Question string
	// WhyItMatters explains why answering this changes a topic
	// page or a confidence band.
	WhyItMatters string
	// WhatEvidenceWouldResolve is the shape of evidence that
	// would close this entry. Required by the schema; opening a
	// question without this field is a hard anti-pattern.
	WhatEvidenceWouldResolve string
	// OpenedAt is the ISO date the question was opened.
	OpenedAt string
	// RelatedEV is the optional list of `EV-###` refs.
	RelatedEV []string
}
