//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_reminder

import (
	"fmt"
	"os"
	"time"

	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	remindCore "github.com/ActiveMemory/ctx/internal/cli/remind/core"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/reminder"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the check-reminders hook logic.
//
// Reads hook input from stdin, loads pending reminders, filters to those
// that are due today or earlier, then emits a relay box with the reminder
// list if any are due. Non-fatal on all errors.
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

	input, _, paused := hook2.Preamble(stdin)
	if paused {
		return nil
	}

	reminders, readErr := remindCore.ReadReminders()
	if readErr != nil {
		return nil // non-fatal: don't break session start
	}

	today := time.Now().Format(cfgTime.DateFormat)
	var due []remindCore.Reminder
	for _, r := range reminders {
		if r.After == nil || *r.After <= today {
			due = append(due, r)
		}
	}

	if len(due) == 0 {
		return nil
	}

	// Build a pre-formatted reminder list for the template variable
	var reminderList string
	for _, r := range due {
		reminderList += fmt.Sprintf(
			desc.Text(text.DescKeyCheckRemindersItemFormat)+token.NewlineLF,
			r.ID, r.Message,
		)
	}

	fallback := reminderList +
		token.NewlineLF +
		desc.Text(text.DescKeyCheckRemindersDismissHint) + token.NewlineLF +
		desc.Text(text.DescKeyCheckRemindersDismissAllHint)
	vars := map[string]any{reminder.VarReminderList: reminderList}
	content := core.LoadMessage(
		hook.CheckReminders, hook.VariantReminders, vars, fallback,
	)
	if content == "" {
		return nil
	}

	writeHook.Nudge(cmd, core.NudgeBox(
		desc.Text(text.DescKeyCheckRemindersRelayPrefix),
		desc.Text(text.DescKeyCheckRemindersBoxTitle),
		content))

	ref := notify.NewTemplateRef(hook.CheckReminders, hook.VariantReminders, vars)
	nudgeMsg := fmt.Sprintf(
		desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckReminders,
		fmt.Sprintf(desc.Text(text.DescKeyCheckRemindersNudgeFormat), len(due)),
	)
	core.NudgeAndRelay(nudgeMsg, input.SessionID, ref)

	return nil
}
