//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package score computes relevance scores for context entries to
// prioritize budget allocation.
//
// [Recency] scores by age (7d=1.0, 30d=0.7, 90d=0.4, older=0.2).
// [Relevance] scores by keyword overlap with active tasks (0.0-1.0).
// [Score] combines both into a 0.0-2.0 range. [All] scores a
// batch of entries. [ExtractTaskKeywords] builds the keyword set
// from active task text.
package score
