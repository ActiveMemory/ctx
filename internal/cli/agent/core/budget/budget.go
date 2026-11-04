//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/agent/core/score"
	"github.com/ActiveMemory/ctx/internal/config/agent"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
)

// Split divides a token budget between two scored sections.
//
// Each section gets at least 30% of the budget (if content exists).
// The remaining 40% is allocated proportionally to content size.
//
// Parameters:
//   - total: Total tokens to split
//   - a: First section's scored entries
//   - b: Second section's scored entries
//
// Returns:
//   - int: Budget for section a
//   - int: Budget for section b
func Split(total int, a, b []score.Entry) (int, int) {
	if len(a) == 0 && len(b) == 0 {
		return 0, 0
	}
	if len(a) == 0 {
		return 0, total
	}
	if len(b) == 0 {
		return total, 0
	}

	aTokens := TotalEntryTokens(a)
	bTokens := TotalEntryTokens(b)
	totalContent := aTokens + bTokens

	if totalContent == 0 {
		return total / 2, total - total/2
	}

	// If everything fits, give each section what it needs
	if totalContent <= total {
		return aTokens, bTokens
	}

	// Minimum 30% each, proportional split of the rest
	minA := total * agent.SplitMinPct / stats.PercentMultiplier
	minB := total * agent.SplitMinPct / stats.PercentMultiplier
	flex := total - minA - minB

	aProportion := float64(aTokens) / float64(totalContent)
	aFlex := int(float64(flex) * aProportion)

	return minA + aFlex, total - (minA + aFlex)
}

// FillSection selects scored entries to fill a budget,
// with graceful degradation.
//
// Includes full entries by score order until ~80% of the budget is consumed.
// Remaining entries get title-only summaries.
//
// Parameters:
//   - entries: Scored entries sorted by score descending
//   - budget: Token budget for this section
//
// Returns:
//   - []string: Full entry bodies that fit in the budget
//   - []string: Title-only summaries for entries that didn't fit
func FillSection(entries []score.Entry, budget int) ([]string, []string) {
	if len(entries) == 0 || budget <= 0 {
		return nil, nil
	}

	fullBudget := budget * agent.FullEntryPct / stats.PercentMultiplier
	used := 0
	var full []string
	var summaries []string

	for i := range entries {
		if entries[i].Score == 0.0 {
			// Superseded entries: skip entirely
			continue
		}
		body := entries[i].BlockContent()
		tokens := entries[i].Tokens
		if used+tokens <= fullBudget {
			full = append(full, body)
			used += tokens
		} else {
			// Title-only summary
			summaries = append(summaries, entries[i].Entry.Title)
		}
	}

	return full, summaries
}

// FitHandover trims a handover body on a line boundary to fit within budget.
//
// Parameters:
//   - body: Raw handover markdown body
//   - budget: Maximum token budget for the handover section
//
// Returns:
//   - string: Handover body that fits within budget, or empty string
func FitHandover(body string, budget int) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" || budget <= 0 {
		return ""
	}
	if ctxToken.EstimateString(trimmed) <= budget {
		return trimmed
	}
	lines := strings.Split(trimmed, token.NewlineLF)
	suffix := desc.Text(text.DescKeyAgentHandoverTruncated)
	suffixTokens := ctxToken.EstimateString(suffix)
	var kept []string
	used := suffixTokens
	for _, line := range lines {
		t := ctxToken.EstimateString(line) + 1
		if used+t > budget {
			break
		}
		kept = append(kept, line)
		used += t
	}
	if len(kept) == 0 {
		return ""
	}
	kept = append(kept, suffix)
	return strings.Join(kept, token.NewlineLF)
}
