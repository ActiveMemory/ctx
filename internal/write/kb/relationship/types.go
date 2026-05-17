//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package relationship

// Row is one relationship-map entry. Fields match the schema at
// `internal/assets/kb/templates/ingest/schemas/relationship-map.md`,
// rendered as a single markdown table row.
type Row struct {
	// From is the originating topic slug or `EV-###` ID.
	From string
	// To is the destination topic slug or `EV-###` ID.
	To string
	// Kind is one of "depends-on", "refines", "contradicts",
	// "supersedes", "relates-to". Introducing a new kind
	// requires a `domain-decisions.md` entry naming the
	// rationale.
	Kind string
	// Summary is a one-line description of the relationship.
	Summary string
}
