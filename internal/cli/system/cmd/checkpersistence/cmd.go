//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkpersistence

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system check-persistence" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-persistence subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-persistence",
		Short: "Persistence nudge hook",
		Long: `Tracks prompts since the last .context/ file modification and nudges
the agent to persist learnings, decisions, or task updates.

  Prompts  1-10: silent (too early)
  Prompts 11-25: nudge once at prompt 20 since last modification
  Prompts   25+: every 15th prompt since last modification

Hook event: UserPromptSubmit
Output: agent directive (when triggered), silent otherwise
Silent when: context files were recently modified`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckPersistence(cmd, os.Stdin)
		},
	}
}
