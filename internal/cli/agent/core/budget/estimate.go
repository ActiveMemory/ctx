//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

import (
	"github.com/ActiveMemory/ctx/internal/cli/agent/core/score"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
)

// FitItems returns items that fit within a token budget.
//
// Items are included in order until the budget would be exceeded.
//
// Parameters:
//   - items: String items to include
//   - budget: Maximum token budget
//
// Returns:
//   - []string: Items that fit within the budget
func FitItems(items []string, budget int) []string {
	if len(items) == 0 {
		return nil
	}
	used := 0
	var result []string
	for _, item := range items {
		tokens := ctxToken.EstimateString(item)
		if used+tokens > budget {
			break
		}
		result = append(result, item)
		used += tokens
	}
	// Always include at least one item if there are any
	if len(result) == 0 && len(items) > 0 {
		result = append(result, items[0])
	}
	return result
}

// EstimateSliceTokens sums token estimates for a string slice.
//
// Parameters:
//   - items: Strings to estimate
//
// Returns:
//   - int: Total estimated tokens
func EstimateSliceTokens(items []string) int {
	total := 0
	for _, item := range items {
		total += ctxToken.EstimateString(item)
	}
	return total
}

// TotalEntryTokens sums pre-computed token counts for scored entries.
//
// Parameters:
//   - entries: Scored entries with token estimates
//
// Returns:
//   - int: Total tokens
func TotalEntryTokens(entries []score.Entry) int {
	total := 0
	for _, e := range entries {
		total += e.Tokens
	}
	return total
}
