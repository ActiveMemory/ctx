//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkjournal

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system check-journal" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-journal subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-journal",
		Short: "Journal export/enrich reminder hook",
		Long: `Detects unexported Claude Code sessions and unenriched journal entries,
then prints actionable commands. Throttled to once per day.

Hook event: UserPromptSubmit
Output: VERBATIM relay with export/enrich commands, silent otherwise
Silent when: no unexported sessions and no unenriched entries`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckJournal(cmd, os.Stdin)
		},
	}
}
