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
}
