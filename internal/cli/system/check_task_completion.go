//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// checkTaskCompletionCmd returns the "ctx system check-task-completion" command.
//
// Fires after Edit/Write tool use (PostToolUse). Counts calls and
// periodically nudges the agent to mark completed tasks in TASKS.md.
// The interval is configurable via task_nudge_interval in .ctxrc
// (default 5, 0 = disabled).
func checkTaskCompletionCmd() *cobra.Command {
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
	if !isInitialized() {
		return nil
	}
	input := readInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = sessionUnknown
	}
	if paused(sessionID) > 0 {
		return nil
	}

	interval := rc.TaskNudgeInterval()
	if interval <= 0 {
		return nil
	}

	counterPath := filepath.Join(stateDir(), "task-nudge-"+sessionID)
	count := readCounter(counterPath)
	count++

	if count < interval {
		writeCounter(counterPath, count)
		return nil
	}

	// Threshold reached — reset and nudge.
	writeCounter(counterPath, 0)

	fallback := "If you completed a task, mark it [x] in TASKS.md."
	msg := loadMessage("check-task-completion", "nudge", nil, fallback)
	if msg == "" {
		return nil
	}
	printHookContext(cmd, "PostToolUse", msg)

	ref := notify.NewTemplateRef("check-task-completion", "nudge", nil)
	_ = notify.Send("relay", "check-task-completion: task completion nudge", input.SessionID, ref)
	eventlog.Append("relay", "check-task-completion: task completion nudge", input.SessionID, ref)

	return nil
}
