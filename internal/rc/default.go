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

// DefaultEntryCountLearnings is the entry count threshold for LEARNINGS.md.
// Learnings are situational; many become stale. Warn above this count.
const DefaultEntryCountLearnings = 30

// DefaultEntryCountDecisions is the entry count threshold for DECISIONS.md.
// Decisions are more durable but still compound. Warn above this count.
const DefaultEntryCountDecisions = 20

// DefaultConventionLineCount is the line count threshold for CONVENTIONS.md.
// Conventions lack dated entry headers, so line count is used instead.
const DefaultConventionLineCount = 200
