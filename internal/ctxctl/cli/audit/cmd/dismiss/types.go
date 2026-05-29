//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

// Strings carries the English user-facing text for the audit
// dismiss subcommand. ctxctl supplies these from its own Go
// constants; the logic package holds no copy of its own.
type Strings struct {
	// Use is the cobra Use string.
	Use string
	// Short is the one-line command description.
	Short string
	// Example is the example-usage block.
	Example string
	// AllFlag is the --all flag description.
	AllFlag string
	// Dismissed is the single-dismissal confirmation format
	// string (one verb: id).
	Dismissed string
	// DismissedAll is the bulk-dismissal confirmation format
	// string (one verb: count).
	DismissedAll string
}
