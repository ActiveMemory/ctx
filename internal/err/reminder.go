//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// ReadReminders wraps a failure to read the reminders file.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read reminders: <cause>"
func ReadReminders(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrReminderReadReminders), cause,
	)
}

// ParseReminders wraps a failure to parse the reminders file.
//
// Parameters:
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "parse reminders: <cause>"
func ParseReminders(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrReminderParseReminders), cause,
	)
}

// InvalidID returns an error for an unparseable ID string.
//
// Parameters:
//   - value: the invalid ID string.
//
// Returns:
//   - error: "invalid ID <value>"
func InvalidID(value string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrReminderInvalidID), value,
	)
}

// ReminderNotFound returns an error when no reminder matches the given ID.
//
// Parameters:
//   - id: the ID that was not found.
//
// Returns:
//   - error: "no reminder with ID <id>"
func ReminderNotFound(id int) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrReminderNotFound), id,
	)
}

// ReminderIDRequired returns an error when no reminder ID is provided.
//
// Returns:
//   - error: "provide a reminder ID or use --all"
func ReminderIDRequired() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrReminderIDRequired))
}
