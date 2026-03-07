//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package qareminder

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Cmd returns the "ctx system qa-reminder" subcommand.
//
// Returns:
//   - *cobra.Command: Configured qa-reminder subcommand
func Cmd() *cobra.Command {
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
			if !core.IsInitialized() {
				return nil
			}
			input := core.ReadInput(os.Stdin)
			sessionID := input.SessionID
			if sessionID == "" {
				sessionID = core.SessionUnknown
			}
			if core.Paused(sessionID) > 0 {
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
			msg := core.LoadMessage("qa-reminder", "gate", nil, fallback)
			if msg == "" {
				return nil
			}
			if line := core.ContextDirLine(); line != "" {
				msg += " [" + line + "]"
			}
			core.PrintHookContext(cmd, "PreToolUse", msg)
			ref := notify.NewTemplateRef("qa-reminder", "gate", nil)
			_ = notify.Send("relay", "qa-reminder: QA gate reminder emitted", input.SessionID, ref)
			eventlog.Append("relay", "qa-reminder: QA gate reminder emitted", input.SessionID, ref)
			return nil
		},
	}
}
