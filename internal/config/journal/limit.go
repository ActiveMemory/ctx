//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// Journal site generation constants.
const (
	// PopularityThreshold is the minimum number of entries to
	// mark a topic or key file as "popular" (gets its own dedicated page).
	PopularityThreshold = 2
	// LineWrapWidth is the soft wrap target column for journal
	// content.
	LineWrapWidth = 80
	// MaxRecentSessions is the maximum number of sessions shown
	// in the zensical navigation sidebar.
	MaxRecentSessions = 20
	// MaxNavTitleLen is the maximum title length before
	// truncation in the zensical navigation sidebar.
	MaxNavTitleLen = 40
	// DatePrefixLen is the length of a YYYY-MM-DD date prefix.
	DatePrefixLen = 10
	// MonthPrefixLen is the length of a YYYY-MM month prefix.
	MonthPrefixLen = 7
	// TimePrefixLen is the length of an HH:MM time prefix.
	TimePrefixLen = 5
	// MaxTitleLen is the maximum character length for a journal title.
	// Keeps H1 headings and link text on a single line (below wrap width).
	MaxTitleLen = 75
	// ShortIDLen is the truncation length for session IDs in filenames.
	ShortIDLen = 8
	// DetailsThreshold is the line count above which tool output is
	// wrapped in a collapsible <details> block.
	DetailsThreshold = 10
	// DefaultRecallListLimit is the default number of sessions
	// shown by recall list.
	DefaultRecallListLimit = 20
	// MultipartSuffix is the separator between the base slug and the part
	// number in multipart journal filenames (e.g. "slug-p2.md").
	MultipartSuffix = "-p"
)
