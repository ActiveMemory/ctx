//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/reminder"
	errReminder "github.com/ActiveMemory/ctx/internal/err/reminder"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ReadReminders reads all reminders from the JSON file.
//
// Returns:
//   - []Reminder: The parsed reminders (nil when file absent)
//   - error: Non-nil on read or parse failure
func ReadReminders() ([]Reminder, error) {
	data, readErr := io.SafeReadUserFile(RemindersPath())
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errReminder.Read(readErr)
	}
	var reminders []Reminder
	if parseErr := json.Unmarshal(data, &reminders); parseErr != nil {
		return nil, errReminder.Parse(parseErr)
	}
	return reminders, nil
}

// WriteReminders writes all reminders to the JSON file.
//
// Parameters:
//   - reminders: The reminder slice to persist
//
// Returns:
//   - error: Non-nil on marshal or write failure
func WriteReminders(reminders []Reminder) error {
	data, marshalErr := json.MarshalIndent(reminders, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	return io.SafeWriteFile(RemindersPath(), data, fs.PermFile)
}

// NextID returns the next available reminder ID (max existing + 1).
//
// Parameters:
//   - reminders: Existing reminders to scan
//
// Returns:
//   - int: The next sequential ID
func NextID(reminders []Reminder) int {
	mx := 0
	for _, r := range reminders {
		if r.ID > mx {
			mx = r.ID
		}
	}
	return mx + 1
}

// RemindersPath returns the full path to the reminders JSON file.
//
// Returns:
//   - string: Absolute path to reminders.json
func RemindersPath() string {
	return filepath.Join(rc.ContextDir(), reminder.Reminders)
}
