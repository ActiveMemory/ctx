//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_task_completion

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/counter"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	coreNudge "github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/nudge"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
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
	if !state.Initialized() {
		return nil
	}
	input, sessionID, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	interval := rc.TaskNudgeInterval()
	if interval <= 0 {
		return nil
	}

	counterPath := filepath.Join(state.StateDir(), nudge.PrefixTask+sessionID)
	count := counter.Read(counterPath)
	count++

	if count < interval {
		counter.Write(counterPath, count)
		return nil
	}

	// Threshold reached - reset and nudge.
	counter.Write(counterPath, 0)

	fallback := desc.Text(text.DescKeyCheckTaskCompletionFallback)
	msg := message.LoadMessage(
		hook.CheckTaskCompletion, hook.VariantNudge, nil, fallback,
	)
	if msg == "" {
		return nil
	}
	writeHook.HookContext(
		cmd, coreSession.FormatContext(hook.EventPostToolUse, msg),
	)

	nudgeMsg := desc.Text(text.DescKeyCheckTaskCompletionNudgeMessage)
	ref := notify.NewTemplateRef(
		hook.CheckTaskCompletion, hook.VariantNudge, nil,
	)
	coreNudge.Relay(
		fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
			hook.CheckTaskCompletion, nudgeMsg), input.SessionID, ref,
	)

	return nil
}
