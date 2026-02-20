//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/index"
)

// Budget tier allocation percentages.
const (
	taskBudgetPct       = 0.40
	conventionBudgetPct = 0.20
)

// assembledPacket holds the budget-aware output sections ready for rendering.
//
// Fields:
//   - ReadOrder: File paths in recommended reading order
//   - Constitution: Constitution rules (always included)
//   - Tasks: Active tasks (budget-capped)
//   - Conventions: Convention items (budget-capped)
//   - Decisions: Full decision entries (scored, budget-fitted)
//   - Learnings: Full learning entries (scored, budget-fitted)
//   - Summaries: Title-only summaries of entries that didn't fit
//   - Instruction: Behavioral instruction for the agent
//   - Budget: Requested token budget
//   - TokensUsed: Actual tokens consumed by the packet
type assembledPacket struct {
	ReadOrder    []string
	Constitution []string
	Tasks        []string
	Conventions  []string
	Decisions    []string
	Learnings    []string
	Summaries    []string
	Instruction  string
	Budget       int
	TokensUsed   int
}

// assembleBudgetPacket builds a context packet respecting the token budget.
//
// Allocation tiers:
//   - Tier 1 (always): constitution, read order, instruction
//   - Tier 2 (40%): active tasks
//   - Tier 3 (20%): conventions
//   - Tier 4+5 (remaining): decisions + learnings, scored by relevance
//
// Parameters:
//   - ctx: Loaded context containing the files
//   - budget: Token budget to respect
//
// Returns:
//   - *assembledPacket: Assembled packet within budget
func assembleBudgetPacket(ctx *context.Context, budget int) *assembledPacket {
	now := time.Now()
	pkt := &assembledPacket{
		Budget: budget,
		Instruction: "Before starting work, confirm to the user: " +
			"\"I have read the required context files and " +
			"I'm following project conventions.\"",
	}

	remaining := budget

	// Tier 1: Always included (constitution, read order, instruction)
	pkt.ReadOrder = getReadOrder(ctx)
	pkt.Constitution = extractConstitutionRules(ctx)

	tier1Tokens := estimateSliceTokens(pkt.ReadOrder) +
		estimateSliceTokens(pkt.Constitution) +
		context.EstimateTokensString(pkt.Instruction)
	remaining -= tier1Tokens

	if remaining <= 0 {
		pkt.TokensUsed = tier1Tokens
		return pkt
	}

	// Tier 2: Tasks (up to 40% of original budget)
	taskCap := int(float64(budget) * taskBudgetPct)
	allTasks := extractActiveTasks(ctx)
	pkt.Tasks = fitItemsInBudget(allTasks, taskCap)
	taskTokens := estimateSliceTokens(pkt.Tasks)
	remaining -= taskTokens

	if remaining <= 0 {
		pkt.TokensUsed = budget - remaining
		return pkt
	}

	// Tier 3: Conventions (up to 20% of original budget)
	convCap := int(float64(budget) * conventionBudgetPct)
	allConventions := extractAllConventions(ctx)
	pkt.Conventions = fitItemsInBudget(allConventions, convCap)
	convTokens := estimateSliceTokens(pkt.Conventions)
	remaining -= convTokens

	if remaining <= 0 {
		pkt.TokensUsed = budget - remaining
		return pkt
	}

	// Extract keywords from tasks for relevance scoring
	keywords := extractTaskKeywords(pkt.Tasks)

	// Tier 4+5: Decisions + Learnings (share remaining budget)
	decisionBlocks := parseEntryBlocks(ctx, config.FileDecision)
	learningBlocks := parseEntryBlocks(ctx, config.FileLearning)

	scoredDecisions := scoreEntries(decisionBlocks, keywords, now)
	scoredLearnings := scoreEntries(learningBlocks, keywords, now)

	// Split remaining budget: proportional to content size, minimum 30% each
	decTokens, learnTokens := splitBudget(
		remaining, scoredDecisions, scoredLearnings,
	)

	pkt.Decisions, pkt.Summaries = fillSection(scoredDecisions, decTokens)

	var learnSummaries []string
	pkt.Learnings, learnSummaries = fillSection(scoredLearnings, learnTokens)
	pkt.Summaries = append(pkt.Summaries, learnSummaries...)

	pkt.TokensUsed = tier1Tokens + taskTokens + convTokens +
		estimateSliceTokens(pkt.Decisions) +
		estimateSliceTokens(pkt.Learnings) +
		estimateSliceTokens(pkt.Summaries)

	return pkt
}

// extractAllConventions extracts all bullet items from CONVENTIONS.md
// (not limited to 5 like the old implementation).
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: All convention bullet items; nil if the file is not found
func extractAllConventions(ctx *context.Context) []string {
	if f := ctx.File(config.FileConvention); f != nil {
		return extractBulletItems(string(f.Content), 1000)
	}
	return nil
}

// parseEntryBlocks parses a context file into entry blocks.
//
// Parameters:
//   - ctx: Loaded context
//   - fileName: Name of the file to parse (e.g., config.FileDecision)
//
// Returns:
//   - []index.EntryBlock: Parsed entry blocks; nil if file not found
func parseEntryBlocks(ctx *context.Context, fileName string) []index.EntryBlock {
	if f := ctx.File(fileName); f != nil {
		return index.ParseEntryBlocks(string(f.Content))
	}
	return nil
}

