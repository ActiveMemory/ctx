//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agent implements the "ctx agent" command for generating
// AI-ready context packets with budget-aware content selection.
//
// The agent command reads context files from .context/ and produces
// a token-budgeted output optimized for AI consumption. Output can
// be in Markdown (default) or JSON format.
//
// # Budget-Aware Assembly
//
// The core algorithm in [assembleBudgetPacket] allocates tokens across
// five tiers, filling higher tiers before lower ones. Each tier has a
// budget cap expressed as a percentage of the total budget:
//
//   - Tier 1 (uncapped): Constitution rules, file read order, and the
//     behavioral instruction. These are always included in full because
//     they are small and inviolable. Their token cost is subtracted
//     from the remaining budget before any other allocation.
//   - Tier 2 (40%): Active tasks from TASKS.md. Tasks are included in
//     file order until the 40% cap is reached. At least one task is
//     always included regardless of budget.
//   - Tier 3 (20%): Convention items from CONVENTIONS.md. Same
//     fill-until-cap strategy as tasks.
//   - Tier 4+5 (remaining): Decisions and learnings share whatever
//     budget remains after tiers 1–3. The split is proportional to
//     content size with a 30% minimum guarantee for each section
//     (see [splitBudget]).
//
// # Entry Scoring
//
// Decisions and learnings are scored before budget fitting. Each entry
// receives a combined score in the range 0.0–2.0, computed by
// [scoreEntry] as:
//
//	score = recencyScore + relevanceScore
//
// Recency ([recencyScore]) uses age brackets: entries from the last
// 7 days score 1.0, 8–30 days score 0.7, 31–90 days score 0.4, and
// older entries score 0.2. This ensures recent context is preferred
// without completely excluding older entries that may still be
// relevant.
//
// Relevance ([relevanceScore]) counts keyword overlap between the
// entry's text and keywords extracted from active tasks. Keywords are
// extracted by [extractTaskKeywords], which splits task text on
// whitespace/punctuation, lowercases, removes stop words, and
// deduplicates. The overlap count is normalized to 1.0 at 3+ matches.
//
// Superseded entries (those containing "~~Superseded") always receive
// a score of 0.0 and are excluded from output unless the budget
// accommodates everything.
//
// # Graceful Degradation
//
// When scored entries exceed their budget allocation, [fillSection]
// applies a two-tier degradation strategy:
//
//  1. Full entries are included in score order until ~80% of the
//     section budget is consumed.
//  2. Remaining entries are reduced to title-only summaries, grouped
//     under an "Also Noted" section in the output.
//
// This ensures the packet always communicates what context exists,
// even when it cannot include the full rationale.
//
// # Token Estimation
//
// All budget accounting uses [context.EstimateTokensString], which
// applies a len/4 heuristic. This deliberately overestimates, which
// is correct for budgeting — it is better to include slightly less
// than to overflow the context window.
package agent
