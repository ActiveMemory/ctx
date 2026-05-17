//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contradiction

// Markdown rendering constants for the contradictions
// artifact. Structural literals stay as Go consts; error
// wrapping format strings have moved to
// commands/text/errors.yaml.
const (
	// IDPrefix is the stable prefix for contradiction IDs.
	IDPrefix = "C-"
	// IDFormat is the Printf format for a zero-padded
	// `C-NNN` identifier.
	IDFormat = "%s%03d"
	// TableHeader is the markdown table header row plus its
	// delimiter row.
	TableHeader = "| ID | Evidence | Summary | Demotion " +
		"applied | Status |\n" +
		"|----|----------|---------|------------------|" +
		"--------|\n"
)
