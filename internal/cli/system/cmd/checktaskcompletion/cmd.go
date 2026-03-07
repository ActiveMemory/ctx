//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checktaskcompletion

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the "ctx system check-task-completion" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-task-completion subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-task-completion",
		Short: "Task completion nudge after edits",
		Long: `Counts Edit/Write tool calls and periodically nudges the agent
to check whether any tasks should be marked done in TASKS.md.

Hook event: PostToolUse (Edit, Write)
Output: agent directive every N edits, silent otherwise
Silent when: counter below threshold, interval is 0, or session is paused`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckTaskCompletion(cmd, os.Stdin)
		},
	}
}

func runCheckTaskCompletion(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}
	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
		return nil
	}

	interval := rc.TaskNudgeInterval()
	if interval <= 0 {
		return nil
	}

	counterPath := filepath.Join(core.StateDir(), "task-nudge-"+sessionID)
	count := core.ReadCounter(counterPath)
	count++

	if count < interval {
		core.WriteCounter(counterPath, count)
		return nil
	}

	// Threshold reached — reset and nudge.
	core.WriteCounter(counterPath, 0)

	fallback := "If you completed a task, mark it [x] in TASKS.md."
	msg := core.LoadMessage("check-task-completion", "nudge", nil, fallback)
	if msg == "" {
		return nil
	}
	core.PrintHookContext(cmd, "PostToolUse", msg)

	ref := notify.NewTemplateRef("check-task-completion", "nudge", nil)
	_ = notify.Send("relay", "check-task-completion: task completion nudge", input.SessionID, ref)
	eventlog.Append("relay", "check-task-completion: task completion nudge", input.SessionID, ref)

	return nil
}
