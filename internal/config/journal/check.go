//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// Check-journal configuration.
const (
	// CheckJournalThrottleID is the state file name for daily throttle of journal checks.
	CheckJournalThrottleID = "journal-reminded"
	// CheckJournalClaudeProjectsSubdir is the relative path under $HOME to
	// the Claude Code projects directory scanned for unexported sessions.
	CheckJournalClaudeProjectsSubdir = ".claude/projects"
)
