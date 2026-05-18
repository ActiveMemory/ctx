//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package glossary

// Row is one glossary entry. Fields match the schema at
// `internal/assets/kb/templates/ingest/schemas/glossary.md`,
// rendered as a single markdown table row.
type Row struct {
	// Term is the canonical term; lowercase preferred.
	Term string
	// Definition is one declarative paragraph; no hedging beyond
	// the confidence band.
	Definition string
	// Confidence is one of "high", "medium", "low",
	// "speculative".
	Confidence string
	// EVRefs is the list of `EV-###` IDs that ground the
	// definition. At least one is required by the schema; the
	// writer does not enforce that; the caller does.
	EVRefs []string
	// RelatedTerms is the optional list of other glossary terms
	// cross-referenced from this entry.
	RelatedTerms []string
}
