//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/entity"
)

// OutputAgentJSON writes the context packet as pretty-printed JSON.
//
// Uses budget-aware assembly to score entries and respect the token budget.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing the files
//   - budget: Token budget for content selection
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func OutputAgentJSON(
	cmd *cobra.Command, ctx *entity.Context, budget int,
) error {
	pkt := AssembleBudgetPacket(ctx, budget)

	packet := Packet{
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
// decisions (full body), learnings (full body), and title-only summaries.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing the files
//   - budget: Token budget for content selection
//
// Returns:
//   - error: Always nil (included for interface consistency)
func OutputAgentMarkdown(
	cmd *cobra.Command, ctx *entity.Context, budget int,
) error {
	pkt := AssembleBudgetPacket(ctx, budget)
	cmd.Print(RenderMarkdownPacket(pkt))
	return nil
}
