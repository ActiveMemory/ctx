//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package qa_reminder

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Run executes the qa-reminder hook logic.
//
// Fires before any git command to inject a hard gate reminding the agent
// to lint, test, and verify a clean working tree before committing.
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
	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}
	if !strings.Contains(input.ToolInput.Command, "git") {
		return nil
	}
	fallback := assets.TextDesc(assets.TextDescKeyQaReminderFallback)
	msg := core.LoadMessage(
		config.HookQAReminder, config.VariantGate, nil, fallback,
	)
	if msg == "" {
		return nil
	}
	msg = core.AppendContextDir(msg)

	core.PrintHookContext(cmd, config.HookEventPreToolUse, msg)

	ref := notify.NewTemplateRef(config.HookQAReminder, config.VariantGate, nil)
	core.Relay(config.HookQAReminder+": "+
		assets.TextDesc(assets.TextDescKeyQaReminderRelayMessage),
		input.SessionID, ref,
	)
	return nil
}
