//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_task_completion

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
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
	if !core.IsInitialized() {
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

	counterPath := filepath.Join(core.StateDir(), config.TaskNudgePrefix+sessionID)
	count := core.ReadCounter(counterPath)
	count++

	if count < interval {
		core.WriteCounter(counterPath, count)
		return nil
	}

	// Threshold reached — reset and nudge.
	core.WriteCounter(counterPath, 0)

	fallback := assets.TextDesc(assets.TextDescKeyCheckTaskCompletionFallback)
	msg := core.LoadMessage(
		config.HookCheckTaskCompletion, config.VariantNudge, nil, fallback,
	)
	if msg == "" {
		return nil
	}
	core.PrintHookContext(cmd, config.HookEventPostToolUse, msg)

	nudgeMsg := assets.TextDesc(assets.TextDescKeyCheckTaskCompletionNudgeMessage)
	ref := notify.NewTemplateRef(
		config.HookCheckTaskCompletion, config.VariantNudge, nil,
	)
	core.Relay(
		config.HookCheckTaskCompletion+": "+nudgeMsg, input.SessionID, ref,
	)

	return nil
}
