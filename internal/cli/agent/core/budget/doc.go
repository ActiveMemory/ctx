//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package budget implements the **token-budgeted context
// assembly algorithm** behind `ctx agent`. Given a token budget
// (`--budget N`, default 8000) and a loaded [entity.Context],
// it produces an AI-ready packet that maximizes information
// density without exceeding the budget — the single most
// performance-sensitive operation in ctx because it runs at
// the head of every prompt in tool integrations that use the
// hook+MCP pipeline.
//
// # The Seven-Tier Allocation
//
// [AssemblePacket] walks the seven content tiers in priority
// order, each with its own share of the budget:
//
//  1. **CONSTITUTION** — always full; never truncated.
//  2. **TASKS** — current and pending work.
//  3. **CONVENTIONS** — coding patterns the AI must follow.
//  4. **DECISIONS** — index table, then full entries as
//     budget permits.
//  5. **LEARNINGS** — same shape as decisions.
//  6. **STEERING** — matched files for this prompt.
//  7. **SKILL** — bundled instructions if a skill matched.
//
// Lower tiers see whatever budget the higher tiers leave
// behind. The constitution invariant — "context loading is the
// first step of every session" — translates into "the
// constitution is always in the packet, no exceptions".
//
// # Two-Tier Degradation
//
// [FillSection] handles the per-section degradation: when full
// entries do not fit, it **falls back to title-only summaries**
// (the index-table form) so the AI still sees that an entry
// exists and can request it by ID. The degradation point is
// chosen to maximize the count of entries the AI sees,
// trading depth for breadth.
//
// # Splitting Between Two Sections
//
// [Split] divides the remaining budget between two scored
// sections (typically DECISIONS vs LEARNINGS) using a
// score-weighted ratio: a section with twice the relevance
// score gets twice the budget share. Score comes from
// [internal/cli/agent/core/score]; budget enforces.
//
// # Token Accounting
//
// [EstimateSliceTokens] is the rough-but-stable estimator
// used throughout: ~4 chars per token for English Markdown,
// with adjustments for code-fence-heavy content. It is not
// the exact count the model will see but is consistent
// enough to keep the assembled packet under budget.
// [FitItems] is the greedy item picker: takes the highest-
// scored items first, stops when the next one would push
// the running total over budget.
//
// # Render Path
//
// [render.go] formats the assembled tiers into the final
// markdown packet with section headers, separators, and
// the read-order preamble the AI uses to navigate the
// content. [out.go] writes the packet to stdout (or to the
// MCP response, depending on caller).
//
// # Concurrency
//
// All functions are pure data transformations over the
// loaded context. Concurrent callers never race; the
// algorithm holds no module-level state.
//
// # Related Packages
//
//   - [internal/cli/agent]            — the `ctx agent`
//     CLI surface that drives this package.
//   - [internal/cli/agent/core/score] — the relevance
//     scorer that ranks entries before [FillSection]
//     consumes them.
//   - [internal/context/load]         — produces the
//     [entity.Context] this package consumes.
//   - [internal/index]                — produces the
//     index-table form used as the fallback.
//   - [internal/steering]             — supplies matched
//     steering files for the steering tier.
//   - [internal/skill]                — supplies matched
//     skill bundles for the skill tier.
package budget
