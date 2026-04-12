//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reminder

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Read wraps a failure to read the reminders file.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read reminders: <cause>"
func Read(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrReminderReadReminders), cause,
	)
}

// Parse wraps a failure to parse the reminders file.
//
// Parameters:
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "parse reminders: <cause>"
func Parse(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrReminderParseReminders), cause,
	)
}

// NotFound returns an error when no reminder matches the given ID.
//
// Parameters:
//   - id: the ID that was not found.
//
// Returns:
//   - error: "no reminder with ID <id>"
func NotFound(id int) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrReminderNotFound), id,
	)
}

// IDRequired returns an error when no reminder ID is provided.
//
// Returns:
//   - error: "provide a reminder ID or use --all"
func IDRequired() error {
	return errors.New(desc.Text(text.DescKeyErrReminderIDRequired))
}
