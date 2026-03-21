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
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
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

	short, long := desc.CommandDesc(cmd.DescKeyAgent)

	c := &cobra.Command{
		Use:   cmd.UseAgent,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed(cflag.Budget) {
				budget = rc.TokenBudget()
			}
			return Run(cmd, budget, format, cooldown, session)
		},
	}

	c.Flags().IntVar(
		&budget,
		cflag.Budget, rc.DefaultTokenBudget,
		desc.FlagDesc(flag.DescKeyAgentBudget),
	)
	c.Flags().StringVar(
		&format, cflag.Format, fmt.FormatMarkdown,
		desc.FlagDesc(flag.DescKeyAgentFormat),
	)
	c.Flags().DurationVar(
		&cooldown, cflag.Cooldown, agent.DefaultCooldown,
		desc.FlagDesc(flag.DescKeyAgentCooldown),
	)
	c.Flags().StringVar(
		&session, cflag.Session, "",
		desc.FlagDesc(flag.DescKeyAgentSession),
	)

	return c
}
