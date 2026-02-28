//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// specsNudgeCmd returns the "ctx system specs-nudge" command.
//
// Prints a short directive reminding the agent to save plans to specs/
// for release tracking. Fires on EnterPlanMode via PreToolUse hook.
//
// Returns:
//   - *cobra.Command: Hidden subcommand for the specs nudge hook
func specsNudgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "specs-nudge",
		Short: "Plan-to-specs directory nudge",
		Long: `Emits a directive reminding the agent to save plans to specs/
for release tracking. Fires on EnterPlanMode tool use.

Hook event: PreToolUse (EnterPlanMode)
Output: agent directive (always, when .context/ is initialized)
Silent when: .context/ not initialized`,
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
			fallback := "Save your plan to specs/ — these documents track what was designed" +
				" for the current release. Use specs/feature-name.md naming. If this" +
				" is a quick fix that doesn't need a spec, proceed without one."
			msg := loadMessage("specs-nudge", "nudge", nil, fallback)
			if msg == "" {
				return nil
			}
			if line := contextDirLine(); line != "" {
				msg += " [" + line + "]"
			}
			printHookContext(cmd, "PreToolUse", msg)
			ref := notify.NewTemplateRef("specs-nudge", "nudge", nil)
			_ = notify.Send("relay", "specs-nudge: plan-to-specs nudge emitted", input.SessionID, ref)
			eventlog.Append("relay", "specs-nudge: plan-to-specs nudge emitted", input.SessionID, ref)
			return nil
		},
	}
}
