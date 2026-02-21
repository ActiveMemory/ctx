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
		Use:   "qa-reminder",
		Short: "QA reminder hook",
		Long: `Emits a hard reminder to lint and test the entire project before
every commit. Fires on every Edit tool use — the repetition is
intentional reinforcement at the point of action.

Hook event: PreToolUse (Edit)
Output: agent directive (always, when .context/ is initialized)
Silent when: .context/ not initialized`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !isInitialized() {
				return nil
			}
			cmd.Println(
				"HARD GATE — DO NOT COMMIT without completing ALL of these steps first:" +
					" (1) lint the ENTIRE project," +
					" (2) test the ENTIRE project," +
					" (3) verify a clean working tree (no modified or untracked files left behind)." +
					" Not just the files you changed — the whole branch." +
					" If unrelated modified files remain," +
					" offer to commit them separately, stash them," +
					" or get explicit confirmation to leave them." +
					" Do NOT say 'I'll do that at the end' or 'I'll handle that after committing.'" +
					" Run lint and tests BEFORE every git commit, every time, no exceptions.",
			)
			return nil
		},
	}
}
