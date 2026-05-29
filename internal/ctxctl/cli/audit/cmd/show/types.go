//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

// Strings carries the English user-facing text for the audit
// show subcommand. ctxctl supplies these from its own Go
// constants; the logic package holds no copy of its own.
type Strings struct {
	// Use is the cobra Use string.
	Use string
	// Short is the one-line command description.
	Short string
	// Example is the example-usage block.
	Example string
}
