//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

import (
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/agent/core/extract"
	"github.com/ActiveMemory/ctx/internal/cli/agent/core/score"
	"github.com/ActiveMemory/ctx/internal/cli/agent/core/sort"
	"github.com/ActiveMemory/ctx/internal/config/agent"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// AssemblePacket builds a context packet respecting the token budget.
//
// Allocation tiers:
//   - Tier 1 (always): constitution, read order, instruction
//   - Tier 2 (40%): active tasks
//   - Tier 3 (20%): conventions
//   - Tier 4+5 (remaining): decisions and learnings, scored by relevance
//   - Tier 6 (remaining after 4+5): steering files
//   - Tier 7 (remaining after 6): skill content (--skill flag)
//
// Parameters:
//   - ctx: Loaded context containing the files
//   - budget: Token budget to respect
//   - steeringBodies: Pre-filtered steering file bodies to include
//   - skillBody: Skill content to include (empty string if none)
//
// Returns:
//   - *AssembledPacket: Assembled packet within budget
func AssemblePacket(
	ctx *entity.Context,
	budget int,
	steeringBodies []string,
	skillBody string,
	sharedBodies []string,
) *AssembledPacket {
	now := time.Now()
	pkt := &AssembledPacket{
		Budget:      budget,
		Instruction: desc.Text(text.DescKeyAgentInstruction),
	}

	remaining := budget

	// Tier 1: Always included (constitution, read order, instruction)
	pkt.ReadOrder = sort.ReadOrder(ctx)
	pkt.Constitution = extract.ConstitutionRules(ctx)

	tier1Tokens := EstimateSliceTokens(pkt.ReadOrder) +
		EstimateSliceTokens(pkt.Constitution) +
		ctxToken.EstimateString(pkt.Instruction)
	remaining -= tier1Tokens

	if remaining <= 0 {
		pkt.TokensUsed = tier1Tokens
		return pkt
	}

	// Tier 2: Tasks (up to 40% of the original budget)
	taskCap := int(float64(budget) * agent.TaskBudgetPct)
	allTasks := extract.ActiveTasks(ctx)
	pkt.Tasks = FitItems(allTasks, taskCap)
	taskTokens := EstimateSliceTokens(pkt.Tasks)
	remaining -= taskTokens

	if remaining <= 0 {
		pkt.TokensUsed = budget - remaining
		return pkt
	}

	// Tier 3: Conventions (up to 20% of the original budget)
	convCap := int(float64(budget) * agent.ConventionBudgetPct)
	allConventions := ExtractAllConventions(ctx)
	pkt.Conventions = FitItems(allConventions, convCap)
	convTokens := EstimateSliceTokens(pkt.Conventions)
	remaining -= convTokens

	if remaining <= 0 {
		pkt.TokensUsed = budget - remaining
		return pkt
	}

	// Extract keywords from tasks for relevance scoring
	keywords := score.ExtractTaskKeywords(pkt.Tasks)

	// Tier 4+5: Decisions + Learnings (share remaining budget)
	decisionBlocks := ParseEntryBlocks(ctx, cfgCtx.Decision)
	learningBlocks := ParseEntryBlocks(ctx, cfgCtx.Learning)

	scoredDecisions := score.All(decisionBlocks, keywords, now)
	scoredLearnings := score.All(learningBlocks, keywords, now)

	// Split the remaining budget: proportional to content size, minimum 30% each
	decTokens, learnTokens := Split(
		remaining, scoredDecisions, scoredLearnings,
	)

	pkt.Decisions, pkt.Summaries = FillSection(scoredDecisions, decTokens)

	var learnSummaries []string
	pkt.Learnings, learnSummaries = FillSection(scoredLearnings, learnTokens)
	pkt.Summaries = append(pkt.Summaries, learnSummaries...)

	usedSoFar := tier1Tokens + taskTokens + convTokens +
		EstimateSliceTokens(pkt.Decisions) +
		EstimateSliceTokens(pkt.Learnings) +
		EstimateSliceTokens(pkt.Summaries)
	remaining = budget - usedSoFar

	// Tier 6: Steering files (from remaining budget)
	if remaining > 0 && len(steeringBodies) > 0 {
		pkt.Steering = FitItems(steeringBodies, remaining)
		steeringTokens := EstimateSliceTokens(pkt.Steering)
		remaining -= steeringTokens
		usedSoFar += steeringTokens
	}

	// Tier 7: Skill content (from remaining budget)
	if remaining > 0 && skillBody != "" {
		skillTokens := ctxToken.EstimateString(skillBody)
		if skillTokens <= remaining {
			pkt.Skill = skillBody
			remaining -= skillTokens
			usedSoFar += skillTokens
		}
	}

	// Tier 8: Shared hub entries (from remaining budget)
	if remaining > 0 && len(sharedBodies) > 0 {
		pkt.Shared = FitItems(sharedBodies, remaining)
		sharedTokens := EstimateSliceTokens(pkt.Shared)
		usedSoFar += sharedTokens
	}

	pkt.TokensUsed = usedSoFar

	return pkt
}
