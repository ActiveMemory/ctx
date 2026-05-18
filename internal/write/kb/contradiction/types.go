//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contradiction

// Row is one contradictions entry. Fields match the schema at
// `internal/assets/kb/templates/ingest/schemas/contradictions.md`,
// rendered as a single markdown table row. The ID is allocated
// by [Append]; callers do not supply it.
type Row struct {
	// EVRefs is the list of `EV-###` IDs that disagree. The
	// schema requires at least two; the writer does not enforce
	// that; the caller does.
	EVRefs []string
	// Summary is a one-line statement of what the rows
	// disagree about.
	Summary string
	// DemotionApplied is the EV id that got demoted plus its
	// new band, e.g. "EV-031 -> low".
	DemotionApplied string
	// Status is one of "open", "resolved".
	Status string
}
