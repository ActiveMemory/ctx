//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// qaReminderCmd returns the "ctx system qa-reminder" command.
//
// Prints a short reminder to lint and test the entire project before
// committing. Fires on PreToolUse:Bash when the command contains "git",
// placing the reminder at the point of action — the commit sequence —
// rather than during every edit.
//
// Returns:
//   - *cobra.Command: Hidden subcommand for the QA reminder hook
func qaReminderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "qa-reminder",
		Short: "QA reminder hook",
		Long: `Emits a hard reminder to lint and test the entire project before
committing. Fires on Bash tool use when the command contains "git",
placing reinforcement at the commit sequence rather than during edits.

Hook event: PreToolUse (Bash)
Output: agent directive (when command contains "git" and .context/ is initialized)
Silent when: .context/ not initialized or command does not contain "git"`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !isInitialized() {
				return nil
			}
			input := readInput(os.Stdin)
			sessionID := input.SessionID
			if sessionID == "" {
				sessionID = sessionUnknown
			}
			if paused(sessionID) > 0 {
				return nil
			}
			if !strings.Contains(input.ToolInput.Command, "git") {
				return nil
			}
			fallback := "HARD GATE — DO NOT COMMIT without completing ALL of these steps first:" +
				" (1) lint the ENTIRE project," +
				" (2) test the ENTIRE project," +
				" (3) verify a clean working tree (no modified or untracked files left behind)." +
				" Not just the files you changed — the whole branch." +
				" If unrelated modified files remain," +
				" offer to commit them separately, stash them," +
				" or get explicit confirmation to leave them." +
				" Do NOT say 'I'll do that at the end' or 'I'll handle that after committing.'" +
				" Run lint and tests BEFORE every git commit, every time, no exceptions."
			msg := loadMessage("qa-reminder", "gate", nil, fallback)
			if msg == "" {
				return nil
			}
			if line := contextDirLine(); line != "" {
				msg += " [" + line + "]"
			}
			printHookContext(cmd, "PreToolUse", msg)
			ref := notify.NewTemplateRef("qa-reminder", "gate", nil)
			_ = notify.Send("relay", "qa-reminder: QA gate reminder emitted", input.SessionID, ref)
			eventlog.Append("relay", "qa-reminder: QA gate reminder emitted", input.SessionID, ref)
			return nil
		},
	}
}
