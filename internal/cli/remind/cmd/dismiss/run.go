//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core"
)

// runDismiss removes a single reminder by ID and prints confirmation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - idStr: String representation of the reminder ID
//
// Returns:
//   - error: Non-nil on invalid ID, missing reminder, or write failure
func runDismiss(cmd *cobra.Command, idStr string) error {
	id, parseErr := strconv.Atoi(idStr)
	if parseErr != nil {
		return fmt.Errorf("invalid ID %q", idStr)
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
		return fmt.Errorf("no reminder with ID %d", id)
	}

	cmd.Println(fmt.Sprintf("  - [%d] %s", reminders[found].ID, reminders[found].Message))
	reminders = append(reminders[:found], reminders[found+1:]...)
	return core.WriteReminders(reminders)
}

// runDismissAll removes all reminders and prints confirmation.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read or write failure
func runDismissAll(cmd *cobra.Command) error {
	reminders, readErr := core.ReadReminders()
	if readErr != nil {
		return readErr
	}

	if len(reminders) == 0 {
		cmd.Println("No reminders.")
		return nil
	}

	for _, r := range reminders {
		cmd.Println(fmt.Sprintf("  - [%d] %s", r.ID, r.Message))
	}
	cmd.Println(fmt.Sprintf("Dismissed %d reminders.", len(reminders)))

	return core.WriteReminders([]core.Reminder{})
}