// splitBudget divides a token budget between two scored sections.
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
func splitBudget(total int, a, b []ScoredEntry) (int, int) {
	if len(a) == 0 && len(b) == 0 {
		return 0, 0
	}
	if len(a) == 0 {
		return 0, total
	}
	if len(b) == 0 {
		return total, 0
	}

	aTokens := totalEntryTokens(a)
	bTokens := totalEntryTokens(b)
	totalContent := aTokens + bTokens

	if totalContent == 0 {
		return total / 2, total - total/2
	}

	// If everything fits, give each section what it needs
	if totalContent <= total {
		return aTokens, bTokens
	}

	// Minimum 30% each, proportional split of the rest
	minA := total * 30 / 100
	minB := total * 30 / 100
	flex := total - minA - minB

	aProportion := float64(aTokens) / float64(totalContent)
	aFlex := int(float64(flex) * aProportion)

	return minA + aFlex, total - (minA + aFlex)
}

// fillSection selects scored entries to fill a budget, with graceful degradation.
//
// Includes full entries by score order until ~80% of budget is consumed.
// Remaining entries get title-only summaries.
//
// Parameters:
//   - entries: Scored entries sorted by score descending
//   - budget: Token budget for this section
//
// Returns:
//   - []string: Full entry bodies that fit in the budget
//   - []string: Title-only summaries for entries that didn't fit
func fillSection(entries []ScoredEntry, budget int) ([]string, []string) {
	if len(entries) == 0 || budget <= 0 {
		return nil, nil
	}

	fullBudget := budget * 80 / 100
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

// fitItemsInBudget returns items that fit within a token budget.
//
// Items are included in order until the budget would be exceeded.
//
// Parameters:
//   - items: String items to include
//   - budget: Maximum token budget
//
// Returns:
//   - []string: Items that fit within the budget
func fitItemsInBudget(items []string, budget int) []string {
	if len(items) == 0 {
		return nil
	}
	used := 0
	var result []string
	for _, item := range items {
		tokens := context.EstimateTokensString(item)
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

// estimateSliceTokens sums token estimates for a string slice.
//
// Parameters:
//   - items: Strings to estimate
//
// Returns:
//   - int: Total estimated tokens
func estimateSliceTokens(items []string) int {
	total := 0
	for _, item := range items {
		total += context.EstimateTokensString(item)
	}
	return total
}

// totalEntryTokens sums pre-computed token counts for scored entries.
//
// Parameters:
//   - entries: Scored entries with token estimates
//
// Returns:
//   - int: Total tokens
func totalEntryTokens(entries []ScoredEntry) int {
	total := 0
	for _, e := range entries {
		total += e.Tokens
	}
	return total
}

// renderMarkdownPacket renders an assembled packet as Markdown.
//
// Parameters:
//   - pkt: Assembled packet to render
//
// Returns:
//   - string: Formatted Markdown output
func renderMarkdownPacket(pkt *assembledPacket) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString("# Context Packet" + nl)
	sb.WriteString(
		fmt.Sprintf(
			"Generated: %s | Budget: %d tokens | Used: ~%d",
			time.Now().UTC().Format(time.RFC3339), pkt.Budget, pkt.TokensUsed,
		) + nl + nl,
	)

	// Read order
	sb.WriteString("## Read These Files (in order)" + nl)
	for i, path := range pkt.ReadOrder {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, path) + nl)
	}
	sb.WriteString(nl)

	// Constitution
	if len(pkt.Constitution) > 0 {
		sb.WriteString("## Constitution (NEVER VIOLATE)" + nl)
		for _, rule := range pkt.Constitution {
			sb.WriteString(fmt.Sprintf("- %s", rule) + nl)
		}
		sb.WriteString(nl)
	}

	// Tasks
	if len(pkt.Tasks) > 0 {
		sb.WriteString("## Current Tasks" + nl)
		for _, t := range pkt.Tasks {
			sb.WriteString(t + nl)
		}
		sb.WriteString(nl)
	}

	// Conventions
	if len(pkt.Conventions) > 0 {
		sb.WriteString("## Key Conventions" + nl)
		for _, conv := range pkt.Conventions {
			sb.WriteString(fmt.Sprintf("- %s", conv) + nl)
		}
		sb.WriteString(nl)
	}

	// Decisions (full body)
	if len(pkt.Decisions) > 0 {
		sb.WriteString("## Recent Decisions" + nl)
		for _, dec := range pkt.Decisions {
			sb.WriteString(dec + nl + nl)
		}
	}

	// Learnings (full body)
	if len(pkt.Learnings) > 0 {
		sb.WriteString("## Key Learnings" + nl)
		for _, learn := range pkt.Learnings {
			sb.WriteString(learn + nl + nl)
		}
	}

	// Summaries
	if len(pkt.Summaries) > 0 {
		sb.WriteString("## Also Noted" + nl)
		for _, s := range pkt.Summaries {
			sb.WriteString(fmt.Sprintf("- %s", s) + nl)
		}
		sb.WriteString(nl)
	}

	sb.WriteString(pkt.Instruction + nl)

	return sb.String()
}
