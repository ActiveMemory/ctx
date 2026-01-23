//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

// EstimateTokens provides a rough token count estimate for content.
// Uses a simple heuristic: ~4 characters per token for English text.
// This is a conservative estimate for Claude/GPT-style tokenizers.
func EstimateTokens(content []byte) int {
	if len(content) == 0 {
		return 0
	}
	// Rough estimate: 1 token per 4 characters
	// This tends to slightly overestimate, which is safer for budgeting
	return (len(content) + 3) / 4
}

// EstimateTokensString estimates tokens for a string.
func EstimateTokensString(s string) int {
	return EstimateTokens([]byte(s))
}
