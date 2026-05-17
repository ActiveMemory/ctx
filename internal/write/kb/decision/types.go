//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision

// Row is one domain-decisions entry. Fields match the schema at
// `internal/assets/kb/templates/ingest/schemas/domain-decisions.md`,
// rendered as a single markdown table row. The ID is allocated
// by [Append]; callers do not supply it.
type Row struct {
	// Date is the ISO date the decision was recorded in the kb.
	Date string
	// Context describes what in the domain prompted the
	// decision; observable facts only.
	Context string
	// Decision is the position taken, in one sentence.
	Decision string
	// Rationale records why this position over the
	// alternatives that were on the table.
	Rationale string
	// Consequence records what now changes for topic pages,
	// glossary, or downstream claims.
	Consequence string
	// SupportingEV is the list of `EV-###` refs that ground
	// the decision. The schema requires at least one; the
	// writer does not enforce that; the caller does.
	SupportingEV []string
}
