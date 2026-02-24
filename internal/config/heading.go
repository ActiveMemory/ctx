//	/    ctx:                         https://ctx.ist
//
// ,'`./    do you remember?
//
//	`.,'\
//	  \    Copyright 2026-present Context contributors.
//	                SPDX-License-Identifier: Apache-2.

package config

// Learnings
const (
	// HeadingLearningStart is the Markdown heading for entries in LEARNINGS.md
	HeadingLearningStart = "## ["
	// HeadingLearnings is the Markdown heading for LEARNINGs.md
	HeadingLearnings = "# Learnings"
	// ColumnLearning is the singular column header for learning index tables.
	ColumnLearning = "Learning"
)

// Task sections in TASKS.md
const (
	// HeadingInProgress is the section heading for in-progress tasks.
	HeadingInProgress = "## In Progress"
	// HeadingNextUp is the section heading for upcoming tasks.
	HeadingNextUp = "## Next Up"
	// HeadingCompleted is the section heading for completed tasks.
	HeadingCompleted = "## Completed"
	// HeadingArchivedTasks is the heading for archived task files.
	HeadingArchivedTasks = "# Archived Tasks"
)

// Decisions
const (
	// HeadingDecisions is the Markdown heading for DECISIONS.md
	HeadingDecisions = "# Decisions"
	// ColumnDecision is the singular column header for decision index tables.
	ColumnDecision = "Decision"
)

// Journal index headings
const (
	// JournalHeadingSessionJournal is the main journal index title.
	JournalHeadingSessionJournal = "# Session Journal"
	// JournalHeadingTopics is the topics index title.
	JournalHeadingTopics = "# Topics"
	// JournalHeadingPopularTopics is the popular topics section heading.
	JournalHeadingPopularTopics = "## Popular Topics"
	// JournalHeadingLongtailTopics is the long-tail topics section heading.
	JournalHeadingLongtailTopics = "## Long-tail Topics"
	// JournalHeadingKeyFiles is the key files index title.
	JournalHeadingKeyFiles = "# Key Files"
	// JournalHeadingFrequentlyTouched is the popular key files section heading.
	JournalHeadingFrequentlyTouched = "## Frequently Touched"
	// JournalHeadingSingleSession is the single-session key files section heading.
	JournalHeadingSingleSession = "## Single Session"
	// JournalHeadingSessionTypes is the session types index title.
	JournalHeadingSessionTypes = "# Session Types"
	// JournalHeadingSuggestions is the suggestions section heading.
	JournalHeadingSuggestions = "## Suggestions"
	// JournalHeadingRecentSessions is the nav section title for recent entries.
	JournalHeadingRecentSessions = "Recent Sessions"
)

// Recall/export headings used in journal entry Markdown.
const (
	// RecallHeadingSummary is the summary section heading in journal entries.
	RecallHeadingSummary = "## Summary"
	// RecallHeadingToolUsage is the tool usage section heading.
	RecallHeadingToolUsage = "## Tool Usage"
	// RecallHeadingConversation is the conversation section heading.
	RecallHeadingConversation = "## Conversation"
)

// Load command headings
const (
	// LoadHeadingContext is the top-level heading for assembled context output.
	LoadHeadingContext = "# Context"
)

// Loop command headings
const (
	// LoopHeadingStart is the heading shown after loop script generation.
	LoopHeadingStart = "To start the loop:"
)

// Journal navigation labels used in the zensical site nav bar.
const (
	// JournalLabelHome is the nav label for the index page.
	JournalLabelHome = "Home"
	// JournalLabelTopics is the nav label for the topics index.
	JournalLabelTopics = "Topics"
	// JournalLabelFiles is the nav label for the key files index.
	JournalLabelFiles = "Files"
	// JournalLabelTypes is the nav label for the session types index.
	JournalLabelTypes = "Types"
)
