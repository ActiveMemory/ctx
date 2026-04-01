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

// Run dismisses one or all reminders based on the all flag.
//
// When all is true, removes every reminder. Otherwise removes the
// single reminder identified by idStr.
//
// Parameters:
//   - cmd: Cobra command for output
//   - idStr: String reminder ID (ignored when all is true)
//   - all: When true, dismiss all reminders
//
// Returns:
//   - error: Non-nil on invalid ID, missing reminder, or write failure
func Run(cmd *cobra.Command, idStr string, all bool) error {
	if all {
		return dismissAll(cmd)
	}
	return dismissOne(cmd, idStr)
}

// dismissOne removes a single reminder by its numeric ID.
//
// Parameters:
//   - cmd: Cobra command for status output
//   - idStr: String representation of the reminder ID
//
// Returns:
//   - error: Non-nil on invalid ID, missing reminder, or write failure
func dismissOne(cmd *cobra.Command, idStr string) error {
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

	remind.Dismissed(cmd, reminders[found].ID, reminders[found].Message)
	reminders = append(reminders[:found], reminders[found+1:]...)
	return core.WriteReminders(reminders)
}

// dismissAll removes every active reminder.
//
// Parameters:
//   - cmd: Cobra command for status output
//
// Returns:
//   - error: Non-nil on read or write failure
func dismissAll(cmd *cobra.Command) error {
	reminders, readErr := core.ReadReminders()
	if readErr != nil {
		return readErr
	}

	if len(reminders) == 0 {
		remind.None(cmd)
		return nil
	}

	for _, r := range reminders {
		remind.Dismissed(cmd, r.ID, r.Message)
	}
	remind.DismissedAll(cmd, len(reminders))

	return core.WriteReminders([]core.Reminder{})
}
