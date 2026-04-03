//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

import (
	"encoding/json"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/entity"
	writeAgent "github.com/ActiveMemory/ctx/internal/write/agent"
)

// OutputAgentJSON writes the context packet as pretty-printed JSON.
//
// Uses budget-aware assembly to score entries and respect the token budget.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing the files
//   - budget: Token budget for content selection
//   - steeringBodies: Pre-filtered steering file bodies
//   - skillBody: Skill content (empty if none)
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func OutputAgentJSON(
	cmd *cobra.Command,
	ctx *entity.Context,
	budget int,
	steeringBodies []string,
	skillBody string,
) error {
	pkt := AssemblePacket(ctx, budget, steeringBodies, skillBody)

	packet := packet{
		Generated:    time.Now().UTC().Format(time.RFC3339),
		Budget:       pkt.Budget,
		TokensUsed:   pkt.TokensUsed,
		ReadOrder:    pkt.ReadOrder,
		Constitution: pkt.Constitution,
		Tasks:        pkt.Tasks,
		Conventions:  pkt.Conventions,
		Decisions:    pkt.Decisions,
		Learnings:    pkt.Learnings,
		Summaries:    pkt.Summaries,
		Steering:     pkt.Steering,
		Skill:        pkt.Skill,
		Instruction:  pkt.Instruction,
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(packet)
}

// OutputAgentMarkdown writes the context packet as formatted Markdown.
//
// Uses budget-aware assembly to score entries and respect the token budget.
// Output includes sections for constitution, tasks, conventions,
// decisions (full body), learnings (full body), title-only summaries,
// steering files, and skill content.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing the files
//   - budget: Token budget for content selection
//   - steeringBodies: Pre-filtered steering file bodies
//   - skillBody: Skill content (empty if none)
//
// Returns:
//   - error: Always nil (included for interface consistency)
func OutputAgentMarkdown(
	cmd *cobra.Command,
	ctx *entity.Context,
	budget int,
	steeringBodies []string,
	skillBody string,
) error {
	pkt := AssemblePacket(ctx, budget, steeringBodies, skillBody)
	writeAgent.Packet(cmd, RenderMarkdownPacket(pkt))
	return nil
}
