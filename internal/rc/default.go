//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// DefaultTokenBudget is the default token budget when not configured.
const DefaultTokenBudget = 8000

// DefaultArchiveAfterDays is the default days before archiving.
const DefaultArchiveAfterDays = 7

// DefaultArchiveKnowledgeAfterDays is the default days before archiving
// decisions and learnings.
const DefaultArchiveKnowledgeAfterDays = 90

// DefaultArchiveKeepRecent is the default number of recent entries to keep
// when archiving decisions and learnings.
const DefaultArchiveKeepRecent = 5

// DefaultEntryCountLearnings is the entry count threshold for LEARNINGS.md.
// Learnings are situational; many become stale. Warn above this count.
const DefaultEntryCountLearnings = 30

// DefaultEntryCountDecisions is the entry count threshold for DECISIONS.md.
// Decisions are more durable but still compound. Warn above this count.
const DefaultEntryCountDecisions = 20
