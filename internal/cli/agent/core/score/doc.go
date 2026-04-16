//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package score computes **per-entry relevance scores** so
// the budgeted context-assembly algorithm in
// [internal/cli/agent/core/budget] can decide which
// decisions, learnings, and conventions to inject when there
// is not enough budget to inject all of them.
//
// The score is a deliberately simple two-component number:
// **recency** plus **relevance to current work**. Either
// component alone produces a poor ranking; together they
// approximate "what would a helpful colleague pull off the
// shelf?".
//
// # The Two Components
//
//   - **[Recency](entry)** — bucketed by age:
//
//     ≤  7 days   → 1.0
//     ≤ 30 days   → 0.7
//     ≤ 90 days   → 0.4
//     older       → 0.2
//
//     Buckets (rather than a continuous decay) keep the
//     ordering stable across small input shifts and make
//     the scoring trivially debuggable.
//
//   - **[Relevance](entry, taskKeywords)** — fraction of
//     the entry's salient tokens that overlap with
//     [ExtractTaskKeywords](activeTasks). Range 0.0–1.0.
//     Stop words come from the embedded list in
//     [internal/assets/read/lookup.StopWords].
//
// [Score](entry, taskKeywords) sums the two for a 0.0–2.0
// composite. [All](entries, taskKeywords) is the bulk
// scorer that returns parallel slices for the budget
// allocator.
//
// # Why Bucketed Recency
//
// A continuous exponential decay would be technically
// purer but produces "score jitter" — entries reorder
// minute-to-minute as their ages cross the decimal
// boundary. Bucketed recency means an entry's relative
// rank only changes when it crosses a real threshold
// (week, month, quarter), which is the cadence at which
// users actually expect their context to age.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never race.
package score
