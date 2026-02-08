//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// outputAgentJSON writes the context packet as pretty-printed JSON.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing the files
//   - budget: Token budget to include in the packet metadata
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func outputAgentJSON(
	cmd *cobra.Command, ctx *context.Context, budget int,
) error {
	packet := Packet{
		Generated:    time.Now().UTC().Format(time.RFC3339),
		Budget:       budget,
		TokensUsed:   ctx.TotalTokens,
		ReadOrder:    getReadOrder(ctx),
		Constitution: extractConstitutionRules(ctx),
		Tasks:        extractActiveTasks(ctx),
		Conventions:  extractConventions(ctx),
		Decisions:    extractRecentDecisions(ctx, 3),
		Instruction: "Before starting work, confirm to the user: " +
			"\"I have read the required context files and " +
			"I'm following project conventions.\"",
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(packet)
}

// outputAgentMarkdown writes the context packet as formatted Markdown.
//
// The output includes sections for read order, constitution rules,
// current tasks, conventions, and recent decisions.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context containing the files
//   - budget: Token budget to display in the header
//
// Returns:
//   - error: Always nil (included for interface consistency)
func outputAgentMarkdown(
	cmd *cobra.Command, ctx *context.Context, budget int,
) error {
	var sb strings.Builder

	timestamp := time.Now().UTC().Format(time.RFC3339)
	sb.WriteString("# Context Packet" + config.NewlineLF)
	sb.WriteString(
		fmt.Sprintf(
			"Generated: %s | Budget: %d tokens | Used: %d",
			timestamp, budget, ctx.TotalTokens,
		) + config.NewlineLF + config.NewlineLF,
	)

	// Read order
	sb.WriteString("## Read These Files (in order)" + config.NewlineLF)
	for i, path := range getReadOrder(ctx) {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, path) + config.NewlineLF)
	}
	sb.WriteString(config.NewlineLF)

	// Constitution
	rules := extractConstitutionRules(ctx)
	if len(rules) > 0 {
		sb.WriteString("## Constitution (NEVER VIOLATE)" + config.NewlineLF)
		for _, rule := range rules {
			sb.WriteString(fmt.Sprintf("- %s", rule) + config.NewlineLF)
		}
		sb.WriteString(config.NewlineLF)
	}

	// Current tasks
	tasks := extractActiveTasks(ctx)
	if len(tasks) > 0 {
		sb.WriteString("## Current Tasks" + config.NewlineLF)
		for _, task := range tasks {
			sb.WriteString(task + config.NewlineLF)
		}
		sb.WriteString(config.NewlineLF)
	}

	// Conventions
	conventions := extractConventions(ctx)
	if len(conventions) > 0 {
		sb.WriteString("## Key Conventions" + config.NewlineLF)
		for _, conv := range conventions {
			sb.WriteString(fmt.Sprintf("- %s", conv) + config.NewlineLF)
		}
		sb.WriteString(config.NewlineLF)
	}

	// Recent decisions
	decisions := extractRecentDecisions(ctx, 3)
	if len(decisions) > 0 {
		sb.WriteString("## Recent Decisions" + config.NewlineLF)
		for _, dec := range decisions {
			sb.WriteString(fmt.Sprintf("- %s", dec) + config.NewlineLF)
		}
		sb.WriteString(config.NewlineLF)
	}

	sb.WriteString(
		"Before starting work, confirm to the user: " +
			"\"I have read the required context files and I'm " +
			"following project conventions.\"" + config.NewlineLF,
	)

	cmd.Print(sb.String())
	return nil
}
