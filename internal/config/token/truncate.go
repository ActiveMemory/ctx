//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

// Display truncation lengths for CLI and MCP output.
const (
	// TruncateLen is the max display length for task text in compact output.
	TruncateLen = 50
	// TruncateContentLen is the max display length for pending update content.
	TruncateContentLen = 60
)

// Token estimation constants.
const (
	// CharsPerToken is the heuristic character-to-token ratio.
	// Deliberately overestimates (safer for budgeting).
	CharsPerToken = 4
)
