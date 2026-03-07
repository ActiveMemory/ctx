//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package specsnudge

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Cmd returns the "ctx system specs-nudge" subcommand.
//
// Returns:
//   - *cobra.Command: Configured specs-nudge subcommand
func Cmd() *cobra.Command {
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
			fallback := "Save your plan to specs/ — these documents track what was designed" +
				" for the current release. Use specs/feature-name.md naming. If this" +
				" is a quick fix that doesn't need a spec, proceed without one."
			msg := core.LoadMessage("specs-nudge", "nudge", nil, fallback)
			if msg == "" {
				return nil
			}
			if line := core.ContextDirLine(); line != "" {
				msg += " [" + line + "]"
			}
			core.PrintHookContext(cmd, "PreToolUse", msg)
			ref := notify.NewTemplateRef("specs-nudge", "nudge", nil)
			_ = notify.Send("relay", "specs-nudge: plan-to-specs nudge emitted", input.SessionID, ref)
			eventlog.Append("relay", "specs-nudge: plan-to-specs nudge emitted", input.SessionID, ref)
			return nil
		},
	}
}
