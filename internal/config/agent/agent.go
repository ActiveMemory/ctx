//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import "time"

// Budget allocation.
const (
	// TaskBudgetPct is the fraction of the token budget allocated to tasks.
	TaskBudgetPct = 0.40
	// ConventionBudgetPct is the fraction of the token budget allocated to conventions.
	ConventionBudgetPct = 0.20
)

// Cooldown configuration.
const (
	// DefaultCooldown is the default cooldown between agent context packet emissions.
	DefaultCooldown = 10 * time.Minute
	// TombstonePrefix is the filename prefix for agent cooldown tombstone files.
	TombstonePrefix = "ctx-agent-"
)

// Scoring configuration.
const (
	// RecencyDaysWeek is the threshold for "recent" entries (0-7 days).
	RecencyDaysWeek = 7
	// RecencyDaysMonth is the threshold for "this month" entries (8-30 days).
	RecencyDaysMonth = 30
	// RecencyDaysQuarter is the threshold for "this quarter" entries (31-90 days).
	RecencyDaysQuarter = 90
	// RecencyScoreWeek is the recency score for entries within a week.
	RecencyScoreWeek = 1.0
	// RecencyScoreMonth is the recency score for entries within a month.
	RecencyScoreMonth = 0.7
	// RecencyScoreQuarter is the recency score for entries within a quarter.
	RecencyScoreQuarter = 0.4
	// RecencyScoreOld is the recency score for entries older than a quarter.
	RecencyScoreOld = 0.2
	// RelevanceMatchCap is the keyword match count that yields maximum relevance (1.0).
	RelevanceMatchCap = 3
)
