//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"github.com/spf13/cobra"
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
//
// Returns:
//   - *cobra.Command: Configured agent command with flags registered
func Cmd() *cobra.Command {
	var (
		budget int
		format string
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

Use --budget to limit token output (default 8000).
Use --format to choose between markdown (md) or JSON output.

Examples:
  ctx agent                    # Default 8000 token budget, markdown output
  ctx agent --budget 4000      # Smaller context packet for limited contexts
  ctx agent --format json      # JSON output for programmatic use
  ctx agent --budget 2000 --format json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgent(cmd, budget, format)
		},
	}

	cmd.Flags().IntVar(&budget, "budget", 8000, "Token budget for context packet")
	cmd.Flags().StringVar(&format, "format", "md", "Output format: md or json")

	return cmd
}
