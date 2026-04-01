//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package budget implements the token-budgeted context assembly
// algorithm for the agent command.
//
// [AssemblePacket] allocates tokens across five tiers (constitution,
// tasks, conventions, decisions, learnings). [Split] divides
// remaining budget between two scored sections. [FillSection]
// applies two-tier degradation: full entries then title-only
// summaries. [FitItems] and [EstimateSliceTokens] handle
// per-item token accounting.
package budget
