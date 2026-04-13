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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	remindStore "github.com/ActiveMemory/ctx/internal/cli/remind/core/store"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreProv "github.com/ActiveMemory/ctx/internal/cli/system/core/provenance"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/reminder"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Run executes the check-reminders hook logic.
//
// Reads hook input from stdin, loads pending reminders, filters to those
// that are due today or earlier, then emits a relay box with provenance
// (session, branch, commit) and the reminder list. Provenance is always
// emitted even when no reminders are due. Non-fatal on all errors.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input, _, paused := coreCheck.Preamble(stdin)

	// Provenance is unconditional — always print first,
	// regardless of initialized/paused state.
	coreProv.Emit(cmd, input.SessionID)

	if !state.Initialized() || paused {
		return nil
	}

	reminders, readErr := remindStore.Read()
	if readErr != nil || len(reminders) == 0 {
		return nil
	}

	today := time.Now().Format(cfgTime.DateFormat)
	var due []remindStore.Reminder
	for _, r := range reminders {
		if r.After == nil || *r.After <= today {
			due = append(due, r)
		}
	}

	if len(due) == 0 {
		return nil
	}

	var reminderList string
	for _, r := range due {
		reminderList += fmt.Sprintf(
			desc.Text(text.DescKeyCheckReminderItemFormat)+
				token.NewlineLF,
			r.ID, r.Message,
		)
	}

	fallback := reminderList +
		token.NewlineLF +
		desc.Text(text.DescKeyCheckReminderDismissHint) +
		token.NewlineLF +
		desc.Text(text.DescKeyCheckReminderDismissAllHint)
	vars := map[string]any{reminder.VarList: reminderList}
	relayMsg := fmt.Sprintf(
		desc.Text(text.DescKeyCheckReminderNudgeFormat),
		len(due),
	)
	nudge.LoadAndEmit(cmd,
		hook.CheckReminder, hook.VariantReminders,
		vars, fallback,
		desc.Text(text.DescKeyCheckReminderRelayPrefix),
		desc.Text(text.DescKeyCheckReminderBoxTitle),
		relayMsg, input.SessionID, "",
	)

	return nil
}
