//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import "github.com/ActiveMemory/ctx/internal/config/runtime"

// Aliases re-exported from config/runtime for use within rc.
const (
	DefaultTokenBudget         = runtime.DefaultTokenBudget
	DefaultArchiveAfterDays    = runtime.DefaultArchiveAfterDays
	DefaultEntryCountLearnings = runtime.DefaultEntryCountLearnings
	DefaultEntryCountDecisions = runtime.DefaultEntryCountDecisions
	DefaultConventionLineCount = runtime.DefaultConventionLineCount
	DefaultInjectionTokenWarn  = runtime.DefaultInjectionTokenWarn
	DefaultContextWindow       = runtime.DefaultContextWindow
	DefaultTaskNudgeInterval   = runtime.DefaultTaskNudgeInterval
	DefaultKeyRotationDays     = runtime.DefaultKeyRotationDays
)
