//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package remind

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Added prints the confirmation for a newly added reminder.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: reminder ID.
//   - message: reminder text.
//   - after: optional date gate (nil if none).
func Added(cmd *cobra.Command, id int, message string, after *string) {
	if cmd == nil {
		return
	}
	suffix := ""
	if after != nil {
		suffix = fmt.Sprintf(desc.Text(text.DescKeyWriteReminderAfterSuffix), *after)
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteReminderAdded),
		id, message, suffix))
}

// Item prints a single reminder in the list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: reminder ID.
//   - message: reminder text.
//   - after: optional date gate (nil if none).
//   - today: current date in YYYY-MM-DD format.
func Item(
	cmd *cobra.Command,
	id int,
	message string,
	after *string,
	today string,
) {
	if cmd == nil {
		return
	}
	annotation := ""
	if after != nil && *after > today {
		annotation = fmt.Sprintf(desc.Text(text.DescKeyWriteReminderNotDue), *after)
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteReminderItem),
		id, message, annotation))
}

// Dismissed prints the confirmation for a dismissed reminder.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: reminder ID.
//   - message: reminder text.
func Dismissed(cmd *cobra.Command, id int, message string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteReminderDismissed),
		id, message))
}

// None prints the message when there are no reminders.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func None(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteReminderNone))
}

// DismissedAll prints the summary after dismissing all reminders.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of dismissed reminders.
func DismissedAll(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteReminderDismissedAll),
		count))
}
