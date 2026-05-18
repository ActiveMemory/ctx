//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcemap

// Row is one source-map entry. Fields match the schema at
// `internal/assets/kb/templates/ingest/schemas/source-map.md`,
// rendered as a single markdown table row.
type Row struct {
	// ShortName is the stable identifier used in
	// `evidence-index.md` cites; never renamed in place.
	ShortName string
	// Kind is one of "transcript", "code", "doc", "url", "mcp",
	// "kb".
	Kind string
	// Locator is the URL, repo path, MCP resource ID, or kb
	// pointer for this source.
	Locator string
	// AdmissionStatus is one of "admitted", "skipped",
	// "pending". (The schema also names "rejected"; the brief
	// pins this writer's vocabulary to admitted|skipped per the
	// kb-editorial-pipeline spec.)
	AdmissionStatus string
	// AdmissionRationale is one sentence explaining the
	// admission decision.
	AdmissionRationale string
	// Dated is the optional ISO date the source itself is
	// dated (empty when unknown).
	Dated string
}
