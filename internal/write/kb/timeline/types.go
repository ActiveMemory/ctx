//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package timeline

// Row is one timeline entry. Fields match the schema at
// `internal/assets/kb/templates/ingest/schemas/timeline.md`,
// rendered as a single markdown table row.
type Row struct {
	// Date is the ISO 8601 date the event occurred (not the
	// date it was ingested).
	Date string
	// Event is a one-paragraph description of what happened.
	Event string
	// SourceEV is the comma-list of `EV-###` refs that ground
	// the event. At least one is required by the schema; the
	// writer does not enforce that; the caller does.
	SourceEV []string
	// RelatedTopics is the optional list of topic slugs the
	// event touches.
	RelatedTopics []string
}
