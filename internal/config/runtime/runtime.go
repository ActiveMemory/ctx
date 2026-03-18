//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package runtime

// Runtime configuration defaults (overridable via .ctxrc).
const (
	// DefaultTokenBudget is the default token budget for context assembly.
	DefaultTokenBudget = 8000
	// DefaultArchiveAfterDays is the default days before archiving completed tasks.
	DefaultArchiveAfterDays = 7
	// DefaultEntryCountLearnings is the entry count threshold for LEARNINGS.md.
	DefaultEntryCountLearnings = 30
	// DefaultEntryCountDecisions is the entry count threshold for DECISIONS.md.
	DefaultEntryCountDecisions = 20
	// DefaultConventionLineCount is the line count threshold for CONVENTIONS.md.
	DefaultConventionLineCount = 200
	// DefaultInjectionTokenWarn is the token threshold for oversize injection warning.
	DefaultInjectionTokenWarn = 15000
	// DefaultContextWindow is the default context window size in tokens.
	DefaultContextWindow = 200000
	// DefaultTaskNudgeInterval is the Edit/Write calls between task completion nudges.
	DefaultTaskNudgeInterval = 5
	// DefaultKeyRotationDays is the days before encryption key rotation nudge.
	DefaultKeyRotationDays = 90
	// DefaultStaleAgeDays is the days before a context file is flagged as stale by drift detection.
	DefaultStaleAgeDays = 30
)
