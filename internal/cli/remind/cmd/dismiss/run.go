//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core"
	errReminder "github.com/ActiveMemory/ctx/internal/err/reminder"
	"github.com/ActiveMemory/ctx/internal/write/remind"
)

// RunDismiss removes a single reminder by ID and prints confirmation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - idStr: String representation of the reminder ID
//
// Returns:
//   - error: Non-nil on invalid ID, missing reminder, or write failure
func RunDismiss(cmd *cobra.Command, idStr string) error {
	id, parseErr := strconv.Atoi(idStr)
	if parseErr != nil {
		return errReminder.InvalidID(idStr)
	}

	reminders, readErr := core.ReadReminders()
	if readErr != nil {
		return readErr
	}

	found := -1
	for i, r := range reminders {
		if r.ID == id {
			found = i
			break
		}
	}

	if found < 0 {
		return errReminder.NotFound(id)
	}

	remind.ReminderDismissed(cmd, reminders[found].ID, reminders[found].Message)
	reminders = append(reminders[:found], reminders[found+1:]...)
	return core.WriteReminders(reminders)
}

// RunDismissAll removes all reminders and prints confirmation.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read or write failure
func RunDismissAll(cmd *cobra.Command) error {
	reminders, readErr := core.ReadReminders()
	if readErr != nil {
		return readErr
	}

	if len(reminders) == 0 {
		remind.ReminderNone(cmd)
		return nil
	}

	for _, r := range reminders {
		remind.ReminderDismissed(cmd, r.ID, r.Message)
	}
	remind.ReminderDismissedAll(cmd, len(reminders))

	return core.WriteReminders([]core.Reminder{})
}
