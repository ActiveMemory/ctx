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
	coreSteering "github.com/ActiveMemory/ctx/internal/cli/agent/core/steering"
	"github.com/ActiveMemory/ctx/internal/config/agent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/fmt"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the "ctx agent" command for generating AI-ready context packets.
//
// The command reads context files from .context/ and outputs a concise packet
// optimized for AI consumption, including constitution rules, active tasks,
// conventions, recent decisions, steering files, and optional skill content.
//
// Flags:
//   - --budget: Token budget for the context packet (default 8000)
//   - --format: Output format, "md" for Markdown or "json" (default "md")
//   - --cooldown: Suppress repeated output within this duration (default 10m)
//   - --session: Session identifier for cooldown tombstone isolation
//   - --skill: Include named skill content in context packet
//
// Returns:
//   - *cobra.Command: Configured agent command with flags registered
func Cmd() *cobra.Command {
	var (
		budget    int
		format    string
		cooldown  time.Duration
		session   string
		skillName string
	)

	short, long := desc.Command(cmd.DescKeyAgent)

	c := &cobra.Command{
		Use:     cmd.UseAgent,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyAgent),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed(cFlag.Budget) {
				budget = rc.TokenBudget()
			}

			// Tier 6: Load applicable steering files.
			steeringBodies := coreSteering.LoadBodies()

			// Tier 7: Load skill content if --skill is provided.
			var skillBody string
			if skillName != "" {
				sk, loadErr := coreSteering.LoadSkill(skillName)
				if loadErr != nil {
					return loadErr
				}
				skillBody = sk
			}

			return Run(
				cmd, budget, format, cooldown, session,
				steeringBodies, skillBody,
			)
		},
	}

	flagbind.IntFlag(
		c, &budget,
		cFlag.Budget, rc.DefaultTokenBudget,
		flag.DescKeyAgentBudget,
	)
	flagbind.StringFlagDefault(
		c, &format,
		cFlag.Format, fmt.FormatMarkdown,
		flag.DescKeyAgentFormat,
	)
	flagbind.DurationFlag(
		c, &cooldown,
		cFlag.Cooldown, agent.DefaultCooldown,
		flag.DescKeyAgentCooldown,
	)
	flagbind.StringFlag(
		c, &session,
		cFlag.Session, flag.DescKeyAgentSession,
	)
	flagbind.StringFlag(
		c, &skillName,
		cFlag.Skill, flag.DescKeyAgentSkill,
	)

	return c
}
