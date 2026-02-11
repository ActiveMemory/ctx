//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// MaxDecisionsToSummarize is the number of recent decisions to include
// in summaries.
const MaxDecisionsToSummarize = 3

// MaxLearningsToSummarize is the number of recent learnings to include
// in summaries.
const MaxLearningsToSummarize = 5

// MaxPreviewLen is the maximum length for preview lines before truncation.
const MaxPreviewLen = 60

// Content detection constants.
const (
	// MinContentLen is the minimum byte length for a file to be considered
	// non-empty by the effectively-empty heuristic.
	MinContentLen = 20
)

// Insight extraction constants.
const (
	// InsightMaxLen is the maximum character length for an extracted insight.
	InsightMaxLen = 150
	// InsightWordBoundaryMin is the minimum cut position when truncating
	// at a word boundary.
	InsightWordBoundaryMin = 100
)

// Recall/export constants.
const (
	// RecallShortIDLen is the truncation length for session IDs in filenames.
	RecallShortIDLen = 8
	// RecallDetailsThreshold is the line count above which tool output is
	// wrapped in a collapsible <details> block.
	RecallDetailsThreshold = 10
)

// Journal site generation constants.
const (
	// JournalPopularityThreshold is the minimum number of entries to
	// mark a topic or key file as "popular" (gets its own dedicated page).
	JournalPopularityThreshold = 2
	// JournalLineWrapWidth is the soft wrap target column for journal
	// content.
	JournalLineWrapWidth = 80
	// JournalMaxRecentSessions is the maximum number of sessions shown
	// in the zensical navigation sidebar.
	JournalMaxRecentSessions = 20
	// JournalMaxNavTitleLen is the maximum title length before
	// truncation in the zensical navigation sidebar.
	JournalMaxNavTitleLen = 40
	// JournalDatePrefixLen is the length of a YYYY-MM-DD date prefix.
	JournalDatePrefixLen = 10
	// JournalMonthPrefixLen is the length of a YYYY-MM month prefix.
	JournalMonthPrefixLen = 7
	// JournalTimePrefixLen is the length of an HH:MM time prefix.
	JournalTimePrefixLen = 5
)
