//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package question

// Markdown rendering constants for the outstanding-questions
// artifact. Structural literals stay as Go consts; error
// wrapping format strings have moved to
// commands/text/errors.yaml.
const (
	// IDPrefix is the stable prefix for outstanding-question
	// IDs.
	IDPrefix = "Q-"
	// IDFormat is the Printf format for a zero-padded
	// `Q-NNN` identifier.
	IDFormat = "%s%03d"
	// TableHeader is the markdown table header row plus its
	// delimiter row.
	TableHeader = "| ID | Question | Why it matters |" +
		" What evidence would resolve | Opened at |" +
		" Related EV |\n" +
		"|----|----------|----------------|" +
		"-----------------------------|-----------|" +
		"------------|\n"
)
