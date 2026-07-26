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
	// DefaultArchiveAfterDays is the default days before
	// archiving completed tasks.
	DefaultArchiveAfterDays = 7
	// DefaultEntryCountLearnings is the entry count threshold for LEARNINGS.md.
	DefaultEntryCountLearnings = 30
	// DefaultEntryCountDecisions is the entry count threshold for DECISIONS.md.
	DefaultEntryCountDecisions = 20
	// DefaultConventionSectionCount is the staged-section count threshold
	// for CONVENTIONS.md — the foldability watermark unit (M5). It replaces
	// the former line-count measure: a section count is what the digest
	// pass acts on, so it is the honest "should I fold?" signal.
	DefaultConventionSectionCount = 12
	// DefaultThemePageByteCeiling is the byte ceiling above which a page
	// (a root or a theme file) is flagged heavy (M5). Bytes, not lines:
	// a line hides 10 or 200 characters. Past this, the advice is split or
	// extract-to-tooling — an LLM is a poor linter.
	DefaultThemePageByteCeiling = 65536
	// DefaultInjectionTokenWarn is the token threshold for
	// oversize injection warning.
	DefaultInjectionTokenWarn = 15000
	// DefaultContextWindow is the default context window size in tokens.
	DefaultContextWindow = 200000
	// DefaultTaskNudgeInterval is the Edit/Write calls between
	// task completion nudges.
	DefaultTaskNudgeInterval = 5
	// DefaultKeyRotationDays is the days before encryption key rotation nudge.
	DefaultKeyRotationDays = 90
	// DefaultStaleAgeDays is the days before a context file is
	// flagged as stale by drift detection.
	DefaultStaleAgeDays = 30
	// DefaultPruneDays is the default age in days for state file pruning.
	DefaultPruneDays = 7
)
