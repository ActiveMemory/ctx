//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/agent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/fmt"
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

	short, long := desc.Command(cmd.DescKeyAgent)

	c := &cobra.Command{
		Use:   cmd.UseAgent,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed(cFlag.Budget) {
				budget = rc.TokenBudget()
			}
			return Run(cmd, budget, format, cooldown, session)
		},
	}

	c.Flags().IntVar(
		&budget,
		cFlag.Budget, rc.DefaultTokenBudget,
		desc.Flag(flag.DescKeyAgentBudget),
	)
	c.Flags().StringVar(
		&format, cFlag.Format, fmt.FormatMarkdown,
		desc.Flag(flag.DescKeyAgentFormat),
	)
	c.Flags().DurationVar(
		&cooldown, cFlag.Cooldown, agent.DefaultCooldown,
		desc.Flag(flag.DescKeyAgentCooldown),
	)
	c.Flags().StringVar(
		&session, cFlag.Session, "",
		desc.Flag(flag.DescKeyAgentSession),
	)

	return c
}
