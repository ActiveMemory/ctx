//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage

// Markdown rendering constants for the source-coverage ledger.
// Structural literals stay as Go consts.
const (
	// TitleHeading is the H1 written above the ledger table.
	TitleHeading = "# Source coverage"
	// LeadParagraph1 is the first instructional paragraph.
	LeadParagraph1 = "State-machine ledger over every source " +
		"the kb has touched."
	// LeadParagraph2 is the second instructional paragraph.
	LeadParagraph2 = "Updated automatically by ctx kb " +
		"commands; hand-edits are honored"
	// LeadParagraph3 is the third instructional paragraph.
	LeadParagraph3 = "but the next pass will re-write rows " +
		"it touches."
	// TableHeader is the markdown table header row plus its
	// delimiter row.
	TableHeader = "| Source | Topic | State | EV coverage |" +
		" Residue | Next action | Updated |\n" +
		"|---|---|---|---|---|---|---|"
	// HeaderCellSource is the header cell that identifies a
	// header row during ledger parsing.
	HeaderCellSource = "Source"
	// DelimRowPrefix is the prefix that identifies a markdown
	// table delimiter row ("|---|...|").
	DelimRowPrefix = "|---"
	// ExpectedCellCount is the column count for a valid ledger
	// row (Source, Topic, State, EVCoverage, Residue,
	// NextAction, Updated).
	ExpectedCellCount = 7
)

// Column indexes into a parsed ledger row. The same indexes
// drive the row binder in
// [github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage]
// when mapping cells back to Row fields.
const (
	// ColSource is the index of the Source cell.
	ColSource = 0
	// ColTopic is the index of the Topic cell.
	ColTopic = 1
	// ColState is the index of the State cell.
	ColState = 2
	// ColEVCoverage is the index of the EV coverage cell.
	ColEVCoverage = 3
	// ColResidue is the index of the Residue cell.
	ColResidue = 4
	// ColNextAction is the index of the Next action cell.
	ColNextAction = 5
	// ColUpdated is the index of the Updated cell.
	ColUpdated = 6
)
