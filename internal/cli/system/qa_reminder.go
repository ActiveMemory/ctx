//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/spf13/cobra"
)

// qaReminderCmd returns the "ctx system qa-reminder" command.
//
// Prints a short reminder to lint and test the entire project before
// declaring work complete. Fires on every Edit via PreToolUse hook —
// the repetition is intentional reinforcement at the point of action.
//
// Returns:
//   - *cobra.Command: Hidden subcommand for the QA reminder hook
func qaReminderCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "qa-reminder",
		Short:  "QA reminder hook",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !isInitialized() {
				return nil
			}
			cmd.Println(
				"IMPORTANT: Before declaring code complete," +
					" lint and test the ENTIRE project —" +
					" not just the files you changed." +
					" You own the whole branch.",
			)
			return nil
		},
	}
}
