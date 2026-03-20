//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_task_completion

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/nudge"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the check-task-completion hook logic.
//
// Tracks a per-session prompt counter and emits a task completion nudge
// when the counter reaches the configured interval. The counter resets
// after each nudge. Disabled when the nudge interval is zero or negative.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !core.Initialized() {
		return nil
	}
	input, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	interval := rc.TaskNudgeInterval()
	if interval <= 0 {
		return nil
	}

	counterPath := filepath.Join(core.StateDir(), nudge.PrefixTask+sessionID)
	count := core.ReadCounter(counterPath)
	count++

	if count < interval {
		core.WriteCounter(counterPath, count)
		return nil
	}

	// Threshold reached — reset and nudge.
	core.WriteCounter(counterPath, 0)

	fallback := desc.TextDesc(text.DescKeyCheckTaskCompletionFallback)
	msg := core.LoadMessage(
		hook.CheckTaskCompletion, hook.VariantNudge, nil, fallback,
	)
	if msg == "" {
		return nil
	}
	core.PrintHookContext(cmd, hook.EventPostToolUse, msg)

	nudgeMsg := desc.TextDesc(text.DescKeyCheckTaskCompletionNudgeMessage)
	ref := notify.NewTemplateRef(
		hook.CheckTaskCompletion, hook.VariantNudge, nil,
	)
	core.Relay(
		hook.CheckTaskCompletion+": "+nudgeMsg, input.SessionID, ref,
	)

	return nil
}
