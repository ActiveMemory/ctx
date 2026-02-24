//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package remind

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Reminder represents a single session-scoped reminder.
type Reminder struct {
	ID      int     `json:"id"`
	Message string  `json:"message"`
	Created string  `json:"created"`
	After   *string `json:"after"` // nullable YYYY-MM-DD
}

// ReadReminders reads all reminders from the JSON file.
// Returns (nil, nil) if the file does not exist.
func ReadReminders() ([]Reminder, error) {
	path := RemindersPath()
	data, err := os.ReadFile(path) //nolint:gosec // project-local path
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read reminders: %w", err)
	}
	var reminders []Reminder
	if err := json.Unmarshal(data, &reminders); err != nil {
		return nil, fmt.Errorf("parse reminders: %w", err)
	}
	return reminders, nil
}

// WriteReminders writes all reminders to the JSON file.
func WriteReminders(reminders []Reminder) error {
	data, err := json.MarshalIndent(reminders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(RemindersPath(), data, config.PermFile)
}

// NextID returns the next available reminder ID (max existing + 1).
func NextID(reminders []Reminder) int {
	max := 0
	for _, r := range reminders {
		if r.ID > max {
			max = r.ID
		}
	}
	return max + 1
}

// RemindersPath returns the full path to the reminders JSON file.
func RemindersPath() string {
	return filepath.Join(rc.ContextDir(), config.FileReminders)
}
