//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

// EstimateTokens provides a rough token count estimate for content.
//
// Uses a simple heuristic of ~4 characters per token for English text.
// This is a conservative estimate for Claude/GPT-style tokenizers that
// tends to slightly overestimate, which is safer for budgeting.
//
// Parameters:
//   - content: Byte slice to estimate tokens for
//
// Returns:
//   - int: Estimated token count (0 for empty content)
func EstimateTokens(content []byte) int {
	if len(content) == 0 {
		return 0
	}
	// Rough estimate: 1 token per 4 characters
	// This tends to slightly overestimate, which is safer for budgeting
	return (len(content) + 3) / 4
}

// EstimateTokensString estimates tokens for a string.
//
// Convenience wrapper around EstimateTokens for string input.
//
// Parameters:
//   - s: String to estimate tokens for
//
// Returns:
//   - int: Estimated token count
func EstimateTokensString(s string) int {
	return EstimateTokens([]byte(s))
}
