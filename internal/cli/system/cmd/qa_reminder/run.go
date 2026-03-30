//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package qa_reminder

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
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
	if !state.Initialized() {
		return nil
	}
	input, _, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}
	if !strings.Contains(input.ToolInput.Command, cfgGit.Binary) {
		return nil
	}
	fallback := desc.Text(text.DescKeyQaReminderFallback)
	msg := message.Load(
		hook.QAReminder, hook.VariantGate, nil, fallback,
	)
	if msg == "" {
		return nil
	}
	msg = ctxContext.AppendDir(msg)

	writeHook.Context(cmd, coreSession.FormatContext(hook.EventPreToolUse, msg))

	ref := notify.NewTemplateRef(hook.QAReminder, hook.VariantGate, nil)
	nudge.Relay(fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.QAReminder, desc.Text(text.DescKeyQaReminderRelayMessage)),
		input.SessionID, ref,
	)
	return nil
}
