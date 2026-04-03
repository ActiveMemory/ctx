//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core/store"
	errReminder "github.com/ActiveMemory/ctx/internal/err/reminder"
	"github.com/ActiveMemory/ctx/internal/write/remind"
)

// One removes a single reminder by its numeric ID.
//
// Parameters:
//   - cmd: Cobra command for status output
//   - idStr: String representation of the reminder ID
//
// Returns:
//   - error: Non-nil on invalid ID, missing reminder, or
//     write failure
func One(
	cmd *cobra.Command, idStr string,
) error {
	id, parseErr := strconv.Atoi(idStr)
	if parseErr != nil {
		return errReminder.InvalidID(idStr)
	}

	reminders, readErr := store.Read()
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

	remind.Dismissed(
		cmd,
		reminders[found].ID,
		reminders[found].Message,
	)
	reminders = append(
		reminders[:found], reminders[found+1:]...,
	)
	return store.Write(reminders)
}

// All removes every active reminder.
//
// Parameters:
//   - cmd: Cobra command for status output
//
// Returns:
//   - error: Non-nil on read or write failure
func All(cmd *cobra.Command) error {
	reminders, readErr := store.Read()
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

	return store.Write([]store.Reminder{})
}
