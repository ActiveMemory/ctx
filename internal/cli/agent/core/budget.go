//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/agent"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/index"
)

// AssembleBudgetPacket builds a context packet respecting the token budget.
//
// Allocation tiers:
//   - Tier 1 (always): constitution, read order, instruction
//   - Tier 2 (40%): active tasks
//   - Tier 3 (20%): conventions
//   - Tier 4+5 (remaining): decisions and learnings, scored by relevance
//
// Parameters:
//   - ctx: Loaded context containing the files
//   - budget: Token budget to respect
//
// Returns:
//   - *AssembledPacket: Assembled packet within budget
func AssembleBudgetPacket(ctx *entity.Context, budget int) *AssembledPacket {
	now := time.Now()
	pkt := &AssembledPacket{
		Budget:      budget,
		Instruction: desc.TextDesc(text.DescKeyAgentInstruction),
	}

	remaining := budget

	// Tier 1: Always included (constitution, read order, instruction)
	pkt.ReadOrder = ReadOrder(ctx)
	pkt.Constitution = ExtractConstitutionRules(ctx)

	tier1Tokens := EstimateSliceTokens(pkt.ReadOrder) +
		EstimateSliceTokens(pkt.Constitution) +
		ctxToken.EstimateTokensString(pkt.Instruction)
	remaining -= tier1Tokens

	if remaining <= 0 {
		pkt.TokensUsed = tier1Tokens
		return pkt
	}

	// Tier 2: Tasks (up to 40% of the original budget)
	taskCap := int(float64(budget) * agent.TaskBudgetPct)
	allTasks := ExtractActiveTasks(ctx)
	pkt.Tasks = FitItemsInBudget(allTasks, taskCap)
	taskTokens := EstimateSliceTokens(pkt.Tasks)
	remaining -= taskTokens

	if remaining <= 0 {
		pkt.TokensUsed = budget - remaining
		return pkt
	}

	// Tier 3: Conventions (up to 20% of the original budget)
	convCap := int(float64(budget) * agent.ConventionBudgetPct)
	allConventions := ExtractAllConventions(ctx)
	pkt.Conventions = FitItemsInBudget(allConventions, convCap)
	convTokens := EstimateSliceTokens(pkt.Conventions)
	remaining -= convTokens

	if remaining <= 0 {
		pkt.TokensUsed = budget - remaining
		return pkt
	}

	// Extract keywords from tasks for relevance scoring
	keywords := ExtractTaskKeywords(pkt.Tasks)

	// Tier 4+5: Decisions + Learnings (share remaining budget)
	decisionBlocks := ParseEntryBlocks(ctx, ctxCfg.Decision)
	learningBlocks := ParseEntryBlocks(ctx, ctxCfg.Learning)

	scoredDecisions := ScoreEntries(decisionBlocks, keywords, now)
	scoredLearnings := ScoreEntries(learningBlocks, keywords, now)

	// Split the remaining budget: proportional to content size, minimum 30% each
	decTokens, learnTokens := SplitBudget(
		remaining, scoredDecisions, scoredLearnings,
	)

	pkt.Decisions, pkt.Summaries = FillSection(scoredDecisions, decTokens)

	var learnSummaries []string
	pkt.Learnings, learnSummaries = FillSection(scoredLearnings, learnTokens)
	pkt.Summaries = append(pkt.Summaries, learnSummaries...)

	pkt.TokensUsed = tier1Tokens + taskTokens + convTokens +
		EstimateSliceTokens(pkt.Decisions) +
		EstimateSliceTokens(pkt.Learnings) +
		EstimateSliceTokens(pkt.Summaries)

	return pkt
}

// ExtractAllConventions extracts all bullet items from CONVENTIONS.md
// (not limited to 5 like the old implementation).
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: All convention bullet items; nil if the file is not found
func ExtractAllConventions(ctx *entity.Context) []string {
	if f := ctx.File(ctxCfg.Convention); f != nil {
		return ExtractBulletItems(string(f.Content), 1000)
	}
	return nil
}

// ParseEntryBlocks parses a context file into entry blocks.
//
// Parameters:
//   - ctx: Loaded context
//   - fileName: Name of the file to parse (e.g., config.Decision)
//
// Returns:
//   - []index.EntryBlock: Parsed entry blocks; nil if the file is not found
func ParseEntryBlocks(ctx *entity.Context, fileName string) []index.EntryBlock {
	if f := ctx.File(fileName); f != nil {
		return index.ParseEntryBlocks(string(f.Content))
	}
	return nil
}

