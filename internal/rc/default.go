//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/runtime"
)

// Aliases re-exported from config/runtime for use within rc.
const (
	// DefaultTokenBudget is the default agent token budget.
	DefaultTokenBudget = runtime.DefaultTokenBudget
	// DefaultArchiveAfterDays is the default task archive age.
	DefaultArchiveAfterDays = runtime.DefaultArchiveAfterDays
	// DefaultEntryCountLearnings is the max learnings shown.
	DefaultEntryCountLearnings = runtime.DefaultEntryCountLearnings
	// DefaultEntryCountDecisions is the max decisions shown.
	DefaultEntryCountDecisions = runtime.DefaultEntryCountDecisions
	// DefaultConventionLineCount is the max convention lines.
	DefaultConventionLineCount = runtime.DefaultConventionLineCount
	// DefaultInjectionTokenWarn is the injection warn threshold.
	DefaultInjectionTokenWarn = runtime.DefaultInjectionTokenWarn
	// DefaultContextWindow is the default context window size.
	DefaultContextWindow = runtime.DefaultContextWindow
	// DefaultTaskNudgeInterval is the default nudge interval.
	DefaultTaskNudgeInterval = runtime.DefaultTaskNudgeInterval
	// DefaultKeyRotationDays is the default key rotation age.
	DefaultKeyRotationDays = runtime.DefaultKeyRotationDays
	// DefaultStaleAgeDays is the default stale entry age.
	DefaultStaleAgeDays = runtime.DefaultStaleAgeDays
)

// Hooks & Steering defaults.
const (
	// DefaultSteeringDir is the default steering directory path.
	DefaultSteeringDir = cfgDir.DefaultSteeringPath
	// DefaultHooksDir is the default hooks directory path.
	DefaultHooksDir = cfgDir.DefaultHooksPath
	// DefaultHookTimeout is the default per-hook execution timeout in seconds.
	DefaultHookTimeout = 10
)
