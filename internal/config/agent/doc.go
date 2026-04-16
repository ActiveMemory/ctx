//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agent centralizes budget, cooldown, and scoring
// constants that govern how AI agents consume and prioritize
// context files.
//
// The ctx agent command compresses project context (tasks,
// conventions, decisions, learnings) into a token-budgeted
// packet. This package defines the knobs that control that
// compression pipeline.
//
// # Budget Allocation
//
// Token budgets are split across context sections using
// percentage-based constants:
//
//   - [TaskBudgetPct] reserves 40% of the budget for tasks.
//   - [ConventionBudgetPct] reserves 20% for conventions.
//   - [SplitMinPct] ensures each section receives at least
//     30% of its allocated share.
//   - [FullEntryPct] dedicates 80% of each section budget
//     to full entries; the remainder becomes title-only
//     summaries.
//
// These ratios keep task context dominant while giving
// conventions enough room for actionable detail.
//
// # Cooldown
//
// Agent packets are expensive to generate. The cooldown
// system prevents redundant emissions:
//
//   - [DefaultCooldown] enforces a 10-minute gap between
//     successive agent runs.
//   - [TombstonePrefix] names the state files that track
//     the last emission timestamp.
//
// # Recency Scoring
//
// Entries are ranked by age to surface recent work first.
// The scoring tiers decay from 1.0 (within a week) down
// to 0.2 (older than a quarter):
//
//   - [RecencyDaysWeek] / [RecencyScoreWeek]   (0-7 days)
//   - [RecencyDaysMonth] / [RecencyScoreMonth] (8-30 days)
//   - [RecencyDaysQuarter] / [RecencyScoreQuarter] (31-90)
//   - [RecencyScoreOld] for everything older.
//
// [RelevanceMatchCap] caps keyword-match relevance at 3
// hits, preventing a single heavily-tagged entry from
// consuming the entire budget.
//
// # Why Centralized
//
// Budget ratios and scoring thresholds are tuned together.
// Keeping them in one package lets callers adjust the
// pipeline without hunting across multiple files.
package agent
