//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkcontextsize

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system check-context-size" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-context-size subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-context-size",
		Short: "Context size checkpoint hook",
		Long: `Counts prompts per session and emits VERBATIM relay reminders at
adaptive intervals, prompting the user to consider wrapping up.

  Prompts  1-15: silent
  Prompts 16-30: every 5th prompt
  Prompts   30+: every 3rd prompt

Also monitors actual context window token usage from session JSONL data.
Fires an independent warning when context window exceeds 80%, regardless
of prompt count.

Hook event: UserPromptSubmit
Output: VERBATIM relay (when triggered), silent otherwise
Silent when: early in session or between checkpoints`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckContextSize(cmd, os.Stdin)
		},
	}
}
