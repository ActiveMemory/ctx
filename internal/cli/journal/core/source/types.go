//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

// Opts holds all flags for the source subcommand.
//
// Fields:
//   - ShowID: Session ID or slug to show
//   - Latest: Show most recent session
//   - Full: Show full conversation (not preview)
//   - Limit: Maximum sessions to list
//   - Project: Filter by project name
//   - Tool: Filter by tool name
//   - Since: Show sessions after this date
//   - Until: Show sessions before this date
//   - AllProjects: Include all projects
//   - Turns: 1-based inclusive turn range START:END
//   - Grep: Case-insensitive regex pattern to search turns
//   - MaxChars: Maximum characters to emit in bounded show mode
type Opts struct {
	ShowID      string
	Latest      bool
	Full        bool
	Limit       int
	Project     string
	Tool        string
	Since       string
	Until       string
	AllProjects bool
	Turns       string
	Grep        string
	MaxChars    int
}

// TurnMatch represents a single turn matching a --grep query.
//
// Fields:
//   - Turn: 1-based turn index in session.Messages
//   - Role: Display role label (User, Assistant, Tool Output)
//   - Timestamp: Formatted HH:MM:SS timestamp
//   - Snippet: Context snippet around the match
type TurnMatch struct {
	Turn      int
	Role      string
	Timestamp string
	Snippet   string
}
