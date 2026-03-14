//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"
)

// ReminderAdded prints the confirmation for a newly added reminder.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: reminder ID.
//   - message: reminder text.
//   - after: optional date gate (nil if none).
func ReminderAdded(cmd *cobra.Command, id int, message string, after *string) {
	if cmd == nil {
		return
	}
	suffix := ""
	if after != nil {
		suffix = fmt.Sprintf(config.TplReminderAfterSuffix, *after)
	}
	cmd.Println(fmt.Sprintf(config.TplReminderAdded, id, message, suffix))
}

// ReminderItem prints a single reminder in the list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: reminder ID.
//   - message: reminder text.
//   - after: optional date gate (nil if none).
//   - today: current date in YYYY-MM-DD format.
func ReminderItem(cmd *cobra.Command, id int, message string, after *string, today string) {
	if cmd == nil {
		return
	}
	annotation := ""
	if after != nil && *after > today {
		annotation = fmt.Sprintf(config.TplReminderNotDue, *after)
	}
	cmd.Println(fmt.Sprintf(config.TplReminderItem, id, message, annotation))
}

// ReminderDismissed prints the confirmation for a dismissed reminder.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: reminder ID.
//   - message: reminder text.
func ReminderDismissed(cmd *cobra.Command, id int, message string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(config.TplReminderDismissed, id, message))
}

// ReminderNone prints the message when there are no reminders.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func ReminderNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplReminderNone)
}

// ReminderDismissedAll prints the summary after dismissing all reminders.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of dismissed reminders.
func ReminderDismissedAll(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(config.TplReminderDismissedAll, count))
}
