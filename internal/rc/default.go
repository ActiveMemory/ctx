//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import "github.com/ActiveMemory/ctx/internal/config"

// Aliases for backward compatibility with external references.
const (
	DefaultTokenBudget         = config.DefaultRcTokenBudget
	DefaultArchiveAfterDays    = config.DefaultRcArchiveAfterDays
	DefaultEntryCountLearnings = config.DefaultRcEntryCountLearnings
	DefaultEntryCountDecisions = config.DefaultRcEntryCountDecisions
	DefaultConventionLineCount = config.DefaultRcConventionLineCount
	DefaultInjectionTokenWarn  = config.DefaultRcInjectionTokenWarn
	DefaultContextWindow       = config.DefaultRcContextWindow
	DefaultTaskNudgeInterval   = config.DefaultRcTaskNudgeInterval
	DefaultKeyRotationDays     = config.DefaultRcKeyRotationDays
)