// SplitBudget divides a token budget between two scored sections.
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
func SplitBudget(total int, a, b []ScoredEntry) (int, int) {
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
	minA := total * 30 / 100
	minB := total * 30 / 100
	flex := total - minA - minB

	aProportion := float64(aTokens) / float64(totalContent)
	aFlex := int(float64(flex) * aProportion)

	return minA + aFlex, total - (minA + aFlex)
}

// FillSection selects scored entries to fill a budget, with graceful degradation.
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
func FillSection(entries []ScoredEntry, budget int) ([]string, []string) {
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

// FitItemsInBudget returns items that fit within a token budget.
//
// Items are included in order until the budget would be exceeded.
//
// Parameters:
//   - items: String items to include
//   - budget: Maximum token budget
//
// Returns:
//   - []string: Items that fit within the budget
func FitItemsInBudget(items []string, budget int) []string {
	if len(items) == 0 {
		return nil
	}
	used := 0
	var result []string
	for _, item := range items {
		tokens := ctxToken.EstimateTokensString(item)
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
		total += ctxToken.EstimateTokensString(item)
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
func TotalEntryTokens(entries []ScoredEntry) int {
	total := 0
	for _, e := range entries {
		total += e.Tokens
	}
	return total
}

// RenderMarkdownPacket renders an assembled packet as Markdown.
//
// Parameters:
//   - pkt: Assembled packet to render
//
// Returns:
//   - string: Formatted Markdown output
func RenderMarkdownPacket(pkt *AssembledPacket) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(desc.TextDesc(text.DescKeyAgentPacketTitle) + nl)
	sb.WriteString(
		fmt.Sprintf(
			desc.TextDesc(text.DescKeyAgentPacketMeta),
			time.Now().UTC().Format(time.RFC3339), pkt.Budget, pkt.TokensUsed,
		) + nl + nl,
	)

	// Read order
	sb.WriteString(desc.TextDesc(text.DescKeyAgentSectionReadOrder) + nl)
	for i, path := range pkt.ReadOrder {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, path) + nl)
	}
	sb.WriteString(nl)

	// Constitution
	if len(pkt.Constitution) > 0 {
		sb.WriteString(desc.TextDesc(text.DescKeyAgentSectionConstitution) + nl)
		for _, rule := range pkt.Constitution {
			sb.WriteString(fmt.Sprintf("- %s", rule) + nl)
		}
		sb.WriteString(nl)
	}

	// Tasks
	if len(pkt.Tasks) > 0 {
		sb.WriteString(desc.TextDesc(text.DescKeyAgentSectionTasks) + nl)
		for _, t := range pkt.Tasks {
			sb.WriteString(t + nl)
		}
		sb.WriteString(nl)
	}

	// Conventions
	if len(pkt.Conventions) > 0 {
		sb.WriteString(desc.TextDesc(text.DescKeyAgentSectionConventions) + nl)
		for _, conv := range pkt.Conventions {
			sb.WriteString(fmt.Sprintf("- %s", conv) + nl)
		}
		sb.WriteString(nl)
	}

	// Decisions (full body)
	if len(pkt.Decisions) > 0 {
		sb.WriteString(desc.TextDesc(text.DescKeyAgentSectionDecisions) + nl)
		for _, dec := range pkt.Decisions {
			sb.WriteString(dec + nl + nl)
		}
	}

	// Learnings (full body)
	if len(pkt.Learnings) > 0 {
		sb.WriteString(desc.TextDesc(text.DescKeyAgentSectionLearnings) + nl)
		for _, learn := range pkt.Learnings {
			sb.WriteString(learn + nl + nl)
		}
	}

	// Summaries
	if len(pkt.Summaries) > 0 {
		sb.WriteString(desc.TextDesc(text.DescKeyAgentSectionSummaries) + nl)
		for _, s := range pkt.Summaries {
			sb.WriteString(fmt.Sprintf("- %s", s) + nl)
		}
		sb.WriteString(nl)
	}

	sb.WriteString(pkt.Instruction + nl)

	return sb.String()
}
