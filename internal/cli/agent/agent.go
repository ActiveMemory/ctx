//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"time"

	"github.com/spf13/cobra"

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

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Print AI-ready context packet",
		Long: `Print a concise context packet optimized for AI consumption.

The output is designed to be copy-pasted into an AI chat
or piped to a system prompt. It includes:
  - Constitution rules (NEVER VIOLATE)
  - Current tasks (budget-capped)
  - Key conventions (budget-capped)
  - Recent decisions (scored by relevance, full body)
  - Key learnings (scored by relevance, full body)

The --budget flag controls content selection. Entries are scored by
recency and relevance to active tasks, then included in priority order
until the budget is consumed. Entries that don't fit get title-only
summaries in an "Also Noted" section.

Use --budget to set token budget (default from .ctxrc or 8000).
Use --format to choose between Markdown (md) or JSON output.

Cooldown (for hooks and automation):
  --session identifies the caller (e.g., $PPID). Without it, cooldown
  is disabled and every call produces output. When --session is set,
  repeated calls within the --cooldown window (default 10m) are suppressed.

Examples:
  ctx agent                              # Default budget, Markdown output
  ctx agent --budget 4000                # Smaller context packet
  ctx agent --format json                # JSON output for programmatic use
  ctx agent --session $PPID              # Cooldown scoped to calling process`,
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
		&format, "format", config.FormatMarkdown, "Output format: md or json",
	)
	cmd.Flags().DurationVar(
		&cooldown, "cooldown", defaultCooldown,
		"Suppress repeated output within this duration (0 to disable)",
	)
	cmd.Flags().StringVar(
		&session, "session", "",
		"Session identifier for cooldown isolation (e.g., $PPID)",
	)

	return cmd
}
