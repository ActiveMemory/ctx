//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence

// Markdown rendering constants for the evidence-index file.
// Structural literals (headings, table shape, ID format)
// stay as Go consts.
const (
	// TitleHeading is the H1 written above the evidence table
	// when the file is first created.
	TitleHeading = "# Evidence index"
	// LeadParagraph1 is the first instructional paragraph
	// rendered after the title.
	LeadParagraph1 = "Append-only EV-### rows. " +
		"Never renumber, never delete."
	// LeadParagraph2 is the second instructional paragraph
	// rendered after the title.
	LeadParagraph2 = "Demote in place when reconciliation " +
		"requires it (see KB-RULES.md)."
	// TableHeader is the markdown table header row plus its
	// delimiter row.
	TableHeader = "| ID | Claim | Source | Locator | sha |" +
		" Confidence | Tags | Occurred | Extracted |\n" +
		"|---|---|---|---|---|---|---|---|---|"
	// RowFormat is the Printf format consumed by the row
	// renderer; operands map to the Row struct fields.
	RowFormat = "| %s | %s | %s | %s | %s | %s |" +
		" %s | %s | %s |"
	// IDFormat is the Printf format used by the ID minter
	// (prefix, digit-width, integer).
	IDFormat = "%s-%0*d"
)
