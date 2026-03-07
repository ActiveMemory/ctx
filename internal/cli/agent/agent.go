//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	agentroot "github.com/ActiveMemory/ctx/internal/cli/agent/cmd/root"
	"github.com/ActiveMemory/ctx/internal/cli/agent/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the "ctx agent" command for generating AI-ready context packets.
//
// The command reads context files from .context/ and outputs a concise packet
// optimized for AI consumption, including constitution rules, active tasks,
// conventions, and recent decisions.
//
// Flags:
//   - --budget: Token budget for the context packet (default 8000)
//   - --format: Output format, "md" for Markdown or "json" (default "md")
//   - --cooldown: Suppress repeated output within this duration (default 10m)
//   - --session: Session identifier for cooldown tombstone isolation
//
// Returns:
//   - *cobra.Command: Configured agent command with flags registered
func Cmd() *cobra.Command {
	var (
		budget   int
		format   string
		cooldown time.Duration
		session  string
	)

	short, long := assets.CommandDesc("agent")

	cmd := &cobra.Command{
		Use:   "agent",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("budget") {
				budget = rc.TokenBudget()
			}
			return agentroot.Run(cmd, budget, format, cooldown, session)
		},
	}

	cmd.Flags().IntVar(
		&budget,
		"budget", rc.DefaultTokenBudget, "Token budget for context packet",
	)
	cmd.Flags().StringVar(
		&format, "format", config.FormatMarkdown, "Output format: md or json",
	)
	cmd.Flags().DurationVar(
		&cooldown, "cooldown", core.DefaultCooldown,
		"Suppress repeated output within this duration (0 to disable)",
	)
	cmd.Flags().StringVar(
		&session, "session", "",
		"Session identifier for cooldown isolation (e.g., $PPID)",
	)

	return cmd
}
