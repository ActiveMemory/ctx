//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core/store"
	errReminder "github.com/ActiveMemory/ctx/internal/err/reminder"
	"github.com/ActiveMemory/ctx/internal/write/remind"
)

// Many removes one or more reminders by ID. All IDs are resolved
// before any deletion to avoid ordering issues.
//
// Parameters:
//   - cmd: Cobra command for status output
//   - ids: Reminder IDs to dismiss
//
// Returns:
//   - error: Non-nil on missing reminder or write failure
func Many(cmd *cobra.Command, ids []int) error {
	reminders, readErr := store.Read()
	if readErr != nil {
		return readErr
	}

	removeSet := make(map[int]bool, len(ids))
	for _, id := range ids {
		found := false
		for _, r := range reminders {
			if r.ID == id {
				found = true
				break
			}
		}
		if !found {
			return errReminder.NotFound(id)
		}
		removeSet[id] = true
	}

	var remaining []store.Reminder
	for _, r := range reminders {
		if removeSet[r.ID] {
			remind.Dismissed(cmd, r.ID, r.Message)
		} else {
			remaining = append(remaining, r)
		}
	}

	return store.Write(remaining)
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
