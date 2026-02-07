//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"time"

	"github.com/spf13/cobra"

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

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Print AI-ready context packet",
		Long: `Print a concise context packet optimized for AI consumption.

The output is designed to be copy-pasted into an AI chat
or piped to a system prompt. It includes:
  - Constitution rules (NEVER VIOLATE)
  - Current tasks
  - Key conventions
  - Recent decisions

Use --budget to limit token output (default from .contextrc or 8000).
Use --format to choose between markdown (md) or JSON output.
Use --cooldown to suppress repeated output within a time window (default 10m).
Use --session to isolate the cooldown per session (e.g., $PPID).

Examples:
  ctx agent                              # Default budget, markdown output
  ctx agent --budget 4000                # Smaller context packet
  ctx agent --format json                # JSON output for programmatic use
  ctx agent --cooldown 15m --session 123 # With cooldown, session-scoped`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("budget") {
				budget = rc.TokenBudget()
			}
			return runAgent(cmd, budget, format, cooldown, session)
		},
	}

	cmd.Flags().IntVar(
		&budget,
		"budget", rc.DefaultTokenBudget, "Token budget for context packet",
	)
	cmd.Flags().StringVar(
		&format, "format", "md", "Output format: md or json",
	)
	cmd.Flags().DurationVar(
		&cooldown, "cooldown", DefaultCooldown,
		"Suppress repeated output within this duration (0 to disable)",
	)
	cmd.Flags().StringVar(
		&session, "session", "",
		"Session identifier for cooldown isolation (e.g., $PPID)",
	)

	return cmd
}
