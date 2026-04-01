//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// Check-journal configuration.
const (
	// ThrottleID is the state file name for daily throttle of journal checks.
	ThrottleID = "journal-reminded"
	// ClaudeProjectsSubdir is the relative path under $HOME to
	// the Claude Code projects directory scanned for unimported sessions.
	ClaudeProjectsSubdir = ".claude/projects"
)
