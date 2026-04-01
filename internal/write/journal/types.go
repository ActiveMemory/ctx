//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// SessionInfo holds pre-formatted session metadata for display.
//
// Fields:
//   - Slug: URL-safe session identifier
//   - ID: Full session UUID
//   - Tool: AI tool name (e.g. "claude-code")
//   - Project: Project name
//   - Branch: Git branch (empty to omit)
//   - Model: Model ID (empty to omit)
//   - Started: Formatted start time
//   - Duration: Formatted duration
//   - Turns: Number of conversation turns
//   - Messages: Total message count
//   - TokensIn: Formatted input token count
//   - TokensOut: Formatted output token count
//   - TokensAll: Formatted total token count
type SessionInfo struct {
	Slug      string
	ID        string
	Tool      string
	Project   string
	Branch    string // empty to omit
	Model     string // empty to omit
	Started   string
	Duration  string
	Turns     int
	Messages  int
	TokensIn  string
	TokensOut string
	TokensAll string
}
