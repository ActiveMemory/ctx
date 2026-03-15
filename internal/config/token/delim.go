//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

const emDash = "—"
const enDash = "-"

const (
	// Colon is the colon character used as a key-value separator.
	Colon = ":"
	// Dash is a hyphen used as a timestamp segment separator.
	Dash = "-"
	// KeyValueSep is the equals sign used as a key-value separator in state files.
	KeyValueSep = "="
	// Separator is a Markdown horizontal rule used between sections.
	Separator = "---"
	// Ellipsis is a Markdown ellipsis.
	Ellipsis = "..."
)

// TopicSeparators are the delimiters between a date and topic in session
// headers (e.g., "2026-01-15 - Fix API").
var TopicSeparators = []string{" " + emDash + " ", " " + enDash + " "}
