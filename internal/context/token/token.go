//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

import cfgToken "github.com/ActiveMemory/ctx/internal/config/token"

// Estimate provides a rough token count estimate for content.
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
func Estimate(content []byte) int {
	if len(content) == 0 {
		return 0
	}
	// Rough estimate: 1 token per CharsPerToken characters.
	// Ceiling division ensures slight overestimate, safer for budgeting.
	return (len(content) + cfgToken.CharsPerToken - 1) / cfgToken.CharsPerToken
}

// EstimateString estimates tokens for a string.
//
// Convenience wrapper around Estimate for string input.
//
// Parameters:
//   - s: String to estimate tokens for
//
// Returns:
//   - int: Estimated token count
func EstimateString(s string) int {
	return Estimate([]byte(s))
}
