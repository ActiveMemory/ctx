//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/reminder"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errReminder "github.com/ActiveMemory/ctx/internal/err/reminder"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Read reads all reminders from the JSON file.
//
// Returns:
//   - []Reminder: The parsed reminders (nil when file absent)
//   - error: Non-nil on read or parse failure
func Read() ([]Reminder, error) {
	path, pathErr := Path()
	if pathErr != nil {
		return nil, pathErr
	}
	data, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errReminder.Read(readErr)
	}
	var reminders []Reminder
	if parseErr := json.Unmarshal(
		data, &reminders,
	); parseErr != nil {
		return nil, errReminder.Parse(parseErr)
	}
	return reminders, nil
}

// Write writes all reminders to the JSON file.
//
// Parameters:
//   - reminders: The reminder slice to persist
//
// Returns:
//   - error: Non-nil on marshal or write failure
func Write(reminders []Reminder) error {
	data, marshalErr := json.MarshalIndent(
		reminders, "", token.Indent2,
	)
	if marshalErr != nil {
		return marshalErr
	}
	path, pathErr := Path()
	if pathErr != nil {
		return pathErr
	}
	return io.SafeWriteFile(path, data, fs.PermFile)
}

// NextID returns the next available reminder ID
// (max existing + 1).
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

// Path returns the full path to the reminders JSON file.
//
// Returns:
//   - string: Absolute path to reminders.json
//   - error: non-nil when the context directory is not declared
func Path() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, reminder.File), nil
}
