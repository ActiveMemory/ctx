//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_reminders

import (
	"fmt"
	"os"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/file"
	time2 "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	remindcore "github.com/ActiveMemory/ctx/internal/cli/remind/core"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
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
	if !core.IsInitialized() {
		return nil
	}

	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	reminders, readErr := remindcore.ReadReminders()
	if readErr != nil {
		return nil // non-fatal: don't break session start
	}

	today := time.Now().Format(time2.DateFormat)
	var due []remindcore.Reminder
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
		reminderList += fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckRemindersItemFormat)+config.NewlineLF, r.ID, r.Message)
	}

	fallback := reminderList +
		config.NewlineLF + assets.TextDesc(assets.TextDescKeyCheckRemindersDismissHint) + config.NewlineLF +
		assets.TextDesc(assets.TextDescKeyCheckRemindersDismissAllHint)
	vars := map[string]any{file.TplVarReminderList: reminderList}
	content := core.LoadMessage(file.HookCheckReminders, file.VariantReminders, vars, fallback)
	if content == "" {
		return nil
	}

	cmd.Println(core.NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckRemindersRelayPrefix),
		assets.TextDesc(assets.TextDescKeyCheckRemindersBoxTitle),
		content))

	ref := notify.NewTemplateRef(file.HookCheckReminders, file.VariantReminders, vars)
	nudgeMsg := file.HookCheckReminders + ": " + fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckRemindersNudgeFormat), len(due))
	core.NudgeAndRelay(nudgeMsg, input.SessionID, ref)

	return nil
}
