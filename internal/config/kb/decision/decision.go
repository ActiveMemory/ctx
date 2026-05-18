//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision

// Markdown rendering constants for the domain-decisions
// artifact. Structural literals stay as Go consts; error
// wrapping format strings have moved to
// commands/text/errors.yaml.
const (
	// IDPrefix is the stable prefix for domain-decision IDs.
	// The canonical schema pins this to `DD-` to keep the
	// namespace distinct from any future project-side `D-###`
	// series.
	IDPrefix = "DD-"
	// IDFormat is the Printf format for a zero-padded
	// `DD-NNN` identifier.
	IDFormat = "%s%03d"
	// TableHeader is the markdown table header row plus its
	// delimiter row.
	TableHeader = "| ID | Date | Context | Decision |" +
		" Rationale | Consequence | Supporting EV |\n" +
		"|----|------|---------|----------|-----------|" +
		"-------------|---------------|\n"
)
