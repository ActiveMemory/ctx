//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkaudit

// Strings carries the English user-facing text for the
// audit-relay hook and the relay box it renders. ctxctl
// supplies these from its own Go constants; the logic
// package holds no copy of its own.
type Strings struct {
	// Use is the cobra Use string (the relay verb).
	Use string
	// Short is the one-line command description.
	Short string
	// Example is the example-usage block.
	Example string
	// RelayLabel labels the relay event in the event log /
	// template ref (provenance, not user-facing copy).
	RelayLabel string
	// RelayVariant is the relay event variant (provenance).
	RelayVariant string
	// BoxTitle titles the relay nudge box.
	BoxTitle string
	// RelayPrefix is the VERBATIM-relay instruction line.
	RelayPrefix string
	// NudgeFormat is the human-readable relay-message
	// suffix (one verb: count).
	NudgeFormat string
	// DismissHint is the per-id dismissal hint in the box.
	DismissHint string
	// DismissAllHint is the dismiss-all hint in the box.
	DismissAllHint string
	// ReportSeparator is the rule between multiple reports.
	ReportSeparator string
	// StalePrefix is the per-report STALE prefix format
	// (verbs: commit-range, age).
	StalePrefix string
}
