//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package token estimates LLM token counts from byte
// content using a characters-per-token heuristic.
//
// # Estimation
//
// Estimate takes a byte slice and returns a rough
// token count. It uses approximately 4 characters per
// token for English text, which is a conservative
// estimate for Claude/GPT-style tokenizers. Ceiling
// division ensures slight overestimation, which is
// safer for budget enforcement.
//
//	tokens := token.Estimate(content)
//
// EstimateString is a convenience wrapper that accepts
// a string instead of a byte slice.
//
//	tokens := token.EstimateString(text)
//
// # Budget Enforcement
//
// Token estimates are used throughout ctx to enforce
// context budgets. When assembling the agent packet,
// files are added in priority order until the budget
// is exhausted. Overestimation ensures the assembled
// output stays within the actual token limit.
//
// # Accuracy
//
// The heuristic is intentionally simple. Real
// tokenizers produce different counts depending on
// vocabulary, but the 4-chars-per-token ratio is a
// well-known approximation that errs on the safe side.
package token
